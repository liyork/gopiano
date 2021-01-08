package module

import (
	"fmt"
	"github.com/liyork/gopiano/com/wolf/piano/gopcp/ch6/webcrawler/errors"
	"net/http"
	"sync"
)

// 数据请求的类型
type Request struct {
	// HTTP请求
	httpReq *http.Request
	// 请求的深度
	depth uint32
}

func NewRequest(httpReq *http.Request, depth uint32) *Request {
	return &Request{httpReq: httpReq, depth: depth}
}

func (req *Request) HTTPReq() *http.Request {
	return req.httpReq
}

func (req *Request) Depth() uint32 {
	return req.depth
}

// 数据响应的类型
type Response struct {
	// HTTP响应
	httpResp *http.Response
	depth    uint32
}

func NewResponse(httpResp *http.Response, depth uint32) *Response {
	return &Response{httpResp: httpResp, depth: depth}
}

func (resp *Response) HTTPResp() *http.Response {
	return resp.httpResp
}

func (resp *Response) Depth() uint32 {
	return resp.depth
}

// 条目的类型
type Item map[string]interface{}

// 数据的接口类型
type Data interface {
	// 判断数据是否有效
	Valid() bool
}

func (req *Request) Valid() bool {
	return req.httpReq != nil && req.httpReq.URL != nil
}

func (resp *Response) Valid() bool {
	return resp.httpResp != nil && resp.httpResp.Body != nil
}

func (item *Item) Valid() bool {
	return item != nil
}

type Counts struct {
	CalledCount    uint64
	AcceptedCount  uint64
	CompletedCount uint64
	HandlingNumber uint64
}

// 代表组件的基础接口类型，实现必须是并发安全的。
type Module interface {
	// 当前组件的ID
	ID() MID
	// 获取当前组件的网络地址的字符串形式
	Addr()
	// 获取当前组件的评分
	Score() uint64
	SetScore(score uint64)
	// 获取评分计算器
	ScoreCalculator() CalculateScore
	// 获取当前组件被调用的计数
	CalledCount() uint64
	// 当前组件接受的调用的计数，组件一般会由于超负荷或参数有误而拒绝调用
	AcceptedCount() uint64
	// 获取当前组件已成功完成的调用计数
	CompletedCount() uint64
	// 获取当前组件正在处理的调用的数量
	HandlingNumber() uint64
	// 一次性获取所有计数
	Counts() Counts
	// 获取组件摘要
	Summary() SummaryStruct
}

type MID string

// 组件ID模板
var midTemplate = "%s%d|%s"

// 组件的类型
type Type string

// 组件类型常量
const (
	// 下载器
	TYPE_DOWNLOADER Type = "downloader"
	// 分析器
	TYPE_ANALYZER Type = "analyzer"
	// 条目处理管道
	TYPE_PIPELINE Type = "pipeline"
)

// 组件类型--字母的映射
var legalTypeLetterMap = map[Type]string{
	TYPE_DOWNLOADER: "D",
	TYPE_ANALYZER:   "A",
	TYPE_PIPELINE:   "P",
}

var legalLetterTypeMap = map[string]Type{}

// 序列号生成器的接口类型
type SNGenertor interface {
	// 获取预设的最小序列号
	Start() uint64
	// 获取预设的最大序列号
	Max() uint64
	// 获取下一个序列号
	Next() uint64
	// 获取循环可使用次数
	CycleCount() uint64
	// 获取一个序列号并准备下一个序列号
	Get() uint64
}

func NewSNGenertor(a, b int) SNGenertor {
	return nil
}

// 组件注册器的接口
type Registrar interface {
	Register(module Module) (bool, error)
	Unregister(mid MID) (bool, error)
	// 获取指定类型的组件的实例，基于负载均衡策略返回实例
	Get(moduleType Type) (Module, error)
	// 获取指定类型的所有组件实例
	GetAllByType(moduleType Type) (map[MID]Module, error)
	// 获取所有组件实例
	GetAll() map[MID]Module
	// 清除所有的组件注册记录
	Clear()
}

// 用于计算组件评分的函数类型
type CalculateScore func(counts Counts) uint64

