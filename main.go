package main

import (
	"sync"
	"time"
)

const limit = 5 // Maximum requests in the period

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

func main() {}
