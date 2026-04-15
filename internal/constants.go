package internal

import "runtime/debug"

// Version information (SSOT)
const AppName = "cake" // Application name

// AppVersion is injected at build time via ldflags.
// Falls back to module version (populated by go install module@version),
// then "dev" for local builds.
var AppVersion = getVersion()

func getVersion() string {
	info, ok := debug.ReadBuildInfo()
	if ok && info.Main.Version != "(devel)" && info.Main.Version != "" {
		return info.Main.Version
	}
	return "dev"
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
