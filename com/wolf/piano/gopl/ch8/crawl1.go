package main

import (
	"fmt"
	"log"
	"os"
)

func crawl(url string) []string {
	fmt.Println(url)
	list, err := Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

func Extract(s string) ([]string, error) {
	return nil, nil
}

// go run crawl1.go http://gopl.io
func main() {
	worklist := make(chan []string)

	// 起始
	go func() {
		worklist <- os.Args[1:]
	}()

	seen := make(map[string]bool)
	// 主程
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				go func(link string) { //另一个协程，计算得到新的list放入worklist
					worklist <- crawl(link)
				}(link)
			}
		}
	}
}
