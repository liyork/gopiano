package gotest

import (
	"flag"
	"testing"
)

// 单元测试（unit testing），是指对软件中的最小可测试单元进行检查和验证。
//单元就是人为规定的最小的被测功能模块。
//单元测试是在软件开发过程中要进行的最低级别的测试活动，软件的独立单元将在与程序的其他部分相隔离的情况下进行测试。
//以_test结尾
//默认的情况下，go test 命令不需要任何的参数，它会自动把你源码包下面所有 test 文件测试完毕，当然你也可以带上参数。
//
//常用的参数：
//-bench regexp 执行相应的 benchmarks，例如 -bench=.；
//-cover 开启测试覆盖率；
//-run regexp 只运行 regexp 匹配的函数，例如 -run=Array 那么就执行包含有 Array 开头的函数；
//-v 显示测试的详细命令。
//
//单元测试源码文件可以由多个测试用例组成，每个测试用例函数需要以Test为前缀
//测试用例文件不会参与正常源码编译，不会被包含到可执行文件中。
//测试用例文件使用 go test 指令来执行，没有也不需要 main() 作为函数入口。所有在以_test结尾的源码内以Test开头的函数会自动被执行。
//测试用例可以不传入 *testing.T 参数。
//
//测试这个文件里的所有测试用例
//go test -v gotestdemo_test.go
//-run跟随的测试用例的名称支持正则表达式，使用-run TestA$即可只执行 TestA 测试用例
//go test -v -run TestA gotestdemo_test.go
//go test -v -run TestA$ gotestdemo_test.go

func TestHelloWorld(t *testing.T) {
	t.Log("hello world")
}

func TestA(t *testing.T) {
	t.Log("A1")
	t.Log("A")
}
func TestAK(t *testing.T) {
	t.Log("AK")
}
func TestB(t *testing.T) {
	t.Log("B")
}
func TestC(t *testing.T) {
	t.Log("C")
}

// go test把-args后面的参数(以字符串形式传入)带到测试中去
// go test -run TestArgs -v -args "cloud"
func TestArgs(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	argList := flag.Args() // flag.Args() 返回 -args 后面的所有参数，以切片表示，每个元素代表一个参数
	for _, arg := range argList {
		if arg == "cloud" {
			t.Log("Running in cloud.")
		} else {
			t.Log("Running in other mode.")
		}
	}
}

// go test -run TestAdd -json
// -o <file>,指定生成的二进制可执行程序，并执行测试，测试结束不会删除该程序
// -bench regexp,是包含关系
// -bench . 或

// -benchtime <t>s,指定每个性能测试的执行时间，如果不指定，则使用默认时间1s。
// -cpu 1,2,4  每个测试执行多次
// -count n,指定每个测试执行的次数，默认执行一次。
// -failfast,指定如果有测试出现失败，则立即停止测试
// -list regexp,只是列出匹配成功的测试函数，并不真正执行
// -parallel n,指定测试的最大并发数。
// -run regexp,根据正则表达式执行单元测试和示例测试
// -timeout d,默认情况下，测试执行超过10分钟就会超时而退出
// -v,详细
// -benchmem,默认情况下，性能测试结果只打印运行次数、每个操作耗时。使用-benchmem则可以打印每个操作分配的字节数、每个操作分配的对象数
