package chapter1

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// A deadlocked program is one in which all concurrent processes are waiting on one
// another. In this state, the program will never recover without outside intervention.
//
// Deadlock Condition according to Coffman Conditions
// 1. Mutual Exclusion
//      A concurrent process holds exclusive rights to a resource at ony one time
//
// 2. Wait For Condition
//      A concurrent process must simultaneously hold a resource
//      and be waiting for an additional resource.
//
// 3. No Preemption
//      A resource held by a concurrent process can only be released
//      by that process, it fulfills this condition.
//
// 4. Circural Wait
//      A concurrent process (P1) must be waiting on a chain of other concurrent pro‐
//      cesses (P2), which are in turn waiting on it (P1), so it fulfills this final condition
//      too.

func TestDeadlock(t *testing.T) {
	type value struct {
		mu    sync.Mutex
		value int
	}

	var wg sync.WaitGroup
	printSum := func(v1, v2 *value) {
		defer wg.Done()
		v1.mu.Lock()
		defer v1.mu.Unlock()

		time.Sleep(2 * time.Second)
		v2.mu.Lock()
		defer v2.mu.Unlock()

		fmt.Printf("sum=%v\n", v1.value+v2.value)
	}

	var a, b value
	wg.Add(2)
	go printSum(&a, &b)
	go printSum(&b, &a)
	wg.Wait()
}
