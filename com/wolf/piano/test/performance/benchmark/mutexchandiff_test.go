package benchmark

import (
	"fmt"
	"sync"
	"testing"
)

//func BenchmarkTest(b *testing.B) {
//	var m = make(map[int]int, 10000)
//	var mu = &sync.Mutex{}
//	var c = make(chan int, 200)
//	var w = &sync.WaitGroup{}
//
//	w.Add(b.N)
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		c <- 1
//		go func(index int) {
//			mu.Lock()
//			m[index] = 0
//			mu.Unlock()
//			w.Done()
//			<-c
//		}(i)
//	}
//	w.Wait()
//}
//
//func BenchmarkTestBlock(b *testing.B) {
//	var m = make(map[int]int, 10000)
//	var c = make(chan int, 200)
//	var dataChan = make(chan int, 200)
//	var w = &sync.WaitGroup{}
//
//	w.Add(b.N + 1)
//	go func() {
//		for i := 0; i < b.N; i++ {
//			m[<-dataChan] = 0
//		}
//		w.Done()
//	}()
//
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		c <- 1
//		go func(index int) {
//			dataChan <- index
//			w.Done()
//			<-c
//		}(i)
//	}
//	w.Wait()
//}

// go test -v  -bench=MapGo -count=1 -benchtime=2s  mutexchandiff_test.go

// 每个go协成进行map操作次数
var iteTimes = 10000

// 1~5000，写占比
var writePer = 10000

// 0~1
var mod = 1

//func TestQexsdf(t *testing.T) {
//	for i := 0; i < iteTimes; i++ {
//		if i%writePer == mod {
//			fmt.Println("i%writePer=1, i:", i)
//		} else {
//			fmt.Println("i%writePer!=1, i:", i)
//		}
//	}
//}

// sync.Map的性能高体现在读操作远多于写操作的时候
//建议使用sync.Map时一定要考虑读定比例。当写操作只占总操作的<=1/10的时候，使用sync.Map性能会明显高很多。10次中有9次是读
func BenchmarkSyncMapGo(b *testing.B) {
	var mp sync.Map
	var wg sync.WaitGroup
	fmt.Println("i:", b.N)
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < iteTimes; j++ {
				if i%writePer == mod {
					mp.Store(0, 0)
				} else {
					_, _ = mp.Load(0)
				}
			}

		}(i)
	}
	wg.Wait()
}

func BenchmarkMapGo(b *testing.B) {
	var mp = make(map[int]int)
	var wg sync.WaitGroup
	var lock sync.Mutex
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			for j := 0; j < iteTimes; j++ {
				lock.Lock()
				if i%writePer == mod {
					mp[0] = 0
				} else {
					i = mp[0]
				}
				lock.Unlock()
			}
		}(i)
	}
	wg.Wait()
}
