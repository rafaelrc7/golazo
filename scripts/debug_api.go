package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/0xjuanma/golazo/internal/data"
)

func main() {
	// Get API key
	apiKey, err := data.FootballDataAPIKey()
	if err != nil {
		fmt.Printf("❌ Error getting API key: %v\n", err)
		fmt.Println("\nMake sure FOOTBALL_DATA_API_KEY is set:")
		fmt.Println("  export FOOTBALL_DATA_API_KEY=\"your-key-here\"")
		os.Exit(1)
	}

	fmt.Printf("✓ API key found (length: %d)\n\n", len(apiKey))

	// Test the API endpoint directly
	baseURL := "https://v3.football.api-sports.io"

	// Test 1: Get all fixtures for today (no filters)
	fmt.Println("Test 1: Fetching all fixtures for today...")
	today := time.Now().Format("2006-01-02")
	url1 := fmt.Sprintf("%s/fixtures?date=%s", baseURL, today)
	testEndpoint(url1, apiKey)

	// Test 2: Get finished matches for today (all leagues)
	fmt.Println("\nTest 2: Fetching finished matches for today (all leagues)...")
	url2 := fmt.Sprintf("%s/fixtures?date=%s&status=FT", baseURL, today)
	testEndpoint(url2, apiKey)

	// Test 3: Get finished matches for today - Premier League only
	fmt.Println("\nTest 3: Fetching finished Premier League matches for today...")
	url3 := fmt.Sprintf("%s/fixtures?date=%s&status=FT&league=39", baseURL, today)
	testEndpoint(url3, apiKey)

	// Test 4: Get finished matches for today - La Liga only
	fmt.Println("\nTest 4: Fetching finished La Liga matches for today...")
	url4 := fmt.Sprintf("%s/fixtures?date=%s&status=FT&league=140", baseURL, today)
	testEndpoint(url4, apiKey)

	// Test 5: Get finished matches for today - Bundesliga only
	fmt.Println("\nTest 5: Fetching finished Bundesliga matches for today...")
	url5 := fmt.Sprintf("%s/fixtures?date=%s&status=FT&league=78", baseURL, today)
	testEndpoint(url5, apiKey)

	// Test 6: Test all supported leagues in one go (simulating actual app behavior)
	fmt.Println("\nTest 6: Testing all supported leagues for today (simulating app behavior)...")
	supportedLeagues := []int{39, 140, 78, 135, 61} // Premier League, La Liga, Bundesliga, Serie A, Ligue 1
	totalMatches := 0
	for _, leagueID := range supportedLeagues {
		url := fmt.Sprintf("%s/fixtures?date=%s&status=FT&league=%d", baseURL, today, leagueID)
		fmt.Printf("  League %d: ", leagueID)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("x-apisports-key", apiKey)
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			bodyBytes, _ := io.ReadAll(resp.Body)
			json.Unmarshal(bodyBytes, &result)
			resp.Body.Close()
			if response, ok := result["response"].([]interface{}); ok {
				fmt.Printf("%d matches\n", len(response))
				totalMatches += len(response)
			} else {
				fmt.Printf("0 matches\n")
			}
		} else {
			fmt.Printf("error\n")
		}
	}
	fmt.Printf("  Total finished matches across all supported leagues: %d\n", totalMatches)
}

func testEndpoint(url string, apiKey string) {
	fmt.Printf("  URL: %s\n", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("  ❌ Failed to create request: %v\n", err)
		return
	}

	req.Header.Set("x-apisports-key", apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("  ❌ Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	fmt.Printf("  Status: %d\n", resp.StatusCode)
	fmt.Printf("  Headers: %+v\n", resp.Header)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("  ❌ Error response: %s\n", bodyStr[:min(1000, len(bodyStr))])
		return
	}

	// Try to parse JSON
	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		fmt.Printf("  ❌ Failed to parse JSON: %v\n", err)
		fmt.Printf("  Response: %s\n", bodyStr[:min(1000, len(bodyStr))])
		return
	}

	// Print full response structure for debugging
	fmt.Printf("  Response keys: %v\n", getKeys(result))

	// Check for errors in response
	if errors, ok := result["errors"].([]interface{}); ok && len(errors) > 0 {
		fmt.Printf("  ⚠ API Errors: %+v\n", errors)
	}

	// Check response structure
	if response, ok := result["response"].([]interface{}); ok {
		fmt.Printf("  ✓ Found %d matches\n", len(response))
		if len(response) > 0 {
			// Show first match summary
			if firstMatch, ok := response[0].(map[string]interface{}); ok {
				fmt.Printf("  Sample match ID: %v\n", firstMatch["fixture"])
				if teams, ok := firstMatch["teams"].(map[string]interface{}); ok {
					if home, ok := teams["home"].(map[string]interface{}); ok {
						if away, ok := teams["away"].(map[string]interface{}); ok {
							fmt.Printf("  Sample: %v vs %v\n", home["name"], away["name"])
						}
					}
				}
			}
		} else {
			fmt.Printf("  ⚠ No matches found in response\n")
		}
	} else {
		fmt.Printf("  ⚠ Unexpected response structure\n")
		fmt.Printf("  Full response: %s\n", bodyStr[:min(500, len(bodyStr))])
	}
}

func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
