package internal

import (
	"runtime/debug"
)

// automatically set by goreleaser
var (
	Version = "DEV"
	Commit  string
)

func init() {
	if Version == "DEV" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "(devel)" {
			Version = info.Main.Version
		}
	}
}
