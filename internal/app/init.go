package app

import (
	"cake/internal/config"
	"cake/internal/state"
	"cake/internal/ui"
	"time"
)

func NewApplication() *Application {
	// Create theme files if missing
	ui.CreateDefaultThemeIfMissing()

	// Load configuration
	cfg, _ := config.Load()

	// Load theme from config (default: gfx)
	var theme ui.Theme
	if cfg != nil {
		theme, _ = ui.LoadThemeByName(cfg.Theme())
		if theme.MainBackgroundColor == "" {
			theme, _ = ui.LoadDefaultTheme()
		}
	} else {
		theme, _ = ui.LoadDefaultTheme()
	}

	// Create project state and restore last chosen options from config
	projectState := state.NewProjectState()
	if cfg != nil {
		// Restore last chosen project if available
		if lastProject := cfg.LastProject(); lastProject != "" {
			projectState.SetSelectedProject(lastProject)
		}
		// Restore last chosen configuration
		if lastConfig := cfg.LastConfiguration(); lastConfig != "" {
			projectState.SetConfiguration(lastConfig)
		}
	}

	return &Application{
		width:           80,
		height:          24,
		theme:           theme,
		sizing:          ui.NewDynamicSizing(),
		mode:            ModeMenu,
		selectedIndex:   0,
		menuItems:       []ui.MenuRow{},
		projectState:    projectState,
		config:          cfg,
		outputBuffer:    ui.GetBuffer(),
		footerHint:      FooterHints["menu_navigate"],
		quitConfirmTime: time.Now(),
	}
}
