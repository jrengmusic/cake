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

	return &Application{
		width:           80,
		height:          24,
		theme:           theme,
		sizing:          ui.NewDynamicSizing(),
		mode:            ModeMenu,
		selectedIndex:   0,
		menuItems:       []PreferenceRow{},
		projectState:    state.NewProjectState(),
		config:          cfg,
		outputBuffer:    ui.GetBuffer(),
		footerHint:      FooterHints["menu_navigate"],
		quitConfirmTime: time.Now(),
	}
}
