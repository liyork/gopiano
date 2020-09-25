package base

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"strings"
	"testing"
)

// interface(接口),可以定义一组方法，但是这些不需要实现。并且interface不能包含任何变量。
// interface是方法的集合
//interface是一种类型，并且是指针类型
//interface的更重要的作用在于多态实现(一种事物的多种形态,同一个interface，不同的类型实现，都可以进行调用，它们都按照统一接口进行操作)
//接口的使用不仅仅针对结构体，自定义类型、变量等等都可以实现接口。
//如果一个接口没有任何方法，我们称为空接口，由于空接口没有方法，所以任何类型都实现了空接口。
//要实现一个接口，必须实现该接口里面的所有方法。
//interface是一组method的组合，我们通过interface来定义对象的一组行为。换句话说，一个 interface 类型定义了一个“方法集合”作为其接口。 interface类型的变量可以保存含有属于这个interface类型的任何类型的值，这时我们就说这个类型实现了这个接口。未被初始化的interface类型变量的零值为空（nil）
//如果一个interface1 作为 interface2 的一个嵌入字段，那么 interface2 隐式的包含了interface1 里面的method

//非侵入式接口

//接口是golang中实现多态性的唯一好途径。
//接口是一组方法签名的集合，然后我们可以定义一个结构体实现该接口所有方法。因此，接口就是定义了对象的行为。
//例如，结构体Dog可以walk和bark， 如果一个接口声明了walk和bark的方法签名，而Dog实现了walk和bark方法，那么Dog就实现了该接口。
//接口的主要工作是仅提供由方法名称，输入参数和返回类型组成的方法签名集合。 由类型（例如struct结构体）来声明方法并实现它们。
//如果一个类型实现了在接口中定义的签名方法，则称该类型实现该接口

// 接口可以采用不同的形式。 这就是多态性的定义。 接口在需要传递给它们的许多类型的参数的函数和方法的情况下非常有用，例如接受所有类型的值的Println函数。 如果你看到Println函数的语法，就像
//当多个类型实现相同的接口时，使用相同的代码可以很容易地使用它们。 因此，只要我们可以使用接口，我们就应该尽量使用它。

// 反射是建立在类型系统(type system)上的
//Go是一门静态类型的语言。每个变量都有一个静态类型，类型在编译的时后被知晓并确定了下来
//虽然在运行时中，接口变量存储的值也许会变，但接口变量的类型是永不会变的
func TestTypeStatic(t *testing.T) {
	type MyInt int

	var i int
	var j MyInt

	// 虽然它们有着相同的基本类型，但静态类型却不一样，在没有类型转换的情况下，它们之间无法互相赋值。
	fmt.Println(reflect.TypeOf(i), reflect.TypeOf(j))
}

// 接口是一个重要的类型，它意味着一个确定的的方法集合。一个接口变量可以存储任何实现了接口的方法的具体值(除了接口本身)
func TestInterfaceType(t *testing.T) {
	// 不管变量r中的具体值是什么，r的类型永远是io.Reader：Go是静态类型的，r的静态类型就是io.Reader。
	var r io.Reader
	r = os.Stdin
	r = bufio.NewReader(r)
	r = new(bytes.Buffer)
}

// 在接口类型中有一个极为重要的例子——空接口：
//interface{}
//表示了一个空的方法集，一切值都可以满足它，因为它们都有零值或方法。

//定义接口
type Skills interface {
	Running()
	Getname() string
}

type Student struct {
	Name string
	Age  int
}

// Student实现接口
func (p Student) Running() { // 实现 Running方法
	fmt.Printf("Student %s running\n", p.Name)
}

func (p Student) Getname() string { //实现Getname方法
	fmt.Println(p.Name)
	return p.Name
}

type Teacher struct {
	Name   string
	Salary int
}

// Teacher实现接口
func (p *Teacher) Getname() string { //实现Getname方法
	fmt.Println(p.Name)
	return p.Name
}

func (p *Teacher) Running() { // 实现 Running方法
	fmt.Printf("Teacher %s running\n", p.Name)
}

func TestIntefaceBase(t *testing.T) {
	var stu1 Student
	stu1.Name = "wd"
	stu1.Age = 22

	var t1 *Teacher
	t1.Name = "wang"

	// 提升类型
	var skill Skills
	skill = stu1
	skill.Running() //调用接口

	skill = t1
	skill.Running()

}

// 继承接口
type Test interface {
	sleeping()
	Skills //嵌套，继承Skills
}

func TestIntefaceType(t *testing.T) {
	var x interface{}

	s := "WD"
	x = s
	y, ok := x.(int)
	z, ok1 := x.(string)
	fmt.Println(y, ok)
	fmt.Println(z, ok1)

	x = Student{}
	student, ok := x.(Student)
	fmt.Println(student, ok)

	teacher, ok := x.(Teacher)
	fmt.Println(teacher, ok)

}

