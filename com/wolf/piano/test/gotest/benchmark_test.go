package gotest

import (
	"fmt"
	"testing"
)

//基准测试可以测试一段程序的运行性能及耗费 CPU 的程度。
//测试代码需要保证函数可重入性及无状态
//测试代码不使用全局变量等带有记忆性质的数据结构。避免多次运行同一段代码时的环境不一致，不能假设 N 值范围。

// go help test
// go test [-c] [-i] [build flags] [packages] [flags for test binary]
//-c : 编译go test成为可执行的二进制文件，但是不运行测试。
//-i : 安装测试包依赖的package，但是不运行测试。
//关于build flags， 调用go help build，这些是编译运行过程中需要使用到的参数，一般设置为空
//关于packages， 调用go help packages，这些是关于包的管理，一般设置为空
//关于flags for test binary，调用go help testflag，这些是go test过程中经常使用到的参数

//常用参数:
//-test.v : 是否输出全部的单元测试用例（不管成功或者失败），默认没有加上，所以只输出失败的单元测试用例。
//-test.run pattern: 只跑哪些单元测试用例
//-test.bench patten: 只跑那些性能测试用例
//-test.benchmem : 是否在性能测试的时候输出内存情况
//-test.benchtime t : 性能测试运行的时间，默认是1s
//-test.cpuprofile cpu.out : 是否输出cpu性能分析文件
//-test.memprofile mem.out : 是否输出内存性能分析文件
//-test.blockprofile block.out : 是否输出内部goroutine阻塞的性能分析文件
//-test.memprofilerate n : 内存性能分析的时候有一个分配了多少的时候才打点记录的问题。这个参数就是设置打点的内存分配间隔，也就是profile中一个sample代表的内存大小。默认是设置为512 * 1024的。如果你将它设置为1，则每分配一个内存块就会在profile中有个打点，那么生成的profile的sample就会非常多。如果你设置为0，那就是不做打点了。
//你可以通过设置memprofilerate=1和GOGC=off来关闭内存回收，并且对每个内存块的分配进行观察。
//-test.blockprofilerate n: 基本同上，控制的是goroutine阻塞时候打点的纳秒数。默认不设置就相当于-test.blockprofilerate=1，每一纳秒都打点记录一下
//-test.parallel n : 性能测试的程序并行cpu数，默认等于GOMAXPROCS。
//-test.timeout t : 如果测试用例运行时间超过t，则抛出panic
//-test.cpu 1,2,4 : 程序运行在哪些CPU上面，使用二进制的1所在位代表，和nginx的nginx_worker_cpu_affinity是一个道理
//-test.short : 将那些运行时间较长的测试用例运行时间缩短

//基准测试原理
//基准测试框架对一个测试用例的默认测试时间是 1 秒。开始测试时，当以 Benchmark 开头的基准测试用例函数返回时还不到 1 秒，那么 testing.B 中的 N 值将按 1、2、5、10、20、50……递增，同时以递增后的值重新调用基准测试用例函数。
//-benchtime参数可以自定义测试时间
//go test -v -bench=. -benchtime=5s benchmark_test.go

//-bench=.表示运行 benchmark_test.go 文件里的所有基准测试，和单元测试中的-run类似
//2000000000 表示测试的次数，也就是 testing.B 结构中提供给程序使用的 N。“0.33 ns/op”表示每一个操作耗费多少时间（纳秒）
//go test -v -bench=. benchmark_test.go
// b.N 由基准测试框架提供。
func Benchmark_Add(b *testing.B) {
	var n int
	for i := 0; i < b.N; i++ {
		n++
	}
}

//基准测试可以对一段代码可能存在的内存分配进行统计
//-benchmem参数以显示内存分配情况
//-bench后添加了 Alloc，指定只测试 Benchmark_Alloc() 函数
//“16 B/op”表示每一次调用需要分配 16 个字节，“2 allocs/op”表示每一次调用有两次分配。
//go test -v -bench=Alloc -benchmem benchmark_test.go
func Benchmark_Alloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("%d", i)
	}
}

//有些测试需要一定的启动和初始化时间，如果从 Benchmark() 函数开始计时会很大程度上影响测试结果的精准性。
// testing.B 提供了一系列的方法可以方便地控制计时器，从而让计时器只在需要的区间进行测试
//从 Benchmark() 函数开始，Timer 就开始计数。StopTimer() 可以停止这个计数过程，做一些耗时的操作，通过 StartTimer() 重新开始计时。ResetTimer() 可以重置计数器的数据
//计数器内部不仅包含耗时数据，还包括内存分配的数据。
//go test -v -bench=Add_TimerControl -benchmem benchmark_test.go
func Benchmark_Add_TimerControl(b *testing.B) {
	// 重置计时器
	b.ResetTimer()
	// 停止计时器
	b.StopTimer()
	// 开始计时器
	b.StartTimer()
	var n int
	for i := 0; i < b.N; i++ {
		n++
	}
}
