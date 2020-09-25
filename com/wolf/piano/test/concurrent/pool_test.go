package concurrent

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 容器是自动伸缩的、高效的，同时也是并发安全的
// 特性1:容器把由其中的对象值产生的存储压力进行分摊。会专门为每一个与操作它的Goroutine相关联的P都生成一个本地池。
// 在Get方法被调用的时候，一般会先尝试从与本地P对应的那个本地池中获取一个对象值。如果获取失败，就会试图从其他P的本地池中
// 偷一个对象值并直接返回给调用方。如果依然未果，使用New生成函数了，产生的对象值永远不会被放置到池中,直接返回给调用方
// Put方法会把它的参数值存放到与当前P对应的那个本地池中。每个P的本地池中的绝大多数对象值都是被同一个容器中的所有本地池所共享的，随时可能会被偷走
// 特性2:对垃圾回收友好。垃圾回收的执行一般会使容器中的所有对象值被全部移除。
// 使用原则：
// 1.不能对通过Get方法获取到的对象值有任何假设,到底哪一个值会被取出是完全不确定的。不能得知操作临时对象池的Goroutine在哪一时刻会与哪一个P相关联
// 所以，我们给予临时对象池的对象值生成函数所产生的值以及通过调用它的Put方法放入到池中的值都应该是无状态的或者状态一致的。
// 我们在取出并使用这些值的时候也不应该以其中的任何状态作为先决条件
// 2.任何对象值都有可能在任何时候被移除掉，并且根本不会通知该池的使用方。这种情况常常会发生在垃圾回收器即将开始回收内存垃圾的时候。
// 如果这时临时对象池中的某个对象值仅被该池引用，那么它还可能会在垃圾回收的时候被回收掉
// 因此，我们也就不能假设之前放入到临时对象池的某个对象值会一直待在池中，即使我们没有显式的把它从池中取出
// 适用场景
// 1.比如作为临时且状态无关的数据的暂存处
// 不适用的场景
// 比如用来存放数据库连接的实例
// 思考，一个池，可以放和取，还是任意的，就不能有状态，任何一个都可以，实现重用。可是对于没有还能再通过new来一个新的。
func TestPoolBase(t *testing.T) {
	// 禁用GC，并保证在函数执行结束前恢复GC
	percent := debug.SetGCPercent(-1)
	defer debug.SetGCPercent(percent)

	var count int32
	newFunc := func() interface{} {
		//fmt.Println("new func...", )
		return atomic.AddInt32(&count, 1)
	}
	// New:用来创建对象值,仅在池中无可用对象值的时候才会被调用
	pool := sync.Pool{New: newFunc}
	v1 := pool.Get() // 池中没有会用New函数
	fmt.Printf("v1: %v\n", v1)
	// 放入
	pool.Put(newFunc())
	pool.Put(newFunc())
	pool.Put(newFunc())
	// 获取到的值是任意的,在该方法返回它之前就一定会把它从池中删除掉。
	v2 := pool.Get()
	fmt.Printf("v2: %v\n", v2)
	// 垃圾回收对临时对象池的影响
	debug.SetGCPercent(0)
	runtime.GC()
	v3 := pool.Get()
	fmt.Printf("v3: %v\n", v3)
	pool.New = nil // 之后get就为nil
	v4 := pool.Get()
	fmt.Printf("v4: %v\n", v4)
}

// 	输出是：
//v1: 1
//v2: 2
//v3: 5
//v4: <nil>
// 为啥我本地是
//v1: 1
//v2: 2
//v3: 3
//v4: 4

// 用于等待都执行done让WaitGroup值为0时唤醒Wait
func TestWaitGroupBase(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		fmt.Println("sub routine is running")
		time.Sleep(2 * time.Second)
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("main routine is running")
}

// 值传递，有问题，等和减的不是一个对象，需要将指针传递过去
func TestWaitGroupPass(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go recv(&wg)
	fmt.Printf("main wg:%p\n", &wg)
	wg.Wait()
	fmt.Println("finish")
}

func recv(wg *sync.WaitGroup) {
	fmt.Printf("recv wg:%p\n", &wg)
	defer wg.Done()

	time.Sleep(2 * time.Second)
	fmt.Println("rece is runnning")
}
