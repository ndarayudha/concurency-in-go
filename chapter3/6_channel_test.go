package chapter3

import (
	"fmt"
	"testing"
)

// Declaring Channel
//
// var dataStream chan interface {}
// dataStream = make(chan interface{})
//
// unidirectional channel (received only)
// var dataStream <- chan interface {}
// dataStream := make(<-chan interface{})
//
// unidirectional channel (send only)
// var dataStream chan <- interface{}
// dataStream := make(chan <- interface{})
//
// Go will implicitly convert bidirectional channels to unidirectional channel when needed
// var receiveChan <-chan interface{}
// var sendChan chan<- interface{}
// dataStream := make(chan interface{})
//
// receiveChan = dataStream
// sendChan = dataStream

func TestChannel(t *testing.T) {
	stringStream := make(chan string)
	go func() {
		stringStream <- "Hello channels!"
	}()
	fmt.Println(<-stringStream)
}
