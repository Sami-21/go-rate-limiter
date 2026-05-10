# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go test ./...                                       # run all tests
go test ./rate/tokenbucket -run TestBucketRefillsOverTime  # run a single test
go run ./cmd/demo                                   # run the demo program
go build ./...                                      # build everything
```

Module path: `go-rate-limiter` (Go 1.25.1, no external dependencies).

## Project vision

`project.md` is the source of truth for scope. The package aims to provide multiple rate-limiting strategies behind a common interface so users can pick what fits:

1. **Token Bucket** — implemented under `rate/tokenbucket/` (the only strategy currently in the tree).
2. **Leaky Bucket** — not yet implemented; expected to add fixed-rate processing with a bounded queue.
3. **Fixed Window** — not yet implemented; counter resets on each window boundary.
4. **Sliding Window** — not yet implemented; rolling time window.
5. **Adaptive / Dynamic** — not yet implemented; adjusts limits based on load, tier, or behavior.

When adding a new strategy, place it in its own subpackage under `rate/` (e.g. `rate/leakybucket/`, `rate/fixedwindow/`, `rate/slidingwindow/`, `rate/adaptive/`) following the shape of `rate/tokenbucket/`. A shared `Allow()`-style interface across strategies has not been introduced yet — when the second strategy lands, decide whether to add one in `rate/` (e.g. `rate.Limiter`) rather than letting each package diverge.

## Architecture

Token-bucket implementation lives in `rate/tokenbucket/` (package `tokenbucket`):

- `rate/tokenbucket/bucket.go` — `Bucket` is the core token bucket, constructed with `New(capacity int, ratePerSecond float64)`. Tokens refill lazily on each `Allow()` call from elapsed wall time (`time.Now().Sub(last).Seconds() * rate`), capped at `capacity`. Synchronization is a single `sync.Mutex`; the bucket holds floating-point tokens but only allows when `tokens >= 1`.
- `rate/tokenbucket/keyed.go` — `Keyed`, constructed with `NewKeyed(capacity, ratePerSecond, ttl)`, wraps a `map[string]*entry` of per-key `Bucket`s, lazily creating one on first `Allow(key)`. Each entry tracks `lastSeen` for TTL-based eviction via `cleanup()`. **Note:** `cleanup()` is unexported and not wired to any goroutine/ticker — the map will grow unbounded until cleanup is invoked. New work in this area should decide whether to start a janitor goroutine, call cleanup inline, or expose it.
- `cmd/demo/main.go` — minimal example that imports `rate/tokenbucket` and exercises `Bucket`.
