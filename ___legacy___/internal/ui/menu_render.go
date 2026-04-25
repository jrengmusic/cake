package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderCakeMenu renders cake menu with shortcut column
// Columns: SHORTCUT(10) | EMOJI(3) | LABEL(18) | VALUE(14)
// Shortcut column is wider to accommodate "ctrl + k" style labels
func RenderCakeMenu(rows []MenuRow, selectedIndex int, theme Theme, contentHeight int, contentWidth int) string {
	shortcutColWidth, emojiColWidth, labelColWidth, valueColWidth := calcMenuColWidths(contentWidth)
	menuBoxWidth := shortcutColWidth + emojiColWidth + labelColWidth + valueColWidth

	lines := buildMenuLines(rows, selectedIndex, theme, shortcutColWidth, emojiColWidth, labelColWidth, valueColWidth, menuBoxWidth)

	return assembleMenuOutput(lines, contentHeight, contentWidth, menuBoxWidth)
}

func calcMenuColWidths(contentWidth int) (shortcutColWidth, emojiColWidth, labelColWidth, valueColWidth int) {
	shortcutColWidth = 10
	emojiColWidth = 3
	labelColWidth = 18
	valueColWidth = 14

	if contentWidth < (shortcutColWidth + emojiColWidth + labelColWidth + valueColWidth) {
		labelColWidth = contentWidth - shortcutColWidth - emojiColWidth - valueColWidth
		if labelColWidth < 0 {
			labelColWidth = 0
		}
	}
	return
}

func buildMenuLines(rows []MenuRow, selectedIndex int, theme Theme, shortcutColWidth, emojiColWidth, labelColWidth, valueColWidth, menuBoxWidth int) []string {
	var lines []string
	visibleSelectableIndex := 0

	for _, row := range rows {
		if row.ID == "separator" {
			lines = append(lines, renderMenuSeparator(theme, menuBoxWidth))
			continue
		}

		isSelectable := row.IsSelectable
		isSelected := isSelectable && visibleSelectableIndex == selectedIndex

		shortcutCol := renderShortcutCol(row, shortcutColWidth)
		emojiCol := renderEmojiCol(row.Emoji, emojiColWidth)
		labelCol := renderLabelCol(row.Label, labelColWidth)
		valueCol := renderValueCol(row.Value, valueColWidth)

		styledLine := styleMenuRow(shortcutCol, emojiCol, labelCol, valueCol, isSelected, isSelectable, theme)
		lines = append(lines, styledLine)

		if isSelectable {
			visibleSelectableIndex++
		}
	}
	return lines
}

func renderMenuSeparator(theme Theme, menuBoxWidth int) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.SeparatorColor)).
		Render(strings.Repeat("─", menuBoxWidth))
}

func renderShortcutCol(row MenuRow, colWidth int) string {
	display := row.ShortcutLabel
	if display == "" {
		display = row.Shortcut
	}
	w := lipgloss.Width(display)
	pad := colWidth - w
	if pad < 0 {
		pad = 0
	}
	return strings.Repeat(" ", pad) + display + " "
}

func renderEmojiCol(emoji string, colWidth int) string {
	w := lipgloss.Width(emoji)
	leftPad := (colWidth - w) / 2
	rightPad := colWidth - w - leftPad
	if leftPad < 0 {
		leftPad = 0
	}
	if rightPad < 0 {
		rightPad = 0
	}
	return strings.Repeat(" ", leftPad) + emoji + strings.Repeat(" ", rightPad)
}

func renderLabelCol(label string, colWidth int) string {
	w := lipgloss.Width(label)
	if w > colWidth {
		label = label[:colWidth]
		w = colWidth
	}
	return label + strings.Repeat(" ", colWidth-w)
}

func renderValueCol(value string, colWidth int) string {
	w := lipgloss.Width(value)
	if w > colWidth {
		value = value[:colWidth]
		w = colWidth
	}
	return strings.Repeat(" ", colWidth-w) + value
}

func styleMenuRow(shortcutCol, emojiCol, labelCol, valueCol string, isSelected, isSelectable bool, theme Theme) string {
	if isSelected {
		return styleSelectedMenuRow(shortcutCol, emojiCol, labelCol, valueCol, theme)
	}
	if !isSelectable {
		return styleUnselectableMenuRow(shortcutCol, emojiCol, labelCol, valueCol, theme)
	}
	return styleNormalMenuRow(shortcutCol, emojiCol, labelCol, valueCol, theme)
}

func styleSelectedMenuRow(shortcutCol, emojiCol, labelCol, valueCol string, theme Theme) string {
	shortcutStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.AccentTextColor)).Bold(true)
	emojiStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.LabelTextColor))
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.MainBackgroundColor)).
		Background(lipgloss.Color(theme.MenuSelectionBackground)).
		Bold(true)
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.AccentTextColor)).Bold(true)
	return shortcutStyle.Render(shortcutCol) + emojiStyle.Render(emojiCol) + labelStyle.Render(labelCol) + valueStyle.Render(valueCol)
}

func styleUnselectableMenuRow(shortcutCol, emojiCol, labelCol, valueCol string, theme Theme) string {
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.DimmedTextColor))
	return dimStyle.Render(shortcutCol + emojiCol + labelCol + valueCol)
}

func styleNormalMenuRow(shortcutCol, emojiCol, labelCol, valueCol string, theme Theme) string {
	shortcutStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.AccentTextColor)).Bold(true)
	emojiStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.LabelTextColor))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.LabelTextColor))
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ContentTextColor))
	return shortcutStyle.Render(shortcutCol) + emojiStyle.Render(emojiCol) + labelStyle.Render(labelCol) + valueStyle.Render(valueCol)
}

func assembleMenuOutput(lines []string, contentHeight, contentWidth, menuBoxWidth int) string {
	innerHeight := contentHeight - 2
	menuHeight := len(lines)
	topPad := (innerHeight - menuHeight) / 2
	if topPad < 0 {
		topPad = 0
	}
	bottomPad := innerHeight - menuHeight - topPad
	if bottomPad < 0 {
		bottomPad = 0
	}
	leftPad := (contentWidth - menuBoxWidth) / 2
	if leftPad < 0 {
		leftPad = 0
	}

	var result strings.Builder

	for i := 0; i < topPad; i++ {
		result.WriteString(strings.Repeat(" ", contentWidth))
		if i < topPad-1 || menuHeight > 0 {
			result.WriteString("\n")
		}
	}

	for i, line := range lines {
		centeredLine := strings.Repeat(" ", leftPad) + line
		lineWidth := lipgloss.Width(centeredLine)
		if lineWidth < contentWidth {
			centeredLine = centeredLine + strings.Repeat(" ", contentWidth-lineWidth)
		}
		result.WriteString(centeredLine)
		if i < len(lines)-1 || bottomPad > 0 {
			result.WriteString("\n")
		}
	}

	for i := 0; i < bottomPad; i++ {
		result.WriteString(strings.Repeat(" ", contentWidth))
		if i < bottomPad-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}
