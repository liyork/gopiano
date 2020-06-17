package routine

import (
	"fmt"
	"github.com/liyork/gopiano/com/wolf/piano/test/common"
	"runtime/debug"
	"testing"
	"time"
)

func testStackRoutine1() {
	fmt.Println("testStackRoutine1 ...")
	testStackRoutine2()
}
func testStackRoutine2() {
	fmt.Println("testStackRoutine1 ...")
	// To print the stack trace for the current goroutine
	debug.PrintStack()
}

func TestStackRoutine(t *testing.T) {
	testStackRoutine1()
}

// To print the stack trace for all goroutines
// Each Profile has a unique name. A few profiles are predefined:
//goroutine - stack traces of all current goroutines
//heap - a sampling of all heap allocations
//threadcreate - stack traces that led to the creation of new OS threads
//block - stack traces that led to blocking on synchronization primitives
func TestPrintAllRoutine(t *testing.T) {
	common.StartRoutine()

	fmt.Println("1111")
	time.Sleep(2222222 * time.Second)
	fmt.Println("2222")
}
