package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// GfxTheme is the default TIT theme - all other themes derive from this reference
const GfxTheme = `name = "GFX"
description = "TIT default theme - reference for all other themes"

[palette]
mainBackgroundColor = "#090D12"
inlineBackgroundColor = "#1B2A31"
selectionBackgroundColor = "#0D141C"

contentTextColor = "#4E8C93"
labelTextColor = "#8CC9D9"
dimmedTextColor = "#33535B"
accentTextColor = "#01C2D2"
highlightTextColor = "#D1D5DA"
terminalTextColor = "#999999"

cwdTextColor = "#67DFEF"
footerTextColor = "#519299"

boxBorderColor = "#8CC9D9"
separatorColor = "#1B2A31"

menuSelectionBackground = "#7EB8C5"
buttonSelectedTextColor = "#0D1418"
confirmationDialogBackground = "#112130"  # trappedDarkness

outputStdoutColor = "#999999"
outputStderrColor = "#FC704C"
outputStatusColor = "#01C2D2"
outputWarningColor = "#F2AB53"
outputDebugColor = "#33535B"
outputInfoColor = "#01C2D2"
`

// SpringTheme is a spring-themed color palette
const SpringTheme = `name = "Spring"
description = "Fresh spring greens with vibrant energy"

[palette]
mainBackgroundColor = "#323B9E"
inlineBackgroundColor = "#0972BB"
selectionBackgroundColor = "#090D12"

contentTextColor = "#179CA8"
labelTextColor = "#90D88D"
dimmedTextColor = "#C8E189"
accentTextColor = "#FEEA85"
highlightTextColor = "#D1D5DA"
terminalTextColor = "#999999"

cwdTextColor = "#FEEA85"
footerTextColor = "#58C9BA"

boxBorderColor = "#90D88D"
separatorColor = "#0972BB"

menuSelectionBackground = "#5BCF90"
buttonSelectedTextColor = "#3F2894"
confirmationDialogBackground = "#244DA8"  # ceruleanBlue

outputStdoutColor = "#999999"
outputStderrColor = "#FD5B68"
outputStatusColor = "#4ECB71"
outputWarningColor = "#F67F78"
outputDebugColor = "#C8E189"
outputInfoColor = "#37CB9F"
`

// SummerTheme is a summer-themed color palette
const SummerTheme = `name = "Summer"
description = "Warm summer blues and bright sunshine"

[palette]
mainBackgroundColor = "#000000"
inlineBackgroundColor = "#4D88D1"
selectionBackgroundColor = "#090D12"

contentTextColor = "#3CA7E0"
labelTextColor = "#19E5FF"
dimmedTextColor = "#5E68C1"
accentTextColor = "#FFBF16"
highlightTextColor = "#D1D5DA"
terminalTextColor = "#999999"

cwdTextColor = "#FFBF16"
footerTextColor = "#8667BF"

boxBorderColor = "#19E5FF"
separatorColor = "#4D88D1"

menuSelectionBackground = "#FE62B9"
buttonSelectedTextColor = "#8667BF"
confirmationDialogBackground = "#2BC6F0"  # pictonBlue

outputStdoutColor = "#999999"
outputStderrColor = "#FF3469"
outputStatusColor = "#00FFFF"
outputWarningColor = "#FF9700"
outputDebugColor = "#5E68C1"
outputInfoColor = "#2BC6F0"
`

// AutumnTheme is an autumn-themed color palette
const AutumnTheme = `name = "Autumn"
description = "Rich autumn oranges and warm earth tones"

[palette]
mainBackgroundColor = "#3E0338"
inlineBackgroundColor = "#5E063E"
selectionBackgroundColor = "#090D12"

contentTextColor = "#E78C79"
labelTextColor = "#F9C94D"
dimmedTextColor = "#F09D06"
accentTextColor = "#F5BB09"
highlightTextColor = "#D1D5DA"
terminalTextColor = "#999999"

cwdTextColor = "#F5BB09"
footerTextColor = "#CD5861"

boxBorderColor = "#F9C94D"
separatorColor = "#5E063E"

menuSelectionBackground = "#F1AE37"
buttonSelectedTextColor = "#3E0338"
confirmationDialogBackground = "#7D0E36"  # roseBudCherry

outputStdoutColor = "#999999"
outputStderrColor = "#DC3003"
outputStatusColor = "#F5BB09"
outputWarningColor = "#E85C03"
outputDebugColor = "#F09D06"
outputInfoColor = "#F48C06"
`

