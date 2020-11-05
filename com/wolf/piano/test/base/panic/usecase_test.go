package panic

import (
	"errors"
	"fmt"
	"log"
	"net"
	"testing"
	"time"
)

// Use Case 1: Avoid Panics Crashing Programs
// telnet localhost 12345
func TestUsecase1(t *testing.T) {
	listener, err := net.Listen("tcp", ":12345")
	if err != nil {
		log.Fatalln(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		// Handle each client connection
		// in a new goroutine.
		go ClientHandler(conn)
	}
}

func ClientHandler(c net.Conn) {
	// 若没有defer捕获panic则整个应用都挂了
	defer func() {
		if v := recover(); v != nil {
			log.Println("capture a panic:", v)
			log.Println("avoid crashing the program")
		}
		c.Close()
	}()
	panic(errors.New("just a demo.")) // a demo-purpose panic
}

// Use Case 2: Automatically Restart a Crashed Goroutine
func TestUsecase2(t *testing.T) {
	log.SetFlags(0)
	go NeverExit("job#A", shouldNotExit)
	go NeverExit("job#B", shouldNotExit)
	select {} // block here for ever
}

func shouldNotExit() {
	for {
		// Simulate a workload.
		time.Sleep(time.Second)

		// Simulate an unexpected panic.
		if time.Now().UnixNano()&0x3 == 0 {
			panic("unexpected situation")
		}
	}
}

func NeverExit(name string, f func()) {
	defer func() {
		if v := recover(); v != nil {
			// A panic is detected.
			log.Println(name, "is crashed. Restart it now.")
			go NeverExit(name, f) // restart
		}
	}()
	f()
}

// 跳出深层嵌套，不推荐
// Use Case 3: Use panic/recover Calls to Simulate Long Jump Statements
func TestUsecase3(t *testing.T) {
	n := func() (result int) {
		defer func() {
			if v := recover(); v != nil {
				if n, ok := v.(int); ok {
					result = n
				}
			}
		}()

		func() {
			func() {
				func() {
					// ...
					panic(123) // panic on succeeded
				}()
				// ...
			}()
		}()
		// ...
		return 0
	}()
	fmt.Println(n) // 123
}

// Use Case 4: Use panic/recover Calls to Reduce Error Checks
// 同样不推荐
/*func doSomething() (err error) {
	defer func() {
		err = recover()
	}()

	doStep1()
	doStep2()
	doStep3()
	doStep4()
	doStep5()

	return
}

// In reality, the prototypes of the doStepN functions
// might be different. For each of them,
// * panic with nil for success and no needs to continue.
// * panic with error for failure and no needs to contine.
// * not panic for continuing.
func doStepN() {
	...
	if err != nil {
		panic(err)
	}
	...
	if done {
		panic(nil)
	}
}

// The above code is less verbose than the following one.
func doSomething() (err error) {
	shouldContinue, err := doStep1()
	if !shouldContinue {
		return err
	}
	shouldContinue, err = doStep2()
	if !shouldContinue {
		return err
	}
	shouldContinue, err = doStep3()
	if !shouldContinue {
		return err
	}
	shouldContinue, err = doStep4()
	if !shouldContinue {
		return err
	}
	shouldContinue, err = doStep5()
	if !shouldContinue {
		return err
	}

	return
}

// If err is not nil, then shouldContinue must be true.
// If shouldContinue is true, err might be nil or non-nil.
func doStepN() (shouldContinue bool, err error) {
	...
	if err != nil {
		return false, err
	}
	...
	if done {
		return false, nil
	}
	return true, nil
}*/
