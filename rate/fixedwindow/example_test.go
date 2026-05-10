package fixedwindow_test

import (
	"fmt"
	"time"

	"github.com/sami-21/go-rate-limiter/rate/fixedwindow"
)

// A Bucket admits up to the configured limit during each fixed window. Once
// the limit is reached, additional requests are denied until the window resets.
func ExampleBucket() {
	b := fixedwindow.New(3, time.Minute)

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
