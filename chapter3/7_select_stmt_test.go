package chapter3

import (
	"fmt"
	"testing"
	"time"
)

func TestExample(t *testing.T) {
	var c1, c2 <-chan interface{}
	var c3 chan<- interface{}

	// Unlike switch blocks, case statements in a select
	// block aren’t tested sequentially, and execution won’t automatically fall through if none
	// of the criteria are met.
	//
	// If none of the channels are ready, the entire
	// select statement blocks. Then when one the channels is ready, that operation will
	// proceed, and its corresponding statements
	select {
	case <-c1:
		// Do something
	case <-c2:
		// Do something
	case c3 <- struct{}{}:
		// Do something
	}
}

func TestSelectV1(t *testing.T) {
	start := time.Now()
	c := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(c)
	}()

	fmt.Println("Blocking on read...")

	select {
	case <-c:
		fmt.Printf("Unblocked %v later.\n", time.Since(start))
	}
}

func TestSelectV2(t *testing.T) {
	// What happens when multiple channels have something to read?
	c1 := make(chan interface{})
	close(c1)
	c2 := make(chan interface{})
	close(c2)

	var c1Count, c2Count int
	for i := 1000; i >= 0; i-- {
		select {
		case <-c1:
			c1Count++
		case <-c2:
			c2Count++
		}
	}

	// The Go runtime will perform a pseudo-random uniform selection over the set of case statements.
	// Based on the select statements above, 'each select has an equal chance of being selected'
	//
	// Go runtime cannot know anything about the intent of select statement
	// Go create a random variable in select statemtn, by weighting the chance of each channel
	// being utilized equally and all Go program that utilize select statement will perform
	// well in average case.

	fmt.Printf("c1Count: %d\nc2Count: %d\n", c1Count, c2Count)

	// c1Count: 492
	// c2Count: 509
}

func TestSelectV3(t *testing.T) {
	// What about the second question: what happens if there are never any channels that
	// become ready?

	start := time.Now()
	var c <-chan int
	select {
	case <-c:
	// Never happened
	case <-time.After(1 * time.Second): // timeout handler
		fmt.Println("Timed out")
	default: // default handler (exit select without blocking)
		fmt.Printf("In default after %v\n\n", time.Since(start))
	}
}

func TestSelectV4(t *testing.T) {
	done := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	workCounter := 0
loop:
	for {
		select {
		case <-done:
			break loop
		default:
		}

		workCounter++
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("Achieved %v cycles of work before signalled to stop.\n", workCounter)
}
