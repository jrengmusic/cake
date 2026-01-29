package app

import (
	"cake/internal/ops"
	"cake/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

// startOpenIDEOperation opens the IDE for the selected project
func (a *Application) startOpenIDEOperation() (tea.Model, tea.Cmd) {
	a.asyncState.Start()
	a.outputBuffer.Clear()
	a.footerHint = "Opening IDE..."

	return a, a.cmdOpenIDE()
}

// cmdOpenIDE executes the open IDE command
func (a *Application) cmdOpenIDE() tea.Cmd {
	return func() tea.Msg {
		outputCallback := func(line string, lineType ui.OutputLineType) {
			a.outputBuffer.Append(line, lineType)
		}

		project := a.projectState.SelectedGenerator
		config := a.projectState.Configuration
		projectRoot := a.projectState.WorkingDirectory

		result := ops.ExecuteOpenIDE(
			project,
			config,
			projectRoot,
			outputCallback,
		)

		return OpenIDECompleteMsg{
			Success: result.Success,
			Error:   result.Error,
		}
	}
}

// startOpenEditorOperation opens the editor in the build directory
func (a *Application) startOpenEditorOperation() (tea.Model, tea.Cmd) {
	a.asyncState.Start()
	a.outputBuffer.Clear()
	a.footerHint = "Opening editor..."

	return a, a.cmdOpenEditor()
}

// cmdOpenEditor executes the open editor command
func (a *Application) cmdOpenEditor() tea.Cmd {
	return func() tea.Msg {
		outputCallback := func(line string, lineType ui.OutputLineType) {
			a.outputBuffer.Append(line, lineType)
		}

		buildPath := a.projectState.GetBuildPath()

		result := ops.ExecuteOpenEditor(buildPath, outputCallback)

		return OpenEditorCompleteMsg{
			Success: result.Success,
			Error:   result.Error,
		}
	}
}
