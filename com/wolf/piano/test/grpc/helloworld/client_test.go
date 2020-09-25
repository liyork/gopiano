package helloworld

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	pb "github.com/liyork/gopiano/com/wolf/piano/test/grpc/helloworld/proto"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50052"
	defaultName = "world"
)

// 基本流程：
// grpc.Dial
// defer conn.Close
// NewGreeterClient
// c.SayHello
func TestClientBase(t *testing.T) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	// 超时,could not greet: rpc error: code = DeadlineExceeded desc = context deadline exceeded
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	fmt.Println("time1:", time.Now())
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	fmt.Println("time2:", time.Now())
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())

	r, err = c.SayHelloAgain(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
