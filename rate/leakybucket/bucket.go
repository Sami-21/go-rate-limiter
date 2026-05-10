package leakybucket

import (
	"math"
	"sync"
	"time"
)

// Bucket is a leaky-bucket rate limiter. It is safe for concurrent use.
//
// Accepted requests are scheduled at a steady leak rate. Requests that arrive
// faster than the leak rate wait in a bounded virtual queue; once the in-flight
// request plus queued requests reaches capacity, subsequent calls are rejected
// until enough time has elapsed for queued work to drain.
type Bucket struct {
	mu       sync.Mutex
	capacity int
	interval time.Duration
	next     time.Time
	now      func() time.Time
}

// New returns a Bucket with the given queue limit and leak rate in requests per
// second. queueLimit is the number of requests that may wait behind the request
// currently being processed, so the total burst admitted at one instant is
// queueLimit + 1.
//
// New panics if queueLimit is negative or ratePerSecond is not positive.
func New(queueLimit int, ratePerSecond float64) *Bucket {
	if queueLimit < 0 {
		panic("leakybucket: queueLimit must be >= 0")
	}
	if ratePerSecond <= 0 {
		panic("leakybucket: ratePerSecond must be > 0")
	}

	now := time.Now
	return &Bucket{
		capacity: queueLimit + 1,
		interval: durationPerRequest(ratePerSecond),
		next:     now(),
		now:      now,
	}
}

// Allow reports whether the request can enter the bucket. If the bucket has
// capacity, the request is scheduled at the next leak interval and Allow returns
// true. If the virtual queue is full, Allow returns false and leaves the bucket
// unchanged apart from observing elapsed time.
func (b *Bucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := b.now()
	if !b.next.After(now) {
		b.next = now
	}

	if b.queued(now) >= b.capacity {
		return false
	}

	b.next = b.next.Add(b.interval)
	return true
}

func (b *Bucket) queued(now time.Time) int {
	if !b.next.After(now) {
		return 0
	}

	return int(math.Ceil(float64(b.next.Sub(now)) / float64(b.interval)))
}

func durationPerRequest(ratePerSecond float64) time.Duration {
	seconds := 1 / ratePerSecond
	interval := time.Duration(seconds * float64(time.Second))
	if interval < time.Nanosecond {
		return time.Nanosecond
	}

	return interval
}
