package app

import (
	"time"

	"github.com/jrengmusic/cake/internal/ui"
)

type TickMsg time.Time

type BuildCompleteMsg struct {
	Success  bool
	ExitCode int
	Error    string
}

type CleanCompleteMsg struct {
	Success bool
	Error   string
}

type CleanAllCompleteMsg struct {
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

type RegenerateCompleteMsg struct {
	Success bool
	Error   string
}

// OutputRefreshMsg triggers UI re-render to show updated console output
// Sent periodically during long-running operations to display streaming output
type OutputRefreshMsg struct{}

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

var footerMessageText = map[FooterMessageType]string{
	MessageNone:              "",
	MessageCtrlCConfirm:      "Press Ctrl+C again to quit (3s timeout)",
	MessageSetupInProgress:   "Setting up CMake... (ESC to abort)",
	MessageBuildInProgress:   "Building project... (ESC to abort)",
	MessageCleanInProgress:   "Cleaning project... (ESC to abort)",
	MessageOperationComplete: "Operation completed. Press ESC to return.",
	MessageOperationFailed:   "Operation failed. Press ESC to return.",
	MessageExitBlocked:       "Operation in progress. Cannot quit.",
}

func GetFooterMessageText(msgType FooterMessageType) string {
	if msg, exists := footerMessageText[msgType]; exists {
		return msg
	}
	return ""
}

var FooterHints = map[string]string{
	"menu_navigate":    "[g] Generate [b] Build [c] Clean [o] Open [/] Config ↑↓ select",
	"setup_gen_choose": "↑↓ choose project │ Enter select │ ESC back",
	"ide_choose":       "↑↓ choose IDE project │ Enter select │ ESC back",
	"editor_choose":    "↑↓ choose build dir │ Enter select │ ESC back",
	"build_choose":     "↑↓ choose build dir │ Enter select │ ESC back",
	"clean_choose":     "↑↓ choose build dir │ Enter select │ ESC back",
	"no_build_dir":     "Build directory not found. Run Setup first.",
	"operation_wait":   "Operation in progress. Please wait...",
	"scanning":         "[Scanning...]",
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
