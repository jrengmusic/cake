package state

import (
	"github.com/jrengmusic/cake/internal"
	"github.com/jrengmusic/cake/internal/utils"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// DetectAvailableProjects checks which generators are available on the system
func (ps *ProjectState) DetectAvailableProjects() {
	ps.AvailableProjects = []Generator{}

	// Check Xcode (macOS only)
	if runtime.GOOS == "darwin" {
		if ps.checkCommandExists("xcodebuild") {
			ps.AvailableProjects = append(ps.AvailableProjects, Generator{
				Name:  "Xcode",
				IsIDE: true,
			})
		}
	}

	// Check Ninja (cross-platform)
	if ps.checkCommandExists("ninja") {
		ps.AvailableProjects = append(ps.AvailableProjects, Generator{
			Name:  "Ninja",
			IsIDE: false,
		})
	}

	// Check Visual Studio (Windows only)
	if runtime.GOOS == "windows" {
		vsGenerators := ps.checkVSGeneratorAvailable()
		for _, vsGen := range vsGenerators {
			ps.AvailableProjects = append(ps.AvailableProjects, Generator{
				Name:  vsGen,
				IsIDE: true,
			})
		}
	}
}

// checkCommandExists tests if a command is available in PATH
func (ps *ProjectState) checkCommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// checkVSGeneratorAvailable returns generator strings for all installed VS versions
func (ps *ProjectState) checkVSGeneratorAvailable() []string {
	return utils.DetectInstalledVSVersions()
}

// scanBuildDirectories scans for existing build directories
func (ps *ProjectState) scanBuildDirectories(rootDir string) {
	ps.Builds = make(map[string]BuildInfo)

	buildsDir := filepath.Join(rootDir, internal.BuildsDirName)

	entries, err := os.ReadDir(buildsDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dirName := entry.Name()
		buildPath := filepath.Join(buildsDir, dirName)

		buildInfo := BuildInfo{
			Generator: utils.GetGeneratorNameFromDirectory(dirName),
			Path:      buildPath,
			Exists:    true,
		}

		// Check if configured (CMakeCache.txt exists)
		cachePath := filepath.Join(buildPath, "CMakeCache.txt")
		if _, err := os.Stat(cachePath); err == nil {
			buildInfo.IsConfigured = true

			// All projects are multi-config, detect available configurations
			buildInfo.Configs = ps.detectConfigurations(buildPath)
		}

		ps.Builds[utils.GetGeneratorNameFromDirectory(dirName)] = buildInfo
	}
}

// detectConfigurations scans for available build configurations
func (ps *ProjectState) detectConfigurations(buildPath string) []string {
	var configs []string

	entries, err := os.ReadDir(buildPath)
	if err != nil {
		return configs
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// Check if this looks like a configuration directory
			configName := entry.Name()
			if configName == internal.ConfigDebug || configName == internal.ConfigRelease ||
				configName == "RelWithDebInfo" || configName == "MinSizeRel" {
				configs = append(configs, configName)
			}
		}
	}

	return configs
}
