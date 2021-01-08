package monitor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/liyork/gopiano/com/wolf/piano/gopcp/ch6/webcrawler/examples/finder/lib"
	"github.com/liyork/gopiano/com/wolf/piano/gopcp/ch6/webcrawler/logger"
	sched "github.com/liyork/gopiano/com/wolf/piano/gopcp/ch6/webcrawler/scheduler"
	"runtime"
	"time"
)

func Monitor(
	scheduler sched.Scheduler,
	// 检查间隔时间，ns
	checkInterval time.Duration,
	// 摘要获取间隔时间，ns
	summarizeInterval time.Duration,
	maxIdleCount uint,
	// 是否在调度器空闲足够长的时间之后自行停止调度器
	autoStop bool,
	// 日志记录函数
	record lib.Record) <-chan uint64 {
	if scheduler == nil {
		panic(errors.New("The scheduler is invalid!"))
	}
	if checkInterval < time.Millisecond*100 {
		checkInterval = time.Millisecond * 100
	}
	if summarizeInterval < time.Second {
		summarizeInterval = time.Second
	}
	if maxIdleCount < 10 {
		maxIdleCount = 10
	}
	logger.Logger.Infof("Monitor parameters: checkInterval: %s, summarizeInterval: %s, maxIdleCount: %d, autoStop: %v",
		checkInterval, summarizeInterval, maxIdleCount, autoStop)
	stopNotifier, stopFunc := context.WithCancel(context.Background())
	// 不断记录scheduler中的错误
	reportError(scheduler, record, stopNotifier)
	recordSummary(scheduler, summarizeInterval, record, stopNotifier)
	checkCountChan := make(chan uint64, 2)
	checkStatus(scheduler, checkInterval, maxIdleCount, autoStop, checkCountChan, record, stopFunc, stopNotifier)
	return checkCountChan
}

func reportError(scheduler sched.Scheduler, record lib.Record, stopNotifier context.Context) {
	go func() {
		waitForSchedulerStart(scheduler)
		errorChan := scheduler.ErrorChan()
		for {
			select {
			case <-stopNotifier.Done():
				return
			default:
			}
			err, ok := <-errorChan
			if ok {
				errMsg := fmt.Sprintf("Received an error from error channel: %s", err)
				record(2, errMsg)
			}
			time.Sleep(time.Microsecond)
		}
	}()
}

func waitForSchedulerStart(scheduler sched.Scheduler) {

}

func recordSummary(scheduler sched.Scheduler, summarizeInterval time.Duration, record lib.Record, stopNotifier context.Context) {
	go func() {
		waitForSchedulerStart(scheduler)
		for {
			select {
			case <-stopNotifier.Done():
				return
			default:
			}

			// 用于区分两次前后数据不一致时才会记录
			var prevSchedSummaryStruct sched.SummaryStruct
			var prevNumGoroutine int
			// 记录的次数
			var recordCount uint64 = 1
			startTime := time.Now()

			currNumGoroutine := runtime.NumGoroutine()
			currSchedSummaryStruct := scheduler.Summary().Struct()
			if currNumGoroutine != prevNumGoroutine || !currSchedSummaryStruct.Same(prevSchedSummaryStruct) {
				summary := summary{
					NumGoroutine: runtime.NumGoroutine(),
					SchedSummary: currSchedSummaryStruct,
					EscapedTime:  time.Since(startTime).String(),
				}
				b, err := json.MarshalIndent(summary, "", "    ")
				if err != nil {
					logger.Logger.Errorf("Occur error when generate shceduler summary: %s\n", err)
					continue
				}
				msg := fmt.Sprintf("Monitory summary[%d]:\n%s", recordCount, b)
				record(0, msg)
				prevNumGoroutine = currNumGoroutine
				prevSchedSummaryStruct = currSchedSummaryStruct
				recordCount++
			}

			time.Sleep(time.Microsecond)
		}
	}()
}

// 监控结果摘要的结构
type summary struct {
	NumGoroutine int                 `json:"goroutine_number"`
	SchedSummary sched.SummaryStruct `json:"sched_summary"`
	EscapedTime  string              `json:"escaped_time"`
}

var msgReachMaxIdleCount = "The scheduler has been idle for a period of time (about %s). Consider to stop it now."
var msgStopScheduler = ""

func checkStatus(scheduler sched.Scheduler, checkInterval time.Duration, maxIdleCount uint, autoStop bool,
	checkCountChan chan<- uint64, record lib.Record, stopFunc context.CancelFunc, stopNotifier context.Context) {
	go func() {
		var checkCount uint64
		defer func() {
			stopFunc()
			checkCountChan <- checkCount
		}()

		waitForSchedulerStart(scheduler)
		for {
			select {
			case <-stopNotifier.Done():
				return
			default:
			}

			var idleCount uint
			var firstIdleTime time.Time

			if scheduler.Idle() {
				idleCount++
				if idleCount == 1 {
					firstIdleTime = time.Now()
				}
				if idleCount >= maxIdleCount {
					msg := fmt.Sprintf(msgReachMaxIdleCount, time.Since(firstIdleTime).String())
					record(0, msg)
					if scheduler.Idle() { // 再次检查
						if autoStop {
							var result string
							if err := scheduler.Stop(); err == nil {
								result = "success"
							} else {
								result = fmt.Sprintf("failing(%s)", err)
							}
							msg = fmt.Sprintf(msgStopScheduler, result)
							record(0, msg)
						}
						break
					} else {
						if idleCount > 0 {
							idleCount = 0
						}
					}
				}
			} else { // 重新计数
				if idleCount > 0 {
					idleCount = 0
				}
			}
			checkCount++
			time.Sleep(time.Microsecond)
		}
	}()
}
