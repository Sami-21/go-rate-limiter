package tokenbucket_test

import (
	"fmt"
	"time"

	"github.com/sami-21/go-rate-limiter/rate/tokenbucket"
)

// A Bucket starts full. With capacity 3 and a slow refill rate, the first
// three Allow calls succeed and the fourth is denied because no token is
// yet available.
func ExampleBucket() {
	b := tokenbucket.New(3, 1)

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

// Keyed isolates state per key. Exhausting one key's bucket does not
// affect others. Stop halts the janitor goroutine started by NewKeyed.
func ExampleKeyed() {
	k := tokenbucket.NewKeyed(1, 1, time.Hour, 0)
	defer k.Stop()

	fmt.Println(k.Allow("alice"))
	fmt.Println(k.Allow("alice"))
	fmt.Println(k.Allow("bob"))

	// Output:
	// true
	// false
	// true
}
