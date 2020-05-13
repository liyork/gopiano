package base

import (
	"fmt"
	"testing"
)

// interface(接口),可以定义一组方法，但是这些不需要实现。并且interface不能包含任何变量。
// interface是方法的集合
//interface是一种类型，并且是指针类型
//interface的更重要的作用在于多态实现(一种事物的多种形态,同一个interface，不同的类型实现，都可以进行调用，它们都按照统一接口进行操作)
//接口的使用不仅仅针对结构体，自定义类型、变量等等都可以实现接口。
//如果一个接口没有任何方法，我们称为空接口，由于空接口没有方法，所以任何类型都实现了空接口。
//要实现一个接口，必须实现该接口里面的所有方法。

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
