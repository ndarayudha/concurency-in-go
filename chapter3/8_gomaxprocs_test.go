package chapter3

import (
	"fmt"
	"runtime"
	"testing"
)

func TestGMP(t *testing.T) {
	// this function controls the number of OS threads that will host so-called “work queues.”
	fmt.Println(runtime.GOMAXPROCS(runtime.NumCPU())) // print default logical core of machine
}