// WinterTheme is a winter-themed color palette
const WinterTheme = `name = "Winter"
description = "Cool winter purples with subtle elegance"

[palette]
mainBackgroundColor = "#233253"
inlineBackgroundColor = "#334676"
selectionBackgroundColor = "#090D12"

contentTextColor = "#CAD0E6"
labelTextColor = "#7F95D6"
dimmedTextColor = "#9BA9D0"
accentTextColor = "#F6F5FA"
highlightTextColor = "#D1D5DA"
terminalTextColor = "#999999"

cwdTextColor = "#F6F5FA"
footerTextColor = "#9BA9D0"

boxBorderColor = "#7F95D6"
separatorColor = "#334676"

menuSelectionBackground = "#7F95D6"
buttonSelectedTextColor = "#F6F5FA"
confirmationDialogBackground = "#233253"  # cloudBurst

outputStdoutColor = "#999999"
outputStderrColor = "#E0BACF"
outputStatusColor = "#435A98"
outputWarningColor = "#CEBAC5"
outputDebugColor = "#9BA9D0"
outputInfoColor = "#435A98"
`

type ThemeDefinition struct {
	Name        string `toml:"name"`
	Description string `toml:"description"`
	Palette     struct {
		MainBackgroundColor          string `toml:"mainBackgroundColor"`
		InlineBackgroundColor        string `toml:"inlineBackgroundColor"`
		SelectionBackgroundColor     string `toml:"selectionBackgroundColor"`
		ContentTextColor             string `toml:"contentTextColor"`
		LabelTextColor               string `toml:"labelTextColor"`
		DimmedTextColor              string `toml:"dimmedTextColor"`
		AccentTextColor              string `toml:"accentTextColor"`
		HighlightTextColor           string `toml:"highlightTextColor"`
		TerminalTextColor            string `toml:"terminalTextColor"`
		CwdTextColor                 string `toml:"cwdTextColor"`
		FooterTextColor              string `toml:"footerTextColor"`
		BoxBorderColor               string `toml:"boxBorderColor"`
		SeparatorColor               string `toml:"separatorColor"`
		ConflictPaneUnfocusedBorder  string `toml:"conflictPaneUnfocusedBorder"`
		ConflictPaneFocusedBorder    string `toml:"conflictPaneFocusedBorder"`
		ConflictSelectionForeground  string `toml:"conflictSelectionForeground"`
		ConflictSelectionBackground  string `toml:"conflictSelectionBackground"`
		ConflictPaneTitleColor       string `toml:"conflictPaneTitleColor"`
		StatusClean                  string `toml:"statusClean"`
		StatusDirty                  string `toml:"statusDirty"`
		TimelineSynchronized         string `toml:"timelineSynchronized"`
		TimelineLocalAhead           string `toml:"timelineLocalAhead"`
		TimelineLocalBehind          string `toml:"timelineLocalBehind"`
		OperationReady               string `toml:"operationReady"`
		OperationNotRepo             string `toml:"operationNotRepo"`
		OperationTimeTravel          string `toml:"operationTimeTravel"`
		OperationConflicted          string `toml:"operationConflicted"`
		OperationMerging             string `toml:"operationMerging"`
		OperationRebasing            string `toml:"operationRebasing"`
		OperationDirtyOp             string `toml:"operationDirtyOp"`
		MenuSelectionBackground      string `toml:"menuSelectionBackground"`
		ButtonSelectedTextColor      string `toml:"buttonSelectedTextColor"`
		ConfirmationDialogBackground string `toml:"confirmationDialogBackground"`
		SpinnerColor                 string `toml:"spinnerColor"`
		DiffAddedLineColor           string `toml:"diffAddedLineColor"`
		DiffRemovedLineColor         string `toml:"diffRemovedLineColor"`
		OutputStdoutColor            string `toml:"outputStdoutColor"`
		OutputStderrColor            string `toml:"outputStderrColor"`
		OutputStatusColor            string `toml:"outputStatusColor"`
		OutputWarningColor           string `toml:"outputWarningColor"`
		OutputDebugColor             string `toml:"outputDebugColor"`
		OutputInfoColor              string `toml:"outputInfoColor"`
	} `toml:"palette"`
}

