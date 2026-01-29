package ops

import (
	"cake/internal/ui"
	"cake/internal/utils"
	"os/exec"
	"path/filepath"
)

type BuildResult struct {
	Success  bool
	ExitCode int
	Error    string
}

func ExecuteBuildProject(generator, config, projectRoot string, outputCallback func(string, ui.OutputLineType)) BuildResult {
	buildDir := filepath.Join(projectRoot, "Builds", utils.GetDirectoryName(generator))

	args := []string{"--build", buildDir, "--config", config}

	outputCallback("Building: "+buildDir, ui.TypeInfo)
	outputCallback("Project: "+generator, ui.TypeInfo)
	outputCallback("Configuration: "+config, ui.TypeInfo)
	outputCallback("", ui.TypeStdout)

	cmd := exec.Command("cmake", args...)
	cmd.Dir = projectRoot

	// Stream stdout/stderr using helper
	if err := utils.StreamCommand(cmd, outputCallback); err != nil {
		outputCallback("ERROR: "+err.Error(), ui.TypeStderr)
		return BuildResult{Success: false, Error: err.Error()}
	}

	// Wait for command to complete
	err := cmd.Wait()

	if err != nil {
		outputCallback("", ui.TypeStdout)
		outputCallback("ERROR: Build failed", ui.TypeStderr)
		return BuildResult{Success: false, Error: err.Error()}
	}

	outputCallback("", ui.TypeStdout)
	outputCallback("Build completed successfully", ui.TypeStatus)
	return BuildResult{Success: true, ExitCode: 0}
}
