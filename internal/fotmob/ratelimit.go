package fotmob

import (
	"sync"
	"time"
)

// RateLimiter provides conservative rate limiting for API requests.
type RateLimiter struct {
	mu              sync.Mutex
	lastRequestTime time.Time
	minInterval     time.Duration
}

// NewRateLimiter creates a new rate limiter.
// minInterval: minimum time between requests (no minimum enforced, use as specified)
func NewRateLimiter(minInterval time.Duration) *RateLimiter {
	// Allow any interval, including very short ones for concurrent requests
	if minInterval < 0 {
		minInterval = 0 // Allow no delay if requested
	}
	return &RateLimiter{
		minInterval: minInterval,
	}
}

// Wait ensures minimum time has passed since last request.
func (rl *RateLimiter) Wait() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastRequestTime)

	if elapsed < rl.minInterval {
		waitTime := rl.minInterval - elapsed
		time.Sleep(waitTime)
	}

	rl.lastRequestTime = time.Now()
}
