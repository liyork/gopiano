package base

import (
	"fmt"
	"reflect"
	"testing"
)

type astruct struct {
}
type teststruct struct {
	b   int
	c   string
	as1 *astruct
	as2 astruct
}

func TestInitialize(t *testing.T) {
	var a int //局部变量只要声明也会被初始化默认值
	fmt.Println("a:", a)

	//声明也会初始化，tt类型是teststruct
	var tt teststruct
	fmt.Println("tt.b:", tt.b)
	fmt.Println("tt.c:", tt.c)
	fmt.Println("tt.as1:", tt.as1) //指针初始值为nil
	fmt.Println("tt.as2:", tt.as2)
	fmt.Printf("tt type:%v\n", reflect.TypeOf(tt))

	//new返回指向已分配内存的指针,newas是astruct类型指针
	newas := new(astruct)
	fmt.Println("newas:", newas)
	fmt.Printf("newas type:%v\n", reflect.TypeOf(newas))

	// new的简写
	ts := &teststruct{}
	fmt.Println("ts.b:", ts.b)
	fmt.Println("ts.c:", ts.c)
	fmt.Println("ts.as1:", ts.as1) //指针初始值为nil
	fmt.Println("ts.as2:", ts.as2)
}

// Go 语言中，结构体和它所包含的数据在内存中是以连续块的形式存在的，即使结构体中嵌套有其他的结构体，这在性能上带来了很大的优势。
func TestInitDiff(t *testing.T) {
	// 使用var声明
	var p1 Person
	p1.name = "zhangsan"
	p1.age = 18
	fmt.Printf("This is %s, %d years old\n", p1.name, p1.age)
	fmt.Printf("p1 type:%v\n", reflect.TypeOf(p1)) //p1就是指向数据Person的变量

	// 使用new构造,返回的是指针
	p2 := new(Person)
	p2.name = "lisi"
	p2.age = 20
	// 虽然 p2 是指针类型，但我们仍然可以像 p2.age = 23 。Go 会自动进行转换
	//(*p2).age = 23 // 也可以先通过*操作符来获取指针所指向的内容，再进行赋值
	fmt.Printf("This is %s, %d years old\n", p2.name, p2.age)
	fmt.Printf("p2 type:%v\n", reflect.TypeOf(p2)) //p2是一个指针，里面内容指向数据Person

	// 直接构造
	p3 := Person{"wangwu", "man", 25}
	fmt.Printf("This is %s, %d years old\n", p3.name, p3.age)
	fmt.Printf("p3 type:%v\n", reflect.TypeOf(p3))

	p4 := &Person{"wangwu", "man", 25} //与p3的区别
	fmt.Printf("This is %s, %d years old\n", p4.name, p4.age)
	fmt.Printf("p4 type:%v\n", reflect.TypeOf(p4))
}
