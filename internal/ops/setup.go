package ops

import (
	"os/exec"
	"path/filepath"
	"strings"
)

type SetupResult struct {
	Success bool
	Error   string
}

func ExecuteSetupProject(workingDir, generator, config string, isMultiConfig bool, outputCallback func(string)) SetupResult {
	if workingDir == "" {
		return SetupResult{Success: false, Error: "Working directory is empty"}
	}

	if generator == "" {
		return SetupResult{Success: false, Error: "Generator is empty"}
	}

	// Determine build path based on generator type
	// Multi-config generators (Xcode, VS): Builds/<Generator>/
	// Single-config generators (Ninja, Makefiles): Builds/<Generator>/<Config>/
	var buildDir string
	if isMultiConfig {
		buildDir = filepath.Join(workingDir, "Builds", generator)
	} else {
		buildDir = filepath.Join(workingDir, "Builds", generator, config)
	}

	args := []string{
		"-G", generator,
		"-S", workingDir,
		"-B", buildDir,
	}

	// Add CMAKE_BUILD_TYPE for single-config generators
	if !isMultiConfig && config != "" {
		args = append(args, "-DCMAKE_BUILD_TYPE="+config)
	}

	outputCallback("Running: cmake " + strings.Join(args, " "))
	outputCallback("")

	cmd := exec.Command("cmake", args...)
	cmd.Dir = workingDir

	output, err := cmd.CombinedOutput()

	for _, line := range strings.Split(string(output), "\n") {
		if line != "" {
			outputCallback(line)
		}
	}

	if err != nil {
		outputCallback("")
		outputCallback("ERROR: " + err.Error())
		return SetupResult{Success: false, Error: err.Error()}
	}

	outputCallback("")
	outputCallback("Setup completed successfully: " + buildDir)
	return SetupResult{Success: true}
}
