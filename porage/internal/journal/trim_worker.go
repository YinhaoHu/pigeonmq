package journal

import (
	"os"
	"path"
	"porage/internal/pkg"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var (
	enableTrimming = atomic.Bool{}
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
			// If in recovering, do not trim.
			if !enableTrimming.Load() {
				pkg.Logger.Warningf("Trimming is disabled")
				continue
			}

			segmentFilePathList, err := getSegmentFilePathList()
			if err != nil {
				pkg.Logger.Errorf("Failed to get segment files: %v", err)
				continue
			}

			minFlushTimeInLedgers := registeredLedgers.getLedgersMinFlushTime()
			for nextSegmentIndex := 1; nextSegmentIndex < len(segmentFilePathList); nextSegmentIndex++ {
				nextSegmentFilePath := segmentFilePathList[nextSegmentIndex]
				nextSegmentCreateTimestamp, err := getSegmentFilePathCreateTimestamp(nextSegmentFilePath)
				if err != nil {
					pkg.Logger.Errorf("Failed to parse segment create time of segment file %s: %v", nextSegmentFilePath, err)
					continue
				}
				if nextSegmentCreateTimestamp < minFlushTimeInLedgers {
					thisSegmentFilePath := segmentFilePathList[nextSegmentIndex-1]
					if err := os.Remove(thisSegmentFilePath); err != nil {
						pkg.Logger.Errorf("Failed to remove segment file %s: %v", nextSegmentFilePath, err)
						continue
					}
					pkg.Logger.Infof("Removed segment file %s", nextSegmentFilePath)
				}
			}
		}
	}
}

func getSegmentFilePathCreateTimestamp(segmentFilePath string) (int64, error) {
	segmentFilename := path.Base(segmentFilePath)
	segmentCreateTimeString := strings.TrimSuffix(segmentFilename, journalFileSuffix)
	segmentCreateTimestamp, err := strconv.ParseInt(segmentCreateTimeString, 10, 64)
	if err != nil {
		return 0, err
	}
	return segmentCreateTimestamp, nil
}
