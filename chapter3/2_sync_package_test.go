package chapter3

import (
	"fmt"
	"math"
	"os"
	"sync"
	"testing"
	"text/tabwriter"
	"time"
)

// The sync package contains the concurrency primitives that useful for low level
// memory synchronization

// WaitGroup is useful when dealing with set of concurrent operations to complete
// when either don't care about the result of concurrent operation or with other
// way to collect the result
//
// If those true, it a better way to use channel and select statement instead
func TestWaitGroupV1(t *testing.T) {
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
//
// A Mutex provides a concurrent-safe way to express exclusive access to these shared resources (memory)
//
// a Mutex shares memory by creating a convention developers must follow to synchronize access to the
// memory. You are responsible for coordinating access to this memory by guarding
// access to it with a mutex.
func TestMutext(t *testing.T) {
	var count int // Shared resources that goroutine need to access

	var lock sync.Mutex

	increment := func() {
		lock.Lock()
		defer lock.Unlock()
		count++
		fmt.Printf("incrementing: %d\n", count)
	}

	decrement := func() {
		lock.Lock()
		defer lock.Unlock()
		count--
		fmt.Printf("decrementing: %d\n", count)
	}

	var arithmetic sync.WaitGroup

	for i := 0; i <= 5; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			increment()
		}()
	}

	for i := 0; i <= 5; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			decrement()
		}()
	}

	arithmetic.Wait()
	fmt.Printf("Aritmethic complete, current value of count is: %d", count)
}

// Critical sections are so named because they reflect a bottleneck in your program. It is
// somewhat expensive to enter and exit a critical section, and so generally people
// attempt to minimize the time spent in critical sections.
//
// One strategy for doing so is to reduce the cross-section of the critical section. There
// may be memory that needs to be shared between multiple concurrent processes, but
// perhaps not all of these processes will read and write to this memory. If this is the
// case, you can take advantage of a different type of mutex: sync.RWMutex.
func TestRWMutex(t *testing.T) {
	producer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		for i := 5; i > 0; i-- {
			l.Lock()
			l.Unlock()
			time.Sleep(1)
		}
	}

	observer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		l.Lock()
		defer l.Unlock()
	}

	test := func(count int, mutex, rwMutex sync.Locker) time.Duration {
		var wg sync.WaitGroup
		wg.Add(count + 1)
		beginTestTime := time.Now()
		go producer(&wg, mutex)
		for i := count; i > 0; i-- {
			go observer(&wg, rwMutex)
		}

		wg.Wait()
		return time.Since(beginTestTime)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	defer tw.Flush()

	var m sync.RWMutex
	fmt.Fprintf(tw, "Readers\tRWMutext\tMutex\n")
	for i := 0; i < 20; i++ {
		count := int(math.Pow(2, float64(i)))
		fmt.Fprintf(
			tw,
			"%d\t%v\t%v\n",
			count,
			test(count, &m, m.RLocker()),
			test(count, &m, &m),
		)
	}
}
