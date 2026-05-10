package main

import (
	"fmt"
	"time"

	"go-rate-limiter/rate/tokenbucket"
)

func main() {
	b := tokenbucket.New(3, 1)

	for i := 1; i <= 10; i++ {
		if b.Allow() {
			fmt.Println("request", i, "allowed")
		} else {
			fmt.Println("request", i, "blocked")
		}

		time.Sleep(300 * time.Millisecond)
	}
}
