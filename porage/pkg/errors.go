package pkg

import "errors"

var (
	// ErrLedgerExisted is the error when the ledger already exists.
	ErrLedgerExisted = errors.New("ledger existed")
	// ErrLedgerNotFound is the error when the ledger is not found.
	ErrLedgerNotFound = errors.New("ledger not found")
	// ErrEntryNotFound is the error when the entry is not found.
	ErrEntryNotFound = errors.New("entry not found")
)
