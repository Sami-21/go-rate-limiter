package tokenbucket

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestKeyedLazyCreation(t *testing.T) {
	k := NewKeyed(2, 1, time.Minute, 0)
	defer k.Stop()

	if got := len(k.entries); got != 0 {
		t.Fatalf("expected no entries before first Allow, got %d", got)
	}

	if !k.Allow("alice") {
		t.Fatal("expected Allow on new key to succeed")
	}

	if got := len(k.entries); got != 1 {
		t.Fatalf("expected 1 entry after first Allow, got %d", got)
	}

	if !k.Allow("alice") {
		t.Fatal("expected second Allow on same key to succeed (capacity 2)")
	}

	if got := len(k.entries); got != 1 {
		t.Fatalf("expected entry to be reused, got %d entries", got)
	}
}

func TestKeyedPerKeyIsolation(t *testing.T) {
	k := NewKeyed(1, 1, time.Minute, 0)
	defer k.Stop()

	if !k.Allow("alice") {
		t.Fatal("expected alice's first request to be allowed")
	}

	if k.Allow("alice") {
		t.Fatal("expected alice's second request to be blocked (capacity 1)")
	}

	if !k.Allow("bob") {
		t.Fatal("expected bob's first request to be allowed independently of alice")
	}
}

func TestKeyedTTLEviction(t *testing.T) {
	k, advance := newKeyedWithFakeClock(1, 1, 5*time.Second)
	defer k.Stop()

	k.Allow("stale")
	k.Allow("fresh")

	advance(3 * time.Second)
	k.Allow("fresh")

	advance(3 * time.Second)

	k.cleanup()

	if _, ok := k.entries["stale"]; ok {
		t.Fatal("expected stale entry to be evicted")
	}

	if _, ok := k.entries["fresh"]; !ok {
		t.Fatal("expected fresh entry to survive eviction")
	}
}

func TestKeyedStopIdempotent(t *testing.T) {
	k := NewKeyed(1, 1, time.Minute, 10*time.Millisecond)

	k.Stop()
	k.Stop()
}

func TestKeyedConcurrentAllow(t *testing.T) {
	k := NewKeyed(100, 1000, time.Minute, 0)
	defer k.Stop()

	const goroutines = 50
	const perGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := range goroutines {
		go func(id int) {
			defer wg.Done()
			key := "user-" + strconv.Itoa(id%10)
			for range perGoroutine {
				k.Allow(key)
			}
		}(i)
	}

	wg.Wait()
}

func TestNewKeyedPanicsOnInvalidArgs(t *testing.T) {
	cases := []struct {
		name     string
		capacity int
		rate     float64
		ttl      time.Duration
	}{
		{"zero capacity", 0, 1, time.Minute},
		{"negative capacity", -1, 1, time.Minute},
		{"zero rate", 1, 0, time.Minute},
		{"negative rate", 1, -1, time.Minute},
		{"zero ttl", 1, 1, 0},
		{"negative ttl", 1, 1, -1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("expected NewKeyed(%d, %v, %v, 0) to panic", tc.capacity, tc.rate, tc.ttl)
				}
			}()
			NewKeyed(tc.capacity, tc.rate, tc.ttl, 0)
		})
	}
}

func newKeyedWithFakeClock(capacity int, ratePerSecond float64, ttl time.Duration) (*Keyed, func(time.Duration)) {
	now := time.Unix(0, 0)
	k := NewKeyed(capacity, ratePerSecond, ttl, 0)
	k.now = func() time.Time { return now }

	advance := func(d time.Duration) {
		now = now.Add(d)
	}

	return k, advance
}
