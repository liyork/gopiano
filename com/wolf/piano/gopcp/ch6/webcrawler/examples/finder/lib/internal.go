package lib

import (
	"github.com/liyork/gopiano/com/wolf/piano/gopcp/ch6/webcrawler/module"
	"github.com/liyork/gopiano/com/wolf/piano/gopcp/ch6/webcrawler/module/local/downloader"
	"net"
	"net/http"
	"os"
	"path"
	"time"
)

func genHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   5,
			IdleConnTimeout:       60 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

func genResponseParsers() []module.ParseResponse {
	// 在HTML格式的HTTP响应体中查找新的请求
	parseLink := func(httpResp *http.Response, respDepth uint32) ([]module.Data, []error) {
		// 检查响应有效性
		// 检查响应Header中COntent-Type，必须text/html
		// 解析响应体,查找a标签，提取它的href属性，若是URL，将其封装成请求并追加到数据列表。再查找img标签，提取src属性，若是一个URL，将其封装成请求并追加到数据
		// 列表。若解析过程有错误，追加到错误列表。
		return nil, nil
	}
	parseImg := func(httpResp *http.Response, respDepth uint32) ([]module.Data, []error) {
		item := make(map[string]interface{})
		item["reader"] = httpResp.Body
		item["name"] = path.Base(httpResp.Request.URL.Path)
		//item["ext"]
		//pictureFormat
		var dataList []module.Data
		items := module.Item(item)
		dataList = append(dataList, &items)
		return dataList, nil
	}
	return []module.ParseResponse{parseLink, parseImg}
}

func genItemProcessors(dirPath string) []module.ProcessItem {
	var filePath string
	savePicture := func(item module.Item) (result module.Item, err error) {
		result = make(map[string]interface{})
		for k, v := range item {
			result[k] = v
		}
		result["file_path"] = filePath
		file, _ := os.Open(filePath)
		fileInfo, err := file.Stat()
		if err != nil {
			return nil, err
		}
		result["file_size"] = fileInfo.Size()
		return result, nil
	}
	// 将上述信息记录在日志中
	recordPicture := func(item module.Item) (result module.Item, err error) {
		return nil, nil
	}
	return []module.ProcessItem{savePicture, recordPicture}
}

// 组件序列号生成器
var snGen = module.NewSNGenertor(1, 0)

func GetDownloaders(number uint8) ([]module.Downloader, error) {
	var downloaders []module.Downloader
	if number == 0 {
		return downloaders, nil
	}

	for i := uint8(0); i < number; i++ {
		mid, err := module.GenMID(module.TYPE_DOWNLOADER, snGen.Get(), "")
		if err != nil {
			return downloaders, err
		}
		d, err := downloader.New(mid, genHTTPClient(), module.CalculateScoreSimple)
		if err != nil {
			return downloaders, err
		}
		downloaders = append(downloaders, d)
	}
	return downloaders, nil
}

// 用到analyzer包的New和本包genResponseParsers
func GetAnalyzers(number uint8) ([]module.Analyzer, error) { return nil, nil }

func GetPipelines(number uint8, dirPath string) ([]module.Pipeline, error) { return nil, nil }

// 日志记录函数的类型,level日志级别：0-普通；1-警告；2-错误
type Record func(level uint8, content string)
