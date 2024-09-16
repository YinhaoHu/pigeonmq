package ledger

import (
	"fmt"
	"porage/internal/index"
	"porage/internal/journal"
	"porage/internal/pkg"
	"time"
)

var (
	localWorkerControl = pkg.NewLocalWorkerControl()
)

func (l *Ledger) startWorkers() {
	go l.persistenceWorker()
}

func (l *Ledger) closeWorkers() {
	pkg.Logger.Debugf("Closing ledger %d worker", l.ledgerID)
	l.persistenceWorkerDescription.Stop()
}

func (l *Ledger) persistenceWorker() {
	workerName := fmt.Sprintf("ledger-%d-persistence-worker", l.ledgerID)
	workerDescriptionString := fmt.Sprintf("Ledger %d persistence worker", l.ledgerID)
	workerDescription := pkg.NewWorkerDescription(workerDescriptionString)
	l.persistenceWorkerDescription = workerDescription
	localWorkerControl.RegisterWorker(workerName, workerDescription)

	nWrittenEntry := uint64(0)
	shouldFlushInterval := time.Duration(myConfig.EntryLogger.FlushInterval) * time.Second
	shouldFlushIntervalTicker := time.NewTicker(shouldFlushInterval)
	closed := false
	for !closed {
		pkg.Logger.Debugf("Ledger %d worker is running", l.ledgerID)
		shouldFlush := false
		select {
		case entry := <-l.messageBuffer:
			pkg.Logger.Debugf("Ledger %d worker is handling entry: %d", l.ledgerID, entry.EntryID)
			// Write to entry logger
			if err := l.entryLogger.Write(entry); err != nil {
				pkg.Logger.Errorf("Ledger %d failed to write entry: %v", l.ledgerID, err)
				continue
			}
			nWrittenEntry += 1

			// Flush to disk
			if nWrittenEntry >= myConfig.EntryLogger.FlushRate {
				shouldFlush = true
			}
			pkg.Logger.Debugf("Ledger %d worker handled entry: %d", l.ledgerID, entry.EntryID)
		case <-shouldFlushIntervalTicker.C:
			shouldFlushIntervalTicker.Stop()
			shouldFlush = true
			pkg.Logger.Debugf("Ledger %d worker: should flush triggerd by flush interval", l.ledgerID)
		case <-l.persistenceWorkerDescription.StopChannel():
			pkg.Logger.Infof("%s: stopped", workerName)
			l.persistenceWorkerDescription.StopResponseChannel() <- struct{}{}
			localWorkerControl.UnregisterWorker(workerName)
			closed = true
		}

		if shouldFlush {
			if err := l.flush(); err != nil {
				continue
			}
			shouldFlush = false
			nWrittenEntry = 0
			shouldFlushIntervalTicker.Reset(shouldFlushInterval)
		}
	}
}

func (l *Ledger) flush() error {
	flushedEntryMetadata, err := l.entryLogger.Flush()
	if err != nil {
		pkg.Logger.Errorf("Ledger %d failed to flush entry logger: %v", l.ledgerID, err)
		return err
	}

	if len(flushedEntryMetadata) == 0 {
		pkg.Logger.Debugf("Ledger %d no entry to flush", l.ledgerID)
		return nil
	}

	// Write to index
	for _, entryMetadata := range flushedEntryMetadata {
		indexValue := index.IndexValue{
			Offset: entryMetadata.Offset,
			Size:   entryMetadata.Size,
		}
		if err := l.index.Put(entryMetadata.EntryID, &indexValue); err != nil {
			pkg.Logger.Errorf("Ledger %d failed to write index: %v", l.ledgerID, err)
			return err
		}
		l.lastFlushedEntryID.Store(int64(entryMetadata.EntryID))
	}
	journal.UpdateLedgerFlushTime(l.ledgerID)
	return nil
}
