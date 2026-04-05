package app

import (
	"github.com/jrengmusic/cake/internal/ops"
	"context"
	tea "github.com/charmbracelet/bubbletea"
)

// startBuildOperation begins the build operation
func (a *Application) startBuildOperation() (tea.Model, tea.Cmd) {
	a.enterConsoleMode(GetFooterMessageText(MessageBuildInProgress))
	ctx, cancel := context.WithCancel(context.Background())
	a.cancelContext = cancel
	return a, tea.Batch(a.cmdBuildProject(ctx), a.cmdRefreshConsole())
}

// cmdBuildProject executes the build command
func (a *Application) cmdBuildProject(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		appendCallback, replaceCallback := a.outputCallbacks()

		project := a.projectState.SelectedProject
		config := a.projectState.Configuration
		projectRoot := a.projectState.WorkingDirectory

		result := ops.ExecuteBuildProject(
			ctx,
			project,
			config,
			projectRoot,
			a.vsEnv,
			appendCallback,
			replaceCallback,
		)

		return BuildCompleteMsg{
			Success:  result.Success,
			ExitCode: result.ExitCode,
			Error:    result.Error,
		}
	}
}
