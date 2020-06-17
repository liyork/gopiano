package concurrent

import (
	"fmt"
	"github.com/liyork/gopiano/com/wolf/piano/utils"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

//Go's concurrency primitives - goroutines and channels - provide an elegant and distinct means of structuring concurrent software.
//Instead of explicitly using locks to mediate access to shared data, Go encourages the use of channels to pass references to data between goroutines. This approach ensures that only one goroutine has access to the data at a given time.
//Do not communicate by sharing memory; instead, share memory by communicating.

const (
	numPollers     = 2                // number of Poller goroutines to launch
	pollInterval   = 3 * time.Second  // how often to poll each URL
	statusInterval = 10 * time.Second // how often to log status to stdout
	errTimeout     = 10 * time.Second // back-off timeout on error
)

var urls = []string{
	"http://www.google.com/",
	"http://golang.org/",
	"http://blog.golang.org/",
}

// State represents the last-known state of a URL.
type State struct {
	url    string
	status string
}

// StateMonitor maintains a map that stores the state of the URLs being
// polled, and prints the current state every updateInterval nanoseconds.
// It returns a chan State to which resource state should be sent.
func StateMonitor(updateInterval time.Duration) chan<- State {
	updates := make(chan State)
	urlStatus := make(map[string]string)
	ticker := time.NewTicker(updateInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				logState(urlStatus)
			case s := <-updates:
				urlStatus[s.url] = s.status
			}
		}
	}()
	return updates
}

// logState prints a state map.
func logState(s map[string]string) {
	log.Println("Current state:")
	for k, v := range s {
		log.Printf(" %s %s", k, v)
	}
}

// Resource represents an HTTP URL to be polled by this program.
type Resource struct {
	url      string
	errCount int
}

// Poll executes an HTTP HEAD request for url
// and returns the HTTP status string or an error string.
func (r *Resource) Poll() string {
	resp, err := http.Head(r.url)
	if err != nil {
		log.Println("Error", r.url, err)
		r.errCount++
		return err.Error()
	}
	r.errCount = 0
	return resp.Status
}

// Sleep sleeps for an appropriate interval (dependent on error state)
// before sending the Resource to done.
func (r *Resource) Sleep(done chan<- *Resource) {
	time.Sleep(pollInterval + errTimeout*time.Duration(r.errCount))
	done <- r
}

func Poller(in <-chan *Resource, out chan<- *Resource, status chan<- State) {
	for r := range in {
		s := r.Poll()
		status <- State{r.url, s}
		fmt.Printf("gid %v is running. \n", utils.GetGID())
		out <- r
	}
}

// 使用chan方式进行交换数据
// urls->pending->Poller->status
//                      ->complete，收到后sleep然后将Resource->pending
func TestUrlPoll(t *testing.T) {
	// Create our input and output channels.
	pending, complete := make(chan *Resource), make(chan *Resource)

	// Launch the StateMonitor.
	status := StateMonitor(statusInterval)

	// Launch some Poller goroutines.
	for i := 0; i < numPollers; i++ {
		go Poller(pending, complete, status)
	}

	// Send some Resources to the pending queue.
	go func() {
		for _, url := range urls {
			pending <- &Resource{url: url}
		}
	}()

	// 循环complete，开启go执行：等待+将Resource放入pending，重新出发上面Poller
	for r := range complete {
		go r.Sleep(pending)
	}
}

// 样例，可能共享内存使用锁的方式
type Resource2 struct {
	url        string
	polling    bool
	lastPolled int64
}

type Resources struct {
	data []*Resource2
	lock *sync.Mutex
}

func PollerInJava(res *Resources) {
	for {
		// get the least recently-polled Resource
		// and mark it as being polled
		res.lock.Lock()
		var r *Resource2
		for _, v := range res.data {
			if v.polling {
				continue
			}
			if r == nil || v.lastPolled < r.lastPolled {
				r = v
			}
		}
		if r != nil {
			r.polling = true
		}
		res.lock.Unlock()
		if r == nil {
			continue
		}

		// poll the URL

		// update the Resource's polling and lastPolled
		res.lock.Lock()
		r.polling = false
		r.lastPolled = time.Now().UnixNano()
		res.lock.Unlock()
	}
}
