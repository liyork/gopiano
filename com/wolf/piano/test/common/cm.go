package common

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

func test1() {
	fmt.Println("111111")
	test2()
}

func test2() {
	fmt.Println("222222")
}

func startGoRoutine1() {
	fmt.Println("startGoRoutine1 ...")
	go func() {
		for {
			time.Sleep(2 * time.Second)
			fmt.Println("111")
			test1()
		}
	}()
}

func startGoRoutine2() {
	fmt.Println("startGoRoutine2 ...")
	go func() {
		for {
			time.Sleep(2 * time.Second)
			fmt.Println("222")
		}
	}()
}

func StartRoutine() {
	startGoRoutine1()
	startGoRoutine2()
	time.Sleep(2 * time.Second)
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
	fmt.Println("333")
	//// http://localhost:6060/debug/pprof
	// func() {
	//	log.Println(http.ListenAndServe("localhost:8082", nil))
	//}()
}
