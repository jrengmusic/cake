package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ConsoleOutState holds the scrolling state for console output
type ConsoleOutState struct {
	ScrollOffset int
	LinesPerPage int
	MaxScroll    int // Cached max scroll position
}

// NewConsoleOutState creates a new console output state with default values
func NewConsoleOutState() ConsoleOutState {
	return ConsoleOutState{
		ScrollOffset: 0,
		LinesPerPage: 18, // Default for content area
	}
}

// Reset resets the scroll state
func (s *ConsoleOutState) Reset() {
	s.ScrollOffset = 0
	s.MaxScroll = 0
}

// ScrollUp moves the viewport up by one line
func (s *ConsoleOutState) ScrollUp() {
	if s.ScrollOffset > 0 {
		s.ScrollOffset--
	}
}

// ScrollDown moves the viewport down by one line
func (s *ConsoleOutState) ScrollDown() {
	if s.ScrollOffset < s.MaxScroll {
		s.ScrollOffset++
	}
}

// RenderConsoleOutput renders console output for full-screen mode (footer handled externally)
// Takes terminal dimensions directly, returns content that occupies full terminal
func RenderConsoleOutput(
	state *ConsoleOutState,
	buffer *OutputBuffer,
	palette Theme,
	maxWidth int,
	totalHeight int,
	operationInProgress bool,
	abortConfirmActive bool,
	autoScroll bool,
) string {
	if maxWidth <= 0 || totalHeight <= 0 {
		return ""
	}

	consoleHeight := totalHeight
	titleHeight := 1
	contentHeight := consoleHeight - titleHeight
	if contentHeight < 1 {
		contentHeight = 1
	}
	wrapWidth := maxWidth - 2

	state.LinesPerPage = contentHeight

	allOutputLines := formatBufferLines(buffer, palette, wrapWidth)
	applyScrollState(state, allOutputLines, contentHeight, autoScroll)

	visibleLines := extractVisibleWindow(allOutputLines, state.ScrollOffset, contentHeight)
	visibleLines = padLinesToWidth(visibleLines, wrapWidth)
	visibleLines = padLinesToHeight(visibleLines, contentHeight, wrapWidth)

	panel := assembleConsolePanel(visibleLines, palette, wrapWidth, consoleHeight)
	return lipgloss.NewStyle().Padding(0, 1).Render(panel)
}

func consoleLineColorMap(palette Theme) map[OutputLineType]string {
	return map[OutputLineType]string{
		TypeStdout:  palette.OutputStdoutColor,
		TypeStderr:  palette.OutputStderrColor,
		TypeCommand: palette.OutputStdoutColor,
		TypeStatus:  palette.OutputStatusColor,
		TypeWarning: palette.OutputWarningColor,
		TypeDebug:   palette.OutputDebugColor,
		TypeInfo:    palette.OutputInfoColor,
	}
}

func formatBufferLines(buffer *OutputBuffer, palette Theme, wrapWidth int) []string {
	snapshotLines, totalBufferLines := buffer.GetSnapshot()

	var allOutputLines []string

	if totalBufferLines == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette.DimmedTextColor)).
			Italic(true)
		allOutputLines = append(allOutputLines, emptyStyle.Render("(no output yet)"))
		return allOutputLines
	}

	colorMap := consoleLineColorMap(palette)
	for _, line := range snapshotLines {
		color := colorMap[line.Type]
		if color == "" {
			color = palette.OutputStdoutColor
		}
		lineStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
		formatted := fmt.Sprintf("[%s] %s", line.Time, line.Text)
		renderedLine := lineStyle.Width(wrapWidth).Render(formatted)
		allOutputLines = append(allOutputLines, strings.Split(renderedLine, "\n")...)
	}
	return allOutputLines
}

func applyScrollState(state *ConsoleOutState, allOutputLines []string, contentHeight int, autoScroll bool) {
	totalOutputLines := len(allOutputLines)
	maxScroll := totalOutputLines - contentHeight
	if maxScroll < 0 {
		maxScroll = 0
	}
	state.MaxScroll = maxScroll

	if autoScroll {
		state.ScrollOffset = maxScroll
		return
	}

	if state.ScrollOffset > maxScroll {
		state.ScrollOffset = maxScroll
	}
	if state.ScrollOffset < 0 {
		state.ScrollOffset = 0
	}
}

func extractVisibleWindow(allOutputLines []string, scrollOffset int, contentHeight int) []string {
	totalOutputLines := len(allOutputLines)
	start := scrollOffset
	end := start + contentHeight
	if start < 0 {
		start = 0
	}
	if end > totalOutputLines {
		end = totalOutputLines
	}

	var visibleLines []string
	for i := start; i < end; i++ {
		visibleLines = append(visibleLines, allOutputLines[i])
	}
	return visibleLines
}

func padLinesToWidth(lines []string, wrapWidth int) []string {
	for i := range lines {
		lineWidth := lipgloss.Width(lines[i])
		if lineWidth < wrapWidth {
			lines[i] = lines[i] + strings.Repeat(" ", wrapWidth-lineWidth)
		}
	}
	return lines
}

func padLinesToHeight(lines []string, contentHeight int, wrapWidth int) []string {
	emptyLine := strings.Repeat(" ", wrapWidth)
	for len(lines) < contentHeight {
		lines = append(lines, emptyLine)
	}
	return lines
}

func assembleConsolePanel(visibleLines []string, palette Theme, wrapWidth int, consoleHeight int) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(palette.OutputInfoColor)).
		Bold(true)

	title := titleStyle.Render("OUTPUT")
	titleWidth := lipgloss.Width(title)
	if titleWidth < wrapWidth {
		title = title + strings.Repeat(" ", wrapWidth-titleWidth)
	}

	blankLine := strings.Repeat(" ", wrapWidth)
	contentBox := strings.Join(visibleLines, "\n")

	panel := lipgloss.JoinVertical(lipgloss.Left, title, contentBox)

	panelLines := strings.Split(panel, "\n")
	for len(panelLines) < consoleHeight {
		panelLines = append(panelLines, blankLine)
	}
	if len(panelLines) > consoleHeight {
		panelLines = panelLines[:consoleHeight]
	}
	return strings.Join(panelLines, "\n")
}
