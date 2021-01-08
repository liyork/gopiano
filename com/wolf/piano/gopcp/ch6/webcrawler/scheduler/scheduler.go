package scheduler

import (
	"context"
	syserrors "errors"
	"fmt"
	"github.com/liyork/gopiano/com/wolf/piano/gopcp/ch5/cmap"
	"github.com/liyork/gopiano/com/wolf/piano/gopcp/ch6/webcrawler/buffer"
	"github.com/liyork/gopiano/com/wolf/piano/gopcp/ch6/webcrawler/errors"
	"github.com/liyork/gopiano/com/wolf/piano/gopcp/ch6/webcrawler/logger"
	"github.com/liyork/gopiano/com/wolf/piano/gopcp/ch6/webcrawler/module"
	"net/http"
	"strings"
	"sync"
)

// 调度器的接口类型
type Scheduler interface {
	Init(requestArgs RequestArgs, dataArgs DataArgs, moduleArgs ModuleArgs) (err error)
	// 启动调度器并执行爬取流程，参数代表首次请求，调度器会以此为起始点开始执行
	Start(firstHTTPReq *http.Request) (err error)
	// 停止调度器的运行，所有处理模块执行的流程都被中止
	Stop() (err error)
	Status() Status
	// 获得错误通道。调度器以及各个处理模块运行过程中出现的所有错误都会被发送到该通道。若结果值为nil，说明通道不可用或调度器已停止
	ErrorChan() <-chan error
	// 判断所有处理模块是否都处于空闲状态
	Idle() bool
	// 获取摘要实例
	Summary() SchedSummary
}

// 请求相关的参数容器的类型
type RequestArgs struct {
	// 可以接受的URL的主域名的列表，不在主域名列表中的请求被忽略
	AcceptedDomains []string `json:"accepted_primary_domains"`
	// 需要爬取的最大深度。实际深度大于此值的请求被忽略
	MaxDepth uint32 `json:"max_depthh"`
}

func (args RequestArgs) Check() error {
	return nil
}

// 数据相关的参数容器的类型
type DataArgs struct {
	// 请求缓冲器的容量
	ReqBufferCap uint32 `json:"req_buffer_cap"`
	// 请求缓冲器的最大数量
	ReqMaxBufferNumber uint32 `json:"req_max_buffer_number"`
	// 响应缓冲器的容量
	RespBufferCap uint32 `json:"resp_buffer_cap"`
	// 响应缓冲器的最大数量
	RespMaxBufferNumber uint32 `json:"resp_max_buffer_number"`
	// 条目缓冲器的容量
	ItemBufferCap uint32 `json:"item_buffer_cap"`
	// 条目缓冲器的最大数量
	ItemMaxBufferNumber uint32 `json:"item_max_buffer_number"`
	// 错误缓冲器的容量
	ErrorBufferCap uint32 `json:"item_max_buffer_number"`
	// 错误缓冲器的最大数量
	ErrMaxBufferNumber uint32 `json:"error_max_buffer_number"`
}

func (args DataArgs) Check() error {
	return nil
}

// 组件相关的参数容器的类型
type ModuleArgs struct {
	Downloads []module.Downloader
	Analyzers []module.Analyzer
	Pipelines []module.Pipeline
}

func (args ModuleArgs) Check() error {
	return nil
}

// 参数容器的接口类型
type Args interface {
	// 用于自检参数的有效性,若结果为nil，说明未发现问题，否则意味着有错误
	Check() error
}

type Status uint8

const (
	//未初始化状态
	SCHED_STATUS_UNINITIALIZED Status = 0
	// 正在初始化的状态
	SCHED_STATUS_INITIALIZING Status = 1
	// 已初始化的状态
	SCHED_STATUS_INITIALIZED Status = 2
	// 正在启动的状态
	SCHED_STATUS_STARTING Status = 3
	// 已启动的状态
	SCHED_STATUS_STARTED Status = 4
	// 正在停止的状态
	SCHED_STATUS_STOPPING Status = 5
	// 已停止的状态
	SCHED_STATUS_STOPPED Status = 6
)

