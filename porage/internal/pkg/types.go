package pkg

import (
	"encoding/binary"
	"sync"

	pb "porage/proto"
)

type NotificationChannel = chan Notification
type NotificationTx = chan<- Notification
type NotificationRx = <-chan Notification

type Notification struct {
	Err error
}

// WriteRequest represents the message type from the network goroutine to ledger worker.
type WriteRequest struct {
	Entry          *JournalEntryPayload
	NotificationTx NotificationTx
}

type LedgerEntry struct {
	EntryID int
	Payload []byte
}

// Serialize the ledger entry to a byte slice.
//
// Format: [EntryID][Payload]
func (e *LedgerEntry) Serialize() []byte {
	data := make([]byte, 8+len(e.Payload))
	binary.BigEndian.PutUint64(data[:8], uint64(e.EntryID))
	copy(data[8:], e.Payload)
	return data
}

// Deserialize the ledger entry from a byte slice.
func (le *LedgerEntry) Deserialize(data []byte) {
	le.EntryID = int(binary.BigEndian.Uint64(data[:8]))
	le.Payload = data[8:]
}

// WorkerDescription is the description of a worker(goroutine).
type WorkerDescription struct {
	Description         string
	stopChannel         chan struct{}
	stopResponseChannel chan struct{}
}

func NewWorkerDescription(description string) *WorkerDescription {
	return &WorkerDescription{
		Description:         description,
		stopChannel:         make(chan struct{}),
		stopResponseChannel: make(chan struct{}),
	}
}

// ToPb converts the WorkerDescription to a protobuf message.
func (wd *WorkerDescription) ToPb() *pb.WorkerDescription {
	return &pb.WorkerDescription{
		Description: wd.Description,
	}
}

func (wd *WorkerDescription) StopChannel() <-chan struct{} {
	return wd.stopChannel
}

func (wd *WorkerDescription) StopResponseChannel() chan<- struct{} {
	return wd.stopResponseChannel
}

func (wd *WorkerDescription) Stop() {
	wd.stopChannel <- struct{}{}
	<-wd.stopResponseChannel
}

type LocalWorkerControl struct {
	workers map[string]*WorkerDescription
	rwlock  sync.RWMutex
}

func NewLocalWorkerControl() *LocalWorkerControl {
	return &LocalWorkerControl{
		workers: make(map[string]*WorkerDescription),
	}
}

func (lwc *LocalWorkerControl) RegisterWorker(name string, desc *WorkerDescription) {
	lwc.rwlock.Lock()
	defer lwc.rwlock.Unlock()
	lwc.workers[name] = desc
}

func (lwc *LocalWorkerControl) UnregisterWorker(name string) {
	lwc.rwlock.Lock()
	defer lwc.rwlock.Unlock()
	delete(lwc.workers, name)
}

func (lwc *LocalWorkerControl) GetWorkerDescriptions() map[string]*WorkerDescription {
	lwc.rwlock.RLock()
	defer lwc.rwlock.RUnlock()
	return lwc.workers
}
