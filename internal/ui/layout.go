package ui

import (
	"github.com/jrengmusic/cake/internal/banner"
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

	// Header: stick to top, exact height
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

	// Footer: stick to bottom, exact height
	footerSection := lipgloss.Place(
		sizing.TerminalWidth,
		FooterHeight,
		lipgloss.Left,
		lipgloss.Top,
		lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.FooterTextColor)).
			Render(footer),
	)

	// Join sections vertically - no wrapping Place
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
