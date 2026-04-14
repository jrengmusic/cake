package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jrengmusic/cake/internal"
	"github.com/jrengmusic/cake/internal/ops"
	"github.com/jrengmusic/cake/internal/ui"
	"github.com/jrengmusic/cake/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
)

func (a *Application) startRegenerateOperation() (tea.Model, tea.Cmd) {
	a.enterConsoleMode("Regenerating project...")

	// Create cancellable context for process termination
	ctx, cancel := context.WithCancel(context.Background())
	a.cancelContext = cancel

	return a, tea.Batch(a.cmdRegenerateProject(ctx), a.cmdRefreshConsole())
}

func (a *Application) cmdRegenerateProject(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		appendCallback, replaceCallback := a.outputCallbacks()

		project := a.projectState.SelectedProject
		projectRoot := a.projectState.WorkingDirectory
		buildDir := filepath.Join(projectRoot, internal.BuildsDirName, utils.GetDirectoryName(project))

		// Step 1: Clean
		appendCallback("=== Step 1: Clean ===", ui.TypeInfo)
		appendCallback("", ui.TypeStdout)

		cleanSucceeded := true
		var cleanErr error
		if _, statErr := os.Stat(buildDir); statErr == nil {
			appendCallback("Removing: "+buildDir, ui.TypeStdout)
			if removeErr := os.RemoveAll(buildDir); removeErr != nil {
				cleanSucceeded = false
				cleanErr = fmt.Errorf("cmdRegenerateProject: remove build dir: %w", removeErr)
				appendCallback("ERROR: Clean failed", ui.TypeStderr)
			} else {
				appendCallback("Clean completed", ui.TypeStatus)
			}
		} else {
			appendCallback("No build directory to clean", ui.TypeWarning)
		}

		appendCallback("", ui.TypeStdout)
		appendCallback("=== Step 2: Generate ===", ui.TypeInfo)

		// Step 2: Generate
		result := RegenerateCompleteMsg{Success: false}
		if cleanSucceeded {
			setupResult := ops.ExecuteSetupProject(
				ctx,
				projectRoot,
				project,
				"",
				a.vsEnv,
				appendCallback,
				replaceCallback,
				func(tree *utils.ProcessTree) {
					a.killTree = tree.Close
				},
			)
			result.Success = setupResult.Success
			result.Error = setupResult.Error
		} else {
			result.Error = cleanErr.Error()
		}

		return result
	}
}
