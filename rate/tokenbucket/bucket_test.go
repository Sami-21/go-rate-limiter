package tokenbucket

import (
	"testing"
	"time"
)

func TestBucketAllowsWithinCapacity(t *testing.T) {
	b := New(2, 1)

	if !b.Allow() {
		t.Fatal("expected first request to be allowed")
	}

	if !b.Allow() {
		t.Fatal("expected second request to be allowed")
	}
}

func TestBucketBlocksWhenCapacityExceeded(t *testing.T) {
	b := New(1, 1)

	if !b.Allow() {
		t.Fatal("expected first request to be allowed")
	}

	if b.Allow() {
		t.Fatal("expected second request to be blocked")
	}
}

func TestBucketRefillsOverTime(t *testing.T) {
	b := New(1, 10)

	if !b.Allow() {
		t.Fatal("expected first request to be allowed")
	}

	if b.Allow() {
		t.Fatal("expected second request to be blocked")
	}

	time.Sleep(120 * time.Millisecond)

	if !b.Allow() {
		t.Fatal("expected request to be allowed after refill")
	}
}
