package base

import (
	"fmt"
	"testing"
)

type xy struct {
	x, y int
}
type pool struct {
	member xy
}

// for中每次将数组内容赋值给p，然后out中的值是对p取指针，进而添加到out中的都是一个p指针
func TestForError(t *testing.T) {
	pools := []pool{
		{
			member: xy{x: 1, y: 1},
		},
		{
			member: xy{x: 2, y: 2},
		},
		{
			member: xy{x: 3, y: 3},
		},
	}

	var out []*xy
	//新建一个变量p, 每次将该数组的元素赋值给这个ｐ,p的地址恒定的
	for _, p := range pools {
		fmt.Printf("p=%p\t", &(p.member))
		// 指针
		out = append(out, &(p.member))
	}
	// 上面后，out中所有内容(指针)都是p最后一次被赋值的值
	fmt.Println("")

	for _, ele := range out {
		ele.x = ele.x + 1
		ele.y = ele.y + 1
	}

	for _, ele := range out {
		// 三个元素是同一个指针
		fmt.Printf("p=%p, x=%d, y=%d\n", ele, ele.x, ele.y)
	}
}

func TestForCorrect1(t *testing.T) {
	pools := []pool{
		{
			member: xy{x: 1, y: 1},
		},
		{
			member: xy{x: 2, y: 2},
		},
		{
			member: xy{x: 3, y: 3},
		},
	}

	var out []*xy
	//新建一个变量p, 每次将该数组的元素赋值给这个ｐ,p的地址恒定的
	for _, p := range pools {
		p := p //解决方案
		fmt.Printf("p=%p\t", &(p.member))
		// 指针
		out = append(out, &(p.member))
	}
	// 上面后，out中所有内容(指针)都是p最后一次被赋值的值
	fmt.Println("")

	for _, ele := range out {
		ele.x = ele.x + 1
		ele.y = ele.y + 1
	}

	for _, ele := range out {
		// 三个元素是同一个指针
		fmt.Printf("p=%p, x=%d, y=%d\n", ele, ele.x, ele.y)
	}
}

type pool2 struct {
	member *xy
}

// 直接pool2中就是指针，每次赋值给p都是一个数组中元素(是指针)，p只用来做临时变量传递给out中
// 看来不要对for生成的临时变量p取指针，就可以避免问题
func TestForCorrect2(t *testing.T) {
	pools := []pool2{
		{
			member: &xy{x: 1, y: 1},
		},
		{
			member: &xy{x: 2, y: 2},
		},
		{
			member: &xy{x: 3, y: 3},
		},
	}

	var out []*xy
	//新建一个变量p, 每次将该数组的元素赋值给这个ｐ,p的地址恒定的
	for _, p := range pools {
		fmt.Printf("p=%p\t", p.member)
		// 指针
		out = append(out, p.member)
	}
	// 上面后，out中所有内容(指针)都是p最后一次被赋值的值
	fmt.Println("")

	for _, ele := range out {
		ele.x = ele.x + 1
		ele.y = ele.y + 1
	}

	for _, ele := range out {
		// 三个元素是同一个指针
		fmt.Printf("p=%p, x=%d, y=%d\n", ele, ele.x, ele.y)
	}
}

func TestForSlicePerformance(t *testing.T) {
	var slice []int
	for index, value := range slice {
		fmt.Println(index, value)
	}

	// 遍历过程中每次迭代会对index和value进行赋值，如果数据量大或者value类型为string时，对value的赋值操作可能是多余的，
	// 可以忽略value值，使用slice[index]引用value值。
	for index, _ := range slice {
		fmt.Println(slice[index])
	}
}

// 根据key值获取value值，虽然看似减少了一次赋值，但通过key值查找value值的性能消耗可能高于赋值消耗。
// 能否优化取决于map所存储数据结构特征、结合实际情况进行。看查找的快则优化成功，否则就不合适了
func TestForMapPerformance(t *testing.T) {
	var myMap map[int]string
	for key, _ := range myMap {
		_, _ = key, myMap[key]
	}
}

// 正常结束。循环内改变切片的数据和长度，不影响循环次数，循环次数在循环开始前就已经确定了。
func TestForResult(t *testing.T) {
	v := []int{1, 2, 3}
	for i := range v {
		v = append(v, i)
	}
}
