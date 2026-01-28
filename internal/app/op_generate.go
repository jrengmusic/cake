package app

import (
	"cake/internal/ops"
	"cake/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

// startGenerateOperation begins the generate/regenerate operation
func (a *Application) startGenerateOperation() (tea.Model, tea.Cmd) {
	a.mode = ModeConsole
	a.asyncState.Start()
	a.outputBuffer.Clear()
	a.footerHint = GetFooterMessageText(MessageSetupInProgress)
	return a, a.cmdGenerateProject()
}

// cmdGenerateProject executes the generate/regenerate command
func (a *Application) cmdGenerateProject() tea.Cmd {
	return func() tea.Msg {
		outputCallback := func(line string) {
			a.outputBuffer.Append(line, ui.TypeStdout)
		}

		generator := a.projectState.SelectedGenerator
		config := a.projectState.Configuration
		projectRoot := a.projectState.WorkingDirectory

		isMultiConfig := a.projectState.IsGeneratorMultiConfig(generator)

		result := ops.ExecuteSetupProject(
			projectRoot,
			generator,
			config,
			isMultiConfig,
			outputCallback,
		)

		return GenerateCompleteMsg{
			Success: result.Success,
			Error:   result.Error,
		}
	}
}
