package app

import (
	"cake/internal/config"
	"cake/internal/ops"
	"cake/internal/state"
	"cake/internal/ui"
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

	consoleState ui.ConsoleOutState
	outputBuffer *ui.OutputBuffer

	asyncOperationActive  bool
	asyncOperationAborted bool

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
	if a.asyncOperationActive {
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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.sizing = ui.CalculateDynamicSizing(msg.Width, msg.Height)
		return a, nil

	case tea.KeyMsg:
		return a.handleKeyPress(msg)

	case TickMsg:
		if a.quitConfirmActive && time.Since(a.quitConfirmTime) > 3*time.Second {
			a.quitConfirmActive = false
			a.footerHint = FooterHints["menu_navigate"]
		}
		return a, nil

	case AutoScanTickMsg:
		return a.handleAutoScanTick()

	case GenerateCompleteMsg:
		a.asyncOperationActive = false
		a.projectState.ForceRefresh()
		a.menuItems = a.GenerateMenu()
		if msg.Success {
			a.footerHint = GetFooterMessageText(MessageOperationComplete)
		} else {
			a.footerHint = "Generate failed: " + msg.Error
		}
		return a, nil

	case BuildCompleteMsg:
		a.asyncOperationActive = false
		if msg.Success {
			a.footerHint = GetFooterMessageText(MessageOperationComplete)
		} else {
			a.footerHint = "Build failed: " + msg.Error
		}
		return a, nil

	case CleanCompleteMsg:
		a.asyncOperationActive = false
		a.projectState.ForceRefresh()
		a.menuItems = a.GenerateMenu()
		if msg.Success {
			a.footerHint = GetFooterMessageText(MessageOperationComplete)
		} else {
			a.footerHint = "Clean failed: " + msg.Error
		}
		return a, nil

	case OpenIDECompleteMsg:
		a.asyncOperationActive = false
		if msg.Success {
			a.footerHint = "IDE opened successfully"
		} else {
			a.footerHint = "Failed to open IDE: " + msg.Error
		}
		return a, nil

	case OpenEditorCompleteMsg:
		a.asyncOperationActive = false
		if msg.Success {
			a.footerHint = "Editor closed"
		} else {
			a.footerHint = "Failed to open editor: " + msg.Error
		}
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

	footerText := a.footerHint

	// Render confirmation dialog if active (TIT pattern)
	if a.confirmDialog != nil && a.confirmDialog.Active {
		dialogContent := a.confirmDialog.Render(a.sizing.ContentHeight)
		return ui.RenderReactiveLayout(a.sizing, a.theme, headerText, dialogContent, footerText)
	}

	var contentText string
	switch a.mode {
	case ModeMenu:
		contentText = a.renderMenuWithBanner()
	case ModePreferences:
		contentText = a.renderPreferencesWithBanner()
	case ModeConsole:
		contentText = a.renderConsoleMode()
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
	case "y", "Y", "enter":
		// User selected Yes
		if a.confirmDialog.GetSelectedButton() == ButtonYes {
			a.confirmDialog.Active = false
			if a.pendingOperation != "" {
				op := a.pendingOperation
				a.pendingOperation = ""
				switch op {
				case "generate":
					return a.startGenerateOperation()
				case "clean":
					return a.startCleanOperation()
				}
			}
		}
		return a, nil

	case "n", "N", "esc":
		// User selected No or cancelled
		a.confirmDialog.Active = false
		a.confirmDialog = nil
		a.pendingOperation = ""
		return a, nil

	case "left", "h":
		// Move to No button
		a.confirmDialog.SelectNo()
		return a, nil

	case "right", "l":
		// Move to Yes button
		a.confirmDialog.SelectYes()
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

// GetVisibleRows returns visible menu items (all items, no filtering needed with MenuRow)
func (a *Application) GetVisibleRows() []ui.MenuRow {
	return a.menuItems
}

// GetVisiblePreferenceRows returns visible preference rows (stub for preferences mode)
func (a *Application) GetVisiblePreferenceRows() []ui.MenuRow {
	return []ui.MenuRow{}
}

// ToggleRowAtIndex handles menu row toggle/action at given index
func (a *Application) ToggleRowAtIndex(index int) (bool, tea.Cmd) {
	if index < 0 || index >= len(a.menuItems) {
		return false, nil
	}

	row := a.menuItems[index]
	return a.executeRowAction(row.ID)
}

// RowIndexByID finds the index of a row by its ID
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
	case "generator":
		// Cycle to next generator
		a.projectState.CycleToNextGenerator()
		return true, nil
	case "regenerate":
		_, cmd := a.startGenerateOperation()
		return true, cmd
	case "openIde":
		_, cmd := a.startOpenIDEOperation()
		return true, cmd
	case "configuration":
		// Cycle to next configuration
		a.projectState.CycleConfiguration()
		return true, nil
	case "build":
		_, cmd := a.startBuildOperation()
		return true, cmd
	case "clean":
		_, cmd := a.startCleanOperation()
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
	return ui.RenderConsoleOutput(
		&a.consoleState,
		a.outputBuffer,
		a.theme,
		a.sizing.ContentInnerWidth,
		a.sizing.ContentHeight,
		a.asyncOperationActive,
		false,
		a.asyncOperationActive,
	)
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
				return a, cmd
			}
		}
		return a, nil
	case "g", "G":
		// Generate/Regenerate - jump to row and execute
		idx := a.RowIndexByID("generate")
		if idx >= 0 {
			a.selectedIndex = idx
			handled, cmd := a.ToggleRowAtIndex(idx)
			if handled {
				a.menuItems = a.GenerateMenu()
				return a, cmd
			}
		}
		return a, nil
	case "o", "O":
		// Open IDE - jump to row and execute
		idx := a.RowIndexByID("openIde")
		if idx >= 0 {
			a.selectedIndex = idx
			handled, cmd := a.ToggleRowAtIndex(idx)
			if handled {
				a.menuItems = a.GenerateMenu()
				return a, cmd
			}
		}
		return a, nil
	case "b", "B":
		// Build - jump to row and execute
		idx := a.RowIndexByID("build")
		if idx >= 0 {
			a.selectedIndex = idx
			handled, cmd := a.ToggleRowAtIndex(idx)
			if handled {
				a.menuItems = a.GenerateMenu()
				return a, cmd
			}
		}
		return a, nil
	case "c", "C":
		// Clean - jump to row and execute
		idx := a.RowIndexByID("clean")
		if idx >= 0 {
			a.selectedIndex = idx
			handled, cmd := a.ToggleRowAtIndex(idx)
			if handled {
				a.menuItems = a.GenerateMenu()
				return a, cmd
			}
		}
		return a, nil
	case "/":
		// Toggle preferences screen (TIT pattern: / opens preferences directly)
		if a.mode == ModeMenu {
			a.mode = ModePreferences
			a.selectedIndex = 0
			a.footerHint = FooterHints["preferences"]
		} else if a.mode == ModePreferences {
			a.mode = ModeMenu
			a.selectedIndex = 0
			a.footerHint = FooterHints["menu_navigate"]
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
		return a, nil
	case "down":
		a.consoleState.ScrollDown()
		return a, nil
	case "esc":
		if !a.asyncOperationActive {
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
	if a.asyncOperationActive {
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

// startGenerateOperation begins the generate/regenerate operation
func (a *Application) startGenerateOperation() (tea.Model, tea.Cmd) {
	a.mode = ModeConsole
	a.asyncOperationActive = true
	a.outputBuffer.Clear()
	a.footerHint = GetFooterMessageText(MessageSetupInProgress)

	return a, a.cmdGenerateProject()
}

// cmdGenerateProject executes the generate/regenerate command
func (a *Application) cmdGenerateProject() tea.Cmd {
	return func() tea.Msg {
		outputCallback := func(line string) {
			a.outputBuffer.Append(line, ui.TypeStdout)
		}

		generator := a.projectState.SelectedGenerator
		config := a.projectState.Configuration
		projectRoot := a.projectState.WorkingDirectory

		// Determine if generator is multi-config
		isMultiConfig := false
		for _, gen := range a.projectState.AvailableGenerators {
			if gen.Name == generator {
				isMultiConfig = gen.IsMultiConfig
				break
			}
		}

		result := ops.ExecuteSetupProject(
			projectRoot,
			generator,
			config,
			isMultiConfig,
			outputCallback,
		)

		return GenerateCompleteMsg{
			Success: result.Success,
			Error:   result.Error,
		}
	}
}

// startOpenIDEOperation opens the IDE for the selected generator
func (a *Application) startOpenIDEOperation() (tea.Model, tea.Cmd) {
	a.asyncOperationActive = true
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

		buildPath := a.projectState.GetBuildPath()
		generator := a.projectState.SelectedGenerator

		result := ops.ExecuteOpenIDE(generator, buildPath, outputCallback)

		return OpenIDECompleteMsg{
			Success: result.Success,
			Error:   result.Error,
		}
	}
}

// startOpenEditorOperation opens the editor in the build directory
func (a *Application) startOpenEditorOperation() (tea.Model, tea.Cmd) {
	a.asyncOperationActive = true
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

		buildPath := a.projectState.GetBuildPath()

		result := ops.ExecuteOpenEditor(buildPath, outputCallback)

		return OpenEditorCompleteMsg{
			Success: result.Success,
			Error:   result.Error,
		}
	}
}

// startBuildOperation begins the build operation
func (a *Application) startBuildOperation() (tea.Model, tea.Cmd) {
	a.mode = ModeConsole
	a.asyncOperationActive = true
	a.outputBuffer.Clear()
	a.footerHint = GetFooterMessageText(MessageBuildInProgress)

	return a, a.cmdBuildProject()
}

// cmdBuildProject executes the build command
func (a *Application) cmdBuildProject() tea.Cmd {
	return func() tea.Msg {
		outputCallback := func(line string) {
			a.outputBuffer.Append(line, ui.TypeStdout)
		}

		generator := a.projectState.SelectedGenerator
		config := a.projectState.Configuration
		projectRoot := a.projectState.WorkingDirectory

		// Determine if generator is multi-config
		isMultiConfig := false
		for _, gen := range a.projectState.AvailableGenerators {
			if gen.Name == generator {
				isMultiConfig = gen.IsMultiConfig
				break
			}
		}

		result := ops.ExecuteBuildProject(
			generator,
			config,
			projectRoot,
			isMultiConfig,
			outputCallback,
		)

		return BuildCompleteMsg{
			Success:  result.Success,
			ExitCode: result.ExitCode,
			Error:    result.Error,
		}
	}
}

// startCleanOperation begins the clean operation
func (a *Application) startCleanOperation() (tea.Model, tea.Cmd) {
	a.mode = ModeConsole
	a.asyncOperationActive = true
	a.outputBuffer.Clear()
	a.footerHint = GetFooterMessageText(MessageCleanInProgress)

	return a, a.cmdCleanProject()
}

// cmdCleanProject executes the clean command
func (a *Application) cmdCleanProject() tea.Cmd {
	return func() tea.Msg {
		outputCallback := func(line string) {
			a.outputBuffer.Append(line, ui.TypeStdout)
		}

		generator := a.projectState.SelectedGenerator
		config := a.projectState.Configuration
		projectRoot := a.projectState.WorkingDirectory

		// Determine if generator is multi-config
		isMultiConfig := false
		for _, gen := range a.projectState.AvailableGenerators {
			if gen.Name == generator {
				isMultiConfig = gen.IsMultiConfig
				break
			}
		}

		result := ops.ExecuteCleanProject(
			generator,
			config,
			projectRoot,
			isMultiConfig,
			outputCallback,
		)

		return CleanCompleteMsg{
			Success: result.Success,
			Error:   result.Error,
		}
	}
}
