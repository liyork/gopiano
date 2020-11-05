package base

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"
)

// atomic.Load*系列函数只能保证读取的不是正在写入的值（比如只被修改了一半的数据）
// atomic.Store* 只保证写入是原子操作(保证写入操作的完整性(如64位可能写入两次，但是也能保证原子))
// 但是两者结合并不保证原子性
// val := atomic.LoadInt32(&num) + atomic.StoreInt32(&num, val+1)

func TestAtomPointer(t *testing.T) {

	var target = 1
	var store unsafe.Pointer

	if atomic.LoadPointer(&store) != nil {
		fmt.Println("has value1:", *((*int)(store)))
		return
	}
	atomic.StorePointer(&store, unsafe.Pointer(&target))

	if atomic.LoadPointer(&store) != nil {
		fmt.Println("has value2:", *((*int)(store)))
		return
	}
}

func TestAtomIncrease(t *testing.T) {
	var a int32 = 1
	for {
		atomic.AddInt32(&a, 1)
		val := atomic.LoadInt32(&a)
		fmt.Println("val:", val)
	}
}

func TestForAdd(t *testing.T) {
	go func() {
		for {
			fmt.Println("this time", time.Now())
			time.Sleep(2 * time.Second)
		}
	}()
	i := 0
	var count int32 = 0
	go func() {
		for {
			i++
			atomic.AddInt32(&count, 1)
		}
	}()
	go func() {
		for {
			atomic.StoreInt32(&count, 0)
		}
	}()
	go func() {
		for {
			val := atomic.LoadInt32(&count)
			fmt.Println("load value:", val)
		}
	}()

	//time.Sleep(11111 * time.Second)
}

func TestAtomicValue(t *testing.T) {
	value := atomic.Value{}
	value.Store(1)
	x := value.Load()
	fmt.Println("x", x)

}
