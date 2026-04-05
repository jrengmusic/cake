package ops

import (
	"github.com/jrengmusic/cake/internal"
	"github.com/jrengmusic/cake/internal/ui"
	"github.com/jrengmusic/cake/internal/utils"
	"context"
	"os/exec"
	"path/filepath"
)

type BuildResult struct {
	Success  bool
	ExitCode int
	Error    string
}

func buildBuildCommand(ctx context.Context, buildDir, config, projectRoot string, vsEnv []string) *exec.Cmd {
	args := []string{"--build", buildDir, "--config", config}
	cmakePath := utils.FindExecutableInEnv("cmake", vsEnv)
	cmd := exec.CommandContext(ctx, cmakePath, args...)
	cmd.Dir = projectRoot
	if len(vsEnv) > 0 {
		cmd.Env = vsEnv
	}
	return cmd
}

func ExecuteBuildProject(ctx context.Context, generator, config, projectRoot string, vsEnv []string, appendCallback func(string, ui.OutputLineType), replaceCallback func(string, ui.OutputLineType)) BuildResult {
	buildDir := filepath.Join(projectRoot, internal.BuildsDirName, utils.GetDirectoryName(generator))
	cmd := buildBuildCommand(ctx, buildDir, config, projectRoot, vsEnv)

	appendCallback("Building: "+buildDir, ui.TypeInfo)
	appendCallback("Project: "+generator, ui.TypeInfo)
	appendCallback("Configuration: "+config, ui.TypeInfo)
	appendCallback("", ui.TypeStdout)

	if err := utils.StreamCommand(cmd, appendCallback, replaceCallback); err != nil {
		appendCallback("ERROR: "+err.Error(), ui.TypeStderr)
		return BuildResult{Success: false, Error: err.Error()}
	}

	err := cmd.Wait()

	if ctx.Err() == context.Canceled {
		return BuildResult{Success: false, Error: "aborted"}
	}

	if err != nil {
		appendCallback("", ui.TypeStdout)
		appendCallback("ERROR: Build failed", ui.TypeStderr)
		return BuildResult{Success: false, Error: err.Error()}
	}

	appendCallback("", ui.TypeStdout)
	appendCallback("Build completed successfully", ui.TypeStatus)
	return BuildResult{Success: true, ExitCode: 0}
}
