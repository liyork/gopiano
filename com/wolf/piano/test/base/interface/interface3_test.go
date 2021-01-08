package _interface

import (
	"fmt"
	"io"
	"testing"
	"unsafe"
)

// 先定义一个接口，和使用此接口作为参数的函数：
type IGreeting interface {
	sayHello()
}

func sayHello(i IGreeting) {
	i.sayHello()
}

// 定义两个结构体：
type Go struct{}

func (g Go) sayHello() {
	fmt.Println("Hi, I am GO!")
}

type PHP struct{}

func (p PHP) sayHello() {
	fmt.Println("Hi, I am PHP!")
}

// 调用 sayHello() 函数，传入了 golang、 php 对象，它们并没有显式地声明实现了 IGreeting 类型，只是实现了接口所规定的 sayHello() 函数。
// 实际上，编译器在调用 sayHello() 函数时，会隐式地将 golang、php 对象转换成 IGreeting 类型，这也是静态语言的类型检查功能
func TestInterfaceBase3(t *testing.T) {
	golang := Go{}
	php := PHP{}
	sayHello(golang)
	sayHello(php)
}

type Person3 struct {
	age int
}

func (p Person3) howOld() int {
	return p.age
}

func (p *Person3) growUp() {
	p.age += 1
}

// 调用了 growUp 函数后，不管调用者是值类型还是指针类型，它的 Age 值都改变了。
// 实际上，当调用时类型和方法的接收者类型不同时，其实是编译器在背后做了一些工作:
// -			值接收者									指针接收者
// 值类型调用者	方法会使用调用者的一个副本，类似于“传值”		使用值的指针调用方法，上例中，qcrao.growUp() 实际上是 (&qcrao).growUp()
// 指针类型调用者  指针被解引用为值，上例中，stefno.howOld() 实际上是 (*stefno).howOld() 实际上也是“传值”  方法里的操作会影响到调用者，类似于指针传参，拷贝了一份指针
func TestStructReceiver(t *testing.T) {
	// ====调用者是值类型
	qcrao := Person3{age: 18}
	// 值类型 调用接收者也是值类型的方法
	fmt.Println(qcrao.howOld())
	// 值类型 调用接收者是指针类型的方法
	qcrao.growUp()
	fmt.Println(qcrao.howOld())

	// ====调用者是指针类型
	stefno := &Person3{age: 100}
	// 指针类型 调用接收者是值类型的方法
	fmt.Println(stefno.howOld())
	// 指针类型 调用接收者也是指针类型的方法
	stefno.growUp()
	fmt.Println(stefno.howOld())
}

type coder interface {
	code()
	debug()
}

type Gopher struct {
	language string
}

func (p Gopher) code() {
	fmt.Printf("I am coding %s language\n", p.language)
}

func (p *Gopher) debug() {
	fmt.Printf("I am debuging %s language\n", p.language)
}

// 隐式自动方法实现
// Gopher 类型并没有实现 debug 方法；
// 虽然*Gopher 类型也没有实现 code 方法，但是因为 Gopher 类型实现了 code 方法，所以让 *Gopher 类型自动拥有了 code 方法。
func TestInterfaceReceiver(t *testing.T) {
	// ok
	//var c coder = &Gopher{"Go"}
	//c.code()
	//c.debug()

	// 报错
	// Cannot use 'Gopher{"Go"}'(type Gopher)as type coder in assignment
	// Type does not implement 'coder' as 'debug' method has a pointer receiver
	//var c coder = Gopher{"Go"}
	//c.code()
	//c.debug()

	// 直接使用struct而非interface
	//var c = Gopher{"Go"}
	//c.code()
	//c.debug()
}

type Coder2 interface {
	code2()
}

type Gopher2 struct {
	name string
}

func (g Gopher2) code2() {
	fmt.Printf("%s is coding\n", g.name)
}

