package app

import (
	"github.com/jrengmusic/cake/internal/ops"
	tea "github.com/charmbracelet/bubbletea"
)

// startOpenIDEOperation opens the IDE for the selected project
func (a *Application) startOpenIDEOperation() (tea.Model, tea.Cmd) {
	a.startAsyncOperation("Opening IDE...")
	return a, a.cmdOpenIDE()
}

// cmdOpenIDE executes the open IDE command
func (a *Application) cmdOpenIDE() tea.Cmd {
	return func() tea.Msg {
		// replace callback unused: operation does not produce progress lines
		appendCallback, _ := a.outputCallbacks()

		project := a.projectState.SelectedProject
		config := a.projectState.Configuration
		projectRoot := a.projectState.WorkingDirectory

		result := ops.ExecuteOpenIDE(
			project,
			config,
			projectRoot,
			appendCallback,
		)

		return OpenIDECompleteMsg{
			Success: result.Success,
			Error:   result.Error,
		}
	}
}

// startOpenEditorOperation opens the editor in the build directory
func (a *Application) startOpenEditorOperation() (tea.Model, tea.Cmd) {
	a.startAsyncOperation("Opening editor...")
	return a, a.cmdOpenEditor()
}

// cmdOpenEditor executes the open editor command
func (a *Application) cmdOpenEditor() tea.Cmd {
	return func() tea.Msg {
		// replace callback unused: operation does not produce progress lines
		appendCallback, _ := a.outputCallbacks()

		buildPath := a.projectState.GetBuildPath()

		result := ops.ExecuteOpenEditor(buildPath, appendCallback)

		return OpenEditorCompleteMsg{
			Success: result.Success,
			Error:   result.Error,
		}
	}
}
