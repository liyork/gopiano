package cmd

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"testing"
)

func TestCommandBase1(t *testing.T) {
	// returns the Cmd struct to execute the named program with the given arguments.
	dateCmd := exec.Command("date")
	// runs the command and returns its standard output.会调用StartProcess,启动另一个进程
	dateOut, err := dateCmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println("> date")
	fmt.Println(string(dateOut))
}

//从外部进程的stdin输入数据并从 stdout 收集结果
func TestCommandBase2(t *testing.T) {
	grepCmd := exec.Command("grep", "hello")

	// 要在start之前获取
	grepOut, _ := grepCmd.StdoutPipe()
	// The pipe will be closed automatically after Wait sees the command exit.
	// A caller need only call Close to force the pipe to close sooner.
	grepIn, _ := grepCmd.StdinPipe()

	grepCmd.Start()

	grepIn.Write([]byte("hello grep\ngoodbye grep"))
	// 主动触发,不然out总是等待
	grepIn.Close()

	// ReadAll reads from r until an error or EOF and returns the data it read.
	grepBytes, _ := ioutil.ReadAll(grepOut)
	// Wait waits for the command to exit and waits for any copying to
	// stdin or copying from stdout or stderr to complete.
	grepCmd.Wait() // 等待grep命令执行完

	fmt.Println("> grep hello")
	fmt.Println(string(grepBytes))
}

//当需要提供一个明确的命令和参数数组来生成命令，使用通过一个字符串生成一个完整的命令，可以使用 bash命令的 -c 选项：
//       -c string If the -c option is present, then commands are read from string.  If there are arguments after the  string,  they  are
//                 assigned to the positional parameters, starting with $0.
func TestCommandBase3(t *testing.T) {
	lsCmd := exec.Command("bash", "-c", "ls -a -l -h")
	lsOut, err := lsCmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println("> ls -a -l -h")
	fmt.Println(string(lsOut))
}

// 注意 Go 并不提供一个经典的 Unix fork 函数。通常这不是个问题，因为运行 Go 协程，生成进程和执行进程覆盖了fork 的大多数使用用场景。
