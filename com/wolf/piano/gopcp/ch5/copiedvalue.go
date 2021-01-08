package main

import (
	"fmt"
	"sync/atomic"
)

func main() {
	var countValue atomic.Value
	countValue.Store([]int{1, 3, 5, 7})
	anotherStore(countValue)
	fmt.Printf("The count value: %+v \n", countValue.Load())
}

func anotherStore(countValue atomic.Value) {
	countValue.Store([]int{2, 4, 6, 8})
}
