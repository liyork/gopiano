package base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"testing"
	"time"
	"unicode/utf8"
)

//1. 左大括号 { 一般不能单独放一行
// 遵守分号注入规则（automatic semicolon injection）：编译器会在每行代码尾部特定分隔符后加 ;(分号) 来分隔多条语句，比如会在 ) 后加分号：
//// 错误示例
//func main()
//{
//	println("hello world")
//}
//
//// 等效于
//func main();	// 无函数体
//{
//	println("hello world")
//}
//// ./main.go: missing function body ./main.go: syntax error: unexpected semicolon or newline before {
//
//// 正确示例
//func main() {
//	println("hello world")
//}

//注意代码块等特殊情况：
// { 并不遵守分号注入规则，不会在其后边自动加分，此时可换行
//func main() {
//	{
//		println("hello world")
//	}
//}

// 2. 未使用的变量
// 在函数体代码中有未使用的变量，则无法通过编译，不过全局变量声明但不使用是可以的
// 解决方案:直接注释或移除未使用的变量

// 3. 未使用的 import
// 如果你 import 一个包，但包中的变量、函数、接口和结构体一个都没有用到的话，将编译失败。
//可以使用 _ 下划线符号作为别名来忽略导入的包，从而避免编译错误，这只会执行 package 的 init()

// 4. 简短声明的变量只能在函数内部使用

// 5. 不能用简短声明方式来单独为一个变量重复声明,:= 左侧至少有一个新变量，才允许多变量的重复声明：
// 正确示例
//func main() {
//	one := 0
//	one, two := 1, 2	// two 是新变量，允许 one 的重复声明。比如 error 处理经常用同名变量 err
//}

// 6. 不能使用简短声明来设置字段的值
// struct中的变量字段不能使用 := 来赋值，要使用预定义的变量(先声明变量)来避免解决：
// 错误示例
//func main() {
//	var data info
//	data.result, err := work() // error: non-name data.result on left side of :=
//	fmt.Printf("info: %+v\n", data)
//}

type info struct {
	result int
}

func work() (int, error) {
	return 3, nil
}

// 正确示例
func TestAssignStructVariable(t *testing.T) {
	var data info
	var err error // err 需要预声明

	data.result, err = work()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("info: %+v\n", data)
}

// 7. 不小心覆盖了变量
// 使用 vet 工具来诊断这种变量覆盖,默认不做覆盖检查,添加 -shadow 选项来启用
//go tool vet -shadow main.go
//使用 go-nyet 来做进一步的检测：
//$GOPATH/bin/go-nyet main.go
//可见，简短声明就是声明一个新变量了。
func TestOverwriteVariable(t *testing.T) {
	x := 1
	println(x) // 1
	{
		println(x) // 1
		x := 2
		println(x) // 2	// 新的 x 变量的作用域只在代码块内部
	}
	println(x) // 1
}

// 8. 显式类型的变量无法使用 nil 来初始化
// nil 是 interface、function、pointer、map、slice 和 channel 类型变量的默认初始值。
// 但声明时若不指定类型，编译器也无法推断出变量的具体类型。
func TestCorrectUseType(t *testing.T) {
	// 错误示例
	//    var x = nil	// error: use of untyped nil
	//	_ = x

	// 正确示例

	var x interface{} = nil
	_ = x
}

// 9. 直接使用值为 nil 的 slice、map
// 允许对值为 nil 的 slice 添加元素，但对值为 nil 的 map 添加元素则会造成运行时 panic
func TestCorrectUseNil(t *testing.T) {
	// map 错误示例
	//var m map[string]int
	//m["one"] = 1 // error: panic: assignment to entry in nil map
	// map正确示例
	// m := make(map[string]int)// map 的正确声明，分配了实际的内存

	// slice 正确示例
	var s []int
	s = append(s, 1)
}

// 10. map 容量
// 在创建 map 类型的变量时可以指定容量，但不能像 slice 一样使用 cap() 来检测分配空间的大小：
func TestMapCap(t *testing.T) {
	// 错误示例
	//	m := make(map[string]int, 99)
	//	println(cap(m)) 	// error: invalid argument m1 (type map[string]int) for cap
}

// 11. string 类型的变量值不能为 nil
func TestStringDefault(t *testing.T) {
	// 错误示例
	//var s string = nil // cannot use nil as type string in assignment
	//if s == nil {      // invalid operation: s == nil (mismatched types string and nil)
	//	s = "default"
	//}

	// 正确示例
	var s string // 字符串类型的零值是空串 ""
	if s == "" {
		s = "default"
	}
}

// 12. Array 类型的值作为函数参数
// 在 C/C++ 中，数组（名）是指针。将数组作为参数传进函数时，相当于传递了数组内存地址的引用，在函数内部会改变该数组的值。
//在 Go 中，数组是值。作为参数传进函数时，传递的是数组的原始值拷贝，此时在函数内部是无法更新该数组的：
// 数组使用值拷贝传参
func TestPassArray2(t *testing.T) {
	x := [3]int{1, 2, 3}

	func(arr [3]int) {
		arr[0] = 7
		fmt.Println(arr) // [7 2 3]
	}(x)
	fmt.Println(x) // [1 2 3]
}

// 如果想修改参数数组：
//直接传递指向这个数组的指针类型：
func TestPassArrayChangeValue1(t *testing.T) {
	x := [3]int{1, 2, 3}

	func(arr *[3]int) {
		(*arr)[0] = 7
		fmt.Println(arr) // &[7 2 3]
	}(&x) // 指针传递
	fmt.Println(x) // [7 2 3]
}

// 直接使用slice进行修改：
// 即使函数内部得到的是 slice 的值拷贝，但依旧会更新 slice 的原始数据（底层 array）
// 会修改 slice 的底层 array，从而修改 slice
func TestPassArrayChangeValue2(t *testing.T) {
	x := []int{1, 2, 3}
	func(arr []int) {
		arr[0] = 7
		fmt.Println(x) // [7 2 3]
	}(x) // slice拷贝
	fmt.Println(x) // [7 2 3]
}

// 13. range 遍历 slice 和 array 时混淆了返回值
// Go 中的 range 在遍历时会生成 2 个值，第一个是元素索引，第二个是元素的值：
func TestCorrectUseRange(t *testing.T) {
	// 错误示例
	//x := []string{"a", "b", "c"}
	//for v := range x {
	//	fmt.Println(v)	// 1 2 3
	//}

	// 正确示例
	x := []string{"a", "b", "c"}
	for _, v := range x { // 使用 _ 丢弃索引
		fmt.Println(v)
	}
}

// 14. slice 和 array 其实是一维数据
// 看起来 Go 支持多维的 array 和 slice，可以创建数组的数组、切片的切片，但其实并不是。
// 对依赖动态计算多维数组值的应用来说，就性能和复杂度而言，用 Go 实现的效果并不理想。
// 可以使用原始的一维数组、“独立“ 的切片、“共享底层数组”的切片3方式来创建动态的多维数组。
// a.使用原始的一维数组：要做好索引检查、溢出检测、以及当数组满时再添加值时要重新做内存分配。
// b.使用“独立”的切片分两步：
//  创建外部 slice
//	对每个内部 slice 进行内存分配
//	注意内部的 slice 相互独立，使得任一内部 slice 增缩都不会影响到其他的 slice
// 使用各自独立的 6 个 slice 来创建 [2][3] 的动态多维数组
func TestCreateMultiArrayUsingIndepentSlice(t *testing.T) {
	x := 2
	y := 4

	table := make([][]int, x)
	for i := range table {
		table[i] = make([]int, y)
	}
}

