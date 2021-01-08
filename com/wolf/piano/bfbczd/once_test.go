package bfbczd

import (
	"fmt"
	"sync"
	"testing"
)

func TestOnce(t *testing.T) {
	var count int

	increment := func() {
		count++
	}

	var once sync.Once

	var wg sync.WaitGroup

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			once.Do(increment)
		}()
	}

	wg.Wait()
	fmt.Printf("Count is %d\n", count)
}

func TestOnce1(t *testing.T) {
	var count int
	increment := func() { count++ }
	decrement := func() { count-- }

	var once sync.Once
	once.Do(increment)
	once.Do(decrement)

	fmt.Printf("Count: %d\n", count)
}

func TestOnce2(t *testing.T) {
	var onceA, onceB sync.Once
	var initB func()
	initA := func() {
		fmt.Println("in initA before")
		onceB.Do(initB)
		fmt.Println("in initA after")
	}
	initB = func() { // 1
		fmt.Println("in initB before")
		onceA.Do(initA)
		fmt.Println("in initB after")
	}
	onceA.Do(initA) // 2
	// 1这里的调用无法执行，直到2被返回。
}
