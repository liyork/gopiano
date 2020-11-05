package main

import (
	"log"
	"runtime"
	"time"
)

func test() {
	//slice会动态扩容，用它来做堆内存申请
	container := make([]int, 8)
	log.Println(" ===> loop begin.")
	for i := 0; i < 32*1000*1000; i++ {
		container = append(container, i)
	}
	log.Println(" ===> loop end.")
}

// go run memtest1.go
// top -p $(pidof snippet_mem)
// GODEBUG='gctrace=1' ./memtest1
func main() {
	log.Println("Start.")
	test()
	log.Println("force gc.")
	runtime.GC() //强制调用gc回收

	log.Println("Done.")

	time.Sleep(1 * time.Hour)
}