// 调度器摘要的接口类型
type SchedSummary interface {
	// 获取摘要信息的结构化形式
	Struct() SummaryStruct
	// 获取摘要信息的字符串形式
	String() string
}

type BufferPoolSummaryStruct struct {
}

type ModuleArgsSummary struct {
}

// 调度器摘要的结构
type SummaryStruct struct {
	RequestArgs     RequestArgs             `json:"request_args"`
	DataArgs        DataArgs                `json:"data_args"`
	ModuleArgs      ModuleArgsSummary       `json:"module_args"`
	Status          string                  `json:"status"`
	Downloaders     []module.SummaryStruct  `json:"downloaders"`
	Analyzers       []module.SummaryStruct  `json:"analyzers"`
	Pipelines       []module.SummaryStruct  `json:"pipelines"`
	ReqBufferPool   BufferPoolSummaryStruct `json:"request_buffer_pool"`
	RespBufferPool  BufferPoolSummaryStruct `json:"response_buffer_pool"`
	ItemBufferPool  BufferPoolSummaryStruct `json:"item_buffer_pool"`
	ErrorBufferPool BufferPoolSummaryStruct `json:"error_buffer_pool"`
	NumURL          uint64                  `json:"url_number"`
}

func (summaryStruct SummaryStruct) Same(summaryStruct2 SummaryStruct) bool {
	return true

}

// 调度器的实现类型
type myScheduler struct {
	// 爬取的最大深度，首次请求深度为0
	maxDepth uint32
	// 可接受的URL的主域名的字典
	acceptedDomainMap cmap.ConcurrentMap
	// 组件注册器
	registrar module.Registrar
	// 请求的缓冲池
	reqBufferPool buffer.Pool
	// 响应的缓冲池
	respBufferPool buffer.Pool
	// 条目的缓冲池
	itemBufferPool buffer.Pool
	// 错误的缓冲池
	errorBufferPool buffer.Pool
	// 已处理的URL的字典
	urlMap cmap.ConcurrentMap
	// 上下文，用于感知调度器的停止
	ctx context.Context
	// 取消函数，用于停止调度器
	cancelFunc context.CancelFunc
	// 状态
	status Status
	// 专用于状态的读写锁
	statusLock sync.RWMutex
	// 摘要
	summary SchedSummary
}

func (sched *myScheduler) Stop() (err error) {
	panic("implement me")
}

func (sched *myScheduler) Status() Status {
	panic("implement me")
}

func (sched *myScheduler) Idle() bool {
	panic("implement me")
}

func (sched *myScheduler) Summary() SchedSummary {
	panic("implement me")
}

func NewScheduler() *myScheduler {
	return &myScheduler{}
}

func (sched *myScheduler) Init(requestArgs RequestArgs, dataArgs DataArgs, moduleArgs ModuleArgs) (err error) {
	logger.Logger.Info("Check status for initialization...")
	var oldStatus Status
	oldStatus, err = sched.checkAndSetStatus(SCHED_STATUS_INITIALIZING)
	if err != nil {
		return
	}

	defer func() {
		sched.statusLock.Lock()
		if err != nil {
			sched.status = oldStatus
		} else {
			sched.status = SCHED_STATUS_INITIALIZED
		}
		sched.statusLock.Unlock()
	}()

	logger.Logger.Info("Check request arguments...")
	if err = requestArgs.Check(); err != nil {
		return err
	}

	logger.Logger.Info("Check data arguments...")
	if err = dataArgs.Check(); err != nil {
		return err
	}
	logger.Logger.Info("Data arguments are valid.")
	logger.Logger.Info("Check module arguments...")
	if err = moduleArgs.Check(); err != nil {
		return err
	}
	logger.Logger.Info("Module arguments are valid.")

	// ...初始化内部字段

	logger.Logger.Info("Register modules...")
	if err = sched.registerModules(moduleArgs); err != nil {
		return err
	}
	logger.Logger.Info("Scheduler has been initialized.")
	return nil
}

