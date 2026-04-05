package app

type AppMode int

const (
	ModeInvalidProject AppMode = iota // Not a CMake project - show "cake is a lie"
	ModeMenu
	ModePreferences
	ModeConsole
)

var modeNames = map[AppMode]string{
	ModeInvalidProject: "invalid",
	ModeMenu:           "menu",
	ModePreferences:    "preferences",
	ModeConsole:        "console",
}

func (m AppMode) String() string {
	if name, exists := modeNames[m]; exists {
		return name
	}
	return "unknown"
}
