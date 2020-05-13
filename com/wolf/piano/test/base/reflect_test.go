package base

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"testing"
)

// Go语言提供了一种机制在运行时更新和检查变量的值、调用变量的方法和变量支持的内在操作，但是在编译时并不知道这些变量的具体类型，这种机制被称为反射。
// go语言提供了一种机制，在编译时不知道类型的情况下，可更新变量，在运行时查看值，调用方法以及直接对他们的布局进行操作。这种机制称为反射(reflection)
// 在计算机科学领域，反射是指一类应用，它们能够自描述和自控制。也就是说，这类应用通过采用某种机制来实现对自己行为的描述（self-representation）和监测（examination），并能根据自身行为的状态和结果，调整或修改应用所描述行为的状态和相关的语义。

// Golang关于类型设计的一些原则
//
//变量包括（type, value）两部分
//理解这一点就知道为什么nil != nil了
//type 包括 static type和concrete type. 简单来说 static type是你在编码是看见的类型(如int、string)，concrete type是runtime系统看见的类型
//类型断言能否成功，取决于变量的concrete type，而不是static type. 因此，一个 reader变量如果它的concrete type也实现了write方法的话，它也可以被类型断言为writer.
//
// 反射，就是建立在类型之上的，Golang的指定类型的变量的类型是静态的（也就是指定int、string这些的变量，它的type是static type），
// 在创建变量的时候就已经确定，反射主要与Golang的interface类型相关（它的type是concrete type），只有interface类型才有反射一说。
// 在Golang的实现中，每个interface变量都有一个对应pair，pair中记录了实际变量的值和类型:
//(value, type),value是实际变量值，type是实际变量的类型
// 一个interface{}类型的变量包含了2个指针，一个指针指向值的类型【对应concrete type】，另外一个指针指向实际的值【对应value】
// 例如，创建类型为*os.File的变量，然后将其赋给一个接口变量r：
//tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
//var r io.Reader
//r = tty
//接口变量r的pair中将记录如下信息：(tty, *os.File)，这个pair在接口变量的连续赋值过程中是不变的，将接口变量r赋给另一个接口变量w:
//var w io.Writer
//w = r.(io.Writer)
//接口变量w的pair与r的pair相同，都是:(tty, *os.File)，即使w是空接口类型，pair也是不变的。
//interface及其pair的存在，是Golang中实现反射的前提，理解了pair，就更容易理解反射。反射就是用来检测存储在接口变量内部(值value；类型concrete type) pair对的一种机制。
// reflect.TypeOf()是获取pair中的type，reflect.ValueOf()获取pair中的value

// Golang reflect慢主要有两个原因
//涉及到内存分配以及后续的GC；
//reflect实现里面有大量的枚举，也就是for循环，比如类型之类的。

// Type：Type类型用来表示一个go类型。
// Value为go值提供了反射接口

type User struct {
	ID   int
	Name string
}

// 定义一个Enum类型
type Enum int

const (
	Zero Enum = 0
)

type cat struct {
}

