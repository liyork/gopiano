package base

import (
	"fmt"
	"testing"
)

type IPeople interface {
	GetName() string
	SetName(string)
	GetAge() int
	SetAge(int)
	Run()
}

// 利用接口和组合继承，可以实现抽象类
type AbstractPeople struct {
	//IPeople
	name string
	age  int
}

func (a AbstractPeople) GetName() string {
	return a.name
}

func (a *AbstractPeople) SetName(newName string) {
	a.name = newName
}

func (a AbstractPeople) GetAge() int {
	fmt.Println("before GetAge")
	out := a.age
	fmt.Println("after GetAge")
	a.GetNameAfter()
	return out
}

// 抽象方法
func (a *AbstractPeople) GetNameAfter() {
	fmt.Println("AbstractPeople GetNameAfter")
}

func (a *AbstractPeople) SetAge(newAge int) {
	a.age = newAge
}

func (a *AbstractPeople) Run() {
}

type ConcretePeople1 struct {
	AbstractPeople
	age1 int
}

// 重写1
//func (a *ConcretePeople1) GetAge() int {
//	return a.age1
//}

// 实现抽象方法
func (a *ConcretePeople1) GetNameAfter() {
	fmt.Println("ConcretePeople1 GetNameAfter")
}

func TestBaseAbstract(t *testing.T) {
	people1 := &ConcretePeople1{}
	people1.age = 1
	people1.name = "xx1"
	people1.age1 = 1111
	// 由于ConcretePeople1把GetAge给重写了，就不会调用AbstractPeople的GetAgele ....
	fmt.Println("people1.getAge, people1.getName", people1.GetAge(), people1.GetName())
	// 但是注释掉[重写1]部分，ConcretePeople1.GetNameAfter又不能执行。。。并不能动态绑定啊。。

	var p IPeople = people1
	fmt.Println("p.getAge, p.getName", p.GetAge(), p.GetName())
	p1 := p.(*ConcretePeople1)
	fmt.Println("p1.getAge, p1.getName", p1.GetAge(), p1.GetName())
}
