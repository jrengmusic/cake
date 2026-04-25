package ui

// MenuRow represents a single menu row
// Fixed 8 rows: [0]Project [1]Regenerate [2]OpenIDE [3]Separator [4]Configuration [5]Build [6]Clean [7]CleanAll
type MenuRow struct {
	ID            string // "project", "regenerate", "openIde", "separator", "configuration", "build", "clean", "cleanAll"
	Shortcut      string // Actual key for handler: "", "g", "o", "", "", "b", "c", "x"
	ShortcutLabel string // Display label (right-aligned): "", "g", "o", "", "", "b", "c", "x"
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
func GenerateMenuRows(projectLabel string, configuration string, canOpenIDE bool, canClean bool, hasBuild bool, hasBuildsToClean bool, isIDEGenerator bool) []MenuRow {
	regenerateLabel := "Generate"
	if hasBuild {
		regenerateLabel = "Regenerate"
	}
	regenerateHint := "Run initial CMake configuration"
	if hasBuild {
		regenerateHint = "Re-run CMake configuration"
	}

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
			Label:         regenerateLabel,
			Value:         "",
			Visible:       true,
			IsAction:      true,
			IsSelectable:  true,
			Hint:          regenerateHint,
		},
		{
			ID:            "openIde",
			Shortcut:      "o",
			ShortcutLabel: "o",
			Emoji:         "📂",
			Label:         openIdeLabel(isIDEGenerator),
			Value:         "",
			Visible:       true,
			IsAction:      true,
			IsSelectable:  canOpenIDE,
			Hint:          openIdeHint(isIDEGenerator),
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
			Shortcut:      "c",
			ShortcutLabel: "c",
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
			Shortcut:      "x",
			ShortcutLabel: "x",
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

func openIdeLabel(isIDEGenerator bool) string {
	if isIDEGenerator {
		return "Open IDE"
	}
	return "Open Editor"
}

func openIdeHint(isIDEGenerator bool) string {
	if isIDEGenerator {
		return "Open project in IDE"
	}
	return "Open build directory in editor"
}
