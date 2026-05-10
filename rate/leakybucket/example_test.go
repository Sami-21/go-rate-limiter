package leakybucket_test

import (
	"fmt"

	"github.com/sami-21/go-rate-limiter/rate/leakybucket"
)

func ExampleBucket_Allow() {
	b := leakybucket.New(2, 1)

	for i := 1; i <= 4; i++ {
		if b.Allow() {
			fmt.Println("request", i, "accepted")
		} else {
			fmt.Println("request", i, "blocked")
		}
	}

	// Output:
	// request 1 accepted
	// request 2 accepted
	// request 3 accepted
	// request 4 blocked
}
