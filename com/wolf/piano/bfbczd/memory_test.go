package bfbczd

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

func TestMemory(t *testing.T) {
	memConsumed := func() uint64 {
		runtime.GC()
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}

	var wg sync.WaitGroup

	// 没有人写入，只有读取永远阻塞
	var c <-chan interface{}
	noop := func() { wg.Done(); <-c } // 1

	const numGoroutines = 1e4 // 2
	wg.Add(numGoroutines)

	before := memConsumed() // 3
	for i := numGoroutines; i > 0; i-- {
		go noop()
	}
	wg.Wait()

	after := memConsumed() // 4
	fmt.Printf("%.3fkb\n", float64(after-before)/numGoroutines/1000)

	//1.我们需要一个永不退出的goroutine，以便我们可以将它们中的一部分保存在内存中进行测量。
	//2.这里我们定义要创建的goroutines的数量。 我们将使用大数定律渐近地逼近一个goroutine的大小。
	//3.这里测量创建分区之前所消耗的内存量。
	//4.这里测量创建goroutines后消耗的内存量。
}
