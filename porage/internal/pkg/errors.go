package pkg

import "errors"

var (
	ErrLedgerNotFound      = errors.New("ledger not found")
	ErrLedgerAlreadyExists = errors.New("ledger already exists")
	ErrBufferBusy          = errors.New("buffer busy")
)
