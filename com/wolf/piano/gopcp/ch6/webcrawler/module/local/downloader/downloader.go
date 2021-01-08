package downloader

import (
	"fmt"
	"github.com/liyork/gopiano/com/wolf/piano/gopcp/ch6/webcrawler/module"
	"github.com/liyork/gopiano/com/wolf/piano/gopcp/ch6/webcrawler/module/stub"
	"net/http"
)

// 下载器的实现类型
type myDownloader struct {
	// 组件基础实例
	stub.ModuleInternal
	// 下载用的HTTP客户端
	httpClient http.Client
}

func New(mid module.MID, client *http.Client, scoreCalculator module.CalculateScore) (module.Downloader, error) {
	moduleBase, err := stub.NewModuleInternal(mid, scoreCalculator)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, genParameterError("nil http client")
	}
	return &myDownloader{
		ModuleInternal: moduleBase,
		httpClient:     *client,
	}, nil
}

func (downloader *myDownloader) Download(req *module.Request) (*module.Response, error) {
	downloader.ModuleInternal.IncrHandlingNumber() // 实时处理数
	defer downloader.ModuleInternal.DecrHandlingNumber()

	downloader.ModuleInternal.IncrCalledCount() // 调用数
	if req == nil {
		return nil, genParameterError("nil request")
	}

	httpReq := req.HTTPReq()
	if httpReq == nil {
		return nil, genParameterError("nil HTTP request")
	}
	downloader.ModuleInternal.IncrAcceptedCount() // 接受计数
	fmt.Printf("Do the request (URL: %s, depth: %d)...\n", httpReq.URL, req.Depth())
	httpResp, err := downloader.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	downloader.ModuleInternal.IncrCompletedCount() // 完成计数
	return module.NewResponse(httpResp, req.Depth()), nil
}

func genParameterError(s string) error {
	return nil
}