// c.使用“共享底层数组”的切片
//  创建一个存放原始数据的容器 slice
//  创建其他的 slice
//  切割原始 slice 来初始化其他的 slice
func TestCreateMultiArrayUsingShareSlice(t *testing.T) {
	h, w := 2, 4
	raw := make([]int, h*w)

	for i := range raw {
		raw[i] = i
	}

	// 初始化原始 slice
	fmt.Println(raw, &raw[4]) // [0 1 2 3 4 5 6 7] 0xc420012120

	table := make([][]int, h)
	for i := range table {

		// 等间距切割原始 slice，创建动态多维数组 table
		// 0: raw[0*4: 0*4 + 4]
		// 1: raw[1*4: 1*4 + 4]
		table[i] = raw[i*w : i*w+w]
	}

	fmt.Println(table, &table[1][0]) // [[0 1 2 3] [4 5 6 7]] 0xc420012120
}

// 15. 访问 map 中不存在的 key
// Go 则会返回元素对应数据类型的零值，比如 nil、'' 、false 和 0，取值操作总有值返回，故不能通过取出来的值来判断 key 是不是在 map 中。
func TestAccessMapNil(t *testing.T) {
	// 错误的 key 检测方式
	//x := map[string]string{"one": "2", "two": "", "three": "3"}
	//if v := x["two"]; v == "" {
	//	fmt.Println("key two is no entry") // 键 two 存不存在都会返回的空字符串
	//}

	// 正确示例,检查 key 是否存在可以用 map 直接访问，检查返回的第二个参数即可
	x := map[string]string{"one": "2", "two": "", "three": "3"}
	if _, ok := x["two"]; !ok {
		fmt.Println("key two is no entry")
	}
}

// 16. string 类型的值是常量，不可更改
// 尝试使用索引遍历字符串，来更新字符串中的个别字符，是不允许的。
// string 类型的值是只读的二进制 byte slice，如果真要修改字符串中的字符，将 string 转为 []byte 修改后，再转为 string 即可：
func TestUpdateString(t *testing.T) {
	// 错误示例
	//x := "text"
	//x[0] = "T" // error: cannot assign to x[0]
	//fmt.Println(x)

	// 注意： 示例并不是更新字符串的正确姿势，因为一个 UTF8 编码的字符可能会占多个字节，比如汉字就需要 3~4 个字节来存储，此时更新其中的一个字节是错误的。
	x := "text"
	xBytes := []byte(x)
	xBytes[0] = 'T' // 注意此时的 T 是 rune 类型
	x = string(xBytes)
	fmt.Println(x) // Text

	//更新字串的正确姿势：将 string 转为 rune slice（此时 1 个 rune 可能占多个 byte），直接更新 rune 中的字符
	x = "text"
	xRunes := []rune(x)
	xRunes[0] = '我'
	x = string(xRunes)
	fmt.Println(x) // 我ext
}

// 17. string 与 byte slice 之间的转换，参与转换的是拷贝的原始值
//Go 在 string 与 byte slice 相互转换上优化了两点，避免了额外的内存分配：
//在 map[string] 中查找 key 时，使用了对应的原始[]byte，避免做 m[string(key)] 的内存分配
//使用 for range 迭代 string时转换为原始[]byte 的迭代：for i,v := range []byte(str) {...}

// 18. string 与索引操作符
// 对字符串用索引访问返回的不是字符，而是一个索引位置的byte值。
// 如果需要使用 for range 迭代访问字符串中的字符（unicode code point / rune），标准库中有 "unicode/utf8" 包来做 UTF8 的相关解码编码。另外 utf8string 也有像 func (s *String) At(i int) rune 等很方便的库函数。
func TestStringIndex(t *testing.T) {
	x := "ascii"
	fmt.Println("len:", len([]byte(x)))
	fmt.Println(x[0])        // 97
	fmt.Printf("%T\n", x[0]) // uint8
}

// 19. 字符串并不都是 UTF8 文本
// string类型变量的值不必是都是UTF8 文本，可以包含任意的值。
// 只有字符串是文字字面值时才是 UTF8 文本，字串可以通过转义来包含其他数据。
//判断字符串是否是 UTF8 文本，可使用 "unicode/utf8" 包中的 ValidString() 函数：
func TestValidateString(t *testing.T) {
	str1 := "ABC"
	fmt.Println(utf8.ValidString(str1)) // true

	str2 := "A\xfeC"
	fmt.Println(utf8.ValidString(str2)) // false

	str3 := "A\\xfeC"
	fmt.Println(utf8.ValidString(str3)) // true	// 把转义字符转义成字面值
}

// 20. 字符串的长度
func TestStringLen(t *testing.T) {
	char := "♥"
	// Go 的内建函数 len() 返回的是字符串的 byte 数量，而不是像 Python 中那样是计算 Unicode 字符数。
	fmt.Println(len([]byte(char))) // 3
	fmt.Println(len(char))         // 3

	// 如果要得到字符串的字符数，可使用 "unicode/utf8" 包中的 RuneCountInString(str string) (n int)
	fmt.Println(utf8.RuneCountInString(char)) // 1

	// 注意： RuneCountInString 并不总是返回我们看到的字符数，因为有的字符会占用 2 个 rune：
	char = "é"
	fmt.Println(len(char))                    // 3
	fmt.Println(utf8.RuneCountInString(char)) // 2
	fmt.Println("cafe\u0301")                 // café是法文的 cafe，实际上是两个 rune 的组合

}

// 21. 在多行 array、slice、map 语句中缺少 , 号
func TestColonInMap(t *testing.T) {
	// 声明语句中 } 折叠到单行后，尾部的 , 是必需的
	//x := []int{
	//	1,
	//	2 // syntax error: unexpected newline, expecting comma or }。
	//}
	y := []int{1, 2}
	z := []int{1, 2}
	fmt.Println(y, z)
}

// 22. log.Fatal 和 log.Panic 不只是 log
// Go 的 log 包在调用 Fatal*()、Panic*() 时能做更多日志外的事，如中断程序的执行等：
func TestFatal(t *testing.T) {
	log.Fatal("Fatal level log: log entry") // 输出信息后，程序终止执行
	log.Println("Nomal level log: log entry")
}

// 23. 对内建数据结构的操作并不是同步的
// 尽管 Go 本身有大量的特性来支持并发，但并不保证并发的数据安全，用户需自己保证变量等数据以原子操作更新。如map等
//goroutine 和 channel 是进行原子操作的好方法，或使用 "sync" 包中的锁。

// 24. range 迭代 string 得到的值
// range 得到的索引是字符值（Unicode point / rune）第一个字节的位置，与其他编程语言不同，这个索引并不直接是字符在字符串中的位置。
//注意一个字符可能占多个 rune，比如法文单词 café 中的 é。操作特殊字符可使用norm 包。
// for range 迭代会尝试将 string 翻译为 UTF8 文本，对任何无效的码点都直接使用 0XFFFD rune（�）UNicode 替代字符来表示。
// 如果 string 中有任何非 UTF8 的数据，应将 string 保存为 byte slice 再进行操作。
func TestStringRange(t *testing.T) {
	data := "A\xfe\x02\xff\x04"
	for _, v := range data {
		fmt.Printf("%#x ", v) // 0x41 0xfffd 0x2 0xfffd 0x4	// 错误
	}

	fmt.Println("")

	for _, v := range []byte(data) {
		fmt.Printf("%#x ", v) // 0x41 0xfe 0x2 0xff 0x4	// 正确
	}
}

// 25. range 迭代 map
//注意每次迭代都可能产生不一样的顺序

