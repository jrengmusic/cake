package ops

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jrengmusic/cake/internal"
	"github.com/jrengmusic/cake/internal/ui"
	"github.com/jrengmusic/cake/internal/utils"
)

type SetupResult struct {
	Success bool
	Error   string
}

func ExecuteSetupProject(ctx context.Context, workingDir, generator, config string, vsEnv []string, appendCallback func(string, ui.OutputLineType), replaceCallback func(string, ui.OutputLineType), onProcessTreeStarted func(*utils.ProcessTree)) SetupResult {
	if workingDir == "" {
		return SetupResult{Success: false, Error: "Working directory is empty"}
	}

	if generator == "" {
		return SetupResult{Success: false, Error: "Generator is empty"}
	}

	buildDir := filepath.Join(workingDir, internal.BuildsDirName, utils.GetDirectoryName(generator))

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

	tree, streamErr := utils.StreamCommand(cmd, appendCallback, replaceCallback, onProcessTreeStarted)

	result := SetupResult{Success: false}
	if streamErr != nil {
		appendCallback("ERROR: "+streamErr.Error(), ui.TypeStderr)
		result.Error = fmt.Errorf("ExecuteSetupProject: StreamCommand: %w", streamErr).Error()
	} else {
		defer tree.Close()

		waitErr := cmd.Wait()
		abortedByUser := ctx.Err() == context.Canceled

		if abortedByUser {
			result.Error = "aborted"
		} else if waitErr != nil {
			appendCallback("", ui.TypeStdout)
			appendCallback("ERROR: "+waitErr.Error(), ui.TypeStderr)
			result.Error = fmt.Errorf("ExecuteSetupProject: cmake configure failed: %w", waitErr).Error()
		} else {
			appendCallback("", ui.TypeStdout)
			appendCallback("Setup completed successfully: "+buildDir, ui.TypeStatus)
			result.Success = true
		}
	}

	return result
}
