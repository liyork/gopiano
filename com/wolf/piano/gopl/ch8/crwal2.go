package main

import (
	"fmt"
	"log"
	"os"
)

// 最多能放20个
var tokens = make(chan struct{}, 20)

func crawl2(url string) []string {
	fmt.Println(url)
	// 权限
	tokens <- struct{}{}
	list, err := Extract(url)
	<-tokens

	if err != nil {
		log.Print(err)
	}
	return list
}

// 解决终止问题，解决并发问题
func main() {
	worklist := make(chan []string)
	var n int

	// 启动1
	n++
	go func() { worklist <- os.Args[1:] }()

	seen := make(map[string]bool)
	for ; n > 0; n-- { // 有新的则n>0
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++ //每次新的+1
				go func(link string) {
					worklist <- crawl2(link)
				}(link)
			}
		}
	}
}
