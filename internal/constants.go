package internal

// Version information (SSOT)
const AppName = "cake" // Application name

// AppVersion is set at build time via -ldflags "-X github.com/jrengmusic/cake/internal.AppVersion=vX.Y.Z"
// Defaults to "dev" for untagged local builds
var AppVersion = "dev"

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
