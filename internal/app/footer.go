package app

import (
	"cake/internal/ui"
	"fmt"
)

// GetFooterContent returns the rendered footer for current mode.
// Priority: quitConfirm > mode-specific hints
// Returns styled string ready for display.
func (a *Application) GetFooterContent() string {
	width := a.sizing.TerminalWidth

	// Priority 1: Quit confirmation (Ctrl+C) - GLOBAL APP OVERRIDE
	if a.quitConfirmActive {
		return ui.RenderFooterOverride(
			GetFooterMessageText(MessageCtrlCConfirm),
			width,
			&a.theme,
		)
	}

	// Priority 2: Mode-specific footer
	switch a.mode {
	case ModeMenu:
		// Menu mode: show selected menu item's hint/description
		return a.getMenuFooter(width)

	case ModeConsole:
		// Console mode: show scroll shortcuts + scroll status (left/right)
		return a.getConsoleFooter(width)

	case ModePreferences:
		// Preferences mode: navigation shortcuts
		shortcuts := FooterHintShortcuts["preferences"]
		return ui.RenderFooter(shortcuts, width, &a.theme, "")

	default:
		return ""
	}
}

// getMenuFooter returns footer for menu mode (selected item's hint)
func (a *Application) getMenuFooter(width int) string {
	visibleRows := a.GetVisibleRows()
	if a.selectedIndex < 0 || a.selectedIndex >= len(visibleRows) {
		return ""
	}

	selectedRow := visibleRows[a.selectedIndex]

	// If row has a hint, show it
	if selectedRow.Hint != "" {
		return ui.RenderFooterHint(selectedRow.Hint, width, &a.theme)
	}

	return ""
}

// getConsoleFooter returns footer for console mode (shortcuts + scroll status)
func (a *Application) getConsoleFooter(width int) string {
	// Determine which shortcut set to use
	var hintKey string
	if a.asyncState.IsActive() {
		hintKey = "console_running"
	} else {
		hintKey = "console_complete"
	}

	shortcuts := FooterHintShortcuts[hintKey]
	rightContent := a.computeConsoleScrollStatus()

	return ui.RenderFooter(shortcuts, width, &a.theme, rightContent)
}

// computeConsoleScrollStatus returns the right-side scroll status for console mode
func (a *Application) computeConsoleScrollStatus() string {
	state := &a.consoleState

	// Handle case where MaxScroll is 0 or negative
	if state.MaxScroll <= 0 {
		return ""
	}

	atBottom := state.ScrollOffset >= state.MaxScroll
	remainingLines := state.MaxScroll - state.ScrollOffset

	if atBottom {
		return "(at bottom)"
	}
	if remainingLines > 0 {
		return fmt.Sprintf("↓ %d more", remainingLines)
	}
	return "(can scroll up)"
}