// 一个接口类型的变量里有两样东西：变量的的具体值和这个值的类型描述。
// 这个实现了接口的值是一个基础的具体数据项，而类型描述了数据项里的所有类型。
func TestInterfaceInter(t *testing.T) {
	var r io.Reader
	tty, _ := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	r = tty // r包含了(value, type)组合，(tty, *os.File)
	fmt.Println(r, t)

	// tty的类型*os.File实现了Read以外的方法；虽然接口值只提供了Read方法，但它内置了所有的类型信息
	var w io.Writer
	// 断言了r中所包含的数据项实现了io.Writer，所以我们可以用它对w赋值
	w = r.(io.Writer) // w将与r一样，包含(tty, *os.File)组合
	fmt.Println(w)

	// 接口的静态类型决定了接口变量的哪些方法可以被外部调用，即便也许它所含的具体值有一个更大的方法集也只能提供部分。

	// 空接口变量将会在此包含同样的“组合”：(tty, *os.File)。这非常方便：一个空接口可以包含任何值和它的类型信息，我们可以在任何需要的时候了解它
	// 在这里我们无需类型断言是因为w已经满足了空接口。在前面的例子中我们将一个值从一个Reader传到了Writer，因为Writer不是Reader的子集，所以我们需要使用类型断言
	// 接口里“组合”的格式永远是（值，实体类型），而不是（值，接口类型）
	var empty interface{}
	empty = w
	fmt.Println(empty)
}

func TestAppendIface(t *testing.T) {
	//var xx []*Skills = []*Skills{}
	//
	//teacher := &Teacher{
	//	Name:   "",
	//	Salary: 0,
	//}
	//xx := append(xx, teacher)
	//fmt.Println("xx:", xx)
}

// 接口
type IAnimal interface {
	GetName()
}

type Animal struct {
	Name  string
	color int
}

// 继承
func TestInherit(t *testing.T) {
	//实例化machineCat并取指针
	machineCat := &MachineCat{
		Animal: &Animal{
			Name: "test",
		},
	}

	machineCat.GetName()
	// 可以直接访问/设定匿名继承来的属性
	fmt.Println(machineCat.Name)
	//相同属性，这里是赋值machineCat的属性Name
	machineCat.Name = "xx11"
	fmt.Println(machineCat.Name)
	fmt.Println(machineCat.Animal.Name)

	dog := &Dog{Animal{Name: "xxx"}}
	// 对于显示名称的继承只能先获取组合对象再获取属性
	fmt.Println(dog.animal.Name)
}

func TestPolymorphic(t *testing.T) {
	machineCat := newMachineCat("newMachineCat")
	var animal IAnimal = machineCat
	animal.GetName()

	animal1 := NewAnimal("newAnimal")
	animal1.GetName()

	machineCat.Animal.GetName()
}

func TestInheritNil(t *testing.T) {
	cat := &MachineCat{}
	//cat.Name = "aaa"// invalid memory address or nil pointer dereference
	// 因为指针默认值为nil
	cat.Animal = &Animal{color: 111}
	cat.Name = "bb"

	// 继承过来的属性可以直接使用，若是组合，那么需要先[.组合对象名再.属性]
	fmt.Println("color:", cat.color)
}

func NewAnimal(name string) *Animal {
	return &Animal{
		Name: name,
	}
}

func (a *Animal) GetName() {
	fmt.Printf("Animal name: %v\n", a.Name)
}

//继承
type MachineCat struct {
	//匿名
	*Animal
	Name string
}

func newMachineCat(name string) *MachineCat {
	return &MachineCat{
		Animal: NewAnimal(name),
	}
}

//定义一个receive为MachineCat的函数,重写了
func (value *MachineCat) GetName() {
	fmt.Printf("MachineCat: %v\n", value.Name)
}

type Cat struct {
	*Animal
}

//组合
type Dog struct {
	animal Animal
}

// 多态，Cat, Dog实现了Animal2中的全部方法，就是Animal2的子类
type Animal2 interface {
	Sleep()
	Age() int
	Type() string
}

type Cat2 struct {
	MaxAge int
}

func (this *Cat2) Sleep() {
	fmt.Println("Cat need sleep")
}
func (this *Cat2) Age() int {
	return this.MaxAge
}
func (this *Cat2) Type() string {
	return "Cat"
}

type Dog2 struct {
	MaxAge int
}

func (this *Dog2) Sleep() {
	fmt.Println("Dog need sleep")
}
func (this *Dog2) Age() int {
	return this.MaxAge
}
func (this *Dog2) Type() string {
	return "Dog"
}

