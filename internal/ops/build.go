package ops

import (
	"bufio"
	"cake/internal/ui"
	"os/exec"
	"path/filepath"
)

type BuildResult struct {
	Success  bool
	ExitCode int
	Error    string
}

func ExecuteBuildProject(generator, config, projectRoot string, outputCallback func(string, ui.OutputLineType)) BuildResult {
	buildDir := filepath.Join(projectRoot, "Builds", generator)

	args := []string{"--build", buildDir, "--config", config}

	outputCallback("Building: "+buildDir, ui.TypeInfo)
	outputCallback("Project: "+generator, ui.TypeInfo)
	outputCallback("Configuration: "+config, ui.TypeInfo)
	outputCallback("", ui.TypeStdout)

	cmd := exec.Command("cmake", args...)
	cmd.Dir = projectRoot

	// Get stdout pipe for streaming
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		outputCallback("ERROR: Failed to create stdout pipe", ui.TypeStderr)
		return BuildResult{Success: false, Error: err.Error()}
	}

	// Get stderr pipe for streaming
	stderr, err := cmd.StderrPipe()
	if err != nil {
		outputCallback("ERROR: Failed to create stderr pipe", ui.TypeStderr)
		return BuildResult{Success: false, Error: err.Error()}
	}

	// Start command
	if err := cmd.Start(); err != nil {
		outputCallback("ERROR: Failed to start command", ui.TypeStderr)
		return BuildResult{Success: false, Error: err.Error()}
	}

	// Stream stdout in goroutine
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			if line != "" {
				outputCallback(line, ui.TypeStdout)
			}
		}
	}()

	// Stream stderr in goroutine
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			if line != "" {
				outputCallback(line, ui.TypeStderr)
			}
		}
	}()

	// Wait for command to complete
	err = cmd.Wait()

	if err != nil {
		outputCallback("", ui.TypeStdout)
		outputCallback("ERROR: Build failed", ui.TypeStderr)
		return BuildResult{Success: false, Error: err.Error()}
	}

	outputCallback("", ui.TypeStdout)
	outputCallback("Build completed successfully", ui.TypeStatus)
	return BuildResult{Success: true, ExitCode: 0}
}
