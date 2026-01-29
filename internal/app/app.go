package app

import (
	"cake/internal/config"
	"cake/internal/state"
	"cake/internal/ui"
	"context"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Type aliases for TIT confirmation dialog pattern
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

	runningCmd    *exec.Cmd
	cancelContext context.CancelFunc

	footerHint   string
	isScanning   bool
	lastBuildDir string // Track last build dir to detect changes

	confirmDialog *ui.ConfirmationDialog // Confirmation dialog (TIT pattern)

	pendingOperation string // Track operation to execute after confirmation
}

func (a *Application) Init() tea.Cmd {
	a.projectState.ForceRefresh()
	a.menuItems = a.GenerateMenu()
	a.lastBuildDir = a.projectState.GetBuildPath()
	a.asyncState = NewAsyncState()
	// Start with auto-scroll enabled (TIT pattern: scroll follows output until user intervenes)
	a.consoleAutoScroll = true
	a.windowSize = WindowSizeHandler{}
	a.keyDispatcher = NewKeyDispatcher()

	// Register key handlers (wrapped to match dispatcher signature)
	a.keyDispatcher.Register(ModeMenu, func(app *Application, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
		return app.handleMenuKeyPress(msg)
	})
	a.keyDispatcher.Register(ModePreferences, func(app *Application, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
		return app.handlePreferencesKeyPress(msg)
	})
	a.keyDispatcher.Register(ModeConsole, func(app *Application, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
		return app.handleOperationKeyPress(msg)
	})

	// Start auto-scan ticker if enabled
	if a.config != nil && a.config.IsAutoScanEnabled() {
		return a.cmdAutoScanTick()
	}
	return nil
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

// handleAutoScanTick handles periodic auto-scan
func (a *Application) handleAutoScanTick() (tea.Model, tea.Cmd) {
	// Skip scan during async operations
	if a.asyncState.IsActive() {
		return a, a.cmdAutoScanTick()
	}

	// Track if build dir changed
	oldBuildDir := a.lastBuildDir

	// Force refresh project state
	a.projectState.ForceRefresh()
	a.isScanning = true
	a.footerHint = FooterHints["scanning"]

	newBuildDir := a.projectState.GetBuildPath()

	// Check if build dirs changed
	oldBuildExists := oldBuildDir != ""
	newBuildExists := newBuildDir != ""
	buildDirChanged := oldBuildDir != newBuildDir

	// Regenerate menu if build state changed
	if buildDirChanged || (oldBuildExists != newBuildExists) {
		a.menuItems = a.GenerateMenu()
	}

	a.lastBuildDir = newBuildDir
	a.isScanning = false

	// Restore footer hint based on mode
	if a.mode == ModeMenu {
		a.footerHint = FooterHints["menu_navigate"]
	} else if a.mode == ModePreferences {
		a.footerHint = "↑↓ navigate │ Enter change │ / back"
	}

	// Continue ticking
	return a, a.cmdAutoScanTick()
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
			a.footerHint = FooterHints["menu_navigate"]
		}
		return a, nil

	case AutoScanTickMsg:
		return a.handleAutoScanTick()

	case GenerateCompleteMsg:
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
			a.footerHint = "Generate failed: " + msg.Error
		}
		return a, nil

	case BuildCompleteMsg:
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

	case OpenIDECompleteMsg:
		a.asyncState.End()
		if msg.Success {
			a.footerHint = "IDE opened successfully"
		} else {
			a.footerHint = "Failed to open IDE: " + msg.Error
		}
		return a, nil

	case RegenerateCompleteMsg:
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
			return a, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
				return OutputRefreshMsg{}
			})
		}
		// Operation completed, stop sending refresh messages
		return a, nil
	}

	return a, nil
}

