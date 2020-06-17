package base

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var NodeMu = map[string]sync.RWMutex{}

// mutex是非重入锁
func TestMutexBase(t *testing.T) {
	var mutex sync.Mutex
	mutex.Lock()
	// 如果锁不是lock状态，会panic
	mutex.Unlock()
}

func TestMutexIns(t *testing.T) {
	var count int
	var lock sync.Mutex
	var arthmatic sync.WaitGroup
	arthmatic.Add(1)

	Increment := func() {
		lock.Lock()
		defer lock.Unlock()
		count++
		//fmt.Printf("Incrementing: %d\n", count)
	}

	Decrement := func() {
		lock.Lock()
		defer lock.Unlock()
		count--
		//fmt.Printf("Decrementing: %d\n", count)
	}

	for i := 0; i < 5; i++ {
		arthmatic.Add(1)
		go func() {
			defer arthmatic.Done()
			Increment()
		}()
	}

	for i := 0; i < 5; i++ {
		arthmatic.Add(1)
		go func() {
			defer arthmatic.Done()
			Decrement()
		}()
	}

	arthmatic.Wait()
	fmt.Println("Arthmatic completed! ", count)
}

func TestMapLock(t *testing.T) {
	go func() {
		for i := 1; i < 100; i++ {
			mutex := NodeMu["xxx"]
			mutex.Lock()
			time.Sleep(5 * time.Second)
			fmt.Println("xxxx")
			mutex.Unlock()
		}
	}()

	go func() {
		for i := 1; i < 100; i++ {
			mutex := NodeMu["xxx"]
			mutex.Lock()
			time.Sleep(5 * time.Second)
			fmt.Println("xxxx222222")
			mutex.Unlock()
		}
	}()

	time.Sleep(1115 * time.Second)
}
func TestMapErr(t *testing.T) {
	Map := make(map[int]int)

	for i := 0; i < 10; i++ {
		go writeMap(Map, i, i)
		go readMap(Map, i)
	}
}

func readMap(Map map[int]int, key int) int {
	return Map[key]
}

func writeMap(Map map[int]int, key int, value int) {
	Map[key] = value
}

type SafeMap struct {
	sync.RWMutex
	Map map[int]int
}

func TestSafeMap(t *testing.T) {
	safeMap := newSafeMap(10)

	for i := 0; i < 100000; i++ {
		go safeMap.writeMap(i, i)
		go safeMap.readMap(i)
	}
}

func newSafeMap(size int) *SafeMap {
	sm := new(SafeMap)
	sm.Map = make(map[int]int)
	return sm

}

func (sm *SafeMap) readMap(key int) int {
	sm.RLock()
	value := sm.Map[key]
	sm.RUnlock()
	return value
}

func (sm *SafeMap) writeMap(key int, value int) {
	sm.Lock()
	sm.Map[key] = value
	sm.Unlock()
}

// 读多写少场景，可防止map的并发场景异常
func TestSyncMapBase(t *testing.T) {
	// 无须初始化，直接声明即可
	var sm sync.Map

	if _, ok := sm.Load(1); !ok {
		fmt.Println("not load 1")
	}

	sm.Store(1, "a")
	if v, ok := sm.Load(1); ok {
		fmt.Println("load 1, value:", v)
	}

	// key存在则返回原值+true
	// key不存在则存储key:value，返回新值+false
	if vv, ok := sm.LoadOrStore(1, "c"); ok {
		fmt.Println("LoadOrStore 1 c exist, ", vv)
	}
	if vv, ok := sm.LoadOrStore(2, "c"); !ok {
		fmt.Println("LoadOrStore 2 c not exist, ", vv)
	}

	// 两次store覆盖
	sm.Store(3, "d")
	sm.Store(3, "e")

	a := 0
	sm.Range(func(k, v interface{}) bool {
		//if a == 1 {
		//	return false //false 则退出
		//}
		fmt.Println(k, ":", v)
		a++
		return true
	})

}

// 创建sync.Map后不要再拷贝使用
