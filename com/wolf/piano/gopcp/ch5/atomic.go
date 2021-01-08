package main

import (
	"fmt"
	"sync/atomic"
)

func main() {
	var ui32 uint32 = 10
	// -3
	atomic.AddUint32(&ui32, ^uint32(int32(3-1)))
	fmt.Println("ui32:", ui32)
	fmt.Printf("ui32:%32b\n", ui32)
}
