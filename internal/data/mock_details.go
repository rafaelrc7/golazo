package data

import (
	"time"

	"github.com/0xjuanma/golazo/internal/api"
)

// MockMatchDetails returns detailed match information for a specific match ID.
func MockMatchDetails(matchID int) (*api.MatchDetails, error) {
	matches, err := MockMatches()
	if err != nil {
		return nil, err
	}

	// Find the match
	var match *api.Match
	for i := range matches {
		if matches[i].ID == matchID {
			match = &matches[i]
			break
		}
	}

	if match == nil {
		return nil, nil
	}

	// Generate events based on match ID
	events := generateMockEvents(matchID, *match)

	return &api.MatchDetails{
		Match:  *match,
		Events: events,
	}, nil
}

// generateMockEvents generates mock events for a match.
func generateMockEvents(matchID int, match api.Match) []api.MatchEvent {
	events := []api.MatchEvent{}

	switch matchID {
	case 1: // Man Utd vs Liverpool
		events = []api.MatchEvent{
			{ID: 1, Minute: 12, Type: "goal", Team: match.HomeTeam, Player: stringPtr("Rashford"), Timestamp: time.Now()},
			{ID: 2, Minute: 34, Type: "goal", Team: match.AwayTeam, Player: stringPtr("Salah"), Timestamp: time.Now()},
			{ID: 3, Minute: 45, Type: "card", Team: match.HomeTeam, Player: stringPtr("Casemiro"), EventType: stringPtr("yellow"), Timestamp: time.Now()},
			{ID: 4, Minute: 56, Type: "goal", Team: match.HomeTeam, Player: stringPtr("Fernandes"), Assist: stringPtr("Rashford"), Timestamp: time.Now()},
		}
	case 2: // Real Madrid vs Barcelona
		events = []api.MatchEvent{
			{ID: 5, Minute: 8, Type: "goal", Team: match.AwayTeam, Player: stringPtr("Lewandowski"), Timestamp: time.Now()},
			{ID: 6, Minute: 23, Type: "goal", Team: match.HomeTeam, Player: stringPtr("Vinicius Jr"), Timestamp: time.Now()},
		}
	case 3: // AC Milan vs Inter
		events = []api.MatchEvent{
			{ID: 7, Minute: 5, Type: "goal", Team: match.HomeTeam, Player: stringPtr("Leao"), Timestamp: time.Now()},
			{ID: 8, Minute: 18, Type: "goal", Team: match.AwayTeam, Player: stringPtr("Martinez"), Timestamp: time.Now()},
			{ID: 9, Minute: 32, Type: "goal", Team: match.HomeTeam, Player: stringPtr("Giroud"), Timestamp: time.Now()},
			{ID: 10, Minute: 45, Type: "card", Team: match.AwayTeam, Player: stringPtr("Barella"), EventType: stringPtr("yellow"), Timestamp: time.Now()},
			{ID: 11, Minute: 67, Type: "goal", Team: match.AwayTeam, Player: stringPtr("Lautaro"), Timestamp: time.Now()},
			{ID: 12, Minute: 78, Type: "goal", Team: match.HomeTeam, Player: stringPtr("Pulisic"), Timestamp: time.Now()},
			{ID: 13, Minute: 89, Type: "card", Team: match.HomeTeam, Player: stringPtr("Theo"), EventType: stringPtr("red"), Timestamp: time.Now()},
		}
	case 4: // Arsenal vs Chelsea (not started)
		events = []api.MatchEvent{}
	}

	return events
}

func stringPtr(s string) *string {
	return &s
}
