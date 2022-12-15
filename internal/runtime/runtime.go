package runtime

import (
	"flag"
	"runtime/debug"
)

// automatically set by goreleaser
var (
	Version = "DEV"
	Commit  string
)

// Flag (-debug)
var isDebug bool

func init() {
	if Version == "DEV" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "(devel)" {
			Version = info.Main.Version
		}
	}

	flag.BoolVar(&isDebug, "debug", false, "set to true to run the provider with support for debuggers")
}

func Debug() bool {
	if !flag.Parsed() {
		flag.Parse()
	}

	return isDebug
}
