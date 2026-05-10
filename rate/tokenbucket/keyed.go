package tokenbucket

import (
	"sync"
	"time"
)

type Keyed struct {
	mu      sync.Mutex
	entries map[string]*entry

	capacity float64
	rate     float64

	ttl time.Duration
}

type entry struct {
	bucket   *Bucket
	lastSeen time.Time
}

func NewKeyed(capacity int, ratePerSecond float64, ttl time.Duration) *Keyed {
	return &Keyed{
		entries:  make(map[string]*entry),
		capacity: float64(capacity),
		rate:     ratePerSecond,
		ttl:      ttl,
	}
}

func (k *Keyed) Allow(key string) bool {
	k.mu.Lock()
	defer k.mu.Unlock()

	now := time.Now()

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

func (k *Keyed) cleanup() {
	now := time.Now()

	for key, e := range k.entries {
		if now.Sub(e.lastSeen) > k.ttl {
			delete(k.entries, key)
		}
	}
}