func getAnimal2(name string) Animal2 {
	switch name {
	case "dog":
		return &Dog2{MaxAge: 20}
	case "cat":
		return &Cat2{MaxAge: 10}
	default:
		panic("No such animal")
	}
}

func TestAnimal2(t *testing.T) {
	animal := getAnimal2("dog")
	animal.Sleep()
	fmt.Printf("%s max age is: %d\n", animal.Type(), animal.Age())
}

//使用类型别名，通过interface关键字来简化接口声明
//定义了Shape接口，它有两个方法Area和Perimeter，他们不接收任何参数并返回float64。 任何实现这两个方法的类型我们都认为它实现了Shape接口。
type Shape interface {
	Area() float64
	Perimeter() float64
}

//由于interface也是一种类型，我们可以创建这种类型的变量。 在上面的例子中，我们可以创建一个类型为接口Shape的变量s。
//接口有两种类型。
// 静态类型的接口是接口本身，例如程序中的Shape。接口没有静态值，而是指向动态值。
// 接口类型的变量可以保存实现接口的Type的值。该类型的值称为接口的动态值(实际)，该类型称为接口的动态类型(实际)。
// 从结果，我们可以看到接口的值是nil而且接口的类型也是nil。这是因为此时，接口不知道是谁会实现它。 当我们调用Println函数时，它指向接口的动态值，而Printf函数中的％T语法指的是接口的动态类型。但实际上，接口的类型是Shape。
func TestInterfaceBase(t *testing.T) {
	var s Shape
	fmt.Println("value of s is", s)    // 实际值nil
	fmt.Printf("type of s is %T\n", s) // 实际类型nil
}

type Rect struct {
	width  float64
	height float64
}

// 使用Rect接收器类型定义了Area和Perimeter等方法,由于这些方法是由Shape接口定义的，因此Rect实现了Shape接口
func (r Rect) Area() float64 {
	return r.width * r.height
}

func (r Rect) Perimeter() float64 {
	return 2 * (r.width + r.height)
}

func TestInterfaceImpl(t *testing.T) {
	var s Shape
	s = Rect{5.0, 4.0}
	r := Rect{5.0, 4.0}
	// s的动态类型现在是Rect，s的动态值是struct Rect的值，即{5 4}。动态因为，我们可以为不同类型分配新的结构，也实现接口Shape
	fmt.Printf("type of s is %T\n", s)
	fmt.Printf("value of s is %v\n", s)
	// 有时，动态类型的接口也称为具体类型，因为当我们访问接口类型时，它返回它的基础动态值的类型，并且它的静态类型保持隐藏
	//可以在s上调用Area方法，因为s的具体类型是Rect而Rect实现了Area方法
	fmt.Println("area of rectange.s", s.Area())
	// 可以将s与类型为Rect的r struct进行比较，因为它们都具有相同的Rect具体类型且具有相同的值
	fmt.Println("s == r is", s == r)
}

type Circle struct {
	radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.radius
}

// 改变s的动态类型和值
func TestInterfaceChangeType(t *testing.T) {
	var s Shape = Rect{10, 3}
	fmt.Printf("type of s is %T\n", s)
	fmt.Printf("value of s is %v\n", s)
	fmt.Printf("value of s is %0.2f\n", s.Area())

	// 由于新的struct类型Circle也实现了Shape接口，我们可以为s分配一个Circle类型的结构
	// 可以理解为什么接口的(实际)类型和值是动态的
	s = Circle{10}
	fmt.Printf("type of s is %T\n", s)
	fmt.Printf("value of s is %v\n", s)
	fmt.Printf("value of s is %0.2f\n", s.Area())
}

// 为了成功实现接口，您需要实现接口声明的所有方法
// 当接口没有方法时，它被称为空接口。 这由interface{}表示。 由于空接口没有任何方法，因此所有类型都实现了该接口。

// 创建了一个自定义字符串类型MyString
type MyString string

// 接收空接口，可以将MyString和Rect类型的变量传递给它，因为类型为空接口的参数i可以保存任何类型的值，因为所有类型都实现它。
func explain(i interface{}) {
	fmt.Printf("value given to explain function is of type '%T' with value %v\n", i, i)
}

func TestInterfaceNil(t *testing.T) {
	ms := MyString("hello world!")
	r := Rect{5.5, 4.5}
	explain(ms)
	explain(r)
}

// 一个类型可以实现多个接口
type Shape1 interface {
	Area() float64
}

type Object interface {
	Volume() float64
}

// 使用Area方法创建了Shape接口，使用Volume方法创建了Object接口
// 由于struct type Cube实现了这两种方法，因此它实现了这两种接口
type Cube struct {
	side float64
}

