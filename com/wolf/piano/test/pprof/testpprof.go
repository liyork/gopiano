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
	//execMethod()
	execHttp()
}

// go build
// ./pprof -cpuprofile cpu.prof -memprofile mem.prof
// cpu.prof  # cpu运行数据,mem.prof  # 内存数据
// go tool pprof cpu.prof --进入交互
// top10
// tree 5
// brew install -y graphviz
// web ，会在/tmp下生成svg文件，用浏览器打开
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
		defer pprof.StopCPUProfile()
	}
	// 主逻辑区，进行一些简单的代码运算
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

// http://localhost:6060/debug/pprof/
// /debug/pprof/
// /debug/pprof/cmdline
// /debug/pprof/profile
// /debug/pprof/symbol
// /debug/pprof/trace
// 对服务器在一段时间内进行数据采样，然后分析服务器的耗时和性能: 该命令后会对服务进行30s的采样
// go tool pprof http://localhost:6060/debug/pprof/profile
func execHttp() {
	common.StartRoutine()
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	time.Sleep(122222 * time.Second)
}
