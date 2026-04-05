package app

import (
	"github.com/jrengmusic/cake/internal/config"
	"github.com/jrengmusic/cake/internal/state"
	"github.com/jrengmusic/cake/internal/ui"
	"github.com/jrengmusic/cake/internal/utils"
	"time"
)

func loadTheme(cfg *config.Config) ui.Theme {
	if cfg != nil {
		// theme load failure is non-fatal: fall through to default
		theme, _ := ui.LoadThemeByName(cfg.Theme())
		if theme.MainBackgroundColor != "" {
			return theme
		}
	}
	// theme load failure is non-fatal: fall through to zero value
	theme, _ := ui.LoadDefaultTheme()
	return theme
}

func captureVSEnvironment() []string {
	vcVarsAllPath, vsErr := utils.FindVCVarsAll()
	if vsErr != nil {
		return nil
	}
	// VS env capture failure is non-fatal: build operations fall back to system PATH
	capturedVSEnv, _ := utils.CaptureVSEnv(vcVarsAllPath)
	return capturedVSEnv
}

func initialModeAndHint(projectState *state.ProjectState, cfg *config.Config) (AppMode, string) {
	if !projectState.HasCMakeLists {
		return ModeInvalidProject, "The cake is a lie"
	}

	if cfg != nil {
		if lastProject := cfg.LastProject(); lastProject != "" {
			projectState.SetSelectedProject(lastProject)
		}
		if lastConfig := cfg.LastConfiguration(); lastConfig != "" {
			projectState.SetConfiguration(lastConfig)
		}
	}
	return ModeMenu, FooterHints["menu_navigate"]
}

func NewApplication() *Application {
	// Create theme files if missing
	ui.CreateDefaultThemeIfMissing()

	// config load failure is non-fatal: defaults used
	cfg, _ := config.Load()
	theme := loadTheme(cfg)

	projectState := state.NewProjectState()
	projectState.ForceRefresh()

	initialMode, footerHint := initialModeAndHint(projectState, cfg)

	return &Application{
		width:           DefaultTerminalWidth,
		height:          DefaultTerminalHeight,
		theme:           theme,
		sizing:          ui.NewDynamicSizing(),
		mode:            initialMode,
		selectedIndex:   0,
		menuItems:       []ui.MenuRow{},
		projectState:    projectState,
		config:          cfg,
		outputBuffer:    ui.GetBuffer(),
		footerHint:      footerHint,
		quitConfirmTime: time.Now(),
		vsEnv:           captureVSEnvironment(),
	}
}
