package ledger

import (
	"os"
	"porage/internal/pkg"
)

var myConfig *pkg.Config

func Startup(config *pkg.Config) {
	myConfig = config

	err := os.MkdirAll(myConfig.Ledger.StoragePath, 0755)
	if err != nil {
		pkg.Logger.Fatalf("Failed to create ledger storage directory: %v", err)
	}
}

// Stop stops the workers.
func Stop() {
	// Close the workers
	for _, description := range localWorkerControl.GetWorkerDescriptions() {
		description.Stop()
	}
	pkg.Logger.Infof("Ledgers stopped")
}
