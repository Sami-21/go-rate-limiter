# CLAUDE.md

Agent guidance for this repository. For install / usage / dev commands,
see [README.md](README.md). For the strategy roadmap, see [ROADMAP.md](ROADMAP.md).

## Project vision

`ROADMAP.md` is the source of truth for scope. The package aims to provide multiple rate-limiting strategies behind a common interface so users can pick what fits:

1. **Token Bucket** — implemented under `rate/tokenbucket/` (the only strategy currently in the tree).
2. **Leaky Bucket** — not yet implemented; expected to add fixed-rate processing with a bounded queue.
3. **Fixed Window** — not yet implemented; counter resets on each window boundary.
4. **Sliding Window** — not yet implemented; rolling time window.
5. **Adaptive / Dynamic** — not yet implemented; adjusts limits based on load, tier, or behavior.

When adding a new strategy, place it in its own subpackage under `rate/` (e.g. `rate/leakybucket/`, `rate/fixedwindow/`, `rate/slidingwindow/`, `rate/adaptive/`) following the shape of `rate/tokenbucket/`. A shared `Allow()`-style interface across strategies has not been introduced yet — when the second strategy lands, decide whether to add one in `rate/` (e.g. `rate.Limiter`) rather than letting each package diverge.

## Architecture

Token-bucket implementation lives in `rate/tokenbucket/` (package `tokenbucket`):

- `rate/tokenbucket/bucket.go` — `Bucket` is the core token bucket, constructed with `New(capacity int, ratePerSecond float64)`. Tokens refill lazily on each `Allow()` call from elapsed wall time (`elapsed * rate`), capped at `capacity`. Synchronization is a single `sync.Mutex`; the bucket holds floating-point tokens but only allows when `tokens >= 1`. A `now func() time.Time` field is injected (defaults to `time.Now`) so tests can advance time deterministically.
- `rate/tokenbucket/keyed.go` — `Keyed`, constructed with `NewKeyed(capacity, ratePerSecond, ttl, cleanupInterval)`, wraps a `map[string]*entry` of per-key `Bucket`s, lazily creating one on first `Allow(key)`. Each entry tracks `lastSeen` for TTL eviction. When `cleanupInterval > 0`, `NewKeyed` starts a janitor goroutine that calls `cleanup()` on a ticker; halt it via `Stop()` (idempotent, safe to call multiple times). Pass `cleanupInterval = 0` to disable the janitor entirely — useful in tests that drive cleanup manually.
- `rate/tokenbucket/doc.go` — package overview comment for `go doc` / pkg.go.dev.
- `cmd/demo/main.go` — minimal example that imports `rate/tokenbucket` and exercises `Bucket`.

`New` and `NewKeyed` panic on non-positive capacity, rate, or ttl — programmer-error contract, not a runtime condition.
