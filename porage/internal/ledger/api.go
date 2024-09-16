package ledger

import (
	entrylogger "porage/internal/entry_logger"
	"porage/internal/index"
	"porage/internal/journal"
	"porage/internal/memtable"
	"porage/internal/pkg"
	"sync"
	"sync/atomic"
)

type Ledger struct {
	ledgerID         uint64
	entryCounter     int
	entryCounterLock *sync.Mutex

	entryLogger *entrylogger.EntryLogger
	index       *index.Index
	memtable    *memtable.MemTable

	lastFlushedEntryID           *atomic.Int64
	messageBuffer                chan *pkg.LedgerEntry
	persistenceWorkerDescription *pkg.WorkerDescription
}

func NewLedger(ledgerID uint64) (*Ledger, error) {
	entryLogger, err := entrylogger.NewEntryLogger(ledgerID)
	if err != nil {
		return nil, err
	}
	index, err := index.NewIndex(ledgerID)
	if err != nil {
		return nil, err
	}
	memtable, err := memtable.NewMemTable(ledgerID)
	if err != nil {
		return nil, err
	}
	lastFlushedEntryID := &atomic.Int64{}
	lastFlushedEntryID.Store(-1)

	messageBuffer := make(chan *pkg.LedgerEntry, myConfig.EntryLogger.MessageBufferSize)
	ledger := &Ledger{
		ledgerID:           ledgerID,
		entryCounterLock:   &sync.Mutex{},
		entryLogger:        entryLogger,
		index:              index,
		memtable:           memtable,
		lastFlushedEntryID: lastFlushedEntryID,
		messageBuffer:      messageBuffer,
	}
	err = ledger.persistInFileSystem()
	if err != nil {
		return nil, err
	}

	ledger.startWorkers()
	journal.RegisterLedger(ledgerID)
	return ledger, nil
}

// LedgerID returns the ledger ID.
func (l *Ledger) LedgerID() uint64 {
	return l.ledgerID
}

// PutEntry puts the entry with entryID and payload into the ledger.
func (l *Ledger) PutEntry(payload []byte) (int, error) {
	l.entryCounterLock.Lock()
	entryID := l.entryCounter
	l.entryCounter++
	pkg.Logger.Debugf("PutEntry: entryID=%d, payload=%s", entryID, string(payload))
	// Write journal
	pkg.Logger.Debugf("PutEntry: entryID=%d, payload=%s, put into journal", entryID, string(payload))
	journalEntryPayload := pkg.JournalEntryPayload{
		LedgerID: l.ledgerID,
		EntryID:  entryID,
		Payload:  payload,
	}
	notificationRx, err := journal.AppendJournal(&journalEntryPayload)
	if err != nil {
		l.entryCounterLock.Unlock()
		return -1, err
	}

	// Write to memtable
	pkg.Logger.Debugf("PutEntry: entryID=%d, payload=%s, put into memtable", entryID, string(payload))
	ledgerEntry := &pkg.LedgerEntry{
		EntryID: entryID,
		Payload: payload,
	}
	l.memtable.Put(ledgerEntry)

	// Async flush
	pkg.Logger.Debugf("PutEntry: entryID=%d, payload=%s, put into message buffer", entryID, string(payload))
	l.messageBuffer <- ledgerEntry
	l.entryCounterLock.Unlock()

	notification := <-notificationRx
	if l.memtable.MeetTrimThreshold() {
		l.memtable.TrimUntil(int(l.lastFlushedEntryID.Load()))
	}

	pkg.Logger.Debugf("PutEntry: entryID=%d, payload=%s, done", entryID, string(payload))
	return entryID, notification.Err
}

// GetEntry returns the entry with entryID in the ledger. If the entry does not exist, return nil.
func (l *Ledger) GetEntry(entryID int) (*pkg.LedgerEntry, error) {
	// Get from memtable first
	entry, err := l.memtable.Get(entryID)
	if err != nil {
		return nil, err
	}
	if entry != nil {
		return entry, nil
	}

	// Get from entry logger
	index, err := l.index.Get(entryID)
	if err != nil {
		return nil, err
	}
	if index == nil {
		return nil, nil
	}
	return l.entryLogger.Read(index.Offset, index.Size)
}

// Length returns the number of entries in the ledger.
func (l *Ledger) Length() (int, error) {
	lastEntryID, indexValue, err := l.index.LastItem()
	pkg.Logger.Debugf("Length: ledgerID=%d, lastEntryID=%d", l.ledgerID, lastEntryID)
	if err != nil {
		return -1, err
	}
	if indexValue == nil {
		return 0, nil
	}
	return lastEntryID + 1, nil
}

// Close closes the ledger.
func (l *Ledger) Close() error {
	l.closeWorkers()
	if err := l.removePersistenceInFileSystem(); err != nil {
		return err
	}

	// Delete entry logger
	if err := l.entryLogger.Delete(); err != nil {
		return err
	}
	// Delete index
	if err := l.index.Delete(); err != nil {
		return err
	}
	// If persistence file is removed, the recovery will not involve this ledger automatically.
	journal.DeregisterLedger(l.ledgerID)

	return nil
}

// PrepareRecovery prepares the ledger for recovery. Return the fromEntryID for recovery.
func (l *Ledger) PrepareRecovery() (int, error) {
	fromEntryID, lastIndexValue, err := l.index.LastItem()
	pkg.Logger.Infof("PrepareRecovery: ledgerID=%d, fromEntryID=%d", l.ledgerID, fromEntryID)
	if err != nil {
		return -1, err
	}
	if lastIndexValue == nil {
		return 0, nil
	}

	l.entryCounter = fromEntryID
	validSize := lastIndexValue.Offset + lastIndexValue.Size
	if err := l.entryLogger.Truncate(int64(validSize)); err != nil {
		return -1, err
	}
	return fromEntryID, nil
}

// PutEntryOnRecovery puts the entry with payload into the ledger during recovery.
//
// Expected to be called only in recovery.
func (l *Ledger) PutEntryOnRecovery(payload []byte) error {
	entryID := l.entryCounter
	l.entryCounter++
	pkg.Logger.Debugf("PutEntryOnRecovery: entryID=%d, payload=%s", entryID, string(payload))

	// Write to memtable
	pkg.Logger.Debugf("PutEntryOnRecovery: entryID=%d, payload=%s, put into memtable", entryID, string(payload))
	ledgerEntry := &pkg.LedgerEntry{
		EntryID: entryID,
		Payload: payload,
	}
	l.memtable.Put(ledgerEntry)

	// Async flush
	pkg.Logger.Debugf("PutEntryOnRecovery: entryID=%d, payload=%s, put into message buffer", entryID, string(payload))
	l.messageBuffer <- ledgerEntry

	if l.memtable.MeetTrimThreshold() {
		l.memtable.TrimUntil(int(l.lastFlushedEntryID.Load()))
	}

	pkg.Logger.Debugf("PutEntryOnRecovery: entryID=%d, payload=%s, done", entryID, string(payload))
	return nil
}

// GetPersistentLedgerIDList returns the IDs of the ledgers that are persisted in the file system.
//
// This function is expected to be called in recovery.
func GetPersistentLedgerIDList() ([]uint64, error) {
	return getPersistentLedgerIDList()
}

// GetWorkerDescriptions returns the descriptions of the workers in the ledger.
func GetWorkerDescriptions() map[string]*pkg.WorkerDescription {
	return localWorkerControl.GetWorkerDescriptions()
}
