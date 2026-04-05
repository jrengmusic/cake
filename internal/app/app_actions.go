package app

import (
	"github.com/jrengmusic/cake/internal/ui"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

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

// GetVisiblePreferenceRows returns visible preference rows for preferences mode
func (a *Application) GetVisiblePreferenceRows() []ui.MenuRow {
	if a.config == nil {
		return []ui.MenuRow{}
	}

	autoScanValue := "OFF"
	if a.config.IsAutoScanEnabled() {
		autoScanValue = "ON"
	}

	return []ui.MenuRow{
		{
			ID:           "prefs_auto_scan",
			Shortcut:     "",
			Emoji:        "🔄",
			Label:        "Auto-scan",
			Value:        autoScanValue,
			Visible:      true,
			IsAction:     false,
			IsSelectable: true,
			Hint:         "Toggle automatic project scanning",
		},
		{
			ID:           "prefs_interval",
			Shortcut:     "",
			Emoji:        "⏱️",
			Label:        "Scan Interval",
			Value:        fmt.Sprintf("%d min", a.config.AutoScanInterval()),
			Visible:      true,
			IsAction:     false,
			IsSelectable: true,
			Hint:         "Adjust auto-scan interval (+/- 1min, =/_ 10min)",
		},
		{
			ID:           "prefs_theme",
			Shortcut:     "",
			Emoji:        "🎨",
			Label:        "Theme",
			Value:        a.config.Theme(),
			Visible:      true,
			IsAction:     false,
			IsSelectable: true,
			Hint:         "Cycle through available themes",
		},
	}
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

// TogglePreferenceAtIndex toggles preference at given VISIBLE index
// Returns true if preference was toggled successfully
func (a *Application) applyNextTheme() bool {
	nextTheme, err := ui.GetNextTheme(a.config.Theme())
	if err != nil {
		a.footerHint = fmt.Sprintf("Failed to get next theme: %v", err)
		return false
	}
	if err := a.config.SetTheme(nextTheme); err != nil {
		a.footerHint = fmt.Sprintf("Failed to save theme: %v", err)
		return false
	}
	newTheme, err := ui.LoadThemeByName(nextTheme)
	if err != nil {
		a.footerHint = fmt.Sprintf("Failed to load theme: %v", err)
		return false
	}
	a.theme = newTheme
	return true
}

func (a *Application) TogglePreferenceAtIndex(visibleIndex int) bool {
	visibleRows := a.GetVisiblePreferenceRows()
	if visibleIndex < 0 || visibleIndex >= len(visibleRows) {
		return false
	}
	row := visibleRows[visibleIndex]
	switch row.ID {
	case "prefs_auto_scan":
		newValue := !a.config.IsAutoScanEnabled()
		if err := a.config.SetAutoScanEnabled(newValue); err != nil {
			a.footerHint = fmt.Sprintf("Failed to save config: %v", err)
			return false
		}
		return true
	case "prefs_theme":
		return a.applyNextTheme()
	case "prefs_interval":
		return true
	}
	return false
}

// showConfirmationDialog creates and activates a confirmation dialog for a pending operation
func (a *Application) showConfirmationDialog(title, explanation, yesLabel, noLabel, actionID string) {
	a.confirmDialog = ui.NewConfirmationDialogWithDefault(ui.ConfirmationConfig{
		Title:       title,
		Explanation: explanation,
		YesLabel:    yesLabel,
		NoLabel:     noLabel,
		ActionID:    actionID,
	}, a.sizing.ContentInnerWidth, &a.theme, ui.ButtonNo)
	a.confirmDialog.Active = true
	a.pendingOperation = actionID
}

func (a *Application) executeRowActionRegenerate() (bool, tea.Cmd) {
	buildInfo := a.projectState.GetSelectedBuildInfo()
	if !buildInfo.Exists {
		_, cmd := a.startGenerateOperation()
		return true, cmd
	}
	a.showConfirmationDialog("Regenerate Project", "Clean and re-run CMake configuration?", "Yes", "No", "regenerate")
	return true, nil
}

func (a *Application) executeRowActionBuild() (bool, tea.Cmd) {
	if !a.projectState.CanBuild() {
		a.buildAfterGenerate = true
		_, cmd := a.startGenerateOperation()
		return true, cmd
	}
	_, cmd := a.startBuildOperation()
	return true, cmd
}

func (a *Application) showCleanConfirmDialog() {
	a.showConfirmationDialog(
		"Clean Build Directory",
		"Remove all build artifacts for "+a.projectState.SelectedProject+"?",
		"Yes", "No", "clean",
	)
}

func (a *Application) showCleanAllConfirmDialog() {
	a.showConfirmationDialog(
		"Clean All Projects",
		"This will permanently delete the entire 'Builds/' directory, removing ALL build artifacts for ALL projects. This action cannot be undone.",
		"Yes, Delete All", "Cancel", "cleanAll",
	)
}

// executeRowAction executes the action associated with a menu row
func (a *Application) executeRowAction(rowID string) (bool, tea.Cmd) {
	switch rowID {
	case "project":
		a.projectState.CycleToNextProject()
		a.menuItems = a.GenerateMenu()
		return true, nil
	case "regenerate":
		return a.executeRowActionRegenerate()
	case "clean":
		a.showCleanConfirmDialog()
		return true, nil
	case "cleanAll":
		a.showCleanAllConfirmDialog()
		return true, nil
	case "openIde":
		_, cmd := a.startOpenIDEOperation()
		return true, cmd
	case "configuration":
		a.projectState.CycleConfiguration()
		a.menuItems = a.GenerateMenu()
		return true, nil
	case "build":
		return a.executeRowActionBuild()
	}
	return false, nil
}
