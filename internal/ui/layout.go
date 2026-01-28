package ui

import (
	"cake/internal/banner"
	"embed"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

//go:embed assets/cake-logo.svg
var logoFS embed.FS

func RenderBannerDynamic(width, height int) string {
	logoData, err := logoFS.ReadFile("assets/cake-logo.svg")
	if err != nil {
		return strings.Repeat(" ", width) + "\n" +
			strings.Repeat(" ", width) + "\n" +
			strings.Repeat(" ", width)
	}

	svgString := string(logoData)

	canvasWidth := width * 2
	canvasHeight := height * 4

	brailleArray := banner.SvgToBrailleArray(svgString, canvasWidth, canvasHeight)

	var output strings.Builder

	for i, row := range brailleArray {
		for _, bc := range row {
			hex := banner.RGBToHex(bc.Color.R, bc.Color.G, bc.Color.B)
			styledChar := lipgloss.NewStyle().
				Foreground(lipgloss.Color(hex)).
				Render(string(bc.Char))
			output.WriteString(styledChar)
		}
		if i < len(brailleArray)-1 {
			output.WriteString("\n") // Conditional newline
		}
	}

	return output.String()
}

// RenderReactiveLayout combines header/content/footer into full-terminal reactive layout
func RenderReactiveLayout(sizing DynamicSizing, theme Theme, header, content, footer string) string {
	// Too small guard
	if sizing.IsTooSmall {
		return renderTooSmallMessage(sizing, theme)
	}

	contentHeight := sizing.TerminalHeight - HeaderHeight - FooterHeight - 1 // -1 for terminal rendering

	// Header: stick to top, exact height (TIT pattern)
	headerSection := lipgloss.Place(
		sizing.TerminalWidth,
		HeaderHeight,
		lipgloss.Left,
		lipgloss.Top,
		header,
	)

	// Content: fills middle space
	contentSection := lipgloss.NewStyle().
		Width(sizing.TerminalWidth).
		Height(contentHeight).
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(content)

	// Footer: stick to bottom, exact height (TIT pattern)
	footerSection := lipgloss.Place(
		sizing.TerminalWidth,
		FooterHeight,
		lipgloss.Left,
		lipgloss.Top,
		lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.FooterTextColor)).
			Render(footer),
	)

	// Join sections vertically - no wrapping Place (TIT pattern)
	return lipgloss.JoinVertical(lipgloss.Left, headerSection, contentSection, footerSection)
}

func renderTooSmallMessage(sizing DynamicSizing, theme Theme) string {
	msg := "Terminal too small (69×19 minimum)"

	centered := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.DimmedTextColor)).
		Align(lipgloss.Center).
		Render(msg)

	return lipgloss.Place(
		sizing.TerminalWidth,
		sizing.TerminalHeight,
		lipgloss.Center,
		lipgloss.Center,
		centered,
	)
}

// RenderConfirmDialog renders a confirmation dialog overlay
func RenderConfirmDialog(sizing DynamicSizing, theme Theme, header, content, footer string) string {
	// Too small guard
	if sizing.IsTooSmall {
		return renderTooSmallMessage(sizing, theme)
	}

	// Use theme colors for dialog
	dialogBg := lipgloss.NewStyle().
		Background(lipgloss.Color(theme.ConfirmationDialogBackground)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.BoxBorderColor)).
		Padding(2, 4)

	// Button styles using theme colors
	yesButton := lipgloss.NewStyle().
		Background(lipgloss.Color(theme.MenuSelectionBackground)).
		Foreground(lipgloss.Color(theme.ButtonSelectedTextColor)).
		Bold(true).
		Padding(0, 2).
		Render("[Y] YES")

	noButton := lipgloss.NewStyle().
		Background(lipgloss.Color(theme.InlineBackgroundColor)).
		Foreground(lipgloss.Color(theme.ContentTextColor)).
		Bold(true).
		Padding(0, 2).
		Render("[N] NO")

	// Combine message and buttons
	dialogContent := dialogBg.Render(
		content + "\n\n" + yesButton + "  " + noButton,
	)

	// Center the dialog in the content area
	contentHeight := sizing.TerminalHeight - HeaderHeight - FooterHeight

	// Place dialog in center of content area
	centeredDialog := lipgloss.Place(
		sizing.TerminalWidth,
		contentHeight,
		lipgloss.Center,
		lipgloss.Center,
		dialogContent,
	)

	// Header: stick to top, exact height (TIT pattern)
	headerSection := lipgloss.Place(
		sizing.TerminalWidth,
		HeaderHeight,
		lipgloss.Left,
		lipgloss.Top,
		header,
	)

	// Content with dialog overlay
	contentSection := lipgloss.NewStyle().
		Width(sizing.TerminalWidth).
		Height(contentHeight).
		Render(centeredDialog)

	// Footer: single line, centered
	footerSection := lipgloss.NewStyle().
		Width(sizing.TerminalWidth).
		Height(FooterHeight).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color(theme.FooterTextColor)).
		Render(footer)

	// Join sections
	return lipgloss.JoinVertical(lipgloss.Left, headerSection, contentSection, footerSection)
}
