package unaryrpc

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/wrappers"
	pb "github.com/liyork/gopiano/com/wolf/piano/test/grpc/upandrunning/unaryrpc/ecommerce"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"testing"
)

type OrderMap map[string]pb.Order

var (
	orderMap = map[string]pb.Order{}
	port     = "0.0.0.0:23444"
	port2    = "0.0.0.0:23445"
)

type server struct {
}

// 正常
func (s *server) GetOrder(ctx context.Context,
	orderId *pb.OrderRequest) (*pb.Order, error) {
	addr, _ := GetPeerAddr(ctx)
	fmt.Println("client adddr:", addr)

	// Service Implementation.
	//time.Sleep(225 * time.Second)
	//ord := orderMap[orderId.Value]
	//return &ord, nil
	fmt.Println("GetOrder")
	return &pb.Order{Id: "1"}, nil
}

//// 测试重试
//func (s *server) GetOrder(ctx context.Context,
//	orderId *pb.OrderRequest) (*pb.Order, error) {
//	// Service Implementation.
//	//time.Sleep(225 * time.Second)
//	//ord := orderMap[orderId.Value]
//	//return &ord, nil
//	fmt.Println("GetOrder", )
//	// 测试重试
//	//return nil, status.Errorf(codes.Unavailable, "maybeFailRequest: failing it")
//	return nil, status.Errorf(codes.DataLoss, "maybeFailRequest: failing it")
//}

// 测试超时
//func (s *server) GetOrder(ctx context.Context, orderId *pb.OrderRequest) (*pb.Order, error) {
//	addr := GetPeerAddr(ctx)
//	fmt.Println("client adddr:", addr)
//
//	for i := 0; i < 5; i++ {
//		// 监测客户端是否已取消
//		if ctx.Err() == context.Canceled {
//			return nil, status.Errorf(codes.Canceled, "SearchService.Search canceled")
//		}
//		time.Sleep(1 * time.Second)
//	}
//	return &pb.Order{Id: "1"}, nil
//}

func (s *server) GetOrderOccurError(ctx context.Context, orderId *wrappers.StringValue) (*pb.Order, error) {
	if orderId.Value == "-1" {
		log.Printf("Order ID is invalid! -> Received Order ID %s",
			orderId.Value)

		errorStatus := status.New(codes.InvalidArgument, "Invalid information received")
		ds, err := errorStatus.WithDetails(
			&errdetails.BadRequest_FieldViolation{
				Field: "ID",
				Description: fmt.Sprintf(
					"Order ID received is not valid %s : %s",
					orderId.Value, "xxxx"),
			},
		)
		// 有错误用之前的errorStatus
		if err != nil {
			return nil, errorStatus.Err()
		}
		return nil, ds.Err()
	}
	return nil, nil
}

// 使用metadata
func (s *server) GetOrderWithMetadata(ctx context.Context,
	orderId *wrappers.StringValue) (*pb.Order, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	fmt.Println("md,ok", md, ok)
	// Service Implementation.
	ord := orderMap[orderId.Value]

	// create and send header
	header := metadata.Pairs("header-key", "val")
	grpc.SendHeader(ctx, header)
	// create and set trailer
	trailer := metadata.Pairs("trailer-key", "val")
	grpc.SetTrailer(ctx, trailer)

	return &ord, nil
}

func TestBase(t *testing.T) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(orderUnaryServerInterceptor))
	//s := grpc.NewServer()
	pb.RegisterOrderManagementServer(s, &server{})

	log.Printf("Starting gRPC listener on port " + port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func TestBase2(t *testing.T) {
	lis, err := net.Listen("tcp", port2)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(orderUnaryServerInterceptor))
	//s := grpc.NewServer()
	pb.RegisterOrderManagementServer(s, &server{})

	log.Printf("Starting gRPC listener on port " + port2)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type helloServer struct {
}

func (s *helloServer) Hello(context.Context, *pb.HelloRequest) (*pb.HelloResponse, error) {
	return nil, nil
}

func TestRegistryMultiService(t *testing.T) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(orderUnaryServerInterceptor))
	pb.RegisterOrderManagementServer(s, &server{})
	pb.RegisterHelloServerServer(s, &helloServer{})

	log.Printf("Starting gRPC listener on port " + port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// 中间有代理，使用此获取clientip
// 需要在代理转发时把代理对端的ip设定到header的x-real-ip
func GetRealAddr(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	rips := md.Get("x-real-ip")
	if len(rips) == 0 {
		return ""
	}

	return rips[0]
}

// 中间无代理，使用此获取clientip
func GetPeerAddr(ctx context.Context) (string, error) {
	var addr string
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("[getClinetIP] invoke FromContext() failed")
	}
	if pr.Addr == net.Addr(nil) {
		return "", fmt.Errorf("[getClientIP] peer.Addr is nil")
	}

	if tcpAddr, ok := pr.Addr.(*net.TCPAddr); ok {
		addr = tcpAddr.IP.String()
	} else {
		addr = pr.Addr.String()
	}

	return addr, nil
}
