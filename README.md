# go-rate-limiter

A small Go library providing rate-limiting strategies for HTTP servers, API
clients, and anywhere else you need to bound throughput.

> **Status:** early. Only the **token bucket** strategy is implemented today.
> See [`project.md`](project.md) for the planned strategies (leaky bucket,
> fixed window, sliding window, adaptive).

## Install

```bash
go get github.com/sami-21/go-rate-limiter/rate/tokenbucket
```

Requires Go 1.25+. No external dependencies.

## Quickstart

A single bucket with capacity 3 refilling at 1 token/second:

```go
package main

import (
    "fmt"

    "github.com/sami-21/go-rate-limiter/rate/tokenbucket"
)

func main() {
    b := tokenbucket.New(3, 1)

    for i := 1; i <= 5; i++ {
        if b.Allow() {
            fmt.Println("request", i, "allowed")
        } else {
            fmt.Println("request", i, "blocked")
        }
    }
}
```

The bucket is born full, so the first three calls succeed and the next two
are blocked until tokens refill.

## Per-key limiting

Use `Keyed` for per-user, per-IP, or per-tenant limiting. Stale entries are
evicted on a TTL by an optional janitor goroutine:

```go
k := tokenbucket.NewKeyed(
    100,           // capacity per key
    10,            // tokens per second per key
    time.Hour,     // entry TTL
    time.Minute,   // janitor cleanup interval (0 disables it)
)
defer k.Stop()

if k.Allow("user-123") {
    // serve the request
}
```

`Keyed.Stop()` is idempotent and waits for the janitor goroutine to exit.

## API reference

Full docs render on [pkg.go.dev](https://pkg.go.dev/github.com/sami-21/go-rate-limiter/rate/tokenbucket)
once published. Locally:

```bash
go doc ./rate/tokenbucket
```

## Development

```bash
go test ./...           # run all tests
go test -race ./...     # run with race detector (recommended)
go run ./cmd/demo       # run the demo program
```

## License

MIT — see [LICENSE](LICENSE).
