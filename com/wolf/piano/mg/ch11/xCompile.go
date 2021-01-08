package main

import (
	"fmt"
	"runtime"
)

// go run xCompile.go
// 设置 `GOOS` 和 `GOARCH` 环境变量来分别指定目标系统和 CPU 架构
// env GOOS=linux GOARCH=arm go build xCompile.go
// file xCompile
// env GOOS=linux GOARCH=386 go build xCompile.go  --有 Intel 处理器的 Debian Linux 机器
func main() {
	fmt.Println("You are using ", runtime.Compiler, " ")
	fmt.Println("on a", runtime.GOARCH, "machine")
	fmt.Println("with Go version", runtime.Version())
}
