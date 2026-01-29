package ops

import (
	"cake/internal/ui"
	"cake/internal/utils"
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

	buildDir := filepath.Join(workingDir, "Builds", utils.GetDirectoryName(generator))

	args := []string{
		"-G", generator,
		"-S", workingDir,
		"-B", buildDir,
	}

	outputCallback("Running: cmake "+strings.Join(args, " "), ui.TypeInfo)
	outputCallback("", ui.TypeStdout)

	cmd := exec.CommandContext(ctx, "cmake", args...)
	cmd.Dir = workingDir

	// Stream stdout/stderr using helper
	if err := utils.StreamCommand(cmd, outputCallback); err != nil {
		outputCallback("ERROR: "+err.Error(), ui.TypeStderr)
		return SetupResult{Success: false, Error: err.Error()}
	}

	// Wait for command to complete
	err := cmd.Wait()

	if ctx.Err() == context.Canceled {
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
