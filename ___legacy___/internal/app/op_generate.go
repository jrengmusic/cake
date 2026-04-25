package app

import (
	"context"

	"github.com/jrengmusic/cake/internal/ops"
	"github.com/jrengmusic/cake/internal/ui"
	"github.com/jrengmusic/cake/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
)

// startGenerateOperation begins the generate/regenerate operation
func (a *Application) startGenerateOperation() (tea.Model, tea.Cmd) {
	a.enterConsoleMode(ui.OpGenerate, GetFooterMessageText(MessageSetupInProgress))

	// Create cancellable context for process termination
	ctx, cancel := context.WithCancel(context.Background())
	a.cancelContext = cancel

	return a, tea.Batch(a.cmdGenerateProject(ctx), a.cmdRefreshConsole(), a.cmdSpinnerTick())
}

// cmdGenerateProject executes the generate/regenerate command
func (a *Application) cmdGenerateProject(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		appendCallback, replaceCallback := a.outputCallbacks()

		generator := a.projectState.SelectedProject
		config := a.projectState.Configuration
		projectRoot := a.projectState.WorkingDirectory

		result := ops.ExecuteSetupProject(
			ctx,
			projectRoot,
			generator,
			config,
			a.vsEnv,
			appendCallback,
			replaceCallback,
			func(tree *utils.ProcessTree) {
				a.killTree = tree.Close
			},
		)

		return GenerateCompleteMsg{
			Success: result.Success,
			Error:   result.Error,
		}
	}
}
