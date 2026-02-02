package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// MenuRow represents a single menu row
// Fixed 8 rows: [0]Project [1]Regenerate [2]OpenIDE [3]Separator [4]Configuration [5]Build [6]Clean [7]CleanAll
type MenuRow struct {
	ID            string // "project", "regenerate", "openIde", "separator", "configuration", "build", "clean", "cleanAll"
	Shortcut      string // Actual key for handler: "", "g", "o", "", "", "b", "k", "ctrl+k"
	ShortcutLabel string // Display label (right-aligned): "", "g", "o", "", "", "b", "k", "ctrl + k"
	Emoji         string // "⚙️", "🚀", "📂", "", "🏗️", "🔨", "🧹", "💥"
	Label         string // "Project", "Regenerate", "Open IDE", "", "Configuration", "Build", "Clean", "Clean All"
	Value         string // "Xcode", "", "", "", "Debug", "", "", ""
	Visible       bool   // true/false based on conditions
	IsAction      bool   // false for toggles, true for actions
	IsSelectable  bool   // false for separator
	Hint          string // Footer hint/description for this row
}

// GenerateMenuRows returns exactly 8 rows (used by app.go)
// All rows always visible - unavailable options are dimmed and not selectable
func GenerateMenuRows(projectLabel string, configuration string, canOpenIDE bool, canClean bool, hasBuild bool, hasBuildsToClean bool) []MenuRow {
	return []MenuRow{
		{
			ID:            "project",
			Shortcut:      "",
			ShortcutLabel: "",
			Emoji:         "⚙️",
			Label:         "Project",
			Value:         projectLabel,
			Visible:       true,
			IsAction:      false,
			IsSelectable:  true,
			Hint:          "Select project type (Xcode, Ninja, etc.)",
		},
		{
			ID:            "regenerate",
			Shortcut:      "g",
			ShortcutLabel: "g",
			Emoji:         "🚀",
			Label:         map[bool]string{true: "Regenerate", false: "Generate"}[hasBuild],
			Value:         "",
			Visible:       true,
			IsAction:      true,
			IsSelectable:  true,
			Hint:          map[bool]string{true: "Re-run CMake configuration", false: "Run initial CMake configuration"}[hasBuild],
		},
		{
			ID:            "openIde",
			Shortcut:      "o",
			ShortcutLabel: "o",
			Emoji:         "📂",
			Label:         "Open IDE",
			Value:         "",
			Visible:       true,
			IsAction:      true,
			IsSelectable:  canOpenIDE, // Not selectable when unavailable
			Hint:          "Launch IDE for this project",
		},
		{
			ID:            "separator",
			Shortcut:      "",
			ShortcutLabel: "",
			Emoji:         "",
			Label:         "",
			Value:         "",
			Visible:       true,
			IsAction:      false,
			IsSelectable:  false,
			Hint:          "",
		},
		{
			ID:            "configuration",
			Shortcut:      "",
			ShortcutLabel: "",
			Emoji:         "🏗️",
			Label:         "Configuration",
			Value:         configuration,
			Visible:       true,
			IsAction:      false,
			IsSelectable:  true,
			Hint:          "Select build configuration (Debug, Release, etc.)",
		},
		{
			ID:            "build",
			Shortcut:      "b",
			ShortcutLabel: "b",
			Emoji:         "🔨",
			Label:         "Build",
			Value:         "",
			Visible:       true,
			IsAction:      true,
			IsSelectable:  true,
			Hint:          "Build the project",
		},
		{
			ID:            "clean",
			Shortcut:      "k",
			ShortcutLabel: "k",
			Emoji:         "🧹",
			Label:         "Clean",
			Value:         "",
			Visible:       true,
			IsAction:      true,
			IsSelectable:  canClean, // Not selectable when unavailable
			Hint:          "Remove build artifacts for current project",
		},
		{
			ID:            "cleanAll",
			Shortcut:      "ctrl+k",
			ShortcutLabel: "ctrl + k",
			Emoji:         "💥",
			Label:         "Clean All",
			Value:         "",
			Visible:       true,
			IsAction:      true,
			IsSelectable:  hasBuildsToClean, // Not selectable when no builds to clean
			Hint:          "Remove entire Builds/ directory (all projects)",
		},
	}
}

// RenderCakeMenu renders cake menu with shortcut column
// Columns: SHORTCUT(10) | EMOJI(3) | LABEL(18) | VALUE(14)
// Shortcut column is wider to accommodate "ctrl + k" style labels
func RenderCakeMenu(rows []MenuRow, selectedIndex int, theme Theme, contentHeight int, contentWidth int) string {
	// Column widths - shortcut column wider for "ctrl + k" style labels
	shortcutColWidth := 10
	emojiColWidth := 3
	labelColWidth := 18
	valueColWidth := 14

	// Dynamically shrink label if content is too narrow
	if contentWidth < (shortcutColWidth + emojiColWidth + labelColWidth + valueColWidth) {
		labelColWidth = contentWidth - shortcutColWidth - emojiColWidth - valueColWidth
		if labelColWidth < 0 {
			labelColWidth = 0
		}
	}

	menuBoxWidth := shortcutColWidth + emojiColWidth + labelColWidth + valueColWidth

	var lines []string
	visibleSelectableIndex := 0 // Tracks only visible AND selectable items

	for _, row := range rows {
		// Handle separator - render it, but don't count in navigation
		if row.ID == "separator" {
			sepLine := lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.SeparatorColor)).
				Render(strings.Repeat("─", menuBoxWidth))
			lines = append(lines, sepLine)
			continue // Don't increment visibleSelectableIndex
		}

		// This row is visible - check if it's selectable and selected
		isSelectable := row.IsSelectable
		isSelected := isSelectable && visibleSelectableIndex == selectedIndex

		// Column 1: SHORTCUT (right-aligned, using ShortcutLabel for display)
		shortcutDisplay := row.ShortcutLabel
		if shortcutDisplay == "" {
			shortcutDisplay = row.Shortcut
		}
		shortcutW := lipgloss.Width(shortcutDisplay)
		shortcutPad := shortcutColWidth - shortcutW
		if shortcutPad < 0 {
			shortcutPad = 0
		}
		// Right-align the shortcut label
		shortcutCol := strings.Repeat(" ", shortcutPad) + shortcutDisplay + " "

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
			// Selected: highlight label+value, shortcut uses accent color + bold (TIT pattern)
			shortcutStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.AccentTextColor)).
				Bold(true)
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
		} else if !isSelectable {
			// Unselectable (unavailable): dimmed
			dimStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.DimmedTextColor))
			styledLine = dimStyle.Render(shortcutCol + emojiCol + labelCol + valueCol)
		} else {
			// Normal selectable row - shortcut uses accent color + bold (TIT pattern)
			shortcutStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.AccentTextColor)).
				Bold(true)
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
		if isSelectable {
			visibleSelectableIndex++ // Only increment for selectable rows
		}
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
