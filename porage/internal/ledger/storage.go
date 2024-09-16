package ledger

import (
	"fmt"
	"os"
	"path"
)

// persist persists the ledger to the file system.
func (l *Ledger) persistInFileSystem() error {
	// Create a file for the ledger.
	filePath := l.makeLedgerFilePath()
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	return file.Sync()
}

// removePersistenceInFileSystem removes the persistence of the ledger in the file system.
func (l *Ledger) removePersistenceInFileSystem() error {
	filePath := l.makeLedgerFilePath()
	return os.Remove(filePath)
}

// getPersistentLedgerIDList returns the IDs of the ledgers that are persisted in the file system.
func getPersistentLedgerIDList() ([]uint64, error) {
	dirPath := myConfig.Ledger.StoragePath

	ledgerIDList := make([]uint64, 0)
	ledgerFileList, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, file := range ledgerFileList {
		filename := path.Base(file.Name())
		ledgerID := parseLedgerIDFromFileName(filename)
		ledgerIDList = append(ledgerIDList, ledgerID)
	}
	return ledgerIDList, err
}

func (l *Ledger) makeLedgerFilePath() string {
	ledgerName := fmt.Sprintf("ledger_%d", l.ledgerID)
	return path.Join(myConfig.Ledger.StoragePath, ledgerName)
}

// parseLedgerIDFromFileName parses the ledger ID from the file name(No slash).
func parseLedgerIDFromFileName(filename string) uint64 {
	var ledgerID uint64
	fmt.Sscanf(filename, "ledger_%d", &ledgerID)
	return ledgerID
}
