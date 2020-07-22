package unaryrpc

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/wrappers"
	pb "github.com/liyork/gopiano/com/wolf/piano/test/grpc/upandrunning/unaryrpc/ecommerce"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

const (
	address = "localhost:22222"
)

func TestClientBase() {

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
		&wrappers.StringValue{Value: "106"})
	log.Print("GetOrder Response -> : ", retrievedOrder)
}

func TestClientDeadline() {

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrderManagementClient(conn)

	clientDeadline := time.Now().Add(
		time.Duration(2 * time.Second))
	ctx, cancel := context.WithDeadline(
		context.Background(), clientDeadline)

	defer cancel()

	// Add Order
	order1 := pb.Order{Id: "101",
		Items:       []string{"iPhone XS", "Mac Book Pro"},
		Destination: "San Jose, CA",
		Price:       2300.00}
	//If the invocation exceeds the specified deadline, it should return an error of the type DEADLINE_EXCEEDED
	res, addErr := client.GetOrder(ctx, &order1)

	if addErr != nil {
		got := status.Code(addErr)
		log.Printf("Error Occured -> addOrder : , %v:", got)
	} else {
		log.Print("AddOrder Response -> ", res.Id)
	}
}

func TestClientProcessErr() {

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
		&wrappers.StringValue{Value: "106"})
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
func TestClientMultiplex() {

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
		&wrappers.StringValue{Value: "106"})
	log.Print("GetOrder Response -> : ", retrievedOrder)

	helloClient := pb.NewHelloServerClient(conn)
	helloClient.Hello(ctx, &wrappers.StringValue{Value: "111"})
}

func TestClientMetadata() {

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
		&wrappers.StringValue{Value: "106"},
		grpc.Header(&header),
		grpc.Trailer(&trailer))
	log.Print("GetOrder Response -> : ", retrievedOrder, header, trailer)
}

const (
	exampleScheme      = "exampleScheme"
	exampleServiceName = "exampleServiceName"
)

func TestClientLb() {

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

func TestClientCompression() {
	//client.AddOrder(ctx, &order1, grpc.UseCompressor(gzip.Name))
}

// log
//GRPC_GO_LOG_VERBOSITY_LEVEL=99
//GRPC_GO_LOG_SEVERITY_LEVEL=info
