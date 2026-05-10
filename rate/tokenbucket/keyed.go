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

func NewKeyed(capacity int, ratePerSecond float64, ttl, cleanupInterval time.Duration) *Keyed {
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