func (a *Application) View() string {
	// Get project name from CMakeLists.txt or working directory
	projectName := a.projectState.GetProjectName()
	if projectName == "" {
		projectName = filepath.Base(a.projectState.WorkingDirectory)
	}

	// Build header state (TIT pattern)
	headerState := ui.HeaderState{
		ProjectName:  projectName,
		CWD:          a.projectState.WorkingDirectory,
		Version:      "v1.0",
		VersionColor: a.theme.DimmedTextColor,
	}

	// Render header info then header with margins (TIT pattern exactly)
	headerInfo := ui.RenderHeaderInfo(a.sizing, a.theme, headerState)
	headerText := ui.RenderHeader(a.sizing, a.theme, headerInfo)

	footerText := a.GetFooterContent()

	// Render confirmation dialog if active (TIT pattern)
	if a.confirmDialog != nil && a.confirmDialog.Active {
		dialogContent := a.confirmDialog.Render(a.sizing.ContentHeight)
		return ui.RenderReactiveLayout(a.sizing, a.theme, headerText, dialogContent, footerText)
	}

	// Console mode: render directly without layout wrapper (TIT pattern)
	// Console handles its own sizing and footer placement
	if a.mode == ModeConsole {
		return a.renderConsoleMode()
	}

	var contentText string
	switch a.mode {
	case ModeMenu:
		contentText = a.renderMenuWithBanner()
	case ModePreferences:
		contentText = a.renderPreferencesWithBanner()
	default:
		contentText = a.renderMenuWithBanner()
	}

	return ui.RenderReactiveLayout(a.sizing, a.theme, headerText, contentText, footerText)
}

// handleConfirmDialogKeyPress handles Y/N keys for confirmation dialogs (TIT pattern)
func (a *Application) handleConfirmDialogKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if a.confirmDialog == nil {
		return a, nil
	}

	key := msg.String()

	switch key {
	case "ctrl+c":
		// Global quit from confirmation dialog
		if a.quitConfirmActive {
			// Second Ctrl+C - quit the app
			return a, tea.Quit
		}
		// First Ctrl+C - start timeout, keep dialog open
		a.quitConfirmActive = true
		a.quitConfirmTime = time.Now()
		a.footerHint = GetFooterMessageText(MessageCtrlCConfirm)
		return a, tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
			return TickMsg(t)
		})

	case "y", "Y":
		// User pressed Y - always select Yes
		a.confirmDialog.Active = false
		if a.pendingOperation != "" {
			op := a.pendingOperation
			a.pendingOperation = ""
			switch op {
			case "generate":
				return a.startGenerateOperation()
			case "clean":
				return a.startCleanOperation()
			case "regenerate":
				return a.startRegenerateOperation()
			}
		}
		return a, nil

	case "n", "N":
		// User pressed N - always cancel
		a.confirmDialog.Active = false
		a.confirmDialog = nil
		a.pendingOperation = ""
		return a, nil

	case "enter", " ":
		// Execute whichever button is currently selected
		if a.confirmDialog.GetSelectedButton() == ButtonYes {
			// Yes is selected - execute the operation
			a.confirmDialog.Active = false
			if a.pendingOperation != "" {
				op := a.pendingOperation
				a.pendingOperation = ""
				switch op {
				case "generate":
					return a.startGenerateOperation()
				case "clean":
					return a.startCleanOperation()
				case "regenerate":
					return a.startRegenerateOperation()
				}
			}
		} else {
			// No is selected - cancel
			a.confirmDialog.Active = false
			a.confirmDialog = nil
			a.pendingOperation = ""
		}
		return a, nil

	case "esc":
		// ESC always cancels
		a.confirmDialog.Active = false
		a.confirmDialog = nil
		a.pendingOperation = ""
		return a, nil

	case "left", "h":
		// Move to Yes button (Yes is on the left)
		a.confirmDialog.SelectYes()
		return a, nil

	case "right", "l":
		// Move to No button (No is on the right)
		a.confirmDialog.SelectNo()
		return a, nil

	default:
		return a, nil
	}
}

// renderMenuWithBanner renders menu (left 50%) + banner (right 50%)
// Both columns centered H/V, identical to TIT layout
func (a *Application) renderMenuWithBanner() string {
	// 50/50 split
	leftWidth := a.sizing.ContentInnerWidth / 2
	rightWidth := a.sizing.ContentInnerWidth - leftWidth

	// Render menu in left column using new ui.RenderCakeMenu
	menuContent := ui.RenderCakeMenu(a.menuItems, a.selectedIndex, a.theme, a.sizing.ContentHeight, leftWidth)

	menuColumn := lipgloss.NewStyle().
		Width(leftWidth).
		Height(a.sizing.ContentHeight).
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(menuContent)

	// Render banner in right column
	banner := ui.RenderBannerDynamic(rightWidth, a.sizing.ContentHeight)

	bannerColumn := lipgloss.NewStyle().
		Width(rightWidth).
		Height(a.sizing.ContentHeight).
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(banner)

	// Join horizontally
	return lipgloss.JoinHorizontal(lipgloss.Top, menuColumn, bannerColumn)
}

