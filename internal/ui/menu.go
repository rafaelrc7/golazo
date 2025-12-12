// Package ui provides rendering functions for the terminal user interface.
package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Elegant color palette
	textColor      = lipgloss.Color("15")  // White
	accentColor    = lipgloss.Color("11")  // Bright yellow
	selectedColor  = lipgloss.Color("11")  // Bright yellow for selection
	borderColor    = lipgloss.Color("240") // Subtle gray
	dimColor       = lipgloss.Color("244") // Muted gray
	highlightColor = lipgloss.Color("220") // Warm yellow
	liveColor      = lipgloss.Color("196") // Red for live
	goalColor      = lipgloss.Color("46")  // Green for goals
	cardColor      = lipgloss.Color("226") // Yellow for cards

	// Menu styles
	menuItemStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Padding(0, 1)

	menuItemSelectedStyle = lipgloss.NewStyle().
				Foreground(selectedColor).
				Bold(true).
				Padding(0, 1)

	menuTitleStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Align(lipgloss.Center).
			Padding(1, 0)

	menuHelpStyle = lipgloss.NewStyle().
			Foreground(dimColor).
			Align(lipgloss.Center).
			Padding(1, 0)
)

// RenderMainMenu renders the main menu view with navigation options.
// width and height specify the terminal dimensions.
// selected indicates which menu item is currently selected (0-indexed).
func RenderMainMenu(width, height, selected int) string {
	menuItems := []string{
		"Live Matches",
		"Favourites",
	}

	items := make([]string, 0, len(menuItems))
	for i, item := range menuItems {
		if i == selected {
			items = append(items, menuItemSelectedStyle.Render("→ "+item))
		} else {
			items = append(items, menuItemStyle.Render("  "+item))
		}
	}

	menuContent := strings.Join(items, "\n")

	title := menuTitleStyle.Render("⚽ Golazo")
	help := menuHelpStyle.Render("↑/↓: navigate  Enter: select  q: quit")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		strings.Repeat("\n", 2),
		menuContent,
		strings.Repeat("\n", 2),
		help,
	)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
