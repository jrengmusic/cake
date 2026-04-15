package app

import (
	"github.com/jrengmusic/cake/internal"
	"github.com/jrengmusic/cake/internal/config"
	"github.com/jrengmusic/cake/internal/state"
	"github.com/jrengmusic/cake/internal/ui"
	"context"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Type aliases for confirmation dialog
type ButtonSelection = ui.ButtonSelection

const (
	ButtonYes ButtonSelection = ui.ButtonYes
	ButtonNo  ButtonSelection = ui.ButtonNo
)

type Application struct {
	width  int
	height int

	sizing ui.DynamicSizing
	theme  ui.Theme

	mode          AppMode
	selectedIndex int
	menuItems     []ui.MenuRow

	projectState *state.ProjectState
	config       *config.Config

	quitConfirmActive bool
	quitConfirmTime   time.Time

	consoleState      ui.ConsoleOutState
	outputBuffer      *ui.OutputBuffer
	consoleAutoScroll bool // Auto-scroll console to bottom (disabled on manual scroll)

	asyncState    *AsyncState
	windowSize    WindowSizeHandler
	keyDispatcher *KeyDispatcher

	cancelContext context.CancelFunc
	killTree      func()

	footerHint   string
	isScanning   bool
	lastBuildDir string // Track last build dir to detect changes

	confirmDialog *ui.ConfirmationDialog // Confirmation dialog

	pendingOperation   string // Track operation to execute after confirmation
	buildAfterGenerate bool   // Chain build after generate when project not yet generated

	lastActivityTime time.Time // Track last user activity for lazy auto-scan

	vsEnv []string // Captured Visual Studio environment (Windows only)

	spinnerFrame int // Current braille spinner animation frame index
}

func (a *Application) registerKeyHandlers() {
	a.keyDispatcher.Register(ModeMenu, func(app *Application, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
		return app.handleMenuKeyPress(msg)
	})
	a.keyDispatcher.Register(ModePreferences, func(app *Application, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
		return app.handlePreferencesKeyPress(msg)
	})
	a.keyDispatcher.Register(ModeConsole, func(app *Application, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
		return app.handleOperationKeyPress(msg)
	})
	a.keyDispatcher.Register(ModeInvalidProject, func(app *Application, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
		return app.handleInvalidProjectKeyPress(msg)
	})
}

func (a *Application) Init() tea.Cmd {
	a.projectState.ForceRefresh()
	a.menuItems = a.GenerateMenu()
	a.lastBuildDir = a.projectState.GetBuildPath()
	a.asyncState = NewAsyncState()
	a.consoleAutoScroll = true
	a.windowSize = WindowSizeHandler{}
	a.keyDispatcher = NewKeyDispatcher()
	a.lastActivityTime = time.Now()
	a.registerKeyHandlers()

	if a.config != nil && a.config.IsAutoScanEnabled() {
		return a.cmdAutoScanTick()
	}
	return nil
}

func (a *Application) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if a.windowSize.CanHandle(msg) {
		return a.windowSize.Handle(a, msg)
	}

	// Handle confirmation dialog keys FIRST (before mode-specific handlers)
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if a.confirmDialog != nil && a.confirmDialog.Active {
			return a.handleConfirmDialogKeyPress(keyMsg)
		}
	}

	if a.keyDispatcher.CanHandle(msg) {
		return a.keyDispatcher.Handle(a, msg)
	}

	switch msg := msg.(type) {
	case TickMsg:
		if a.quitConfirmActive && time.Since(a.quitConfirmTime) > 3*time.Second {
			a.quitConfirmActive = false
			a.footerHint = a.GetDefaultFooterHint()
		}
		return a, nil

	case AutoScanTickMsg:
		return a.handleAutoScanTick()

	case GenerateCompleteMsg:
		if a.cancelContext != nil {
			a.cancelContext()
			a.cancelContext = nil
		}
		if a.killTree != nil {
			a.killTree()
			a.killTree = nil
		}
		a.asyncState.End()
		if a.asyncState.IsAborted() {
			a.asyncState.ClearAborted()
			a.buildAfterGenerate = false
			a.footerHint = "Operation aborted"
			return a, nil
		}
		a.projectState.ForceRefresh()
		a.menuItems = a.GenerateMenu()
		if msg.Success {
			// If build was requested but project wasn't generated yet, chain into build now
			if a.buildAfterGenerate {
				a.buildAfterGenerate = false
				_, cmd := a.startBuildOperation()
				return a, cmd
			}
			a.footerHint = GetFooterMessageText(MessageOperationComplete)
		} else {
			a.buildAfterGenerate = false
			a.footerHint = "Generate failed: " + msg.Error
		}
		return a, nil

	case BuildCompleteMsg:
		if a.cancelContext != nil {
			a.cancelContext()
			a.cancelContext = nil
		}
		if a.killTree != nil {
			a.killTree()
			a.killTree = nil
		}
		a.asyncState.End()
		if a.asyncState.IsAborted() {
			a.asyncState.ClearAborted()
			a.footerHint = "Operation aborted"
			return a, nil
		}
		if msg.Success {
			a.footerHint = GetFooterMessageText(MessageOperationComplete)
		} else {
			a.footerHint = "Build failed: " + msg.Error
		}
		return a, nil

	case CleanCompleteMsg:
		a.asyncState.End()
		if a.asyncState.IsAborted() {
			a.asyncState.ClearAborted()
			a.footerHint = "Operation aborted"
			return a, nil
		}
		a.projectState.ForceRefresh()
		a.menuItems = a.GenerateMenu()
		if msg.Success {
			a.footerHint = GetFooterMessageText(MessageOperationComplete)
		} else {
			a.footerHint = "Clean failed: " + msg.Error
		}
		return a, nil

	case CleanAllCompleteMsg:
		a.asyncState.End()
		if a.asyncState.IsAborted() {
			a.asyncState.ClearAborted()
			a.footerHint = "Operation aborted"
			return a, nil
		}
		a.projectState.ForceRefresh()
		a.menuItems = a.GenerateMenu()
		if msg.Success {
			a.footerHint = "All builds cleaned successfully"
		} else {
			a.footerHint = "Clean All failed: " + msg.Error
		}
		// Stay in console mode - user presses ESC to return
		return a, nil

	case OpenIDECompleteMsg:
		a.asyncState.End()
		if msg.Success {
			a.footerHint = "IDE opened successfully"
		} else {
			a.footerHint = "Failed to open IDE: " + msg.Error
		}
		return a, nil

	case RegenerateCompleteMsg:
		if a.cancelContext != nil {
			a.cancelContext()
			a.cancelContext = nil
		}
		if a.killTree != nil {
			a.killTree()
			a.killTree = nil
		}
		a.asyncState.End()
		if a.asyncState.IsAborted() {
			a.asyncState.ClearAborted()
			a.footerHint = "Operation aborted"
			return a, nil
		}
		a.projectState.ForceRefresh()
		a.menuItems = a.GenerateMenu()
		if msg.Success {
			a.footerHint = GetFooterMessageText(MessageOperationComplete)
		} else {
			a.footerHint = "Regenerate failed: " + msg.Error
		}
		return a, nil

	case OpenEditorCompleteMsg:
		a.asyncState.End()
		if msg.Success {
			a.footerHint = "Editor closed"
		} else {
			a.footerHint = "Failed to open editor: " + msg.Error
		}
		return a, nil

	case OutputRefreshMsg:
		// Force re-render to display updated console output
		// If operation still active, schedule next refresh tick
		if a.asyncState.IsActive() {
			return a, tea.Tick(CacheRefreshInterval, func(t time.Time) tea.Msg {
				return OutputRefreshMsg{}
			})
		}
		// Operation completed, stop sending refresh messages
		return a, nil

	case SpinnerTickMsg:
		var cmd tea.Cmd
		if a.asyncState.IsActive() {
			a.spinnerFrame++
			cmd = tea.Tick(SpinnerTickInterval, func(t time.Time) tea.Msg {
				return SpinnerTickMsg{}
			})
		}
		return a, cmd
	}

	return a, nil
}

