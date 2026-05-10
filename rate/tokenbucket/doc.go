// Package tokenbucket provides a concurrency-safe token-bucket rate limiter.
//
// A Bucket holds up to a configurable capacity of tokens that refill at a
// fixed rate. Each call to Allow consumes one token if available and reports
// whether the request was admitted. Tokens refill lazily from elapsed wall
// time on each Allow call, so no background goroutine is required for the
// bucket itself.
//
// Keyed wraps a per-key map of buckets for use cases like per-user or per-IP
// limiting, with TTL-based eviction driven by an optional janitor goroutine.
package tokenbucket
