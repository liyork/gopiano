package concurrent

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type maps struct {
	mapinmap sync.Map
}

func TestSyncMapInMap(t *testing.T) {
	m := maps{}

	// 错误
	actualMap, ok := m.mapinmap.LoadOrStore("a", sync.Map{})
	// 正确1
	//actualMap, ok := m.mapinmap.LoadOrStore("a", &sync.Map{})
	fmt.Println("actualMap, ok:", actualMap, ok)

	// 错误
	// ./syncmap_test.go:20:35: cannot call pointer method on actualMap.(sync.Map)
	// 方法签名func (m *Map) LoadOrStore
	//inmap, ok := actualMap.(sync.Map).LoadOrStore("b", "xxx1")
	// 正确2，单独声明一个变量，不报错。原因: 指针不能作为接收者, 需要需要定义变量来接送地址
	//i := actualMap.(sync.Map)
	//inmap, ok := i.LoadOrStore("b", "xxx1")
	// 正确1
	inmap, ok := actualMap.(*sync.Map).LoadOrStore("b", "xxx1")
	fmt.Println("inmap, ok:", inmap, ok)

}

type Counter struct {
	count int32
}

var counterMap sync.Map

func TestSyncMapConcurrent(t *testing.T) {
	go func() {
		for {
			addCounter("a")
			time.Sleep(300 * time.Millisecond)
		}
	}()
	//go func() {
	//	for {
	//		addCounter("a")
	//		time.Sleep(400 * time.Millisecond)
	//	}
	//}()
	go func() {
		for {
			subCounter("a")
			time.Sleep(300 * time.Millisecond)
		}
	}()
	time.Sleep(10 * time.Second)
}

func addCounter(key string) {
	tmpCounter, ok := counterMap.LoadOrStore(key, &Counter{})
	actualCounter := tmpCounter.(*Counter)
	atomic.AddInt32(&actualCounter.count, 1)
	if !ok {
		fmt.Println("first store key:", key, " count:", actualCounter.count)
	}
	fmt.Println("not first store key:", key, " count:", actualCounter.count)
}

func subCounter(key string) {
	value, _ := counterMap.Load(key)
	if value == nil {
		fmt.Println("key is nil ignore")
		return
	}
	count := value.(*Counter).count
	newCount := atomic.AddInt32(&count, -1)
	if newCount == 0 {
		counterMap.Delete(key)
		fmt.Println("key related is nil:", key)
	}
	fmt.Println("key related is not nil:", key)
}
