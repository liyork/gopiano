package base

import (
	"fmt"
	"testing"
)

// 总结：
// 变量关联一个值，存放值就需要占用内存，变量和值用一个地址，用&可以取到，可以操作变量对值进行获取或修改。
// 指针同样也是变量，自己也占用地址，而指针的内容则是一个地址。可以用*p进行解指针，即找到指针的内容对应的地址，后续就可以当做变量操作目标了。

// 变量是一种使用方便的占位符，用于引用计算机内存地址，普通的变量占用的内存地址上放数值，而指针变量占用的内存地址上放地址(地址是其他值的地址)

// 在Go语言中，对于布尔变量或数值类型或字符串类型或数组都是按照值传递的：值在传递给函数或者方法时会被复制一份，然后方法或函数使用的是复制的这份值，也就不会对原值产生什么影响。一般情况下，对于布尔变量或数值类型或字符串类型的按值传递是非常廉价的，Go语言编译器会在传递过程中进行安全优化。
//
//但是在Go语言中，字符串是不可变的，因此在进行修改字符串时（例如使用+=操作），Go语言必须创建一个新的字符串，然后复制原始的字符串并将其添加到新字符串之后，对于大字符串来说，操作的代价可能会比较大。
//
//对于大字符串是这样，对于数组进行值传递也是如此。为了解决可能产生的巨大代价，Go语言使用数组切片来代替数组的使用。传递一个切片的代价跟传递字符串差不多，无论该切片的长度或容量是多大。对切片进行复制修改操作也不会像字符串那样需要创建新的切片，因为切片是可变的，属于引用类型。
//
//Go语言中的引用类型有：映射（map），数组切片（slice），通道（channel），方法与函数。
//如果我们想要管理本地变量的生命周期呢？这时就需要使用指针来管理本地变量，只要该变量至少存在一个指针，那么该变量的生命周期就可以独立于作用域。
//
//使用指针能让我们控制变量的生命周期，不受作用域的影响，另外变量在传递过程中成本最小化，且可以轻易的修改变量的内容，而不是对复制的值进行操作。指针是一个变量，这个变量实际上是保存了另一个变量的内存地址，任何被指针保存了内存地址的变量都可以通过指针来修改内容。指针的传递非常廉价。
//
//
//在使用指针前，我们需要明白两个操作符的含义
//①操作符& : 当作二元操作符时，是按位与操作；当作一元操作符时，是返回该变量的内存地址。
//②操作符* : 当作二元操作符时，是相乘的操作；当作一元操作符（解引用操作符）时，是返回该指针指向的变量的值，其实就是解除变量的指针引用，返回该变量的值。

//#指针
//一个变量存储的是一个值，但是这个值在内存中有一个地址，而指针保存的就是这个地址，通过这个地址，可以获取到值
func TestBasePointer(t *testing.T) {
	//定义的一个a变量
	var a int = 123
	//通过&符号，为a生成一个指针，是*int类型的变量接收
	p := &a
	// 通过 * 符号可以将 p 指针指向的地址的值取出来
	fmt.Println(*p)
	// 指针类型，由于 a 变量是一个 int 类型，所以指针 p 的类型是一个 int 类型的指针，通过打印可以获取到类型：*int
	fmt.Printf("%T\n", p)  // *int
	fmt.Printf("%T\n", *p) // int
	// 通过 new() 函数对特定类型创建一个指针类型
	b := new(string)
	*b = "修改 b 指针的原始值"
	fmt.Printf("%T\n", b) // 输出 *string， 是一个 string 类型的指针

	var cat int = 1                 //声明整型变量 cat。
	var str string = "banana"       //声明字符串变量 str。
	fmt.Printf("%p %p", &cat, &str) //打印 cat 和 str 变量的内存地址，指针的值是带有0x十六进制前缀的一组数据
}

