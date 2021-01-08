package bfbczd

import (
	"fmt"
	"testing"
)

// 使得for range也可以感知到done
func TestForCanReturn(t *testing.T) {
	myChan := make(chan interface{})

	done := make(chan interface{})
	defer close(done)

	for val := range orDone(done, myChan) {
		// Do something with val
		fmt.Println("val:", val)
	}
}

var orDone = func(done, c <-chan interface{}) <-chan interface{} {

	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()

	return valStream
}
