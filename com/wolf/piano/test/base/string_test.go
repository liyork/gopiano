package base

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

//一个字符串是一个不可改变的字节序列。字符串可以包含任意的数据

// 对字符串操作的4个包bytes、strings、strconv、unicode包
//bytes包操作[]byte。因为字符串是只读的，因此逐步构创建字符串会导致很多分配和复制。使用 bytes.Buffer类型会更高。
//strings包提供切割，索引，前缀,查找替换等功能。
//strconv包提供了布尔型、整型数、浮点数和对应字符串的相互转换，还提供了双引号转义相 关的转换。
//unicode包提供了IsDigit、IsLetter、IsUpper和IsLower等类似功能，它们用于给字符分类。

func TestStringNil(t *testing.T) {

	// string不能设为nil
	//var x string = nil //error,cannot use nil as type string in assignment

	// 不能用string和nil比较
	//if x == nil { //error,invalid operation: x == nil (mismatched types string and nil)
	//}

	// 默认值为""而不是nil
	var x string //defaults to "" (zero value)

	fmt.Println("len:", len(x))

	if x == "" {
		x = "default"
	}

	fmt.Println("x", x)
	time.Sleep(2 * time.Second)
}

// 将 string 转为 rune slice（此时 1 个 rune 可能占多个 byte），直接更新 rune 中的字符
func TestUpdate(t *testing.T) {
	x := "text"
	xRunes := []rune(x) // 1个英文字符占用一个rune对应一个byte，一个中文字符可能占用1个rune对应多个byte
	fmt.Println("len xRunes:", len(xRunes))
	xRunes[0] = '我'
	x = string(xRunes)
	fmt.Println(x) // 我ext
	fmt.Println("len xRunes:", len(xRunes))

	fmt.Println(string(65))
}

func TestGet(t *testing.T) {
	x := "我ext"
	fmt.Println(x[0], len(x)) // byte长度

	n := utf8.RuneCountInString(x)
	fmt.Println(n)
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

// string是Go的内建类型，但对它的读写操作并非线程安全的，原因在于它的内部实际上是通过struct存储的，
// 可以在runtime/string.go里面看到它的内部定义。
// type stringStruct struct {
//	str unsafe.Pointer
//	len int
//}
// 在频繁的写入操作中，另一协程可能读到部分写入的结果
// 在并发场景下，string跟interface一样，都是需要使用atomic包来保证读写的原子性。
// 看来很多对象只要涉及多线程的读写就要小心race问题
func TestStringNotSafe(t *testing.T) {
	var a = "0"
	ch := make(chan string)
	go func() {
		i := 1
		for {
			if i%2 == 0 {
				a = "0"
			} else {
				a = "aa"
			}
			time.Sleep(1 * time.Millisecond)
			i++
		}
	}()
	// 检查string变量问题
	go func() {
		for {
			b := a
			if b != "0" && b != "aa" {
				ch <- b
			}
		}
	}()
	for i := 0; i < 10; i++ {
		fmt.Println("get strange string: ", <-ch)
	}
}

func TestStringByte(t *testing.T) {
	s := "hello, world"
	// len()返回字符串字节数目（不是rune数）
	fmt.Println(len(s)) // 12个长度，每个英文字符占用1个字节
	// 通过索引可以访问某个字节值,非ASCII的UTF8字符有多个字节
	fmt.Println(s[0], s[7]) // "104 119" ('h' and 'w')

	// 循环是循环的字节，而非字符
	for i, r := range "Hello, 世界ꡐ" { // 世占用3个字节
		fmt.Printf("%d\t%q\t%d\n", i, r, r)
	}

	fmt.Println(s[:5])        // "hello"
	fmt.Println(s[7:])        // "world"
	fmt.Println(s[:])         // "hello, world"
	fmt.Println("hi" + s[5:]) //hi world
}

func TestStringConcat(t *testing.T) {
	s1 := "a"
	s2 := "b"
	var build strings.Builder
	build.WriteString(s1)
	build.WriteString(s2)
	s3 := build.String()
	fmt.Println("s3:", s3)
}
