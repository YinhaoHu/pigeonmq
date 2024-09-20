package journal

import (
	"porage/internal/pkg"
)

var myConfig *pkg.JournalConfig

func Startup(config *pkg.JournalConfig) {
	myConfig = config
	startStorage()
	startWorkers()
}

// Stop stops the workers.
func Stop() {
	closeWorkers()
	pkg.Logger.Infof("Journal stopped")
}
