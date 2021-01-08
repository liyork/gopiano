package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// 通过信号接收停止事件，用chan+Shutdown等待关闭完成再停止服务器
func main() {
	var addr = ":8080"

	// handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		fmt.Fprintln(w, "hello")
	})

	// server
	srv := http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// make sure idle connections returned
	processed := make(chan struct{})
	go func() {
		// 接收信号
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		log.Println("receive shutdown signal")

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		// Shutdown 方法会阻塞至所有连接闲置或 context 完成
		if err := srv.Shutdown(ctx); nil != err {
			log.Fatalf("server shutdown failed, err: %v\n", err)
		}
		log.Println("server gracefully shutdown")

		close(processed)
	}()

	//在 Shutdown 时会立刻返回
	err := srv.ListenAndServe()
	if http.ErrServerClosed != err {
		log.Fatalf("server not gracefully shutdown, err :%v\n", err)
	}

	// waiting for goroutine above processed
	<-processed
}
