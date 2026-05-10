package fixedwindow

import (
	"sync"
	"time"
)

// Bucket is a fixed-window rate limiter. It is safe for concurrent use.
//
// A Bucket admits up to limit requests during each fixed time window. Once the
// limit is reached, additional requests are rejected until the next window
// starts, at which point the request count resets.
type Bucket struct {
	mu          sync.Mutex
	limit       int
	window      time.Duration
	count       int
	windowStart time.Time
	now         func() time.Time
}

// New returns a Bucket with the given request limit and fixed window duration.
// The first window starts when the bucket is created.
//
// New panics if limit or window is not positive.
func New(limit int, window time.Duration) *Bucket {
	if limit <= 0 {
		panic("fixedwindow: limit must be > 0")
	}
	if window <= 0 {
		panic("fixedwindow: window must be > 0")
	}

	now := time.Now
	return &Bucket{
		limit:       limit,
		window:      window,
		windowStart: now(),
		now:         now,
	}
}

// Allow reports whether the current window has capacity for one request. If
// so, it increments the current window count and returns true. Otherwise it
// returns false. The count resets lazily when Allow observes that a new window
// has started.
func (b *Bucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := b.now()
	if !now.Before(b.windowStart.Add(b.window)) {
		windowsElapsed := int(now.Sub(b.windowStart) / b.window)
		b.windowStart = b.windowStart.Add(time.Duration(windowsElapsed) * b.window)
		b.count = 0
	}

	if b.count >= b.limit {
		return false
	}

	b.count++
	return true
}
