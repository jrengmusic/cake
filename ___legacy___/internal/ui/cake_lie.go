package ui

import (
	"github.com/jrengmusic/cake/internal/banner"
	"embed"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

//go:embed assets/cake-lie.svg
var cakeLieFS embed.FS

// RenderCakeLieBanner renders the "The cake is a lie" banner with braille art from SVG
// Used when cake is launched outside a CMake project
func RenderCakeLieBanner(width, height int) string {
	logoData, err := cakeLieFS.ReadFile("assets/cake-lie.svg")
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
			output.WriteString("\n")
		}
	}

	return output.String()
}
