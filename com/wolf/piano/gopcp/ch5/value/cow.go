package value

import "sync/atomic"

// 并发安全的整数数组接口
type ConcurrentArray interface {
	// 用于设置指定索引上的元素值
	Set(index uint32, elem int) (err error)
	// 获取指定索引上的元素值
	Get(index uint32) (elem int, err error)
	// 获取数组的长度
	Len() uint32
}

// 实现
type concurrentArray struct {
	length uint32
	val    atomic.Value
}

// 创建一个ConcurrentArray类型值
func NewConcurrentArray(length uint32) ConcurrentArray {
	array := concurrentArray{}
	array.length = length
	array.val.Store(make([]int, array.length))
	return &array
}

// copy-on-write
func (array *concurrentArray) Set(index uint32, elem int) (err error) {
	if err = array.checkIndex(index); err != nil {
		return
	}
	if err = array.checkValue(); err != nil {
		return
	}
	newArray := make([]int, array.length)
	// 拷贝一个值
	copy(newArray, array.val.Load().([]int))
	// 设定拷贝值
	newArray[index] = elem
	// 原子写回
	array.val.Store(newArray)
	return
}

func (array *concurrentArray) Get(index uint32) (elem int, err error) {
	if err = array.checkIndex(index); err != nil {
		return
	}
	if err = array.checkValue(); err != nil {
		return
	}
	elem = array.val.Load().([]int)[index]
	return
}

func (array *concurrentArray) Len() uint32 {
	return array.length
}

func (array *concurrentArray) checkIndex(u uint32) error {
	return nil
}

func (array *concurrentArray) checkValue() error {
	return nil
}
