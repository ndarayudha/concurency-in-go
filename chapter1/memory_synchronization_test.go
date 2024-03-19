package chapter1

import (
	"fmt"
	"sync"
	"testing"
)

func TestMemorySynchronization(t *testing.T) {
  var memoryAccess sync.Mutex

  var data int

  // Two concurrent process try to access data, (main and goroutine)
  // Not atomic code

  go func() {
    memoryAccess.Lock() // guarding start
    data++
    memoryAccess.Unlock() // guarding end
  }()

  memoryAccess.Lock() // guarding start
  if data == 0 {
    fmt.Printf("the value is 0\n") 
  } else {
    fmt.Printf("The value is %v\n", data)
  }
  memoryAccess.Unlock() // guarding end

  // There's name for a section of the code that needs EXCLUSIVE ACCESS to SHARED RESOURCE
  // This is Called CRITICAL SECTION
  //
  // There are 3 critical sections
  // 1. The goroutine, which is incrementing the data variables
  // 2. The if statement, which checks whether the value of data is 0
  // 3. The fmt.Printf statement, which retrieve the value of data for output
  //
  // To achive the atomic code, the process of accessing memory 
  // needs to SYNCHRONIZED by guarding CRITICAL SECTION
  // to avoid DATA RACE (not atomic)
  //
  // But the order of operations of this code is still nondeterministic
  // Either the goroutine will execute first, or both if and else section
  // And also it make the code slower, every time perform the operation to accessing memory,
  // the code need to pauses for a period of time
}
