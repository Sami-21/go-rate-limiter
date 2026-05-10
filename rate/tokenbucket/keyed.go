package tokenbucket

import (
	"sync"
	"time"
)

// Keyed is a per-key collection of token buckets, suitable for per-user or
// per-IP rate limiting. Buckets are created lazily on first Allow for a key
// and evicted by TTL via an optional janitor goroutine.
//
// Keyed is safe for concurrent use.
type Keyed struct {
	mu      sync.Mutex
	entries map[string]*entry

	capacity float64
	rate     float64
	ttl      time.Duration

	now      func() time.Time
	done     chan struct{}
	stopOnce sync.Once
	wg       sync.WaitGroup
}

type entry struct {
	bucket   *Bucket
	lastSeen time.Time
}

// NewKeyed returns a Keyed with the given per-key bucket capacity, refill
// rate (tokens per second), and entry TTL. If cleanupInterval is greater
// than zero, a background janitor goroutine is started that evicts entries
// older than ttl at that interval; call Stop to halt it. Pass
// cleanupInterval = 0 to disable the janitor entirely.
//
// NewKeyed panics if capacity, ratePerSecond, or ttl is not positive.
func NewKeyed(capacity int, ratePerSecond float64, ttl, cleanupInterval time.Duration) *Keyed {
	if capacity <= 0 {
		panic("tokenbucket: capacity must be > 0")
	}
	if ratePerSecond <= 0 {
		panic("tokenbucket: ratePerSecond must be > 0")
	}
	if ttl <= 0 {
		panic("tokenbucket: ttl must be > 0")
	}

	k := &Keyed{
		entries:  make(map[string]*entry),
		capacity: float64(capacity),
		rate:     ratePerSecond,
		ttl:      ttl,
		now:      time.Now,
		done:     make(chan struct{}),
	}

	if cleanupInterval > 0 {
		k.wg.Add(1)
		go k.janitor(cleanupInterval)
	}

	return k
}

// Allow reports whether one token is available for the given key. If the
// key has no entry yet, one is created lazily.
func (k *Keyed) Allow(key string) bool {
	k.mu.Lock()
	defer k.mu.Unlock()

	now := k.now()

	e, exists := k.entries[key]
	if !exists {
		e = &entry{
			bucket:   New(int(k.capacity), k.rate),
			lastSeen: now,
		}
		k.entries[key] = e
	}

	e.lastSeen = now

	return e.bucket.Allow()
}

// Stop halts the janitor goroutine started by NewKeyed (if any) and waits
// for it to exit. Stop is idempotent and safe to call from multiple
// goroutines. After Stop, Allow continues to work but stale entries are
// no longer evicted.
func (k *Keyed) Stop() {
	k.stopOnce.Do(func() {
		close(k.done)
	})
	k.wg.Wait()
}

func (k *Keyed) janitor(interval time.Duration) {
	defer k.wg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			k.cleanup()
		case <-k.done:
			return
		}
	}
}

func (k *Keyed) cleanup() {
	k.mu.Lock()
	defer k.mu.Unlock()

	now := k.now()
	for key, e := range k.entries {
		if now.Sub(e.lastSeen) > k.ttl {
			delete(k.entries, key)
		}
	}
}
