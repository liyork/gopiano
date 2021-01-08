package main

import (
	"fmt"
	"os"
	"testing"
)

var ERR error

// go test -bench=. writingBU_test.go writingBU.go
func benchmarkCreate(b *testing.B, buffer, filesize int) {
	var err error
	for i := 0; i < b.N; i++ {
		err = Create("/tmp/random", buffer, filesize)
	}
	ERR = err
	err = os.Remove("/tmp/random")
	if err != nil {
		fmt.Println(err)
	}
}

func Benchmark1Create(b *testing.B) {
	benchmarkCreate(b, 1, 1000000)
}

func Benchmark2Create(b *testing.B) {
	benchmarkCreate(b, 2, 1000000)
}

func Benchmark4Create(b *testing.B) {
	benchmarkCreate(b, 4, 1000000)
}

func Benchmark10Create(b *testing.B) {
	benchmarkCreate(b, 10, 1000000)
}

func Benchmark1000Create(b *testing.B) {
	benchmarkCreate(b, 1000, 1000000)
}

// 使用一个大小为 1 个字节的写缓冲是完全无效的并且缓冲所有的操作。另外，这样的缓冲大小需要更多的内存操作，这也使程序运行的更慢！
// 使用 2 个字节的缓冲可以整个程序速度提升 2 倍，这是好事。然而，这仍然很慢。这同样适用于 4 个字节的写缓冲。
// 当决定用 10 个字节的写缓冲时，这会变的更快。最后，这个结果显示使用 1,000 字节的写缓冲没有比使用 10 字节的快 100 倍，这意味着在速度和写缓冲大小之间的最佳点是在这俩个值之间。