// 26. switch 中的 fallthrough 语句
// switch 语句中的 case 代码块会默认带上 break，但可以使用 fallthrough 来强制执行下一个 case 代码块。
func TestSwitchFallthrough(t *testing.T) {
	isSpace := func(char byte) bool {
		switch char {
		case ' ': // 空格符会直接 break，返回 false
		// fallthrough	// 若打开。在 case 代码块末尾使用 fallthrough，强制执行下一个 case 代码块。
		case '\t':
			return true
		}
		return false
	}

	fmt.Println(isSpace('\t')) // true
	fmt.Println(isSpace(' '))  // false

	// 可以改写 case 为多条件判断：
	isSpace = func(char byte) bool {
		switch char {
		case ' ', '\t':
			return true
		}
		return false
	}
	fmt.Println(isSpace('\t')) // true
	fmt.Println(isSpace(' '))  // true
}

// 27. 自增和自减运算
// Go 特立独行，去掉了前置操作，同时 ++、— 只作为运算符(对操作数进行操作)而非表达式(没有最后求值)。
func TestAddAdd(t *testing.T) {
	// 错误示例
	//data := []int{1, 2, 3}
	//i := 0
	//++i			// syntax error: unexpected ++, expecting }
	//fmt.Println(data[i++])	// syntax error: unexpected ++, expecting :

	// 正确示例
	data := []int{1, 2, 3}
	i := 0
	i++
	fmt.Println(data[i]) // 2
}

// 28. 按位取反
// Go 重用 ^(XOR) 操作符来按位取反：
func TestXor(t *testing.T) {
	// 错误的取反操作
	//fmt.Println(~2)		// bitwise complement operator is ^

	// 正确示例
	var d uint8 = 2
	fmt.Printf("%08b\n", d)  // 00000010
	fmt.Printf("%08b\n", ^d) // 11111101

	//^ 也是按位异或XOR操作符
	// 一个操作符能重用两次，是因为一元的 NOT 操作 NOT 0x02，与二元的 XOR 操作 0x22 XOR 0xff 是一致的。
	//Go 也有特殊的操作符 AND NOT &^ 操作符，不同位才取1。
	var a uint8 = 0x82
	var b uint8 = 0x02
	fmt.Printf("%08b [A]\n", a)
	fmt.Printf("%08b [B]\n", b)

	fmt.Printf("%08b (NOT B)\n", ^b)
	fmt.Printf("%08b ^ %08b = %08b [B XOR 0xff]\n", b, 0xff, b^0xff)

	fmt.Printf("%08b ^ %08b = %08b [A XOR B]\n", a, b, a^b)
	fmt.Printf("%08b & %08b = %08b [A AND B]\n", a, b, a&b)
	fmt.Printf("%08b &^%08b = %08b [A 'AND NOT' B]\n", a, b, a&^b)
	fmt.Printf("%08b&(^%08b)= %08b [A AND (NOT B)]\n", a, b, a&(^b))
}

// 29. 运算符的优先级
// 除了位清除（bit clear）操作符，Go 也有很多和其他语言一样的位操作符，但优先级另当别论。
// 优先级列表：
// Precedence    Operator
//    5             *  /  %  <<  >>  &  &^
//    4             +  -  |  ^
//    3             ==  !=  <  <=  >  >=
//    2             &&
//    1             ||

// 30. 不导出的 struct 字段无法被 encode
// 以小写字母开头的字段成员是无法被外部直接访问的，所以 struct 在进行 json、xml、gob 等格式的 encode 操作时，这些私有字段会被忽略，导出时得到零值：
type MyData struct {
	One int
	two string
}

func TestExportUppercase(t *testing.T) {
	in := MyData{1, "two"}
	fmt.Printf("%#v\n", in) // main.MyData{One:1, two:"two"}

	encoded, _ := json.Marshal(in)
	fmt.Println(string(encoded)) // {"One":1}	// 私有字段 two 被忽略了

	var out MyData
	json.Unmarshal(encoded, &out)
	fmt.Printf("%#v\n", out) // main.MyData{One:1, two:""}
}

// 看来go中方法调用都是值拷贝，只不过值是实际值还是指针地址
// 31. 程序退出时还有 goroutine 在执行
// 程序默认不等所有 goroutine 都执行完才退出，这点需要特别注意：
func TestRoutineQuit(t *testing.T) {
	//workerCount := 2
	//for i := 0; i < workerCount; i++ {
	//	go doIt(i)
	//}
	//time.Sleep(1 * time.Second)
	//// 主程序会直接退出,不等两个 goroutine 执行完就直接退出了
	//fmt.Println("all done!")

	//常用解决办法：使用 "WaitGroup" 变量，它会让主程序等待所有 goroutine 执行完毕再退出。
	//如果你的 goroutine 要做消息的循环处理等耗时操作，可以向它们发送一条 kill 消息来关闭它们。或直接关闭一个它们都等待接收数据的 channel：
	// 等待所有 goroutine 执行完毕
	// 进入死锁
	var wg sync.WaitGroup
	done := make(chan struct{})
	ch := make(chan interface{})

	workerCount := 2
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go doIt2(i, ch, done, &wg) // wg 传指针，doIt() 内部会改变 wg 的值
	}
	for i := 0; i < workerCount; i++ { // 向 ch 中发送数据，关闭 goroutine
		ch <- i
	}
	close(done)
	wg.Wait()
	close(ch)
	fmt.Println("all done!")
}

func doIt(workerID int) {
	fmt.Printf("[%v] is running\n", workerID)
	time.Sleep(3 * time.Second) // 模拟 goroutine 正在执行
	fmt.Printf("[%v] is done\n", workerID)
}

// 之前不能停止原因是 goroutine 得到的 "WaitGroup" 变量是 var wg WaitGroup 的一份拷贝值，即 doIt2() 传参只传值。所以哪怕在每个 goroutine 中都调用了 wg.Done()， 主程序中的 wg 变量并不会受到影响
// 需要使用传址方式为 WaitGroup 变量传参
func doIt2(workerID int, ch <-chan interface{}, done <-chan struct{}, wg *sync.WaitGroup) {
	fmt.Printf("[%v] is running\n", workerID)
	defer wg.Done()
	for {
		select {
		case m := <-ch:
			fmt.Printf("[%v] m => %v\n", workerID, m)
		case <-done:
			fmt.Printf("[%v] is done\n", workerID)
			return
		}
	}
}

// 32. 向无缓冲的 channel 发送数据，只要 receiver 准备好了sender就会立刻返回
// 只有在数据未被 receiver 处理时，sender 才会阻塞
// 因运行环境而异，在 sender 发送完数据后，receiver 的 goroutine 可能没有足够的时间处理下一个数据。???
// 这里测试是，肯定发送时要阻塞直到接收处理时才会可以
func TestBlockChan(t *testing.T) {
	ch := make(chan string)

	go func() {
		for m := range ch {
			fmt.Println("Processed:", m)
			time.Sleep(3 * time.Second) // 模拟需要长时间运行的操作
		}
	}()

	ch <- "cmd.1"
	ch <- "cmd.2" // 不会被接收处理
}

// 33. 向已关闭的 channel 发送数据会造成 panic
// 从已关闭的 channel 接收数据是安全的：
// 接收状态值 ok 是 false 时表明 channel 中已没有数据可以接收了。类似的，从有缓冲的 channel 中接收数据，缓存的数据获取完再没有数据可取时，状态值也是 false
// 向已关闭的 channel 中发送数据会造成 panic
// 可使用一个废弃 channel done 来告诉剩余的 goroutine 无需再向 ch 发送数据。此时 <- done的结果是 {}
func TestStopChan1(t *testing.T) {
	ch := make(chan int, 2)
	go func() {
		ch <- 1
	}()
	go func() {
		ch <- 2
	}()

	time.Sleep(1 * time.Second)
	close(ch)

	// 即使关闭，也会读完才停止range
	for a := range ch {
		fmt.Println("a:", a)
		time.Sleep(1 * time.Second)
	}
}

