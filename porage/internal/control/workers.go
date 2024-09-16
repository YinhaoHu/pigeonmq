package control

import (
	"maps"
	"porage/internal/journal"
	"porage/internal/ledger"
	"porage/internal/pkg"
	"sync"
)

type WorkerControl struct {
	workerRepo     map[string]*pkg.WorkerDescription
	workerRepoLock *sync.RWMutex
}

func NewWorkerControl() *WorkerControl {
	return &WorkerControl{
		workerRepo:     make(map[string]*pkg.WorkerDescription),
		workerRepoLock: &sync.RWMutex{},
	}
}

// List returns a list of worker descriptions.
func (wc *WorkerControl) List() map[string]*pkg.WorkerDescription {
	wc.workerRepoLock.RLock()
	defer wc.workerRepoLock.RUnlock()

	journalWorkerDescriptions := journal.GetWorkerDescriptions()
	for workerName, description := range journalWorkerDescriptions {
		wc.workerRepo[workerName] = description
	}

	ledgerWorkerDescriptions := ledger.GetWorkerDescriptions()
	for workerName, description := range ledgerWorkerDescriptions {
		wc.workerRepo[workerName] = description
	}

	return maps.Clone(wc.workerRepo)
}