func (c Cube) Area() float64 {
	return 6 * (c.side * c.side)
}

func (c Cube) Volume() float64 {
	return c.side * c.side * c.side
}

func TestInterfaceMulti(t *testing.T) {
	c := Cube{3}
	// 可以将struct type Cube的值赋给Shape或Object类型的变量
	// 我们期望s具有c和o的动态值
	var s Shape1 = c
	var o Object = c
	fmt.Println("volume of s of interface type Shape is", s.Area())
	fmt.Println("are of s of interface type Object is", o.Volume())
}

// 使用语法i.(Type)找出接口的基础动态值,i是接口，Type是实现接口i的子类型,Go将检查i的动态类型是否与Type相同
// 类型断言
func TestInterfaceAssert(t *testing.T) {
	var s Shape1 = Cube{3}
	c := s.(Cube) // 现在可以访问变量c中接口s的底层值，它是Cube类型的结构
	// 在类型断言语法i.(Type)中，如果Type没有实现接口（类型）i那么go编译器会抛出错误
	fmt.Println("area of c of type Cube is", c.Area())
	fmt.Println("volume of c of type Cube is", c.Volume())
}

type Skin interface {
	Color() float64
}

// 另一种类型断言语法的变体
//value, ok := i.(Type)
// 可以检查使用ok变量，如果Type实现接口（类型）i。 如果是，那么ok将为true，有具体类型Type。否则为false，value为struct的零值。
// 如何知道接口的底层值是否实现了任何其他接口？ 类型断言也可以这样做。 如果Type in断言语法是interface，那么go将检查i的动态类型是否实现接口Type
func TestInterfaceAssert2(t *testing.T) {
	var s Shape1 = Cube{3}
	value1, ok1 := s.(Object)
	fmt.Printf("dynamic value of Shape 's' with value %v implements interface Object? %v\n", value1, ok1)
	// 由于Cube结构不实现Skin接口。如果我们使用更简单的v := i.(type)语法，那么我们的程序会报错
	value2, ok2 := s.(Skin)
	fmt.Printf("dynamic value of Shape 's' with value %v implements interface Skin? %v\n", value2, ok2)
}

// 类型开关,类型切换
func explain2(i interface{}) {
	switch i.(type) { //  在switch中使用i.(type)语句，我们可以访问该动态类型。
	case string:
		fmt.Println("i stored string", strings.ToUpper(i.(string))) // 类型断言
	case int:
		fmt.Println("I stored int", i)
	default:
		fmt.Println("i stored something else", i)
	}
}

func TestInterfaceTypeSwitch(t *testing.T) {
	explain2("hello world")
	explain2(52)
	explain2(true)
}

// 嵌入式接口
// 在go中，接口不能实现其他接口或扩展它们，但我们可以通过合并两个或多个接口来创建新接口
type Material interface {
	Shape1
	Object
}

// %T返回实际类型
func TestInterfaceNested(t *testing.T) {
	// 由于Cube实现了方法Area和Volume，它实现了Shape和Object接口。 但由于接口Material是这些接口的嵌入式接口，Cube也必定实现它
	// 发生这种情况是因为像匿名嵌套结构一样，嵌套接口的所有方法都被提升为父接口
	c := Cube{3}
	var m Material = c
	var s Shape1 = c
	var o Object = c
	fmt.Printf("dynamic type and value of interface m of static type Material is '%T' and '%v'\n", m, m)
	fmt.Printf("dynamic type and value of interface m of static type Shape is '%T' and '%v'\n", s, s)
	fmt.Printf("dynamic type and value of interface m of static type Object is '%T' and '%v'\n", o, o)
}

// 指针与值接收器
// 到目前为止，我们已经看到了带有值接收器
//对于接受指针接收器的方法，接口是否正常
//指针接收器的方法将对指针或值都起作用
//在接口的情况下，如果方法有指针接收器，那么接口将具有动态类型的指针而不是动态类型的值。
// 当我们为接口变量分配类型值时，我们需要分配类型为value的指针
type Rect2 struct {
	width  float64
	height float64
}

//func (r Rect2) Area() float64 {
//	return r.width * r.height
//}
// 指针实现方法
func (r *Rect2) Area() float64 {
	return r.width * r.height
}

func (r Rect2) Perimeter() float64 {
	return 2 * (r.width + r.height)
}

func TestInterfacePointer(t *testing.T) {
	r := Rect2{5.0, 4.0}
	var s Shape = &r // 这里用指针,s的具体值现在是一个指针
	area := s.Area()
	perimeter := s.Perimeter()
	fmt.Println("area of rectangle is", area)
	fmt.Println("permeter of rectangle is", perimeter)
}
