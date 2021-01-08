package test

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestBase(t *testing.T) {
	t1 := reflect.TypeOf(3)
	fmt.Println(t1.String())
	fmt.Println(t1)

	var w io.Writer = os.Stdout
	fmt.Println(reflect.TypeOf(w))

	fmt.Printf("%T\n", 3)

	v := reflect.ValueOf(3)
	fmt.Println(v)
	fmt.Printf("%T\n", v)
	fmt.Println(v.String())

	// 还原
	x := v.Interface()
	i := x.(int)
	fmt.Printf("%d\n", i)

	fmt.Println("")
	x = 2
	a := reflect.ValueOf(2)
	b := reflect.ValueOf(x)
	c := reflect.ValueOf(&x)
	d := c.Elem()
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(d)
	fmt.Println(a.CanAddr())
	fmt.Println(b.CanAddr())
	fmt.Println(c.CanAddr())
	fmt.Println(d.CanAddr())

	fmt.Println("")
	x1 := 2
	d1 := reflect.ValueOf(&x1).Elem()
	px := d1.Addr().Interface().(*int)
	*px = 3
	fmt.Println(x1)

	d1.Set(reflect.ValueOf(4))
	fmt.Println(x1)

	fmt.Println("")
	// panic: reflect: reflect.flag.mustBeAssignable using unaddressable value
	x3 := 2
	//b3 := reflect.ValueOf(x3)
	//b3.Set(reflect.ValueOf(3))

	b3 := reflect.ValueOf(&x3).Elem()
	b3.SetInt(3)
	fmt.Println(b3.CanAddr(), b3.CanSet())
	fmt.Println(x3)

	fmt.Println("")
	// 有类型的变量，不能设定其他类型
	x4 := 1
	rx := reflect.ValueOf(&x4).Elem()
	rx.SetInt(2)
	rx.Set(reflect.ValueOf(3))
	//rx.SetString("hello")            // panic
	//rx.Set(reflect.ValueOf("hello")) //panic

	// interface{}值需要用valueof设定,类型都行
	var y interface{}
	ry := reflect.ValueOf(&y).Elem()
	//ry.SetInt(2) //panic
	ry.Set(reflect.ValueOf(3))
	//ry.SetString("hello")//panic
	ry.Set(reflect.ValueOf("hello"))
	fmt.Println(ry)

	fmt.Println("")
	stdout := reflect.ValueOf(os.Stdout).Elem()
	fmt.Println(stdout.Type())
	fd := stdout.FieldByName("fd")
	fmt.Println(fd)
	//fmt.Println(fd.Int())
	//fd.SetInt(2) //panic
}
