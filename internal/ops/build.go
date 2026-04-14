package ops

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/jrengmusic/cake/internal"
	"github.com/jrengmusic/cake/internal/ui"
	"github.com/jrengmusic/cake/internal/utils"
)

type BuildResult struct {
	Success  bool
	ExitCode int
	Error    string
}

func ExecuteBuildProject(ctx context.Context, generator, config, projectRoot string, vsEnv []string, appendCallback func(string, ui.OutputLineType), replaceCallback func(string, ui.OutputLineType), onProcessTreeStarted func(*utils.ProcessTree)) BuildResult {
	buildDir := filepath.Join(projectRoot, internal.BuildsDirName, utils.GetDirectoryName(generator))

	args := []string{"--build", buildDir, "--config", config}

	appendCallback("Building: "+buildDir, ui.TypeInfo)
	appendCallback("Project: "+generator, ui.TypeInfo)
	appendCallback("Configuration: "+config, ui.TypeInfo)
	appendCallback("", ui.TypeStdout)

	cmakePath := utils.FindExecutableInEnv("cmake", vsEnv)
	cmd := exec.CommandContext(ctx, cmakePath, args...)
	cmd.Dir = projectRoot
	if len(vsEnv) > 0 {
		cmd.Env = vsEnv
	}

	tree, streamErr := utils.StreamCommand(cmd, appendCallback, replaceCallback, onProcessTreeStarted)

	result := BuildResult{Success: false}
	if streamErr != nil {
		appendCallback("ERROR: "+streamErr.Error(), ui.TypeStderr)
		result.Error = fmt.Errorf("ExecuteBuildProject: StreamCommand: %w", streamErr).Error()
	} else {
		defer tree.Close()

		waitErr := cmd.Wait()
		abortedByUser := ctx.Err() == context.Canceled

		if abortedByUser {
			result.Error = "aborted"
		} else if waitErr != nil {
			appendCallback("", ui.TypeStdout)
			appendCallback("ERROR: Build failed", ui.TypeStderr)
			result.Error = fmt.Errorf("ExecuteBuildProject: cmake --build failed: %w", waitErr).Error()
		} else {
			appendCallback("", ui.TypeStdout)
			appendCallback("Build completed successfully", ui.TypeStatus)
			result.Success = true
			result.ExitCode = 0
		}
	}

	return result
}
