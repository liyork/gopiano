package bfbczd

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// 	看一个使用done通道模式的例子，并比较下切换到使用context文包获得什么好处。 这是一个同时打印问候和告别的程序：
func TestContext(t *testing.T) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background()) //1
	defer cancel()
	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := printGreeting(ctx); err != nil {
			fmt.Printf("cannot print greeting: %v\n", err)
			cancel() //2
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewell(ctx); err != nil {
			fmt.Printf("cannot print farewell: %v\n", err)
		}
	}()

	wg.Wait()
}

func printGreeting(ctx context.Context) error {
	greeting, err := genGreeting(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", greeting)
	return nil
}

func printFarewell(ctx context.Context) error {
	farewell, err := genFarewell(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", farewell)
	return nil
}

func genGreeting(ctx context.Context) (string, error) {
	// 构建自定义的Context.Context以满足其需求，而不必影响父级的Context
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second) //3
	defer cancel()

	switch locale, err := locale(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func genFarewell(ctx context.Context) (string, error) {
	switch locale, err := locale(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func locale(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err() //4
	case <-time.After(1 * time.Minute): // 1分钟
	}
	return "EN/US", nil
}

//1./在main函数中使用context.Background()建立个新的Context，并使用context.WithCancel将其包裹以便对其执行取消操作。
//2.在这一行上，如果从 printGreeting返回错误，main将取消context。
//3.这里genGreeting用context.WithTimeout包装Context。这将在1秒后自动取消返回的context，从而取消它传递context的子进程，即语言环境。
//4.这一行返回为什么Context被取消的原因。 这个错误会一直冒泡到main，这会导致注释2处的取消操作被调用。

// 进一步改进：因为我们知道locale需要大约一分钟的时间才能运行，所以可以在locale中检查是否给出了deadline
// 允许locale函数快速失败
func locale2(ctx context.Context) (string, error) {
	if deadline, ok := ctx.Deadline(); ok { //1
		if deadline.Sub(time.Now().Add(1*time.Minute)) <= 0 {
			return "", context.DeadlineExceeded
		}
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(1 * time.Minute):
	}
	return "EN/US", nil
}

//1.我们在这里检查是否Context提供了deadline。如果提供了，而且我们的程序时间已经越过这个时间线，这时简单的返回一个context包预设的错误——DeadlineExceeded。

func TestContextValue(t *testing.T) {
	ProcessRequest("jane", "abc123")
}

type ctxKey int

const (
	ctxUserID ctxKey = iota
	ctxAuthToken
)

func UserID(c context.Context) string {
	return c.Value(ctxUserID).(string)
}

func AuthToken(c context.Context) string {
	return c.Value(ctxAuthToken).(string)
}

func ProcessRequest(userID, authToken string) {
	ctx := context.WithValue(context.Background(), ctxUserID, userID)
	ctx = context.WithValue(ctx, ctxAuthToken, authToken)
	HandleResponse(ctx)
}

func HandleResponse(ctx context.Context) {
	fmt.Printf(
		"handling response for %v (auth: %v)",
		UserID(ctx),
		AuthToken(ctx),
	)
}
