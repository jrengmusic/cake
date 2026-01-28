package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const EmojiColumnWidth = 3

// HeaderState for CAKE - simplified from TIT (3-line format)
type HeaderState struct {
	ProjectName  string
	CWD          string
	Version      string
	VersionColor string
}

// RenderHeaderInfo renders CAKE header info (3 lines exactly)
// CAKE format:
//
//	Line 1: "    PROJECT" (4 spaces, uppercase project name)
//	Line 2: "📂  CWD"     (emoji, 2 spaces, path)
//	Line 3: separator
func RenderHeaderInfo(sizing DynamicSizing, theme Theme, state HeaderState) string {
	totalWidth := sizing.HeaderInnerWidth

	// === LINE 1: Project Name ===
	// "    PROJECT" - 4 spaces for alignment, uppercase
	projectLine := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(theme.LabelTextColor)).
		Render("    " + strings.ToUpper(state.ProjectName))

	// === LINE 2: CWD with emoji ===
	cwdLine := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(theme.CwdTextColor)).
		Render("📂  " + state.CWD)

	// === LINE 3: Separator ===
	separatorLine := lipgloss.NewStyle().
		Width(totalWidth).
		Foreground(lipgloss.Color(theme.SeparatorColor)).
		Render(strings.Repeat("─", totalWidth))

	// Combine all 3 lines
	result := lipgloss.JoinVertical(lipgloss.Left, projectLine, cwdLine, separatorLine)

	return result
}

// RenderHeader renders header with margins (TIT pattern exactly)
func RenderHeader(sizing DynamicSizing, theme Theme, info string) string {
	marginStyle := lipgloss.NewStyle().
		PaddingLeft(HorizontalMargin).
		PaddingRight(HorizontalMargin)

	infoStyled := lipgloss.NewStyle().
		Width(sizing.HeaderInnerWidth).
		AlignVertical(lipgloss.Top).
		Render(info)

	return marginStyle.Render(infoStyled)
}
