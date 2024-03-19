package chapter3

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

// fork-join model
//
//                        main
//                          |
//                          |   fork
//                          |------------> Child
//                          |                |
//                          |                |
//                          |  work is done  |
//                          |                |
//            join point    |<---------------X
//                          |                |
//                          |                |
//                          |                |
//                        main             child

func sayHello() {
	fmt.Println("hello")
}

func TestGoroutine(t *testing.T) {
	// Declare Goroutine
	// The "go" statement is how Go perform a fork
	// and the forked threads of execution are goroutine
	//
	// The fork-join model is a logical model of how concurrency is performed.
	// It does describe a C program that calls fork and then wait, but only at a logical level.
	// The fork-join model says nothing about howmemory is managed.

	// Here, the sayHello function will be run on its own goroutine, while the rest of the
	// program continues executing. In this example, there is no join point. The goroutine
	// executing sayHello will simply exit at some undetermined time in the future, and the
	// rest of the program will have already continued executing.
	//
	// However, there is one problem with this example: as written, it’s undetermined
	// whether the sayHello function will ever be run at all. The goroutine will be CREATED
	// and SCHEDULED with Go’s runtime to execute, but it may not actually get a chance to
	// run before the main goroutine exits.
	go sayHello()
	// continue doing other things
}

func TestSynchronization(t *testing.T) {
	// In order to create joint point, it need to synchronized the main goroutine and the
	// sayHello goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	sayWorld := func() {
		wg.Done()
		fmt.Println("world")
	}
	go sayWorld()
	// the wg.Wait() block the main goroutine until the goroutine
	// hosting the sayHello function terminate
	wg.Wait() // this is the join point
}

func TestClosure(t *testing.T) {
	var wg sync.WaitGroup
	yeah := "Yeah"

	wg.Add(1)
	go func() {
		defer wg.Done()
		yeah = "Heyyy" // it will change the value of "yeah"
	}()
	wg.Wait()

	fmt.Println(yeah)
}

// IMPORTANT NOTE
// For go 1.22, the result of code bellow is print all the stirng in a slice (not duplicated)
//
// For go < 1.22 it will print just one of the data in a slice
func TestClosureV2(t *testing.T) {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(salutation)
		}()
	}
	wg.Wait()
}

// For Go < 1.22, the fixed code should be like this
//
// Since the multiple goroutines operate against the same address space,
// it still have to consider about synchronization
//
// It can use either synchronization access to shared memory of the goroutine access,
// or use CSP primitives to SHARE MEMORY BY COMMUNICATION(mentioned in golang docs)
func TestClosureV3(t *testing.T) {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func(salutation string) {
			defer wg.Done()
			fmt.Println(salutation)
		}(salutation)
	}
	wg.Wait()
}

// interesting thing about goroutines:
// - lightweight: only 1kb
// - the garbage collector does nothing to collect goroutines that have been abandoned somehow.
//
// Even like the code bellow
// go func() {
// // <operation that will block forever>
// }()
// Do work
func TestNotGarbageCollected(t *testing.T) {
	memConsumed := func() uint64 {
		runtime.GC() // run the garbage collector
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}

	var c <-chan interface{}
	var wg sync.WaitGroup
	noop := func() {
		wg.Done()
		<-c // never exit
	}

	const numGoroutines = 1e4

	wg.Add(numGoroutines)
	before := memConsumed()
	for i := numGoroutines; i > 0; i-- {
		go noop()
	}
	wg.Wait()
	after := memConsumed()

	fmt.Printf("it consumed %.3fkb\n", float64(after-before)/numGoroutines/1000)
}

// go test -bench=. -run=BenchmarkContextSwitch -cpu=1 chapter3/goroutine_test.go
func BenchmarkContextSwitch(b *testing.B) {
	var wg sync.WaitGroup
	begin := make(chan struct{})
	c := make(chan struct{})

	var token struct{}
	sender := func() {
		defer wg.Done()
		<-begin
		for i := 0; i < b.N; i++ {
			c <- token
		}
	}

	reciver := func() {
		defer wg.Done()
		<-begin
		for i := 0; i < b.N; i++ {
			<-c
		}
	}

  wg.Add(2)
  go sender()
  go reciver()
  b.StartTimer()
  close(begin)
  wg.Wait()
}
