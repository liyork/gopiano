package main

import (
	"fmt"
	"runtime"
)

//这个 `GOMAXPROCS` 环境变量允许您限制操作系统线程数，它能同时执行用户级 Go 代码。Go 1.5 版本开始，`GOMAXPROCS` 的默认值应该是您 Unix 机器的内核数。

func getGOMAXPROCS() int {
	return runtime.GOMAXPROCS(0)
}

// go run maxprocs.go
// export GOMAXPROCS=800; go run maxprocs.go
func main() {
	fmt.Printf("GOMAXPROCS:%d\n", getGOMAXPROCS())
}
