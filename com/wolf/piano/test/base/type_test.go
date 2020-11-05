package base

import (
	"fmt"
	"testing"
	"unsafe"
)

// golang是一种强类型语言，虽然在代码中经常看到这种写法，i:=10这其实这是编译器自动做了类型推断在编译期间。
// 编译器会对数据进行类型检查。不同类型的数据不能赋值,不能在函数中传参。强类型语言有一些优势，很多的错误会在编译期间被检查出来
// golang中的类型转换分强制类型转换和类型断言
func TestTypeConvert1(t *testing.T) {
	var a1 float32 = 5.6
	var b1 int = 10
	//fmt.Println (a * b)// 类型不匹配,需要强制类型转换
	fmt.Println(a1 * float32(b1))

	// 普通变量类型int,float,string 都可以使用 type (a)这种形式来进行强制类型转换,比如
	var a2 int32 = 10
	var b2 int64 = int64(a2)
	var c2 float32 = 12.3
	var d2 float64 = float64(c2)
	fmt.Println(b2, d2)

	// 指针也是有类型的,
	var a int = 10
	var p *int = &a
	// c= (*int64)(p)// cannot convert p (type *int) to type *int64
	// 指针的强制类型转换需要用到unsafe包中的函数实现
	var c *int64 = (*int64)(unsafe.Pointer(p))
	fmt.Println(*c)
}

// 类型转换用于将一种数据类型的变量转换为另外一种类型的变量。
// type_name(expression)
func TestTypeConvert2(t *testing.T) {

	var sum int = 17
	var count int = 5
	var mean float32

	mean = float32(sum) / float32(count) // 将整型转化为浮点型
	fmt.Printf("mean 的值为: %f\n", mean)
}

// go的语言中提供了断言的功能。所有程序都实现了interface{}的接口，这意味着，所有的类型如string,int,int64甚至是自定义的struct类型都就
// 此拥有了interface{}的接口，这种做法和java中的Object类型比较类似。
// 在一个数据通过func funcName(interface{})的方式传进来的时候，也就意味着这个参数被自动的转为interface{}的类型。
func TestTypeAssert1(t *testing.T) {
	var a interface{} = 10
	switch a.(type) {
	case int:
		fmt.Println("int")
	case float32:
		fmt.Println("string")
	}
}

func TestTypeAssert2(t *testing.T) {
	var a interface{} = 10
	// 第一个是对应类型的值,第二个是bool类型的,类型判断是否正确
	q, ok := a.(int)
	if ok {
		fmt.Println("int", q)
	}
	t2, ok := a.(float32)
	if !ok {
		fmt.Println("float32", t2)
	}
}

// 通过给已有类型Type起一个别名Alias，然后为Alias增加一些新的方法使其成为一个新的类型，同时Alias将完全拥有Type的所有方法
type MyRect Rect

func (r *MyRect) GetArea() float64 {
	return r.width * r.height
}

func TestNewType(t *testing.T) {
	rect1 := &Rect{100, 200}
	rect1.ModifyRect()
	fmt.Println(rect1.width, rect1.height)  // 1000 1000
	fmt.Println((*MyRect)(rect1).GetArea()) // 1e+06 进行强制类型转换可以调用

	rect2 := &MyRect{100, 200}
	// 强制类型转换,// 如果不确定符号的优先级，可以用括号来约束。
	rect3 := (*Rect)(rect2)
	rect3.ModifyRect()
	fmt.Println(rect2.width, rect2.height) // 1000 1000
	fmt.Println(rect2.GetArea())           // 1e+06
}
