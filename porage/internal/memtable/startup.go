package memtable

import (
	"porage/internal/pkg"
)

var myConfig *pkg.MemtableConfig

// Startup initializes the memtable.
func Startup(config *pkg.MemtableConfig) {
	myConfig = config
}
