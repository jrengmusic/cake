package main

import (
	"fmt"
	"os"

	"cake/internal/app"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	application := app.NewApplication()
	program := tea.NewProgram(application, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
