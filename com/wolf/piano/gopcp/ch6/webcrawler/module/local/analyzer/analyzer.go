package analyzer

import (
	"fmt"
	"github.com/liyork/gopiano/com/wolf/piano/gopcp/ch6/webcrawler/buffer"
	"github.com/liyork/gopiano/com/wolf/piano/gopcp/ch6/webcrawler/errors"
	"github.com/liyork/gopiano/com/wolf/piano/gopcp/ch6/webcrawler/module"
	"github.com/liyork/gopiano/com/wolf/piano/gopcp/ch6/webcrawler/module/stub"
)

// 分析器的实现类型
type myAnalyzer struct {
	stub.ModuleInternal
	// 响应解析器列表
	respParsers []module.ParseResponse
}

func (analyzer *myAnalyzer) RespParsers() []module.ParseResponse {
	return analyzer.respParsers
}

func New(mid module.MID, respParsers []module.ParseResponse, scoreCalculator module.CalculateScore) (module.Analyzer, error) {
	moduleBase, err := stub.NewModuleInternal(mid, scoreCalculator)
	if err != nil {
		return nil, err
	}

	if respParsers == nil {
		return nil, errors.GenParameterError("nil response parsers")
	}
	if len(respParsers) == 0 {
		return nil, errors.GenParameterError("empty response parser list")
	}
	var innerParsers []module.ParseResponse
	for i, parser := range respParsers {
		if parser == nil {
			return nil, errors.GenParameterError(fmt.Sprintf("nil response parser[%d]", i))
		}
		innerParsers = append(innerParsers, parser)
	}
	return &myAnalyzer{
		ModuleInternal: moduleBase,
		respParsers:    innerParsers,
	}, nil
}

func (analyzer *myAnalyzer) Analyze(resp *module.Response) (dataList []module.Data, errList []error) {
	analyzer.ModuleInternal.IncrHandlingNumber()
	defer analyzer.ModuleInternal.DecrHandlingNumber()

	analyzer.ModuleInternal.IncrCalledCount()
	if resp == nil {
		errList = append(errList, errors.GenParameterError("nil response"))
		return
	}
	httpResp := resp.HTTPResp()
	if httpResp == nil {
		errList = append(errList, errors.GenParameterError("nil HTTP response"))
		return
	}
	httpReq := httpResp.Request
	if httpReq == nil {
		errList = append(errList, errors.GenParameterError("nil HTTP request"))
		return
	}
	var reqURL = httpReq.URL
	if reqURL == nil {
		errList = append(errList, errors.GenParameterError("nil HTTP request URL"))
		return
	}
	analyzer.ModuleInternal.IncrAcceptedCount()
	respDepth := resp.Depth()
	fmt.Printf("Parse the response (URL: %s, depth: %d)...\n", reqURL, respDepth)

	if httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	multipleReader, err := buffer.NewMultipleReader(httpResp.Body)
	if err != nil {
		errList = append(errList, genError(err.Error()))
		return
	}
	dataList = []module.Data{}
	for _, respParser := range analyzer.respParsers {
		httpResp.Body = multipleReader.Reader()
		pDataList, pErrorList := respParser(httpResp, respDepth)
		if pDataList != nil {
			for _, pData := range pDataList {
				if pData == nil {
					continue
				}
				dataList = appendDataList(dataList, pData, respDepth)
			}
		}
		if pErrorList != nil {
			for _, pError := range pErrorList {
				if pError == nil {
					continue
				}
				errList = append(errList, pError)
			}
		}
	}
	if len(errList) == 0 {
		analyzer.ModuleInternal.IncrCompletedCount()
	}
	return dataList, errList
}

func appendDataList(data []module.Data, data2 module.Data, u uint32) []module.Data {
	return nil
}

func genError(s string) error {
	return nil
}
