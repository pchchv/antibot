package main

import (
	"sync"
	"time"
)

const (
	limit    = 5                // Maximum requests in the period
	interval = 10 * time.Second // Time period in which there will be a restriction on requests
)

type visitor struct {
	limiter  *rateLimiter // Request counter for this visitor
	lastSeen time.Time    // Time of last request
}

type rateLimiter struct {
	mu     sync.Mutex
	count  int
	limit  int
	window time.Time
}

func newRateLimiter(limit int) *rateLimiter {
	return &rateLimiter{
		count:  0,
		limit:  limit,
		window: time.Now(),
	}
}

func newVisitor(limiter *rateLimiter) *visitor {
	return &visitor{
		limiter:  limiter,
		lastSeen: time.Now(),
	}
}

func (rl *rateLimiter) allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if now.Sub(rl.window) > interval {
		// Counter reset after a time period expires
		rl.count = 0
		rl.window = now
	}

	if rl.count < rl.limit {
		// Counter increase on successful request
		rl.count++
		return true
	}

	return false
}

func main() {}
