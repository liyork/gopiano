package server

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/wrappers"
	pb "github.com/liyork/gopiano/com/wolf/piano/test/grpc/upandrunning/clientstreamingrpc/ecommerce"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
)

type server struct{}

type OrderMap map[string]pb.Order

var (
	orderMap = map[string]pb.Order{}
)

// If the server decides to prematurely stop reading from the client’s stream, the server should cancel the client stream so the client knows to stop producing messages.
func (s *server) UpdateOrders(stream pb.OrderManagement_UpdateOrdersServer) error {

	ordersStr := "Updated Order IDs : "
	for {
		order, err := stream.Recv()
		if err == io.EOF {
			// Finished reading the order stream.
			return stream.SendAndClose(
				&wrappers.StringValue{Value: "Orders processed " + ordersStr})
		}
		// Update order
		orderMap[order.Id] = *order

		log.Printf("Order ID ", order.Id, ": Updated")
		ordersStr += order.Id + ", "
	}
}

func (s *server) UpdateOrdersWithMetadata(stream pb.OrderManagement_UpdateOrdersServer) error {

	// read
	md, ok := metadata.FromIncomingContext(stream.Context())
	fmt.Println("md,ok", md, ok)

	// write
	// create and send header
	header := metadata.Pairs("header-key", "val")
	stream.SendHeader(header)
	// create and set trailer
	trailer := metadata.Pairs("trailer-key", "val")
	stream.SetTrailer(trailer)

	ordersStr := "Updated Order IDs : "
	for {
		order, err := stream.Recv()
		if err == io.EOF {
			// Finished reading the order stream.
			return stream.SendAndClose(
				&wrappers.StringValue{Value: "Orders processed " + ordersStr})
		}
		// Update order
		orderMap[order.Id] = *order

		log.Printf("Order ID ", order.Id, ": Updated")
		ordersStr += order.Id + ", "
	}
}
