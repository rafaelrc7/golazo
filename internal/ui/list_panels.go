package ui

import (
	"fmt"
	"strings"

	"github.com/0xjuanma/golazo/internal/api"
	"github.com/0xjuanma/golazo/internal/constants"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

// RenderLiveMatchesListPanel renders the left panel using bubbletea list component.
// Note: listModel is passed by value, so SetSize must be called before this function.
// Uses Neon design with Golazo red/cyan theme.
func RenderLiveMatchesListPanel(width, height int, listModel list.Model) string {
	// Wrap list in panel with neon styling
	title := neonPanelTitleStyle.Width(width - 6).Render(constants.PanelLiveMatches)
	listView := listModel.View()

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		listView,
	)

	panel := neonPanelStyle.
		Width(width).
		Height(height).
		Render(content)

	return panel
}

// RenderStatsListPanel renders the left panel for stats view using bubbletea list component.
// Note: listModel is passed by value, so SetSize must be called before this function.
// Uses Neon design with Golazo red/cyan theme.
// List titles are only shown when there are items. Empty lists show gray messages instead.
// For 1-day view, shows both finished and upcoming lists stacked vertically.
func RenderStatsListPanel(width, height int, finishedList list.Model, upcomingList list.Model, dateRange int) string {
	// Render date range selector with neon styling
	dateSelector := renderDateRangeSelector(width-6, dateRange)

	emptyStyle := neonEmptyStyle.Width(width - 6)

	var finishedListView string
	finishedItems := finishedList.Items()
	if len(finishedItems) == 0 {
		// No items - show empty message, no list title
		finishedListView = emptyStyle.Render(constants.EmptyNoFinishedMatches + "\n\nTry selecting a different date range (h/l keys)")
	} else {
		// Has items - show list (which includes its title)
		finishedListView = finishedList.View()
	}

	// For 1-day view, show both lists stacked vertically
	if dateRange == 1 {
		var upcomingListView string
		upcomingItems := upcomingList.Items()
		if len(upcomingItems) == 0 {
			// No upcoming matches - show empty message, no list title
			upcomingListView = emptyStyle.Render("No upcoming matches scheduled for today")
		} else {
			// Has items - show list (which includes its title)
			upcomingListView = upcomingList.View()
		}

		// Combine both lists with date selector
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			dateSelector,
			"",
			finishedListView,
			"",
			upcomingListView,
		)
		panel := neonPanelStyle.
			Width(width).
			Height(height).
			Render(content)
		return panel
	}

	// For 3-day view, only show finished matches
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		dateSelector,
		"",
		finishedListView,
	)

	panel := neonPanelStyle.
		Width(width).
		Height(height).
		Render(content)

	return panel
}

