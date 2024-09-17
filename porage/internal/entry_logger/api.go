package entrylogger

import (
	"os"
	"porage/internal/pkg"
)

// EntryLogger is the entry logger for a ledger.
type EntryLogger struct {
	ledgerID uint64
	file     *os.File

	entryMetadata []*EntryMetadata
}

type EntryMetadata struct {
	EntryID int
	Offset  int
	Size    int
}

func NewEntryMetadata(entryID int, offset int, size int) *EntryMetadata {
	return &EntryMetadata{
		EntryID: entryID,
		Offset:  offset,
		Size:    size,
	}
}

// NewEntryLogger creates a new entry logger with the given ledgerID.
func NewEntryLogger(ledgerID uint64) (*EntryLogger, error) {
	filePath := makeFilePathByLedgerID(ledgerID)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	entryMetadata := make([]*EntryMetadata, 0)
	entryLogger := &EntryLogger{
		ledgerID:      ledgerID,
		file:          file,
		entryMetadata: entryMetadata,
	}
	return entryLogger, nil
}

// Write writes the entry to the entry logger. It is not guaranteed that the entry is written to the disk.
func (el *EntryLogger) Write(entry *pkg.LedgerEntry) error {
	pkg.Logger.Debugf("Write entry for ledger %d, entry %d", el.ledgerID, entry.EntryID)
	data := entry.Serialize()

	_, err := el.file.Write(data)
	pkg.Logger.Debugf("Entry written for ledger %d, entry %d", el.ledgerID, entry.EntryID)
	if err != nil {
		return err
	}

	fileInfo, err := el.file.Stat()
	if err != nil {
		return err
	}
	offset := fileInfo.Size() - int64(len(data))
	size := len(data)

	entryMetadata := NewEntryMetadata(entry.EntryID, int(offset), size)
	el.entryMetadata = append(el.entryMetadata, entryMetadata)
	return nil
}

// Read reads the entry from the entry logger.
func (el *EntryLogger) Read(offset int, size int) (*pkg.LedgerEntry, error) {
	binLedgerEntry := make([]byte, size)
	_, err := el.file.ReadAt(binLedgerEntry, int64(offset))
	if err != nil {
		pkg.Logger.Errorf("Read entry failed for ledger %d, offset %d, size %d with err %v", el.ledgerID, offset, size, err)
		return nil, err
	}
	entry := &pkg.LedgerEntry{}
	entry.Deserialize(binLedgerEntry)
	return entry, nil
}

// Delete deletes the entry logger.
func (el *EntryLogger) Delete() error {
	if err := os.Remove(el.file.Name()); err != nil {
		return err
	}
	el.file = nil
	return nil
}

// Flush flushes the entry logger to the file system.
//
// This method is not thread safe. Write and Flush should not be called concurrently.
func (el *EntryLogger) Flush() ([]*EntryMetadata, error) {
	if err := el.file.Sync(); err != nil {
		return nil, err
	}
	entryMetadata := el.entryMetadata
	el.entryMetadata = make([]*EntryMetadata, 0)
	return entryMetadata, nil
}

// Truncate truncates the entry logger to the given size.
//
// This function is expected to be called in recovery.
func (el *EntryLogger) Truncate(size int64) error {
	return el.file.Truncate(size)
}

// Close closes the entry logger.
func (el *EntryLogger) Close() error {
	if el.file == nil {
		return nil
	}
	return el.file.Close()
}
