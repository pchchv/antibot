package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	limit      = 5                // Maximum requests in the period
	interval   = 10 * time.Second // Time period in which there will be a restriction on requests
	staticBody = "Hello, World!"  // Static content that will be rendered on a successful request
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

func (v *visitor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Checking for an IPv4 subnet
	parsedIP := net.ParseIP(ip)
	if parsedIP.To4() == nil {
		http.Error(w, "IPv4 only", http.StatusBadRequest)
		return
	}

	// Restricting the number of requests
	if !v.limiter.allow() {
		http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		return
	}

	// Output static content on successful request
	fmt.Fprint(w, staticBody)
}

func main() {
	limiter := newRateLimiter(limit)
	visitor := newVisitor(limiter)

	server := &http.Server{
		Addr:    ":8080",
		Handler: visitor,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
