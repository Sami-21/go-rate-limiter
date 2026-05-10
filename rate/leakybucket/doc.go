// Package leakybucket provides a concurrency-safe leaky-bucket rate limiter.
//
// A Bucket processes accepted requests at a fixed leak rate. Short bursts can
// be absorbed by a bounded virtual queue, while requests that arrive when the
// queue is full are rejected. This makes the strategy useful when callers want
// to smooth traffic into a steady cadence instead of spending an initial token
// burst immediately.
package leakybucket
