package chapter1

import "testing"

func TestAtomic(t *testing.T) {
  var i int

  // SOMETHING WILL ATOMIC IN SOME CONTEXT BUT NOT IN SOME CONTEXT

  i++
  // Retrieve the value of i
  // Increment the value of i
  // Store the value of i

  // If the context is a program with no concurrent process, the code above is atomic within that context
  // if the context is a goroutine that DOESN'T expose i to other goroutines, the code is atomic
  
  // IF SOMETHING IS ATOMIC, IMPLICITY IT IS SAFE WITHIN CONCURRENT CONTEXT
}
