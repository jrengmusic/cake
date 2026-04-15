package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const spinnerLabelSeparator = "  "
const fallbackOpLabel = "WORKING"
const minContentHeight = 1
const consolePanelHorizontalPadding = 2

var opLabels = map[OpType]string{
	OpBuild:      "BUILDING",
	OpGenerate:   "CONFIGURING",
	OpClean:      "CLEANING",
	OpCleanAll:   "CLEANING ALL",
	OpRegenerate: "REGENERATING",
}

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
	isActive bool,
	spinnerFrame int,
	op OpType,
) string {
	result := ""
	if maxWidth > 0 {
		if totalHeight > 0 {
			consoleHeight := totalHeight
			titleHeight := 1
			contentHeight := consoleHeight - titleHeight
			if contentHeight < minContentHeight {
				contentHeight = minContentHeight
			}
			wrapWidth := maxWidth - consolePanelHorizontalPadding

			state.LinesPerPage = contentHeight

			snapshotLines, totalBufferLines := buffer.GetSnapshot()
			totalDisplayLines := countDisplayLines(snapshotLines, wrapWidth)
			applyScrollState(state, totalDisplayLines, contentHeight, autoScroll)

			visibleLines := formatVisibleLines(snapshotLines, totalBufferLines, palette, wrapWidth, state.ScrollOffset, contentHeight)
			visibleLines = padLinesToWidth(visibleLines, wrapWidth)
			visibleLines = padLinesToHeight(visibleLines, contentHeight, wrapWidth)

			panel := assembleConsolePanel(visibleLines, palette, wrapWidth, consoleHeight, isActive, spinnerFrame, op)
			result = lipgloss.NewStyle().Padding(0, 1).Render(panel)
		}
	}
	return result
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

// countEntryDisplayLines returns the number of display lines an entry occupies at wrapWidth.
// Cheap — no rendering, just visual width calculation.
func countEntryDisplayLines(line OutputLine, wrapWidth int) int {
	formatted := fmt.Sprintf("[%s] %s", line.Time, line.Text)
	visualWidth := lipgloss.Width(formatted)
	displayLines := (visualWidth + wrapWidth - 1) / wrapWidth
	if displayLines < 1 {
		displayLines = 1
	}
	return displayLines
}

// countDisplayLines walks snapshot entries and returns total display line count.
// Cheap — no rendering, just visual width calculation.
func countDisplayLines(snapshotLines []OutputLine, wrapWidth int) int {
	totalDisplayLines := 0
	for _, line := range snapshotLines {
		totalDisplayLines += countEntryDisplayLines(line, wrapWidth)
	}
	return totalDisplayLines
}

// renderEntry renders a single buffer entry and splits it into display lines.
func renderEntry(line OutputLine, palette Theme, wrapWidth int) []string {
	colorMap := consoleLineColorMap(palette)
	color := colorMap[line.Type]
	if color == "" {
		color = palette.OutputStdoutColor
	}
	lineStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	formatted := fmt.Sprintf("[%s] %s", line.Time, line.Text)
	renderedLine := lineStyle.Width(wrapWidth).Render(formatted)
	return strings.Split(renderedLine, "\n")
}

// collectVisibleFromEntry takes rendered display lines for one entry and appends
// those that fall within [scrollOffset, scrollOffset+contentHeight), given the
// entry starts at displayLineBase in the global display-line coordinate space.
// Returns updated collected slice and updated accumulator.
func collectVisibleFromEntry(
	collected []string,
	entryLines []string,
	displayLineBase int,
	scrollOffset int,
	contentHeight int,
) []string {
	for localIdx, displayLine := range entryLines {
		globalIdx := displayLineBase + localIdx
		withinWindow := globalIdx >= scrollOffset && len(collected) < contentHeight
		if withinWindow {
			collected = append(collected, displayLine)
		}
	}
	return collected
}

// formatVisibleLines renders ONLY the buffer entries that intersect the visible window.
// Returns at most contentHeight display lines (or fewer if buffer is small).
// Preserves empty-buffer "(no output yet)" behavior when totalBufferLines == 0.
func formatVisibleLines(
	snapshotLines []OutputLine,
	totalBufferLines int,
	palette Theme,
	wrapWidth int,
	scrollOffset int,
	contentHeight int,
) []string {
	var visibleLines []string

	if totalBufferLines == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette.DimmedTextColor)).
			Italic(true)
		visibleLines = append(visibleLines, emptyStyle.Render("(no output yet)"))
	} else {
		displayLineAccumulator := 0
		for _, line := range snapshotLines {
			entryLineCount := countEntryDisplayLines(line, wrapWidth)
			entryEndLine := displayLineAccumulator + entryLineCount
			entryInWindow := entryEndLine > scrollOffset && len(visibleLines) < contentHeight
			if entryInWindow {
				entryLines := renderEntry(line, palette, wrapWidth)
				visibleLines = collectVisibleFromEntry(visibleLines, entryLines, displayLineAccumulator, scrollOffset, contentHeight)
			}
			displayLineAccumulator += entryLineCount
		}
	}

	return visibleLines
}

func clampScrollOffset(scrollOffset int, maxScroll int) int {
	clamped := scrollOffset
	if clamped > maxScroll {
		clamped = maxScroll
	}
	if clamped < 0 {
		clamped = 0
	}
	return clamped
}

func applyScrollState(state *ConsoleOutState, totalDisplayLines int, contentHeight int, autoScroll bool) {
	maxScroll := totalDisplayLines - contentHeight
	if maxScroll < 0 {
		maxScroll = 0
	}
	state.MaxScroll = maxScroll

	if autoScroll {
		state.ScrollOffset = maxScroll
	} else {
		state.ScrollOffset = clampScrollOffset(state.ScrollOffset, maxScroll)
	}
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

func assembleConsolePanel(visibleLines []string, palette Theme, wrapWidth int, consoleHeight int, isActive bool, spinnerFrame int, op OpType) string {
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(palette.OutputInfoColor)).
		Bold(true)

	title := buildConsoleTitle(labelStyle, palette, wrapWidth, isActive, spinnerFrame, op)

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

func buildConsoleTitle(labelStyle lipgloss.Style, palette Theme, wrapWidth int, isActive bool, spinnerFrame int, op OpType) string {
	var title string
	if isActive {
		spinnerStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette.SpinnerColor)).
			Bold(true)
		frame := spinnerStyle.Render(GetSpinnerFrame(spinnerFrame))
		label, ok := opLabels[op]
		if !ok {
			label = fallbackOpLabel
		}
		title = frame + spinnerLabelSeparator + labelStyle.Render(label+" ...")
	} else {
		title = labelStyle.Render("OUTPUT")
	}
	titleWidth := lipgloss.Width(title)
	if titleWidth < wrapWidth {
		title = title + strings.Repeat(" ", wrapWidth-titleWidth)
	}
	return title
}
