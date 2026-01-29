package ops

import (
	"cake/internal/ui"
	"cake/internal/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type OpenResult struct {
	Success bool
	Error   string
}

// ExecuteOpenIDE opens the IDE project for the given build directory
func ExecuteOpenIDE(generator, config, projectRoot string, outputCallback func(string, ui.OutputLineType)) OpenResult {
	buildDir := filepath.Join(projectRoot, "Builds", utils.GetDirectoryName(generator))

	var projectFile string

	switch generator {
	case "Xcode":
		projectFile = findXcodeProject(buildDir)
		if projectFile == "" {
			return OpenResult{Success: false, Error: "No Xcode project found"}
		}
		cmd := exec.Command("open", projectFile)
		outputCallback("Opening Xcode project: "+projectFile, ui.TypeInfo)
		if err := cmd.Run(); err != nil {
			return OpenResult{Success: false, Error: err.Error()}
		}
		return OpenResult{Success: true}

	case "Visual Studio 18 2026", "Visual Studio 17 2022":
		projectFile = findVisualStudioSolution(buildDir)
		if projectFile == "" {
			return OpenResult{Success: false, Error: "No Visual Studio solution found"}
		}
		cmd := exec.Command("cmd", "/c", "start", projectFile)
		outputCallback("Opening Visual Studio: "+projectFile, ui.TypeInfo)
		if err := cmd.Run(); err != nil {
			return OpenResult{Success: false, Error: err.Error()}
		}
		return OpenResult{Success: true}

	default:
		return OpenResult{Success: false, Error: "Unsupported IDE project: " + generator}
	}
}

// ExecuteOpenEditor opens the editor (Neovim) in the build directory
func ExecuteOpenEditor(buildDir string, outputCallback func(string, ui.OutputLineType)) OpenResult {
	cmd := exec.Command("nvim", ".")
	cmd.Dir = buildDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	outputCallback("Opening Neovim in: "+buildDir, ui.TypeInfo)

	if err := cmd.Run(); err != nil {
		return OpenResult{Success: false, Error: err.Error()}
	}
	return OpenResult{Success: true}
}

// findXcodeProject searches for an Xcode project file (*.xcodeproj) in the build directory
func findXcodeProject(buildDir string) string {
	entries, err := os.ReadDir(buildDir)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), ".xcodeproj") {
			return filepath.Join(buildDir, entry.Name())
		}
	}
	return ""
}

// findVisualStudioSolution searches for a Visual Studio solution file (*.sln) in the build directory
func findVisualStudioSolution(buildDir string) string {
	entries, err := os.ReadDir(buildDir)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sln") {
			return filepath.Join(buildDir, entry.Name())
		}
	}
	return ""
}
