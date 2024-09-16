package index

import (
	"path"
	"strconv"
)

func makeStoragePathByLedgerID(ledgerID uint64) string {
	return path.Join(myConfig.StoragePath, "ledger_"+strconv.FormatUint(ledgerID, 10))
}
