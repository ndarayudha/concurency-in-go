package chapter3

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// The sync package contains the concurrency primitives that useful for low level
// memory synchronization

// WaitGroup is useful when dealing with set of concurrent operations to complete
// when either don't care about the result of concurrent operation or with other
// way to collect the result
//
// If those true, it a better way to use channel and select statement instead
func TestWaitGroup(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1) // define the counter for each goroutine to run
	go func() {
		defer wg.Done() // decrement the counter, signal that the task is done
		fmt.Println("Goroutine 1 sleeping...")
		time.Sleep(1 * time.Second)
	}()

	wg.Add(1) // define the counter for each goroutine to run
	go func() {
		defer wg.Done() // decrement the counter, signal that the task is done
		fmt.Println("Goroutine 2 sleeping...")
		time.Sleep(1 * time.Second)
	}()

	// Joint Point
	wg.Wait() // Block the main goroutine, until all goroutine is finished the task (0 counter)
}

func sayHelloWorld(wg *sync.WaitGroup, id int) {
	defer wg.Done()

	fmt.Printf("Hello World %v\n", id)
}

func TestWaitGroupV2(t *testing.T) {
	var wg sync.WaitGroup

	numGoroutine := 5

	wg.Add(numGoroutine)
	for i := 0; i < numGoroutine; i++ {
		go sayHelloWorld(&wg, i+1)
	}

	wg.Wait()
}

// Mutex "mutual exclution" is to create guard for critical section (memory access)
func TestMutext(t *testing.T) {

}
