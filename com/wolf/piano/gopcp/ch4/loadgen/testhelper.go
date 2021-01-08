package loadgen

import (
	"encoding/json"
	"fmt"
	loadgenlib "github.com/liyork/gopiano/com/wolf/piano/gopcp/ch4/loadgen/lib"
	"math/rand"
	"net"
	"time"
)

// 用于表示服务器请求的结构
type ServerReq struct {
	ID       int64
	Operands []int
	// + - * /
	Operator string
}

// 表示服务器响应的结构
type ServerResp struct {
	ID int64 // 与对应的请求的ID一致
	// 运算式子，如2 + 4 + 5 = 11
	Formula string
	Result  int
	Err     error
}

// 表示TCP通信器的结构,TCP Communicator
type TCPComm struct {
	addr string
}

// 用于表示操作符切片
var operators = []string{"+", "-", "*", "/"}

// 构建一个请求
func (comm *TCPComm) BuildReq() loadgenlib.RawReq {
	id := time.Now().UnixNano() // 为了保持唯一性，使用纳秒级的时间戳
	sreq := ServerReq{
		ID: id,
		Operands: []int{
			int(rand.Int31n(1000) + 1), // [1,1000]
			int(rand.Int31n(1000) + 1),
		},
		Operator: func() string {
			return operators[rand.Int31n(100)%4]
		}(),
	}
	bytes, err := json.Marshal(sreq)
	if err != nil {
		panic(err)
	}
	rawReq := loadgenlib.RawReq{ID: id, Req: bytes}
	return rawReq
}

// 发起一次通信
func (comm *TCPComm) Call(req []byte, timeoutNS time.Duration) ([]byte, error) {
	conn, err := net.DialTimeout("tcp", comm.addr, timeoutNS)
	if err != nil {
		return nil, err
	}
	const DELIM = "\n"
	_, err = write(conn, req, DELIM)
	if err != nil {
		return nil, err
	}
	return read(conn, DELIM)
}

func read(conn net.Conn, s string) ([]byte, error) {
	return nil, nil
}

func write(conn net.Conn, bytes []byte, DELIM interface{}) (interface{}, error) {
	return nil, nil
}

// 检查响应内容
func (comm *TCPComm) CheckResp(rawReq loadgenlib.RawReq, rawResp loadgenlib.RawResp) *loadgenlib.CallResult {
	var commResult loadgenlib.CallResult
	commResult.ID = rawResp.ID
	commResult.Req = rawReq
	commResult.Resp = rawResp
	var sreq ServerReq
	err := json.Unmarshal(rawReq.Req, &sreq)
	if err != nil {
		commResult.Code = loadgenlib.RET_CODE_FATAL_CALL
		commResult.Msg = fmt.Sprintf("Incorrectly formatted Req: %s!\n", string(rawReq.Req))
		return &commResult
	}

	var sresp ServerResp
	err = json.Unmarshal(rawResp.Resp, &sresp)
	if err != nil {
		commResult.Code = loadgenlib.RET_CODE_ERROR_RESPONSE
		commResult.Msg = fmt.Sprintf("Incorrectly formatted Resp: %s!\n", string(rawResp.Resp))
		return &commResult
	}

	commResult.Code = loadgenlib.RET_CODE_SUCCESS
	commResult.Msg = fmt.Sprintf("Success. (%s)", sresp.Formula)

	return &commResult
}
