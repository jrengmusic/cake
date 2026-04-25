package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/jrengmusic/cake/internal/ui"
)

// startAsyncOperation marks the app as running an async operation and clears prior output.
// Used by operations that do not present a console view (e.g. launching an external IDE).
func (a *Application) startAsyncOperation(footerHint string) {
	a.asyncState.Start(ui.OpNone)
	a.outputBuffer.Clear()
	a.footerHint = footerHint
}

// enterConsoleMode starts an async operation AND switches the UI to the streaming console view.
func (a *Application) enterConsoleMode(op ui.OpType, footerHint string) {
	a.spinnerFrame = 0
	a.asyncState.Start(op)
	a.outputBuffer.Clear()
	a.footerHint = footerHint
	a.mode = ModeConsole
	a.consoleAutoScroll = true
}

// cmdRefreshConsole sends periodic refresh messages while async operation is active
// This forces UI re-renders to display streaming output in real-time
func (a *Application) cmdRefreshConsole() tea.Cmd {
	return tea.Tick(CacheRefreshInterval, func(t time.Time) tea.Msg {
		return OutputRefreshMsg{}
	})
}

// cmdSpinnerTick drives the spinner animation at its own rate, independent of console refresh
func (a *Application) cmdSpinnerTick() tea.Cmd {
	return tea.Tick(SpinnerTickInterval, func(t time.Time) tea.Msg {
		return SpinnerTickMsg{}
	})
}

// outputCallbacks returns the standard append and replace callbacks for streaming operations.
// Append adds a new line to the buffer. Replace overwrites the last line (for progress output).
func (a *Application) outputCallbacks() (func(string, ui.OutputLineType), func(string, ui.OutputLineType)) {
	appendFn := func(line string, lineType ui.OutputLineType) {
		a.outputBuffer.Append(line, lineType)
	}
	replaceFn := func(line string, lineType ui.OutputLineType) {
		a.outputBuffer.ReplaceLast(line, lineType)
	}
	return appendFn, replaceFn
}

// cmdAutoScanTick returns a command that sends AutoScanTickMsg at the configured interval
func (a *Application) cmdAutoScanTick() tea.Cmd {
	if a.config == nil {
		return nil
	}
	interval := time.Duration(a.config.AutoScanInterval()) * time.Minute
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return AutoScanTickMsg{}
	})
}

func (a *Application) isAutoScanIdle() bool {
	return !a.asyncState.IsActive() && time.Since(a.lastActivityTime) >= IdleScanThreshold
}

func (a *Application) refreshProjectStateAndMenu() {
	oldBuildDir := a.lastBuildDir
	a.projectState.ForceRefresh()
	a.isScanning = true
	a.footerHint = FooterHints["scanning"]

	newBuildDir := a.projectState.GetBuildPath()
	oldBuildExists := oldBuildDir != ""
	newBuildExists := newBuildDir != ""
	buildDirChanged := oldBuildDir != newBuildDir

	if buildDirChanged || (oldBuildExists != newBuildExists) {
		a.menuItems = a.GenerateMenu()
	}

	a.lastBuildDir = newBuildDir
	a.isScanning = false
}

func (a *Application) restoreFooterHintAfterScan() {
	if a.mode == ModeMenu {
		a.footerHint = FooterHints["menu_navigate"]
	} else if a.mode == ModePreferences {
		a.footerHint = "↑↓ navigate │ Enter change │ / back"
	}
}

// handleAutoScanTick handles periodic auto-scan with lazy update
func (a *Application) handleAutoScanTick() (tea.Model, tea.Cmd) {
	if !a.isAutoScanIdle() {
		return a, a.cmdAutoScanTick()
	}
	a.refreshProjectStateAndMenu()
	a.restoreFooterHintAfterScan()
	return a, a.cmdAutoScanTick()
}
