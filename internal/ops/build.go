package ops

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type BuildResult struct {
	Success  bool
	ExitCode int
	Error    string
}

func ExecuteBuildProject(generator, config string, projectRoot string, isMultiConfig bool, outputCallback func(string)) BuildResult {
	// Determine build path based on generator type
	// Multi-config generators (Xcode, VS): Builds/<Generator>/
	// Single-config generators (Ninja, Makefiles): Builds/<Generator>/<Config>/
	var buildDir string
	var args []string

	if isMultiConfig {
		buildDir = filepath.Join(projectRoot, "Builds", generator)
		args = []string{"--build", buildDir, "--config", config}
	} else {
		buildDir = filepath.Join(projectRoot, "Builds", generator, config)
		args = []string{"--build", buildDir}
	}

	// Verify build directory exists
	if !directoryExists(buildDir) {
		return BuildResult{Success: false, Error: "Build directory not found: " + buildDir}
	}

	outputCallback("Building: " + buildDir)
	outputCallback("Generator: " + generator)
	outputCallback("Configuration: " + config)
	outputCallback("")

	// Use cmake --build for all generators (unified approach)
	cmd := exec.Command("cmake", args...)
	cmd.Dir = projectRoot

	outputCallback("Running: cmake " + strings.Join(args, " "))
	outputCallback("")

	output, err := cmd.CombinedOutput()

	for _, line := range strings.Split(string(output), "\n") {
		if line != "" {
			outputCallback(line)
		}
	}

	if err != nil {
		outputCallback("")
		outputCallback("ERROR: Build failed")
		if exitErr, ok := err.(*exec.ExitError); ok {
			return BuildResult{Success: false, ExitCode: exitErr.ExitCode(), Error: err.Error()}
		}
		return BuildResult{Success: false, ExitCode: 1, Error: err.Error()}
	}

	outputCallback("")
	outputCallback("Build completed successfully")
	return BuildResult{Success: true, ExitCode: 0}
}

func directoryExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