func (sched *myScheduler) checkAndSetStatus(wantedStatus Status) (oldStatus Status, err error) {
	sched.statusLock.Lock()
	defer sched.statusLock.Unlock()

	oldStatus = sched.status
	err = checkStatus(oldStatus, wantedStatus, nil)
	if err == nil {
		sched.status = wantedStatus
	}
	return
}

// 状态检查
// 规则：1.处于正在初始化、正在启动或正在停止状态时，不能从外部改变状态
// 2.想要的状态只能是正在初始化、正在启动或正在停止状态之一
// 3.处于未初始化状态时，不能变为正在启动或正在停止状态
// 4.处于已启动状态时，不能变为正在初始化或正在启动
// 5.只要未处于已启动状态，就不能变为正在停止状态
func checkStatus(currentStatus Status, wantedStatus Status, lock sync.Locker) (err error) {
	return nil
}

func (sched *myScheduler) Start(firstHTTPReq *http.Request) (err error) {
	defer func() {
		if p := recover(); p != nil {
			errMsg := fmt.Sprintf("Fatal scheduler error: %sched", p)
			logger.Logger.Fatal(errMsg)
			err = genError(errMsg)
		}
	}()
	logger.Logger.Info("Start scheduler...")
	logger.Logger.Info("Check status for start...")
	var oldStatus Status
	oldStatus, err = sched.checkAndSetStatus(SCHED_STATUS_STARTING)
	defer func() {
		sched.statusLock.Lock()
		if err != nil {
			sched.status = oldStatus
		} else {
			sched.status = SCHED_STATUS_STARTED
		}
		sched.statusLock.Unlock()
	}()
	if err != nil {
		return
	}

	// 检查参数
	logger.Logger.Info("Check first HTTP request...")
	if firstHTTPReq == nil {
		err = errors.GenParameterError("nil first HTTP request")
		return
	}
	logger.Logger.Info("The first HTTP request is valid.")
	logger.Logger.Info("Get the primary domain...")
	logger.Logger.Info("-- Host: %s", firstHTTPReq.Host)
	var primaryDomain string
	primaryDomain, err = getPrimaryDomain(firstHTTPReq.Host)
	if err != nil {
		return
	}
	logger.Logger.Info("-- Primay domain: %s", primaryDomain)
	sched.acceptedDomainMap.Put(primaryDomain, struct{}{})

	if err = sched.checkBufferPoolForStart(); err != nil {
		return
	}
	sched.download()
	sched.analyze()
	sched.pick()
	logger.Logger.Info("Scheduler has been started.")
	firstReq := module.NewRequest(firstHTTPReq, 0)
	sched.sendReq(firstReq)
	return nil
}

func genError(s string) error {
	return nil
}

func getPrimaryDomain(s string) (string, error) {
	return "", nil
}

func (sched *myScheduler) download() {
	go func() {
		for {
			if sched.canceled() {
				break
			}
			datum, err := sched.reqBufferPool.Get()
			if err != nil {
				logger.Logger.Warnln("The request buffer pool was closed. Break request reception.")
				break
			}
			req, ok := datum.(*module.Request)
			if !ok {
				errMsg := fmt.Sprintf("incorrect request type: %T", datum)
				sendError(syserrors.New(errMsg), "", sched.errorBufferPool)
			}
			sched.downloadOne(req)
		}
	}()
}

func (sched *myScheduler) canceled() bool {
	select {
	case <-sched.ctx.Done():
		return true
	default:
		return false
	}
}

