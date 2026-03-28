package app

import (
	"cake/internal/ops"
	tea "github.com/charmbracelet/bubbletea"
)

// startCleanOperation begins the clean operation
func (a *Application) startCleanOperation() (tea.Model, tea.Cmd) {
	a.mode = ModeConsole
	a.asyncState.Start()
	a.outputBuffer.Clear()
	a.consoleAutoScroll = true // Re-enable auto-scroll for new operation (TIT pattern)
	a.footerHint = GetFooterMessageText(MessageCleanInProgress)
	return a, tea.Batch(a.cmdCleanProject(), a.cmdRefreshConsole())
}

// cmdCleanProject executes the clean command
func (a *Application) cmdCleanProject() tea.Cmd {
	return func() tea.Msg {
		appendCallback, _ := a.outputCallbacks()

		project := a.projectState.SelectedProject
		config := a.projectState.Configuration
		projectRoot := a.projectState.WorkingDirectory

		result := ops.ExecuteCleanProject(
			project,
			config,
			projectRoot,
			appendCallback,
		)

		return CleanCompleteMsg{
			Success: result.Success,
			Error:   result.Error,
		}
	}
}
