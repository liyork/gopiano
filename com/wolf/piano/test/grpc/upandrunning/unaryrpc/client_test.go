package unaryrpc

import (
	"context"
	"fmt"
	pb "github.com/liyork/gopiano/com/wolf/piano/test/grpc/upandrunning/unaryrpc/ecommerce"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	_ "net/http/pprof"
	"testing"
	"time"
)

// grpc的ClientConn对象可以帮我们实现自动重连的机制，并且是并发安全的，因此可以定义一个全局的ClientConn

const (
	address = "10.0.12.76:22222"
)

func TestClientBase(t *testing.T) {

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithUnaryInterceptor(orderUnaryClientInterceptor))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	orderMgtClient := pb.NewOrderManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Get Order
	retrievedOrder, err := orderMgtClient.GetOrder(ctx,
		&pb.OrderRequest{Value: "106"})
	log.Printf("GetOrder Response -> :%v, err:%v \n", retrievedOrder, err)
	time.Sleep(1111 * time.Second)
}

// 在 gRPC 中强调 TL;DR（Too long, Don't read）并建议始终设定截止日期
// 当未设置 Deadlines 时，将采用默认的 DEADLINE_EXCEEDED（这个时间非常大）
// 如果产生了阻塞等待，就会造成大量正在进行的请求都会被保留，并且所有请求都有可能达到最大超时
//这会使服务面临资源耗尽的风险，例如内存，这会增加服务的延迟，或者在最坏的情况下可能导致整个进程崩溃
func TestClientDeadline(t *testing.T) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrderManagementClient(conn)

	clientDeadline := time.Now().Add(1 * time.Second)
	// 第一个形参为父上下文，第二个形参为调整的截止时间。父子时间以最早为准
	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)

	defer cancel()

	// Add Order
	order1 := pb.OrderRequest{Value: "101"}
	//If the invocation exceeds the specified deadline, it should return an error of the type DEADLINE_EXCEEDED
	res, addErr := client.GetOrder(ctx, &order1)

	if addErr != nil {
		statusErr, ok := status.FromError(addErr)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Fatalln("client.Search err: deadline")
			}
		}
		log.Fatalf("client.Search err: %v", addErr)
		//got := status.Code(addErr)
	} else {
		log.Print("AddOrder Response -> ", res.Id)
	}
}

func TestClientProcessErr(t *testing.T) {

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithUnaryInterceptor(orderUnaryClientInterceptor))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	orderMgtClient := pb.NewOrderManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// Get Order
	retrievedOrder, err := orderMgtClient.GetOrder(ctx,
		&pb.OrderRequest{Value: "106"})
	if err != nil {
		errorCode := status.Code(err)
		if errorCode == codes.InvalidArgument {
			log.Printf("Invalid Argument Error : %s", errorCode)
			errorStatus := status.Convert(err)
			for _, d := range errorStatus.Details() {
				switch info := d.(type) {
				case *errdetails.BadRequest_FieldViolation:
					log.Printf("Request Field Invalid: %s", info)
				default:
					log.Printf("Unexpected error type: %s", info)
				}
			}
		} else {
			log.Printf("Unhandled error : %s ", errorCode)
		}
	} else {
		log.Print("AddOrder Response -> ", retrievedOrder.Id)
	}
}

// 重用连接，对两个svc
func TestClientMultiplex(t *testing.T) {

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithUnaryInterceptor(orderUnaryClientInterceptor))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	orderMgtClient := pb.NewOrderManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// Get Order
	retrievedOrder, err := orderMgtClient.GetOrder(ctx,
		&pb.OrderRequest{Value: "106"})
	log.Print("GetOrder Response -> : ", retrievedOrder)

	helloClient := pb.NewHelloServerClient(conn)
	helloClient.Hello(ctx, &pb.HelloRequest{Id: "111"})
}