func TestStopChan2(t *testing.T) {
	ch := make(chan int)
	done := make(chan struct{})

	for i := 0; i < 3; i++ {
		go func(idx int) {
			select {
			case ch <- (idx + 1) * 2:
				fmt.Println(idx, "Send result")
			case <-done: // 信号chan
				fmt.Println(idx, "Exiting")
			}
		}(i)
	}
	fmt.Println("Result: ", <-ch)
	close(done) // 关闭信号
	time.Sleep(3 * time.Second)
}

// 34. 使用了值为 nil 的 channel
// 在一个值为 nil 的 channel 上发送和接收数据将永久阻塞：
func TestChanNil(t *testing.T) {
	var ch chan int // 未初始化，值为 nil
	for i := 0; i < 3; i++ {
		go func(i int) {
			ch <- i
		}(i)
	}
	fmt.Println("Result: ", <-ch)
	time.Sleep(2 * time.Second)
}

// 利用操作nil的chan会阻塞，这个死锁的特性，可以用在 select 中动态的打开和关闭 case 语句块：
func TestChangeCase(t *testing.T) {
	inCh := make(chan int)
	outCh := make(chan int)

	go func() {
		var in <-chan int = inCh
		var out chan<- int
		var val int
		for {
			select {
			case out <- val: // 写入val后设定out=nil关闭此case
				println("--------")
				out = nil
				in = inCh
			case val = <-in: // 先读取，赋值给val，后为out赋值，则上面case激活，可以将val写入out，将in=nil关闭此case
				println("++++++++++")
				out = outCh
				in = nil
			}
		}
	}()
	go func() {
		for r := range outCh {
			fmt.Println("Result: ", r)
		}
	}()
	time.Sleep(0)
	inCh <- 1
	inCh <- 2
	time.Sleep(3 * time.Second)
}

// 34. 若函数 receiver 传参是传值方式，则无法修改参数的原有值
// 方法 receiver 的参数与一般函数的参数类似：如果声明为值，那方法体得到的是一份参数的值拷贝，此时对参数的任何修改都不会对原有值产生影响。
//除非 receiver 参数是 map 或 slice 类型的变量，并且是以指针方式更新 map 中的字段、slice 中的元素的，才会更新原有值:
type data2 struct {
	num   int
	key   *string
	items map[string]bool
}

func (this *data2) pointerFunc() {
	this.num = 7
}

func (this data2) valueFunc() {
	this.num = 8
	*this.key = "valueFunc.key"
	this.items["valueFunc"] = true
}
func TestChangeValueInStruct(t *testing.T) {
	key := "key1"

	d := data2{1, &key, make(map[string]bool)}
	fmt.Printf("num=%v  key=%v  items=%v\n", d.num, *d.key, d.items)

	d.pointerFunc() // 修改 num 的值为 7
	fmt.Printf("num=%v  key=%v  items=%v\n", d.num, *d.key, d.items)

	d.valueFunc() // 修改 key 和 items 的值
	fmt.Printf("num=%v  key=%v  items=%v\n", d.num, *d.key, d.items)
}

