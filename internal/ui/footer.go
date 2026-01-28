package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// FooterShortcut represents a single keyboard shortcut hint
type FooterShortcut struct {
	Key  string // e.g., "↑↓", "Enter", "Esc"
	Desc string // e.g., "scroll", "back", "select"
}

// FooterStyles provides commonly used style objects for footers
type FooterStyles struct {
	shortcutStyle lipgloss.Style
	descStyle     lipgloss.Style
	sepStyle      lipgloss.Style
}

// NewFooterStyles creates a new FooterStyles instance using theme colors
func NewFooterStyles(theme *Theme) FooterStyles {
	return FooterStyles{
		shortcutStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.AccentTextColor)).
			Bold(true),
		descStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.ContentTextColor)),
		sepStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.DimmedTextColor)),
	}
}

// RenderFooter renders footer shortcuts with optional right-side content
// If rightContent is provided: shortcuts on left, rightContent on right
// If rightContent is empty: shortcuts centered (or use hint override)
func RenderFooter(shortcuts []FooterShortcut, width int, theme *Theme, rightContent string) string {
	if theme == nil {
		return ""
	}

	styles := NewFooterStyles(theme)

	// If rightContent provided, render left shortcuts + right content
	if rightContent != "" {
		// Build styled parts: "Key desc"
		var leftParts []string
		for _, sc := range shortcuts {
			part := styles.shortcutStyle.Render(sc.Key) + styles.descStyle.Render(" "+sc.Desc)
			leftParts = append(leftParts, part)
		}
		leftJoined := strings.Join(leftParts, styles.sepStyle.Render("  │  "))

		// Style right content
		rightStyled := styles.descStyle.Render(rightContent)

		// Calculate spacing
		leftWidth := lipgloss.Width(leftJoined)
		rightWidth := lipgloss.Width(rightStyled)
		padding := width - leftWidth - rightWidth
		if padding < 0 {
			padding = 0
		}

		return leftJoined + strings.Repeat(" ", padding) + rightStyled
	}

	// No rightContent: center shortcuts
	var parts []string
	for _, sc := range shortcuts {
		part := styles.shortcutStyle.Render(sc.Key) + styles.descStyle.Render(" "+sc.Desc)
		parts = append(parts, part)
	}

	sep := styles.sepStyle.Render("  │  ")
	content := strings.Join(parts, sep)

	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(content)
}

// RenderFooterOverride renders override message (e.g., Ctrl+C confirm)
// Used for global app-level overrides that take precedence
func RenderFooterOverride(message string, width int, theme *Theme) string {
	if message == "" {
		return ""
	}
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.FooterTextColor)).
		Width(width).
		Align(lipgloss.Center)
	return style.Render(message)
}

// RenderFooterHint renders a simple centered hint (for menu descriptions)
func RenderFooterHint(hint string, width int, theme *Theme) string {
	if hint == "" {
		return ""
	}
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.FooterTextColor)).
		Width(width).
		Align(lipgloss.Center)
	return style.Render(hint)
}