// 使用 reflect.TypeOf() 函数可以获得任意值的类型对象（reflect.Type）,通过类型对象可以访问任意值的类型信息
// Go 程序中的类型（Type）指的是系统原生数据类型，如 int、string、bool、float32 等类型，以及使用 type 关键字定义的类型，这些类型的名称就是其类型本身的名称
// TypeOf returns the reflection Type that represents the dynamic type of i.
// If i is a nil interface value, TypeOf returns nil.
// func TypeOf(i interface{}) Type {...}
func TestReflectBase(t *testing.T) {

	typeOfCat := reflect.TypeOf(cat{})
	// 显示反射类型对象的名称和种类
	fmt.Println(typeOfCat.String(), typeOfCat.Name(), typeOfCat.Kind())
	// 获取Zero常量的反射类型对象
	typeOfA := reflect.TypeOf(Zero)
	// 显示反射类型对象的名称和种类
	fmt.Println(typeOfA.Name(), typeOfA.Kind())

	inf := new(Skills)
	// Array, Chan, Map, Ptr, or Slice.
	inf_type := reflect.TypeOf(inf).Elem() // 引用类型需要用Elem()获取指针所指的对象类型

	stu1 := Student{Name: "wd", Age: 22}
	stu_type := reflect.TypeOf(stu1)
	fmt.Println(stu_type.String())                           //main.Student
	fmt.Println(stu_type.Name())                             //Student
	fmt.Println(stu_type.PkgPath())                          //main
	fmt.Println(stu_type.Kind())                             //struct
	fmt.Println(stu_type.Size())                             //24
	fmt.Println(inf_type.NumMethod())                        //2
	fmt.Println(inf_type.Method(0), inf_type.Method(0).Name) // {reading main func() <invalid Value> 0} reading
	fmt.Println(inf_type.MethodByName("reading"))            //{reading main func() <invalid Value> 0} true

}

// 对指针获取反射对象时，可以通过 reflect.Elem() 方法获取这个指针指向的元素类型，这个获取过程被称为取元素，等效于对指针类型变量做了一个*操作
func TestReflectPoint(t *testing.T) {
	// 创建cat的实例并返回指针
	ins := &cat{}
	// 获取结构体实例的反射类型对象
	typeOfCat := reflect.TypeOf(ins)
	// 显示反射类型对象的名称和种类(指针变量的类型名称和种类),go对所有指针变量的种类都是 Ptr
	fmt.Printf("name:'%v' kind:'%v'\n", typeOfCat.Name(), typeOfCat.Kind())
	// 注意对于：数组、切片、映射、通道、指针、接口 ,取指针类型的元素类型
	typeOfCat = typeOfCat.Elem()
	// 显示反射类型对象的名称和种类(指针变量指向元素的类型名称和种类)
	fmt.Printf("element name: '%v', element kind: '%v'\n", typeOfCat.Name(), typeOfCat.Kind())

	// 通过elem，但必须传递的是指针，修改值
	var a int64 = 100
	reflectSetValue2(&a)
	fmt.Println("a:", a)
}

func TestReflectSetValue(t *testing.T) {
	var num float64 = 1.2345
	fmt.Println("old value of pointer:", num)

	// 通过reflect.ValueOf获取num中的reflect.Value，注意，参数必须是指针才能修改其值
	pointer := reflect.ValueOf(&num)
	newValue := pointer.Elem()

	fmt.Println("type of pointer:", newValue.Type())
	fmt.Println("settability of pointer:", newValue.CanSet())

	// 重新赋值
	newValue.SetFloat(77)
	fmt.Println("new value of pointer:", num)

	// 如果reflect.ValueOf的参数不是指针，会如何？
	pointer = reflect.ValueOf(num)
	fmt.Println("settability of pointer2:", pointer.CanSet())
	//newValue = pointer.Elem() // 如果非指针，这里直接panic
}

func reflectSetValue2(x interface{}) {
	v := reflect.ValueOf(x)
	// 反射中使用 Elem()方法获取指针对应的值
	if v.Elem().Kind() == reflect.Int64 {
		v.Elem().SetInt(200)
	}
}

