package app

import (
	"github.com/jrengmusic/cake/internal/ui"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (a *Application) adjustAutoScanInterval(delta int) {
	newInterval := a.config.AutoScanInterval() + delta
	if newInterval > MaxAutoScanInterval {
		newInterval = MaxAutoScanInterval
	}
	if newInterval < MinAutoScanInterval {
		newInterval = MinAutoScanInterval
	}
	if err := a.config.SetAutoScanInterval(newInterval); err != nil {
		a.footerHint = fmt.Sprintf("Failed to save config: %v", err)
	}
}

func (a *Application) movePreferenceSelectionUp(visibleRows []ui.MenuRow) {
	if a.selectedIndex > 0 {
		nextIndex := a.selectedIndex - 1
		for nextIndex >= 0 && visibleRows[nextIndex].ID == "separator" {
			nextIndex--
		}
		if nextIndex >= 0 {
			a.selectedIndex = nextIndex
		}
	}
}

func (a *Application) movePreferenceSelectionDown(visibleRows []ui.MenuRow) {
	visibleCount := len(visibleRows)
	if a.selectedIndex < visibleCount-1 {
		nextIndex := a.selectedIndex + 1
		for nextIndex < visibleCount && visibleRows[nextIndex].ID == "separator" {
			nextIndex++
		}
		if nextIndex < visibleCount {
			a.selectedIndex = nextIndex
		}
	}
}

func (a *Application) handlePreferencesIntervalKey(key string, visibleRows []ui.MenuRow) {
	if visibleRows[a.selectedIndex].ID != "prefs_interval" {
		return
	}
	intervalKeyMap := map[string]int{
		"+": 1, "=": 1,
		"-": -1, "_": -1,
		"shift++": AutoScanIntervalStep, "shift+=": AutoScanIntervalStep,
		"shift+-": -AutoScanIntervalStep, "shift+_": -AutoScanIntervalStep,
	}
	if delta, ok := intervalKeyMap[key]; ok {
		a.adjustAutoScanInterval(delta)
	}
}

func (a *Application) handlePreferencesKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	a.lastActivityTime = time.Now()
	visibleRows := a.GetVisiblePreferenceRows()

	switch msg.String() {
	case "up", "k":
		a.movePreferenceSelectionUp(visibleRows)
	case "down", "j":
		a.movePreferenceSelectionDown(visibleRows)
	case "enter", " ":
		a.TogglePreferenceAtIndex(a.selectedIndex)
	case "+", "=", "-", "_", "shift++", "shift+=", "shift+-", "shift+_":
		a.handlePreferencesIntervalKey(msg.String(), visibleRows)
	case "/", "esc":
		a.mode = ModeMenu
		a.selectedIndex = 0
		a.footerHint = FooterHints["menu_navigate"]
	case "ctrl+c":
		return a.handleCtrlC()
	}
	return a, nil
}

func (a *Application) abortActiveOperation() {
	if a.cancelContext != nil {
		a.cancelContext()
	}
	if a.killTree != nil {
		a.killTree()
		a.killTree = nil
	}
	a.asyncState.Abort()
	a.outputBuffer.Append("", ui.TypeStdout)
	a.outputBuffer.Append("Operation aborted by user", ui.TypeStderr)
	a.outputBuffer.Append("Press ESC to return to menu", ui.TypeInfo)
}

func (a *Application) returnToMenuFromConsole() {
	a.mode = ModeMenu
	a.selectedIndex = 0
	a.menuItems = a.GenerateMenu()
	a.footerHint = FooterHints["menu_navigate"]
}

func (a *Application) handleOperationKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	a.lastActivityTime = time.Now()
	switch msg.String() {
	case "up":
		a.consoleState.ScrollUp()
		a.consoleAutoScroll = false
		return a, nil
	case "down":
		a.consoleState.ScrollDown()
		a.consoleAutoScroll = false
		return a, nil
	case "esc":
		if a.asyncState.IsActive() {
			a.abortActiveOperation()
		} else {
			a.returnToMenuFromConsole()
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

// handleInvalidProjectKeyPress handles keys in invalid project mode (no CMakeLists.txt)
func (a *Application) handleInvalidProjectKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "Q":
		return a, tea.Quit
	case "esc":
		return a, tea.Quit
	case "ctrl+c":
		return a.handleCtrlC()
	}
	return a, nil
}
