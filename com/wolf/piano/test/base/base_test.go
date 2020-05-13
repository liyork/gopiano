package base

import (
	"fmt"
	"testing"
)

// 指针储存的是一个值的地址，指针本身也需要地址来储存
func Test_point(t *testing.T) {
	var p *int   // 指针声明
	p = new(int) // 指针赋值
	*p = 1       // *表示取指针存储地址的值，赋值
	fmt.Println(p, &p, *p)
	// 0xc0000901a8,p指针存储的地址
	// 0xc000088018,&代表取指针本身的地址
	// *得到指针存储地址的值，1
	// 总结：对于指针p,直接打印p得到的是p存储的地址，&p是p指针本身占用内存地址，*p是指针存储地址的值
	//地址	0xc000088018	0xc0000901a8
	//值	0xc0000901a8	1
}

func Test_pointErr(t *testing.T) {
	var i *int
	// i = new(int) --解决问题
	*i = 1
	fmt.Println(i, &i, *i)

	// panic: runtime error: invalid memory address or nil pointer dereference [recovered]
	//	panic: runtime error: invalid memory address or nil pointer dereference
	// [signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x10f899c]
	// 指针i，占用内存，存储的值为nil，*i为取指针i的存储值对应内存的值，为nil并没有给其分配内存，然后赋值1报错
}