// GetVisibleRows returns only visible AND selectable menu items (excludes hidden rows and separator)
func (a *Application) GetVisibleRows() []ui.MenuRow {
	var visible []ui.MenuRow
	for _, row := range a.menuItems {
		if row.Visible && row.IsSelectable {
			visible = append(visible, row)
		}
	}
	return visible
}

// GetVisibleIndex returns the visible selectable index for a given row ID
// Returns -1 if row is hidden or not selectable
func (a *Application) GetVisibleIndex(rowID string) int {
	visibleIndex := 0
	for _, row := range a.menuItems {
		if row.ID == rowID {
			if row.Visible && row.IsSelectable {
				return visibleIndex
			}
			return -1
		}
		if row.Visible && row.IsSelectable {
			visibleIndex++
		}
	}
	return -1
}

// GetArrayIndex returns the array index for a given visible selectable index
func (a *Application) GetArrayIndex(visibleIdx int) int {
	visibleCount := 0
	for i, row := range a.menuItems {
		if row.Visible && row.IsSelectable {
			if visibleCount == visibleIdx {
				return i
			}
			visibleCount++
		}
	}
	return -1
}

// GetVisiblePreferenceRows returns visible preference rows (stub for preferences mode)
func (a *Application) GetVisiblePreferenceRows() []ui.MenuRow {
	return []ui.MenuRow{}
}

// ToggleRowAtIndex handles menu row toggle/action at given VISIBLE index
func (a *Application) ToggleRowAtIndex(visibleIndex int) (bool, tea.Cmd) {
	arrayIndex := a.GetArrayIndex(visibleIndex)
	if arrayIndex < 0 || arrayIndex >= len(a.menuItems) {
		return false, nil
	}

	row := a.menuItems[arrayIndex]
	return a.executeRowAction(row.ID)
}

// RowIndexByID finds the ARRAY index of a row by its ID (legacy, use GetVisibleIndex)
func (a *Application) RowIndexByID(rowID string) int {
	for i, row := range a.menuItems {
		if row.ID == rowID {
			return i
		}
	}
	return -1
}

// TogglePreferenceAtIndex toggles preference at index (stub for preferences mode)
func (a *Application) TogglePreferenceAtIndex(index int) bool {
	return false
}

// executeRowAction executes the action associated with a menu row
func (a *Application) executeRowAction(rowID string) (bool, tea.Cmd) {
	switch rowID {
	case "project":
		// Cycle to next generator
		a.projectState.CycleToNextGenerator()
		a.menuItems = a.GenerateMenu()
		return true, nil
	case "regenerate":
		// Check if build exists for selected project
		buildInfo := a.projectState.GetSelectedBuildInfo()
		if !buildInfo.Exists {
			// No existing build - just generate without confirmation
			_, cmd := a.startGenerateOperation()
			return true, cmd
		}
		// Build exists - show confirmation dialog for regenerate (default to No for safety)
		a.confirmDialog = ui.NewConfirmationDialogWithDefault(ui.ConfirmationConfig{
			Title:       "Regenerate Project",
			Explanation: "Clean and re-run CMake configuration?",
			YesLabel:    "Yes",
			NoLabel:     "No",
			ActionID:    "regenerate",
		}, a.sizing.ContentInnerWidth, &a.theme, ui.ButtonNo)
		a.confirmDialog.Active = true
		a.pendingOperation = "regenerate"
		return true, nil
	case "clean":
		// Show confirmation dialog for clean (default to No for safety)
		a.confirmDialog = ui.NewConfirmationDialogWithDefault(ui.ConfirmationConfig{
			Title:       "Clean Build Directory",
			Explanation: "Remove all build artifacts?",
			YesLabel:    "Yes",
			NoLabel:     "No",
			ActionID:    "clean",
		}, a.sizing.ContentInnerWidth, &a.theme, ui.ButtonNo)
		a.confirmDialog.Active = true
		a.pendingOperation = "clean"
		return true, nil
	case "openIde":
		_, cmd := a.startOpenIDEOperation()
		return true, cmd
	case "configuration":
		// Cycle to next configuration
		a.projectState.CycleConfiguration()
		a.menuItems = a.GenerateMenu()
		return true, nil
	case "build":
		_, cmd := a.startBuildOperation()
		return true, cmd
	}
	return false, nil
}

