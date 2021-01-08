package main

//
//import (
//	"bufio"
//	"fmt"
//	"os"
//	"strings"
//	"syscall"
//)
//
//var maxSyscalls = 0
//
//// 请注意，`SYSCALLS`文件的信息取自`syscall`包的文档，它将每个系统调用与一个数字相关联，该数字是系统调用的内部`Go`表示形式。该文件主要用于打印被跟踪程序所使用的系统调用的名称。
//const SYSCALLFILE = "SYSCALLS"
//
//// go run traceSyscall.go
//// go run traceSyscall.go ls /tmp/
//// strace -c ls /tmp
//func main() {
//	var SYSTEMCALLS []string
//	f, err := os.Open(SYSCALLFILE)
//	defer f.Close()
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	scanner := bufio.NewScanner(f)
//	for scanner.Scan() {
//		line := scanner.Text()
//		line = strings.Replace(line, " ", "", -1)
//		line = strings.Replace(line, "SYS_", "", -1)
//		temp := strings.ToLower(strings.Split(line, "=")[0])
//		SYSTEMCALLS = append(SYSTEMCALLS, temp)
//		maxSyscalls++
//	}
//
//	COUNTER := make([]int, maxSyscalls)
//	var regs syscall.PtraceRegs
//	cmd := exec.Command(os.Args[1], os.Args[2:]...)
//	cmd.Stdin = os.Stdin
//	cmd.Stdout = os.Stdout
//	cmd.Stderr = os.Stderr
//	cmd.SysProcAttr = &syscall.SysProcAttr{Ptrace: true}
//	err = cmd.Start()
//	err = cmd.Wait()
//	if err != nil {
//		fmt.Println("Wait:", err)
//	}
//	pid := cmd.Process.Pid
//	fmt.Println("Process ID:", pid)
//
//	before := true
//	forCount := 0
//	for {
//		if before {
//			err := syscall.PtraceGetRegs(pid, &regs)
//			if err != nil {
//				break
//			}
//			if regs.Orig_rax > uint64(maxSyscalls) {
//				fmt.Println("Unknown:", regs.Orig_rax)
//				return
//			}
//			COUNTER[regs.Orig_rax]++
//			forCount++
//		}
//		// `syscall.PtraceSyscall()`函数的作用是：告诉`Go`继续执行正在被跟踪的程序，但是当程序进入或退出系统调用时停止执行，这正是我们想要的！由于每个系统调用在被调用之前和完成其工作之后都会被跟踪，因此我们使用`before`变量来计算每个系统调用仅一次。
//		err = syscall.PtraceSyscall(pid, 0)
//		if err != nil {
//			fmt.Println("PtraceSyscall:", err)
//			return
//		}
//		_, err = syscall.Wait4(pid, nil, 0, nil)
//		if err != nil {
//			fmt.Println("Wait4:", err)
//			return
//		}
//		before = !before
//	}
//
//	for i, x := range COUNTER {
//		if x != 0 {
//			fmt.Println(SYSTEMCALLS[i], "->", x)
//		}
//	}
//	fmt.Println("Total System Calls:", forCount)
//}
//
//
//
