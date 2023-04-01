package main

import (
	"testing"
	"time"
)

func TestRateLimiter_Allow(t *testing.T) {
	l := newRateLimiter(5)

	// First request should be allowed
	if !l.allow() {
		t.Error("First request should be allowed")
	}

	// Subsequent requests should be allowed until the limit is reached
	for i := 1; i < 5; i++ {
		if !l.allow() {
			t.Errorf("Request %d should be allowed", i)
		}
	}

	// Requests beyond the limit should not be allowed
	if l.allow() {
		t.Error("Request beyond the limit should not be allowed")
	}

	// Wait for the interval to expire, then requests should be allowed again
	time.Sleep(interval)
	if !l.allow() {
		t.Error("Request after interval should be allowed")
	}
}
