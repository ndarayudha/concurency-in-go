package chapter3

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"testing"
	"time"
)

// At a high level, a the pool pattern is a way to create and make available a fixed num‐
// ber, or pool, of things for use. It’s commonly used to constrain the creation of things
// that are expensive (e.g., database connections) so that only a fixed number of them
// are ever created, but an indeterminate number of operations can still request access to
// these things. In the case of Go’s sync.Pool, this data type can be safely used by multi‐
// ple goroutines.
//
// Pool’s primary interface is its Get method. When called, Get will first check whether
// there are any available instances within the pool to return to the caller, and if not, call
// its New member variable to create a new one. When finished, callers call Put to place
// the instance they were working with back in the pool for use by other processes.
func TestPool(t *testing.T) {
	myPool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating new instance.")
			return struct{}{}
		},
	}

	// Invoke the New function defined on the pool since instances haven't yet been instantiated
	myPool.Get()
	instance := myPool.Get()

	// Put an instance previously retrieved back in the pool.
	// This increaases the available of instances to one
	myPool.Put(instance)

	// Reuse the instance previously allocated and put it back in the pool.
	// The New function will not be invoked.
	myPool.Get()
}

func TestPoolV2(t *testing.T) {
	// Variable to count how many new objects were created by the pool
	var numCalcsCreated int

	// Define a sync.Pool that manages []byte slices of size 1024
	calcPool := &sync.Pool{
		New: func() interface{} {
			// This function is called whenever the pool needs to create a new object
			numCalcsCreated += 1      // Increment the counter for each new object created
			mem := make([]byte, 1024) // Allocate memory for the new object
			return &mem               // Return a pointer to the memory
		},
	}

	// Pre-fill the pool with 4 objects to reduce the need for new object creation later
	calcPool.Put(calcPool.New()) // Add a new object to the pool
	calcPool.Put(calcPool.New()) // Add another object to the pool
	calcPool.Put(calcPool.New()) // Add another object to the pool
	calcPool.Put(calcPool.New()) // Add another object to the pool

	// Define the number of worker goroutines that will use the pool
	const numWorkers = 1024 * 1024 // 1,048,576 workers
	var wg sync.WaitGroup          // WaitGroup to ensure all workers finish before the program ends
	wg.Add(numWorkers)             // Add the number of workers to the WaitGroup

	// Launch all worker goroutines
	for i := numWorkers; i > 0; i-- {
		go func() {
			defer wg.Done() // Mark this worker as done when it finishes

			// Get an object from the pool
			mem := calcPool.Get().(*[]byte)

			// Return the object back to the pool when done
			defer calcPool.Put(mem)
		}()
	}

	// Wait for all workers to complete their tasks
	wg.Wait()

	// Print the total number of objects created by the pool
	fmt.Printf("%d calculators were created.", numCalcsCreated)
}

func connectToService() interface{} {
	time.Sleep(1 * time.Second)
	return struct{}{}
}

func warmServiceConnCache() *sync.Pool {
	p := &sync.Pool{
		New: connectToService,
	}
	for i := 0; i < 10; i++ {
		p.Put(p.New)
	}
	return p
}

func startNetworkDaemon() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		connPool := warmServiceConnCache()

		server, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		defer server.Close()

		wg.Done()

		for {
			conn, err := server.Accept()
			if err != nil {
				log.Printf("cannot accept connection: %v", err)
				continue
			}
			// connectToService()
			// fmt.Fprintln(conn, "")
			// conn.Close()

			svcConn := connPool.Get()
			fmt.Fprintln(conn, "")
			connPool.Put(svcConn)
			conn.Close()
		}
	}()

	return &wg
}

func init() {
	daemonStarted := startNetworkDaemon()
	daemonStarted.Wait()
}

func BenchmarkNetworkRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			b.Fatalf("cannot dial host: %v", err)
		}
		if _, err := io.ReadAll(conn); err != nil {
			b.Fatalf("cannot read: %v", err)
		}
		conn.Close()
	}
}
