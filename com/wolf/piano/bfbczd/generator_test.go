package bfbczd

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

// 创建了一个repeat生成器来生成无限数量的重复生成器，但是只取前10个。
// repeat生成器由take接收。虽然我们可以生成无线数量的流，但只会生成n+1个实例，其中n是我们传入take的数量。
func TestGenerator1(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	for num := range take(done, repeat(done, 1), 10) {
		fmt.Printf("v:%v ", num)
	}
}

// repeat的生成器,重复你传给它的值，直到你告诉它停止
var repeat = func(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	fmt.Println("repeat...")
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for { // 无限循环，每次都用values放入valueStream，直到done
			for _, v := range values {
				select {
				case <-done:
					return
				case valueStream <- v:
					fmt.Println("repeat ", v)
				}
			}
		}
	}()
	return valueStream
}

// 从其传入的valueStream中取出第1个元素
var take = func(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	fmt.Println("take...")
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
				fmt.Println("take ")
			}
		}
	}()
	return takeStream
}

// 创建一个重复调用函数的生成器repeatFn：
func TestGenerator2(t *testing.T) {
	repeatFn := func(done <-chan interface{}, fn func() interface{}) <-chan interface{} {

		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				select {
				case <-done:
					return
				case valueStream <- fn():
				}
			}
		}()
		return valueStream
	}

	done := make(chan interface{})
	defer close(done)

	rand := func() interface{} {
		return rand.Int()
	}

	for num := range take(done, repeatFn(done, rand), 10) {
		fmt.Println(num)
	}
}

// assert
func TestGenerator3(t *testing.T) {
	toString := getToStringFunc()

	done := make(chan interface{})
	defer close(done)

	var message string
	for token := range toString(done, take(done, repeat(done, "I", "am."), 5)) {
		message += token
	}

	fmt.Printf("message: %s...", message)
}

func getToStringFunc() func(done <-chan interface{}, valueStream <-chan interface{}) <-chan string {
	toString := func(done <-chan interface{}, valueStream <-chan interface{}) <-chan string {
		stringStream := make(chan string)
		go func() {
			defer close(stringStream)
			for v := range valueStream {
				select {
				case <-done:
					return
				case stringStream <- v.(string):
				}
			}
		}()
		return stringStream
	}
	return toString
}

func BenchmarkGeneric(b *testing.B) {
	done := make(chan interface{})
	defer close(done)

	b.ResetTimer()
	for range getToStringFunc()(done, take(done, repeat(done, "a"), b.N)) {
	}
}

func BenchmarkTyped(b *testing.B) {

	done := make(chan interface{})
	defer close(done)

	b.ResetTimer()
	for range take(done, repeat(done, "a"), b.N) {
	}
}

// 寻找素数例子--慢
func TestFindPrime1(t *testing.T) {
	rand := func() interface{} { return rand.Intn(50000000) }

	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	randIntStream := toInt(done, repeat(done, rand))
	fmt.Println("Primes:")
	for prime := range take(done, primeFinder(done, randIntStream), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v", time.Since(start))
}

func primeFinder(done chan interface{}, randIntStream interface{}) <-chan interface{} {
	return nil
}

func toInt(done chan interface{}, fn interface{}) interface{} {
	return nil
}

// 寻找素数例子--快，扇入扇出
func TestFindPrime2(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	rand := func() interface{} { return rand.Intn(50000000) }

	randIntStream := toInt(done, repeat(done, rand))

	numFinders := runtime.NumCPU()
	fmt.Printf("Spinning up %d prime finders.\n", numFinders)
	finders := make([]<-chan interface{}, numFinders)
	fmt.Println("Primes:")
	for i := 0; i < numFinders; i++ {
		finders[i] = primeFinder(done, randIntStream)
	}

	for prime := range take(done, fanIn(done, finders...), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v", time.Since(start))
}

var fanIn = func(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} { // 1

	var wg sync.WaitGroup // 2
	multiplexedStream := make(chan interface{})

	multiplex := func(c <-chan interface{}) { // 3
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- i:
			}
		}
	}

	// 从所有的通道中取数据
	wg.Add(len(channels)) // 4
	for _, c := range channels {
		go multiplex(c)
	}

	// 等待所有数据汇总完毕
	go func() { // 5
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

//1.使用done通道来关闭衍生的goroutine，并接收接口类型的通道切片来汇总数据。
//2.这里我们使用sync.WaitGroup以等待全部通道读取完成。
//3.我们在这里建立函数multiplex，它会读取传入的通道，并把该通道的值放入multiplexedStream。
//4.这里增加等待计数。
//5.这里我们建立一个goroutine等待汇总完毕。这样函数块可以快速return，不必等待wg.Wait()。这种用法不多见，但在这里很符合场景需求。
