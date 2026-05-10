package leakybucket

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestBucketAllowsImmediateRequestAndQueue(t *testing.T) {
	b := New(2, 1)

	for i := range 3 {
		if !b.Allow() {
			t.Fatalf("expected request %d to be accepted", i+1)
		}
	}
}

func TestBucketBlocksWhenQueueFull(t *testing.T) {
	b := New(1, 1)

	if !b.Allow() {
		t.Fatal("expected first request to be accepted")
	}

	if !b.Allow() {
		t.Fatal("expected second request to be queued")
	}

	if b.Allow() {
		t.Fatal("expected third request to be blocked while queue is full")
	}
}

func TestBucketDrainsAtFixedRate(t *testing.T) {
	b, advance := newBucketWithFakeClock(1, 2)

	if !b.Allow() {
		t.Fatal("expected first request to be accepted")
	}
	if !b.Allow() {
		t.Fatal("expected second request to be queued")
	}
	if b.Allow() {
		t.Fatal("expected third request to be blocked before one leak interval")
	}

	advance(499 * time.Millisecond)
	if b.Allow() {
		t.Fatal("expected request to stay blocked before the 500ms leak interval")
	}

	advance(1 * time.Millisecond)
	if !b.Allow() {
		t.Fatal("expected request to be accepted after one leak interval")
	}
}

func TestBucketDrainsIdlePeriods(t *testing.T) {
	b, advance := newBucketWithFakeClock(0, 10)

	if !b.Allow() {
		t.Fatal("expected first request to be accepted")
	}
	if b.Allow() {
		t.Fatal("expected second request to be blocked with no queue")
	}

	advance(100 * time.Millisecond)
	if !b.Allow() {
		t.Fatal("expected request to be accepted after idle drain")
	}
}

func TestBucketParallelAllow(t *testing.T) {
	b, _ := newBucketWithFakeClock(100, 1000)

	const goroutines = 50
	const perGoroutine = 200

	var accepted int64
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for range goroutines {
		go func() {
			defer wg.Done()
			for range perGoroutine {
				if b.Allow() {
					atomic.AddInt64(&accepted, 1)
				}
			}
		}()
	}

	wg.Wait()

	if accepted == 0 {
		t.Fatal("expected at least some requests to be accepted")
	}
	if accepted > 101 {
		t.Fatalf("expected at most initial in-flight request plus queue, got %d", accepted)
	}
}

func TestNewPanicsOnInvalidArgs(t *testing.T) {
	cases := []struct {
		name       string
		queueLimit int
		rate       float64
	}{
		{"negative queue limit", -1, 1},
		{"zero rate", 1, 0},
		{"negative rate", 1, -1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("expected New(%d, %v) to panic", tc.queueLimit, tc.rate)
				}
			}()
			New(tc.queueLimit, tc.rate)
		})
	}
}

func newBucketWithFakeClock(queueLimit int, ratePerSecond float64) (*Bucket, func(time.Duration)) {
	now := time.Unix(0, 0)
	b := New(queueLimit, ratePerSecond)
	b.now = func() time.Time { return now }
	b.next = now

	advance := func(d time.Duration) {
		now = now.Add(d)
	}

	return b, advance
}
