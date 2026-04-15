package internal

import "runtime/debug"

// Version information (SSOT)
const AppName = "cake" // Application name

// AppVersion is set at build time via -ldflags "-X github.com/jrengmusic/cake/internal.AppVersion=vX.Y.Z"
// Must be a constant string initializer — function-call initializers
// silently disable -X ldflags (Go issue #64246).
// Defaults to "dev" for untagged local builds.
var AppVersion = "dev"

// init falls back to module version for `go install module@version` builds
// when ldflags were not applied. Leaves AppVersion as "dev" for local builds.
func init() {
	if AppVersion == "dev" {
		info, ok := debug.ReadBuildInfo()
		if ok && info.Main.Version != "(devel)" && info.Main.Version != "" {
			AppVersion = info.Main.Version
		}
	}
}

// Filesystem names (SSOT)
const (
	BuildsDirName  = "Builds"          // Root directory for all build artifacts
	CMakeListsFile = "CMakeLists.txt"  // CMake project definition file
)

// Build configuration names (SSOT)
const (
	ConfigDebug   = "Debug"
	ConfigRelease = "Release"
)

// Auto-scan defaults (SSOT)
const (
	DefaultAutoScanInterval = 10 // minutes
)
