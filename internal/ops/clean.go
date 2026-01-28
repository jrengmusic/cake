package ops

import (
	"os"
	"path/filepath"
)

type CleanResult struct {
	Success bool
	Error   string
}

func ExecuteCleanProject(generator, config string, projectRoot string, isMultiConfig bool, outputCallback func(string)) CleanResult {
	// Determine build path based on generator type
	// Multi-config generators (Xcode, VS): Builds/<Generator>/
	// Single-config generators (Ninja, Makefiles): Builds/<Generator>/<Config>/
	var buildDir string
	if isMultiConfig {
		buildDir = filepath.Join(projectRoot, "Builds", generator)
	} else {
		buildDir = filepath.Join(projectRoot, "Builds", generator, config)
	}

	if buildDir == "" {
		return CleanResult{Success: false, Error: "Build directory is empty"}
	}

	outputCallback("Cleaning: " + buildDir)
	outputCallback("")

	dirToDelete, err := os.Stat(buildDir)
	if err != nil {
		outputCallback("ERROR: Directory not found: " + buildDir)
		return CleanResult{Success: false, Error: err.Error()}
	}

	if !dirToDelete.IsDir() {
		outputCallback("ERROR: Not a directory: " + buildDir)
		return CleanResult{Success: false, Error: "Not a directory"}
	}

	outputCallback("Deleting: " + buildDir)

	err = os.RemoveAll(buildDir)
	if err != nil {
		outputCallback("ERROR: Failed to remove directory")
		outputCallback(err.Error())
		return CleanResult{Success: false, Error: err.Error()}
	}

	outputCallback("Successfully removed build directory")
	outputCallback("")
	outputCallback("Clean completed successfully")
	return CleanResult{Success: true}
}
