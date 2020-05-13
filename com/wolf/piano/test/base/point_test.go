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

func TestPoint(t *testing.T) {
	/* 声明一般变量 */
	var a int = 10
	fmt.Printf("变量a的地址: %x\n", &a) // 取地址符是&，放到一个变量前使用就会返回相应变量的内存地址
	fmt.Printf("变量a的值: %d\n", a)

	// 一个指针变量指向了一个值的内存地址
	// 指针声明格式如下： var var_name *var-type
	// var-type 为指针类型，var_name 为指针变量名，*号用于指定变量是作为一个指针
	var ip *int /* 声明指针变量 */

	ip = &a // 将a变量的地址赋值给指针变量p
	fmt.Printf("ip变量储存的指针地址: %x\n", ip)

	//在指针类型前面加上 * 号（前缀）来获取指针所指向的内容
	fmt.Printf("*ip 变量的值: %d\n", *ip)

	// 当一个指针被定义后没有分配到任何变量时，它的值为 nil。
	// nil 指针也称为空指针。在概念上和其它语言的null、None、nil、NULL一样，都指代零值或空值。
	var ptr *int
	fmt.Printf("ptr地址为%v, ptr 的值为 : %x\n", &ptr, ptr)

	// 空指针判断：
	if ptr == nil {
		fmt.Println("ptr是空指针")
	} else {
		fmt.Println("ptr不是空指针")
	}
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
