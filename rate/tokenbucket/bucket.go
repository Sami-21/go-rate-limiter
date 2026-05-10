package tokenbucket

import (
	"sync"
	"time"
)

// Bucket is a token-bucket rate limiter. It is safe for concurrent use.
//
// Tokens refill lazily on each Allow call from the wall time elapsed since
// the previous call, capped at the configured capacity.
type Bucket struct {
	mu       sync.Mutex
	capacity float64
	tokens   float64
	rate     float64
	last     time.Time
	now      func() time.Time
}

// New returns a Bucket with the given capacity (maximum tokens) and refill
// rate in tokens per second. The bucket is created full.
//
// New panics if capacity or ratePerSecond is not positive.
func New(capacity int, ratePerSecond float64) *Bucket {
	if capacity <= 0 {
		panic("tokenbucket: capacity must be > 0")
	}
	if ratePerSecond <= 0 {
		panic("tokenbucket: ratePerSecond must be > 0")
	}

	now := time.Now
	return &Bucket{
		capacity: float64(capacity),
		tokens:   float64(capacity),
		rate:     ratePerSecond,
		last:     now(),
		now:      now,
	}
}

// Allow reports whether one token is available. If so, it consumes the
// token and returns true; otherwise it returns false and the bucket is
// unchanged apart from the lazy refill.
func (b *Bucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := b.now()
	elapsed := now.Sub(b.last).Seconds()

	b.tokens += elapsed * b.rate
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}

	b.last = now

	if b.tokens >= 1 {
		b.tokens--
		return true
	}

	return false
}
