package version

import (
	"strings"
)

var (
	Version = "0.0.0"
	Commit  = "unknown"
	Date    = "unknown"
)

func init() {
	if !strings.HasPrefix(Version, "v") {
		Version = "v" + Version
	}
}
