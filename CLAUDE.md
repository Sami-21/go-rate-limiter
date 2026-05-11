# CLAUDE.md

Agent guidance for this repository. For install, usage, and development
commands, see [README.md](README.md). For planned features and release scope,
see [ROADMAP.md](ROADMAP.md).

## Project vision

`go-rate-limiter` is a small, dependency-free Go library that provides multiple
concurrency-safe rate-limiting strategies for HTTP servers, API clients, and
other throughput-control use cases.

`ROADMAP.md` is the source of truth for product scope. The currently implemented
strategies are:

1. **Token Bucket** — implemented under `rate/tokenbucket/`; supports bursty
   traffic with lazy token refills.
2. **Leaky Bucket** — implemented under `rate/leakybucket/`; smooths accepted
   requests into a steady cadence with a bounded virtual queue.
3. **Fixed Window** — implemented under `rate/fixedwindow/`; counts accepted
   requests in discrete reset windows.
4. **Sliding Window** — implemented under `rate/slidingwindow/`; counts accepted
   requests in a trailing rolling window.
5. **Adaptive / Dynamic** — planned; should adjust limits based on load, tier,
   behavior, or caller-provided hooks.

All strategies should keep the public `Allow() bool` shape unless the project
intentionally introduces a broader API. If adding shared abstractions, prefer a
small root `rate` package interface such as `type Limiter interface { Allow() bool }`
over coupling the concrete strategy packages to each other.

## Architecture

Each strategy lives in its own subpackage under `rate/` and exposes a `Bucket`
type plus a `New` constructor. Keep packages independent, small, and documented.

- `rate/tokenbucket/`
  - `bucket.go` — `Bucket`, constructed with `New(capacity int, ratePerSecond float64)`. The bucket starts full, refills lazily from elapsed time on each `Allow`, caps tokens at capacity, and uses a `sync.Mutex` for concurrency safety.
  - `keyed.go` — `Keyed`, constructed with `NewKeyed(capacity, ratePerSecond, ttl, cleanupInterval)`, stores per-key token buckets for use cases such as per-user or per-IP limiting. Entries are lazily created and optionally evicted by a janitor goroutine. Call `Stop()` to halt the janitor when cleanup is enabled.
  - `doc.go`, `example_test.go`, and tests provide package documentation and executable examples.
- `rate/leakybucket/`
  - `bucket.go` — `Bucket`, constructed with `New(queueLimit int, ratePerSecond float64)`. The implementation maintains a virtual queue by tracking the next scheduled leak time. `queueLimit` is the number of requests that can wait behind the current in-flight request, so the immediate burst is `queueLimit + 1`.
- `rate/fixedwindow/`
  - `bucket.go` — `Bucket`, constructed with `New(limit int, window time.Duration)`. The implementation tracks the current window start and resets the count lazily when `Allow` observes a new fixed window.
- `rate/slidingwindow/`
  - `bucket.go` — `Bucket`, constructed with `New(limit int, window time.Duration)`. The implementation stores timestamps for accepted events and prunes expired events lazily on each `Allow`.
- `cmd/demo/`
  - Minimal runnable example program. Keep it simple and aligned with README examples.

Constructors currently panic on invalid configuration values such as non-positive
capacity, rate, limit, window, or TTL. Treat that as the package's programmer
error contract unless changing it deliberately across all strategies and docs.

## Go code conventions

Follow idiomatic Go first; avoid clever abstractions unless they make the public
API simpler or the implementation measurably safer.

- Run `gofmt` on every Go file. Do not hand-align code in ways `gofmt` will undo.
- Keep package names short, lowercase, and descriptive (`tokenbucket`, not `tokenBucket`).
- Prefer clear exported API names and concise comments that start with the exported identifier, e.g. `// Bucket is ...` and `// New returns ...`.
- Keep exported surface area minimal. Add exported methods only when they solve a demonstrated user need.
- Return concrete types from constructors when callers benefit from concrete methods, as the current `New` functions do. Accept small interfaces at call sites when useful.
- Use table-driven tests for validation matrices and boundary cases.
- Prefer deterministic tests with injected clocks over sleeps. Avoid time-based flakes.
- Do not introduce external dependencies unless the benefit is substantial and documented.
- Do not add package-level mutable state for limiter behavior. Keep state inside limiter instances.
- Avoid `init` functions unless absolutely necessary.
- Avoid panics outside constructor/configuration validation and impossible internal states.
- Keep error messages and panic messages package-qualified and actionable, e.g. `tokenbucket: capacity must be > 0`.
- Use `time.Duration` for elapsed-time/window APIs and `float64` only where fractional rates are part of the public contract.
- Be careful with integer/float conversions in rate calculations. Add tests for boundaries and rounding behavior.
- Preserve zero goroutine behavior for simple buckets; only `Keyed` with cleanup enabled should start a background goroutine.

## Concurrency and time guidelines

Rate limiters are commonly used in hot paths and concurrent HTTP handlers, so
thread safety and predictable time handling are core requirements.

- All public limiter methods must be safe for concurrent use unless explicitly documented otherwise.
- Protect mutable shared state with `sync.Mutex`, `sync.RWMutex`, atomics, or channels as appropriate. Prefer `sync.Mutex` for simple state machines.
- Do not hold a coarse lock while performing avoidable work that can be done under a narrower per-entry lock.
- Avoid calling user-provided callbacks while holding internal locks.
- Ensure background goroutines have a documented shutdown path and tests that verify shutdown is idempotent.
- Use `time.Now` through an injectable `now func() time.Time` where deterministic tests need to advance time.
- Consider clock rollback and very large elapsed durations when doing refill or prune math.
- Avoid real sleeps in tests unless testing goroutine shutdown or ticker behavior specifically; prefer fake clocks.
- Run `go test -race ./...` after concurrency changes.

## Testing expectations

Every behavior change should include tests. Maintain the current pattern of unit
tests plus executable examples for package documentation.

Minimum checks before committing Go changes:

```bash
go test ./...
go test -race ./...
go vet ./...
```

Add or update tests for:

- constructor validation and panic contracts;
- first request behavior and initial capacity/burst behavior;
- exact boundary conditions around refill intervals and window expiration;
- long idle periods that skip multiple windows or refill beyond capacity;
- concurrent `Allow` calls;
- cleanup and `Stop` behavior for keyed limiters;
- examples whose `// Output:` blocks must remain deterministic.

When changing algorithms, add benchmarks where performance characteristics matter,
especially for hot-path `Allow` methods and high-cardinality keyed use cases.

## Documentation expectations

Documentation is part of the public API for this library.

- Keep README examples synchronized with package examples and `cmd/demo`.
- Update `ROADMAP.md` when feature status changes.
- Update package comments in each `doc.go` when behavior, guarantees, or tradeoffs change.
- Document concurrency safety, constructor panic behavior, and goroutine lifecycle expectations.
- Prefer short examples that compile and run as tests.
- When adding a new strategy, include `doc.go`, `bucket.go`, `bucket_test.go`, and `example_test.go` in the new package.

## Workflow guidance for agents

1. Inspect relevant files before editing; do not assume `CLAUDE.md` is more current than code, README, or roadmap.
2. Keep changes focused and small. Avoid drive-by refactors when fixing a targeted issue.
3. Run the narrowest useful tests during development, then the full checks listed above before finalizing code changes.
4. If a check cannot run because of an environment limitation, report it clearly.
5. Before committing, verify `git diff` and `git status --short`.
6. Commit changes on the current branch with a concise imperative commit message.
