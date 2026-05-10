package tokenbucket

import (
	"sync"
	"time"
)

type Bucket struct {
	mu       sync.Mutex
	capacity float64
	tokens   float64
	rate     float64
	last     time.Time
	now      func() time.Time
}

func New(capacity int, ratePerSecond float64) *Bucket {
	now := time.Now
	return &Bucket{
		capacity: float64(capacity),
		tokens:   float64(capacity),
		rate:     ratePerSecond,
		last:     now(),
		now:      now,
	}
}

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
