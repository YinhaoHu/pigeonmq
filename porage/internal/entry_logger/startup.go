package entrylogger

import (
	"os"
	"porage/internal/pkg"
)

var myConfig *pkg.EntryLoggerConfig

// Startup sets the configuration items for entry logger.
func Startup(config *pkg.EntryLoggerConfig) {
	myConfig = config

	err := os.MkdirAll(myConfig.StoragePath, 0755)
	if err != nil {
		panic(err)
	}
}