// renderDateRangeSelector renders a horizontal date range selector (Today, 3d).
func renderDateRangeSelector(width int, selected int) string {
	options := []struct {
		days  int
		label string
	}{
		{1, "Today"},
		{3, "3d"},
	}

	items := make([]string, 0, len(options))
	for _, opt := range options {
		if opt.days == selected {
			// Selected option - neon red
			item := neonDateSelectedStyle.Render(opt.label)
			items = append(items, item)
		} else {
			// Unselected option - dim
			item := neonDateUnselectedStyle.Render(opt.label)
			items = append(items, item)
		}
	}

	// Join items with separator
	separator := "  "
	selector := strings.Join(items, separator)

	// Center the selector
	selectorStyle := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Padding(0, 1)

	return selectorStyle.Render(selector)
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

	// Create separator with neon red accent
	separatorStyle := neonSeparatorStyle.Height(panelHeight)
	separator := separatorStyle.Render("┃")

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
// Rebuilt to match live view structure exactly: spinner at top, left panel (matches), right panel (details).
func RenderStatsViewWithList(width, height int, finishedList list.Model, upcomingList list.Model, details *api.MatchDetails, randomSpinner *RandomCharSpinner, viewLoading bool, dateRange int) string {
	// Handle edge case: if width/height not set, use defaults
	if width <= 0 {
		width = 80
	}
	if height <= 0 {
		height = 24
	}

	// Reserve 3 lines at top for spinner (always reserve to prevent layout shift)
	// Match live view exactly
	spinnerHeight := 3
	availableHeight := height - spinnerHeight
	if availableHeight < 10 {
		availableHeight = 10 // Minimum height for panels
	}

	// Render spinner centered in reserved space - match live view exactly
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

	// Calculate panel dimensions - match live view exactly (35% left, 65% right)
	leftWidth := width * 35 / 100
	if leftWidth < 25 {
		leftWidth = 25
	}
	rightWidth := width - leftWidth - 1
	if rightWidth < 35 {
		rightWidth = 35
		leftWidth = width - rightWidth - 1
	}

	// Use panelHeight similar to live view to ensure proper spacing
	panelHeight := availableHeight - 2

	// Render left panel (finished matches list) - match live view structure
	// For 1-day view, combine finished and upcoming lists vertically
	leftPanel := RenderStatsListPanel(leftWidth, panelHeight, finishedList, upcomingList, dateRange)

	// Render right panel (match details) - use dedicated stats panel renderer
	rightPanel := renderStatsMatchDetailsPanel(rightWidth, panelHeight, details)

	// Create separator with neon red accent
	separatorStyle := neonSeparatorStyle.Height(panelHeight)
	separator := separatorStyle.Render("┃")

	// Combine panels
	panels := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		separator,
		rightPanel,
	)

	// Combine spinner area and panels - this shifts panels down
	// Match live view exactly - use lipgloss.Left
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		spinnerArea,
		panels,
	)

	return content
}

