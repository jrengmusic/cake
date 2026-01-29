package state

import (
	"cake/internal/utils"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Generator represents a CMake generator with metadata
type Generator struct {
	Name  string // CMake generator name: "Xcode", "Ninja", "Visual Studio 18 2026", "Visual Studio 17 2022"
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
	WorkingDirectory    string
	HasCMakeLists       bool
	AvailableGenerators []Generator          // Generators detected as available on system
	SelectedGenerator   string               // Currently selected generator (cycled by user)
	Builds              map[string]BuildInfo // Build state by generator name (Builds/<Generator>/)
	Configuration       string               // Current configuration: "Debug" or "Release"
	IsPluginProject     bool
	LastRefreshTime     time.Time
	RefreshInterval     time.Duration
	IsConfigured        bool     // Whether CMake has been run (CMakeCache.txt exists)
	Configs             []string // Available configurations (Debug, Release, etc.) - for multi-config generators
}

// NewProjectState creates a new ProjectState instance
func NewProjectState() *ProjectState {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	ps := &ProjectState{
		WorkingDirectory:    cwd,
		HasCMakeLists:       false,
		AvailableGenerators: []Generator{},
		SelectedGenerator:   "",
		Builds:              make(map[string]BuildInfo),
		Configuration:       "Debug",
		IsPluginProject:     false,
		RefreshInterval:     time.Second * 2,
	}

	// Detect available generators on startup
	ps.DetectAvailableGenerators()

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

	cmakePath := filepath.Join(cwd, "CMakeLists.txt")
	_, err = os.Stat(cmakePath)
	ps.HasCMakeLists = (err == nil)

	ps.scanBuildDirectories(cwd)

	// Set default selected generator if none selected
	if ps.SelectedGenerator == "" && len(ps.AvailableGenerators) > 0 {
		ps.SelectedGenerator = ps.AvailableGenerators[0].Name
	}

	ps.LastRefreshTime = time.Now()
}

// ShouldRefresh checks if refresh is needed based on interval
func (ps *ProjectState) ShouldRefresh() bool {
	if time.Since(ps.LastRefreshTime) > ps.RefreshInterval {
		return true
	}
	return false
}

// DetectAvailableGenerators checks which generators are available on the system
func (ps *ProjectState) DetectAvailableGenerators() {
	ps.AvailableGenerators = []Generator{}

	// Check Xcode (macOS only)
	if runtime.GOOS == "darwin" {
		if ps.checkCommandExists("xcodebuild") {
			ps.AvailableGenerators = append(ps.AvailableGenerators, Generator{
				Name:  "Xcode",
				IsIDE: true,
			})
		}
	}

	// Check Ninja (cross-platform)
	if ps.checkCommandExists("ninja") {
		ps.AvailableGenerators = append(ps.AvailableGenerators, Generator{
			Name:  "Ninja",
			IsIDE: false,
		})
	}

	// Check Visual Studio (Windows only)
	if runtime.GOOS == "windows" {
		// Check for Visual Studio via vswhere or cmake generator availability
		if ps.checkVSGeneratorAvailable() {
			ps.AvailableGenerators = append(ps.AvailableGenerators, Generator{
				Name:  "Visual Studio 18 2026",
				IsIDE: true,
			})
			ps.AvailableGenerators = append(ps.AvailableGenerators, Generator{
				Name:  "Visual Studio 17 2022",
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

// checkVSGeneratorAvailable checks if Visual Studio generators are available
func (ps *ProjectState) checkVSGeneratorAvailable() bool {
	// Try to run cmake with --help to list generators
	cmd := exec.Command("cmake", "--help")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "Visual Studio")
}

// CycleToNextGenerator advances to the next available generator
func (ps *ProjectState) CycleToNextGenerator() {
	if len(ps.AvailableGenerators) == 0 {
		return
	}

	currentIndex := -1
	for i, gen := range ps.AvailableGenerators {
		if gen.Name == ps.SelectedGenerator {
			currentIndex = i
			break
		}
	}

	// Move to next, wrapping around
	if currentIndex >= 0 {
		nextIndex := (currentIndex + 1) % len(ps.AvailableGenerators)
		ps.SelectedGenerator = ps.AvailableGenerators[nextIndex].Name
	} else {
		ps.SelectedGenerator = ps.AvailableGenerators[0].Name
	}
}

// CycleToPrevGenerator advances to the previous available generator
func (ps *ProjectState) CycleToPrevGenerator() {
	if len(ps.AvailableGenerators) == 0 {
		return
	}

	currentIndex := -1
	for i, gen := range ps.AvailableGenerators {
		if gen.Name == ps.SelectedGenerator {
			currentIndex = i
			break
		}
	}

	// Move to previous, wrapping around
	if currentIndex >= 0 {
		prevIndex := currentIndex - 1
		if prevIndex < 0 {
			prevIndex = len(ps.AvailableGenerators) - 1
		}
		ps.SelectedGenerator = ps.AvailableGenerators[prevIndex].Name
	} else {
		ps.SelectedGenerator = ps.AvailableGenerators[0].Name
	}
}

// CycleConfiguration toggles between Debug and Release
func (ps *ProjectState) CycleConfiguration() {
	if ps.Configuration == "Debug" {
		ps.Configuration = "Release"
	} else {
		ps.Configuration = "Debug"
	}
}

// GetBuildPath returns the build directory path for the selected generator
func (ps *ProjectState) GetBuildPath() string {
	if ps.SelectedGenerator == "" {
		return ""
	}

	gen := ps.SelectedGenerator
	// All projects are multi-config: Builds/<dir>/
	return filepath.Join(ps.WorkingDirectory, "Builds", utils.GetDirectoryName(gen))
}

// GetBuildDirectory returns the build directory path for given generator and config
func (ps *ProjectState) GetBuildDirectory(generatorName string, config string) string {
	// All projects are multi-config: Builds/<dir>/
	return filepath.Join(ps.WorkingDirectory, "Builds", utils.GetDirectoryName(generatorName))
}

// scanBuildDirectories scans for existing build directories
func (ps *ProjectState) scanBuildDirectories(rootDir string) {
	ps.Builds = make(map[string]BuildInfo)

	buildsDir := filepath.Join(rootDir, "Builds")

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
			if configName == "Debug" || configName == "Release" ||
				configName == "RelWithDebInfo" || configName == "MinSizeRel" {
				configs = append(configs, configName)
			}
		}
	}

	return configs
}

// GetSelectedBuildInfo returns the build info for the selected generator
func (ps *ProjectState) GetSelectedBuildInfo() BuildInfo {
	if buildInfo, exists := ps.Builds[ps.SelectedGenerator]; exists {
		return buildInfo
	}
	return BuildInfo{Generator: ps.SelectedGenerator, Exists: false}
}

// CanGenerate returns true if we can generate the selected generator
func (ps *ProjectState) CanGenerate() bool {
	return ps.SelectedGenerator != "" && ps.HasCMakeLists
}

// CanBuild returns true if we can build (build directory exists and is configured)
func (ps *ProjectState) CanBuild() bool {
	buildInfo := ps.GetSelectedBuildInfo()
	return buildInfo.Exists && buildInfo.IsConfigured
}

// CanOpenIDE returns true if the selected generator is an IDE generator
func (ps *ProjectState) CanOpenIDE() bool {
	for _, gen := range ps.AvailableGenerators {
		if gen.Name == ps.SelectedGenerator && gen.IsIDE {
			return true
		}
	}
	return false
}

// CanOpenEditor returns true if we can open editor (CLI build exists)
func (ps *ProjectState) CanOpenEditor() bool {
	buildInfo := ps.GetSelectedBuildInfo()
	// Can open editor if build exists (even if not configured)
	return buildInfo.Exists
}

// GetProjectLabel returns a display-friendly generator name
func (ps *ProjectState) GetProjectLabel() string {
	name := ps.SelectedGenerator

	// Truncate long names for display
	switch name {
	case "Visual Studio 18 2026":
		return "VS 2026"
	case "Visual Studio 17 2022":
		return "VS 2022"
	default:
		return name
	}
}

// GetProjectName extracts project name using layered file parsing
// Priority: KANJUT Parameters.xml > set(PROJECT_NAME) > project() literal > directory name
func (ps *ProjectState) GetProjectName() string {
	// Layer 1: KANJUT Parameters.xml
	if name := ps.extractFromParametersXML(); name != "" {
		return name
	}

	// Layer 2: set(PROJECT_NAME "...") in CMakeLists.txt
	if name := ps.extractFromSetProjectName(); name != "" {
		return name
	}

	// Layer 3: project(LiteralName) in CMakeLists.txt
	if name := ps.extractFromProjectCall(); name != "" {
		return name
	}

	// Layer 4: Fallback to project working directory
	return filepath.Base(ps.WorkingDirectory)
}

// extractFromParametersXML parses KANJUT's Parameters.xml for PRODUCT name attribute
func (ps *ProjectState) extractFromParametersXML() string {
	paramsPath := filepath.Join(ps.WorkingDirectory, "Source", "layout", "Parameters.xml")

	data, err := os.ReadFile(paramsPath)
	if err != nil {
		return ""
	}

	content := string(data)

	// Find <PRODUCT ... name="..." ...>
	productIdx := strings.Index(content, "<PRODUCT")
	if productIdx == -1 {
		return ""
	}

	// Find closing > of PRODUCT tag
	closeIdx := strings.Index(content[productIdx:], ">")
	if closeIdx == -1 {
		return ""
	}

	productTag := content[productIdx : productIdx+closeIdx]

	// Extract name="value"
	nameIdx := strings.Index(productTag, `name="`)
	if nameIdx == -1 {
		return ""
	}

	valueStart := nameIdx + 6 // len(`name="`)
	valueEnd := strings.Index(productTag[valueStart:], `"`)
	if valueEnd == -1 {
		return ""
	}

	return productTag[valueStart : valueStart+valueEnd]
}

// extractFromSetProjectName parses set(PROJECT_NAME "...") from CMakeLists.txt
func (ps *ProjectState) extractFromSetProjectName() string {
	cmakePath := filepath.Join(ps.WorkingDirectory, "CMakeLists.txt")

	data, err := os.ReadFile(cmakePath)
	if err != nil {
		return ""
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Match: set(PROJECT_NAME "value") or set(PROJECT_NAME "value" ...)
		if strings.HasPrefix(strings.ToLower(trimmed), "set(project_name") {
			// Find first quoted string
			quoteStart := strings.Index(trimmed, `"`)
			if quoteStart == -1 {
				continue
			}
			quoteEnd := strings.Index(trimmed[quoteStart+1:], `"`)
			if quoteEnd == -1 {
				continue
			}
			return trimmed[quoteStart+1 : quoteStart+1+quoteEnd]
		}
	}

	return ""
}

// extractFromProjectCall parses project(Name) from CMakeLists.txt (literal names only)
func (ps *ProjectState) extractFromProjectCall() string {
	cmakePath := filepath.Join(ps.WorkingDirectory, "CMakeLists.txt")

	data, err := os.ReadFile(cmakePath)
	if err != nil {
		return ""
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		lowerTrimmed := strings.ToLower(trimmed)

		// Match: project( at start of line
		if !strings.HasPrefix(lowerTrimmed, "project(") {
			continue
		}

		// Extract content after "project("
		parenStart := strings.Index(trimmed, "(")
		if parenStart == -1 {
			continue
		}

		rest := trimmed[parenStart+1:]

		// Skip if contains variable reference
		if strings.Contains(rest, "$") {
			return ""
		}

		// Handle quoted: project("Name" ...)
		if strings.HasPrefix(strings.TrimSpace(rest), `"`) {
			quoteStart := strings.Index(rest, `"`)
			quoteEnd := strings.Index(rest[quoteStart+1:], `"`)
			if quoteEnd != -1 {
				return rest[quoteStart+1 : quoteStart+1+quoteEnd]
			}
		}

		// Handle unquoted: project(Name) or project(Name VERSION ...)
		// Extract first token (alphanumeric + underscore + hyphen)
		var name strings.Builder
		for _, ch := range strings.TrimSpace(rest) {
			if ch == ' ' || ch == '\t' || ch == ')' || ch == '"' {
				break
			}
			name.WriteRune(ch)
		}

		result := name.String()
		if result != "" && !strings.Contains(result, "$") {
			return result
		}
	}

	return ""
}