// renderPreferencesWithBanner renders preferences (left 50%) + banner (right 50%)
func (a *Application) renderPreferencesWithBanner() string {
	// 50/50 split
	leftWidth := a.sizing.ContentInnerWidth / 2
	rightWidth := a.sizing.ContentInnerWidth - leftWidth

	// Render preferences in left column
	prefsContent := a.renderPreferenceMenuRows(leftWidth, a.GetVisiblePreferenceRows())

	prefsColumn := lipgloss.NewStyle().
		Width(leftWidth).
		Height(a.sizing.ContentHeight).
		Align(lipgloss.Left).
		AlignVertical(lipgloss.Center).
		Render(prefsContent)

	// Render banner in right column
	banner := ui.RenderBannerDynamic(rightWidth, a.sizing.ContentHeight)

	bannerColumn := lipgloss.NewStyle().
		Width(rightWidth).
		Height(a.sizing.ContentHeight).
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(banner)

	// Join horizontally
	return lipgloss.JoinHorizontal(lipgloss.Top, prefsColumn, bannerColumn)
}

// renderPreferenceMenuRows renders preference rows with proper formatting
func (a *Application) renderPreferenceMenuRows(maxWidth int, rows []ui.MenuRow) string {
	var lines []string

	for i, row := range rows {
		if row.ID == "separator" {
			// Render separator line
			line := lipgloss.NewStyle().
				Foreground(lipgloss.Color(a.theme.SeparatorColor)).
				Render(strings.Repeat("─", maxWidth))
			lines = append(lines, line)
			continue
		}

		var line string

		// Format: "EMOJI  LABEL           VALUE"
		emoji := row.Emoji + " "
		label := row.Label
		value := row.Value

		// Calculate spacing to right-align value
		contentWidth := lipgloss.Width(emoji) + lipgloss.Width(label) + lipgloss.Width(value)
		valueSpacing := ""
		if maxWidth > contentWidth {
			valueSpacing = strings.Repeat(" ", maxWidth-contentWidth)
		}

		rowText := emoji + label + valueSpacing + value

		if i == a.selectedIndex {
			// Selected row
			line = lipgloss.NewStyle().
				Foreground(lipgloss.Color(a.theme.MainBackgroundColor)).
				Background(lipgloss.Color(a.theme.MenuSelectionBackground)).
				Bold(true).
				Width(maxWidth).
				Align(lipgloss.Left).
				Render(rowText)
		} else {
			// Normal row
			line = lipgloss.NewStyle().
				Foreground(lipgloss.Color(a.theme.LabelTextColor)).
				Width(maxWidth).
				Align(lipgloss.Left).
				Render(rowText)
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func (a *Application) renderConsoleMode() string {
	// Console height accounts for footer (TIT pattern: console fills space above footer)
	consoleHeight := a.sizing.TerminalHeight - ui.FooterHeight

	consoleContent := ui.RenderConsoleOutput(
		&a.consoleState,
		a.outputBuffer,
		a.theme,
		a.sizing.TerminalWidth,
		consoleHeight,
		a.asyncState.IsActive(),
		false,
		a.consoleAutoScroll,
	)

	footerText := a.GetFooterContent()

	// Footer only - console handles its own OUTPUT title (TIT pattern)
	footerSection := lipgloss.Place(
		a.sizing.TerminalWidth,
		ui.FooterHeight,
		lipgloss.Left,
		lipgloss.Top,
		lipgloss.NewStyle().
			Foreground(lipgloss.Color(a.theme.FooterTextColor)).
			Render(footerText),
	)

	// Join console + footer
	return lipgloss.JoinVertical(lipgloss.Left, consoleContent, footerSection)
}

func (a *Application) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle confirmation dialog keys first (works in any mode)
	if a.confirmDialog != nil && a.confirmDialog.Active {
		return a.handleConfirmDialogKeyPress(msg)
	}

	switch a.mode {
	case ModeMenu:
		return a.handleMenuKeyPress(msg)
	case ModePreferences:
		return a.handlePreferencesKeyPress(msg)
	case ModeConsole:
		return a.handleOperationKeyPress(msg)
	}
	return a, nil
}

func (a *Application) handleMenuKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	visibleRows := a.GetVisibleRows()
	visibleCount := len(visibleRows)

	// Helper: find next visible index (already skips hidden/separator via GetVisibleRows)
	findNextVisibleIndex := func(current, direction int) int {
		next := current + direction
		if next < 0 {
			return 0
		}
		if next >= visibleCount {
			return visibleCount - 1
		}
		return next
	}

	switch msg.String() {
	case "up", "k":
		// Move up using visible indices
		if a.selectedIndex > 0 {
			a.selectedIndex = findNextVisibleIndex(a.selectedIndex, -1)
		}
		return a, nil
	case "down", "j":
		// Move down using visible indices
		if a.selectedIndex < visibleCount-1 {
			a.selectedIndex = findNextVisibleIndex(a.selectedIndex, 1)
		}
		return a, nil
	case "enter", " ":
		// Execute action at selected visible index
		if a.selectedIndex >= 0 && a.selectedIndex < visibleCount {
			handled, cmd := a.ToggleRowAtIndex(a.selectedIndex)
			if handled {
				a.menuItems = a.GenerateMenu()
				// Clamp selectedIndex if visible count changed
				newVisible := a.GetVisibleRows()
				if a.selectedIndex >= len(newVisible) {
					a.selectedIndex = len(newVisible) - 1
				}
				if a.selectedIndex < 0 {
					a.selectedIndex = 0
				}
				return a, cmd
			}
		}
		return a, nil
	case "g", "G":
		// Generate/Regenerate - jump to row and execute
		idx := a.GetVisibleIndex("regenerate")
		if idx >= 0 {
			a.selectedIndex = idx
			handled, cmd := a.ToggleRowAtIndex(idx)
			if handled {
				a.menuItems = a.GenerateMenu()
				// Clamp selectedIndex if visible count changed
				newVisible := a.GetVisibleRows()
				if a.selectedIndex >= len(newVisible) {
					a.selectedIndex = len(newVisible) - 1
				}
				if a.selectedIndex < 0 {
					a.selectedIndex = 0
				}
				return a, cmd
			}
		}
		return a, nil
	case "o", "O":
		// Open IDE - jump to row and execute
		idx := a.GetVisibleIndex("openIde")
		if idx >= 0 {
			a.selectedIndex = idx
			handled, cmd := a.ToggleRowAtIndex(idx)
			if handled {
				a.menuItems = a.GenerateMenu()
				// Clamp selectedIndex if visible count changed
				newVisible := a.GetVisibleRows()
				if a.selectedIndex >= len(newVisible) {
					a.selectedIndex = len(newVisible) - 1
				}
				if a.selectedIndex < 0 {
					a.selectedIndex = 0
				}
				return a, cmd
			}
		}
		return a, nil
	case "b", "B":
		// Build - jump to row and execute
		idx := a.GetVisibleIndex("build")
		if idx >= 0 {
			a.selectedIndex = idx
			handled, cmd := a.ToggleRowAtIndex(idx)
			if handled {
				a.menuItems = a.GenerateMenu()
				// Clamp selectedIndex if visible count changed
				newVisible := a.GetVisibleRows()
				if a.selectedIndex >= len(newVisible) {
					a.selectedIndex = len(newVisible) - 1
				}
				if a.selectedIndex < 0 {
					a.selectedIndex = 0
				}
				return a, cmd
			}
		}
		return a, nil
	case "c", "C":
		// Clean - jump to row and execute
		idx := a.GetVisibleIndex("clean")
		if idx >= 0 {
			a.selectedIndex = idx
			handled, cmd := a.ToggleRowAtIndex(idx)
			if handled {
				a.menuItems = a.GenerateMenu()
				// Clamp selectedIndex if visible count changed
				newVisible := a.GetVisibleRows()
				if a.selectedIndex >= len(newVisible) {
					a.selectedIndex = len(newVisible) - 1
				}
				if a.selectedIndex < 0 {
					a.selectedIndex = 0
				}
				return a, cmd
			}
		}
		return a, nil
	case "/":
		// Toggle preferences screen (TIT pattern: / opens preferences directly)
		if a.mode == ModeMenu {
			a.mode = ModePreferences
			a.selectedIndex = 0
		} else if a.mode == ModePreferences {
			a.mode = ModeMenu
			a.selectedIndex = 0
		}
		return a, nil
	case "ctrl+c":
		return a.handleCtrlC()
	}
	return a, nil
}

func (a *Application) handlePreferencesKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	visibleRows := a.GetVisiblePreferenceRows()
	visibleCount := len(visibleRows)

	switch msg.String() {
	case "up", "k":
		// Move up, skipping separators (robust boundary check)
		if a.selectedIndex > 0 {
			nextIndex := a.selectedIndex - 1
			// Skip separator rows going up
			for nextIndex >= 0 && visibleRows[nextIndex].ID == "separator" {
				nextIndex--
			}
			// Only apply if we found a valid slot
			if nextIndex >= 0 {
				a.selectedIndex = nextIndex
			}
		}
		return a, nil
	case "down", "j":
		// Move down, skipping separators (robust boundary check)
		if a.selectedIndex < visibleCount-1 {
			nextIndex := a.selectedIndex + 1
			// Skip separator rows going down
			for nextIndex < visibleCount && visibleRows[nextIndex].ID == "separator" {
				nextIndex++
			}
			// Only apply if we found a valid slot
			if nextIndex < visibleCount {
				a.selectedIndex = nextIndex
			}
		}
		return a, nil
	case "enter", " ":
		if a.TogglePreferenceAtIndex(a.selectedIndex) {
			return a, nil
		}
		return a, nil
	case "/", "esc":
		// Return to main menu
		a.mode = ModeMenu
		a.selectedIndex = 0
		a.footerHint = FooterHints["menu_navigate"]
		return a, nil
	case "ctrl+c":
		return a.handleCtrlC()
	}
	return a, nil
}

