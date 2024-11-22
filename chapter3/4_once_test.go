package chapter3

import (
	"fmt"
	"sync"
	"testing"
)

// As the name implies, sync.Once is a type that utilizes some sync primitives internally
// to ensure that only one call to Do ever calls the function passed in—even on different
// goroutines. This is indeed because we wrap the call to increment in a sync.Once Do
// method
func TestIncrement(t *testing.T) {
	var count int

	increment := func() {
		count++
	}

	var once sync.Once

	var increments sync.WaitGroup
	increments.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			defer increments.Done()
			once.Do(increment)
		}()
	}

	increments.Wait()
	fmt.Printf("Count is %d\n", count)
}

// sync.Once only counts the number of times Do is called, 
// not how many times unique functions passed into Do are called.
func TestIncrement2(t *testing.T) {
	var count int
	increment := func() { count++ }
	decrement := func() { count-- }

	var once sync.Once
	once.Do(increment)
	once.Do(decrement)

	fmt.Printf("Count: %d\n", count)
}
