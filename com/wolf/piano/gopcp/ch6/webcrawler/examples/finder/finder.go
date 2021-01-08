package main

import (
	"flag"
	"fmt"
	"github.com/liyork/gopiano/com/wolf/piano/gopcp/ch6/webcrawler/examples/finder/lib"
	"github.com/liyork/gopiano/com/wolf/piano/gopcp/ch6/webcrawler/logger"
	"github.com/liyork/gopiano/com/wolf/piano/gopcp/ch6/webcrawler/monitor"
	sched "github.com/liyork/gopiano/com/wolf/piano/gopcp/ch6/webcrawler/scheduler"
	"net/http"
	"os"
	"strings"
	"time"
)

// 命令参数
var (
	firstURL string
	domains  string
	depth    uint
	dirPath  string
)

func init() {
	flag.StringVar(&firstURL, "first", "http://zhihu.sogou.com/zhihu?query=golang+logo", "The first URL which you want to access.")
	flag.StringVar(&domains, "domains", "zhihu.com", "The primary domains which you accepted.Pleader using comma-separated multiple domains.")
	flag.UintVar(&depth, "depth", 3, "The depth for crawling.")
	flag.StringVar(&dirPath, "dir", "./pictures", "The path which you want to save the image files.")
}

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s: \n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\tfinder [flag] \n")
	fmt.Fprintf(os.Stderr, "Flags \n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = Usage
	flag.Parse()

	scheduler := sched.NewScheduler()
	domainParts := strings.Split(domains, ",")
	var acceptedDomains []string
	for _, domain := range domainParts {
		domain = strings.TrimSpace(domain)
		if domain != "" {
			acceptedDomains = append(acceptedDomains, domain)
		}
	}
	requestArgs := sched.RequestArgs{
		AcceptedDomains: acceptedDomains,
		MaxDepth:        uint32(depth),
	}
	dataArgs := sched.DataArgs{
		ReqBufferCap:        50,
		ReqMaxBufferNumber:  1000,
		RespBufferCap:       50,
		RespMaxBufferNumber: 10,
		ItemBufferCap:       50,
		ItemMaxBufferNumber: 100,
		ErrorBufferCap:      50,
		ErrMaxBufferNumber:  1,
	}
	downloaders, err := lib.GetDownloaders(1)
	if err != nil {
		logger.Logger.Fatalf("An error occurs when creating downloaders: %s", err)
	}
	analyzers, err := lib.GetAnalyzers(1)
	if err != nil {
		logger.Logger.Fatalf("An error occurs when creating analyzer: %s", err)
	}
	pipelines, err := lib.GetPipelines(1, dirPath)
	if err != nil {
		logger.Logger.Fatalf("An error occurs when creating pipelines: %s", err)
	}
	moduleArgs := sched.ModuleArgs{
		Downloads: downloaders,
		Analyzers: analyzers,
		Pipelines: pipelines,
	}
	err = scheduler.Init(requestArgs, dataArgs, moduleArgs)
	if err != nil {
		logger.Logger.Fatalf("An error occurs when initializing scheduler: %s", err)
	}

	checkInterval := time.Second
	summarizeInterval := 100 * time.Millisecond
	maxIdleCount := uint(5)
	checkCountChan := monitor.Monitor(scheduler, checkInterval, summarizeInterval, maxIdleCount, true, Record)

	firstHTTPReq, err := http.NewRequest("GET", firstURL, nil)
	if err != nil {
		logger.Logger.Fatalln(err)
		return
	}
	err = scheduler.Start(firstHTTPReq)
	if err != nil {
		logger.Logger.Fatalf("An error occurs when starting scheduler: %s", err)
	}
	// 等待监控结束，就是idle之时
	<-checkCountChan
}

func Record(level byte, content string) {
	if content == "" {
		return
	}
	switch level {
	case 0:
		logger.Logger.Infoln(content)
	case 1:
		logger.Logger.Warnln(content)
	case 2:
		logger.Logger.Errorln(content)
	}
}
