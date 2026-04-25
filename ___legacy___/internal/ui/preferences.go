package ui

import (
	"fmt"

	"github.com/jrengmusic/cake/internal/config"

	"github.com/charmbracelet/lipgloss"
)

// PreferenceRow represents a single preference row with label and value
type PreferenceRow struct {
	Emoji   string
	Label   string
	Value   string
	Enabled bool
}

// BuildPreferenceRows builds preference rows from config
// Returns rows matching menu item order for consistent rendering
func BuildPreferenceRows(cfg *config.Config) []PreferenceRow {
	if cfg == nil {
		return []PreferenceRow{}
	}

	autoUpdateValue := "OFF"
	if cfg.IsAutoScanEnabled() {
		autoUpdateValue = "ON"
	}

	return []PreferenceRow{
		{Emoji: "🔄", Label: "Auto-scan", Value: autoUpdateValue, Enabled: true},
		{Emoji: "⏱️", Label: "Scan Interval", Value: fmt.Sprintf("%d min", cfg.AutoScanInterval()), Enabled: true},
		{Emoji: "🎨", Label: "Theme", Value: cfg.Theme(), Enabled: true},
	}
}

// RenderPreferencesMenu renders preference rows as EMOJI | LABEL | VALUE
// No shortcut column - preferences use navigation only
func RenderPreferencesMenu(rows []PreferenceRow, selectedIndex int, theme Theme, contentHeight int, contentWidth int) string {
	if len(rows) == 0 {
		return ""
	}

	emojiColWidth, labelColWidth, valueColWidth := calcPrefColWidths(contentWidth)
	menuBoxWidth := emojiColWidth + labelColWidth + valueColWidth

	lines := buildPrefLines(rows, selectedIndex, theme, emojiColWidth, labelColWidth, valueColWidth)

	return assembleMenuOutput(lines, contentHeight, contentWidth, menuBoxWidth)
}

func calcPrefColWidths(contentWidth int) (emojiColWidth, labelColWidth, valueColWidth int) {
	emojiColWidth = 3
	labelColWidth = 18
	valueColWidth = 10

	if contentWidth < (emojiColWidth + labelColWidth + valueColWidth) {
		labelColWidth = contentWidth - emojiColWidth - valueColWidth
		if labelColWidth < 0 {
			labelColWidth = 0
		}
	}
	return
}

func buildPrefLines(rows []PreferenceRow, selectedIndex int, theme Theme, emojiColWidth, labelColWidth, valueColWidth int) []string {
	var lines []string
	for i, row := range rows {
		emojiCol := renderEmojiCol(row.Emoji, emojiColWidth)
		labelCol := renderLabelCol(row.Label, labelColWidth)
		valueCol := renderValueCol(row.Value, valueColWidth)

		var styledLine string
		if i == selectedIndex {
			styledLine = styleSelectedPrefRow(emojiCol, labelCol, valueCol, theme)
		} else {
			styledLine = styleNormalPrefRow(emojiCol, labelCol, valueCol, theme)
		}
		lines = append(lines, styledLine)
	}
	return lines
}

func styleSelectedPrefRow(emojiCol, labelCol, valueCol string, theme Theme) string {
	emojiStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.LabelTextColor))
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.MainBackgroundColor)).
		Background(lipgloss.Color(theme.MenuSelectionBackground)).
		Bold(true)
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.AccentTextColor)).Bold(true)
	return emojiStyle.Render(emojiCol) + labelStyle.Render(labelCol) + valueStyle.Render(valueCol)
}

func styleNormalPrefRow(emojiCol, labelCol, valueCol string, theme Theme) string {
	emojiStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.LabelTextColor))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.LabelTextColor))
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ContentTextColor))
	return emojiStyle.Render(emojiCol) + labelStyle.Render(labelCol) + valueStyle.Render(valueCol)
}

// RenderPreferencesWithBanner renders preferences (left) + banner (right)
// 50/50 split, same layout as main menu
func RenderPreferencesWithBanner(cfg *config.Config, selectedIndex int, theme Theme, sizing DynamicSizing) string {
	// 50/50 split
	leftWidth := sizing.ContentInnerWidth / 2
	rightWidth := sizing.ContentInnerWidth - leftWidth

	// Build preference rows from config
	rows := BuildPreferenceRows(cfg)

	// Left column: preferences menu
	menuContent := RenderPreferencesMenu(rows, selectedIndex, theme, sizing.ContentHeight, leftWidth)

	menuColumn := lipgloss.NewStyle().
		Width(leftWidth).
		Height(sizing.ContentHeight).
		Align(lipgloss.Left).
		AlignVertical(lipgloss.Center).
		Render(menuContent)

	// Right column: banner (same as main menu)
	banner := RenderBannerDynamic(rightWidth, sizing.ContentHeight)

	bannerColumn := lipgloss.NewStyle().
		Width(rightWidth).
		Height(sizing.ContentHeight).
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(banner)

	return lipgloss.JoinHorizontal(lipgloss.Top, menuColumn, bannerColumn)
}
