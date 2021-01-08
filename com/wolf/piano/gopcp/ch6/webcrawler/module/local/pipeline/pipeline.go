package pipeline

import (
	"github.com/liyork/gopiano/com/wolf/piano/gopcp/ch6/webcrawler/module"
	"github.com/liyork/gopiano/com/wolf/piano/gopcp/ch6/webcrawler/module/stub"
)

// 条目处理管道的实现类型
type myPipeline struct {
	stub.ModuleInternal
	// 条目处理器的列表
	itemProcessors []module.ProcessItem
	failFast       bool
}

func New(mid module.MID, respParsers []module.ParseResponse, scoreCalculator module.CalculateScore) (module.Analyzer, error) {
	return nil, nil
}

func (pipeline *myPipeline) Send(item module.Item) []error {
	// ...
	var errs []error
	// ...
	var currentItem = item
	for _, processor := range pipeline.itemProcessors {
		processedItem, err := processor(currentItem)
		if err != nil {
			errs = append(errs, err)
			if pipeline.failFast {
				break
			}
		}
		if processedItem != nil {
			currentItem = processedItem
		}
	}
	//...
	return errs
}
