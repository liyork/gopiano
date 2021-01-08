package v1

import (
	"context"
	"golang.org/x/time/rate"
	"log"
	"os"
	"sort"
	"sync"
	"testing"
	"time"
)

// 粗细同时保证

func TestRateLimit(t *testing.T) {
	defer log.Printf("Done.")

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	apiConnection := Open()
	var wg sync.WaitGroup
	wg.Add(20)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ReadFile(context.Background())
			if err != nil {
				log.Printf("cannot ReadFile: %v", err)
			}
		}()

	}
	log.Printf("ReadFile")

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ResolveAddress(context.Background())
			if err != nil {
				log.Printf("cannot ResolveAddress: %v", err)

			}
		}()
	}
	log.Printf("ResolveAddress")

	wg.Wait()

	time.Sleep(10 * time.Second)
}

// 组合
type RateLimiter interface { //1
	Wait(context.Context) error
	Limit() rate.Limit
}

// 排序
func MultiLimiter(limiters ...RateLimiter) *multiLimiter {
	byLimit := func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	}
	sort.Slice(limiters, byLimit) //2
	// minuteLimit在前，因为他的Limit数少
	return &multiLimiter{limiters: limiters}

}

type multiLimiter struct {
	limiters []RateLimiter
}

// 先满足要求高的，即0.166
func (l *multiLimiter) Wait(ctx context.Context) error {
	for _, l := range l.limiters {
		if err := l.Wait(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (l *multiLimiter) Limit() rate.Limit {
	return l.limiters[0].Limit() //3
} //1.这里我们定义一个RateLimiter接口，以便MultiLimiter可以递归地定义其他MultiLimiter实例。
//2.这里我们实现一个优化，并按照每个RateLimiter的Limit()行排序。
//3.因为我们在multiLimiter实例化时对子RateLimiter实例进行排序，所以我们可以简单地返回限制性最高的limit，这将是切片中的第一个元素。

func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}

func Open() *APIConnection {
	// 1s有2个，2/s
	secondLimit := rate.NewLimiter(Per(2, time.Second), 1) //1
	// 1min有10个,60s有10个,10/60s=1/6s=0.166/s
	minuteLimit := rate.NewLimiter(Per(10, time.Minute), 10) //2
	return &APIConnection{
		rateLimiter: MultiLimiter(secondLimit, minuteLimit), //3
	}
}

type APIConnection struct {
	rateLimiter RateLimiter
}

func (a *APIConnection) ReadFile(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	// Pretend we do work here
	return nil
}

func (a *APIConnection) ResolveAddress(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	// Pretend we do work here
	return nil
} //1.我们定义了每秒的极限。
//2.我们定义每分钟的突发极限为10，为用户提供初始池。每秒的限制将确保我们不会因请求而使系统过载。
//3.我们结合这两个限制，并将其设置为APIConnection的主限制器。
//正如您所看到的，我们每秒发出两个请求，直到请求＃11，此时我们开始每隔六秒发出一次请求。 这是因为我们耗尽了我们可用的每分钟请求令牌池，并受此限制。