func TestReflectStruct(t *testing.T) {
	// 声明一个空结构体
	type cat struct {
		Name string
		// 以 ` 开始和结尾的字符串。这个字符串在Go语言中被称为 Tag（标签）。一般用于给字段添加自定义信息，方便其他模块根据信息进行不同功能的处理。
		Type int `json:"type" id:"100"`
	}
	// 创建cat的实例
	ins := cat{Name: "mimi", Type: 1}
	// 获取结构体实例的反射类型对象
	typeOfCat := reflect.TypeOf(ins)
	// 获得一个结构体类型共有多少个字段。如果类型不是结构体，将会触发宕机错误。
	for i := 0; i < typeOfCat.NumField(); i++ {
		// 获取每个成员的结构体字段类型,返回 StructField 结构，这个结构描述结构体的成员信息，通过这个信息可以获取成员与结构体的关系，如偏移、索引、是否为匿名字段、结构体标签（StructTag）等
		fieldType := typeOfCat.Field(i)
		// 输出成员名和tag
		fmt.Printf("name: %v  tag: '%v'\n", fieldType.Name, fieldType.Tag)
	}
	// 通过字段名, 找到字段类型信息
	if catType, ok := typeOfCat.FieldByName("Type"); ok {
		// 从tag中取出需要的tag,StructTag
		fmt.Println(catType.Tag.Get("json"), catType.Tag.Get("id"))
	}
}

func TestReflectType(t *testing.T) {
	type myInt int64

	var a *float32 // 指针
	var b myInt    // 自定义类型
	var c rune     // 类型别名//代表int32
	reflectType(a) // type: kind:ptr
	reflectType(b) // type:myInt kind:int64
	reflectType(c) // type:int32 kind:int32

	type person struct {
		name string
		age  int
	}
	var d = person{
		name: "wang",
		age:  18,
	}

	reflectType(d) // type:person kind:struct
}

func reflectType(x interface{}) {
	t := reflect.TypeOf(x)
	fmt.Printf("type:%v kind:%v\n", t.Name(), t.Kind())
}

// ValueOf returns a new Value initialized to the concrete value
// stored in the interface i.  ValueOf(nil) returns the zero
// func ValueOf(i interface{}) Value {...}
// reflect.ValueOf的返回值也是具体值，不过reflect.Value也可以包含一个接口值
func TestReflectValueBase(t *testing.T) {
	var num float64 = 1.2345
	//得到了一个类型为”relfect.Value”变量
	pointer := reflect.ValueOf(&num)
	value := reflect.ValueOf(num)

	fmt.Println("value: ", value)
	// 获得接口变量的真实内容
	value.Interface()
	// 注意的时候，转换的时候，如果转换的类型不完全符合，则直接panic
	convertPointer := pointer.Interface().(*float64)
	convertValue := value.Interface().(float64)

	fmt.Println(convertPointer)
	fmt.Println(convertValue)

}

func TestReflectGetAll(t *testing.T) {
	user := User1{1, "Allen.Wu", 25}
	DoFiledAndMethod(user)
}

type User1 struct {
	Id   int
	Name string
	Age  int
}

func (u User1) ReflectCallFunc() {
	fmt.Println("Allen.Wu ReflectCallFunc")
}

// 反射字段需要对外
// 通过接口来获取任意参数，然后一一揭晓
func DoFiledAndMethod(input interface{}) {

	getType := reflect.TypeOf(input)
	fmt.Println("get Type is :", getType.Name())

	getValue := reflect.ValueOf(input)
	fmt.Println("get all Fields is:", getValue)

	// 获取方法字段
	// 1. 先获取interface的reflect.Type，然后通过NumField进行遍历
	// 2. 再通过reflect.Type的Field获取其Field
	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i)
		// 取值
		value := getValue.Field(i).Interface()
		fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)
	}

	// 获取方法
	// 1. 先获取interface的reflect.Type，然后通过.NumMethod进行遍历
	for i := 0; i < getType.NumMethod(); i++ {
		m := getType.Method(i)
		fmt.Printf("%s: %v\n", m.Name, m.Type)
	}
	// 输出：ReflectCallFuncHasArgs: func(base.User1, string, int)
	// 结构体的方法其实也是通过函数实现,把base.User1当参数了
}

func TestReflectValue(t *testing.T) {
	var a float32 = 3.14
	var b int64 = 100
	reflectValue(a) // type is float32, value is 3.140000
	reflectValue(b) // type is int64, value is 100
	// 将int类型的原始值转换为reflect.Value类型
	c := reflect.ValueOf(10)
	fmt.Printf("type c :%T\n", c) // type c :reflect.Value
}

