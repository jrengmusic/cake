package app

import (
	"cake/internal/ui"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

// PreferenceRow represents a single row in the preference-style menu
// Format: EMOJI | LABEL | VALUE
type PreferenceRow struct {
	ID        string // Unique identifier for the row
	Emoji     string // Visual indicator (⚙️, 🚀, 📂, etc.)
	Label     string // Left column text
	Value     string // Right column text (for toggles)
	IsToggle  bool   // true if pressing Enter/Space cycles the value
	IsAction  bool   // true if pressing Enter executes an action
	Separator bool   // true if this is a section separator
	Visible   bool   // false if row should be hidden
}

// GenerateMenu creates the preference-style menu for cake
// Format: EMOJI | LABEL | VALUE (like TIT preferences)
func (a *Application) GenerateMenu() []PreferenceRow {
	rows := []PreferenceRow{}

	// === SECTION 1: Generator Selection (ALWAYS VISIBLE) ===

	// Generator row - TOGGLE (cycle available generators)
	rows = append(rows, PreferenceRow{
		ID:       "generator",
		Emoji:    "⚙️",
		Label:    "Generator",
		Value:    a.projectState.GetGeneratorLabel(),
		IsToggle: true,
		IsAction: false,
		Visible:  true,
	})

	// Generate/Regenerate row - ACTION (dynamic label based on build existence)
	generateLabel := "Generate"
	buildInfo := a.projectState.GetSelectedBuildInfo()
	if buildInfo.Exists {
		generateLabel = "Regenerate"
	}

	rows = append(rows, PreferenceRow{
		ID:       "generate",
		Emoji:    "🚀",
		Label:    generateLabel,
		Value:    "",
		IsToggle: false,
		IsAction: true,
		Visible:  a.projectState.CanGenerate(),
	})

	// Open IDE - CONDITIONAL (only for IDE generators with existing build)
	if a.projectState.CanOpenIDE() && buildInfo.Exists {
		rows = append(rows, PreferenceRow{
			ID:       "openIde",
			Emoji:    "📂",
			Label:    "Open IDE",
			Value:    "",
			IsToggle: false,
			IsAction: true,
			Visible:  true,
		})
	}

	// Open Editor - CONDITIONAL (only for CLI generators with existing build)
	if !a.projectState.CanOpenIDE() && buildInfo.Exists {
		rows = append(rows, PreferenceRow{
			ID:       "openEditor",
			Emoji:    "📝",
			Label:    "Open Editor",
			Value:    "",
			IsToggle: false,
			IsAction: true,
			Visible:  true,
		})
	}

	// === SEPARATOR ===
	rows = append(rows, PreferenceRow{
		ID:        "sep1",
		Emoji:     "",
		Label:     "",
		Value:     "",
		IsToggle:  false,
		IsAction:  false,
		Separator: true,
		Visible:   true,
	})

	// === SECTION 2: Configuration & Build ===

	// Configuration row - TOGGLE (Debug <-> Release)
	rows = append(rows, PreferenceRow{
		ID:       "configuration",
		Emoji:    "🏗️",
		Label:    "Configuration",
		Value:    a.projectState.Configuration,
		IsToggle: true,
		IsAction: false,
		Visible:  true,
	})

	// Build row - ACTION (only if build exists and is configured)
	rows = append(rows, PreferenceRow{
		ID:       "build",
		Emoji:    "🔨",
		Label:    "Build",
		Value:    "",
		IsToggle: false,
		IsAction: true,
		Visible:  a.projectState.CanBuild(),
	})

	// Clean row - ACTION (only if build exists)
	rows = append(rows, PreferenceRow{
		ID:       "clean",
		Emoji:    "🧹",
		Label:    "Clean",
		Value:    "",
		IsToggle: false,
		IsAction: true,
		Visible:  buildInfo.Exists,
	})

	return rows
}

// GetVisibleRows returns only the visible rows (for rendering)
func (a *Application) GetVisibleRows() []PreferenceRow {
	rows := a.GenerateMenu()
	visible := []PreferenceRow{}
	for _, row := range rows {
		if row.Visible {
			visible = append(visible, row)
		}
	}
	return visible
}

// GetVisibleRowCount returns the count of visible rows (for selection)
func (a *Application) GetVisibleRowCount() int {
	return len(a.GetVisibleRows())
}

// ToggleRowAtIndex cycles or activates the row at the given index
// Returns (handled bool, cmd tea.Cmd) - cmd is non-nil for async operations
func (a *Application) ToggleRowAtIndex(index int) (bool, tea.Cmd) {
	visibleRows := a.GetVisibleRows()
	if index < 0 || index >= len(visibleRows) {
		return false, nil
	}

	row := visibleRows[index]

	switch row.ID {
	case "generator":
		// Cycle to next generator
		a.projectState.CycleToNextGenerator()
		return true, nil

	case "configuration":
		// Toggle configuration
		a.projectState.CycleConfiguration()
		return true, nil

	case "generate":
		// Check if build exists - if so, show confirmation for regenerate
		buildInfo := a.projectState.GetSelectedBuildInfo()
		if buildInfo.Exists {
			// Show confirmation dialog for regenerate - store pending operation (TIT pattern)
			a.pendingOperation = "generate"
			a.confirmDialog = ui.NewConfirmationDialog(
				ui.ConfirmationConfig{
					Title:       "Regenerate " + a.projectState.SelectedGenerator + "?",
					Explanation: "This will overwrite existing build files.",
					YesLabel:    "Yes",
					NoLabel:     "No",
					ActionID:    "generate",
				},
				a.sizing.ContentInnerWidth,
				&a.theme,
			)
			a.confirmDialog.Active = true
			a.confirmDialog.SelectYes()
			return true, nil
		}
		// No confirmation needed for fresh generate
		_, cmd := a.startGenerateOperation()
		return true, cmd

	case "openIde":
		// Open IDE
		_, cmd := a.startOpenIDEOperation()
		return true, cmd

	case "openEditor":
		// Open editor
		_, cmd := a.startOpenEditorOperation()
		return true, cmd

	case "build":
		// Execute build
		_, cmd := a.startBuildOperation()
		return true, cmd

	case "clean":
		// Show confirmation dialog for clean (TIT pattern)
		buildInfo := a.projectState.GetSelectedBuildInfo()
		a.pendingOperation = "clean"
		a.confirmDialog = ui.NewConfirmationDialog(
			ui.ConfirmationConfig{
				Title:       "Clean " + buildInfo.Generator + "?",
				Explanation: "This will delete all build files.",
				YesLabel:    "Yes",
				NoLabel:     "No",
				ActionID:    "clean",
			},
			a.sizing.ContentInnerWidth,
			&a.theme,
		)
		a.confirmDialog.Active = true
		a.confirmDialog.SelectYes()
		return true, nil
	}

	return false, nil
}

// RowByID returns the preference row with the given ID (if visible)
func (a *Application) RowByID(id string) PreferenceRow {
	visibleRows := a.GetVisibleRows()
	for _, row := range visibleRows {
		if row.ID == id {
			return row
		}
	}
	return PreferenceRow{}
}

// RowIndexByID returns the visible index of the row with the given ID
func (a *Application) RowIndexByID(id string) int {
	visibleRows := a.GetVisibleRows()
	for i, row := range visibleRows {
		if row.ID == id {
			return i
		}
	}
	return -1
}

// GeneratePreferencesMenu creates the preferences screen rows
func (a *Application) GeneratePreferencesMenu() []PreferenceRow {
	rows := []PreferenceRow{}

	// === Auto-Scan Section ===
	autoScanEnabled := "OFF"
	if a.config != nil && a.config.IsAutoScanEnabled() {
		autoScanEnabled = "ON"
	}

	rows = append(rows, PreferenceRow{
		ID:       "autoScan",
		Emoji:    "🔄",
		Label:    "Auto-update",
		Value:    autoScanEnabled,
		IsToggle: true,
		IsAction: false,
		Visible:  true,
	})

	interval := "10 min"
	if a.config != nil {
		interval = fmt.Sprintf("%d min", a.config.AutoScanInterval())
	}

	rows = append(rows, PreferenceRow{
		ID:       "scanInterval",
		Emoji:    "⏱️",
		Label:    "Update Interval",
		Value:    interval,
		IsToggle: true,
		IsAction: false,
		Visible:  true,
	})

	// === Separator ===
	rows = append(rows, PreferenceRow{
		ID:        "sep1",
		Emoji:     "",
		Label:     "",
		Value:     "",
		IsToggle:  false,
		IsAction:  false,
		Separator: true,
		Visible:   true,
	})

	// === Theme Section ===
	theme := "gfx"
	if a.config != nil {
		theme = a.config.Theme()
	}

	rows = append(rows, PreferenceRow{
		ID:       "theme",
		Emoji:    "🎨",
		Label:    "Theme",
		Value:    theme,
		IsToggle: true,
		IsAction: false,
		Visible:  true,
	})

	// === Separator ===
	rows = append(rows, PreferenceRow{
		ID:        "sep2",
		Emoji:     "",
		Label:     "",
		Value:     "",
		IsToggle:  false,
		IsAction:  false,
		Separator: true,
		Visible:   true,
	})

	// === Back to Menu ===
	rows = append(rows, PreferenceRow{
		ID:       "back",
		Emoji:    "←",
		Label:    "Back to Menu",
		Value:    "",
		IsToggle: false,
		IsAction: true,
		Visible:  true,
	})

	return rows
}

// GetVisiblePreferenceRows returns only visible preference rows
func (a *Application) GetVisiblePreferenceRows() []PreferenceRow {
	rows := a.GeneratePreferencesMenu()
	visible := []PreferenceRow{}
	for _, row := range rows {
		if row.Visible {
			visible = append(visible, row)
		}
	}
	return visible
}

// TogglePreferenceAtIndex handles preference row toggles
func (a *Application) TogglePreferenceAtIndex(index int) bool {
	visibleRows := a.GetVisiblePreferenceRows()
	if index < 0 || index >= len(visibleRows) {
		return false
	}

	row := visibleRows[index]

	switch row.ID {
	case "autoScan":
		if a.config != nil {
			current := a.config.IsAutoScanEnabled()
			a.config.SetAutoScanEnabled(!current)
		}
		return true

	case "scanInterval":
		// Cycle through interval options: 5, 10, 15, 30 minutes
		if a.config != nil {
			current := a.config.AutoScanInterval()
			var next int
			switch current {
			case 5:
				next = 10
			case 10:
				next = 15
			case 15:
				next = 30
			default:
				next = 10
			}
			a.config.SetAutoScanInterval(next)
		}
		return true

	case "theme":
		// Cycle through available themes
		themes := []string{"gfx", "spring", "summer", "autumn", "winter"}
		currentTheme := "gfx"
		if a.config != nil {
			currentTheme = a.config.Theme()
		}

		currentIndex := -1
		for i, t := range themes {
			if t == currentTheme {
				currentIndex = i
				break
			}
		}

		nextIndex := (currentIndex + 1) % len(themes)
		nextTheme := themes[nextIndex]

		if a.config != nil {
			a.config.SetTheme(nextTheme)

			// Also update the UI theme
			newTheme, err := ui.LoadThemeByName(nextTheme)
			if err == nil {
				a.theme = newTheme
			}
		}
		return true

	case "back":
		a.mode = ModeMenu
		a.selectedIndex = 0
		a.footerHint = FooterHints["menu_navigate"]
		return true
	}

	return false
}
