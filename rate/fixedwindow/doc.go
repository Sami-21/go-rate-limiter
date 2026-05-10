// Package fixedwindow provides a concurrency-safe fixed-window rate limiter.
//
// A Bucket counts accepted requests within a fixed time window. Calls to Allow
// are admitted until the configured limit is reached, then rejected until the
// next window starts and the count resets. This strategy is useful for simple
// limits with clear reset intervals, such as N requests per minute.
package fixedwindow
