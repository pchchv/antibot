package main

import (
	"net/http"
	"net/http/httptest"
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

func TestVisitor_ServeHTTP(t *testing.T) {
	l := newRateLimiter(2)
	v := newVisitor(l)

	// Create test server
	testServer := httptest.NewServer(v)
	defer testServer.Close()

	// Test successful request
	resp, err := http.Get(testServer.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, resp.StatusCode)
	}

	// Test too many requests
	resp, err = http.Get(testServer.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, resp.StatusCode)
	}

	resp, err = http.Get(testServer.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected status %d; got %d", http.StatusTooManyRequests, resp.StatusCode)
	}

	// Test invalid method
	req, err := http.NewRequest(http.MethodPost, testServer.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d; got %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}

	// Test invalid IP address
	req, err = http.NewRequest(http.MethodGet, testServer.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.RemoteAddr = "notanipaddress"
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d; got %d", http.StatusBadRequest, resp.StatusCode)
	}

	// Test IPv6 address
	req, err = http.NewRequest(http.MethodGet, testServer.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.RemoteAddr = "[::1]:1234"
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d; got %d", http.StatusBadRequest, resp.StatusCode)
	}
}