func sendError(err error, mid module.MID, errorBufferPool buffer.Pool) bool {
	if err == nil || errorBufferPool == nil || errorBufferPool.Closed() {
		return false
	}
	var crawlerError errors.CrawlerError
	var ok bool
	crawlerError, ok = err.(errors.CrawlerError)
	if !ok {
		var moduleType module.Type
		var errorType errors.ErrorType
		ok, moduleType = module.GetType(mid)
		if !ok {
			errorType = errors.ERROR_TYPE_SCHEDULER
		} else {
			switch moduleType {
			case module.TYPE_DOWNLOADER:
				errorType = errors.ERROR_TYPE_DOWNLOADER
			case module.TYPE_ANALYZER:
				errorType = errors.ERROR_TYPE_ANALYZER
			case module.TYPE_PIPELINE:
				errorType = errors.ERROR_TYPE_PIPELINE
			}
		}
		crawlerError = errors.NewCrawlerError(errorType, err.Error())
	}
	if errorBufferPool.Closed() {
		return false
	}
	// 防止由于使用方没有处理过来errorBufferPool导致的阻塞。另开协成
	go func(crawlerError errors.CrawlerError) {
		if err := errorBufferPool.Put(crawlerError); err != nil {
			logger.Logger.Warnln("The error buffer pool was closed. Ignore error sending.")
		}
	}(crawlerError)
	return true
}

func (sched *myScheduler) downloadOne(req *module.Request) {
	if req == nil {
		return
	}
	if sched.canceled() {
		return
	}
	m, err := sched.registrar.Get(module.TYPE_DOWNLOADER)
	if err != nil || m == nil {
		errMsg := fmt.Sprintf("countn't get a downloader: %s", err)
		sendError(syserrors.New(errMsg), "", sched.errorBufferPool)
		// 重新放回
		sched.sendReq(req)
		return
	}

	downloader, ok := m.(module.Downloader)
	if !ok {
		errMsg := fmt.Sprintf("incorrect downloader type: %T (MID: %s)", m, m.ID())
		sendError(syserrors.New(errMsg), m.ID(), sched.errorBufferPool)
		// 重新放回
		sched.sendReq(req)
		return
	}
	resp, err := downloader.Download(req)
	if resp != nil {
		sendResp(resp, sched.respBufferPool)
	}
	if err != nil {
		sendError(err, m.ID(), sched.errorBufferPool)
	}
}

func sendResp(resp *module.Response, respBufferPool buffer.Pool) bool {
	if resp == nil || respBufferPool == nil || respBufferPool.Closed() {
		return false
	}
	go func(resp *module.Response) {
		if err := respBufferPool.Put(resp); err != nil {
			logger.Logger.Warnln("The response buffer pool was closed. Ignore response sending.")
		}
	}(resp)
	return true
}

func (sched *myScheduler) analyze() {
	go func() {
		for {
			if sched.canceld() {
				break
			}
			datum, err := sched.respBufferPool.Get()
			if err != nil {
				logger.Logger.Warnln("The response buffer pool was closed. Break response reception.")
				break
			}
			resp, ok := datum.(*module.Response)
			if !ok {
				errMsg := fmt.Sprintf("incorrect response type: %T", datum)
				sendError(syserrors.New(errMsg), "", sched.errorBufferPool)
			}
			sched.analyzeOne(resp)
		}
	}()
}

func (sched *myScheduler) analyzeOne(resp *module.Response) {
	if resp == nil {
		return
	}
	if sched.canceled() {
		return
	}
	m, err := sched.registrar.Get(module.TYPE_ANALYZER)
	if err != nil || m == nil {
		errMsg := fmt.Sprintf("couldn't get an analyzer: %s", err)
		sendError(syserrors.New(errMsg), "", sched.errorBufferPool)
		sendResp(resp, sched.respBufferPool)
		return
	}
	analyzer, ok := m.(module.Analyzer)
	if !ok {
		errMsg := fmt.Sprintf("incorrect analyzer type: %T (MID: %s)", m, m.ID())
		sendError(syserrors.New(errMsg), m.ID(), sched.errorBufferPool)
		sendResp(resp, sched.respBufferPool)
		return
	}
	dataList, errs := analyzer.Analyze(resp)
	if dataList != nil {
		for _, data := range dataList {
			if data == nil {
				continue
			}
			switch d := data.(type) {
			case *module.Request:
				sched.sendReq(d)
			case *module.Item:
				sendItem(d, sched.itemBufferPool)
			default:
				errMsg := fmt.Sprintf("Unsupported data type %T! (data: %#v)", d, d)
				sendError(syserrors.New(errMsg), m.ID(), sched.errorBufferPool)
			}
		}
	}
	if err != nil {
		for _, err := range errs {
			sendError(err, m.ID(), sched.errorBufferPool)
		}
	}
}

