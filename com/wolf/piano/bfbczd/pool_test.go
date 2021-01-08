package bfbczd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"sync"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	myPool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating new instance.")
			return struct{}{}
		},
	}

	myPool.Get()             //1
	instance := myPool.Get() //1
	myPool.Put(instance)     //2
	myPool.Get()             //3
	//1这里我们调用Get方法，将调用在池中定义的New函数，因为实例尚未实例化。
	//2在这里，我们将先前检索到的实例放回池中。 这时实例的可用数量为1个。
	//3执行此调用时，我们将重新使用先前分配的实例。New函数不会被调用。
}

func TestPool2(t *testing.T) {
	var numCalcsCreated int
	calcPool := &sync.Pool{
		New: func() interface{} {
			numCalcsCreated += 1
			mem := make([]byte, 1024) //1KB
			return &mem               // 1
		},
	}

	// 将池扩充到4KB
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	const numWorkers = 1024 * 1024
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := numWorkers; i > 0; i-- {
		go func() {
			defer wg.Done()

			mem := calcPool.Get().(*[]byte) // 2
			// 一旦加入随机，就表示协程不能取完马上换回去，也就表示其他协程要进行new
			//intn := rand.Intn(3)
			//time.Sleep(time.Duration(intn) * time.Millisecond)
			defer calcPool.Put(mem)

		}()
	}

	// 假设内存中执行了一些快速的操作

	wg.Wait()
	fmt.Printf("%d calculators were created.\n", numCalcsCreated)
}

// 1注意，我们存储了字节切片的指针。
// 2这里我们断言此类型是一个指向字节切片的指针。

func connectToService() interface{} {
	time.Sleep(1 * time.Second)
	return struct{}{}
}

func startNetworkDaemon() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		server, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		defer server.Close()
		wg.Done()
		for {
			conn, err := server.Accept()
			if err != nil {
				log.Printf("cannot accept connection: %v", err)
				continue
			}

			connectToService()
			fmt.Fprintln(conn, "")
			conn.Close()
		}
	}()

	return &wg
}
func init() {
	//daemonStarted := startNetworkDaemon()
	daemonStarted := startNetworkDaemon2()
	daemonStarted.Wait()
}

// go test -benchtime=10s -bench=. pool_test.go
func Benchmark_NetworkRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			b.Fatalf("cannot dial host: %v", err)
		}
		if _, err := ioutil.ReadAll(conn); err != nil {
			b.Fatalf("cannot read: %v", err)
		}
		conn.Close()
	}
}

func warmServiceConnCache() *sync.Pool {
	p := &sync.Pool{
		New: connectToService,
	}
	for i := 0; i < 10; i++ {
		p.Put(p.New())
	}
	return p
}

func startNetworkDaemon2() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		connPool := warmServiceConnCache()

		server, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		defer server.Close()
		wg.Done()
		for {
			conn, err := server.Accept()
			if err != nil {
				log.Printf("cannot accept connection: %v", err)
				continue

			}
			svcConn := connPool.Get()
			fmt.Fprintln(conn, "")
			connPool.Put(svcConn)
			conn.Close()
		}
	}()

	return &wg
}