func TestInterfaceValueNil(t *testing.T) {
	var c Coder2
	// c 的 动态类型和动态值都为 nil
	fmt.Println(c == nil)
	fmt.Printf("c: %T, %v\n", c, c)

	// g 也为 nil
	var g *Gopher2
	fmt.Println(g == nil)

	// 当把 g 赋值给 c 后，c 的动态类型变成了 *main.Gopher，仅管 c 的动态值仍为 nil，当 c 和 nil 作比较的时候，结果就是 false 了
	c = g
	fmt.Println(c == nil)
	fmt.Printf("c: %T, %v\n", c, c)
}

type MyError struct{}

func (i MyError) Error() string {
	return "MyError"
}

// Process 函数返回了一个 error 接口，这块隐含了类型转换。
// 所以，虽然它的动态值是 nil，但它的类型是 *MyError，最后和 nil 比较的时候，结果为 false。
func TestInterfaceValueNil2(t *testing.T) {
	err := Process()
	fmt.Println(err)

	fmt.Println(err == nil)
}

func Process() error {
	var err *MyError = nil
	return err
}

// 直接定义了一个 iface 结构体，用两个指针来描述 itab 和 data
type iface struct {
	itab, data uintptr
}

// 打印出接口的动态类型和值
func TestPrintInterfaceTypeAndValue(t *testing.T) {
	var a interface{} = nil

	var b interface{} = (*int)(nil)

	x := 5
	var c interface{} = (*int)(&x)

	// 将 a, b, c 在内存中的内容强制解释成我们自定义的 iface
	ia := *(*iface)(unsafe.Pointer(&a))
	ib := *(*iface)(unsafe.Pointer(&b))
	ic := *(*iface)(unsafe.Pointer(&c))

	// 打印出动态类型和动态值的地址
	// a 的动态类型和动态值的地址均为 0，也就是 nil；
	// b 的动态类型和 c 的动态类型一致，都是 *int；最后，c 的动态值为 5。
	fmt.Println(ia, ib, ic)

	fmt.Println(*(*int)(unsafe.Pointer(ic.data)))
}

// 下述赋值语句会发生隐式地类型转换，在转换的过程中，编译器会检测等号右边的类型是否实现了等号左边接口所规定的函数。
// 编译器会由此检查 *myWriter 类型是否实现了 io.Writer 接口。
type myWriter struct {
}

func (w myWriter) Write(p []byte) (n int, err error) {
	return
}

func TestCompileCheck(t *testing.T) {
	// 检查 *myWriter 类型是否实现了 io.Writer 接口
	var _ io.Writer = (*myWriter)(nil)

	// 检查 myWriter 类型是否实现了 io.Writer 接口
	var _ io.Writer = myWriter{}
}

// int 型和 float64 是相互兼容的。
// 类型转换
// 转换前后的两个类型要相互兼容才行。语法为： <结果类型> := <目标类型> (<表达式>)
func TestTypeConvert(t *testing.T) {
	var i int = 9

	var f float64
	f = float64(i)
	fmt.Printf("%T, %v\n", f, f)

	f = 10.8
	a := int(f)
	fmt.Printf("%T, %v\n", a, a)

	// 注释去掉，编译器会报告类型不兼容的错误：
	// s := []int(i)
}

type Student struct {
	Name string
	Age  int
}

// Go 语言中不允许隐式类型转换，也就是说 = 两边，不允许出现类型不相同的变量。针对实际类型
// 类型转换、类型断言本质都是把一个类型转换成另外一个类型。
// 类型断言是对接口进行的操作
// 当一个函数的形参是接口，那么在函数中，需要对形参进行断言，从而得到它的真实类型。
// 断言的语法为： 安全类型断言
// <目标类型的值>，<布尔参数> := <表达式>.( 目标类型 )
func TestTypeAssert(t *testing.T) {
	student := new(Student)
	//i2 := student.(Student)
	var i interface{} = student
	s, ok := i.(Student)
	if ok {
		fmt.Println(s)
	}
}
func TestSwitchTypeAssert(t *testing.T) {
	//‍i 是一个 *Student 类型，匹配第三个 case，
	// 从打印的三个地址来看，这三处的变量实际上都是不一样的。在 main 函数里有一个局部变量 i；调用函数时，实际上是复制了一份参数，因此函数里又有一个变量 v，它是 i 的拷贝；断言之后，又生成了一份新的拷贝
	//var i interface{} = new(Student)

	// i 在这里动态类型是 (*Student), 数据为 nil，它的类型并不是 nil，它与 nil 作比较的时候，得到的结果也是 false。
	var i interface{} = (*Student)(nil)

	// i是 nil 类型。
	//var i interface{}

	fmt.Printf("%p %v\n", &i, i)

	judge(i)
}

