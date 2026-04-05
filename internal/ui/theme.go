package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

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
		if _, err := os.Stat(themePath); os.IsNotExist(err) {
			if err := os.WriteFile(themePath, []byte(content), 0644); err != nil {
				return fmt.Errorf("failed to write %s theme: %w", name, err)
			}
		}
	}

	return nil
}

// ThemeDefinition represents a theme file structure
type ThemeDefinition struct {
	Name        string `toml:"name"`
	Description string `toml:"description"`
	Palette     struct {
		// Backgrounds
		MainBackgroundColor      string `toml:"mainBackgroundColor"`
		InlineBackgroundColor    string `toml:"inlineBackgroundColor"`
		SelectionBackgroundColor string `toml:"selectionBackgroundColor"`

		// Text - Content & Body
		ContentTextColor   string `toml:"contentTextColor"`
		LabelTextColor     string `toml:"labelTextColor"`
		DimmedTextColor    string `toml:"dimmedTextColor"`
		AccentTextColor    string `toml:"accentTextColor"`
		HighlightTextColor string `toml:"highlightTextColor"`

		// Special Text
		CwdTextColor    string `toml:"cwdTextColor"`
		FooterTextColor string `toml:"footerTextColor"`

		// Borders
		BoxBorderColor string `toml:"boxBorderColor"`
		SeparatorColor string `toml:"separatorColor"`

		// Confirmation Dialog
		ConfirmationDialogBackground string `toml:"confirmationDialogBackground"`

		// UI Elements / Buttons
		MenuSelectionBackground string `toml:"menuSelectionBackground"`
		ButtonSelectedTextColor string `toml:"buttonSelectedTextColor"`

		// Console Output Colors
		OutputStdoutColor  string `toml:"outputStdoutColor"`
		OutputStderrColor  string `toml:"outputStderrColor"`
		OutputStatusColor  string `toml:"outputStatusColor"`
		OutputWarningColor string `toml:"outputWarningColor"`
		OutputDebugColor   string `toml:"outputDebugColor"`
		OutputInfoColor    string `toml:"outputInfoColor"`
	} `toml:"palette"`
}

// Theme defines all semantic colors from the active theme
type Theme struct {
	// Backgrounds
	MainBackgroundColor      string
	InlineBackgroundColor    string
	SelectionBackgroundColor string

	// Text - Content & Body
	ContentTextColor   string
	LabelTextColor     string
	DimmedTextColor    string
	AccentTextColor    string
	HighlightTextColor string

	// Special Text
	CwdTextColor    string
	FooterTextColor string

	// Borders
	BoxBorderColor string
	SeparatorColor string

	// Confirmation Dialog
	ConfirmationDialogBackground string

	// UI Elements / Buttons
	MenuSelectionBackground string
	ButtonSelectedTextColor string

	// Console Output Colors
	OutputStdoutColor  string
	OutputStderrColor  string
	OutputStatusColor  string
	OutputWarningColor string
	OutputDebugColor   string
	OutputInfoColor    string
}

// LoadTheme loads a theme from a TOML file
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
		// Backgrounds
		MainBackgroundColor:      themeDef.Palette.MainBackgroundColor,
		InlineBackgroundColor:    themeDef.Palette.InlineBackgroundColor,
		SelectionBackgroundColor: themeDef.Palette.SelectionBackgroundColor,

		// Text - Content & Body
		ContentTextColor:   themeDef.Palette.ContentTextColor,
		LabelTextColor:     themeDef.Palette.LabelTextColor,
		DimmedTextColor:    themeDef.Palette.DimmedTextColor,
		AccentTextColor:    themeDef.Palette.AccentTextColor,
		HighlightTextColor: themeDef.Palette.HighlightTextColor,

		// Special Text
		CwdTextColor:    themeDef.Palette.CwdTextColor,
		FooterTextColor: themeDef.Palette.FooterTextColor,

		// Borders
		BoxBorderColor: themeDef.Palette.BoxBorderColor,
		SeparatorColor: themeDef.Palette.SeparatorColor,

		// Confirmation Dialog
		ConfirmationDialogBackground: themeDef.Palette.ConfirmationDialogBackground,

		// UI Elements / Buttons
		MenuSelectionBackground: themeDef.Palette.MenuSelectionBackground,
		ButtonSelectedTextColor: themeDef.Palette.ButtonSelectedTextColor,

		// Console Output Colors
		OutputStdoutColor:  themeDef.Palette.OutputStdoutColor,
		OutputStderrColor:  themeDef.Palette.OutputStderrColor,
		OutputStatusColor:  themeDef.Palette.OutputStatusColor,
		OutputWarningColor: themeDef.Palette.OutputWarningColor,
		OutputDebugColor:   themeDef.Palette.OutputDebugColor,
		OutputInfoColor:    themeDef.Palette.OutputInfoColor,
	}

	return theme, nil
}

// CreateDefaultThemeIfMissing creates or regenerates all 5 themes (gfx + 4 seasons)
// SSOT: Always regenerates from GfxTheme to ensure all colors are current
func CreateDefaultThemeIfMissing() (string, error) {
	return "", EnsureFiveThemesExist()
}

// LoadDefaultTheme loads the default gfx theme
func LoadDefaultTheme() (Theme, error) {
	themeFile := filepath.Join(getConfigDirectory(), "themes", "gfx.toml")
	return LoadTheme(themeFile)
}

// getConfigDirectory returns the cake config directory
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
		return nil, fmt.Errorf("failed to read themes directory: %w", err)
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

// GetNextTheme cycles to the next available theme
func GetNextTheme(currentTheme string) (string, error) {
	themes, err := DiscoverAvailableThemes()
	if err != nil {
		return "", err
	}

	if len(themes) == 0 {
		return "", fmt.Errorf("no themes found")
	}

	// Find current theme index
	currentIndex := -1
	for i, theme := range themes {
		if theme == currentTheme {
			currentIndex = i
			break
		}
	}

	// Cycle to next (or first if current not found)
	nextIndex := (currentIndex + 1) % len(themes)
	return themes[nextIndex], nil
}
