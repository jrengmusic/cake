package app

import (
	"cake/internal/config"
	"cake/internal/state"
	"cake/internal/ui"
	"cake/internal/utils"
	"fmt"
	"os"
	"strings"
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

	// Create project state and refresh to detect CMakeLists.txt
	projectState := state.NewProjectState()
	projectState.ForceRefresh() // This populates HasCMakeLists

	// Capture Visual Studio environment on Windows (enables cmake invocation)
	vcVarsAllPath, vsErr := utils.FindVCVarsAll()
	var capturedVSEnv []string
	debugLog, _ := os.Create(`C:\Users\jreng\cake-debug.txt`)
	if vsErr == nil {
		var captureErr error
		capturedVSEnv, captureErr = utils.CaptureVSEnv(vcVarsAllPath)
		if captureErr != nil {
			fmt.Fprintf(debugLog, "CaptureVSEnv failed: %v\n", captureErr)
		} else {
			fmt.Fprintf(debugLog, "CaptureVSEnv captured %d env vars\n", len(capturedVSEnv))
			fmt.Fprintf(debugLog, "vcVarsAllPath: %s\n", vcVarsAllPath)
			for _, e := range capturedVSEnv {
				if strings.HasPrefix(strings.ToUpper(e), "PATH=") {
					fmt.Fprintf(debugLog, "PATH: %s\n", e[:200])
					break
				}
			}
		}
	} else {
		fmt.Fprintf(debugLog, "FindVCVarsAll failed: %v\n", vsErr)
	}
	debugLog.Close()

	// Check if current directory is a valid CMake project
	isValidProject := projectState.HasCMakeLists

	var initialMode AppMode
	var footerHint string
	if isValidProject {
		// Valid CMake project - start in menu mode
		initialMode = ModeMenu
		footerHint = FooterHints["menu_navigate"]

		// Restore last chosen options from config
		if cfg != nil {
			if lastProject := cfg.LastProject(); lastProject != "" {
				projectState.SetSelectedProject(lastProject)
			}
			if lastConfig := cfg.LastConfiguration(); lastConfig != "" {
				projectState.SetConfiguration(lastConfig)
			}
		}
	} else {
		// Not a CMake project - show "cake is a lie" mode
		initialMode = ModeInvalidProject
		footerHint = "The cake is a lie"
	}

	return &Application{
		width:           80,
		height:          24,
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
		vsEnv:           capturedVSEnv,
	}
}
