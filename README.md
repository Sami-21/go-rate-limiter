# go-rate-limiter

A small Go library providing rate-limiting strategies for HTTP servers, API
clients, and anywhere else you need to bound throughput.

> **Status:** early. The **token bucket**, **leaky bucket**, and **fixed
> window** strategies are implemented today. See [`ROADMAP.md`](ROADMAP.md) for
> the planned strategies (sliding window, adaptive).

## Install

```bash
go get github.com/sami-21/go-rate-limiter/rate/tokenbucket
go get github.com/sami-21/go-rate-limiter/rate/leakybucket
go get github.com/sami-21/go-rate-limiter/rate/fixedwindow
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

## Leaky bucket

Use `leakybucket` when you want to smooth accepted traffic into a steady
cadence. A queue limit of 2 accepts one request immediately and queues two more
at the configured leak rate; additional requests are blocked until the queue
drains:

```go
package main

import (
    "fmt"

    "github.com/sami-21/go-rate-limiter/rate/leakybucket"
)

func main() {
    b := leakybucket.New(2, 1) // queue up to 2 requests, leak 1 request/second

    for i := 1; i <= 4; i++ {
        if b.Allow() {
            fmt.Println("request", i, "accepted")
        } else {
            fmt.Println("request", i, "blocked")
        }
    }
}
```

## Fixed window

Use `fixedwindow` for simple limits with clear reset intervals, such as 3
requests per minute. The counter resets when the next fixed window starts:

```go
package main

import (
    "fmt"
    "time"

    "github.com/sami-21/go-rate-limiter/rate/fixedwindow"
)

func main() {
    b := fixedwindow.New(3, time.Minute)

    for i := 1; i <= 4; i++ {
        if b.Allow() {
            fmt.Println("request", i, "allowed")
        } else {
            fmt.Println("request", i, "blocked")
        }
    }
}
```

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

Full docs render on pkg.go.dev once published:

- [`rate/tokenbucket`](https://pkg.go.dev/github.com/sami-21/go-rate-limiter/rate/tokenbucket)
- [`rate/leakybucket`](https://pkg.go.dev/github.com/sami-21/go-rate-limiter/rate/leakybucket)
- [`rate/fixedwindow`](https://pkg.go.dev/github.com/sami-21/go-rate-limiter/rate/fixedwindow)

Locally:

```bash
go doc ./rate/tokenbucket
go doc ./rate/leakybucket
go doc ./rate/fixedwindow
```

## Development

```bash
go test ./...           # run all tests
go test -race ./...     # run with race detector (recommended)
go run ./cmd/demo       # run the demo program
```

## License

MIT — see [LICENSE](LICENSE).
