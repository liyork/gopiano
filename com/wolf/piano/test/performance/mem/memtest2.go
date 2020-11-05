package main

import (
	"log"
	"runtime"
	"time"
)

func traceMemStats() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	log.Printf(" ===> Alloc:%d(bytes) HeapIdle:%d(bytes) HeapReleased:%d(bytes)", ms.Alloc, ms.HeapIdle, ms.HeapReleased)
}

func test2() {
	container := make([]int, 8)

	log.Println(" ===> loop begin.")
	for i := 0; i < 32*1000*1000; i++ {
		container = append(container, i)
		if i == 16*1000*1000 {
			traceMemStats()
		}
	}

	log.Println(" ===> loop end.")
}

// 利用 runtime库里的ReadMemStats()方法,查看内存
// go run memtest2.go
func main() {
	log.Println(" ===> [Start].")

	traceMemStats()
	test2()
	traceMemStats()

	log.Println(" ===> [force gc].")
	runtime.GC() //强制调用gc回收

	log.Println(" ===> [Done].")
	traceMemStats()

	// 每10秒使用runtime.ReadMemStats获取内存使用情况。
	go func() {
		for {
			traceMemStats()
			time.Sleep(10 * time.Second)
		}
	}()

	time.Sleep(3600 * time.Second)
}
