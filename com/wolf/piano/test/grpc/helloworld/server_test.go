package helloworld

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	pb "github.com/liyork/gopiano/com/wolf/piano/test/grpc/helloworld/proto"
	"google.golang.org/grpc"
)

const (
	port = ":50052"
)

// server is used to implement helloworld.GreeterServer.  实现
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello again " + in.GetName()}, nil
}

// 基本流程：
// 构造net.Listen监听
// grpc.NewServer构造server
// RegisterGreeterServer注册server
// s.Serve
func TestServerBase(t *testing.T) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	time.Sleep(99999 * time.Second)
}