func TestPoint(t *testing.T) {
	/* 声明一般变量 */
	var a int = 10
	fmt.Printf("变量a的地址: %x\n", &a) // 取地址符是&，放到一个变量前使用就会返回相应变量的内存地址
	fmt.Printf("变量a的地址: %p\n", &a) // 相比于%x多了Ox
	fmt.Printf("变量a的值: %d\n", a)

	// 一个指针变量指向了一个值的内存地址
	// 指针声明格式如下： var var_name *var-type
	// var-type 为指针类型，var_name 为指针变量名，*号用于指定变量是作为一个指针
	var ip *int /* 声明指针变量 */

	ip = &a // 将a变量的地址赋值给指针变量p
	fmt.Printf("ip变量储存的指针地址: %x\n", ip)

	//在指针类型前面加上 * 号（前缀）来获取指针所指向的内容
	fmt.Printf("*ip 变量的值: %d\n", *ip)
}

// The of a pointer is nil. That means any uninitialized pointer will have the value nil. Let's see a complete example -
// 变量声明而没有赋值，默认为零值，不同类型零值不同，例如字符串零值为空字符串；
// 指针声明而没有赋值，默认为nil(nil 指针也称为空指针)，即该指针没有任何指向。当指针没有指向的时候，不能对(*point)进行读写操作，否则会报空指针异常。
func TestPointNil(t *testing.T) {
	var ptr *int
	fmt.Printf("ptr地址为%v, ptr 的值为 : %#v\n", &ptr, ptr)
	fmt.Println("ptr = ", ptr)
	// panic, bad access: nil dereference，读写都报错
	//fmt.Println(*ptr, )
	//*ptr = 1

	// 空指针判断：
	if ptr == nil {
		fmt.Println("ptr是空指针")
	} else {
		fmt.Println("ptr不是空指针")
	}
}

// 使用==运算符来比较两个相同类型的指针是否相等。
func TestPointerCompare(t *testing.T) {
	var a = 3014
	var p1 = &a
	var p2 = &a

	if p1 == p2 {
		fmt.Printf("compare p1(%x) = p2(%x) value(%d): \n", p1, p2, *p1)
	}

	//var b = 3014
	//var p3 = &b
	//sprintf := fmt.Sprintf("%x", &a)
	//i, _ := strconv.ParseInt(sprintf, 10, 64)
	//*p3 = int(i)
	//if p1 == p3 {
	//	fmt.Printf("compare p1(%x) = p3(%x) value(%d): \n", p1, p3, *p1)
	//}
}

// 为什么nil访问方法可以，访问成员变量就会crash呢？
// nil is a predeclared identifier representing the zero value for
// a pointer, channel, func, interface, map, or slice type.
//
// var nil Type // Type must be a pointer, channel, func, interface, map, or slice type
// 指针类型的nil没有分配内存空闲，对于方法，不需要存储空间，而成员变量需要内存空间存放，所以当nil访问成员变量时，由于引用了无效的内存，所以crash.
func TestPoint2(t *testing.T) {
	//test1(nil)// 可以
	test2(nil) // 报错
}

type Plane struct {
	Num int
}

func (this *Plane) Fly1() {
	fmt.Println("Fly1......")
}

func (this *Plane) Fly2() {
	fmt.Println("Fly2......Num:", this.Num)
}

func test1(pl *Plane) {
	pl.Fly1()
}

func test2(pl *Plane) {
	pl.Fly2()
}

// 传递一个 int 类型的参数，传递的其实这个参数的一个副本；传递一个指针类型的参数，其实传递的是这个指针的一份拷贝，而不是这个指针指向的值
//任何存放在内存里的东西都有自己的地址，指针也不例外，它虽然指向别的数据，但是也有存放该指针的内存
//声明了一个变量i，值为 10，它的内存存放地址是 0xc000056080，通过这个地址我们可以找到变量 i，
//指针ip是一个指针类型的变量，它也需要内存存放它，它的内存地址是 0xc00007e018，里面的内容是0xc000056080(即i的地址)
//在我们传递指针变量 ip 给 modify 函数的时候，是对这个指针变量的一个拷贝，所以新拷贝的指针变量的内
//存地址就变了，是新的 0xc00007e028，不管是 0xc00007e018 还是 xc00007e028 我们都可以称之为指针的指针
//他们的内容都是0xc000056080，这个 0xc000056080就是变量i的地址，这也就是为什么我们可以修改变量i 的值
// 传递指针
func modify(ip *int) {
	fmt.Printf("函数里接收指针的内存地址是：%p\n", &ip)
	*ip = 1
}

