package index

import (
	"porage/internal/pkg"
)

// Configuration items.
var myConfig *pkg.IndexFileConfig

func Startup(config *pkg.IndexFileConfig) {
	myConfig = config
}
