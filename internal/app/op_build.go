package app

import (
	"cake/internal/ops"
	tea "github.com/charmbracelet/bubbletea"
)

// startBuildOperation begins the build operation
func (a *Application) startBuildOperation() (tea.Model, tea.Cmd) {
	a.mode = ModeConsole
	a.asyncState.Start()
	a.outputBuffer.Clear()
	a.consoleAutoScroll = true // Re-enable auto-scroll for new operation (TIT pattern)
	a.footerHint = GetFooterMessageText(MessageBuildInProgress)
	return a, tea.Batch(a.cmdBuildProject(), a.cmdRefreshConsole())
}

// cmdBuildProject executes the build command
func (a *Application) cmdBuildProject() tea.Cmd {
	return func() tea.Msg {
		appendCallback, replaceCallback := a.outputCallbacks()

		project := a.projectState.SelectedProject
		config := a.projectState.Configuration
		projectRoot := a.projectState.WorkingDirectory

		result := ops.ExecuteBuildProject(
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