func judge(v interface{}) {
	fmt.Printf("%p %v\n", &v, v)

	switch v := v.(type) {
	case nil:
		fmt.Printf("%p %v\n", &v, v)
		fmt.Printf("nil type[%T] %v\n", v, v)
	case Student:
		fmt.Printf("%p %v\n", &v, v)
		fmt.Printf("Student type[%T] %v\n", v, v)
	case *Student:
		fmt.Printf("%p %v\n", &v, v)
		fmt.Printf("*Student type[%T] %v\n", v, v)
	default:
		fmt.Printf("%p %v\n", &v, v)
		fmt.Printf("unknow\n")
	}
}

// String
//func (s Student) String() string {
//	return fmt.Sprintf("[Name: %s], [Age: %d]", s.Name, s.Age)
//}

// 若上面注释，用这里
func (s *Student) String() string {
	return fmt.Sprintf("[Name: %s], [Age: %d]", s.Name, s.Age)
}

// fmt.Println 函数的参数是 interface。对于内置类型，函数内部会用穷举法，得出它的真实类型，然后转换为字符串打印。而对于自定义类型，首先确定该类型是否实现了 String() 方法，如果实现了，则直接打印输出 String() 方法的结果；否则，会通过反射来遍历对象的成员进行打印。
func TestPrintString(t *testing.T) {
	var s = Student{
		Name: "qcrao",
		Age:  18,
	}

	fmt.Println(s)
	fmt.Println(&s)
	// 当Student 结构体定义了接受者类型是值类型的 String() 方法时，通过以上均可以按照自定义的格式来打印。
	//如果 Student 结构体定义了接受者类型是指针类型的 String() 方法时，只有通过指针才能按照自定义的格式打印。
}

// 多态

type Person interface {
	job()
	growUp()
}

type Student2 struct {
	age int
}

// 既可以当Student2又可以当*Student2接收者
func (p Student2) job() {
	fmt.Println("I am a student.")
	return
}

// 只能当*Student2接收者
func (p *Student2) growUp() {
	p.age += 1
	return
}

type Programmer struct {
	age int
}

func (p Programmer) job() {
	fmt.Println("I am a programmer.")
	return
}

func (p Programmer) growUp() {
	p.age += 10
	return
}

func whatJob(p Person) {
	p.job()
}

func growUp(p Person) {
	p.growUp()
}

// 直接调用接口函数，实际执行的时候是看最终传入的实体类型是什么，调用的是实体类型实现的函数。于是，不同对象针对同一消息就有多种表现，多态就实现了。
// 在函数 whatJob() 或者 growUp() 内部，接口 person 绑定了实体类型 *Student 或者 Programmer。这里会直接调用 fun 里保存的函数(已经被动态修改)，类似于：s.tab->fun[0]，而因为 fun 数组里保存的是实体类型实现的函数，所以当函数传入不同的实体类型时，调用的实际上是不同的函数实现，从而实现多态。
func TestPolymorphic(t *testing.T) {
	qcrao := &Student2{age: 18}
	whatJob(qcrao)

	// 注意，*Student2 类型实现了接口， Student 类型却没有，因为
	//qcrao2 := Student2{age: 18}
	//whatJob(qcrao2)

	growUp(qcrao)
	fmt.Println(qcrao)

	stefno := Programmer{age: 100}
	whatJob(stefno)

	growUp(stefno)
	fmt.Println(stefno)
}
