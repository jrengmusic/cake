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
	a.footerHint = GetFooterMessageText(MessageBuildInProgress)
	return a, a.cmdBuildProject()
}

// cmdBuildProject executes the build command
func (a *Application) cmdBuildProject() tea.Cmd {
	return func() tea.Msg {
		outputCallback := func(line string) {
			a.outputBuffer.Append(line, ui.TypeStdout)
		}

		generator := a.projectState.SelectedGenerator
		config := a.projectState.Configuration
		projectRoot := a.projectState.WorkingDirectory

		isMultiConfig := a.projectState.IsGeneratorMultiConfig(generator)

		result := ops.ExecuteBuildProject(
			generator,
			config,
			projectRoot,
			isMultiConfig,
			outputCallback,
		)

		return BuildCompleteMsg{
			Success:  result.Success,
			ExitCode: result.ExitCode,
			Error:    result.Error,
		}
	}
}
