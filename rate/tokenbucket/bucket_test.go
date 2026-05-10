package tokenbucket

import (
	"sync"
	"sync/atomic"
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
	b, advance := newBucketWithFakeClock(1, 10)

	if !b.Allow() {
		t.Fatal("expected first request to be allowed")
	}

	if b.Allow() {
		t.Fatal("expected second request to be blocked")
	}

	advance(120 * time.Millisecond)

	if !b.Allow() {
		t.Fatal("expected request to be allowed after refill")
	}
}

func TestBucketRefillCapsAtCapacity(t *testing.T) {
	b, advance := newBucketWithFakeClock(3, 1)

	for i := range 3 {
		if !b.Allow() {
			t.Fatalf("expected request %d to be allowed (initial capacity)", i+1)
		}
	}

	advance(time.Hour)

	for i := range 3 {
		if !b.Allow() {
			t.Fatalf("expected request %d to be allowed after long refill", i+1)
		}
	}

	if b.Allow() {
		t.Fatal("expected refill to cap at capacity (3) — got a 4th allowed request")
	}
}

func TestBucketParallelAllow(t *testing.T) {
	b := New(100, 1000)

	const goroutines = 50
	const perGoroutine = 200

	var allowed int64
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for range goroutines {
		go func() {
			defer wg.Done()
			for range perGoroutine {
				if b.Allow() {
					atomic.AddInt64(&allowed, 1)
				}
			}
		}()
	}

	wg.Wait()

	if allowed == 0 {
		t.Fatal("expected at least some requests to be allowed")
	}
}

func newBucketWithFakeClock(capacity int, ratePerSecond float64) (*Bucket, func(time.Duration)) {
	now := time.Unix(0, 0)
	b := New(capacity, ratePerSecond)
	b.now = func() time.Time { return now }
	b.last = now

	advance := func(d time.Duration) {
		now = now.Add(d)
	}

	return b, advance
}
