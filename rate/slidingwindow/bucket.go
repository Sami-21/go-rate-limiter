package slidingwindow

import (
	"sync"
	"time"
)

// Bucket is a sliding-window rate limiter. It is safe for concurrent use.
//
// A Bucket admits up to limit requests during the trailing window duration.
// Accepted requests are recorded by timestamp and fall out of the window as
// time advances, allowing new requests without waiting for a fixed reset time.
type Bucket struct {
	mu     sync.Mutex
	limit  int
	window time.Duration
	events []time.Time
	now    func() time.Time
}

// New returns a Bucket with the given request limit and trailing window
// duration.
//
// New panics if limit or window is not positive.
func New(limit int, window time.Duration) *Bucket {
	if limit <= 0 {
		panic("slidingwindow: limit must be > 0")
	}
	if window <= 0 {
		panic("slidingwindow: window must be > 0")
	}

	return &Bucket{
		limit:  limit,
		window: window,
		events: make([]time.Time, 0, limit),
		now:    time.Now,
	}
}

// Allow reports whether fewer than limit requests have been accepted within the
// trailing window. If so, it records the current request and returns true.
// Otherwise it returns false. Expired request timestamps are pruned lazily on
// each call.
func (b *Bucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := b.now()
	b.prune(now)

	if len(b.events) >= b.limit {
		return false
	}

	b.events = append(b.events, now)
	return true
}

func (b *Bucket) prune(now time.Time) {
	cutoff := now.Add(-b.window)
	firstLive := 0
	for firstLive < len(b.events) && !b.events[firstLive].After(cutoff) {
		firstLive++
	}

	if firstLive == 0 {
		return
	}
	if firstLive == len(b.events) {
		b.events = b.events[:0]
		return
	}

	copy(b.events, b.events[firstLive:])
	b.events = b.events[:len(b.events)-firstLive]
}
