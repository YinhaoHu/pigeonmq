package entrylogger

import (
	"path"
	"strconv"
)

// makeFilePathByLedgerID makes the file path by the ledger ID.
func makeFilePathByLedgerID(ledgerID uint64) string {
	filePath := path.Join(myConfig.StoragePath, "ledger_"+strconv.FormatUint(ledgerID, 10)+".logger")
	return filePath
}
