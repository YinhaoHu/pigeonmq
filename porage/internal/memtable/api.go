package memtable

import (
	"porage/internal/pkg"
	"sync"
)

type MemTable struct {
	ledgerID uint64

	entryContainer     map[int]*pkg.LedgerEntry
	entryContainerLock sync.RWMutex
	minEntryInMem      int
}

func NewMemTable(ledgerID uint64) (*MemTable, error) {
	return &MemTable{
		ledgerID:           ledgerID,
		entryContainer:     make(map[int]*pkg.LedgerEntry),
		entryContainerLock: sync.RWMutex{},
	}, nil
}

// Get the entry with entryID from the memtable. If the entry does not exist, nil will be returned.
//
// An entry which is not in memtable is possibly in the entry logger.
func (m *MemTable) Get(entryID int) (*pkg.LedgerEntry, error) {
	m.entryContainerLock.RLock()
	defer m.entryContainerLock.RUnlock()

	entry := m.entryContainer[entryID]
	return entry, nil
}

// Put puts the entry into the memtable.
func (m *MemTable) Put(entry *pkg.LedgerEntry) {
	m.entryContainerLock.Lock()
	defer m.entryContainerLock.Unlock()

	m.entryContainer[entry.EntryID] = entry
}

// TrimUntil trims the memtable until the entryID.
func (m *MemTable) TrimUntil(entryID int) {
	pkg.Logger.Infof("Begin to trim memtable of ledger %d", m.ledgerID)
	m.entryContainerLock.Lock()
	defer m.entryContainerLock.Unlock()
	for m.minEntryInMem <= entryID && len(m.entryContainer) >= myConfig.TrimThreshold {
		delete(m.entryContainer, m.minEntryInMem)
		m.minEntryInMem++
	}
	pkg.Logger.Infof("Accomplished to trim memtable of ledger %d", m.ledgerID)
}

// MeetTrimThreshold returns true if the memtable meets the trim threshold.
func (m *MemTable) MeetTrimThreshold() bool {
	return len(m.entryContainer) > myConfig.TrimThreshold
}
