package base

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

// 在 Go 语言中当一个变量被声明的时候，系统会自动初始化它的默认值，比如 int 被初始化为 0，指针为 nil。
// var 声明同样也会为结构体类型的数据分配内存，所以才能在声明了 var s T 之后就能直接给他的字段进行赋值
// 表达式 new(Type) 和 &Type{} 是等价的

// 有两个原因需要使用指针接收者。首先避免在每个方法调用中拷贝值。其次，方法可以修改接收者指向的值。

type A struct {
	name string
	age  int
}

func (a A) IsEmpty() bool {
	return reflect.DeepEqual(a, A{})
}

func TestStructNil(t *testing.T) {
	var a A

	if a == (A{}) { // 括号不能去
		fmt.Println("a == A{} empty")
	}

	if a.IsEmpty() {
		fmt.Println("reflect deep is empty")
	}
}

// 结构体实例化后字段的默认值是字段类型的默认值，例如 ，数值为 0、字符串为 ""（空字符串）、布尔为 false、指针为 nil 等
func TestInitStruct(t *testing.T) {

	type People struct {
		name  string
		child *People
	}

	// 指定字段初始化
	relation := &People{ // 由 People类型取地址后，形成类型为 *People 的实例
		name: "爷爷",
		child: &People{
			name: "爸爸",
			child: &People{
				name: "我",
			},
		},
	}

	fmt.Println("relation:", relation)

	// name默认"",child为nil
	p := &People{} // 一种简写，底层仍会调用 new()
	fmt.Println("p:", p)

	// 全部字段初始化
	type Address struct {
		Province    string
		City        string
		ZipCode     int
		PhoneNumber string
	}

	addr := Address{
		"四川",
		"成都",
		610000,
		"0",
	}
	fmt.Println(addr)
}

func TestAnoyStruct(t *testing.T) {

	// 实例化一个匿名结构体
	msg := &struct { // 定义部分
		id   int
		data string
	}{ // 值初始化部分
		1024,
		"hello",
	}
	printMsgType(msg)

}

// 打印消息类型, 传入匿名结构体
func printMsgType(msg *struct {
	id   int
	data string
}) {
	// 使用动词%T打印msg的类型
	fmt.Printf("%T\n", msg)
}

// 空结构体的特点:不占用内存,地址不变
// 空结构体作用:建议用于传递信号的通道，因为不占用内存
func TestEmptyStruct(t *testing.T) {
	//空结构体的特点:1、不占用内存；2、地址不变
	var s struct{}
	var s1 struct{}
	fmt.Println("空结构体占用内存的情况：", unsafe.Sizeof(s))
	fmt.Printf("空结构体指针指向情况:s = %p, s1 = %p,两个指针的比较结果：%v\n", &s, &s1, &s == &s1)

	cacheChan := make(chan string, 3)
	receiveChan := make(chan struct{}, 1)
	finishChan := make(chan struct{}, 2) //操作完成信号

	go func() {
		<-receiveChan //阻塞协程，直到接收到值
		for value := range cacheChan {
			fmt.Println("接收到值为：", value)
		}
		finishChan <- struct{}{}
	}()

	go func() {
		// 模拟发送数据
		for index, value := range []string{"1", "2", "3"} {
			fmt.Println("发送数据：", value)
			cacheChan <- value
			if index == 2 {
				receiveChan <- struct{}{}
			}
		}
		close(cacheChan)
		finishChan <- struct{}{}
	}()

	fmt.Println("阻塞，等待上面两个协程运行结束")
	<-finishChan
	<-finishChan
}

