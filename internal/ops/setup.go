package ops

import (
	"bufio"
	"cake/internal/ui"
	"context"
	"os/exec"
	"path/filepath"
	"strings"
)

type SetupResult struct {
	Success bool
	Error   string
}

func ExecuteSetupProject(ctx context.Context, workingDir, generator, config string, outputCallback func(string, ui.OutputLineType)) SetupResult {
	if workingDir == "" {
		return SetupResult{Success: false, Error: "Working directory is empty"}
	}

	if generator == "" {
		return SetupResult{Success: false, Error: "Generator is empty"}
	}

	buildDir := filepath.Join(workingDir, "Builds", generator)

	args := []string{
		"-G", generator,
		"-S", workingDir,
		"-B", buildDir,
	}

	outputCallback("Running: cmake "+strings.Join(args, " "), ui.TypeInfo)
	outputCallback("", ui.TypeStdout)

	cmd := exec.CommandContext(ctx, "cmake", args...)
	cmd.Dir = workingDir

	// Get stdout pipe for streaming
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		outputCallback("ERROR: Failed to create stdout pipe", ui.TypeStderr)
		return SetupResult{Success: false, Error: err.Error()}
	}

	// Get stderr pipe for streaming
	stderr, err := cmd.StderrPipe()
	if err != nil {
		outputCallback("ERROR: Failed to create stderr pipe", ui.TypeStderr)
		return SetupResult{Success: false, Error: err.Error()}
	}

	// Start command
	if err := cmd.Start(); err != nil {
		outputCallback("ERROR: Failed to start command", ui.TypeStderr)
		return SetupResult{Success: false, Error: err.Error()}
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

	if ctx.Err() == context.Canceled {
		// Message already printed by ESC handler, just return
		return SetupResult{Success: false, Error: "aborted"}
	}

	if err != nil {
		outputCallback("", ui.TypeStdout)
		outputCallback("ERROR: "+err.Error(), ui.TypeStderr)
		return SetupResult{Success: false, Error: err.Error()}
	}

	outputCallback("", ui.TypeStdout)
	outputCallback("Setup completed successfully: "+buildDir, ui.TypeStatus)
	return SetupResult{Success: true}
}