func (a *Application) handleOperationKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up":
		a.consoleState.ScrollUp()
		a.consoleAutoScroll = false // Disable auto-scroll on manual scroll (TIT pattern)
		return a, nil
	case "down":
		a.consoleState.ScrollDown()
		a.consoleAutoScroll = false // Disable auto-scroll on manual scroll (TIT pattern)
		return a, nil
	case "esc":
		if a.asyncState.IsActive() {
			// Kill running process (like Ctrl+C)
			if a.runningCmd != nil && a.runningCmd.Process != nil {
				a.runningCmd.Process.Kill()
			}
			if a.cancelContext != nil {
				a.cancelContext()
			}
			a.asyncState.Abort()
			// Print abort message to console using stderr color from theme
			a.outputBuffer.Append("", ui.TypeStdout)
			a.outputBuffer.Append("Operation aborted by user", ui.TypeStderr)
			a.outputBuffer.Append("Press ESC to return to menu", ui.TypeInfo)
		} else {
			// Return to menu when idle
			a.mode = ModeMenu
			a.selectedIndex = 0
			a.menuItems = a.GenerateMenu()
			a.footerHint = FooterHints["menu_navigate"]
		}
		return a, nil
	case "ctrl+c":
		return a.handleCtrlC()
	}
	return a, nil
}

func (a *Application) handleCtrlC() (tea.Model, tea.Cmd) {
	if a.asyncState.IsActive() {
		a.footerHint = FooterHints["operation_wait"]
		return a, nil
	}

	if a.quitConfirmActive {
		return a, tea.Quit
	}

	a.quitConfirmActive = true
	a.quitConfirmTime = time.Now()
	a.footerHint = GetFooterMessageText(MessageCtrlCConfirm)

	return a, tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// cmdRefreshConsole sends periodic refresh messages while async operation is active
// This forces UI re-renders to display streaming output in real-time
func (a *Application) cmdRefreshConsole() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return OutputRefreshMsg{}
	})
}
