// Package slidingwindow provides a concurrency-safe sliding-window rate limiter.
//
// A Bucket counts accepted requests in a trailing time window. Unlike a fixed
// window limiter, individual requests expire as they age out of the window, so
// new requests can be admitted gradually instead of only at reset boundaries.
package slidingwindow