func TestPointBase(t *testing.T) {
	i := 10
	ip := &i

	fmt.Printf("原始指针的内存地址是：%p\n", &ip)

	modify(ip)
	fmt.Println("int 值被修改了,新值是：", i)

}

//在 Go 语言中没有引用传递,Go 语言中所有的传参都是值传递，都是一个副本、一个拷贝。
// 因为拷贝的内容是非引用类型(int、string、struct等这些)，这样就在函数中的数据是一个副本，无法修改原内容数据的；
// 若传递的是指针，虽然对指针拷贝进方法，但是两个指针指向的数据是一个，所以修改原内容数据是否可以修改原内容数据，
type Person1 struct {
	Name string
}

func modify1(p *Person1) {
	fmt.Printf("函数里接收到的Person内存地址：%p\n", &p)
	p.Name = "李四"
}

func TestPerson1(t *testing.T) {
	p := &Person1{"张三"}
	fmt.Printf("原始Person的内存地址是：%p\n", &p)

	modify1(p)
	fmt.Println(p)
}

func TestBaseValuePass(t *testing.T) {
	a := 3
	b := 4
	c := "abc"
	d := [3]int{1, 2, 3}
	fmt.Printf("main方法：a的值为 %v,b的值为 %v,c的值为 %v,d的值为 %v \n", a, b, c, d)
	demo(a, b, c, d)
	fmt.Printf("main方法：a的值为 %v,b的值为 %v,c的值为 %v,d的值为 %v \n", a, b, c, d)
}

func demo(a, b int, c string, d [3]int) {
	a = 5
	b = 6
	c = "efg"
	d[0] = 0
	fmt.Printf("demo函数：a的值为 %v,b的值为 %v,c的值为 %v,d的值为 %v\n", a, b, c, d)
}

// *p和变量a的值是相等的,两者都与同一块内存地址相关联，任意一个变量进行修改操作都会影响到另一个变量的值，但是若变量p被赋值其他变量的指针就不行了。
func TestBasePoint(t *testing.T) {
	a := 3
	p := &a //这里是获取变量a的内存地址，并将其赋值给变量p
	fmt.Printf("a的值为:%v, a的地址是:%v ，p的值是:%v，p的地址是:%v, p指向的变量的值:%v\n",
		a, &a, p, &p, *p)
}

func TestPointUsage(t *testing.T) {
	a := 3
	b := 4
	p1 := &a //获取变量a的内存地址，并将其赋值给变量p1
	p2 := &b //获取变量b的内存地址，并将其赋值给变量p2
	fmt.Printf("a的值为 %v, a的指针是 %v ，p1指向的变量的值为 %v\n", a, p1, *p1)
	fmt.Printf("b的值为 %v, b的指针是 %v ，p2指向的变量的值为 %v\n", b, p2, *p2)

	fmt.Println(demo2(p1, p2))
	fmt.Printf("a的值为 %v, a的指针是 %v ，p1指向的变量的值为 %v\n", a, p1, *p1)
	fmt.Printf("b的值为 %v, b的指针是 %v ，p2指向的变量的值为 %v\n", b, p2, *p2)
}

func demo2(a, b *int) int {
	*a = 5
	*b = 6
	return *a + *b
}

func TestMultiRef(t *testing.T) {
	a := 3
	p1 := &a  //p1是指向变量a内存地址的指针
	p2 := &p1 //p2是指向变量p1内存地址的指针
	fmt.Printf("a:%v, p1:%v, *p1:%v, p2:%v, **p2:%v\n", a, p1, *p1, p2, **p2)
	a = 4
	fmt.Printf("a:%v, p1:%v, *p1:%v, p2:%v, **p2:%v\n", a, p1, *p1, p2, **p2)
}

