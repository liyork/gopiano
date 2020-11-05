package base

import (
	"fmt"
	"reflect"
	"testing"
)

// 传递的是值的拷贝，不会影响原值
func passValue1(a int, b int) {
	a = b
}

func TestPassValue1(t *testing.T) {
	x := 5
	y := 10
	passValue1(x, y)
	fmt.Println(x)
}

// 传递地址的拷贝
func passPoint(a *int, b *int) {
	*a = *b
}

func TestPassPoint(t *testing.T) {
	x := 5
	y := 10
	passPoint(&x, &y)
	fmt.Println(x)
}

// go的数组传递也是值传递，进行拷贝
func passArray(a [1]int) {
	a[0] = 1111
}

func passArrayPointer(a *[1]int) {
	(*a)[0] = 1111
}

func TestPassArray(t *testing.T) {
	arr := [1]int{1}
	passArray(arr)
	fmt.Println(arr)
}

// 传递数组的引用
func TestPassArrayPointer(t *testing.T) {
	arr := [1]int{1}
	passArrayPointer(&arr)
	fmt.Println(arr)
}

// slice作为参数传递的是指针，因此传递速度会快很多
func TestPassSlice(t *testing.T) {
	arr := []int{1}
	passSlice(arr)
	fmt.Printf("in TestPassSlice, p:%p, v:%v\n", &arr, arr)
}

// slice传递的是切片拷贝，但是底层用的一个数组
func passSlice(a []int) {
	// 修改底层array
	a[0] = 1111
	fmt.Printf("in passSlice, p:%p, v:%v\n", &a, a)
}

// 值类型和引用类型
// 引用不绑定特定对象（c++中引用是要绑定特定对象），例如有两个同类型引用a和b，它们可以引用各自的对象A和B；但如果a和b的引用都指向A，那么通过b修改对象内容可以反应到a引用之中。
// golang从本质上说，一切皆是值类型，并没有内建一个类似java或c#等语言中的reference类型。
// golang可以使用指针，可以在行为上达到类似java或c#等语言中reference的效果。
// golang中从行为表现上看，数组属于值类型，数组切片、字典、通道和接口属于引用类型。

func TestPointUpdate(t *testing.T) {
	type Rect struct {
		width, height float64
	}

	var a Rect = Rect{100, 200}
	// 指针修改
	var b *Rect = &a

	b.width = 300
	b.height = 400

	fmt.Println(a.width, a.height) // 300 400
}

// 数组属于值类型
func TestArrayUpdate(t *testing.T) {
	var a [3]int = [3]int{1, 2, 3}
	var b [3]int = a
	// 数组是值拷贝，不影响a
	b[1]++

	fmt.Println(a) // [1 2 3]
	fmt.Println(b) // [1 3 3]
}

func TestSliceRefUpdate(t *testing.T) {
	var a []int = []int{1, 2, 3}
	var b []int = a

	// 底层公用一个数组，修改一个内容
	b[1]++

	fmt.Println(a) // [1 3 3]
	fmt.Println(b) // [1 3 3]
}

// map为引用传递
func TestMapRefUpdate(t *testing.T) {
	var a, b map[string]string
	a = make(map[string]string)
	a["1"] = "haha"
	a["2"] = "hehe"

	b = a
	b["2"] = "shit"

	fmt.Println(a) // map[1:haha 2:shit]
	fmt.Println(b) // map[1:haha 2:shit]
}

func ModifyRect(r Rect) {
	r.width = 1000
	r.height = 1000
	return
}

// 结构体属于值类型，跨函数传递某个对象将只能得到其副本，倘若要在另一个函数中修改对象的内容，那么结果只是修改了副本内容，原对象的内容将没有改变
func TestStructUpdate(t *testing.T) {
	rect := Rect{100, 200}
	ModifyRect(rect)
	fmt.Println(rect.width, rect.height) // 100, 200
}

func ModifyRect2(r *Rect) {
	r.width = 1000
	r.height = 1000
	return
}

// 只有通过传递指针，才能跨函数修改对象内容。
// 由于访问对象成员无论是实例还是指针，使用的都是点操作符，这就带来了类型的隐蔽性，因此在创建对象和传递实例时，尽量使用对象指针
func TestStructUpdateUsingPoint(t *testing.T) {
	rect := &Rect{100, 200}
	ModifyRect2(rect)
	fmt.Println(rect.width, rect.height) // 1000, 1000
}

// 可以给任意类型（包括内置类型，但不包括指针类型）添加相应的成员方法。
// 没有隐藏的this或self指针，即方法施加的目标对象将被显式传递，没有被隐藏起来。
// 成员对象定义时候指定作用的目标对象是对象实例还是对象指针。
// 无论是对象实例还是对象指针都可以调用成员函数，不管成员函数作用对象类型，但是不影响最后的结果，即结果只由作用的目标类型决定（是对象实例还是对象指针）
func (r *Rect) ModifyRect() {
	r.width = 1000
	r.height = 1000
	fmt.Println(r.width, r.height) // 1000, 1000
}
func TestMemberMethod(t *testing.T) {
	rect1 := Rect{100, 200}
	fmt.Println(reflect.ValueOf(rect1).Type())
	// 定义是指针的方法,使用对象实例调用依然有效
	rect1.ModifyRect()
	fmt.Println(rect1.width, rect1.height) // 1000, 1000

	rect2 := &Rect{100, 200}
	fmt.Println(reflect.ValueOf(rect2).Type())
	rect2.ModifyRect()                     // 使用对象指针调用依然有效
	fmt.Println(rect2.width, rect2.height) // 1000, 1000
}
