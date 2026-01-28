package app

type AppMode int

const (
	ModeMenu AppMode = iota
	ModePreferences
	ModeConsole
)

type ModeMetadata struct {
	Name         string
	Description  string
	AcceptsInput bool
	IsAsync      bool
}

var modeDescriptions = map[AppMode]ModeMetadata{
	ModeMenu: {
		Name:         "menu",
		Description:  "Main preference menu",
		AcceptsInput: true,
		IsAsync:      false,
	},
	ModePreferences: {
		Name:         "preferences",
		Description:  "Application preferences",
		AcceptsInput: true,
		IsAsync:      false,
	},
	ModeConsole: {
		Name:         "console",
		Description:  "Async operation console display",
		AcceptsInput: true,
		IsAsync:      true,
	},
}

func GetModeMetadata(m AppMode) ModeMetadata {
	if meta, exists := modeDescriptions[m]; exists {
		return meta
	}
	return ModeMetadata{Name: "unknown", Description: "Unknown mode"}
}

func (m AppMode) String() string {
	return GetModeMetadata(m).Name
}