func sendItem(items *module.Item, pool buffer.Pool) {

}

// 向请求缓冲池发送请求，不符合要求的请求会被过滤掉
func (sched *myScheduler) sendReq(req *module.Request) bool {
	if req == nil {
		return false
	}
	if sched.canceled() {
		return false
	}
	httpReq := req.HTTPReq()
	if httpReq == nil {
		logger.Logger.Warnln("Ignore the request! Its HTTP request is invalid!")
		return false
	}
	reqURL := httpReq.URL
	if reqURL == nil {
		logger.Logger.Warnln("Ignore the request! Its URL is invalid!")
		return false
	}
	scheme := strings.ToLower(reqURL.Scheme)
	if scheme != "http" && scheme != "https" {
		logger.Logger.Warnf("Ignore the request! Its URL scheme is %q, but should be %q or %q. (URL: %s)\n", scheme, "http", "https", reqURL)
		return false
	}
	// 重复请求
	if v := sched.urlMap.Get(reqURL.String()); v != nil {
		logger.Logger.Warnf("Ignore the request! Its URL is repeated. (URL: %s)\n", reqURL)
		return false
	}
	pd, _ := getPrimaryDomain(httpReq.Host)
	if sched.acceptedDomainMap.Get(pd) == nil {
		if pd == "bing.net" {
			panic(httpReq.URL)
		}
		logger.Logger.Warnf("Ignore the request! Its host %q is not in accepted primary domain map. (URL: %s)\n", httpReq.Host, reqURL)
		return false
	}
	if req.Depth() > sched.maxDepth {
		logger.Logger.Warnf("Ignore the request! Its depth %d is greater than %d. (URL: %s)\n", req.Depth(), sched.maxDepth, reqURL)
		return false
	}
	go func(req interface{}) {
		if err := sched.reqBufferPool.Put(req); err != nil {
			logger.Logger.Warnln("The request buffer pool was closed. Ignore request sending.")
		}
	}(req)
	sched.urlMap.Put(reqURL.String(), struct{}{})
	return true
}

// 持续从errorBufferPool中获取error放入errCh中
func (sched *myScheduler) ErrorChan() <-chan error {
	errBuffer := sched.errorBufferPool
	errCh := make(chan error, errBuffer.BufferCap())
	go func(errBuffer buffer.Pool, errCh chan error) {
		for {
			if sched.canceled() {
				close(errCh)
				break
			}
			datum, err := errBuffer.Get()
			if err != nil {
				logger.Logger.Warnln("The error buffer pool was closed. Break error reception.")
				close(errCh)
				break
			}
			err, ok := datum.(error)
			if !ok {
				errMsg := fmt.Sprintf("incorrect error type: %T", datum)
				sendError(syserrors.New(errMsg), "", sched.errorBufferPool)
				continue
			}
			if sched.canceld() {
				close(errCh)
				break
			}
			errCh <- err
		}
	}(errBuffer, errCh)
	return errCh
}

func (sched *myScheduler) checkBufferPoolForStart() error {
	return nil
}

func (sched *myScheduler) canceld() bool {
	return true
}

func (sched *myScheduler) registerModules(args ModuleArgs) error {
	return nil
}

func (sched *myScheduler) pick() {

}