// 结构体是没有位段的结构体
// 宽度描述了存储一个数据类型实例需要占用的字节数，由于进程的内存空间是一维的，我更倾向于将宽度理解为Size
// 宽度是数据类型的一个属性。Go程序中所有的实例都是一种数据类型，一个实例的宽度是由他的数据类型决定的，通常是8bit的整数倍。
func TestUseEmptyStruct1(t *testing.T) {
	// 数组的宽度是他元素宽度的整数倍。
	var a [3]uint32
	fmt.Println(unsafe.Sizeof(a)) // prints 12

	// 结构体提供了定义组合类型的灵活方式，组合类型的宽度是字段宽度的和，然后再加上填充宽度。
	var s1 S
	fmt.Println(unsafe.Sizeof(s1)) // prints 8, not 6

	// 空结构体的宽度是0，他占用了0字节的内存空间。
	var s3 struct{}
	fmt.Println(unsafe.Sizeof(s3)) // prints 0

	// 由于空结构体占用0字节，那么空结构体也不需要填充字节。所以空结构体组成的组合数据类型也不会占用内存空间。
	type S struct {
		A struct{}
		B struct{}
	}

	var s4 S
	fmt.Println(unsafe.Sizeof(s4)) // prints 0

	//由于Go的正交性，空结构体可以像其他结构体一样正常使用。正常结构体拥有的属性，空结构体一样具有。
	//可以定义一个空结构体组成的数组，当然这个切片不占用内存空间。
	var x [1000000000]struct{}
	fmt.Println(unsafe.Sizeof(x)) // prints 0

	// 空结构体组成的切片的宽度只是他的头部数据的长度，就像上例展示的那样，切片元素不占用内存空间。
	var x1 = make([]struct{}, 1000000000)
	fmt.Println(unsafe.Sizeof(x1))

	// 切片的内置子切片、长度和容量等属性依旧可以工作。
	var x2 = make([]struct{}, 100)
	var y = x2[:50]
	fmt.Println(len(y), cap(y)) // prints 50 100

	// 可以寻址一个空结构体，空结构体是可寻址的，就像其他类型的实例一样。
	var a1 struct{}
	var b1 = &a1
	fmt.Println(b1)
	// 两个空结构体的地址可以相等。
	var a2 struct{}
	var b2 struct{}
	fmt.Println("&a2 == &b2", &a2 == &b2)

	// 空结构体的元素也具有一样的属性。
	a3 := make([]struct{}, 10)
	b3 := make([]struct{}, 20)
	fmt.Println(&a3 == &b3)       // false, a and b are different slices
	fmt.Println(&a3[0] == &b3[0]) // true, their backing arrays are the same

	a4 := struct{}{} // not the zero value, a real new struct{} instance
	b4 := struct{}{}
	fmt.Println(a4 == b4) // true
}

func TestUseEmptyStruct2(t *testing.T) {
	// 空结构体可以作为方法的接收者。
	var a, b S
	a.addr() // 0x1beeb0
	b.addr() // 0x1beeb0

	//golang的struct{}
	//做控制而非数据信息： chan struct{}
	//实现set: map[string]struct{}

}

type S struct{}

func (s *S) addr() { fmt.Printf("%p\n", s) }

// connReader 连接reader
type connReader struct {
	closeChan chan struct{}
}

// Close 发送close信号 关闭包decoder
func (r *connReader) Close() error {
	if r.closeChan == nil {
		return nil
	}
	r.closeChan <- struct{}{} //空结构体还可以写入chan管道，作为触发器，用其他冗余数据的话，会影响程序性能
	r.closeChan = nil
	return nil
}

func TestIsEmpty(t *testing.T) {
	//cr := &connReader{}
	var cr *connReader
	if cr == nil {
		fmt.Println("cr is null")
	} else {
		fmt.Println("cr is not null")
	}
}

//在介绍面向对象时，讲过可以通过属性和方法（函数）来描述对象
type Integer int //:表示的意思是给int类型指定了一个别名叫Integer,后面可以用Integer来代替int 来使用
//定义了一个方法，方法的定义与函数的区别:
// 第一：在关键字func后面加上(a Integer), 这个在方法中称之为接收者，即接收传递过来的第一个参数，然后复制a，a的类型是Integer ,
// 由于Integer是int的别名，所以a的类型为int
//第二：在表示参数的类型时，都使用了对应的别名。
//通过方法的定义，可以看出方法其实就是给某个类型绑定的函数。在该案例中，是为整型绑定的函数，只不过在给整型绑定函数(方法)时，
// 一定要通过type来指定一个别名，因为int类型是系统已经规定好了，无法直接绑定函数，所以只能通过别名的方式。
//第三:调用方式不同
//var result Interger=3
//表示定义一个整型变量result，并赋值为3.
//result.Test( 3)
//通过result变量，完成方法的调用。因为，Test方法，是为int类型绑定的函数，而result变量为int类型。所以可以调用Test( )方法。
// result变量的值会传递给Test( )方法的接受者，也就是参数a, 而实参Test( 3),会传递形参b.
//当然，我们也可以将Test( )方法，理解成是为int类型扩展了，追加了的方法。因为系统在int类型时，是没有改方法的。
func (a Integer) Test(b Integer) Integer {
	fmt.Println("b:", b)
	return b + 1
}

func TestIntegerMethod(t *testing.T) {
	var sum Integer = 6
	test := sum.Test(1)
	fmt.Println(test)
}

type Student1 struct {
	id   int
	name string
}

func (s Student1) print1() {
	fmt.Println(s)
}

func (s *Student1) print2() {
	fmt.Println(s)
}

func TestStructMethod(t *testing.T) {
	student := Student1{101, "xxx"}
	student.print1()
	// 为什么结构体指针变量，可以调用print2方法呢？
	// 原因是：先转换成*student再调用
	student.print2()

	studentPoint := &student
	studentPoint.print1()
	studentPoint.print2()
}

func TestPrintfStruct(t *testing.T) {
	type Q struct {
		a int
		b float64
		c string
	}
	q := &Q{7, -2.35, "abc\tdef"}
	fmt.Printf("%v\n", q)
	fmt.Printf("%+v\n", q)
	fmt.Printf("%#v\n", q)
}
