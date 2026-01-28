package ops

import (
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
func ExecuteOpenIDE(generator, buildDir string, outputCallback func(string)) OpenResult {
	var projectFile string

	switch generator {
	case "Xcode":
		projectFile = findXcodeProject(buildDir)
		if projectFile == "" {
			return OpenResult{Success: false, Error: "No Xcode project found"}
		}
		cmd := exec.Command("open", projectFile)
		outputCallback("Opening Xcode project: " + projectFile)
		if err := cmd.Run(); err != nil {
			return OpenResult{Success: false, Error: err.Error()}
		}
		return OpenResult{Success: true}

	case "Visual Studio 17 2022", "Visual Studio 16 2019":
		projectFile = findVisualStudioSolution(buildDir)
		if projectFile == "" {
			return OpenResult{Success: false, Error: "No Visual Studio solution found"}
		}
		cmd := exec.Command("cmd", "/c", "start", projectFile)
		outputCallback("Opening Visual Studio: " + projectFile)
		if err := cmd.Run(); err != nil {
			return OpenResult{Success: false, Error: err.Error()}
		}
		return OpenResult{Success: true}

	default:
		return OpenResult{Success: false, Error: "Unsupported IDE generator: " + generator}
	}
}

// ExecuteOpenEditor opens the editor (Neovim) in the build directory
func ExecuteOpenEditor(buildDir string, outputCallback func(string)) OpenResult {
	cmd := exec.Command("nvim", ".")
	cmd.Dir = buildDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	outputCallback("Opening Neovim in: " + buildDir)

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
