# Roadmap

The package aims to provide multiple rate-limiting strategies behind a
common interface so users can pick the best fit for their scenario.

Legend: ✅ implemented · 🚧 in progress · ⏳ planned

## Strategies

### 1. Token Bucket ✅
A fixed number of tokens refills over time. Each request consumes a token.
- Use cases: API rate limiting, burst handling.
- Lives at `rate/tokenbucket/`.

### 2. Leaky Bucket ✅
Requests are processed at a fixed rate, and excess requests are queued
(up to a limit).
- Use cases: smoothing out traffic, ensuring steady processing rates.
- Lives at `rate/leakybucket/`.

### 3. Fixed Window ⏳
Counts requests within a fixed timeframe (e.g., per minute). Once the
limit is hit, no more requests are allowed until the next window starts.
- Use cases: simple rate limits with clear reset intervals.

### 4. Sliding Window ⏳
Continuously checks the recent time window (e.g., last 60 seconds). As
time moves, requests "fall out" of the window, allowing new ones.
- Use cases: smoother rate limiting that avoids bursts.

### 5. Adaptive (Dynamic) ⏳
Adjusts limits based on system load, user behavior, or usage tier. Can
increase or decrease limits dynamically.
- Use cases: prioritize premium users, adjust limits during peak traffic,
  slow down potential abuse.

## Goals

1. Easy interface to switch between strategies.
2. User-customizable limits and behavior.
3. Optional adaptive logic hooks.
4. Per-strategy examples.

## Implementation order

1. Token Bucket ✅
2. Leaky Bucket ✅
3. Fixed Window
4. Sliding Window
5. Adaptive
6. Thorough tests and examples (rolling, per strategy)
7. Publish v1
