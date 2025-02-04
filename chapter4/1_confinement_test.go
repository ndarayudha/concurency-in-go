package chapter4

import (
	"fmt"
	"testing"
)

// Safe Operation
// 1. Synchronization primitive for sharing memory (e.g., sync.Mutex)
// 2. Synchronization via communicating (e.g., channels)
//
// 3. Immutable data
// 4. Data protected by confinement

func TestAdHoc(t *testing.T) {
	data := make([]int, 4)

	loopData := func(handleData chan<- int) {
		defer close(handleData)
		for i := range data {
			handleData <- i
		}
	}

	handleData := make(chan int)
	go loopData(handleData)

	for num := range handleData {
		fmt.Println(num)
	}
}

func TestLexical(t *testing.T) {
  // Confines write access
	chanOwner := func() <-chan int {
		result := make(chan int, 5)
    go func() {
      defer close(result)
      for i := 0; i <=5; i++ {
        result <- i
      }
    }()

    return result
	}

  // Confines read access
  consumer := func(results <-chan int) {
    for result := range results {
      fmt.Printf("Received %d\n", result)
    }
    fmt.Println("Done receiving")
  }

  result := chanOwner()
  consumer(result)
}
