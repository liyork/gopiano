package base

import (
	"fmt"
	"log"
	"os"
	"testing"
)

//====== 非指针方式组合
type Base struct {
}

func (b *Base) Func1() {
	fmt.Println("Base.Func1() was invoked!")
}

func (b *Base) Func2() {
	fmt.Println("Base.Func2() was invoked!")
}

type Derived struct {
	Base
}

func (d *Derived) Func2() {
	fmt.Println("Derived.Func2() was invoked!")
}

func (d *Derived) Func3() {
	fmt.Println("Derived.Func3() was invoked!")
}

//  继承规则
// 在派生类没有改写基类的成员方法时，相应的成员方法被派生类继承。
// 派生类可以直接调用基类的成员方法，譬如基类有个成员方法为Base.Func()，那么Derived.Func()等同于Derived.Base.Func()
// 倘若派生类的成员方法名与基类的成员方法名相同，那么基类方法将被覆盖或叫隐藏，譬如基类和派生类都有成员方法Func()，那么Derived.Func()将只能调用派生类的Func()方法，如果要调用基类版本，可以通过Derived.Base.Func()来调用。
func TestBaseRefInherit(t *testing.T) {
	d := &Derived{}
	// 都用的base
	d.Func1()
	d.Base.Func1()

	// 重写了
	d.Func2()
	// base
	d.Base.Func2()

	// 仅子类有
	d.Func3()
}

// 内存布局
// golang很清晰地展示类的内存布局是怎样的，Base的位置即基类成员展开的位置。还可以修改内存布局，Base的位置可以出现在派生类的任何位置
//type Derived struct {
//	DerivedName string
//	Base// 这里
//}

// 名字覆盖
type Base1 struct {
	Name string
}

type Derived1 struct {
	Base1
	Name string
}

func TestNameOverwrite(t *testing.T) {
	d := &Derived1{}
	// 子类
	d.Name = "Derived"
	d.Base1.Name = "Base"

	// 优先用子类的覆盖名称
	fmt.Println(d.Name)
	fmt.Println(d.Base1.Name)
}

//====== 指针方式组合
type MyJob struct {
	Command string
	// 指针
	*log.Logger
}

func (job *MyJob) Start() {
	job.Println("job started!") // job.Logger.Println
	fmt.Println(job.Command)
	job.Println("job finished!") // job.Logger.Println
}

// 继承规则
// 派生类对基类采用指针方式的组合，依然具有派生的效果，只是派生类创建实例的时候需要外部提供一个基类实例的指针
// 其他规则与非指针方式组合一致
func TestBasePointRefInherit(t *testing.T) {
	logFile, err := os.OpenFile("./job.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer logFile.Close()

	logger := log.New(logFile, "[info]", log.Ldate|log.Ltime|log.Llongfile)
	// 匿名组合相当于以其类型名称（去掉包名部分）作为成员变量的名字。
	job := MyJob{
		Command: "programming",
		Logger:  logger,
	}

	// 在经过合适的赋值后，MyJob类型的所有成员方法可以很方便地借用所有log.Logger提供的方法。匿名组合
	job.Start()
	job.Println("test finished!") // job.Logger.Println
}

//接口继承
type MyInterface1 interface {
	Func1() bool
}

type MyInterface2 interface {
	//Func1() bool // duplicate method Func1
	Func2() bool
}

type MyInterface3 interface {
	MyInterface1
	MyInterface2

	//Func1() bool // duplicate method Func1
	Func3() bool
}

type MyClass2 struct {
}

func (p *MyClass2) Func1() bool {
	fmt.Println("MyClass.Func1()")
	return true
}

func (p *MyClass2) Func2() bool {
	fmt.Println("MyClass.Func2()")
	return true
}

func (p *MyClass2) Func3() bool {
	fmt.Println("MyClass.Func3()")
	return true
}

func TestInterfaceInherit(t *testing.T) {
	var myInterface3 MyInterface3 = new(MyClass)
	myInterface3.Func1() // MyClass.Func1()
	myInterface3.Func2() // MyClass.Func2()
	myInterface3.Func3() // MyClass.Func3()
}
