# GoRateLimiter

A flexible Go library providing multiple rate-limiting strategies for different needs. This package will implement various rate-limiting techniques so users can choose the best approach for their scenario.

## Planned Features

### 1. Token Bucket
- Description: A fixed number of tokens refills over time. Each request consumes a token.
- Example use cases: API rate limiting, burst handling.

### 2. Leaky Bucket
- Description: Requests are processed at a fixed rate, and excess requests are queued (up to a limit).
- Example use cases: Smoothing out traffic, ensuring steady processing rates.

### 3. Fixed Window
- Description: Counts requests within a fixed timeframe (e.g., per minute). Once the limit is hit, no more requests are allowed until the next window starts.
- Example use cases: Simple rate limits with clear reset intervals.

### 4. Sliding Window
- Description: Continuously checks the recent time window (e.g., last 60 seconds). As time moves, requests “fall out” of the window, allowing new ones.
- Example use cases: Smoother rate limiting that avoids bursts.

### 5. Adaptive (Dynamic) Rate Limiting
- Description: Adjusts limits based on system load, user behavior, or usage tier. Can increase or decrease limits dynamically.
- Example use cases: Prioritize premium users, adjust limits during peak traffic, or slow down potential abuse.

## Goals for the Package

1. Provide an easy interface to switch between strategies.
2. Allow users to customize limits and behavior.
3. Offer optional adaptive logic hooks.
4. Include examples of each method in use.

## Roadmap

1. Implement the core Token Bucket.
2. Add Leaky Bucket with queueing.
3. Implement Fixed Window counting.
4. Implement Sliding Window logic.
5. Add adaptive rate limiting logic.
6. Add thorough tests and examples.
7. Publish and document.

Stay tuned as the project evolves!