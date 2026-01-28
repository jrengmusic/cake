package app

import (
	"cake/internal/ops"
	"cake/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

// startOpenIDEOperation opens the IDE for the selected generator
func (a *Application) startOpenIDEOperation() (tea.Model, tea.Cmd) {
	a.asyncState.Start()
	a.outputBuffer.Clear()
	a.footerHint = "Opening IDE..."

	return a, a.cmdOpenIDE()
}

// cmdOpenIDE executes the open IDE command
func (a *Application) cmdOpenIDE() tea.Cmd {
	return func() tea.Msg {
		outputCallback := func(line string) {
			a.outputBuffer.Append(line, ui.TypeStdout)
		}

		generator := a.projectState.SelectedGenerator
		config := a.projectState.Configuration

		buildPath := a.projectState.GetBuildDirectory(generator, config)

		result := ops.ExecuteOpenIDE(generator, buildPath, outputCallback)

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
		outputCallback := func(line string) {
			a.outputBuffer.Append(line, ui.TypeStdout)
		}

		generator := a.projectState.SelectedGenerator
		config := a.projectState.Configuration

		buildPath := a.projectState.GetBuildDirectory(generator, config)

		result := ops.ExecuteOpenEditor(buildPath, outputCallback)

		return OpenEditorCompleteMsg{
			Success: result.Success,
			Error:   result.Error,
		}
	}
}