func reflectValue(x interface{}) {
	v := reflect.ValueOf(x)
	k := v.Kind()
	switch k {
	case reflect.Int64:
		// v.Int()从反射中获取整型的原始值，然后通过int64()强制类型转换
		fmt.Printf("type is int64, value is %d\n", int64(v.Int()))
	case reflect.Float32:
		// v.Float()从反射中获取浮点型的原始值，然后通过float32()强制类型转换
		fmt.Printf("type is float32, value is %f\n", float32(v.Float()))
	case reflect.Float64:
		// v.Float()从反射中获取浮点型的原始值，然后通过float64()强制类型转换
		fmt.Printf("type is float64, value is %f\n", float64(v.Float()))
	}
}

// 反射可以将“接口类型变量”转换为“反射类型对象”，反射类型指的是reflect.Type和reflect.Value这两种
func TestReflectTypeAndValue(t *testing.T) {
	var num float64 = 1.2345

	fmt.Println("type: ", reflect.TypeOf(num))
	fmt.Println("value: ", reflect.ValueOf(num))
}

func TestReflectMethodInfo(t *testing.T) {
	of := reflect.TypeOf(TestIsNull)
	in := of.In(0)
	fmt.Println(in)
	if in.Kind() == reflect.Ptr {
		elem := in.Elem()
		fmt.Println(elem.Name(), elem.Kind())
	}
}

func TestReflectOptValue(t *testing.T) {
	var a = 2
	vof := reflect.ValueOf(a)
	i := vof.Int()
	fmt.Println("i:", i)

	tof := reflect.TypeOf(a)
	i2 := vof.Type()
	// 一样
	fmt.Println("type:", i2, tof)

	// panic: reflect: reflect.flag.mustBeAssignable using unaddressable value [recovered]
	// 报错，因为我们的a是一个值类型，而值类型的传递是拷贝了一个副本，当 vof := reflect.ValueOf(a) 函数通过传递一个 a 拷贝创建了 vof，
	// 那么 vof 的改变并不能更改原始的 a。要想 vof 的更改能作用到 a，那就必须传递 a 的地址 v = reflect.ValueOf(&a)
	// vof.SetInt(333)
	vof2 := reflect.ValueOf(&a)
	// 通过 Elem() 方法进行取地址
	vof2.Elem().SetInt(333)
	fmt.Println("i:", a)
}

func TestBaseReflect(t *testing.T) {
	user := User{}
	a := reflect.TypeOf(user)
	if _, ok := a.FieldByName("Name"); ok {
		println("存在")
	} else {
		println("不存在")
	}
	fmt.Println(a.Name(), a.Kind(), a.PkgPath())
	field, b := a.FieldByName("ID")
	fmt.Println(field, b)

	user2 := &User{}
	var user3 interface{}
	user3 = user2
	i := reflect.TypeOf(user3)
	elem := i.Elem()
	fmt.Println("elem:", elem.Name())

	typeOf := reflect.TypeOf(test)
	fmt.Println(typeOf)

	name := runtime.FuncForPC(reflect.ValueOf(test).Pointer()).Name()
	fmt.Println("name:", name)
}

type Handler func(obj interface{}) error

func test(qqqq interface{}) error {
	fmt.Println("1111")
	return nil
}

