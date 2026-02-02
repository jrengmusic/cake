package app

import (
	"context"
	"os"
	"path/filepath"

	"cake/internal/ops"
	"cake/internal/ui"
	"cake/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
)

func (a *Application) startRegenerateOperation() (tea.Model, tea.Cmd) {
	a.mode = ModeConsole
	a.asyncState.Start()
	a.outputBuffer.Clear()
	a.consoleAutoScroll = true // Re-enable auto-scroll for new operation (TIT pattern)
	a.footerHint = "Regenerating project..."

	// Create cancellable context for process termination
	ctx, cancel := context.WithCancel(context.Background())
	a.cancelContext = cancel

	return a, tea.Batch(a.cmdRegenerateProject(ctx), a.cmdRefreshConsole())
}

func (a *Application) cmdRegenerateProject(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		outputCallback := func(line string, lineType ui.OutputLineType) {
			a.outputBuffer.Append(line, lineType)
		}

		project := a.projectState.SelectedProject
		projectRoot := a.projectState.WorkingDirectory
		buildDir := filepath.Join(projectRoot, "Builds", utils.GetDirectoryName(project))

		// Step 1: Clean
		outputCallback("=== Step 1: Clean ===", ui.TypeInfo)
		outputCallback("", ui.TypeStdout)

		if _, err := os.Stat(buildDir); err == nil {
			outputCallback("Removing: "+buildDir, ui.TypeStdout)
			if err := os.RemoveAll(buildDir); err != nil {
				outputCallback("ERROR: Clean failed", ui.TypeStderr)
				return RegenerateCompleteMsg{Success: false, Error: err.Error()}
			}
			outputCallback("Clean completed", ui.TypeStatus)
		} else {
			outputCallback("No build directory to clean", ui.TypeWarning)
		}

		outputCallback("", ui.TypeStdout)
		outputCallback("=== Step 2: Generate ===", ui.TypeInfo)

		// Step 2: Generate
		result := ops.ExecuteSetupProject(
			ctx,
			projectRoot,
			project,
			"",
			outputCallback,
		)

		return RegenerateCompleteMsg{
			Success: result.Success,
			Error:   result.Error,
		}
	}
}
