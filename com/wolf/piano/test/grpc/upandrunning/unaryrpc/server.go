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
	"google.golang.org/grpc/status"
	"log"
	"net"
)

type OrderMap map[string]pb.Order

var (
	orderMap = map[string]pb.Order{}
	port     = "22222"
)

type server struct {
}

func (s *server) GetOrder(ctx context.Context,
	orderId *wrappers.StringValue) (*pb.Order, error) {
	// Service Implementation.
	ord := orderMap[orderId.Value]
	return &ord, nil
}

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

func TestBase() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(orderUnaryServerInterceptor))
	pb.RegisterOrderManagementServer(s, &server{})

	log.Printf("Starting gRPC listener on port " + port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type helloServer struct {
}

func (s *helloServer) Hello(context.Context, *wrappers.StringValue) (*wrappers.StringValue, error) {
	return nil, nil
}

func TestRegistryMultiService() {
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
