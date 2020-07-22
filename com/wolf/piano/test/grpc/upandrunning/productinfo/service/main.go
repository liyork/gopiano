package main

import (
	"context"
	"github.com/gofrs/uuid"
	pb "github.com/liyork/gopiano/com/wolf/piano/test/grpc/upandrunning/productinfo/ecommerce"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

const (
	port = ":22222"
)

type server struct {
	productMap map[string]*pb.Product
}

// Add product remote method
func (s *server) AddProduct(ctx context.Context, in *pb.Product) (
	*pb.ProductID, error) {
	out, err := uuid.NewV4()
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"Error while generating Product ID", err)
	}
	in.Id = out.String()
	if s.productMap == nil {
		s.productMap = make(map[string]*pb.Product)
	}
	s.productMap[in.Id] = in
	return &pb.ProductID{Value: in.Id}, status.New(codes.OK, "").Err()
}

// A Context object contains metadata such as the identity of the end user authorization tokens and
// the request’s deadline, and it will exist during the lifetime of the request.
// Get product remote method
func (s *server) GetProduct(ctx context.Context, in *pb.ProductID) (
	*pb.Product, error) {
	value, exists := s.productMap[in.Value]
	if exists {
		return value, status.New(codes.OK, "").Err()
	}
	return nil, status.Errorf(codes.NotFound, "Product does not exist.", in.Value)
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProductInfoServer(s, &server{})

	log.Printf("Starting gRPC listener on port " + port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
