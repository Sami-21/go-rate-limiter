package slidingwindow

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestBucketAllowsWithinLimit(t *testing.T) {
	b := New(2, time.Minute)

	if !b.Allow() {
		t.Fatal("expected first request to be allowed")
	}

	if !b.Allow() {
		t.Fatal("expected second request to be allowed")
	}
}

func TestBucketBlocksWhenLimitExceeded(t *testing.T) {
	b := New(1, time.Minute)

	if !b.Allow() {
		t.Fatal("expected first request to be allowed")
	}

	if b.Allow() {
		t.Fatal("expected second request to be blocked")
	}
}

func TestBucketExpiresRequestsIndividually(t *testing.T) {
	b, advance := newBucketWithFakeClock(2, time.Minute)

	if !b.Allow() {
		t.Fatal("expected first request to be allowed")
	}

	advance(30 * time.Second)
	if !b.Allow() {
		t.Fatal("expected second request to be allowed")
	}

	advance(30*time.Second - time.Nanosecond)
	if b.Allow() {
		t.Fatal("expected request to remain blocked before oldest event expires")
	}

	advance(time.Nanosecond)
	if !b.Allow() {
		t.Fatal("expected request to be allowed when oldest event expires")
	}

	if b.Allow() {
		t.Fatal("expected request to be blocked while second event is still live")
	}
}

func TestBucketPrunesAllExpiredRequests(t *testing.T) {
	b, advance := newBucketWithFakeClock(3, time.Second)

	for i := range 3 {
		if !b.Allow() {
			t.Fatalf("expected request %d to be allowed", i+1)
		}
	}

	if b.Allow() {
		t.Fatal("expected request to be blocked after limit is reached")
	}

	advance(time.Second)

	for i := range 3 {
		if !b.Allow() {
			t.Fatalf("expected request %d to be allowed after all events expire", i+1)
		}
	}

	if b.Allow() {
		t.Fatal("expected request to be blocked after new limit is reached")
	}
}

func TestBucketParallelAllow(t *testing.T) {
	b := New(100, time.Minute)

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

	if allowed != 100 {
		t.Fatalf("expected exactly 100 requests to be allowed, got %d", allowed)
	}
}

func TestNewPanicsOnInvalidArgs(t *testing.T) {
	cases := []struct {
		name   string
		limit  int
		window time.Duration
	}{
		{"zero limit", 0, time.Minute},
		{"negative limit", -1, time.Minute},
		{"zero window", 1, 0},
		{"negative window", 1, -time.Minute},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("expected New(%d, %v) to panic", tc.limit, tc.window)
				}
			}()
			New(tc.limit, tc.window)
		})
	}
}

func newBucketWithFakeClock(limit int, window time.Duration) (*Bucket, func(time.Duration)) {
	now := time.Unix(0, 0)
	b := New(limit, window)
	b.now = func() time.Time { return now }

	advance := func(d time.Duration) {
		now = now.Add(d)
	}

	return b, advance
}
