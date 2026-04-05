package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// menuShortcutMap maps keyboard shortcuts to their menu row IDs
var menuShortcutMap = map[string]string{
	"g":      "regenerate",
	"G":      "regenerate",
	"o":      "openIde",
	"O":      "openIde",
	"b":      "build",
	"B":      "build",
	"c":      "clean",
	"C":      "clean",
	"x":      "cleanAll",
	"X":      "cleanAll",
}

// executeMenuShortcut jumps to a menu row by ID and executes its action
func (a *Application) executeMenuShortcut(rowID string) (tea.Model, tea.Cmd) {
	idx := a.GetVisibleIndex(rowID)
	if idx >= 0 {
		a.selectedIndex = idx
		handled, cmd := a.ToggleRowAtIndex(idx)
		if handled {
			a.menuItems = a.GenerateMenu()
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
}

// executePendingOperation executes the operation stored in pendingOperation
func (a *Application) executePendingOperation() (tea.Model, tea.Cmd) {
	op := a.pendingOperation
	a.pendingOperation = ""

	switch op {
	case "generate":
		return a.startGenerateOperation()
	case "clean":
		return a.startCleanOperation()
	case "cleanAll":
		return a.startCleanAllOperation()
	case "regenerate":
		return a.startRegenerateOperation()
	}
	return a, nil
}

func (a *Application) cancelConfirmDialog() {
	a.confirmDialog.Active = false
	a.confirmDialog = nil
	a.pendingOperation = ""
}

func (a *Application) confirmAndExecute() (tea.Model, tea.Cmd) {
	a.confirmDialog.Active = false
	if a.pendingOperation != "" {
		return a.executePendingOperation()
	}
	return a, nil
}

func (a *Application) handleConfirmDialogEnter() (tea.Model, tea.Cmd) {
	if a.confirmDialog.GetSelectedButton() == ButtonYes {
		return a.confirmAndExecute()
	}
	a.cancelConfirmDialog()
	return a, nil
}

func (a *Application) handleConfirmDialogCtrlC() (tea.Model, tea.Cmd) {
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

// handleConfirmDialogKeyPress handles Y/N keys for confirmation dialogs
func (a *Application) handleConfirmDialogKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if a.confirmDialog == nil {
		return a, nil
	}
	switch msg.String() {
	case "ctrl+c":
		return a.handleConfirmDialogCtrlC()
	case "y", "Y":
		return a.confirmAndExecute()
	case "n", "N":
		a.cancelConfirmDialog()
		return a, nil
	case "enter", " ":
		return a.handleConfirmDialogEnter()
	case "esc":
		a.cancelConfirmDialog()
		return a, nil
	case "left", "h":
		a.confirmDialog.SelectYes()
		return a, nil
	case "right", "l":
		a.confirmDialog.SelectNo()
		return a, nil
	default:
		return a, nil
	}
}

func (a *Application) clampSelectedIndex() {
	newVisible := a.GetVisibleRows()
	if a.selectedIndex >= len(newVisible) {
		a.selectedIndex = len(newVisible) - 1
	}
	if a.selectedIndex < 0 {
		a.selectedIndex = 0
	}
}

func (a *Application) executeMenuSelection(visibleCount int) (tea.Model, tea.Cmd) {
	if a.selectedIndex >= 0 && a.selectedIndex < visibleCount {
		handled, cmd := a.ToggleRowAtIndex(a.selectedIndex)
		if handled {
			a.menuItems = a.GenerateMenu()
			a.clampSelectedIndex()
			return a, cmd
		}
	}
	return a, nil
}

func clampToRange(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func (a *Application) togglePreferencesMode() {
	if a.mode == ModeMenu {
		a.mode = ModePreferences
		a.selectedIndex = 0
	} else if a.mode == ModePreferences {
		a.mode = ModeMenu
		a.selectedIndex = 0
	}
}

func (a *Application) handleMenuKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	a.lastActivityTime = time.Now()
	visibleCount := len(a.GetVisibleRows())

	switch msg.String() {
	case "up", "k":
		if a.selectedIndex > 0 {
			a.selectedIndex = clampToRange(a.selectedIndex-1, 0, visibleCount-1)
		}
		return a, nil
	case "down", "j":
		if a.selectedIndex < visibleCount-1 {
			a.selectedIndex = clampToRange(a.selectedIndex+1, 0, visibleCount-1)
		}
		return a, nil
	case "enter", " ":
		return a.executeMenuSelection(visibleCount)
	case "/":
		a.togglePreferencesMode()
		return a, nil
	case "ctrl+c":
		return a.handleCtrlC()
	default:
		if rowID, ok := menuShortcutMap[msg.String()]; ok {
			return a.executeMenuShortcut(rowID)
		}
	}
	return a, nil
}

