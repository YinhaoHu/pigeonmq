package journal

import (
	"os"
	"path"
	"porage/internal/pkg"
	"strconv"
	"strings"
	"time"
)

// trim_worker is the worker that trims the journal entries.
func trim_worker() {
	workerName := "trim_worker"
	workerDescription := pkg.NewWorkerDescription("Trim the journal entries")
	localWorkerControl.RegisterWorker(workerName, workerDescription)

	for {
		select {
		case <-workerDescription.StopChannel():
			pkg.Logger.Infof("%s: stopped", workerName)
			localWorkerControl.UnregisterWorker(workerName)
			workerDescription.StopResponseChannel() <- struct{}{}
			return
		case <-time.After(time.Duration(myConfig.TrimInterval) * time.Second):
			segmentFilePathList, err := getSegmentFilePathList()
			if err != nil {
				pkg.Logger.Errorf("Failed to get segment files: %v", err)
				continue
			}
			if path.Base(segmentFilePathList[len(segmentFilePathList)-1]) == path.Base(currentSegmentFile.Name()) {
				// Do not trim the current segment file.
				segmentFilePathList = segmentFilePathList[:len(segmentFilePathList)-1]
			}

			minFlushTimeInLedgers := registeredLedgers.getLedgersMinFlushTime()
			for _, segmentFilePath := range segmentFilePathList {
				segmentFilename := path.Base(segmentFilePath)
				segmentCreateTimeString := strings.TrimSuffix(segmentFilename, journalFileSuffix)
				segmentCreateTimestamp, err := strconv.ParseInt(segmentCreateTimeString, 10, 64)
				if err != nil {
					pkg.Logger.Errorf("Failed to parse segment create time of segment file %s: %v", segmentFilename, err)
					continue
				}
				if segmentCreateTimestamp < minFlushTimeInLedgers {
					if err := os.Remove(segmentFilePath); err != nil {
						pkg.Logger.Errorf("Failed to remove segment file %s: %v", segmentFilePath, err)
						continue
					}
					pkg.Logger.Infof("Removed segment file %s", segmentFilePath)
				}
			}
		}

	}
}
