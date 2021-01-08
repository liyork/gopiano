package loadgen

import (
	"context"
	"errors"
	"fmt"
	"github.com/liyork/gopiano/com/wolf/piano/gopcp/ch4/loadgen/lib"
	"sync/atomic"
	"time"
)

// 用于表示载荷发生器的实现
type myGenerator struct {
	caller      lib.Caller           // 调用器,具体实现者
	timeoutNS   time.Duration        // 响应超时时间，ns
	lps         uint32               // 每秒载荷量(Loads per second)
	durationNS  time.Duration        // 负载持续时间，ns
	concurrency uint32               // 载荷并发量
	tickets     lib.GoTickets        // goroutine 票池
	ctx         context.Context      // 上下文
	cancelFunc  context.CancelFunc   // 取消函数
	callCount   int64                // 调用计数
	status      uint32               // 载荷发生器状态
	resultCh    chan *lib.CallResult // 调用结果通道
}

func (gen *myGenerator) Status() uint32 {
	return gen.status
}

func (gen *myGenerator) CallCount() int64 {
	return gen.callCount
}

// 把所有参数都内置到了ParamSet的结构体中。这样，若需变动NewGenerator函数的参数时，无需改变它的声明了，变动只会影响ParamSet类型。
type ParamSet struct {
	Caller     lib.Caller
	TimeoutNS  time.Duration
	LPS        uint32
	DurationNS time.Duration
	ResultCh   chan *lib.CallResult
}

func (set ParamSet) Check() error {
	return nil
}

// 新建一个载荷发生器
func NewGenerator(pset ParamSet) (lib.Generator, error) {
	logger.Infoln("New a load generator...")
	if err := pset.Check(); err != nil {
		return nil, err
	}
	gen := &myGenerator{
		caller:     pset.Caller,
		timeoutNS:  pset.TimeoutNS,
		lps:        pset.LPS,
		durationNS: pset.DurationNS,
		status:     lib.STATUS_ORIGINAL,
		resultCh:   pset.ResultCh,
	}
	if err := gen.init(); err != nil { // 初始化concurrency和tickets
		return nil, err
	}
	return gen, nil
}

func (gen *myGenerator) init() error {
	// 1e9/gen.lps表示发送频率
	gen.concurrency = uint32(int64(gen.timeoutNS)/int64(1e9/gen.lps) + 1)
	tickets, err := lib.NewGoTickets(gen.concurrency)
	if err != nil {
		return err
	}
	gen.tickets = tickets
	return nil
}

func (gen *myGenerator) Start() bool {
	// 设定节流阀
	var throttle <-chan time.Time
	if gen.lps > 0 {
		interval := time.Duration(1e9 / gen.lps)
		logger.Infof("Setting throttle (%v)...", interval)
		throttle = time.Tick(interval)
	}

	// 初始化上下文和取消函数
	gen.ctx, gen.cancelFunc = context.WithTimeout(context.Background(), gen.durationNS)

	// 初始化调用计数
	gen.callCount = 0
	// 设置状态为已启动
	atomic.StoreUint32(&gen.status, lib.STATUS_STARTED)

	go func() {
		// 生成载荷
		logger.Infoln("Generating loads...")
		gen.genLoad(throttle)
		logger.Infof("Stopped. (call count: %d)", gen.callCount)
	}()

	return true
}

func (gen *myGenerator) Stop() bool {
	if !atomic.CompareAndSwapUint32(&gen.status, lib.STATUS_STARTED, lib.STATUS_STOPPING) { //未成功设定STATUS_STOPPING则返回
		return false
	}
	// 让ctx字段发出停止“信号”
	gen.cancelFunc()
	for { // 不断检查状态变更。
		if atomic.LoadUint32(&gen.status) == lib.STATUS_STOPPED { // 若为已停止，表明prepareToStop方法执行完毕。
			break
		}
		time.Sleep(time.Microsecond)
	}
	return true
}

// 产生载荷并向承受方发送
func (gen *myGenerator) genLoad(throttle <-chan time.Time) {
	for {
		select {
		case <-gen.ctx.Done(): // 返回一个接收通道，在上下文超时或取消时关闭，这时针对他的接收操作立即返回。
			gen.prepareToStop(gen.ctx.Err())
			return
		default:
		}

		gen.asyncCall()
		if gen.lps > 0 { // 节流阀有效并需要使用
			select {
			case <-throttle: // 周期长短由节流阀控制
			case <-gen.ctx.Done():
				gen.prepareToStop(gen.ctx.Err())
				return
			}
		}
	}
}

