package errors

import (
	"bytes"
	"fmt"
	"strings"
)

type ErrorType string

// 错误类型常量
const (
	// 下载器错误
	ERROR_TYPE_DOWNLOADER ErrorType = "downloader error"
	// 分析器错误
	ERROR_TYPE_ANALYZER ErrorType = "analyzer error"
	// 条目处理管道错误
	ERROR_TYPE_PIPELINE ErrorType = "pipeline error"
	// 调度器错误
	ERROR_TYPE_SCHEDULER ErrorType = "scheduler error"
)

// 爬虫错误的接口类型
type CrawlerError interface {
	// 获取错误的类型
	Type() ErrorType
	// 获取错误提示信息
	Error() string
}

// 爬虫错误的实现类型
type myCrawlerError struct {
	errType ErrorType
	// 错误的提示信息
	errMsg string
	// 完整的错误提示信息
	fullErrMsg string
}

func NewCrawlerError(errType ErrorType, errMsg string) CrawlerError {
	return &myCrawlerError{
		errType: errType,
		errMsg:  strings.TrimSpace(errMsg),
	}
}

func (ce *myCrawlerError) Type() ErrorType {
	return ce.errType
}

func (ce *myCrawlerError) Error() string {
	if ce.fullErrMsg == "" { // 为空才生成，避免频繁生成
		ce.genFullErrMsg()
	}
	return ce.fullErrMsg
}

func (ce *myCrawlerError) genFullErrMsg() {
	var buffer bytes.Buffer
	buffer.WriteString("crawler error: ")
	if ce.errType != "" {
		buffer.WriteString(string(ce.errType))
		buffer.WriteString(": ")
	}
	buffer.WriteString(ce.errMsg)
	ce.fullErrMsg = fmt.Sprintf("%s", buffer.String())
	return
}

func NewIllegalParameterError(msg string) error {
	return nil
}

func GenParameterError(s string) error {
	return nil
}
