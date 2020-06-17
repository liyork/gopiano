package base

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
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
func (p Teacher) Getname() string { //实现Getname方法
	fmt.Println(p.Name)
	return p.Name
}

func (p Teacher) Running() { // 实现 Running方法
	fmt.Printf("Teacher %s running\n", p.Name)
}

func TestIntefaceBase(t *testing.T) {
	var stu1 Student
	stu1.Name = "wd"
	stu1.Age = 22

	var t1 Teacher
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
