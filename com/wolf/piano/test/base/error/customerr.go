package main

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"path/filepath"
	"runtime"
	"strings"
)

// 业务代码通用的错误
type ErrorCoder interface {
	Error() string
	Code() uint32
	Msg() string
	Where() string // 第一次生成这个错的地方, 第一次: 当newCoder和wrap一个非errorCoder的时候
}

// Grpc的错误
type GRPCStatuser interface {
	GRPCStatus() *status.Status
	Error() string
}

// 定义错误
type ErrorCode struct {
	code  uint32
	msg   string
	where string
}

// 错误，附带code
func (e *ErrorCode) Error() string {
	return fmt.Sprintf("code = %d ; msg = %s", e.code, e.msg)
}

// 不带code的错误消息
func (e *ErrorCode) Msg() string {
	return e.msg
}
func (e *ErrorCode) Code() uint32 {
	return e.code
}
func (e *ErrorCode) Where() string {
	return e.where
}

func New(msg string) *ErrorCode {
	where := caller(1, false)
	return &ErrorCode{code: 0, msg: msg, where: where}
}

// 构造一个错误
func NewCoder(code uint32, msg string, extMsg ...string) *ErrorCode {
	if len(extMsg) != 0 {
		msg = strings.Join(extMsg, " : ") + " : " + msg
	}
	where := caller(1, false)
	return &ErrorCode{code: code, msg: msg, where: where}
}

func NewCoderWhere(code uint32, callDepth int, msg string, extMsg ...string) *ErrorCode {
	if len(extMsg) != 0 {
		msg = strings.Join(extMsg, " : ") + " : " + msg
	}
	where := caller(callDepth, false)
	return &ErrorCode{code: code, msg: msg, where: where}
}

func NewCodere(code uint32, err error, extMsg ...string) *ErrorCode {
	var msg string
	if err != nil {
		msg = err.Error()
	}
	if len(extMsg) != 0 {
		msg = strings.Join(extMsg, " : ") + " : " + msg
	}
	where := caller(1, false)
	return &ErrorCode{code: code, msg: msg, where: where}
}

// 获取源代码行数
func caller(calldepth int, short bool) string {
	_, file, line, ok := runtime.Caller(calldepth + 1)
	if !ok {
		file = "???"
		line = 0
	} else if short {
		file = filepath.Base(file)
	}

	return fmt.Sprintf("%s:%d", file, line)
}

// Wrap 为error添加一个说明, 当这个err不确定是否应该报500或者是由其他服务返回时使用
// 如果err是ErrorCoder或者GRPCStatuser, code将继承, 否则code为0
func Wrap(err error, extMsg ...string) *ErrorCode {
	var msg string
	var code uint32
	var where string
	switch v := err.(type) {
	case ErrorCoder:
		msg = v.Msg()
		code = v.Code()
		where = v.Where()
	case GRPCStatuser:
		s := v.GRPCStatus()
		if s.Code() == codes.Unknown {
			code = 0
		} else if s.Code() < 20 {
			// 只要是grpc自带的错误就说明是系统错误
			code = 500
		} else {
			code = uint32(s.Code())
		}
		msg = s.Message()
		where = caller(1, false)
	default:
		msg = v.Error()
		code = 0
		where = caller(1, false)
	}
	if len(extMsg) != 0 {
		msg = strings.Join(extMsg, " : ") + " : " + msg
	}
	return &ErrorCode{code: code, msg: msg, where: where}
}
