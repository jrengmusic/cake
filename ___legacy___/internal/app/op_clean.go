package app

import (
	"github.com/jrengmusic/cake/internal/ops"
	"github.com/jrengmusic/cake/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

// startCleanOperation begins the clean operation
func (a *Application) startCleanOperation() (tea.Model, tea.Cmd) {
	a.enterConsoleMode(ui.OpClean, GetFooterMessageText(MessageCleanInProgress))
	return a, tea.Batch(a.cmdCleanProject(), a.cmdRefreshConsole(), a.cmdSpinnerTick())
}

// cmdCleanProject executes the clean command
func (a *Application) cmdCleanProject() tea.Cmd {
	return func() tea.Msg {
		// replace callback unused: operation does not produce progress lines
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