func (a *Application) renderModeContent() string {
	switch a.mode {
	case ModeInvalidProject:
		return ui.RenderCakeLieBanner(a.sizing.ContentInnerWidth, a.sizing.ContentHeight)
	case ModeMenu:
		return a.renderMenuWithBanner()
	case ModePreferences:
		return a.renderPreferencesWithBanner()
	default:
		return a.renderMenuWithBanner()
	}
}

func (a *Application) View() string {
	projectName := a.projectState.GetProjectName()
	if projectName == "" {
		projectName = filepath.Base(a.projectState.WorkingDirectory)
	}

	headerState := ui.HeaderState{
		ProjectName: projectName,
		CWD:         a.projectState.WorkingDirectory,
		Version:     internal.AppVersion,
	}

	headerInfo := ui.RenderHeaderInfo(a.sizing, a.theme, headerState)
	headerText := ui.RenderHeader(a.sizing, a.theme, headerInfo)
	footerText := a.GetFooterContent()

	if a.confirmDialog != nil && a.confirmDialog.Active {
		dialogContent := a.confirmDialog.Render(a.sizing.ContentHeight)
		return ui.RenderReactiveLayout(a.sizing, a.theme, headerText, dialogContent, footerText)
	}

	if a.mode == ModeConsole {
		return a.renderConsoleMode()
	}

	return ui.RenderReactiveLayout(a.sizing, a.theme, headerText, a.renderModeContent(), footerText)
}
