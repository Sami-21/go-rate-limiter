package ratelimiter

import (
	"sync"
	"time"
)

type Limiter struct {
	mu       sync.Mutex
	capacity float64
	tokens   float64
	rate     float64
	last     time.Time
}

func NewLimiter(capacity int, ratePerSecond float64) *Limiter {
	return &Limiter{
		capacity: float64(capacity),
		tokens:   float64(capacity),
		rate:     ratePerSecond,
		last:     time.Now(),
	}
}

func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(l.last).Seconds()

	l.tokens += elapsed * l.rate
	if l.tokens > l.capacity {
		l.tokens = l.capacity
	}

	l.last = now

	if l.tokens >= 1 {
		l.tokens--
		return true
	}

	return false
}