type Theme struct {
	MainBackgroundColor          string
	InlineBackgroundColor        string
	SelectionBackgroundColor     string
	ContentTextColor             string
	LabelTextColor               string
	DimmedTextColor              string
	AccentTextColor              string
	HighlightTextColor           string
	TerminalTextColor            string
	CwdTextColor                 string
	FooterTextColor              string
	BoxBorderColor               string
	SeparatorColor               string
	ConflictPaneUnfocusedBorder  string
	ConflictPaneFocusedBorder    string
	ConflictSelectionForeground  string
	ConflictSelectionBackground  string
	ConflictPaneTitleColor       string
	StatusClean                  string
	StatusDirty                  string
	TimelineSynchronized         string
	TimelineLocalAhead           string
	TimelineLocalBehind          string
	OperationReady               string
	OperationNotRepo             string
	OperationTimeTravel          string
	OperationConflicted          string
	OperationMerging             string
	OperationRebasing            string
	OperationDirtyOp             string
	MenuSelectionBackground      string
	ButtonSelectedTextColor      string
	ConfirmationDialogBackground string
	SpinnerColor                 string
	DiffAddedLineColor           string
	DiffRemovedLineColor         string
	OutputStdoutColor            string
	OutputStderrColor            string
	OutputStatusColor            string
	OutputWarningColor           string
	OutputDebugColor             string
	OutputInfoColor              string
}

