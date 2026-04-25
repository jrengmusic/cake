package app

import (
	"github.com/jrengmusic/cake/internal/ui"

	"github.com/charmbracelet/lipgloss"
)

// renderMenuWithBanner renders menu (left 65%) + banner (right 35%)
func (a *Application) renderMenuWithBanner() string {
	// 65/35 split — menu needs more room for 2-column rows
	leftWidth := a.sizing.ContentInnerWidth * MenuBannerSplitPct / 100
	rightWidth := a.sizing.ContentInnerWidth - leftWidth

	// Render menu in left column using new ui.RenderCakeMenu
	menuContent := ui.RenderCakeMenu(a.menuItems, a.selectedIndex, a.theme, a.sizing.ContentHeight, leftWidth)

	menuColumn := lipgloss.NewStyle().
		Width(leftWidth).
		Height(a.sizing.ContentHeight).
		Align(lipgloss.Left).
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

// renderPreferencesWithBanner renders preferences (left 50%) + banner (right 50%)
func (a *Application) renderPreferencesWithBanner() string {
	return ui.RenderPreferencesWithBanner(a.config, a.selectedIndex, a.theme, a.sizing)
}

func (a *Application) renderConsoleMode() string {
	// Console height accounts for footer
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
		a.asyncState.IsActive(),
		a.spinnerFrame,
		a.asyncState.CurrentOp(),
	)

	footerText := a.GetFooterContent()

	// Footer only - console handles its own OUTPUT title
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
