package journal

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"porage/internal/pkg"
	"sort"
	"sync"
	"time"
)

const journalFileSuffix = ".journal"

var (
	currentSegmentFile *os.File = nil

	registeredLedgers = newLedgerRegisterationCenter()
)

type ledgerInfo struct {
	lastFlushTime int64 // Unix nano timestamp
}

type ledgerRegisterationCenter struct {
	lock    *sync.RWMutex
	mapping map[uint64]*ledgerInfo
}

func newLedgerRegisterationCenter() *ledgerRegisterationCenter {
	return &ledgerRegisterationCenter{
		lock:    &sync.RWMutex{},
		mapping: make(map[uint64]*ledgerInfo),
	}
}

func (lrc *ledgerRegisterationCenter) register(ledgerID uint64) {
	lrc.lock.Lock()
	defer lrc.lock.Unlock()

	lrc.mapping[ledgerID] = &ledgerInfo{
		lastFlushTime: 0,
	}
}

func (lrc *ledgerRegisterationCenter) deregister(ledgerID uint64) {
	lrc.lock.Lock()
	defer lrc.lock.Unlock()

	delete(lrc.mapping, ledgerID)
}

func (lrc *ledgerRegisterationCenter) updateFlushTime(ledgerID uint64, flushTime int64) {
	lrc.lock.Lock()
	defer lrc.lock.Unlock()

	if ledgerInfo, ok := lrc.mapping[ledgerID]; ok {
		ledgerInfo.lastFlushTime = flushTime
	}
}

func (lrc *ledgerRegisterationCenter) getLedgersMinFlushTime() int64 {
	lrc.lock.RLock()
	defer lrc.lock.RUnlock()

	minFlushTime := int64(0)
	for _, ledgerInfo := range lrc.mapping {
		if minFlushTime == 0 || ledgerInfo.lastFlushTime < minFlushTime {
			minFlushTime = ledgerInfo.lastFlushTime
		}
	}
	return minFlushTime
}

// JournalEntry represents a journal entry.
//
// Export for testing.
type JournalEntry struct {
	// Timestamp sequence id. Valid after Complete is called.
	sequenceID uint64
	// Size of the entry. Valid after Complete is called.
	Size uint64
	// Payload
	Entry *pkg.JournalEntryPayload
	// Payload in binary format. Valid after Complete is called.
	BinPayload []byte
}

func NewJournalEntry(entry *pkg.JournalEntryPayload) *JournalEntry {
	return &JournalEntry{
		Entry: entry,
	}
}

// complete completes the journal entry by setting the sequence ID and the size.
func (je *JournalEntry) complete() {
	je.sequenceID = uint64(time.Now().UnixNano())
	je.Size = je.Entry.Size()
	je.BinPayload = je.Entry.Serialize()
}

// WriteTo writes the journal entry to the given file.
//
// Export for testing.
func (je *JournalEntry) WriteTo(w *os.File) error {
	je.complete()
	// Write the size
	err := binary.Write(w, binary.BigEndian, je.Size)
	if err != nil {
		return err
	}
	// Write the sequence ID
	err = binary.Write(w, binary.BigEndian, je.sequenceID)
	if err != nil {
		return err
	}
	// Write the entry
	_, err = w.Write(je.BinPayload)
	return err
}

// writeEntry writes a journal entry to the current segment file. The entry might be lost if the process crashes
// before Commit is called.
func writeEntry(entry *JournalEntry) error {
	return entry.WriteTo(currentSegmentFile)
}

// commit flushes the current segment file to disk.
func commit() error {
	return currentSegmentFile.Sync()
}

// createNewSegmentIfNeed creates a new segment file if the current segment file exceeds the soft threshold.
// The current segment file is closed without the guarantee to flush the data to disk. The current design idea
// is to call this function after commit.
func createNewSegmentIfNeed() error {
	fileInfo, err := currentSegmentFile.Stat()
	if err != nil {
		return err
	}

	if uint64(fileInfo.Size()) > myConfig.SegmentSoftThreshold {
		// Close the current segment file.
		currentSegmentFile.Close()
		// Create a new segment file
		segmentFileName := makeSegmentFilePath()
		currentSegmentFile, err = os.OpenFile(segmentFileName, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

// makeSegmentFilePath creates a new segment file path with the current time.
func makeSegmentFilePath() string {
	return fmt.Sprintf("%s/%d%s", myConfig.StoragePath, time.Now().UnixNano(), journalFileSuffix)
}

// getSegmentFilePathList returns all the segment files in the storage path. The returned slice is sorted by the file name.
func getSegmentFilePathList() ([]string, error) {
	files, err := os.ReadDir(myConfig.StoragePath)
	if err != nil {
		return nil, err
	}

	var segmentFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == journalFileSuffix {
			segmentFiles = append(segmentFiles, filepath.Join(myConfig.StoragePath, file.Name()))
		}
	}

	sort.Strings(segmentFiles)
	return segmentFiles, nil
}

// StartStorage initializes the journal storage.
func startStorage() {
	// Set the storage path
	err := os.MkdirAll(myConfig.StoragePath, 0755)
	if err != nil {
		pkg.Logger.Fatalf("Failed to create the storage path: %v", err)
	}

	// Open the current segment file
	segmentFileName := makeSegmentFilePath()
	pkg.Logger.Debugf("Current segment file: %s", segmentFileName)
	currentSegmentFile, err = os.OpenFile(segmentFileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		pkg.Logger.Fatalf("Failed to open the current segment file: %v", err)
	}
}
