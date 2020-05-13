package base

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestStringnil(t *testing.T) {

	// string不能设为nil
	//var x string = nil //error,cannot use nil as type string in assignment

	// 不能用string和nil比较
	//if x == nil { //error,invalid operation: x == nil (mismatched types string and nil)
	//}

	// 默认值为""而不是nil
	var x string //defaults to "" (zero value)

	if x == "" {
		x = "default"
	}

	fmt.Println("x", x)
	time.Sleep(2 * time.Second)
}

// golang 里面的字符串都是不可变的，每次运算都会产生一个新的字符串
func BenchmarkAddStringWithOperator(b *testing.B) {
	hello := "hello"
	world := "world"
	for i := 0; i < b.N; i++ {
		_ = hello + "," + world
	}
}

// 内部使用 []byte 实现，不像直接运算符这种会产生很多临时的字符串，但是内部的逻辑比较复杂，有很多额外的判断，
// 还用到了 interface，所以性能也不是很好
// 优点是可以拼接其他类型
func BenchmarkAddStringWithSprintf(b *testing.B) {
	hello := "hello"
	world := "world"
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%s,%s", hello, world)
	}
}

// join会先根据字符串数组的内容，计算出一个拼接之后的长度，然后申请对应大小的内存，一个一个字符串填入，
// 在已有一个数组的情况下，这种效率会很高，但是本来没有，去构造这个数据的代价也不小
func BenchmarkAddStringWithJoin(b *testing.B) {
	hello := "hello"
	world := "world"
	for i := 0; i < b.N; i++ {
		_ = strings.Join([]string{hello, world}, ",")
	}
}

// 这个比较理想，可以当成可变字符使用，对内存的增长也有优化，如果能预估字符串的长度，还可以用 buffer.Grow() 接口来设置 capacity
func BenchmarkAddStringWithBuffer(b *testing.B) {
	hello := "hello"
	world := "world"
	for i := 0; i < 1000; i++ {
		var buffer bytes.Buffer
		buffer.WriteString(hello)
		buffer.WriteString(",")
		buffer.WriteString(world)
		_ = buffer.String()
	}
}

func TestSubString(t *testing.T) {
	tracer := "aye bye"
	comma := strings.Index(tracer, "b") // index从0开始
	fmt.Println("comma:", comma)
	fmt.Println("tracer[comma:]", tracer[comma:]) // startIndex包含
	pos := strings.Index(tracer[comma:], "a")
	fmt.Println("pos:", pos)
}

func TestReplace(t *testing.T) {
	s := "abcabcabcxxxx"
	// 替换所有
	fmt.Println("n=-1: ", strings.Replace(s, "a", "d", -1))
	// 替换2个
	fmt.Println("n=-1: ", strings.Replace(s, "a", "d", 2))
}
