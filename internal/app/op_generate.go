package app

import (
	"cake/internal/ops"
	"cake/internal/ui"
	"context"
	tea "github.com/charmbracelet/bubbletea"
)

// startGenerateOperation begins the generate/regenerate operation
func (a *Application) startGenerateOperation() (tea.Model, tea.Cmd) {
	a.mode = ModeConsole
	a.asyncState.Start()
	a.outputBuffer.Clear()
	a.consoleAutoScroll = true // Re-enable auto-scroll for new operation (TIT pattern)
	a.footerHint = GetFooterMessageText(MessageSetupInProgress)

	// Create cancellable context for process termination
	ctx, cancel := context.WithCancel(context.Background())
	a.cancelContext = cancel

	return a, tea.Batch(a.cmdGenerateProject(ctx), a.cmdRefreshConsole())
}

// cmdGenerateProject executes the generate/regenerate command
func (a *Application) cmdGenerateProject(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		outputCallback := func(line string, lineType ui.OutputLineType) {
			a.outputBuffer.Append(line, lineType)
		}

		generator := a.projectState.SelectedGenerator
		config := a.projectState.Configuration
		projectRoot := a.projectState.WorkingDirectory

		result := ops.ExecuteSetupProject(
			ctx,
			projectRoot,
			generator,
			config,
			outputCallback,
		)

		return GenerateCompleteMsg{
			Success: result.Success,
			Error:   result.Error,
		}
	}
}
