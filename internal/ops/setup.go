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

func ExecuteSetupProject(ctx context.Context, workingDir, generator, config string, vsEnv []string, appendCallback func(string, ui.OutputLineType), replaceCallback func(string, ui.OutputLineType)) SetupResult {
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

	appendCallback("Running: cmake "+strings.Join(args, " "), ui.TypeInfo)
	appendCallback("", ui.TypeStdout)

	cmakePath := utils.FindExecutableInEnv("cmake", vsEnv)
	cmd := exec.CommandContext(ctx, cmakePath, args...)
	cmd.Dir = workingDir
	if len(vsEnv) > 0 {
		cmd.Env = vsEnv
	}

	// Stream stdout/stderr using helper
	if err := utils.StreamCommand(cmd, appendCallback, replaceCallback); err != nil {
		appendCallback("ERROR: "+err.Error(), ui.TypeStderr)
		return SetupResult{Success: false, Error: err.Error()}
	}

	// Wait for command to complete
	err := cmd.Wait()

	if ctx.Err() == context.Canceled {
		return SetupResult{Success: false, Error: "aborted"}
	}

	if err != nil {
		appendCallback("", ui.TypeStdout)
		appendCallback("ERROR: "+err.Error(), ui.TypeStderr)
		return SetupResult{Success: false, Error: err.Error()}
	}

	appendCallback("", ui.TypeStdout)
	appendCallback("Setup completed successfully: "+buildDir, ui.TypeStatus)
	return SetupResult{Success: true}
}
