package control

import (
	"porage/internal/ledger"
	porage "porage/pkg"
	"sync"
)

type LedgerControl struct {
	ledgerRepo     map[uint64]*ledger.Ledger
	ledgerRepoLock *sync.RWMutex
}

func NewLedgerControl() *LedgerControl {
	return &LedgerControl{
		ledgerRepo:     make(map[uint64]*ledger.Ledger),
		ledgerRepoLock: &sync.RWMutex{},
	}
}

// AddLedger adds a ledger to the ledger control.
//
// This function is thread-safe. Expected to be called in recovery.
func (lc *LedgerControl) AddLedger(l *ledger.Ledger) {
	lc.ledgerRepoLock.Lock()
	defer lc.ledgerRepoLock.Unlock()

	lc.ledgerRepo[l.LedgerID()] = l
}

// CreateLedger creates a new ledger with the given ledgerID. If the ledger already exists, it returns error.
//
// This function is thread-safe.
func (lc *LedgerControl) CreateLedger(ledgerID uint64) error {
	lc.ledgerRepoLock.Lock()
	defer lc.ledgerRepoLock.Unlock()

	if _, ok := lc.ledgerRepo[ledgerID]; ok {
		return porage.ErrLedgerExisted
	}

	l, err := ledger.NewLedger(ledgerID)
	if err != nil {
		return err
	}
	lc.ledgerRepo[ledgerID] = l
	return nil
}

// GetLedger returns the ledger with the given ledgerID. If the ledger does not exist, it returns nil.
//
// This function is thread-safe.
func (lc *LedgerControl) GetLedger(ledgerID uint64) *ledger.Ledger {
	lc.ledgerRepoLock.RLock()
	defer lc.ledgerRepoLock.RUnlock()

	l := lc.ledgerRepo[ledgerID]
	return l
}

// RemoveLedger closes the ledger with the given ledgerID and removes it from the ledger control.
// If the ledger does not exist, it returns error.
//
// This function is thread-safe.
func (lc *LedgerControl) RemoveLedger(ledgerID uint64) error {
	lc.ledgerRepoLock.Lock()
	defer lc.ledgerRepoLock.Unlock()

	ledger, ok := lc.ledgerRepo[ledgerID]
	if !ok {
		return porage.ErrLedgerNotFound
	}
	err := ledger.Close()
	if err != nil {
		return err
	}
	delete(lc.ledgerRepo, ledgerID)
	return nil
}

// ListLedgers returns a list of ledgerIDs in the ledger control.
//
// This function is thread-safe.
func (lc *LedgerControl) ListLedgers() []uint64 {
	lc.ledgerRepoLock.RLock()
	defer lc.ledgerRepoLock.RUnlock()

	ledgerIDs := make([]uint64, 0, len(lc.ledgerRepo))
	for ledgerID := range lc.ledgerRepo {
		ledgerIDs = append(ledgerIDs, ledgerID)
	}
	return ledgerIDs
}
