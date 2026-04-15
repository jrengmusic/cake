package app

import (
	"os"
	"path/filepath"

	"github.com/jrengmusic/cake/internal"
	"github.com/jrengmusic/cake/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

// startCleanAllOperation begins the clean all operation (removes entire Builds/ directory)
func (a *Application) startCleanAllOperation() (tea.Model, tea.Cmd) {
	a.enterConsoleMode(ui.OpCleanAll, "Removing all build artifacts...")
	return a, tea.Batch(a.cmdCleanAllProject(), a.cmdRefreshConsole(), a.cmdSpinnerTick())
}

// cmdCleanAllProject executes the clean all command (removes entire Builds/ directory)
func (a *Application) cmdCleanAllProject() tea.Cmd {
	return func() tea.Msg {
		// replace callback unused: operation does not produce progress lines
		appendCallback, _ := a.outputCallbacks()

		buildsDir := filepath.Join(a.projectState.WorkingDirectory, internal.BuildsDirName)

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
