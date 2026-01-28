package app

import (
	"cake/internal/ops"
	"cake/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

// startCleanOperation begins the clean operation
func (a *Application) startCleanOperation() (tea.Model, tea.Cmd) {
	a.mode = ModeConsole
	a.asyncState.Start()
	a.outputBuffer.Clear()
	a.footerHint = GetFooterMessageText(MessageCleanInProgress)
	return a, a.cmdCleanProject()
}

// cmdCleanProject executes the clean command
func (a *Application) cmdCleanProject() tea.Cmd {
	return func() tea.Msg {
		outputCallback := func(line string) {
			a.outputBuffer.Append(line, ui.TypeStdout)
		}

		generator := a.projectState.SelectedGenerator
		config := a.projectState.Configuration
		projectRoot := a.projectState.WorkingDirectory

		isMultiConfig := a.projectState.IsGeneratorMultiConfig(generator)

		result := ops.ExecuteCleanProject(
			generator,
			config,
			projectRoot,
			isMultiConfig,
			outputCallback,
		)

		return CleanCompleteMsg{
			Success: result.Success,
			Error:   result.Error,
		}
	}
}
