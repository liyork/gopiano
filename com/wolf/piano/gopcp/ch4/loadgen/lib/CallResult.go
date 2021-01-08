package lib

import (
	"errors"
	"fmt"
	"time"
)

// 用于表示调用结果的结构
type CallResult struct {
	ID     int64         // ID
	Req    RawReq        // 原生请求
	Resp   RawResp       // 原生响应
	Code   RetCode       // 响应代码
	Msg    string        // 结果成因的简述
	Elapse time.Duration // 耗时
}

// 用于表示原生请求的结构
type RawReq struct {
	ID  int64
	Req []byte
}

// 用于表示原生响应的结构
type RawResp struct {
	ID     int64
	Resp   []byte
	Err    error
	Elapse time.Duration
}

// 声明代表载荷发生器状态的常量
const (
	// 原值
	STATUS_ORIGINAL uint32 = 0
	// 正在启动
	STATUS_STARTING uint32 = 1
	// 已启动
	STATUS_STARTED uint32 = 2
	// 正在停止
	STATUS_STOPPING uint32 = 3
	// 已停止
	STATUS_STOPPED uint32 = 4
)

// 用于表示调用器的接口
type Caller interface {
	// 构建请求
	BuildReq() RawReq
	// 调用
	Call(req []byte, timeoutNS time.Duration) ([]byte, error)
	// 检查响应
	CheckResp(rawReq RawReq, rawResp RawResp) *CallResult
}

// 用于表示载荷发生器的接口
type Generator interface {
	// 启动，结果值代表是否已启动成功
	Start() bool
	// 停止，结果值代表是否已成功停止
	Stop() bool
	// 获取状态
	Status() uint32
	// 获取调用计数。每次启动会重置该计数
	CallCount() int64
}

// 用于表示goroutine票池的接口
type GoTickets interface {
	Take()
	Return()
	// 票池是否已被激活
	Active() bool
	Total() uint32
	Remainder() uint32
}

// 用于表示goroutine票池的实现
type myGoTickets struct {
	total    uint32        // 票的总数
	ticketCh chan struct{} // 承载goroutine票的容器
	active   bool          // 票池是否已被激活
}

func (gt *myGoTickets) Take() {
	panic("implement me")
}

func (gt *myGoTickets) Return() {
	panic("implement me")
}

func (gt *myGoTickets) Active() bool {
	panic("implement me")
}

func (gt *myGoTickets) Total() uint32 {
	panic("implement me")
}

func (gt *myGoTickets) Remainder() uint32 {
	panic("implement me")
}

func (gt *myGoTickets) init(total uint32) bool {
	if gt.active {
		return false
	}

	if total == 0 {
		return false
	}

	// 通道中缓冲的元素值的个数代表了还没有被获取和已被归还的goroutine票的总和。
	ch := make(chan struct{}, total)
	n := int(total)
	for i := 0; i < n; i++ {
		ch <- struct{}{}
	}
	gt.ticketCh = ch
	gt.total = total
	gt.active = true
	return true
}

// 用于新建一个goroutine票池
func NewGoTickets(total uint32) (GoTickets, error) {
	gt := myGoTickets{}
	if !gt.init(total) {
		errMsg := fmt.Sprintf("The goroutine ticket pool can not be initialized! (total=%d)\n", total)
		return nil, errors.New(errMsg)
	}
	return &gt, nil
}

type RetCode int

// 保留1~1000给载荷承受方使用
const (
	RET_CODE_SUCCESS              RetCode = 0    // 成功
	RET_CODE_WARNING_CALL_TIMEOUT         = 1001 // 调用超时警告
	RET_CODE_ERROR_CALL                   = 2001 // 调用错误
	RET_CODE_ERROR_RESPONSE               = 2002 // 响应内容错误
	RET_CODE_ERROR_CALLEE                 = 2003 // 被调用方(被测软件)的内部错误
	RET_CODE_FATAL_CALL                   = 3001 // 调用过程中发生了致命错误
)
