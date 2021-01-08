package base

import (
	"testing"
	"time"
)

// 并发子测试
func parallelTest1(t *testing.T) {
	t.Parallel()
	time.Sleep(3 * time.Second)
}

// 并发子测试
func parallelTest2(t *testing.T) {
	t.Parallel()
	time.Sleep(2 * time.Second)
}

// 并发子测试
func parallelTest3(t *testing.T) {
	t.Parallel()
	time.Sleep(1 * time.Second)
}

// TestSubParallel 通过把多个子测试放到一个组中并发执行，同时多个子测试可以共享setup和tear-down
// go test subparallel_test.go -v -run SubParallel
func TestSubParallel(t *testing.T) {
	t.Logf("Setup")

	t.Run("group", func(t *testing.T) {
		t.Run("Test1", parallelTest1)
		t.Run("Test2", parallelTest2)
		t.Run("Test3", parallelTest3)
	})

	// tear down
	t.Logf("teardown")
}

func TestSubParallel2(t *testing.T) {
	t.Logf("Setup2")

	t.Run("group", func(t *testing.T) {
		t.Run("Test12", parallelTest1)
		t.Run("Test22", parallelTest2)
		t.Run("Test32", parallelTest3)
	})

	// tear down
	t.Logf("teardown2")
}
