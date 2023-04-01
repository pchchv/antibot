package main

import (
	"sync"
	"time"
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
