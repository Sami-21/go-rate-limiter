package limiter

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
}

func New(capacity int, ratePerSecond float64) *Bucket {
	return &Bucket{
		capacity: float64(capacity),
		tokens:   float64(capacity),
		rate:     ratePerSecond,
		last:     time.Now(),
	}
}

func (b *Bucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
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
