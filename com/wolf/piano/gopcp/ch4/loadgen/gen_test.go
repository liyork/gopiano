package loadgen

import (
	loadgenlib "github.com/liyork/gopiano/com/wolf/piano/gopcp/ch4/loadgen/lib"
	"testing"
	"time"
)

var printDetail bool

// go test -v -run=TestStart
func TestStart(t *testing.T) {
	pset, server, _ := startServerAndGen(t)
	defer server.Close()

	total := 0
	countMap := make(map[loadgenlib.RetCode]int)
	for r := range pset.ResultCh {
		total++
		countMap[r.Code] = countMap[r.Code] + 1
		if printDetail {
			t.Logf("Result: ID=%d, Code=%d, Msg=%d, Elapse=%v.\n", r.ID, r.Code, r.Msg, r.Elapse)
		}
	}

	t.Logf("Total: %d.\n", total)
	successCount := countMap[loadgenlib.RET_CODE_SUCCESS]
	// pset.DurationNS/1e9得到s
	tps := float64(successCount) / float64(pset.DurationNS/1e9)
	// tps，被测软件平均每秒有效的处理(或称响应)载荷的数量
	t.Logf("Loads per second: %d; Treatments per second: %f.\n", pset.LPS, tps)
}

func startServerAndGen(t *testing.T) (ParamSet, TCPServer, loadgenlib.Generator) {
	server := NewTCPServer()

	serverAddr := "127.0.0.1:8080"
	t.Logf("Startup TCP server(%s)...\n", serverAddr)
	err := server.Listen(serverAddr)
	if err != nil {
		t.Fatalf("TCP Server startup failing! (addr=%s)!\n", serverAddr)
		t.FailNow()
	}
	pset := ParamSet{
		Caller:     NewTCPComm(serverAddr),
		TimeoutNS:  50 * time.Millisecond,
		LPS:        uint32(1000),
		DurationNS: 10 * time.Second,
		ResultCh:   make(chan *loadgenlib.CallResult, 50),
	}
	t.Logf("Initialize laod generator (timeoutNS=%v, lps=%d, durationNS=%v)...", pset.TimeoutNS, pset.LPS, pset.DurationNS)
	gen, err := NewGenerator(pset)
	if err != nil {
		t.Fatalf("Load generator initialization failing: %s\n", err)
		t.FailNow()
	}
	t.Log("Start load generator...")
	gen.Start()
	return pset, server, gen
}

// go test -v -run=TestStop
// 手动停止
func TestStop(t *testing.T) {
	_, server, gen := startServerAndGen(t)
	defer server.Close()

	timeoutNS := 2 * time.Second
	time.AfterFunc(timeoutNS, func() {
		gen.Stop()
	})
}

type TCPServer struct {
}

func (server TCPServer) Close() {

}

func (server TCPServer) Listen(s string) error {
	return nil
}

func NewTCPServer() TCPServer {
	return TCPServer{}
}

func NewTCPComm(addr string) loadgenlib.Caller {
	return &TCPComm{addr: addr}
}
