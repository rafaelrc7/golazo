package ui

import (
	"strings"

	"github.com/0xjuanma/golazo/internal/api"
	"github.com/0xjuanma/golazo/internal/constants"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

// RenderLiveMatchesListPanel renders the left panel using bubbletea list component.
// Note: listModel is passed by value, so SetSize must be called before this function.
func RenderLiveMatchesListPanel(width, height int, listModel list.Model) string {
	// Wrap list in panel
	title := panelTitleStyle.Width(width - 6).Render(constants.PanelLiveMatches)
	listView := listModel.View()

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		listView,
	)

	panel := panelStyle.
		Width(width).
		Height(height).
		Render(content)

	return panel
}

// RenderStatsListPanel renders the left panel for stats view using bubbletea list component.
// Note: listModel is passed by value, so SetSize must be called before this function.
func RenderStatsListPanel(width, height int, listModel list.Model) string {
	// Wrap list in panel
	title := panelTitleStyle.Width(width - 6).Render(constants.PanelFinishedMatches)
	listView := listModel.View()

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		listView,
	)

	panel := panelStyle.
		Width(width).
		Height(height).
		Render(content)

	return panel
}

// RenderMultiPanelViewWithList renders the live matches view with list component.
func RenderMultiPanelViewWithList(width, height int, listModel list.Model, details *api.MatchDetails, liveUpdates []string, sp spinner.Model, loading bool, randomSpinner *RandomCharSpinner, viewLoading bool) string {
	// Handle edge case: if width/height not set, use defaults
	if width <= 0 {
		width = 80
	}
	if height <= 0 {
		height = 24
	}

	// Reserve 3 lines at top for spinner (always reserve to prevent layout shift)
	spinnerHeight := 3
	availableHeight := height - spinnerHeight
	if availableHeight < 10 {
		availableHeight = 10 // Minimum height for panels
	}

	// Render spinner centered in reserved space
	var spinnerArea string
	if viewLoading && randomSpinner != nil {
		spinnerView := randomSpinner.View()
		if spinnerView != "" {
			// Center the spinner horizontally using style with width and alignment
			spinnerStyle := lipgloss.NewStyle().
				Width(width).
				Height(spinnerHeight).
				Align(lipgloss.Center).
				AlignVertical(lipgloss.Center)
			spinnerArea = spinnerStyle.Render(spinnerView)
		} else {
			// Fallback if spinner view is empty
			spinnerStyle := lipgloss.NewStyle().
				Width(width).
				Height(spinnerHeight).
				Align(lipgloss.Center).
				AlignVertical(lipgloss.Center)
			spinnerArea = spinnerStyle.Render("Loading...")
		}
	} else {
		// Reserve space with empty lines - ensure it takes up exactly spinnerHeight lines
		spinnerArea = strings.Repeat("\n", spinnerHeight)
	}

	// Calculate panel dimensions
	leftWidth := width * 35 / 100
	if leftWidth < 25 {
		leftWidth = 25
	}
	rightWidth := width - leftWidth - 1
	if rightWidth < 35 {
		rightWidth = 35
		leftWidth = width - rightWidth - 1
	}

	// Use panelHeight similar to stats view to ensure proper spacing
	panelHeight := availableHeight - 2

	// Render left panel (matches list) - shifted down
	leftPanel := RenderLiveMatchesListPanel(leftWidth, panelHeight, listModel)

	// Render right panel (match details with live updates) - shifted down
	rightPanel := renderMatchDetailsPanel(rightWidth, panelHeight, details, liveUpdates, sp, loading)

	// Create separator
	separatorStyle := lipgloss.NewStyle().
		Foreground(borderColor).
		Height(panelHeight).
		Padding(0, 1)
	separator := separatorStyle.Render("│")

	// Combine panels
	panels := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		separator,
		rightPanel,
	)

	// Combine spinner area and panels - this shifts panels down
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		spinnerArea,
		panels,
	)

	return content
}

// RenderStatsViewWithList renders the stats view with list component.
func RenderStatsViewWithList(width, height int, listModel list.Model, details *api.MatchDetails, randomSpinner *RandomCharSpinner, viewLoading bool) string {
	// Reserve 3 lines at top for spinner (always reserve to prevent layout shift)
	spinnerHeight := 3
	availableHeight := height - spinnerHeight

	// Render spinner centered in reserved space
	var spinnerArea string
	if viewLoading && randomSpinner != nil {
		spinnerView := randomSpinner.View()
		if spinnerView != "" {
			// Center the spinner horizontally
			spinnerArea = lipgloss.Place(width, spinnerHeight, lipgloss.Center, lipgloss.Center, spinnerView)
		} else {
			// Fallback if spinner view is empty
			spinnerArea = lipgloss.Place(width, spinnerHeight, lipgloss.Center, lipgloss.Center, "Loading...")
		}
	} else {
		// Reserve space with empty lines
		spinnerArea = strings.Repeat("\n", spinnerHeight)
	}

	// Calculate panel dimensions
	leftWidth := width * 40 / 100
	if leftWidth < 30 {
		leftWidth = 30
	}
	rightWidth := width - leftWidth - 1
	if rightWidth < 40 {
		rightWidth = 40
		leftWidth = width - rightWidth - 1
	}

	panelHeight := availableHeight - 2

	// Render left panel (finished matches list) - shifted down
	leftPanel := RenderStatsListPanel(leftWidth, panelHeight, listModel)

	// Render right panel (match stats) - shifted down
	rightPanel := renderMatchStatsPanel(rightWidth, panelHeight, details)

	// Create separator
	separatorStyle := lipgloss.NewStyle().
		Foreground(borderColor).
		Height(panelHeight).
		Padding(0, 1)
	separator := separatorStyle.Render("│")

	// Combine panels
	panels := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		separator,
		rightPanel,
	)

	// Combine spinner area and panels
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		spinnerArea,
		panels,
	)

	return content
}