func TestIsNull(t *testing.T) {

	// 声明一个 *int 类型的指针，初始值为 nil。
	var a *int
	fmt.Println("var a *int:", reflect.ValueOf(a).IsNil())
	// nil值
	fmt.Println("nil:", reflect.ValueOf(nil).IsValid())
	// (*int)(nil) 的含义是将 nil 转换为 *int，也就是*int 类型的空指针
	fmt.Println("(*int)(nil):", reflect.ValueOf((*int)(nil)).Elem().IsValid())
	// 实例化一个结构体
	s := struct{}{}
	// 尝试从结构体中查找一个不存在的字段
	fmt.Println("不存在的结构体成员:", reflect.ValueOf(s).FieldByName("").IsValid())
	// 尝试从结构体中查找一个不存在的方法
	fmt.Println("不存在的结构体方法:", reflect.ValueOf(s).MethodByName("").IsValid())
	// 实例化一个map,与 make 方式创建的 map 等效
	m := map[int]int{}
	// 尝试从map中查找一个不存在的键
	fmt.Println("不存在的键：", reflect.ValueOf(m).MapIndex(reflect.ValueOf(3)).IsValid())
}

func TestReflectInvokeFun(t *testing.T) {
	of := reflect.ValueOf(testFunc)
	of.Call(getValues())
}

// 注意，又不符合则panic
// 如何通过反射来进行方法的调用？
// 本来可以用u.ReflectCallFuncXXX直接调用的，但是如果要通过反射，那么首先要将方法注册，也就是MethodByName，然后通过反射调动mv.Call
func TestReflectInvokeMethod(t *testing.T) {
	user := User1{1, "Allen.Wu", 25}

	// 带有参数的调用方法
	args := getValues("wudebao", 30)
	invokeMethod(user, "ReflectCallFuncHasArgs", args)

	// 无参数的调用方法
	args = getValues()
	invokeMethod(user, "ReflectCallFuncNoArgs", args)

	// 调用指针类型方法
	invokeMethod(&user, "ReflectCallFuncNoArgs2", args)
}

func testFunc() {
	fmt.Println("testFunc..")
}

//根据参数获取对应的Values
func getValues(param ...interface{}) []reflect.Value {
	vals := make([]reflect.Value, 0, len(param))
	for i := range param {
		vals = append(vals, reflect.ValueOf(param[i]))
	}
	return vals
}

func invokeMethod(obj interface{}, funcInter string, paramsValue []reflect.Value) {
	getValue := reflect.ValueOf(obj)
	method := getValue.MethodByName(funcInter)

	if method.Kind() != reflect.Func {
		log.Fatal("funcInter is not func")
		return
	}

	if method.Type().NumIn() > 0 {
		in := method.Type().In(0)
		fmt.Println("in:", in)
	}

	if method.Type().NumOut() > 0 {
		out := method.Type().Out(0)
		fmt.Println("out:", out)
	}

	values := method.Call(paramsValue) //方法调用并返回值
	for i := range values {
		fmt.Println(values[i])
	}
}

func (u User1) ReflectCallFuncHasArgs(name string, age int) {
	fmt.Println("ReflectCallFuncHasArgs name: ", name, ", age:", age, "and origal User.Name:", u.Name)
}

func (u User1) ReflectCallFuncNoArgs() {
	fmt.Println("ReflectCallFuncNoArgs")
}

func (u *User1) ReflectCallFuncNoArgs2() {
	fmt.Println("ReflectCallFuncNoArgs2")
}

// 空接口相当于一个容器，能接受任何东西。.ValueOf(
// reflect.Value.Elem(),返回一个 interface 或者 pointer 的值
// Elem returns the value that the interface v contains or that the pointer v points to. It panics if v's Kind is not Interface or Ptr. It returns the zero Value if v is nil.
// reflect.Type.Elem(),返回一个类型（如：Array，Map，Chan等）的元素的类型
// Elem returns a type's element type. It panics if the type's Kind is not Array, Chan, Map, Ptr, or Slice.
func TestReflectNew(t *testing.T) {
	var a string
	of := reflect.TypeOf(a)
	// New returns a Value representing a pointer to a new zero value for the specified type. That is, the returned Value's Type is PtrTo(typ).
	sptr := reflect.New(of)
	fmt.Println("sptr:", sptr)
	// 返回值类型：reflect.Value
	sval := sptr.Elem()
	ss := sval.Interface().(string)
	fmt.Println("ss:", ss) // 空字符串
}
