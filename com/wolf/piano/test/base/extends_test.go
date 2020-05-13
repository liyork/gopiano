package base

import (
	"fmt"
	"testing"
)

// 接口
type IAnimal interface {
	GetName()
}

type Animal struct {
	Name string
}

// 继承
func TestInherit(t *testing.T) {
	//实例化machineCat并取指针
	machineCat := &MachineCat{
		&Animal{
			"test",
		},
	}

	machineCat.GetName()
	// 可以直接访问匿名继承来的属性
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
	cat.Animal = &Animal{}
	cat.Name = "bb"
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
}

func newMachineCat(name string) *MachineCat {
	return &MachineCat{
		NewAnimal(name),
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
