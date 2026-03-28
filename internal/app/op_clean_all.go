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
		appendCallback, _ := a.outputCallbacks()

		buildsDir := filepath.Join(a.projectState.WorkingDirectory, "Builds")

		appendCallback("", ui.TypeStdout)
		appendCallback("Cleaning all projects...", ui.TypeInfo)
		appendCallback("Target: "+buildsDir, ui.TypeStdout)
		appendCallback("", ui.TypeStdout)

		// Remove the entire Builds directory
		err := os.RemoveAll(buildsDir)
		if err != nil {
			appendCallback("Error: Failed to remove Builds directory: "+err.Error(), ui.TypeStderr)
			return CleanAllCompleteMsg{
				Success: false,
				Error:   err.Error(),
			}
		}

		appendCallback("✓ All build artifacts removed successfully", ui.TypeInfo)
		appendCallback("", ui.TypeStdout)
		appendCallback("Press ESC to return to menu", ui.TypeInfo)

		return CleanAllCompleteMsg{
			Success: true,
			Error:   "",
		}
	}
}