// 35. 关闭 HTTP 的响应体
// 使用 HTTP 标准库发起请求、获取响应时，即使你不从响应中读取任何数据或响应为空，也需要手动关闭响应体
// 应该先检查 HTTP 响应错误为 nil，再调用 resp.Body.Close() 来关闭响应体：
func TestCloseHttpBody(t *testing.T) {
	// 绝大多数情况下的正确关闭方式
	resp, err := http.Get("https://api.ipify.org?format=json")
	checkError(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	checkError(err)

	fmt.Println(string(body))
}

func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// 绝大多数上述请求失败的情况下，resp 的值为 nil 且 err 为 non-nil。
// 但如果你得到的是重定向错误，那它俩的值都是 non-nil，最后依旧可能发生内存泄露
// 2 个解决办法：
func TestCloseHttpBody2(t *testing.T) {
	resp, err := http.Get("http://www.baidu.com")

	// 1.直接在处理 HTTP 响应错误的代码块中，直接关闭非 nil 的响应体。.关闭 resp.Body 的正确姿势
	if resp != nil {
		defer resp.Body.Close()
	}

	// 2.或者,手动调用 defer 来关闭响应体：
	checkError(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	checkError(err)

	fmt.Println(string(body))

	// resp.Body.Close() 早先版本的实现是读取响应体的数据之后丢弃，保证了 keep-alive 的 HTTP 连接能重用处理不止一个请求。
	// 但 Go 的最新版本将读取并丢弃数据的任务交给了用户，如果你不处理，HTTP 连接可能会直接关闭而非重用，参考在 Go 1.5 版本文档。
	//如果程序大量重用 HTTP 长连接，你可能要在处理完成响应体之后加入：
	_, err = io.Copy(ioutil.Discard, resp.Body) // 手动丢弃读取完毕的数据
	//如果你需要完整读取响应，上边的代码是需要写的。比如在解码 API 的 JSON 响应数据：
	//json.NewDecoder(resp.Body).Decode(&data)
}

// 36. 关闭 HTTP 连接
// 一些支持 HTTP1.1 或 HTTP1.0 配置了 connection: keep-alive 选项的服务器会保持一段时间的长连接。
// 但标准库 "net/http" 的连接默认只在服务器主动要求关闭时才断开，所以你的程序可能会消耗完 socket 描述符
// 解决办法有 2 个，请求结束后：
//   直接设置请求变量的 Close 字段值为 true，每次请求结束后就会主动关闭连接。
//   设置 Header 请求头部选项 Connection: close，然后服务器返回的响应头部也会有这个选项，此时 HTTP 标准库会主动断开连接。
func TestCloseHttp1(t *testing.T) {
	// 主动关闭连接
	req, err := http.NewRequest("GET", "http://golang.org", nil)
	checkError(err)

	req.Close = true
	//req.Header.Add("Connection", "close")	// 等效的关闭方式

	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	checkError(err)

	body, err := ioutil.ReadAll(resp.Body)
	checkError(err)

	fmt.Println(string(body))
}

// 可以创建一个自定义配置的 HTTP transport 客户端，用来取消 HTTP 全局的复用连接：
func TestCloseHttp2(t *testing.T) {
	tr := http.Transport{DisableKeepAlives: true}
	client := http.Client{Transport: &tr}

	resp, err := client.Get("https://golang.google.cn/")
	if resp != nil {
		defer resp.Body.Close()
	}
	checkError(err)

	fmt.Println(resp.StatusCode) // 200

	body, err := ioutil.ReadAll(resp.Body)
	checkError(err)

	fmt.Println(len(string(body)))
}

// 根据需求选择使用场景：
//若你的程序要向同一服务器发大量请求，使用默认的保持长连接。
//若你的程序要连接大量的服务器，且每台服务器只请求一两次，那收到请求后直接关闭连接。或增加最大文件打开数 fs.file-max 的值。

//37. 将 JSON 中的数字解码为 interface 类型
// 在 encode/decode JSON 数据时，Go 默认会将数值当做 float64 处理，比如下边的代码会造成 panic：
func TestJsonFloat1(t *testing.T) {
	var data = []byte(`{"status": 200}`)
	var result map[string]interface{}

	if err := json.Unmarshal(data, &result); err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%T\n", result["status"]) // float64
	//var status = result["status"].(int)  // 类型断言错误

	// 如果你尝试 decode 的 JSON 字段是整型，可以：
	//将 int 值转为 float 统一使用
	//将 decode 后需要的 float 值转为 int 使用
	//将 decode 的值转为 int 使用
	var status = uint64(result["status"].(float64))
	fmt.Println("Status value: ", status)
}

// 使用 Decoder 类型来 decode JSON 数据，明确表示字段的值类型
// 指定字段类型
func TestJsonFloat2(t *testing.T) {
	var data = []byte(`{"status": 200}`)
	var result map[string]interface{}

	var decoder = json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	if err := decoder.Decode(&result); err != nil {
		log.Fatalln(err)
	}

	var status, _ = result["status"].(json.Number).Int64()
	fmt.Println("Status value: ", status)
}

// 你可以使用 string 来存储数值数据，在 decode 时再决定按 int 还是 float 使用
// 将数据转为 decode 为 string
func TestJsonFloat3(t *testing.T) {
	var data = []byte(`{"status": 200}`)
	var result map[string]interface{}
	var decoder = json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&result); err != nil {
		log.Fatalln(err)
	}
	var status uint64
	err := json.Unmarshal([]byte(result["status"].(json.Number).String()), &status)
	checkError(err)
	fmt.Println("Status value: ", status)
}

// 使用 struct 类型将你需要的数据映射为数值型
// struct 中指定字段类型
func TestJsonFloat4(t *testing.T) {
	var data = []byte(`{"status": 200}`)
	var result struct {
		Status uint64 `json:"status"`
	}

	err := json.NewDecoder(bytes.NewReader(data)).Decode(&result)
	checkError(err)
	fmt.Printf("Result: %+v", result)
}

// 可以使用 struct 将数值类型映射为 json.RawMessage 原生数据类型
// 适用于如果 JSON 数据不着急 decode 或 JSON 某个字段的值类型不固定等情况：
// 状态名称可能是 int 也可能是 string，指定为 json.RawMessage 类型
func TestJsonFloat5(t *testing.T) {
	records := [][]byte{
		[]byte(`{"status":200, "tag":"one"}`),
		[]byte(`{"status":"ok", "tag":"two"}`),
	}

	var result struct {
		StatusCode uint64
		StatusName string
		Status     json.RawMessage `json:"status"`
		Tag        string          `json:"tag"`
	}

	for idx, record := range records {
		err := json.NewDecoder(bytes.NewReader(record)).Decode(&result)
		checkError(err)

		var name string
		err = json.Unmarshal(result.Status, &name)
		if err == nil {
			result.StatusName = name
		}

		var code uint64
		err = json.Unmarshal(result.Status, &code)
		if err == nil {
			result.StatusCode = code
		}

		fmt.Printf("[%v] result => %+v\n", idx, result)
	}
}

// 38. struct、array、slice 和 map 的值比较
type data3 struct {
	num     int
	fp      float32
	complex complex64
	str     string
	char    rune
	yes     bool
	events  <-chan string
	handler interface{}
	ref     *byte
	raw     [10]byte
}

// 可以使用相等运算符 == 来比较结构体变量，前提是两个结构体的成员都是可比较的类型：
func TestMapEqual1(t *testing.T) {
	v1 := data3{}
	v2 := data3{}
	fmt.Println("v1 == v2: ", v1 == v2) // true
}

type data4 struct {
	num    int
	checks [10]func() bool   // 无法比较
	doIt   func() bool       // 无法比较
	m      map[string]string // 无法比较
	bytes  []byte            // 无法比较
}

//如果两个结构体中有任意成员是不可比较的，将会造成编译错误。注意数组成员只有在数组元素可比较时候才可比较
func TestMapEqual2(t *testing.T) {
	//v1 := data4{}
	//v2 := data4{}
	//编译错误
	//fmt.Println("v1 == v2: ", v1 == v2)

	// Go 提供了一些库函数来比较那些无法使用 == 比较的变量，比如使用 "reflect" 包的 DeepEqual() ：
	// 用来比较，相等运算符无法比较的元素，这种比较方式可能比较慢，根据你的程序需求来使用。
	// struct
	v1 := data4{}
	v2 := data4{}
	fmt.Println("v1 == v2: ", reflect.DeepEqual(v1, v2)) // true
	// map
	m1 := map[string]string{"one": "a", "two": "b"}
	m2 := map[string]string{"two": "b", "one": "a"}
	fmt.Println("v1 == v2: ", reflect.DeepEqual(m1, m2)) // true
	// slice
	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	// 注意两个 slice 相等，值和顺序必须一致，否则[]int{1, 3, 2}就false了
	fmt.Println("v1 == v2: ", reflect.DeepEqual(s1, s2)) // true
}

func TestDeepEqual(t *testing.T) {
	// DeepEqual() 还有其他用法：
	var b1 []byte = nil
	b2 := []byte{}
	fmt.Println("b1 == b2: ", reflect.DeepEqual(b1, b2)) // false

	// 注意：DeepEqual() 并不总适合于比较 slice
	var str = "one"
	var in interface{} = "one"
	fmt.Println("str == in: ", reflect.DeepEqual(str, in)) // true

	// slice顺序有要求
	v1 := []string{"one", "two"}
	v2 := []string{"two", "one"}
	fmt.Println("v1 == v2: ", reflect.DeepEqual(v1, v2)) // false

	data := map[string]interface{}{
		"code":  200,
		"value": []string{"one", "two"},
	}
	encoded, _ := json.Marshal(data)
	var decoded map[string]interface{}
	json.Unmarshal(encoded, &decoded)
	fmt.Println("data:", data)
	fmt.Println("decoded:", decoded)
	// 再反解析出来的map比较时不一样
	fmt.Println("data == decoded: ", reflect.DeepEqual(data, decoded)) // false
	data2 := map[string]interface{}{
		"code":  200,
		"value": []string{"one", "two"},
	}
	//直接构造的map比较则一样。
	fmt.Println("data == decoded: ", reflect.DeepEqual(data, data2)) // true
}

// 如果要大小写不敏感来比较 byte 或 string 中的英文文本，可以使用 "bytes" 或 "strings" 包的 ToUpper() 和 ToLower() 函数都转换一样后再比较。
// 比较其他语言的 byte 或 string，应使用 bytes.EqualFold() 和 strings.EqualFold()
//如果 byte slice 中含有验证用户身份的数据（密文哈希、token 等），不应再使用 reflect.DeepEqual()、bytes.Equal()、 bytes.Compare()。
// 这三个函数容易对程序造成 timing attacks，此时应使用 "crypto/subtle" 包中的 subtle.ConstantTimeCompare() 等函数

// reflect.DeepEqual() 认为空 slice 与 nil slice 并不相等，但是byte.Equal() 会认为二者相等：
func TestBytesEqual(t *testing.T) {
	var b1 []byte = nil
	b2 := []byte{}

	// b1 与 b2 长度相等、有相同的字节序
	// nil 与 slice 在字节上是相同的
	fmt.Println("b1 == b2: ", bytes.Equal(b1, b2)) // true
}

// 39. 从 panic 中恢复
// 在一个 defer 延迟执行的函数中调用 recover() ，它便能捕捉 / 中断 panic
func TestRecoverUse(t *testing.T) {
	// 错误的 recover 调用示例
	//recover()	// 什么都不会捕捉
	//panic("not good")	// 发生 panic，主程序退出
	//recover()	// 不会被执行
	//println("ok")

	// 正确的 recover 调用示例，recover() 仅在 defer 执行的函数中调用才会生效。
	//defer func() {
	//	fmt.Println("recovered: ", recover())
	//}()
	//panic("not good")

	// 错误的调用示例，doRecover方法中的recover()不行
	//defer func() {
	//	doRecover()
	//}()
	//panic("not good")

	// 这样可以，看来defer的方法就直接用recover不要再调用其他函数了再来一层了
	defer doRecover()
	panic("good")
}

func doRecover() {
	fmt.Println("recobered: ", recover())
}

// 40. 在 range 迭代 slice、array、map 时通过更新引用来更新元素
// 在 range 迭代中，得到的值其实是元素的一份值拷贝，更新拷贝并不会更改原来的元素，即是拷贝的地址并不是原有元素的地址
func TestRangeUpdate(t *testing.T) {
	data := []int{1, 2, 3}
	for _, v := range data {
		v += 10 // data 中原有元素是不会被修改的
	}
	fmt.Println("data: ", data) // data:  [1 2 3]

	// 如果要修改原有元素的值，应该使用索引直接访问：
	data = []int{1, 2, 3}
	for i, v := range data {
		data[i] = v * 10
	}
	fmt.Println("data: ", data) // data:  [10 20 30]

	// 如果集合保存的是指向值的指针，需稍作修改。依旧需要使用索引访问元素，不过可以使用 range 出来的元素直接更新原有值：
	data1 := []*struct{ num int }{{1}, {2}, {3}} // 定义匿名结构体并直接初始化
	for _, v := range data1 {
		v.num = v.num * 10 // 直接使用指针更新
	}
	fmt.Println(data1[0], data1[1], data1[2]) // &{10} &{20} &{30}
}

// 41. slice 中隐藏的数据
// 从 slice 中重新切出新 slice 时，新 slice 会引用原 slice 的底层数组。
// 如果跳了这个坑，程序可能会分配大量的临时 slice 来指向原底层数组的部分数据，将导致难以预料的内存使用。
func TestSliceFromSliceErr(t *testing.T) {
	data := get()
	fmt.Println(len(data), cap(data), &data[0]) // 3 10000 0xc420080000  // 重新分配容量为 10000 的 slice
}

func get() []byte {
	raw := make([]byte, 10000)
	fmt.Println(len(raw), cap(raw), &raw[0]) // 10000 10000 0xc420080000
	return raw[:3]                           // 重新分配一个slice，容量为 10000
}

// 可以通过拷贝临时 slice 的数据，而不是重新切片来解决：
// 这个似乎也并没有解决多slice问题，应该是底层指向不同数组了吧。???
func TestSliceFromSliceCorrect(t *testing.T) {
	data := get2()
	fmt.Println(len(data), cap(data), &data[0]) // 3 3 0xc4200160b8
}

func get2() (res []byte) {
	raw := make([]byte, 10000)
	fmt.Println(len(raw), cap(raw), &raw[0]) // 10000 10000 0xc420080000
	res = make([]byte, 3)
	copy(res, raw[:3])
	return
}

// 42. Slice 中数据的误用
// 错误使用 slice 的拼接示例
func TestSliceUseErr(t *testing.T) {
	path := []byte("AAAA/BBBBBBBBB")
	sepIndex := bytes.IndexByte(path, '/') // 4
	println(sepIndex)

	dir1 := path[:sepIndex]
	dir2 := path[sepIndex+1:]
	println("dir1: ", string(dir1)) // AAAA
	println("dir2: ", string(dir2)) // BBBBBBBBB

	dir1 = append(dir1, "suffix"...)
	println("dir1: ", string(dir1))         // AAAAsuffix
	println("dir2: ", string(dir2))         // uffixBBBB
	println("current path: ", string(path)) // AAAAsuffixBBBB，底层被改了

	// 拼接的结果不是正确的 AAAAsuffix/BBBBBBBBB，因为 dir1、 dir2 两个 slice 引用的数据都是 path 的底层数组，修改 dir1(通过append)同时也修改了path，也导致了 dir2 的修改
	path = bytes.Join([][]byte{dir1, dir2}, []byte{'/'})
	println("new path: ", string(path)) // AAAAsuffix/uffixBBBB	// 错误结果
}

// 解决上面错误方法：
//重新分配新的 slice 并拷贝需要的数据
//使用完整的 slice 表达式：input[low:high:max]，容量便调整为 max - low
func TestSliceUseCorrect(t *testing.T) {
	// 使用 full slice expression
	path := []byte("AAAA/BBBBBBBBB")
	sepIndex := bytes.IndexByte(path, '/') // 4
	// 第三个参数是用来控制 dir1 的新容量，再往 dir1 中 append 超额元素时，将分配新的 buffer 来保存。而不是覆盖原来的 path 底层数组
	dir1 := path[:sepIndex:sepIndex] // 此时dir1的cap指定为4， 而不是先前的 16
	dir2 := path[sepIndex+1:]
	dir1 = append(dir1, "suffix"...) // 扩容后为指向新数组

	path = bytes.Join([][]byte{dir1, dir2}, []byte{'/'})
	println("dir1: ", string(dir1))     // AAAAsuffix
	println("dir2: ", string(dir2))     // BBBBBBBBB
	println("new path: ", string(path)) // AAAAsuffix/BBBBBBBBB
}

// 43. 旧 slice
// 当你从一个已存在的 slice 创建新 slice 时，二者的数据指向相同的底层数组。
// 如果你的程序使用这个特性，那需要注意 "旧"（stale） slice 问题。
// 某些情况下，向一个 slice 中追加元素而它指向的底层数组容量不足时，将会重新分配一个新数组来存储数据。而其他 slice 还指向原来的旧底层数组。
func TestOldSlice(t *testing.T) {
	// 超过容量将重新分配数组来拷贝值、重新存储
	oldSlice := []int{1, 2, 3}
	fmt.Println(len(oldSlice), cap(oldSlice), oldSlice) // 3 3 [1 2 3 ]

	newSlice := oldSlice[1:]
	fmt.Println(len(newSlice), cap(newSlice), newSlice) // 2 2 [2 3]

	for i := range newSlice {
		newSlice[i] += 20
	}
	// 此时的 oldSlice 与 newSlice 是指向同一个底层数组的
	fmt.Println(oldSlice) // [1 22 23]
	fmt.Println(newSlice) // [22 23]

	newSlice = append(newSlice, 4) // 向容量为 2 的 newSlice 中再追加元素，此时将分配新数组来存，并赋值给newSlice

	for i := range newSlice {
		newSlice[i] += 10
	}
	fmt.Println(oldSlice) // [1 22 23]	// 此时的 oldSlice 不再更新，为旧数据
	fmt.Println(newSlice) // [32 33 14]
}

// 44. 类型声明与方法
// 从一个现有的非 interface 类型创建新类型时，并不会继承原有的方法：
// 定义 Mutex 的自定义类型
type myMutex sync.Mutex

func TestExtendFromNoInterface(t *testing.T) {
	//sync.Mutex是结构体
	//var mtx myMutex
	//mtx.Lock()
	//mtx.UnLock()

	// 如果需要使用原类型的方法，可将原类型以匿名字段的形式嵌到你定义的新 struct 中：
	// 类型以字段形式直接嵌入
	type myLocker struct {
		sync.Mutex
	}

	var locker myLocker
	locker.Lock()
	locker.Unlock()

	// interface 类型声明也保留它的方法集：
	type myLocker2 sync.Locker
	var locker2 myLocker2
	locker2.Lock()
	locker2.Unlock()
}

// 45. 跳出 for-switch 和 for-select 代码块
// 没有指定标签的 break 只会跳出 switch/select 语句，若不能使用 return 语句跳出的话，可为 break 跳出标签指定的代码块：
func TestBreakFromFor(t *testing.T) {
	// break 配合 label 跳出指定代码块
loop:
	for {
		switch {
		case true:
			time.Sleep(time.Second)
			fmt.Println("breaking out...")
			//break	// 这样写会死循环，一直打印 breaking out...

			// goto 虽然也能跳转到指定位置，但依旧会再次进入 for-switch，死循环。
			//goto loop

			// 正确方式，跳出到loop之后不执行for
			break loop
		}
	}
	fmt.Println("out...")
}

// 46. for 语句中的迭代变量与闭包函数
// for 语句中的`迭代变量`在每次迭代中都会重用，即 for 中创建的闭包函数接收到的参数始终是同一个变量，在 goroutine 开始执行时都会得到同一个迭代值：
func TestForClosureErr(t *testing.T) {
	data := []string{"one", "two", "three"}
	for _, v := range data {
		go func() {
			// 输出 three three three
			fmt.Println(v)
		}()
	}

	time.Sleep(3 * time.Second)
}

// 最简单的解决方法：无需修改 goroutine 函数，在 for 内部使用局部变量保存迭代值，再传参：
func TestForClosureCorrect1(t *testing.T) {
	data := []string{"one", "two", "three"}

	for _, v := range data {
		vCopy := v
		go func() {
			fmt.Println(vCopy)
		}()
	}

	time.Sleep(3 * time.Second)
	// 输出 one two three
}

// 另一个解决方法：直接将当前的迭代值以参数形式传递给匿名函数：
func TestForClosureCorrect2(t *testing.T) {
	data := []string{"one", "two", "three"}

	for _, v := range data {
		go func(in string) {
			fmt.Println(in)
		}(v)
	}

	time.Sleep(3 * time.Second)
	// 输出 one two three
}

// 注意下边这个稍复杂的 3 个示例区别：
type field struct {
	name string
}

func (p *field) print() {
	fmt.Println(p.name)
}

func TestForClosureCorrect3(t *testing.T) {
	// 错误示例
	data := []field{{"one"}, {"two"}, {"three"}}
	for _, v := range data {
		go v.print()
	}
	time.Sleep(3 * time.Second)
	// 输出 three three three

	// 正确示例，值
	data = []field{{"one"}, {"two"}, {"three"}}
	for _, v := range data {
		v := v
		go v.print()
	}
	time.Sleep(3 * time.Second)
	// 输出 one two three

	fmt.Println("")

	// 正确示例，指针，到了地址就不会复用v了??
	data2 := []*field{{"one"}, {"two"}, {"three"}}
	for _, v := range data2 { // 此时迭代值 v 是三个元素值的地址，每次 v 指向的值不同
		go v.print()
	}
	time.Sleep(3 * time.Second)
	// 输出 one two three
}

// 47. defer 函数的参数值
// 对 defer 延迟执行的函数，它的参数会在声明时候就会求出具体值，而不是在执行时才求值：
func TestEvaluateInDefer(t *testing.T) {
	// 在 defer 函数中参数会提前求值
	var i = 1
	defer fmt.Println("result: ", func() int { return i * 2 }())
	i++
}

// 48. defer 函数的执行时机
// 对 defer 延迟执行的函数，会在调用它的函数结束时执行，而不是在调用它的语句块结束时执行，注意区分开。
// 比如在一个长时间执行的函数里，内部 for 循环中使用 defer 来清理每次迭代产生的资源调用，就会出现问题：
func TestDeferErr(t *testing.T) {
	targets := walkDir()

	for _, target := range targets {
		f, err := os.Open(target)
		if err != nil {
			fmt.Println("bad target:", target, "error:", err) //error:too many open files
			break
		}
		defer f.Close() // 错误：在每次 for 语句块结束时，不会关闭文件资源
		// 使用 f 资源...
	}
}

func walkDir() []string {
	if len(os.Args) != 2 {
		os.Exit(1)
	}
	// 命令行参数指定目录名
	dir := os.Args[1]
	start, err := os.Stat(dir)
	if err != nil || !start.IsDir() {
		os.Exit(2)
	}
	// 遍历读取目录下的文件
	var targets []string
	filepath.Walk(dir, func(fPath string, fInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !fInfo.Mode().IsRegular() {
			return nil
		}

		targets = append(targets, fPath)
		return nil
	})
	return targets
}

// 先创建 10000 个文件：
// #!/bin/bash
//for n in {1..10000}; do
//	echo content > "file${n}.txt"
//done

// 解决办法：defer 延迟执行的函数写入匿名函数中：
func TestDeferCorrect(t *testing.T) {
	targets := walkDir()

	for _, target := range targets {
		func() { // 匿名函数，每次for内执行
			f, err := os.Open(target)
			if err != nil {
				fmt.Println("bad target:", target, "error:", err)
				return // 在匿名函数内使用 return 代替 break 即可
			}
			defer f.Close() // 匿名函数执行结束，调用关闭文件资源。也可以直接调用 f.Close() 来关闭，不用defer。
			// 使用 f 资源...
		}()
	}
}

// 49. 失败的类型断言
// 在类型断言语句中，断言失败则会返回目标类型的“零值”，断言变量与原来变量混用可能出现异常情况：
func TestAssert(t *testing.T) {
	// 错误示例
	var data interface{} = "great"

	// data 混用，若是断言错误，则源data值被修改
	if data, ok := data.(int); ok {
		fmt.Println("[is an int], data: ", data)
	} else {
		fmt.Println("[not an int], data: ", data) // [isn't a int], data:  0
	}

	// 正确示例
	var data2 interface{} = "great"

	if res, ok := data2.(int); ok {
		fmt.Println("[is an int], data: ", res)
	} else {
		fmt.Println("[not an int], data: ", data2) // [not an int], data:  great
	}
}

// 50. 阻塞的 gorutinue 与资源泄露
func TestRoutineLeak(t *testing.T) {
	// 在搜索重复时依旧每次都起一个 goroutine 去处理，每个 goroutine 都把它的搜索结果发送到结果 channel 中，
	// 但是channel 中收到的第一条数据会直接返回。
	// 返回完第一条数据后，其他 goroutine 的搜索结果怎么处理？他们自己的协程如何处理？他们被阻塞

	// 在 First() 中的结果 channel 是无缓冲的，这意味着只有第一个 goroutine 能返回，由于没有 receiver，
	// 其他的 goroutine 会在发送上一直阻塞。如果你大量调用，则可能造成资源泄露。
}

func First(query string, replicas []func(string) string) string {
	c := make(chan string)
	replicaSearch := func(i int) { c <- replicas[i](query) }
	for i := range replicas {
		go replicaSearch(i)
	}
	return <-c
}

// 为避免上述泄露，你应该确保所有的 goroutine 都能正确退出，
// 有 2 个解决方法：
// a.使用带缓冲的 channel，确保能接收全部 goroutine 的返回结果：
func FirstCorrect1(query string, replicas []func(string) string) string {
	// 带缓冲
	c := make(chan string, len(replicas))
	searchReplica := func(i int) { c <- replicas[i](query) }
	for i := range replicas {
		go searchReplica(i)
	}
	return <-c // 有一个结果退出
}

// b.使用 select 语句，+ 能保存一个值的缓冲channel default 语句：
//default 的缓冲 channel 保证了即使结果 channel 收不到数据，也不会阻塞 goroutine
// 这个不会有随机select问题吗?而且大家都没有执行完都执行default，那么就没有返回值了
func FirstCorrect2(query string, replicas []func(string) string) string {
	c := make(chan string, 1)
	searchReplica := func(i int) {
		select {
		case c <- replicas[i](query):
		default:
		}
	}
	for i := range replicas {
		go searchReplica(i)
	}
	return <-c
}

// c.使用特殊的废弃（cancellation） channel 来中断剩余 goroutine 的执行：
func FirstCorrect3(query string, replicas []func(string) string) string {
	c := make(chan string)
	// 特殊chan
	done := make(chan struct{})
	defer close(done)
	searchReplica := func(i int) {
		select {
		case c <- replicas[i](query):
		case <-done:
		}
	}
	for i := range replicas {
		go searchReplica(i)
	}
	return <-c
}

// 51. 使用指针作为方法的 receiver
// 只要值是可寻址的，就可以在值上直接调用指针方法。即使对一个方法，它的 receiver 是指针就足矣。
// 但不是所有值都是可寻址的，比如 map中的元素、通过 interface 引用的变量：
type data5 struct {
	name string
}

type printer5 interface {
	print5()
}

func (p *data5) print5() {
	fmt.Println("name: ", p.name)
}

func TestSeekValue(t *testing.T) {
	d1 := data5{"one"}
	d1.print5() // d1 变量可寻址，可直接调用指针 receiver 的方法

	// cannot use data literal (type data5) as type printer5 in assignment:
	//data does not implement printer5 (print5 method has pointer5 receiver)
	// 类型不匹配
	//var in printer5 = data5{"two"}
	//in.print5()

	// cannot call pointer method on m["x"]
	// cannot take the address of m["x"]
	//m := map[string]data{
	//	"x": data{"three"},
	//}
	//m["x"].print() // m["x"] 是不可寻址的	// 变动频繁
}

// 52. 更新 map 字段的值
// 如果 map 一个字段的值是 struct 类型，则无法直接更新该 struct 的单个字段：
func TestUpdateStructInMap(t *testing.T) {
	m := map[string]data5{
		"x": {"Tom"},
	}
	// 因为 map 中的元素是不可寻址的。
	//m["x"].name = "Jerry"

	// 更新 map 中 struct 元素的字段值，有 2 个方法：
	// a.使用局部变量
	// 提取整个 struct 到局部变量中，修改字段值后再整个赋值
	r := m["x"]
	r.name = "Jerry"
	m["x"] = r
	fmt.Println(m) // map[x:{Jerry}]

	// b.使用指向元素的指针
	m1 := map[string]*data5{
		"x": {"Tom"},
	}

	m1["x"].name = "Jerry" // 直接修改 m["x"] 中的字段
	fmt.Println(m["x"])    // &{Jerry}

	// 但是要注意下边这种误用：
	m1["z"].name = "what???" //这个值不在map中存在，nil
	fmt.Println(m["x"])

	// 需区分开的是，slice 中的元素可寻址：
	s := []data{{"Tom"}}
	s[0].name = "Jerry"
	fmt.Println(s) // [{Jerry}]
}

// 53. nil interface 和 nil interface 值
// 虽然 interface 看起来像指针类型，但它不是。
// interface 类型的变量只有在类型和值均为 nil 时才为 nil
// 如果你的 interface 变量的值是跟随其他变量变化的（雾），与 nil 比较相等时小心：
func TestInterfaceNil2(t *testing.T) {
	var data *byte
	var in interface{}

	fmt.Println(data, data == nil) // <nil> true
	fmt.Println(in, in == nil)     // <nil> true

	fmt.Printf("type:%T\n", in)
	in = data
	fmt.Println(in, in == nil) // <nil> false	// data 值为 nil，但 in 值不为 nil,in的类型为*byte
	fmt.Printf("type:%T\n", in)
}

// 如果你的函数返回值类型是 interface，更要小心这个坑：
func TestInterfaceNil3(t *testing.T) {
	// 错误示例
	doIt := func(arg int) interface{} {
		var result *struct{} = nil
		if arg > 0 {
			result = &struct{}{}
		}
		// 返回了interface{}，类型为*struct{}，值为nil
		return result
	}

	if res := doIt(-1); res != nil {
		fmt.Println("Good result: ", res) // Good result:  <nil>
		fmt.Printf("%T\n", res)           // *struct {}	// res 不是 nil，它的值为 nil
		fmt.Printf("%v\n", res)           // <nil>
	}

	// 正确示例
	doIt1 := func(arg int) interface{} {
		var result *struct{} = nil
		if arg > 0 {
			result = &struct{}{}
		} else {
			return nil // 明确指明返回 nil
		}
		return result
	}

	if res := doIt1(-1); res != nil {
		fmt.Println("Good result: ", res)
	} else {
		fmt.Println("Bad result: ", res) // Bad result:  <nil>
	}
}

// 54. 堆栈变量
// go你并不总是清楚你的变量是分配到了堆还是栈。
// 在 C++ 中使用 new 创建的变量总是分配到堆内存上的，但在 Go 中即使使用 new()、make() 来创建变量，变量为内存分配位置依旧归 Go 编译器管。
// Go 编译器会根据变量的大小及其 "escape analysis" 的结果来决定变量的存储位置，故能准确返回本地变量的地址，这在 C/C++ 中是不行的。
// 在 go build 或 go run 时，加入 -m 参数，能准确分析程序的变量分配位置：

// 55. GOMAXPROCS、Concurrency（并发）and Parallelism（并行）
// Go 1.4 及以下版本，程序只会使用 1 个执行上下文 / OS 线程，即任何时间都最多只有 1 个 goroutine 在执行。
// Go 1.5 版本将可执行上下文的数量设置为runtime.NumCPU()返回的逻辑 CPU 核心数，这个数与系统实际总的 CPU 逻辑核心数是否一致，
// 取决于你的 CPU 分配给程序的核心数，可以使用 GOMAXPROCS 环境变量或者动态的使用 runtime.GOMAXPROCS()来调整。
// 误区：GOMAXPROCS 表示执行 goroutine 的 CPU 核心数，参考文档
// GOMAXPROCS 的值是可以超过 CPU 的实际数量的，在 1.5 中最大为 256上下文，表示可以让m同时运行的上下文
func TestGOMAXPROCS(t *testing.T) {
	fmt.Println(runtime.GOMAXPROCS(-1)) // 4
	fmt.Println(runtime.NumCPU())       // 4
	runtime.GOMAXPROCS(20)
	fmt.Println(runtime.GOMAXPROCS(-1)) // 20
	runtime.GOMAXPROCS(300)
	fmt.Println(runtime.GOMAXPROCS(-1)) // Go 1.9.2 // 300
}

// 56. 读写操作的重新排序
// Go 可能会重排一些操作的执行顺序，可以保证在一个 goroutine 中操作是顺序执行的，但不保证多 goroutine 的执行顺序：
var _ = runtime.GOMAXPROCS(3)
var ax, bx int

func u1() {
	ax = 1
	bx = 2
}

func u2() {
	ax = 3
	bx = 4
}

func px() {
	println(ax)
	println(bx)
}

// 如果想保持多 goroutine 像代码中的那样顺序执行，可以使用 channel 或 sync 包中的锁机制等
func TestMultiRoutineSeq(t *testing.T) {
	go u1() // 多个 goroutine 的执行顺序不定
	go u2()
	go px()
	time.Sleep(1 * time.Second)
}

// 57. 优先调度
func TestPrioritySchedule(t *testing.T) {
	// 你的程序可能出现一个 goroutine 在运行时阻止了其他 goroutine 的运行，比如程序中有一个不让调度器运行的 for 循环：
	done := false

	go func() {
		done = true
	}()

	// for 的循环体不必为空，但如果代码不会触发调度器执行，将出现问题。
	for !done {
	}

	println("done !")

	// 调度器会在 GC、Go 声明、阻塞 channel、阻塞系统调用和锁操作后再执行，也会在非内联函数调用时执行：
	// 可以添加 -m 参数来分析 for 代码块中调用的内联函数：
	done1 := false

	go func() {
		done1 = true
	}()

	for !done1 {
		println("not done !") // 并不内联执行
	}

	println("done !")

	// 你也可以使用 runtime 包中的 Gosched() 来 手动启动调度器：
	done2 := false

	go func() {
		done2 = true
	}()

	for !done2 {
		runtime.Gosched()
	}

	println("done !")
}
