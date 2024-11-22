package chapter3

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// Cond
//
// a rendezvous point for goroutines waiting for or announcing the occurrence
// of an event.
//
// Signal finds the goroutine that’s been waiting the longest and notifies that, whereas Broadcast sends a
// signal to all goroutines that are waiting.
func TestCond(t *testing.T) {
	c := sync.NewCond(&sync.Mutex{})    // create a condition with a mutex as the locker
	queue := make([]interface{}, 0, 10) // initialize a slice with length 0 and capacity 10

	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()        // lock the condition's mutex to safely access shared resources
		queue = queue[1:] // simulate dequeuing by removing the front element
		fmt.Println("Removed from queue")
		c.L.Unlock() // unlock the mutex after modifying the queue
		c.Signal()   // notify one goroutine waiting on the condition that an item was removed
	}

	for i := 0; i < 10; i++ {
		c.L.Lock() // lock the condition's mutex to access the queue safely

		for len(queue) == 2 { // pause if the queue has reached a size of 2
			c.Wait() // suspend main goroutine until receiving a signal that space is available
		}
		fmt.Println("Adding to queue")
		queue = append(queue, struct{}{})   // add an item to the queue
		go removeFromQueue(1 * time.Second) // dequeue an item after 1 second in a separate goroutine
		c.L.Unlock()                        // unlock the mutex after adding to the queue
	}
}

type Button struct {
	Clicked *sync.Cond
}

func TestBroadcast(t *testing.T) {
	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}

	// define a convenience function that will allow us to register functions to
	// handle signals from a condition. Each handler is run on its own goroutine, and
	// subscribe will not exit until that goroutine is confirmed to be running.
	subscribe := func(c *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)

		go func() {
			goroutineRunning.Done() // Notify that the goroutine has running
			c.L.Lock()              // Lock the mutex to wait safely
			defer c.L.Unlock()      // Unlock it when done
			c.Wait()                // Wait for the condition to be signaled
			fn()                    // Execute the callback function
		}()

		goroutineRunning.Wait() // Ensure the goroutine starts before moving on
	}

	var clickRegistered sync.WaitGroup
	clickRegistered.Add(3)

	subscribe(button.Clicked, func() {
		fmt.Println("Maximizing window.")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Displaying annoying dialog box!")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Mouse clicked.")
		clickRegistered.Done()
	})

	// Main goroutine send signal to 3 waiting goroutine
	// Simulate a user raising the mouse button from having clicked the application's button
	button.Clicked.Broadcast()

	clickRegistered.Wait()
}
