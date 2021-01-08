package v1

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

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

func Open() *APIConnection {
	return &APIConnection{
		// rate.Limit(1)--每秒event,b=1
		rateLimiter: rate.NewLimiter(rate.Limit(1), 1), //1
	}
}

type APIConnection struct {
	rateLimiter *rate.Limiter
}

func (a *APIConnection) ReadFile(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil { //2
		return err
	}
	fmt.Println("ReadFile...")
	// Pretend we do work here
	return nil
}

func (a *APIConnection) ResolveAddress(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil { //2
		return err
	}
	fmt.Println("ResolveAddress...")
	// Pretend we do work here
	return nil
} //1.我们将所有API连接的速率限制设置为每秒一个事件。
//2.我们等待速率限制器有足够的权限来完成我们的请求。
