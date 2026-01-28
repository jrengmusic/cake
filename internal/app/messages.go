package app

import (
	"time"

	"cake/internal/ui"
)

type TickMsg time.Time

type SetupCompleteMsg struct {
	Success bool
	Error   string
}

type BuildCompleteMsg struct {
	Success  bool
	ExitCode int
	Error    string
}

type CleanCompleteMsg struct {
	Success bool
	Error   string
}

type OpenIDECompleteMsg struct {
	Success bool
	Error   string
}

type OpenEditorCompleteMsg struct {
	Success bool
	Error   string
}

type GenerateCompleteMsg struct {
	Success bool
	Error   string
}

// AutoScanTickMsg is sent periodically to trigger auto-scan
type AutoScanTickMsg struct{}

type FooterMessageType int

const (
	MessageNone FooterMessageType = iota
	MessageCtrlCConfirm
	MessageSetupInProgress
	MessageBuildInProgress
	MessageCleanInProgress
	MessageOperationComplete
	MessageOperationFailed
	MessageExitBlocked
)

func GetFooterMessageText(msgType FooterMessageType) string {
	messages := map[FooterMessageType]string{
		MessageNone:              "",
		MessageCtrlCConfirm:      "Press Ctrl+C again to quit (3s timeout)",
		MessageSetupInProgress:   "Setting up CMake... (ESC to abort)",
		MessageBuildInProgress:   "Building project... (ESC to abort)",
		MessageCleanInProgress:   "Cleaning project... (ESC to abort)",
		MessageOperationComplete: "Operation completed. Press ESC to return.",
		MessageOperationFailed:   "Operation failed. Press ESC to return.",
		MessageExitBlocked:       "Operation in progress. Cannot quit.",
	}
	if msg, exists := messages[msgType]; exists {
		return msg
	}
	return ""
}

var FooterHints = map[string]string{
	"menu_navigate":    "[g] Generate [b] Build [c] Clean [/] Config ↑↓ select",
	"setup_gen_choose": "↑↓ choose generator │ Enter select │ ESC back",
	"ide_choose":       "↑↓ choose IDE project │ Enter select │ ESC back",
	"editor_choose":    "↑↓ choose build dir │ Enter select │ ESC back",
	"build_choose":     "↑↓ choose build dir │ Enter select │ ESC back",
	"clean_choose":     "↑↓ choose build dir │ Enter select │ ESC back",
	"no_build_dir":     "Build directory not found. Run Setup first.",
	"operation_wait":   "Operation in progress. Please wait...",
	"scanning":         "[Scanning...]",
}

var ErrorMessages = map[string]string{
	"cmake_not_found":     "CMake not found in PATH",
	"invalid_generator":   "Invalid generator: %s",
	"setup_failed":        "Setup failed: %s",
	"build_failed":        "Build failed: %s",
	"clean_failed":        "Clean failed: %s",
	"invalid_build_dir":   "Build directory does not exist: %s",
	"project_detect_fail": "Failed to detect project: %s",
	"cwd_read_failed":     "Failed to get current directory",
}

var OutputMessages = map[string]string{
	"detecting_project":   "Detecting project...",
	"scanning_generators": "Scanning available generators...",
	"setup_starting":      "Configuring CMake...",
	"build_starting":      "Compiling project...",
	"clean_starting":      "Removing build artifacts...",
}

// FooterHintShortcuts defines all mode-specific footer shortcuts (SSOT)
// Key = mode identifier, Value = list of shortcuts
var FooterHintShortcuts = map[string][]ui.FooterShortcut{
	// Console mode - shows scroll controls
	"console_running": {
		{Key: "↑↓", Desc: "scroll"},
		{Key: "Esc", Desc: "abort"},
	},
	"console_complete": {
		{Key: "↑↓", Desc: "scroll"},
		{Key: "Esc", Desc: "back"},
	},

	// Preferences mode
	"preferences": {
		{Key: "↑↓", Desc: "navigate"},
		{Key: "Enter", Desc: "change"},
		{Key: "/", Desc: "back"},
	},

	// Confirmation dialog
	"confirmation": {
		{Key: "←→", Desc: "select"},
		{Key: "Enter", Desc: "confirm"},
		{Key: "Esc", Desc: "cancel"},
	},
}
