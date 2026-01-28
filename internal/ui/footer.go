package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// RenderFooter renders the footer hint with proper styling
func RenderFooter(hint string, theme Theme, width int) string {
	if hint == "" {
		hint = " "
	}

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.FooterTextColor)).
		Background(lipgloss.Color(theme.MainBackgroundColor)).
		Width(width).
		Align(lipgloss.Center)

	return footerStyle.Render(hint)
}
