package bfbczd

import (
	"fmt"
	"math"
	"os"
	"sync"
	"testing"
	"text/tabwriter"
	"time"
)

func TestMutex(t *testing.T) {
	var count int
	var lock sync.Mutex

	increment := func() {
		lock.Lock()         // 1
		defer lock.Unlock() // 2
		count++
		fmt.Printf("Incrementing: %d\n", count)
	}

	decrement := func() {
		lock.Lock()         // 1
		defer lock.Unlock() // 2
		count--
		fmt.Printf("Decrementing: %d\n", count)
	}

	// Increment
	var arithmetic sync.WaitGroup
	for i := 0; i <= 5; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			increment()
		}()

	}

	// Decrement
	for i := 0; i <= 5; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			decrement()
		}()
	}

	arithmetic.Wait()
	fmt.Println("Arithmetic complete.")
	//1.在这里，我们要求独占使用关键部分 - 在这种情况下，count变量由互斥锁保护。
	//2.这里表明我们已经完成了对共享部分的锁定。

}

func TestRWMutex(t *testing.T) {
	// 生产用写锁
	producer := func(wg *sync.WaitGroup, l sync.Locker) { //1
		defer wg.Done()
		for i := 5; i > 0; i-- {
			l.Lock()
			l.Unlock()
			time.Sleep(1) //2
		}
	}

	// 观察者用读锁
	observer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		l.Lock()
		defer l.Unlock()
	}

	test := func(count int, mutex, rwMutex sync.Locker) time.Duration {
		var wg sync.WaitGroup
		wg.Add(count + 1) // 1生产，count消费

		beginTestTime := time.Now()
		go producer(&wg, mutex)
		for i := count; i > 0; i-- {
			go observer(&wg, rwMutex)
		}

		wg.Wait()
		return time.Since(beginTestTime)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	defer tw.Flush()

	var m sync.RWMutex
	fmt.Fprintf(tw, "Readers\tRWMutext\tMutex\n")
	for i := 0; i < 20; i++ {
		count := int(math.Pow(2, float64(i)))
		fmt.Fprintf( // 							写，读锁								写，写锁
			tw, "%d\t%v\t%v\n", count, test(count, &m, m.RLocker()), test(count, &m, &m),
		)
	}
	//1.producer函数的第二个参数是类型sync.Locker。 该接口有两种方法，锁定和解锁，互斥和RWMutex类型都适用。
	//2.在这里，我们让producer休眠一秒钟，使其不那么活跃。
	//你可以通过这个例子看到，RWMutext在大量级上相对于Mutex是有性能优势的，不过这同样取决于你在锁住的部分做了什么。通常建议在逻辑上合理的情况下使用RWMutex而不是Mutex。
}
