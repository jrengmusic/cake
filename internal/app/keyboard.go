package app

import tea "github.com/charmbracelet/bubbletea"

type KeyHandler func(*Application) (tea.Model, tea.Cmd)
