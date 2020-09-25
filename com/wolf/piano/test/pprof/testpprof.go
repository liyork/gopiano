package main

import (
	"flag"
	"github.com/liyork/gopiano/com/wolf/piano/test/common"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

const (
	col = 10000
	row = 10000
)

func main() {
	// 适用于执行一会可以完成的应用
	//execMethod()

	// 适用于一直跑的程序
	execHttp()
}

// go build -o tpprof
// ./tpprof -cpuprofile cpu.prof -memprofile mem.prof
// 得到:  cpu.prof  # cpu运行数据   mem.prof  # 内存数据

// 进入交互   go tool pprof xx(应用程序) yy（应用程序的prof文件）
// go tool pprof tpprof cpu.prof
// top10
// tree 5
// 输入web ，会在/tmp下生成svg文件，用浏览器打开(前提brew install -y graphviz)

// 生成调用图
// go tool pprof -svg ./tpprof cpu.profpwd
// 火焰图
// go-torch --binaryname=./tpprof --binaryinput=cpu.prof

func execMethod() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		// 开始监控cpu
		if err := pprof.StartCPUProfile(f); err != nil { //监控cpu
			log.Fatal("could not start CPU profile: ", err)
		}
		// 停止监控cpu，使用StartCPUProfile后一定要调用该函数停止监控。
		// 注意，若pprof.StopCPUProfile() 未执行则会看到 prof 文件是空的，所以需要进行关闭或者对于长时间运行的服务使用http方式
		defer pprof.StopCPUProfile()
	}
	logic()
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		// GC，获取最新的数据信息
		runtime.GC()
		// 把堆中的内存分配信息写入分析文件中
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}
}

// 主逻辑区，进行一些简单的代码运算
func logic() {
	x := [row][col]int{}
	s := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			x[i][j] = s.Intn(100000)
		}
	}
	for i := 0; i < row; i++ {
		tmp := 0
		for j := 0; j < col; j++ {
			tmp += x[i][j]
		}
	}
}

// 安装graphviz：
// brew install graphviz

// http://localhost:6060/debug/pprof/
// /debug/pprof/
// /debug/pprof/goroutine  go tool pprof -http=":8081" goroutine
// /debug/pprof/cmdline
// /debug/pprof/profile?seconds=5  go tool pprof -http=":8081" profile
// /debug/pprof/trace?seconds=5  go tool trace -http=":8081" trace
// /debug/pprof/symbol
// /debug/pprof/heap
// 对服务器在一段时间内进行数据采样，然后分析服务器的耗时和性能: 该命令后会对服务进行30s的采样
// 先启动服务，后go tool pprof http://localhost:6060/debug/pprof/profile，之后输入top
// 或者wget http://localhost:6060/debug/pprof/profile 然后go tool pprof tpprof profile

// 火焰图:
// 安装:go get github.com/uber/go-torch

// 生成文件后用浏览器打开：go-torch -u http://localhost:6060 -p > torch.svg
// 直接生成：go-torch -u http://localhost:6060 -t 30
// 对tpprof文件生成：go-torch tpprof profile
// 内存火焰图：go-torch http://localhost:6060/debug/pprof/heap --colors mem  -f mem.svg
// wget  http://localhost:6060/debug/pprof/heap
// go tool pprof testpprof(应用程序) heap(heap文件)

func execHttp() {
	common.StartRoutine()
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	time.Sleep(122222 * time.Second)
}
