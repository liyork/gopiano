package server

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/wrappers"
	pb "github.com/liyork/gopiano/com/wolf/piano/test/grpc/upandrunning/serverstreamingrpc/ecommerce"
	"google.golang.org/grpc"
	"log"
	"net"
	"strings"
)

type server struct {
}

type OrderMap map[string]pb.Order

var (
	orderMap = map[string]pb.Order{}
	port     = "2222"
)

// OrderManagement_SearchOrdersServer acts as a reference object to the stream that we can write multiple responses to.
func (s *server) SearchOrders(searchQuery *wrappers.StringValue,
	stream pb.OrderManagement_SearchOrdersServer) error {

	for key, order := range orderMap {
		log.Print(key, order)
		for _, itemStr := range order.Items {
			log.Print(itemStr)
			if strings.Contains(
				itemStr, searchQuery.Value) {
				// Send the matching orders in a stream
				err := stream.Send(&order)
				if err != nil {
					return fmt.Errorf(
						"error sending message to stream : %v",
						err)
				}
				log.Print("Matching Order Found : " + key)
				break
			}
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.StreamInterceptor(orderServerStreamInterceptor))
	pb.RegisterOrderManagementServer(s, &server{})

	log.Printf("Starting gRPC listener on port " + port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
