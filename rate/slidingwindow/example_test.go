package slidingwindow_test

import (
	"fmt"
	"time"

	"github.com/sami-21/go-rate-limiter/rate/slidingwindow"
)

// A Bucket admits up to the configured limit during the trailing window.
// Additional requests are denied until earlier accepted requests age out.
func ExampleBucket() {
	b := slidingwindow.New(3, time.Minute)

	fmt.Println(b.Allow())
	fmt.Println(b.Allow())
	fmt.Println(b.Allow())
	fmt.Println(b.Allow())

	// Output:
	// true
	// true
	// true
	// false
}