// Go语言打印指向person的指针时，会打印person属性的具体内容，并且在前缀上加上&表示该变量是一个指针。
func TestNew(t *testing.T) {
	person1 := Person2{"zhangsan", "man", 25} //创建一个person1对象

	person2 := new(Person2) //使用new创建一个person2对象，同时获得指针
	person2.name, person2.sex, person2.age = "wangwu", "man", 25

	person3 := &Person2{"lisi", "man", 25} //使用&创建一个person3对象，同时获得指针
	fmt.Printf("person1:%v, person2:%v, person3:%v\n", person1, person2, person3)
}

func TestPassPoint2(t *testing.T) {
	person1 := Person2{"zhangsan", "man", 25} //创建一个person1对象
	fmt.Printf("person1:%v\n", person1)
	demo3(&person1)
	fmt.Printf("person1:%v\n", person1)
}

//两种解引用（就是将一个指针转化为原对象）的方式，
// 第一种： (*person).age是显示的解引用，
// 第二种是使用 "." 操作符自动的将指针解引用，更简单。
func demo3(person *Person2) {
	(*person).age = 18     //显示的解引用
	person.name = "GoLang" //隐式的解引用
}

// 交换的是 a 和 b 的地址，在交换完毕后，a 和 b 的变量值确实被交换。但和 a、b 关联的两个变量并没有实际关联。
func TestSwap(t *testing.T) {
	a := 1
	b := 2
	// 取出 x 和 y 的地址作为参数传给 swap() 函数进行调用
	fmt.Println(&a, a, &b, b)
	//swap(&a, &b)
	swap2(&a, &b)
	fmt.Println(&a, a, &b, b)
}

func swap(a, b *int) { //类型都为 *int 指针类型
	// *a在等号右边为取地址的值,取a指针的值, 赋给临时变量t,t 此时是 int 类型
	t := *a
	// *a在等号左边为用值赋值给指定地址
	// 取 b 的指针值，赋给指针 a 指向的变量。注意，此时*a的意思不是取 a 指针的值，而是“a 指向的变量”
	*a = *b //
	// 将 t 的值赋给指针 b 指向的变量。
	*b = t
}

// 参数到这里被复制，x,y是另外两个指针变量，
func swap2(x, y *int) {
	fmt.Println("param1:==>", &x, x, &y, y)
	// 交换了下x和y的值，由于x和y是a，b的拷贝，所以并不会影响a和b以及a、b指向的内容
	x, y = y, x
	fmt.Println("param2:==>", &x, x, &y, y)
}

// 涉及取值然后赋值指针所指向的地址的内容
func swap3(a, b *int) {
	*b, *a = *a, *b
}

type Teacher1 struct {
	High int32
	Sex  int32
}

type Teacher2 struct {
	*Teacher1
}

func (t *Teacher1) String() string {
	return fmt.Sprintf("{High:%d,Sex:%d}", t.High, t.Sex)
}

func TestPrintPointer(t *testing.T) {
	// 结构体指针的数组2
	arr2 := []*Teacher1{
		&Teacher1{High: 170, Sex: 17},
		&Teacher1{High: 180, Sex: 18},
	}
	fmt.Printf("打印结构体指针数组2：%+v \n", arr2)

	teacher2 := Teacher2{Teacher1: &Teacher1{High: 1111, Sex: 111}}
	fmt.Printf("teacher:%+v\n", teacher2)
}

type person struct {
	name string
}

// tom是person对象的指针
var p = person{name: "xxxx"}
var tom *person = &p

func TestPointer(t *testing.T) {
	// 使用%+v来打印tom指针(打印的是tom指向的内容)
	fmt.Printf("%+v\n", tom)
	// &tom创建一个指针，指向tom，打印是&tom指针指向的内容
	fmt.Printf("%+v\n", &tom)

	// %p打印指针地址
	fmt.Printf("%p\n", tom)
	fmt.Printf("%p\n", &tom)

	//fmt.Printf("%+v\n", p)
	//fmt.Printf("%p\n", &p)
}

