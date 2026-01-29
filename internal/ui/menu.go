package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// MenuRow represents a single menu row
// Fixed 7 rows: [0]Generator [1]Regenerate [2]OpenIDE [3]Separator [4]Configuration [5]Build [6]Clean
type MenuRow struct {
	ID           string // "generator", "regenerate", "openIde", "separator", "configuration", "build", "clean"
	Shortcut     string // "", "g", "o", "", "", "b", "c"
	Emoji        string // "⚙️", "🚀", "📂", "", "🏗️", "🔨", "🧹"
	Label        string // "Generator", "Regenerate", "Open IDE", "", "Configuration", "Build", "Clean"
	Value        string // "Xcode", "", "", "", "Debug", "", ""
	Visible      bool   // true/false based on conditions
	IsAction     bool   // false for toggles, true for actions
	IsSelectable bool   // false for separator
	Hint         string // Footer hint/description for this row
}

// GenerateMenuRows returns exactly 7 rows (used by app.go)
func GenerateMenuRows(projectLabel string, configuration string, canOpenIDE bool, canClean bool, hasBuild bool) []MenuRow {
	return []MenuRow{
		{
			ID:           "project",
			Shortcut:     "",
			Emoji:        "⚙️",
			Label:        "Project",
			Value:        projectLabel,
			Visible:      true,
			IsAction:     false,
			IsSelectable: true,
			Hint:         "Select project type (Xcode, Ninja, etc.)",
		},
		{
			ID:           "regenerate",
			Shortcut:     "g",
			Emoji:        "🚀",
			Label:        map[bool]string{true: "Regenerate", false: "Generate"}[hasBuild],
			Value:        "",
			Visible:      true,
			IsAction:     true,
			IsSelectable: true,
			Hint:         map[bool]string{true: "Re-run CMake configuration", false: "Run initial CMake configuration"}[hasBuild],
		},
		{
			ID:           "openIde",
			Shortcut:     "o",
			Emoji:        "📂",
			Label:        "Open IDE",
			Value:        "",
			Visible:      canOpenIDE,
			IsAction:     true,
			IsSelectable: true,
			Hint:         "Launch IDE for this project",
		},
		{
			ID:           "separator",
			Shortcut:     "",
			Emoji:        "",
			Label:        "",
			Value:        "",
			Visible:      true,
			IsAction:     false,
			IsSelectable: false,
			Hint:         "",
		},
		{
			ID:           "configuration",
			Shortcut:     "",
			Emoji:        "🏗️",
			Label:        "Configuration",
			Value:        configuration,
			Visible:      true,
			IsAction:     false,
			IsSelectable: true,
			Hint:         "Select build configuration (Debug, Release, etc.)",
		},
		{
			ID:           "build",
			Shortcut:     "b",
			Emoji:        "🔨",
			Label:        "Build",
			Value:        "",
			Visible:      true,
			IsAction:     true,
			IsSelectable: true,
			Hint:         "Build the project with selected generator",
		},
		{
			ID:           "clean",
			Shortcut:     "c",
			Emoji:        "🧹",
			Label:        "Clean",
			Value:        "",
			Visible:      canClean,
			IsAction:     true,
			IsSelectable: true,
			Hint:         "Remove build artifacts",
		},
	}
}

// RenderCakeMenu renders cake menu with shortcut column
// Columns: SHORTCUT(3) | EMOJI(3) | LABEL(18) | VALUE(12)
func RenderCakeMenu(rows []MenuRow, selectedIndex int, theme Theme, contentHeight int, contentWidth int) string {
	// Column widths
	const (
		shortcutColWidth = 3
		emojiColWidth    = 3
		labelColWidth    = 18
		valueColWidth    = 14
		menuBoxWidth     = shortcutColWidth + emojiColWidth + labelColWidth + valueColWidth
	)

	var lines []string
	visibleSelectableIndex := 0 // Tracks only visible AND selectable items

	for _, row := range rows {
		// Handle hidden rows - render empty line, don't count in navigation
		if !row.Visible {
			lines = append(lines, strings.Repeat(" ", menuBoxWidth))
			continue // Don't increment visibleSelectableIndex
		}

		// Handle separator - render it, but don't count in navigation
		if row.ID == "separator" {
			sepLine := lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.SeparatorColor)).
				Render(strings.Repeat("─", menuBoxWidth))
			lines = append(lines, sepLine)
			continue // Don't increment visibleSelectableIndex
		}

		// This row is visible and selectable
		// Check if it's selected
		isSelected := visibleSelectableIndex == selectedIndex

		// Column 1: SHORTCUT (left-aligned)
		shortcutStr := row.Shortcut
		shortcutW := lipgloss.Width(shortcutStr)
		shortcutPad := shortcutColWidth - shortcutW
		if shortcutPad < 0 {
			shortcutPad = 0
		}
		shortcutCol := shortcutStr + strings.Repeat(" ", shortcutPad)

		// Column 2: EMOJI (center-aligned)
		emojiStr := row.Emoji
		emojiW := lipgloss.Width(emojiStr)
		emojiLeftPad := (emojiColWidth - emojiW) / 2
		emojiRightPad := emojiColWidth - emojiW - emojiLeftPad
		if emojiLeftPad < 0 {
			emojiLeftPad = 0
		}
		if emojiRightPad < 0 {
			emojiRightPad = 0
		}
		emojiCol := strings.Repeat(" ", emojiLeftPad) + emojiStr + strings.Repeat(" ", emojiRightPad)

		// Column 3: LABEL (left-aligned)
		labelStr := row.Label
		labelW := lipgloss.Width(labelStr)
		if labelW > labelColWidth {
			labelStr = labelStr[:labelColWidth]
			labelW = labelColWidth
		}
		labelCol := labelStr + strings.Repeat(" ", labelColWidth-labelW)

		// Column 4: VALUE (right-aligned)
		valueStr := row.Value
		valueW := lipgloss.Width(valueStr)
		if valueW > valueColWidth {
			valueStr = valueStr[:valueColWidth]
			valueW = valueColWidth
		}
		valueCol := strings.Repeat(" ", valueColWidth-valueW) + valueStr

		// Build styled line
		var styledLine string

		if isSelected {
			// Selected: highlight label+value
			shortcutStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.LabelTextColor))
			emojiStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.LabelTextColor))
			labelStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.MainBackgroundColor)).
				Background(lipgloss.Color(theme.MenuSelectionBackground)).
				Bold(true)
			valueStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.AccentTextColor)).
				Bold(true)

			styledLine = shortcutStyle.Render(shortcutCol) +
				emojiStyle.Render(emojiCol) +
				labelStyle.Render(labelCol) +
				valueStyle.Render(valueCol)
		} else {
			// Normal
			shortcutStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.LabelTextColor))
			emojiStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.LabelTextColor))
			labelStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.LabelTextColor))
			valueStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.ContentTextColor))
			styledLine = shortcutStyle.Render(shortcutCol) +
				emojiStyle.Render(emojiCol) +
				labelStyle.Render(labelCol) +
				valueStyle.Render(valueCol)
		}

		lines = append(lines, styledLine)
		visibleSelectableIndex++ // Only increment for visible, selectable rows
	}

	// Center vertically and horizontally
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

	// Top padding
	for i := 0; i < topPad; i++ {
		result.WriteString(strings.Repeat(" ", contentWidth))
		if i < topPad-1 || menuHeight > 0 {
			result.WriteString("\n")
		}
	}

	// Menu lines
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

	// Bottom padding
	for i := 0; i < bottomPad; i++ {
		result.WriteString(strings.Repeat(" ", contentWidth))
		if i < bottomPad-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}
