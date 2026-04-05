// Package ui provides user interface components for the cake application
package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ConfirmationConfig defines the configuration for a confirmation dialog
type ConfirmationConfig struct {
	Title       string
	Explanation string
	YesLabel    string
	NoLabel     string
	ActionID    string
}

// ButtonSelection represents which button is currently selected
type ButtonSelection string

const (
	ButtonYes ButtonSelection = "yes"
	ButtonNo  ButtonSelection = "no"
)

// ConfirmationDialog represents a confirmation dialog state
type ConfirmationDialog struct {
	Config         ConfirmationConfig
	Width          int
	Theme          *Theme
	Active         bool
	Context        map[string]string
	SelectedButton ButtonSelection
}

// NewConfirmationDialog creates a new confirmation dialog
func NewConfirmationDialog(config ConfirmationConfig, width int, theme *Theme) *ConfirmationDialog {
	return &ConfirmationDialog{
		Config:         config,
		Width:          width,
		Theme:          theme,
		Active:         false,
		Context:        make(map[string]string),
		SelectedButton: ButtonYes, // Default to Yes button
	}
}

// NewConfirmationDialogWithDefault creates a confirmation dialog with specified default button
func NewConfirmationDialogWithDefault(config ConfirmationConfig, width int, theme *Theme, defaultButton ButtonSelection) *ConfirmationDialog {
	return &ConfirmationDialog{
		Config:         config,
		Width:          width,
		Theme:          theme,
		Active:         false,
		Context:        make(map[string]string),
		SelectedButton: defaultButton,
	}
}

// SelectYes selects the Yes button
func (c *ConfirmationDialog) SelectYes() {
	c.SelectedButton = ButtonYes
}

// SelectNo selects the No button
func (c *ConfirmationDialog) SelectNo() {
	c.SelectedButton = ButtonNo
}

// GetSelectedButton returns the currently selected button
func (c *ConfirmationDialog) GetSelectedButton() ButtonSelection {
	return c.SelectedButton
}

// ApplyContext applies context placeholders to the config
func (c *ConfirmationDialog) ApplyContext() ConfirmationConfig {
	config := c.Config

	// Apply context to title
	if c.Context != nil {
		config.Title = applyPlaceholders(config.Title, c.Context)
		config.Explanation = applyPlaceholders(config.Explanation, c.Context)
	}

	return config
}

// applyPlaceholders replaces {placeholder} with context values
func applyPlaceholders(text string, context map[string]string) string {
	result := text
	for key, value := range context {
		placeholder := "{" + key + "}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// Render renders the confirmation dialog centered within the given height using DynamicSizing
func (c *ConfirmationDialog) Render(height int) string {
	config := c.ApplyContext()
	dialogWidth := c.Width - 10

	dialogStyle := buildDialogStyle(c.Theme, dialogWidth)
	explanationStyle := buildExplanationStyle(c.Theme, dialogWidth)
	titleStyle := buildDialogTitleStyle(c.Theme, dialogWidth)
	yesButtonStyle, noButtonStyle := buildButtonStyles(c.Theme, c.SelectedButton)

	content := buildDialogContent(config, titleStyle, explanationStyle, yesButtonStyle, noButtonStyle, lipgloss.Color(c.Theme.ConfirmationDialogBackground))
	dialog := dialogStyle.Render(content)

	return lipgloss.Place(c.Width, height, lipgloss.Center, lipgloss.Center, dialog)
}

func buildDialogStyle(theme *Theme, dialogWidth int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(dialogWidth).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.BoxBorderColor)).
		Background(lipgloss.Color(theme.ConfirmationDialogBackground)).
		Padding(1, 2).
		Align(lipgloss.Center)
}

func buildExplanationStyle(theme *Theme, dialogWidth int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.ContentTextColor)).
		Background(lipgloss.Color(theme.ConfirmationDialogBackground)).
		Width(dialogWidth - 4).
		Align(lipgloss.Left)
}

func buildDialogTitleStyle(theme *Theme, dialogWidth int) lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Width(dialogWidth - 4).
		Background(lipgloss.Color(theme.ConfirmationDialogBackground))
}

func buildButtonStyles(theme *Theme, selected ButtonSelection) (yesStyle lipgloss.Style, noStyle lipgloss.Style) {
	selectedFg := lipgloss.Color(theme.ButtonSelectedTextColor)
	selectedBg := lipgloss.Color(theme.MenuSelectionBackground)
	unselectedFg := lipgloss.Color(theme.ContentTextColor)
	unselectedBg := lipgloss.Color(theme.InlineBackgroundColor)

	if selected == ButtonYes {
		yesStyle = lipgloss.NewStyle().Foreground(selectedFg).Background(selectedBg).Bold(true).Padding(0, 2)
	} else {
		yesStyle = lipgloss.NewStyle().Foreground(unselectedFg).Background(unselectedBg).Bold(true).Padding(0, 2)
	}

	if selected == ButtonNo {
		noStyle = lipgloss.NewStyle().Foreground(selectedFg).Background(selectedBg).Bold(true).Padding(0, 2)
	} else {
		noStyle = lipgloss.NewStyle().Foreground(unselectedFg).Background(unselectedBg).Bold(true).Padding(0, 2)
	}
	return
}

func buildDialogContent(config ConfirmationConfig, titleStyle, explanationStyle, yesButtonStyle, noButtonStyle lipgloss.Style, dialogBg lipgloss.Color) string {
	var content strings.Builder

	content.WriteString(titleStyle.Render(config.Title) + "\n")
	content.WriteString("\n")
	content.WriteString(explanationStyle.Render(config.Explanation) + "\n")
	content.WriteString("\n")

	yesButton := yesButtonStyle.Render(strings.ToUpper(config.YesLabel))

	var buttonRow string
	if config.NoLabel == "" {
		buttonRow = yesButton
	} else {
		noButton := noButtonStyle.Render(strings.ToUpper(config.NoLabel))
		buttonGap := lipgloss.NewStyle().Background(dialogBg).Render("  ")
		buttonRow = lipgloss.JoinHorizontal(lipgloss.Center, yesButton, buttonGap, noButton)
	}

	buttonContainer := lipgloss.NewStyle().Align(lipgloss.Center)
	content.WriteString(buttonContainer.Render(buttonRow))

	return content.String()
}
