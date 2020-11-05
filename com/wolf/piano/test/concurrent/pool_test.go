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
// 重用对象,sync.Pool的缓存对象数量是没有限制的（只受限于内存），清理时间比较随意

// 如何在多个goroutine之间使用同一个pool做到高效呢？官方的做法就是尽量减少竞争，因为sync.pool为每个P都分配了一个子池
// 当执行一个pool的get或者put操作的时候都会先把当前的goroutine固定到某个P的子池上面，然后再对该子池进行操作。每个子池里面有一个私有对象和共享列表对象，私有对象是只有对应的P能够访问，因为一个P同一时间只能执行一个goroutine，因此对私有对象存取操作是不需要加锁的。共享列表是和其他P分享的，因此操作共享列表是需要加锁的。
// 获取对象过程是：
// 1）固定到某个P，尝试从私有对象获取，如果私有对象非空则返回该对象，并把私有对象置空；
// 2）如果私有对象是空的时候，就去当前子池的共享列表获取（需要加锁）
// 3）如果当前子池的共享列表也是空的，那么就尝试去其他P的子池的共享列表偷取一个（需要加锁）
// 4）如果其他子池都是空的，最后就用用户指定的New函数产生一个新的对象返回
// 可以看到一次get操作最少0次加锁，最大N（N等于MAXPROCS）次加锁
// 可以看到一次get操作最少0次加锁，最大N（N等于MAXPROCS）次加锁
// 1）固定到某个P，如果私有对象为空则放到私有对象
// 2）否则加入到该P子池的共享列表中（需要加锁）
// 可以看到一次put操作最少0次加锁，最多1次加锁

// 使用原则：
// 1.不能对通过Get方法获取到的对象值有任何假设,到底哪一个值会被取出是完全不确定的。不能得知操作临时对象池的Goroutine在哪一时刻会与哪一个P相关联
// 所以，我们给予临时对象池的对象值生成函数所产生的值以及通过调用它的Put方法放入到池中的值都应该是无状态的或者状态一致的。
// 我们在取出并使用这些值的时候也不应该以其中的任何状态作为先决条件
// 2.任何对象值都有可能在任何时候被移除掉，并且根本不会通知该池的使用方。这种情况常常会发生在垃圾回收器即将开始回收内存垃圾的时候。
// 如果这时临时对象池中的某个对象值仅被该池引用，那么它还可能会在垃圾回收的时候被回收掉
// 因此，我们也就不能假设之前放入到临时对象池的某个对象值会一直待在池中，即使我们没有显式的把它从池中取出
// 适用场景
// 1.比如作为临时且状态无关的数据的暂存处
// 2.Pool用于存储那些被分配了但是没有被使用，而未来可能会使用的值，以减小垃圾回收的压力
// 不适用的场景
// 比如用来存放数据库连接的实例
// 思考，一个池，可以放和取，还是任意的，就不能有状态，任何一个都可以，实现重用。可是对于没有还能再通过new来一个新的。
func TestPoolBase(t *testing.T) {
	var count int32
	newFunc := func() interface{} {
		i := atomic.AddInt32(&count, 1)
		//fmt.Println("new value:", i)
		return i
	}
	// New:用来创建对象值,仅在池中无可用对象值的时候才会被调用
	pool := sync.Pool{New: newFunc}
	v1 := pool.Get() // 池中没有会用New函数,若有则返回
	fmt.Println("get value1:", v1)
	v1 = pool.Get()
	fmt.Println("get value2:", v1)
	pool.Put(v1)
	v1 = pool.Get()
	fmt.Println("get value3:", v1)
}

func TestPoolBaseGC(t *testing.T) {
	// 禁用GC，并保证在函数执行结束前恢复GC
	percent := debug.SetGCPercent(-1)
	defer debug.SetGCPercent(percent)

	var count int32
	newFunc := func() interface{} {
		return atomic.AddInt32(&count, 1)
	}
	pool := sync.Pool{New: newFunc}
	v1 := pool.Get() //1
	fmt.Printf("v1: %v\n", v1)
	// 放入
	pool.Put(newFunc()) //2
	pool.Put(newFunc()) //3
	pool.Put(newFunc()) //4
	// 在该方法返回它之前就一定会把它从池中删除掉，即删掉后返回。
	v2 := pool.Get() //随机
	fmt.Printf("v2: %v\n", v2)

	// 垃圾回收对临时对象池的影响
	debug.SetGCPercent(0)
	runtime.GC()

	v3 := pool.Get()
	fmt.Printf("v3: %v\n", v3)
	pool.New = nil // 之后get就为nil
	v4 := pool.Get()
	fmt.Printf("v4: %v\n", v4)
	v5 := pool.Get()
	fmt.Printf("v5: %v\n", v5)
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