type simpleHost struct {
	name string
}

func TestPoint22222(t *testing.T) {
	// s是个变量，有自己的地址
	s := simpleHost{}
	// s的地址，打印指针的值
	fmt.Printf("%p\n", &s)

	host := &s //host也是变量，也有地址，只不过里面的值是s的地址
	// host的地址，host里面存放的是s的地址
	fmt.Printf("%p, %p\n", &host, host)
	test2222(host)
}

func test2222(host *simpleHost) {
	// host的地址--这里由于传参导致值被复制，host里面存放的是s的地址
	fmt.Printf("%p, %p\n", &host, host)
}

func TestPointer2(t *testing.T) {
	var house = "Malibu Point 10880, 90265"
	// 对字符串取地址, 将指针保存到变量 ptr 中,ptr类型为*string
	ptr := &house
	// 打印ptr的类型
	fmt.Printf("ptr type: %T\n", ptr)
	// ptr变量占用的地址
	fmt.Printf("address &ptr: %v\n", &ptr)
	// 打印ptr的指针地址
	fmt.Printf("address p: %p\n", ptr)
	fmt.Printf("address v: %v\n", ptr)
	// 对 ptr 指针变量进行取值操作
	value := *ptr
	// 取值后的类型
	fmt.Printf("value type: %T\n", value)
	// 指针取值后就是指向变量的值
	fmt.Printf("value: %s\n", value)
}

// 如果一个指针变量存放的又是另一个指针变量的地址，则称这个指针变量为指向指针的指针变量。
//当定义一个指向指针的指针变量时，第一个指针存放第二个指针的地址，第二个指针存放变量的地址：
//Pointer      Pointer     Variable
//Address -->  Address --> Value
//
//指向指针的指针变量声明格式如下：指向指针的指针变量为整型
//var ptr **int;
func TestPointerPointer(t *testing.T) {
	// 变量
	var a int = 3000
	// *int类型指针
	var ptr *int
	// **int类型指针
	var pptr **int

	// 指针ptr，持有a内存地址
	ptr = &a

	// 指针pptr，持有ptr内存地址
	pptr = &ptr

	fmt.Printf("变量 a = %d\n", a)
	fmt.Printf("取出指针变量*ptr指向地址的值 = %d\n", *ptr)
	fmt.Printf("指向指针的指针变量，取值 **pptr = %d\n", **pptr)
}

type PersonPointer struct {
	name string
	age  int
}

func (p PersonPointer) sayHi() {
	fmt.Printf("SayHi -- This is %s, my age is %d\n", p.name, p.age)
}
func (p PersonPointer) ModifyAge(age int) {
	p.age = age
	// 一个拷贝
	fmt.Printf("ModifyAge person:%p, <%s:%d>\n", &p, p.name, p.age)
}

func (p *PersonPointer) ChangeAge(age int) {
	p.age = age
	// 也是一个拷贝，不过p这个调用时的p都是指针，都指向的一个内存地址
	fmt.Printf("ChangeAge person:%p, %p, <%s:%d>\n", &p, p, p.name, p.age)
}

// 方法即为有接受者的函数，接受者可以是类型的实例变量或者是类型的实例指针变量
// 函数方法的接受者，也可以是指针变量。无论普通接收者还是指针接收者都会被拷贝传入方法中，不同在于拷贝的指针，其指向的地方都一样，只是其自身的地址不一样。
func TestVariableAndPointerReceiver(t *testing.T) {
	person := PersonPointer{"xxx", 21}
	fmt.Printf("person:%p, %s, %d \n", &person, person.name, person.age)

	person.sayHi()
	person.ModifyAge(210)
	// 编译器会改成(&person).ChangeAge
	person.ChangeAge(210)
	person.sayHi()
}