// renderStatsMatchDetailsPanel renders the right panel for stats view with match details.
// Uses Neon design with Golazo red/cyan theme.
// Displays expanded match information including statistics, lineups, and more.
func renderStatsMatchDetailsPanel(width, height int, details *api.MatchDetails) string {
	if details == nil {
		emptyMessage := neonDimStyle.
			Align(lipgloss.Center).
			Width(width - 6).
			PaddingTop(height / 4).
			Render("Select a match to view details")

		return neonPanelCyanStyle.
			Width(width).
			Height(height).
			Render(emptyMessage)
	}

	contentWidth := width - 6 // Account for border padding
	var lines []string

	// Team names
	homeTeam := details.HomeTeam.ShortName
	if homeTeam == "" {
		homeTeam = details.HomeTeam.Name
	}
	awayTeam := details.AwayTeam.ShortName
	if awayTeam == "" {
		awayTeam = details.AwayTeam.Name
	}

	// ═══════════════════════════════════════════════
	// MATCH HEADER
	// ═══════════════════════════════════════════════
	lines = append(lines, neonHeaderStyle.Render("Match Info"))
	lines = append(lines, "")

	// Score line - centered with large emphasis
	var scoreDisplay string
	if details.HomeScore != nil && details.AwayScore != nil {
		scoreDisplay = fmt.Sprintf("%s  %s  %s",
			neonTeamStyle.Render(homeTeam),
			neonScoreStyle.Render(fmt.Sprintf("%d - %d", *details.HomeScore, *details.AwayScore)),
			neonTeamStyle.Render(awayTeam))
	} else {
		scoreDisplay = fmt.Sprintf("%s  vs  %s",
			neonTeamStyle.Render(homeTeam),
			neonTeamStyle.Render(awayTeam))
	}
	lines = append(lines, lipgloss.NewStyle().Width(contentWidth).Align(lipgloss.Center).Render(scoreDisplay))

	// Status + Half-time in one line
	var statusLine string
	switch details.Status {
	case api.MatchStatusFinished:
		statusLine = neonFinishedStyle.Render("FT")
	case api.MatchStatusLive:
		if details.LiveTime != nil {
			statusLine = neonLiveStyle.Render(*details.LiveTime)
		} else {
			statusLine = neonLiveStyle.Render("LIVE")
		}
	default:
		statusLine = neonDimStyle.Render(string(details.Status))
	}
	if details.HalfTimeScore != nil && details.HalfTimeScore.Home != nil && details.HalfTimeScore.Away != nil {
		statusLine += neonDimStyle.Render(fmt.Sprintf("  (HT: %d-%d)", *details.HalfTimeScore.Home, *details.HalfTimeScore.Away))
	}
	lines = append(lines, lipgloss.NewStyle().Width(contentWidth).Align(lipgloss.Center).Render(statusLine))
	lines = append(lines, "")

	// Match context row
	if details.League.Name != "" {
		lines = append(lines, neonLabelStyle.Render("League:      ")+neonValueStyle.Render(details.League.Name))
	}
	if details.Venue != "" {
		lines = append(lines, neonLabelStyle.Render("Venue:       ")+neonValueStyle.Render(truncateString(details.Venue, contentWidth-14)))
	}
	if details.MatchTime != nil {
		lines = append(lines, neonLabelStyle.Render("Date:        ")+neonValueStyle.Render(details.MatchTime.Format("02 Jan 2006, 15:04")))
	}
	if details.Referee != "" {
		lines = append(lines, neonLabelStyle.Render("Referee:     ")+neonValueStyle.Render(details.Referee))
	}
	if details.Attendance > 0 {
		lines = append(lines, neonLabelStyle.Render("Attendance:  ")+neonValueStyle.Render(formatNumber(details.Attendance)))
	}

	// ═══════════════════════════════════════════════
	// GOALS TIMELINE
	// ═══════════════════════════════════════════════
	var homeGoals, awayGoals []api.MatchEvent
	for _, event := range details.Events {
		if event.Type == "goal" {
			if event.Team.ID == details.HomeTeam.ID {
				homeGoals = append(homeGoals, event)
			} else {
				awayGoals = append(awayGoals, event)
			}
		}
	}

	if len(homeGoals) > 0 || len(awayGoals) > 0 {
		lines = append(lines, "")
		lines = append(lines, neonHeaderStyle.Render("Goals"))

		if len(homeGoals) > 0 {
			lines = append(lines, neonTeamStyle.Render(homeTeam))
			for _, g := range homeGoals {
				goalLine := renderGoalLine(g, contentWidth-2)
				lines = append(lines, "  "+goalLine)
			}
		}

		if len(awayGoals) > 0 {
			lines = append(lines, neonTeamStyle.Render(awayTeam))
			for _, g := range awayGoals {
				goalLine := renderGoalLine(g, contentWidth-2)
				lines = append(lines, "  "+goalLine)
			}
		}
	}

	// ═══════════════════════════════════════════════
	// CARDS SUMMARY
	// ═══════════════════════════════════════════════
	var homeYellow, awayYellow, homeRed, awayRed int
	for _, event := range details.Events {
		isHome := event.Team.ID == details.HomeTeam.ID
		// Check for card events - FotMob uses lowercase "card" type
		if event.Type == "card" {
			if event.EventType != nil {
				switch *event.EventType {
				case "yellow", "yellowcard":
					if isHome {
						homeYellow++
					} else {
						awayYellow++
					}
				case "red", "redcard", "secondyellow":
					if isHome {
						homeRed++
					} else {
						awayRed++
					}
				}
			}
		}
	}

	if homeYellow > 0 || awayYellow > 0 || homeRed > 0 || awayRed > 0 {
		lines = append(lines, "")
		lines = append(lines, neonHeaderStyle.Render("Cards"))

		// Get short team names (3 chars max)
		homeShort := homeTeam
		if len(homeShort) > 3 {
			homeShort = homeShort[:3]
		}
		awayShort := awayTeam
		if len(awayShort) > 3 {
			awayShort = awayShort[:3]
		}

		cardLine := fmt.Sprintf("  Yellow: %s %d - %d %s",
			neonDimStyle.Render(homeShort),
			homeYellow,
			awayYellow,
			neonDimStyle.Render(awayShort))
		lines = append(lines, neonTeamStyle.Render(cardLine))
		if homeRed > 0 || awayRed > 0 {
			redLine := fmt.Sprintf("  Red:    %s %d - %d %s",
				neonDimStyle.Render(homeShort),
				homeRed,
				awayRed,
				neonDimStyle.Render(awayShort))
			lines = append(lines, neonLiveStyle.Render(redLine))
		}
	}

	// ═══════════════════════════════════════════════
	// MATCH STATISTICS
	// ═══════════════════════════════════════════════
	if len(details.Statistics) > 0 {
		lines = append(lines, "")
		lines = append(lines, neonHeaderStyle.Render("Statistics"))

		// Priority labels/keys to show first (check both key and label)
		priorityPatterns := []string{
			"possession", "ball possession",
			"expected_goals", "xg", "expected goals",
			"total_shots", "shots", "total shots",
			"shots_on_target", "on target", "shots on target",
			"passes", "accurate passes",
			"pass_accuracy", "pass accuracy", "passes %",
		}
		shownIndices := make(map[int]bool)

		// Show priority stats first
		for _, pattern := range priorityPatterns {
			for i, stat := range details.Statistics {
				if shownIndices[i] {
					continue
				}
				keyLower := strings.ToLower(stat.Key)
				labelLower := strings.ToLower(stat.Label)
				if strings.Contains(keyLower, pattern) || strings.Contains(labelLower, pattern) {
					statLine := renderStatLine(stat, contentWidth-2)
					lines = append(lines, "  "+statLine)
					shownIndices[i] = true
					break
				}
			}
		}

		// Show remaining stats (limit to avoid overflow)
		remaining := 0
		maxRemaining := 6 // Show more stats if space allows
		for i, stat := range details.Statistics {
			if !shownIndices[i] && remaining < maxRemaining {
				statLine := renderStatLine(stat, contentWidth-2)
				lines = append(lines, "  "+statLine)
				shownIndices[i] = true
				remaining++
			}
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)

	return neonPanelCyanStyle.
		Width(width).
		Height(height).
		Render(content)
}

// renderGoalLine renders a single goal with scorer, minute, and assist
func renderGoalLine(g api.MatchEvent, maxWidth int) string {
	player := "Unknown"
	if g.Player != nil {
		player = *g.Player
	}

	minuteStr := neonScoreStyle.Render(fmt.Sprintf("%d'", g.Minute))
	playerStr := neonValueStyle.Render(truncateString(player, maxWidth-10))

	line := fmt.Sprintf("%s %s", minuteStr, playerStr)

	// Add assist if available
	if g.Assist != nil && *g.Assist != "" {
		line += neonDimStyle.Render(fmt.Sprintf(" (%s)", truncateString(*g.Assist, 15)))
	}

	return line
}

// renderStatLine renders a single statistic comparing home vs away
func renderStatLine(stat api.MatchStatistic, maxWidth int) string {
	// Format: "Home Value | Label | Away Value"
	labelWidth := 16
	valueWidth := (maxWidth - labelWidth - 3) / 2

	label := truncateString(stat.Label, labelWidth)
	homeVal := truncateString(stat.HomeValue, valueWidth)
	awayVal := truncateString(stat.AwayValue, valueWidth)

	return fmt.Sprintf("%s %s %s",
		neonValueStyle.Render(fmt.Sprintf("%*s", valueWidth, homeVal)),
		neonDimStyle.Render(fmt.Sprintf("%-*s", labelWidth, label)),
		neonValueStyle.Render(fmt.Sprintf("%-*s", valueWidth, awayVal)))
}

// truncateString truncates a string to maxLen, adding "..." if truncated
func truncateString(s string, maxLen int) string {
	if maxLen <= 3 {
		return s
	}
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// formatNumber formats a number with thousand separators
func formatNumber(n int) string {
	s := fmt.Sprintf("%d", n)
	if n < 1000 {
		return s
	}

	// Insert commas from right to left
	result := ""
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result += ","
		}
		result += string(c)
	}
	return result
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