// 组件摘要结构的类型
type SummaryStruct struct {
	ID        MID         `json:"id"`
	Called    uint64      `json:"called"`
	Accepted  uint64      `json:"accepted"`
	Completed uint64      `json:"completed"`
	Handling  uint64      `json:"handling"`
	Extra     interface{} `json:"extra,omitempty"`
}

// 下载器的接口类型，实现类型必须是并发安全的
type Downloader interface {
	Module
	Download(req *Request) (*Response, error)
}

// 分析器的接口类型，实现类型必须是并发安全的
type Analyzer interface {
	Module
	// 返回当前分析器使用的响应解析函数的列表
	RespParsers() []ParseResponse
	// 根据规则分析响应并返回请求和条目，响应需要分别经过若干响应解析函数的处理，然后合并结果
	Analyze(resp *Response) ([]Data, []error)
}

// 代表条目处理管道的接口类型，实现类必须是并发安全的
type Pipeline interface {
	Module
	// 返回当前条目处理管道使用的条目处理函数的列表
	ItemProcessors() []ProcessItem
	// 向条目处理管道发送条目，条目需要依次经过若干条目处理函数的处理
	Send(item Item) []error
	// 表示当前条目处理管道是否快速失败的，指：只要在处理某个条目时在某一个步骤上出错，那么条目处理管道就会忽略掉后续的所有处理步骤并报告错误
	FailFast() bool
	SetFailFast(failFast bool)
}

// 用于处理条目的函数的类型
type ProcessItem func(item Item) (result Item, err error)

// 组件注册器的实现类型
type myRegistrar struct {
	// 组件类型与对应组件实例的映射
	moduleTypeMap map[Type]map[MID]Module
	// 组件注册专用读写锁
	rwlock sync.RWMutex
}

func (registrar *myRegistrar) Unregister(mid MID) (bool, error) {
	panic("implement me")
}

func (registrar *myRegistrar) GetAllByType(moduleType Type) (map[MID]Module, error) {
	panic("implement me")
}

func (registrar *myRegistrar) GetAll() map[MID]Module {
	panic("implement me")
}

func (registrar *myRegistrar) Register(module Module) (bool, error) {
	if module == nil {
		return false, errors.NewIllegalParameterError("nil module instance")
	}
	mid := module.ID()
	parts, err := SplitMID(mid)
	if err != nil {
		return false, err
	}
	moduleType := legalLetterTypeMap[parts[0]]
	if !CheckType(moduleType, module) {
		errMsg := fmt.Sprintf("incorrect module type: %s", moduleType)
		return false, errors.NewIllegalParameterError(errMsg)
	}
	return true, nil
}

func CheckType(s Type, module Module) bool {
	return true
}

// 用于获取一个指定类型的组件的实例，基于负载均衡策略返回实例
func (registrar *myRegistrar) Get(moduleType Type) (Module, error) {
	modules, err := registrar.GetAllByType(moduleType)
	if err != nil {
		return nil, err
	}
	minScore := uint64(0)
	var selectedModule Module
	for _, module := range modules { // 找到评分最低
		err := SetScore(module)
		if err != nil {
			return nil, err
		}
		score := module.Score()
		if minScore == 0 || score < minScore {
			selectedModule = module
			minScore = score
		}
	}
	return selectedModule, nil
}

func SetScore(module Module) error {
	return nil
}

func (registrar *myRegistrar) Clear() {
}

func NewRegistrar() Registrar {
	return &myRegistrar{
		moduleTypeMap: map[Type]map[MID]Module{},
	}
}

func SplitMID(mid MID) ([]string, error) {
	return nil, nil
}

type ParseResponse func(httpResp *http.Response, respDepth uint32) ([]Data, []error)

func GetType(mid MID) (bool, Type) {
	return true, ""
}

func CalculateScoreSimple(counts Counts) uint64 {
	return counts.CalledCount + counts.AcceptedCount<<1 + counts.CompletedCount<<2 + counts.HandlingNumber<<4
}

func GenMID(t Type, a uint64, s string) (MID, error) {
	return "", nil
}