// 用于为停止载荷发生器做准备
func (gen *myGenerator) prepareToStop(ctxError error) {
	logger.Infof("Prepare to stop load generator (cause: %s)...", ctxError)
	atomic.CompareAndSwapUint32(&gen.status, lib.STATUS_STARTED, lib.STATUS_STOPPING)
	logger.Infof("Closing result channel...")
	close(gen.resultCh)
	atomic.StoreUint32(&gen.status, lib.STATUS_STOPPED)
}

// 异步地调用承受方接口
func (gen *myGenerator) asyncCall() {
	gen.tickets.Take()
	go func() {
		defer func() {
			if p := recover(); p != nil {
				err, ok := interface{}(p).(error)
				var errMsg string
				if ok {
					errMsg = fmt.Sprintf("Async Call Panic! (error: %s)", err)
				} else {
					errMsg = fmt.Sprintf("Async Call Panic! (clue: %#v)", p)
				}
				logger.Errorln(errMsg)
				result := &lib.CallResult{
					ID:   -1,
					Code: lib.RET_CODE_FATAL_CALL,
					Msg:  errMsg,
				}
				gen.sendResult(result)
			}
			gen.tickets.Return()
		}()
		rawReq := gen.caller.BuildReq()

		// time包中的定时器可以用来设定某一个操作或任务的超时时间。要做到实时的超时判断，最好的方式就是与通道和select语句联用，
		// 不过就需要再启用一个goroutine来封装调用操作。如此，goroutine票池就收效甚微。
		//可以不用额外启用goroutine的情况下实现实时的超时判断
		//先声明代表调用状态的变量，并保证仅其上实施原子操作。
		// 调用状态：0-未调用或调用中；1-调用完成；2-调用超时
		var callStatus uint32
		timer := time.AfterFunc(gen.timeoutNS, func() { // 不过这个是由系统启动一个协成触发
			if !atomic.CompareAndSwapUint32(&callStatus, 0, 2) { // 若失败则说明载荷响应接收操作已完成，忽略超时处理
				return
			}

			result := &lib.CallResult{
				ID:     rawReq.ID,
				Req:    rawReq,
				Code:   lib.RET_CODE_WARNING_CALL_TIMEOUT,
				Msg:    fmt.Sprintf("Timeout! (expected: < %v)", gen.timeoutNS),
				Elapse: gen.timeoutNS,
			}
			gen.sendResult(result)
		})

		rawResp := gen.callOne(&rawReq)
		// 看谁先争抢到设定状态
		if !atomic.CompareAndSwapUint32(&callStatus, 0, 1) { // 若不成功，则说明调用已超时，
			return
		}
		// 开始响应处理
		timer.Stop()

		var result *lib.CallResult
		if rawResp.Err != nil {
			result = &lib.CallResult{
				ID:     rawResp.ID,
				Req:    rawReq,
				Code:   lib.RET_CODE_ERROR_CALL,
				Msg:    rawResp.Err.Error(),
				Elapse: rawResp.Elapse,
			}
		} else {
			result = gen.caller.CheckResp(rawReq, *rawResp)
			result.Elapse = rawResp.Elapse
		}
		gen.sendResult(result)
	}()
}

// 用于发送调用结果
func (gen *myGenerator) sendResult(result *lib.CallResult) bool {
	if atomic.LoadUint32(&gen.status) != lib.STATUS_STARTED { // 若状态不是已启动，就不发送了。
		gen.printIgnoredResult(result, "stopped load generator")
		return false
	}

	select {
	case gen.resultCh <- result:
		return true
	default: //防止chan已满而阻塞
		gen.printIgnoredResult(result, "full result channel")
		return false
	}
}

// 向载荷承受方发起一次调用
func (gen *myGenerator) callOne(rawReq *lib.RawReq) *lib.RawResp {
	atomic.AddInt64(&gen.callCount, 1)
	if rawReq == nil {
		return &lib.RawResp{ID: -1, Err: errors.New("Invalid raw request.")}
	}
	// 返回一个代表了当前时刻的纳秒数，从1970年1月1日的零时整开始算起。
	start := time.Now().UnixNano()
	resp, err := gen.caller.Call(rawReq.Req, gen.timeoutNS)
	end := time.Now().UnixNano()
	elapsedTime := time.Duration(end - start)
	var rawResp lib.RawResp
	if err != nil {
		errMsg := fmt.Sprintf("Sync Call Error: %s.", err)
		rawResp = lib.RawResp{
			ID:     rawReq.ID,
			Err:    errors.New(errMsg),
			Elapse: elapsedTime,
		}
	} else {
		rawResp = lib.RawResp{
			ID:     rawReq.ID,
			Resp:   resp,
			Elapse: elapsedTime,
		}
	}
	return &rawResp
}

func (gen *myGenerator) printIgnoredResult(result *lib.CallResult, s string) {

}