func TestClientMetadata(t *testing.T) {

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithUnaryInterceptor(orderUnaryClientInterceptor))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	orderMgtClient := pb.NewOrderManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 1
	//md := metadata.Pairs(
	//	"timestamp", time.Now().Format(time.StampNano),
	//	"kn", "vn",
	//)
	//mdCtx := metadata.NewOutgoingContext(context.Background(), md)

	// 2
	ctxA := metadata.AppendToOutgoingContext(ctx,
		"k1", "v1", "k1", "v2", "k2", "v3")

	var header, trailer metadata.MD
	retrievedOrder, err := orderMgtClient.GetOrder(ctxA,
		&pb.OrderRequest{Value: "106"},
		grpc.Header(&header),
		grpc.Trailer(&trailer))
	log.Print("GetOrder Response -> : ", retrievedOrder, header, trailer)
}

const (
	exampleScheme      = "exampleScheme"
	exampleServiceName = "exampleServiceName"
)

func TestClientLb(t *testing.T) {

	pickfirstConn, err := grpc.Dial(
		fmt.Sprintf("%s:///%s",
			// 	exampleScheme      = "example"
			//	exampleServiceName = "lb.example.grpc.io"
			exampleScheme, exampleServiceName),
		// "pick_first" is the default option.
		grpc.WithBalancerName("pick_first"),
		grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer pickfirstConn.Close()

	log.Println("==== Calling helloworld.Greeter/SayHello " +
		"with pick_first ====")
	makeRPCs(pickfirstConn, 10)

	// Make another ClientConn with round_robin policy.
	roundrobinConn, err := grpc.Dial(
		fmt.Sprintf("%s:///%s", exampleScheme, exampleServiceName),
		// "example:///lb.example.grpc.io"
		grpc.WithBalancerName("round_robin"),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer roundrobinConn.Close()

	log.Println("==== Calling helloworld.Greeter/SayHello " +
		"with round_robin ====")
	makeRPCs(roundrobinConn, 10)
}

func makeRPCs(conn *grpc.ClientConn, i int) {

}

func TestClientCompression(t *testing.T) {
	//client.AddOrder(ctx, &order1, grpc.UseCompressor(gzip.Name))
}

// log
//GRPC_GO_LOG_VERBOSITY_LEVEL=99
//GRPC_GO_LOG_SEVERITY_LEVEL=info

// 连接不上，则3s超时返回
// 链接上，会在3s内进行指定次数重试，一定要符合RetryableStatusCodes，也就是服务端有响应且err对应上code
// 对于连不上，通过wireshark看，确实会有重试syn(条件:ip.src==10.0.12.76)，而且服务端在超时范围内启动确实会有成功
func TestWithRetry(t *testing.T) {
	//pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
	// http://localhost:6060/debug/pprof
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	//os.Setenv("GRPC_GO_RETRY", "on")

	retryPolicy := `{
	   "methodConfig": [{
	     "name": [{"service": "ecommerce.OrderManagement"}],
	     "waitForReady": true,
	     "retryPolicy": {
	             "MaxAttempts": 4,
	             "InitialBackoff": ".01s",
	             "MaxBackoff": ".01s",
	             "BackoffMultiplier": 1.0,
	             "RetryableStatusCodes": [ "UNAVAILABLE","DATA_LOSS" ]
	     }
	   }]}`

	conn, err := grpc.Dial("10.0.12.76:22222", grpc.WithInsecure(), grpc.WithDefaultServiceConfig(retryPolicy))
	//conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//defer conn.Close()

	client := pb.NewOrderManagementClient(conn)

	//clientDeadline := time.Now().Add(
	//	time.Duration(2 * time.Second))
	//ctx, cancel := context.WithDeadline(
	//	context.Background(), clientDeadline)

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	//ctx, cancel := context.WithCancel(context.Background());
	//defer cancel()

	// Add Order
	order1 := pb.OrderRequest{Value: "101"}
	//If the invocation exceeds the specified deadline, it should return an error of the type DEADLINE_EXCEEDED
	res, addErr := client.GetOrder(ctx, &order1)

	if addErr != nil {
		got := status.Code(addErr)
		log.Printf("Error Occured -> addOrder : , %v:", got)
	} else {
		log.Print("AddOrder Response -> ", res.Id)
	}
	time.Sleep(1111 * time.Second)
}

func TestResolve(t *testing.T) {

	//r := lb.NewResolver(*serv)
	//b := grpc.RoundRobin(r)

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithUnaryInterceptor(orderUnaryClientInterceptor))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	orderMgtClient := pb.NewOrderManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Get Order
	retrievedOrder, err := orderMgtClient.GetOrder(ctx,
		&pb.OrderRequest{Value: "106"})
	log.Printf("GetOrder Response -> :%v, err:%v \n", retrievedOrder, err)
	time.Sleep(1111 * time.Second)
}

func TestLongTimeBreak(t *testing.T) {

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	orderMgtClient := pb.NewOrderManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Get Order
	retrievedOrder, err := orderMgtClient.GetOrder(ctx, &pb.OrderRequest{Value: "106"})
	log.Printf("GetOrder Response1 -> :%v, err:%v \n", retrievedOrder, err)

	time.Sleep(20 * time.Second)
	ctx2, _ := context.WithTimeout(context.Background(), 5*time.Second)
	retrievedOrder, err = orderMgtClient.GetOrder(ctx2, &pb.OrderRequest{Value: "106"})
	log.Printf("GetOrder Response2 -> :%v, err:%v \n", retrievedOrder, err)

	time.Sleep(1111 * time.Second)
}

// state := conn.GetState()
//开始是IDLE
//有请求或者连接中就是state: CONNECTING
//成功了就是state: READY
//对方断开连接了，然后就是state: CONNECTING
//然后就是state: TRANSIENT_FAILURE
//之后又是state: CONNECTING
//成功了就是state: READY
func TestState(t *testing.T) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithUnaryInterceptor(orderUnaryClientInterceptor))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	orderMgtClient := pb.NewOrderManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		for {
			state := conn.GetState()
			fmt.Println("state:", state)
			time.Sleep(1 * time.Second)
		}
	}()

	time.Sleep(5 * time.Second)
	// Get Order
	retrievedOrder, err := orderMgtClient.GetOrder(ctx,
		&pb.OrderRequest{Value: "106"})

	log.Printf("GetOrder Response -> :%v, err:%v \n", retrievedOrder, err)

	time.Sleep(10 * time.Second)
	orderMgtClient.GetOrder(ctx,
		&pb.OrderRequest{Value: "106"})

	time.Sleep(1111 * time.Second)
}

