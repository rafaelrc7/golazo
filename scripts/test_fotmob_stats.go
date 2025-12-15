package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/0xjuanma/golazo/internal/fotmob"
)

func main() {
	fmt.Println("Testing FotMob Stats Data Retrieval")
	fmt.Println("===================================\n")

	// Create FotMob client
	client := fotmob.NewClient()

	// Test 1: Recent finished matches (1 day)
	fmt.Println("Test 1: Fetching finished matches from last 1 day...")
	today := time.Now().UTC()
	dateFrom1d := today.AddDate(0, 0, -(1 - 1))
	fmt.Printf("  Date range: %s to %s\n", dateFrom1d.Format("2006-01-02"), today.Format("2006-01-02"))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	matches1d, err := client.RecentFinishedMatches(ctx, 1)
	if err != nil {
		fmt.Printf("❌ Error fetching 1-day matches: %v\n", err)
	} else {
		fmt.Printf("✓ Found %d finished matches (1 day)\n", len(matches1d))
		if len(matches1d) > 0 {
			fmt.Println("\nSample matches (1 day):")
			for i, match := range matches1d {
				if i >= 5 {
					break
				}
				homeScore := "?"
				awayScore := "?"
				if match.HomeScore != nil {
					homeScore = fmt.Sprintf("%d", *match.HomeScore)
				}
				if match.AwayScore != nil {
					awayScore = fmt.Sprintf("%d", *match.AwayScore)
				}
				dateStr := "?"
				if match.MatchTime != nil {
					dateStr = match.MatchTime.Format("Jan 2 15:04")
				}
				fmt.Printf("  %d. %s %s-%s %s (%s) - %s [Status: %s]\n",
					i+1,
					match.HomeTeam.ShortName,
					homeScore,
					awayScore,
					match.AwayTeam.ShortName,
					dateStr,
					match.League.Name,
					match.Status,
				)
			}
		}
	}

	fmt.Println("\n" + strings.Repeat("-", 50) + "\n")

	// Test 2: Recent finished matches (3 days)
	fmt.Println("Test 2: Fetching finished matches from last 3 days...")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	matches3d, err := client.RecentFinishedMatches(ctx2, 3)
	if err != nil {
		fmt.Printf("❌ Error fetching 3-day matches: %v\n", err)
	} else {
		fmt.Printf("✓ Found %d finished matches (3 days)\n", len(matches3d))
	}

	fmt.Println("\n" + strings.Repeat("-", 50) + "\n")

	// Test 3: Recent finished matches (7 days)
	fmt.Println("Test 3: Fetching finished matches from last 7 days...")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel3()

	matches7d, err := client.RecentFinishedMatches(ctx3, 7)
	if err != nil {
		fmt.Printf("❌ Error fetching 7-day matches: %v\n", err)
	} else {
		fmt.Printf("✓ Found %d finished matches (7 days)\n", len(matches7d))
	}

	fmt.Println("\n" + strings.Repeat("-", 50) + "\n")

	// Test 4: Upcoming matches
	fmt.Println("Test 4: Fetching upcoming matches for today...")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel4()

	upcoming, err := client.UpcomingMatches(ctx4)
	if err != nil {
		fmt.Printf("❌ Error fetching upcoming matches: %v\n", err)
	} else {
		fmt.Printf("✓ Found %d upcoming matches\n", len(upcoming))
		if len(upcoming) > 0 {
			fmt.Println("\nSample upcoming matches:")
			for i, match := range upcoming {
				if i >= 5 {
					break
				}
				dateStr := "?"
				if match.MatchTime != nil {
					dateStr = match.MatchTime.Format("Jan 2 15:04")
				}
				fmt.Printf("  %d. %s vs %s (%s) - %s [Status: %s]\n",
					i+1,
					match.HomeTeam.ShortName,
					match.AwayTeam.ShortName,
					dateStr,
					match.League.Name,
					match.Status,
				)
			}
		}
	}

	fmt.Println("\n" + strings.Repeat("-", 50) + "\n")

	// Test 5: Test date range query
	fmt.Println("Test 5: Testing FinishedMatchesByDateRange...")
	today5 := time.Now().UTC()
	yesterday := today5.AddDate(0, 0, -1)
	ctx5, cancel5 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel5()

	dateRangeMatches, err := client.FinishedMatchesByDateRange(ctx5, yesterday, today)
	if err != nil {
		fmt.Printf("❌ Error fetching date range matches: %v\n", err)
	} else {
		fmt.Printf("✓ Found %d finished matches in date range (%s to %s)\n",
			len(dateRangeMatches),
			yesterday.Format("2006-01-02"),
			today5.Format("2006-01-02"),
		)
	}

	fmt.Println("\n" + strings.Repeat("-", 50) + "\n")

	// Test 6: Test match details if we have matches
	if len(matches1d) > 0 {
		fmt.Printf("Test 6: Fetching match details for first match (ID: %d)...\n", matches1d[0].ID)
		ctx6, cancel6 := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel6()

		details, err := client.MatchDetails(ctx6, matches1d[0].ID)
		if err != nil {
			fmt.Printf("❌ Error fetching match details: %v\n", err)
		} else {
			fmt.Printf("✓ Match details retrieved\n")
			fmt.Printf("  Home: %s, Away: %s\n", details.HomeTeam.Name, details.AwayTeam.Name)
			fmt.Printf("  Events: %d\n", len(details.Events))
			if details.Venue != "" {
				fmt.Printf("  Venue: %s\n", details.Venue)
			}
		}
	} else {
		fmt.Println("Test 6: Skipped (no matches available)")
	}

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("Summary:")
	fmt.Printf("  1-day matches: %d\n", len(matches1d))
	fmt.Printf("  3-day matches: %d\n", len(matches3d))
	fmt.Printf("  7-day matches: %d\n", len(matches7d))
	fmt.Printf("  Upcoming matches: %d\n", len(upcoming))
	fmt.Println(strings.Repeat("=", 50))
}
