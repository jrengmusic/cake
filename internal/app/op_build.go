package app

import (
	"cake/internal/ops"
	"cake/internal/ui"
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
		outputCallback := func(line string, lineType ui.OutputLineType) {
			a.outputBuffer.Append(line, lineType)
		}

		project := a.projectState.SelectedGenerator
		config := a.projectState.Configuration
		projectRoot := a.projectState.WorkingDirectory

		result := ops.ExecuteBuildProject(
			project,
			config,
			projectRoot,
			outputCallback,
		)

		return BuildCompleteMsg{
			Success:  result.Success,
			ExitCode: result.ExitCode,
			Error:    result.Error,
		}
	}
}
