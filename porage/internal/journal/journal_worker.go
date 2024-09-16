package journal

import (
	"porage/internal/pkg"
	"time"
)

var (
	messageBuffer chan *pkg.WriteRequest
)

// journal_worker is a goroutine that processes messages from the message buffer.
// It receives messages from that buffer and dispatch them to the corresponding ledger channel.
func journal_worker() {
	workerName := "journal_worker"
	workerDescriptionString := "Write the journal entries and group commit them"
	workerDescription := pkg.NewWorkerDescription(workerDescriptionString)
	localWorkerControl.RegisterWorker(workerName, workerDescription)

	notificationChannelArray := make([]pkg.NotificationTx, 0, myConfig.GroupCommitThreasold)
	groupCommitInterval := time.Duration(myConfig.GroupCommitInterval) * time.Millisecond
	groupCommitIntervalTicker := time.NewTicker(groupCommitInterval)

	shouldCommit := false
	for {
		select {
		case message := <-messageBuffer:
			pkg.Logger.Debugf("Worker receives a message: %v", message)
			// Write the entry to the journal storage
			journalEntry := NewJournalEntry(message.Entry)
			if err := writeEntry(journalEntry); err != nil {
				message.NotificationTx <- pkg.Notification{
					Err: err,
				}
				continue
			}
			notificationChannelArray = append(notificationChannelArray, message.NotificationTx)
			if uint64(len(notificationChannelArray)) >= myConfig.GroupCommitThreasold {
				shouldCommit = true
				pkg.Logger.Debugf("Worker: should group commit triggerd by GroupCommitThreasold")
			}
		case <-groupCommitIntervalTicker.C:
			groupCommitIntervalTicker.Stop()
			shouldCommit = true
			pkg.Logger.Debugf("Worker: should group commit triggerd by groupCommitInterval")
		case <-workerDescription.StopChannel():
			pkg.Logger.Infof("%s: stopped", workerName)
			localWorkerControl.UnregisterWorker(workerName)
			workerDescription.StopResponseChannel() <- struct{}{}
			return
		}

		if shouldCommit {
			pkg.Logger.Debugf("Worker: group commit")
			commit()
			// Notify and clear the notification array
			for _, notificationChannel := range notificationChannelArray {
				notificationChannel <- pkg.Notification{
					Err: nil,
				}
			}
			notificationChannelArray = notificationChannelArray[:0]
			createNewSegmentIfNeed()
			shouldCommit = false
			groupCommitIntervalTicker.Reset(groupCommitInterval)
		}
	}
}
