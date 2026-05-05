package main

import (
	"fmt"
	"time"

	"go-rate-limiter/rateLimiter"
)

func main() {
	limiter := ratelimiter.NewLimiter(3, 1)

	for i := 1; i <= 10; i++ {
		if limiter.Allow() {
			fmt.Println("request", i, "allowed")
		} else {
			fmt.Println("request", i, "blocked")
		}

		time.Sleep(300 * time.Millisecond)
	}
}