// 在TRANSIENT_FAILURE和CONNECTING之间变更，每次变更都notify返回一次
// WaitForStateChange当前状态不是sourceState则直接返回true，即不用等，若是相同则等待超时或取消则返回false，若是状态变更则返回true
// 名字就是等待状态变更，true变了(一上来就不一致或之后收到变更了不一致)，false未变(超时或取消)
func TestWaitForChange(t *testing.T) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithUnaryInterceptor(orderUnaryClientInterceptor))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	orderMgtClient := pb.NewOrderManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5111*time.Second)
	defer cancel()

	ctx1, _ := context.WithCancel(context.Background())
	for { // 循环处理，若是ready则返回，若是超时则处理超时，之后继续循环
		s := conn.GetState()
		fmt.Println("state:", s)
		if s == connectivity.Ready {
			fmt.Println("connection is ready")
			break
		}
		if !conn.WaitForStateChange(ctx1, s) {
			// ctx got timeout or canceled.
			// handle timeout
			fmt.Println("WaitForStateChange is false")
		}
	}

	time.Sleep(5 * time.Second)
	// Get Order
	retrievedOrder, err := orderMgtClient.GetOrder(ctx,
		&pb.OrderRequest{Value: "106"})

	log.Printf("GetOrder Response -> :%v, err:%v \n", retrievedOrder, err)

	time.Sleep(555 * time.Second)
}
