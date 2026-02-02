package app

import (
	"os"
	"path/filepath"

	"cake/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

// startCleanAllOperation begins the clean all operation (removes entire Builds/ directory)
func (a *Application) startCleanAllOperation() (tea.Model, tea.Cmd) {
	a.mode = ModeConsole
	a.asyncState.Start()
	a.outputBuffer.Clear()
	a.consoleAutoScroll = true // Re-enable auto-scroll for new operation (TIT pattern)
	a.footerHint = "Removing all build artifacts..."
	return a, tea.Batch(a.cmdCleanAllProject(), a.cmdRefreshConsole())
}

// cmdCleanAllProject executes the clean all command (removes entire Builds/ directory)
func (a *Application) cmdCleanAllProject() tea.Cmd {
	return func() tea.Msg {
		outputCallback := func(line string, lineType ui.OutputLineType) {
			a.outputBuffer.Append(line, lineType)
		}

		buildsDir := filepath.Join(a.projectState.WorkingDirectory, "Builds")

		outputCallback("", ui.TypeStdout)
		outputCallback("Cleaning all projects...", ui.TypeInfo)
		outputCallback("Target: "+buildsDir, ui.TypeStdout)
		outputCallback("", ui.TypeStdout)

		// Remove the entire Builds directory
		err := os.RemoveAll(buildsDir)
		if err != nil {
			outputCallback("Error: Failed to remove Builds directory: "+err.Error(), ui.TypeStderr)
			return CleanAllCompleteMsg{
				Success: false,
				Error:   err.Error(),
			}
		}

		outputCallback("✓ All build artifacts removed successfully", ui.TypeInfo)
		outputCallback("", ui.TypeStdout)
		outputCallback("Press ESC to return to menu", ui.TypeInfo)

		return CleanAllCompleteMsg{
			Success: true,
			Error:   "",
		}
	}
}
