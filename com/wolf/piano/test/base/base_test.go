package base

import (
	"fmt"
	"math/rand"
	"testing"
)

// 指针储存的是一个值的地址，指针本身也需要地址来储存
func Test_point(t *testing.T) {
	var p *int   // 指针声明
	p = new(int) // 指针赋值
	*p = 1       // *表示取指针存储地址的值，赋值
	fmt.Println(p, &p, *p)
	// 0xc0000901a8,p指针存储的地址
	// 0xc000088018,&代表取指针本身的地址
	// *得到指针存储地址的值，1
	// 总结：对于指针p,直接打印p得到的是p存储的地址，&p是p指针本身占用内存地址，*p是指针存储地址的值
	//地址	0xc000088018	0xc0000901a8
	//值	0xc0000901a8	1
}

func Test_pointErr(t *testing.T) {
	var i *int
	// i = new(int) --解决问题
	*i = 1
	fmt.Println(i, &i, *i)

	// panic: runtime error: invalid memory address or nil pointer dereference [recovered]
	//	panic: runtime error: invalid memory address or nil pointer dereference
	// [signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x10f899c]
	// 指针i，占用内存，存储的值为nil，*i为取指针i的存储值对应内存的值，为nil并没有给其分配内存，然后赋值1报错
}

func TestRand(t *testing.T) {
	// 生成随机数，但是数值不会变，每次生成都是固定数字
	//a := rand.Int()
	//fmt.Println(a)
	//b := rand.Intn(100) //[0,100)
	//fmt.Println(b)

	// 生成可变随机数
	//rand.Seed(time.Now().UnixNano())
	////生成10个0-99之间的随机数
	//for i := 0; i < 10; i++ {
	//	fmt.Println(rand.Intn(100))
	//}
}

func TestAssign(t *testing.T) {
	var x1 = []string{"1", "2"}
	x2 := x1
	fmt.Println(x1, x2)
	fmt.Println(&x1, &x2)
	fmt.Printf("x1 pointer add:%p, x2 pointer add:%p\n", &x1, &x2)
	x1[0] = "a"
	fmt.Printf("x1:%v,x2:%v\n", x1, x2) // 这不是拷贝了，是相同的底层数据，拷贝的是上面的slice

	var x3 = []person{{"1"}, {"2"}}
	x4 := x3
	fmt.Printf("x3 pointer add:%p, x4 pointer add:%p\n", &x3, &x4)
	x3[0].name = "333"
	// 看来:=是浅拷贝
	fmt.Printf("x3 value:%v, x4 value:%v\n", x3, x4)
}

// 报错
func TestRand0(t *testing.T) {
	intn := rand.Intn(0)
	fmt.Println("int:", intn)
}

// 使用 _ 下划线符号作为别名来忽略导入的包，从而避免编译错误，这只会执行 package 的 init()

// 简短声明(赋值)
func TestShortDeclare(t *testing.T) {
	// := 左侧至少有一个新变量，才允许多变量的重复声明
	one := 0
	one, two := 1, 2
	one, two = two, one // 交换两个变量值的简写
}

func TestShortDeclareOverwrite(t *testing.T) {
	x := 1
	println(x)

	{
		println(x)
		// 新的 x 变量的作用域只在代码块内部
		x := 2
		println(x)
	}

	println(x)
}

// nil 是 interface、function、pointer、map、slice 和 channel 类型变量的默认初始值。
// 但声明时不指定类型，编译器也无法推断出变量的具体类型。
// 允许对值为 nil 的 slice 添加元素，但对值为 nil 的 map 添加元素则会造成运行时 panic
