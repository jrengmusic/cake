package app

import (
	"cake/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

// MessageHandler handles specific message types
type MessageHandler interface {
	CanHandle(msg tea.Msg) bool
	Handle(a *Application, msg tea.Msg) (tea.Model, tea.Cmd)
}

// WindowSizeHandler handles terminal resize
type WindowSizeHandler struct{}

func (h WindowSizeHandler) CanHandle(msg tea.Msg) bool {
	_, ok := msg.(tea.WindowSizeMsg)
	return ok
}

func (h WindowSizeHandler) Handle(a *Application, msg tea.Msg) (tea.Model, tea.Cmd) {
	m := msg.(tea.WindowSizeMsg)
	a.width = m.Width
	a.height = m.Height
	a.sizing = ui.CalculateDynamicSizing(m.Width, m.Height)
	return a, nil
}

// KeyDispatcher routes keyboard input
type KeyDispatcher struct {
	handlers map[AppMode]func(*Application, tea.KeyMsg) (tea.Model, tea.Cmd)
}

func NewKeyDispatcher() *KeyDispatcher {
	return &KeyDispatcher{
		handlers: make(map[AppMode]func(*Application, tea.KeyMsg) (tea.Model, tea.Cmd)),
	}
}

func (h *KeyDispatcher) Register(mode AppMode, handler func(*Application, tea.KeyMsg) (tea.Model, tea.Cmd)) {
	h.handlers[mode] = handler
}

func (h *KeyDispatcher) CanHandle(msg tea.Msg) bool {
	_, ok := msg.(tea.KeyMsg)
	return ok
}

func (h *KeyDispatcher) Handle(a *Application, msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg := msg.(tea.KeyMsg)
	if handler, ok := h.handlers[a.mode]; ok {
		return handler(a, keyMsg)
	}
	return a, nil
}
