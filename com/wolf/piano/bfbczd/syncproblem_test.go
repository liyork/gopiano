package bfbczd

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 保护value的数据竞争
func TestProtectMemory(t *testing.T) {
	var memoryAccess sync.Mutex //1
	var value int
	go func() {
		memoryAccess.Lock() //2 独占，上锁
		value++
		memoryAccess.Unlock() //3 解锁
	}()

	memoryAccess.Lock() //4 独占，解锁
	if value == 0 {
		fmt.Printf("the value is %v.\n", value)
	} else {
		fmt.Printf("the value is %v.\n", value)
	}
	memoryAccess.Unlock() //5
}

// 死锁
func TestDeadlock(t *testing.T) {
	type value struct {
		mu    sync.Mutex
		value int
	}

	var wg sync.WaitGroup
	printSum := func(v1, v2 *value) {
		defer wg.Done()

		v1.mu.Lock()         //1
		defer v1.mu.Unlock() //2

		time.Sleep(2 * time.Second) //3
		v2.mu.Lock()
		defer v2.mu.Unlock()

		fmt.Printf("sum=%v\n", v1.value+v2.value)
	}

	var a, b value
	wg.Add(2)
	go printSum(&a, &b)
	go printSum(&b, &a)
	wg.Wait()

}

func TestAlivelock(t *testing.T) {
	// 公用一个锁
	cadence := sync.NewCond(&sync.Mutex{})
	go func() {
		// 每1毫秒，广播一下，通知大家唤醒
		for range time.Tick(1 * time.Millisecond) {
			cadence.Broadcast()
		}
	}()

	// 上锁等待
	takeStep := func() {
		cadence.L.Lock()
		cadence.Wait() // 睡眠并释放锁，当被唤醒后还会再尝试获取锁
		cadence.L.Unlock()
	}

	// +1，等待被通知，若还是1则true返回，否则等待被通知，-1，返回false
	tryDir := func(dirName string, direction *int32, out *bytes.Buffer) bool { //1
		fmt.Fprintf(out, " %v", dirName)
		atomic.AddInt32(direction, 1) //2
		takeStep()                    //3
		//fmt.Println(dirName, " ", *direction)// 值永远是2，两个协程一起执行add为2，然后一个协程获锁等待，唤醒后是2，此时不会产生第一个协程再获取锁然后等待然后再-1?应该是goruntime保证等待的有先后
		if atomic.LoadInt32(direction) == 1 {
			fmt.Fprint(out, ". Success!")
			return true
		}
		takeStep()
		atomic.AddInt32(direction, -1) //4
		return false
	}

	var left, right int32
	tryLeft := func(out *bytes.Buffer) bool { return tryDir("left", &left, out) }
	tryRight := func(out *bytes.Buffer) bool { return tryDir("right", &right, out) }
	//1.tryDir 允许一个人尝试向某个方向移动并返回，无论他们是否成功。 每个方向都表示为试图朝这个方向移动的次数。
	//2.首先，我们通过将该方向递增1来朝着某个方向移动
	//3.每个人必须以相同的速度或节奏移动。 takeStep模拟所有动作之间的恒定节奏。
	//4.在这里，这个人意识到他们不能在这个方向上放弃。 我们通过将该方向递减1来表示这一点。

	walk := func(walking *sync.WaitGroup, name string) {
		var out bytes.Buffer
		// 退出时打印
		defer func() { fmt.Println(out.String()) }()
		defer walking.Done()

		fmt.Fprintf(&out, "%v is trying to scoot:", name)
		for i := 0; i < 5; i++ { //1
			if tryLeft(&out) || tryRight(&out) { //2
				return
			}
		}
		fmt.Fprintf(&out, "\n%v tosses her hands up in exasperation!", name)
	}

	var peopleInHallway sync.WaitGroup //3
	peopleInHallway.Add(2)
	// 两个协程
	go walk(&peopleInHallway, "Alice")
	go walk(&peopleInHallway, "Barbara")

	peopleInHallway.Wait()

	//1.我对尝试次数进行了人为限制，以便该程序结束。 在一个有活锁的程序中，可能没有这种限制，这就是为什么它是一个现实工作中的问题。
	//2.首先，这个人会试图向左走，如果失败了，会尝试向右走。
	//3.这个变量为程序提供了等待，直到两个人都能够相互通过或放弃。
}

func TestStarvation(t *testing.T) {
	var wg sync.WaitGroup
	var sharedLock sync.Mutex
	const runtime = 1 * time.Second

	greedyWorker := func() {
		defer wg.Done()

		var count int
		// 小于1s则执行
		for begin := time.Now(); time.Since(begin) <= runtime; {
			sharedLock.Lock()
			time.Sleep(3 * time.Nanosecond)
			sharedLock.Unlock()
			count++
		}

		fmt.Printf("Greedy worker was able to execute %v work loops\n", count)
	}

	politeWorker := func() {
		defer wg.Done()

		var count int
		for begin := time.Now(); time.Since(begin) <= runtime; {
			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			count++
		}
		fmt.Printf("Polite worker was able to execute %v work loops.\n", count)
	}

	wg.Add(2)
	go greedyWorker()
	go politeWorker()

	wg.Wait()

}
