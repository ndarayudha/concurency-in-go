package chapter1

import (
	"fmt"
	"testing"
)

func TestRaceCondition(t *testing.T) {
  var data int

  go func() {
    data++
  }()

  if data == 0 {
    fmt.Printf("the value is %v", data)
  }
}

