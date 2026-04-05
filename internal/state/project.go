package state

import (
	"github.com/jrengmusic/cake/internal"
	"os"
	"path/filepath"
	"time"
)

// Generator represents a CMake generator with metadata
type Generator struct {
	Name  string // CMake project name: "Xcode", "Ninja", "Visual Studio 18 2026", "Visual Studio 17 2022"
	IsIDE bool   // true for Xcode, VS; false for Ninja, Makefiles
}

// BuildInfo represents the state of a build directory
type BuildInfo struct {
	Generator    string   // Generator name
	Path         string   // Full path to build directory
	Exists       bool     // Whether build directory exists
	IsConfigured bool     // Whether CMake has been run (CMakeCache.txt exists)
	Configs      []string // Available configurations (Debug, Release, etc.) - for multi-config generators
}

// ProjectState represents the current state of the CMake project
type ProjectState struct {
	WorkingDirectory  string
	HasCMakeLists     bool
	AvailableProjects []Generator          // Projects detected as available on system
	SelectedProject   string               // Currently selected project (cycled by user)
	Builds            map[string]BuildInfo // Build state by project name (Builds/<Generator>/)
	Configuration     string               // Current configuration: "Debug" or "Release"
	IsPluginProject   bool
	LastRefreshTime   time.Time
	RefreshInterval   time.Duration
	IsConfigured      bool     // Whether CMake has been run (CMakeCache.txt exists)
	Configs           []string // Available configurations (Debug, Release, etc.) - for multi-config generators
}

// NewProjectState creates a new ProjectState instance
func NewProjectState() *ProjectState {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	ps := &ProjectState{
		WorkingDirectory:  cwd,
		HasCMakeLists:     false,
		AvailableProjects: []Generator{},
		SelectedProject:   "",
		Builds:            make(map[string]BuildInfo),
		Configuration:     internal.ConfigDebug,
		IsPluginProject:   false,
		RefreshInterval:   time.Second * 2,
	}

	// Detect available projects on startup
	ps.DetectAvailableProjects()

	return ps
}

// Refresh updates project state from filesystem
func (ps *ProjectState) Refresh() {
	if !ps.ShouldRefresh() {
		return
	}
	ps.ForceRefresh()
}

// ForceRefresh immediately refreshes project state
func (ps *ProjectState) ForceRefresh() {
	cwd, err := os.Getwd()
	if err != nil {
		ps.WorkingDirectory = "."
		ps.HasCMakeLists = false
		return
	}
	ps.WorkingDirectory = cwd

	cmakePath := filepath.Join(cwd, internal.CMakeListsFile)
	_, err = os.Stat(cmakePath)
	ps.HasCMakeLists = (err == nil)

	ps.scanBuildDirectories(cwd)

	// Set default selected project if none selected
	if ps.SelectedProject == "" && len(ps.AvailableProjects) > 0 {
		ps.SelectedProject = ps.AvailableProjects[0].Name
	}

	ps.LastRefreshTime = time.Now()
}

// ShouldRefresh checks if refresh is needed based on interval
func (ps *ProjectState) ShouldRefresh() bool {
	return time.Since(ps.LastRefreshTime) > ps.RefreshInterval
}

// CycleToNextProject advances to the next available project
func (ps *ProjectState) CycleToNextProject() {
	if len(ps.AvailableProjects) == 0 {
		return
	}

	currentIndex := -1
	for i, gen := range ps.AvailableProjects {
		if gen.Name == ps.SelectedProject {
			currentIndex = i
			break
		}
	}

	// Move to next, wrapping around
	if currentIndex >= 0 {
		nextIndex := (currentIndex + 1) % len(ps.AvailableProjects)
		ps.SelectedProject = ps.AvailableProjects[nextIndex].Name
	} else {
		ps.SelectedProject = ps.AvailableProjects[0].Name
	}
}

// CycleToPrevProject advances to the previous available project
func (ps *ProjectState) CycleToPrevProject() {
	if len(ps.AvailableProjects) == 0 {
		return
	}

	currentIndex := -1
	for i, gen := range ps.AvailableProjects {
		if gen.Name == ps.SelectedProject {
			currentIndex = i
			break
		}
	}

	// Move to previous, wrapping around
	if currentIndex >= 0 {
		prevIndex := currentIndex - 1
		if prevIndex < 0 {
			prevIndex = len(ps.AvailableProjects) - 1
		}
		ps.SelectedProject = ps.AvailableProjects[prevIndex].Name
	} else {
		ps.SelectedProject = ps.AvailableProjects[0].Name
	}
}

// CycleConfiguration toggles between Debug and Release
func (ps *ProjectState) CycleConfiguration() {
	if ps.Configuration == internal.ConfigDebug {
		ps.Configuration = internal.ConfigRelease
	} else {
		ps.Configuration = internal.ConfigDebug
	}
}

// SetConfiguration sets the configuration directly (used for restoring from config)
func (ps *ProjectState) SetConfiguration(configuration string) {
	if configuration == internal.ConfigDebug || configuration == internal.ConfigRelease {
		ps.Configuration = configuration
	}
}

// SetSelectedProject sets the selected project directly (used for restoring from config)
func (ps *ProjectState) SetSelectedProject(generator string) {
	// Validate that the project is available
	for _, gen := range ps.AvailableProjects {
		if gen.Name == generator {
			ps.SelectedProject = generator
			return
		}
	}
	// If not found in available projects, don't set it
	// It will fall back to first available on ForceRefresh
}

// GetBuildPath returns the build directory path for the selected project
func (ps *ProjectState) GetBuildPath() string {
	if ps.SelectedProject == "" {
		return ""
	}
	return ps.GetBuildDirectory(ps.SelectedProject)
}

// GetSelectedBuildInfo returns the build info for the selected project
func (ps *ProjectState) GetSelectedBuildInfo() BuildInfo {
	if buildInfo, exists := ps.Builds[ps.SelectedProject]; exists {
		return buildInfo
	}
	return BuildInfo{Generator: ps.SelectedProject, Exists: false}
}

// CanGenerate returns true if we can generate the selected project
func (ps *ProjectState) CanGenerate() bool {
	return ps.SelectedProject != "" && ps.HasCMakeLists
}

// CanBuild returns true if we can build (build directory exists and is configured)
func (ps *ProjectState) CanBuild() bool {
	buildInfo := ps.GetSelectedBuildInfo()
	return buildInfo.Exists && buildInfo.IsConfigured
}

// CanOpenIDE returns true if any project is available (IDE or editor)
func (ps *ProjectState) CanOpenIDE() bool {
	return len(ps.AvailableProjects) > 0
}

// HasBuildsToClean returns true if Builds/ directory exists and has content
// This is used to determine if Clean All should be available
func (ps *ProjectState) HasBuildsToClean() bool {
	buildsDir := filepath.Join(ps.WorkingDirectory, internal.BuildsDirName)

	// Check if Builds directory exists
	info, err := os.Stat(buildsDir)
	if err != nil || !info.IsDir() {
		return false
	}

	// Check if Builds directory has any content
	entries, err := os.ReadDir(buildsDir)
	if err != nil {
		return false
	}

	return len(entries) > 0
}
