package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// 有时候，我们只想用其他的程序（也许是非 Go 程序）来完全替代当前的 Go 进程。可以使用经典的 exec方法的Go实现
// go run exec.go
func main() {
	fmt.Println("main start")
	// 将执行 ls 命令。Go 需要提供我们需要执行的可执行文件的绝对路径，所以我们将使用exec.LookPath 来得到它（大概是 /bin/ls）。
	binary, lookErr := exec.LookPath("ls")
	if lookErr != nil {
		panic(lookErr)
	}
	//Exec 需要的参数是切片的形式的。给 ls 一些基本的参数。注意，第一个参数需要是程序名。
	args := []string{"ls", "-a", "-l", "-h"}
	//Exec 同样需要使用环境变量。这里我们仅提供当前的环境变量。
	env := os.Environ()
	fmt.Println("exec before")
	//os.Exec 调用。如果这个调用成功，那么我们的进程将在这里被替换成 /bin/ls -a -l -h 进程。如果存在错误，那么我们将会得到一个返回值。
	// 执行后，当前进程的代码都是args[0]了，就不会执行下面的命令了
	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}
	fmt.Println("exec after")
}