func LoadTheme(themeFilePath string) (Theme, error) {
	fileData, err := os.ReadFile(themeFilePath)
	if err != nil {
		return Theme{}, fmt.Errorf("failed to read theme file: %w", err)
	}

	var themeDef ThemeDefinition
	if err := toml.Unmarshal(fileData, &themeDef); err != nil {
		return Theme{}, fmt.Errorf("failed to parse theme file: %w", err)
	}

	theme := Theme{
		MainBackgroundColor:          themeDef.Palette.MainBackgroundColor,
		InlineBackgroundColor:        themeDef.Palette.InlineBackgroundColor,
		SelectionBackgroundColor:     themeDef.Palette.SelectionBackgroundColor,
		ContentTextColor:             themeDef.Palette.ContentTextColor,
		LabelTextColor:               themeDef.Palette.LabelTextColor,
		DimmedTextColor:              themeDef.Palette.DimmedTextColor,
		AccentTextColor:              themeDef.Palette.AccentTextColor,
		HighlightTextColor:           themeDef.Palette.HighlightTextColor,
		TerminalTextColor:            themeDef.Palette.TerminalTextColor,
		CwdTextColor:                 themeDef.Palette.CwdTextColor,
		FooterTextColor:              themeDef.Palette.FooterTextColor,
		BoxBorderColor:               themeDef.Palette.BoxBorderColor,
		SeparatorColor:               themeDef.Palette.SeparatorColor,
		ConflictPaneUnfocusedBorder:  themeDef.Palette.ConflictPaneUnfocusedBorder,
		ConflictPaneFocusedBorder:    themeDef.Palette.ConflictPaneFocusedBorder,
		ConflictSelectionForeground:  themeDef.Palette.ConflictSelectionForeground,
		ConflictSelectionBackground:  themeDef.Palette.ConflictSelectionBackground,
		ConflictPaneTitleColor:       themeDef.Palette.ConflictPaneTitleColor,
		StatusClean:                  themeDef.Palette.StatusClean,
		StatusDirty:                  themeDef.Palette.StatusDirty,
		TimelineSynchronized:         themeDef.Palette.TimelineSynchronized,
		TimelineLocalAhead:           themeDef.Palette.TimelineLocalAhead,
		TimelineLocalBehind:          themeDef.Palette.TimelineLocalBehind,
		OperationReady:               themeDef.Palette.OperationReady,
		OperationNotRepo:             themeDef.Palette.OperationNotRepo,
		OperationTimeTravel:          themeDef.Palette.OperationTimeTravel,
		OperationConflicted:          themeDef.Palette.OperationConflicted,
		OperationMerging:             themeDef.Palette.OperationMerging,
		OperationRebasing:            themeDef.Palette.OperationRebasing,
		OperationDirtyOp:             themeDef.Palette.OperationDirtyOp,
		MenuSelectionBackground:      themeDef.Palette.MenuSelectionBackground,
		ButtonSelectedTextColor:      themeDef.Palette.ButtonSelectedTextColor,
		ConfirmationDialogBackground: themeDef.Palette.ConfirmationDialogBackground,
		SpinnerColor:                 themeDef.Palette.SpinnerColor,
		DiffAddedLineColor:           themeDef.Palette.DiffAddedLineColor,
		DiffRemovedLineColor:         themeDef.Palette.DiffRemovedLineColor,
		OutputStdoutColor:            themeDef.Palette.OutputStdoutColor,
		OutputStderrColor:            themeDef.Palette.OutputStderrColor,
		OutputStatusColor:            themeDef.Palette.OutputStatusColor,
		OutputWarningColor:           themeDef.Palette.OutputWarningColor,
		OutputDebugColor:             themeDef.Palette.OutputDebugColor,
		OutputInfoColor:              themeDef.Palette.OutputInfoColor,
	}

	return theme, nil
}

// EnsureFiveThemesExist creates/regenerates all 5 themes at startup
func EnsureFiveThemesExist() error {
	configThemeDir := filepath.Join(getConfigDirectory(), "themes")
	if err := os.MkdirAll(configThemeDir, 0755); err != nil {
		return fmt.Errorf("failed to create themes directory: %w", err)
	}

	themes := map[string]string{
		"gfx":    GfxTheme,
		"spring": SpringTheme,
		"summer": SummerTheme,
		"autumn": AutumnTheme,
		"winter": WinterTheme,
	}

	for name, content := range themes {
		themePath := filepath.Join(configThemeDir, name+".toml")
		if err := os.WriteFile(themePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s theme: %w", name, err)
		}
	}

	return nil
}

func CreateDefaultThemeIfMissing() (string, error) {
	return "", EnsureFiveThemesExist()
}

func LoadDefaultTheme() (Theme, error) {
	themeFile := filepath.Join(getConfigDirectory(), "themes", "gfx.toml")
	return LoadTheme(themeFile)
}

func getConfigDirectory() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".cake"
	}
	return filepath.Join(home, ".config", "cake")
}

// DiscoverAvailableThemes returns a list of available theme names
func DiscoverAvailableThemes() ([]string, error) {
	themesDir := filepath.Join(getConfigDirectory(), "themes")

	files, err := os.ReadDir(themesDir)
	if err != nil {
		return nil, err
	}

	var themes []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".toml" {
			themeName := file.Name()[:len(file.Name())-5] // Remove .toml extension
			themes = append(themes, themeName)
		}
	}

	return themes, nil
}

// LoadThemeByName loads a theme by name from the themes directory
func LoadThemeByName(themeName string) (Theme, error) {
	themeFile := filepath.Join(getConfigDirectory(), "themes", themeName+".toml")
	return LoadTheme(themeFile)
}
