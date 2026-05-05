package ratelimiter

import (
	"testing"
	"time"
)

func TestLimiterAllowsWithinCapacity(t *testing.T) {
	limiter := NewLimiter(2, 1)

	if !limiter.Allow() {
		t.Fatal("expected first request to be allowed")
	}

	if !limiter.Allow() {
		t.Fatal("expected second request to be allowed")
	}
}

func TestLimiterBlocksWhenCapacityExceeded(t *testing.T) {
	limiter := NewLimiter(1, 1)

	if !limiter.Allow() {
		t.Fatal("expected first request to be allowed")
	}

	if limiter.Allow() {
		t.Fatal("expected second request to be blocked")
	}
}

func TestLimiterRefillsOverTime(t *testing.T) {
	limiter := NewLimiter(1, 10)

	if !limiter.Allow() {
		t.Fatal("expected first request to be allowed")
	}

	if limiter.Allow() {
		t.Fatal("expected second request to be blocked")
	}

	time.Sleep(120 * time.Millisecond)

	if !limiter.Allow() {
		t.Fatal("expected request to be allowed after refill")
	}
}
