package chapter4

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

// Goroutine come with cost resources, they're not garbage collected by the runtime.

// Goroutine Terminations:
// 1. When it has completed it work
// 2. When it cannot continue its work due to an unrecoverable error
// 3. When it's told to stop working

func TestCleanupV1(t *testing.T) {
	doWork := func(strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(completed)
			for s := range strings {
				// Do something interesting
				fmt.Println(s)
			}
		}()
		return completed
	}

	// Stings channel will never actually gets any strings written onto it, and the goroutine
	// containing doWork will remain in memory for the life time of this process
	// The worst case, the main goroutine could continue to spin up goroutines throughout its life,
	// causing creep in memory utilization.
	doWork(nil)
	// Perhaps more work is done here
	fmt.Println("Done.")
}

func TestCleanupV2(t *testing.T) {
	before := runtime.NumGoroutine()

	// establish a signal between the parent goroutine and it's children to allows the parent
	// to signal cancelation to its children

	doWork := func(done <-chan interface{}, strings <-chan string) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)
			for {
				select {
				case s := <-strings:
					fmt.Println(s)
				case <-done: // handling canelation signal
					return
				}
			}
		}()
		return terminated
	}

	done := make(chan interface{})
	terminated := doWork(done, nil)

	go func() {
		// Cancel operation after 1 second
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()

	// join the goroutine spawned from doWork with the main goroutine
	<-terminated // join point
	fmt.Println("Done.")

	after := runtime.NumGoroutine()
	fmt.Printf("Number of Goroutine, Before %d, After %d", before, after)
}

func TestCleanupV3(t *testing.T) {
	before := runtime.NumGoroutine()

	newRandStream := func() <-chan int {
		randStream := make(chan int)

		go func() {
			defer fmt.Println("newRandStream closure exited.") // never shows
			defer close(randStream)
			for {
				randStream <- rand.Int() // trying to keep writing value
			}
		}()

		return randStream
	}

	randStream := newRandStream()
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream) // but it's no longer being read after 3 iteration
	}

	after := runtime.NumGoroutine()
	fmt.Printf("Number of Goroutine, Before %d, After %d", before, after)
}

func TestCleanupV4(t *testing.T) {
	before := runtime.NumGoroutine()

	newRandStream := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)

		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)
			for {
				select {
				case <-done:
					return
				case randStream <- rand.Int():
				}
			}
		}()

		return randStream
	}

	done := make(chan interface{}) // cancelation channel
	randStream := newRandStream(done)

	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream) // but it's no longer being read after 3 iteration
	}

	close(done) // signal childs goroutine to stop

	fmt.Println("Processing other tasks")
	time.Sleep(1 * time.Second)

	after := runtime.NumGoroutine()
	fmt.Printf("Number of Goroutine, Before %d, After %d", before, after)
}
