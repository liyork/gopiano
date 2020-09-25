package concurrent

import (
	"sync"
	"testing"
)

var l sync.Mutex
var l1 sync.RWMutex
var a string

func flock() {
	a = "hello, world"
	l.Unlock()
}

// For any sync.Mutex or sync.RWMutex variable l and n < m,
// call n of l.Unlock() happens before call m of l.Lock() returns.
// unlock优先于lock
// is guaranteed to print "hello, world". The first call to l.Unlock() (in f) happens before the second call
// to l.Lock() (in main) returns, which happens before the print.
// For any call to l.RLock on a sync.RWMutex variable l,
// there is an n such that the l.RLock happens (returns) after call n to l.Unlock
// and the matching l.RUnlock happens before call n+1 to l.Lock.
func TestLockHappenBefore(t *testing.T) {
	l.Lock()
	go flock()
	l.Lock()
	println(a)

	l.Unlock() // n
	l1.RLock()
	l1.RUnlock()
	l1.Lock() // n+1
}
