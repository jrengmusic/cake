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

func ExecuteBuildProject(generator, config, projectRoot string, vsEnv []string, appendCallback func(string, ui.OutputLineType), replaceCallback func(string, ui.OutputLineType)) BuildResult {
	buildDir := filepath.Join(projectRoot, "Builds", utils.GetDirectoryName(generator))

	args := []string{"--build", buildDir, "--config", config}

	appendCallback("Building: "+buildDir, ui.TypeInfo)
	appendCallback("Project: "+generator, ui.TypeInfo)
	appendCallback("Configuration: "+config, ui.TypeInfo)
	appendCallback("", ui.TypeStdout)

	cmakePath := utils.FindExecutableInEnv("cmake", vsEnv)
	cmd := exec.Command(cmakePath, args...)
	cmd.Dir = projectRoot
	if len(vsEnv) > 0 {
		cmd.Env = vsEnv
	}

	// Stream stdout/stderr using helper
	if err := utils.StreamCommand(cmd, appendCallback, replaceCallback); err != nil {
		appendCallback("ERROR: "+err.Error(), ui.TypeStderr)
		return BuildResult{Success: false, Error: err.Error()}
	}

	// Wait for command to complete
	err := cmd.Wait()

	if err != nil {
		appendCallback("", ui.TypeStdout)
		appendCallback("ERROR: Build failed", ui.TypeStderr)
		return BuildResult{Success: false, Error: err.Error()}
	}

	appendCallback("", ui.TypeStdout)
	appendCallback("Build completed successfully", ui.TypeStatus)
	return BuildResult{Success: true, ExitCode: 0}
}
