package journal

import (
	"porage/internal/pkg"
)

var (
	localWorkerControl = pkg.NewLocalWorkerControl()
)

func startWorkers() {
	messageBuffer = make(chan *pkg.WriteRequest, myConfig.MessageBufferSize)
	go journal_worker()
	go trim_worker()
}

func closeWorkers() {
	descriptions := localWorkerControl.GetWorkerDescriptions()
	for _, description := range descriptions {
		description.Stop()
	}
}
