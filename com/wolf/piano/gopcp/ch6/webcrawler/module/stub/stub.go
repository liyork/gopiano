package stub

import (
	"fmt"
	"github.com/liyork/gopiano/com/wolf/piano/gopcp/ch6/webcrawler/errors"
	"github.com/liyork/gopiano/com/wolf/piano/gopcp/ch6/webcrawler/module"
)

// 组件的内部基础接口的类型
type ModuleInternal interface {
	module.Module
	// 调用计数+1
	IncrCalledCount()
	// 接受计数+1
	IncrAcceptedCount()
	// 成功完成数+1
	IncrCompletedCount()
	// 实时处理数+1
	IncrHandlingNumber()
	// 实时处理数-1
	DecrHandlingNumber()
	Clear()
}

// 组件内部基础接口的实现类型
type myModule struct {
	// 组件ID
	mid module.MID
	// 组件的网络地址
	addr string
	// 组件评分
	score uint64
	// 评分计数器
	scoreCalculator module.CalculateScore
	// 调用计数
	calledCount uint64
	// 接受计数
	acceptedCount uint64
	// 成功完成计数
	completedCount uint64
	// 实时处理数
	handlingNumber uint64
}

func (m *myModule) ID() module.MID {
	panic("implement me")
}

func (m *myModule) Addr() {
	panic("implement me")
}

func (m *myModule) Score() uint64 {
	panic("implement me")
}

func (m *myModule) SetScore(score uint64) {
	panic("implement me")
}

func (m *myModule) ScoreCalculator() module.CalculateScore {
	panic("implement me")
}

func (m *myModule) CalledCount() uint64 {
	panic("implement me")
}

func (m *myModule) AcceptedCount() uint64 {
	panic("implement me")
}

func (m *myModule) CompletedCount() uint64 {
	panic("implement me")
}

func (m *myModule) HandlingNumber() uint64 {
	panic("implement me")
}

func (m *myModule) Counts() module.Counts {
	panic("implement me")
}

func (m *myModule) Summary() module.SummaryStruct {
	panic("implement me")
}

func (m *myModule) IncrCalledCount() {
	panic("implement me")
}

func (m *myModule) IncrAcceptedCount() {
	panic("implement me")
}

func (m *myModule) IncrCompletedCount() {
	panic("implement me")
}

func (m *myModule) IncrHandlingNumber() {
	panic("implement me")
}

func (m *myModule) DecrHandlingNumber() {
	panic("implement me")
}

func (m *myModule) Clear() {
	panic("implement me")
}

func NewModuleInternal(mid module.MID, scoreCalculator module.CalculateScore) (ModuleInternal, error) {
	parts, err := module.SplitMID(mid)
	if err != nil {
		return nil, errors.NewIllegalParameterError(fmt.Sprintf("illegal ID %q: %s", mid, err))
	}
	return &myModule{
		mid:             mid,
		addr:            parts[2],
		scoreCalculator: scoreCalculator,
	}, nil
}
