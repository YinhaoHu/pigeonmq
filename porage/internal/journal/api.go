package journal

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"porage/internal/pkg"
	"time"
)

// AppendJournal appends a journal entry to the journal storage. The journal entry is guaranteed to
// be written to the storage only if the notification channel notifies.
func AppendJournal(entry *pkg.JournalEntryPayload) (pkg.NotificationRx, error) {
	if uint64(len(messageBuffer)) > myConfig.MessageBufferBusyThreshold {
		return nil, pkg.ErrBufferBusy
	}
	notificationChannel := make(pkg.NotificationChannel, 1)
	messageBuffer <- &pkg.WriteRequest{
		Entry:          entry,
		NotificationTx: notificationChannel,
	}
	return notificationChannel, nil
}

func RegisterLedger(ledgerID uint64) {
	registeredLedgers.register(ledgerID)
}

func DeregisterLedger(ledgerID uint64) {
	registeredLedgers.deregister(ledgerID)
}

func UpdateLedgerFlushTime(ledgerID uint64) {
	registeredLedgers.updateFlushTime(ledgerID, time.Now().UnixNano())
}

// ReadJournal reads the all the journal entries in the first segmenf file before the current segment file.
// If there are more segment files before the current segment file, `nextSegmentIdx` will be negative. Start
// Use the returned `nextSegmentIdx` to read the next segment file. Start with `nextSegmentIdx` as 0.
//
// Expected to be called after the local storage is initialized for recovery.
func ReadJournal(segmentIdx int) (entries []*JournalEntry, nextSegmentIdx int, err error) {
	segmentFiles, _ := os.ReadDir(myConfig.StoragePath)
	// We do not read the current segment file.
	if len(segmentFiles) == 1 {
		return nil, -1, nil
	}
	pkg.Logger.Debugf("Reading journal entries from %d segment files.", len(segmentFiles)-1)

	segmentFile := segmentFiles[segmentIdx]
	segmentPath := filepath.Join(myConfig.StoragePath, segmentFile.Name())
	file, err := os.Open(segmentPath)
	if err != nil {
		return nil, -1, err
	}
	defer file.Close()

	entries = make([]*JournalEntry, 0)
	for {
		entry := &JournalEntry{}
		err = binary.Read(file, binary.BigEndian, &entry.Size)
		if err != nil {
			break
		}
		err = binary.Read(file, binary.BigEndian, &entry.sequenceID)
		if err != nil {
			break
		}
		entryData := make([]byte, entry.Size)
		_, err = file.Read(entryData)
		if err != nil {
			break
		}
		entry.Entry = pkg.DeserializeJournalEntry(entryData)
		entries = append(entries, entry)
	}
	if err.Error() != "EOF" {
		return nil, -1, err
	}

	// Do not read the current segment file.
	nextSegmentIdx = segmentIdx + 1
	if nextSegmentIdx == len(segmentFiles)-1 {
		return entries, -1, nil
	}

	return entries, nextSegmentIdx, nil
}

// EnableTrimWorker enables the trim worker.
func EnableTrimWorker() {
	enableTrimming.Store(true)
}

// GetWorkerDescriptions returns the workers description.
func GetWorkerDescriptions() map[string]*pkg.WorkerDescription {
	return localWorkerControl.GetWorkerDescriptions()
}
