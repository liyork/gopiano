package server

import (
	"context"
	"fmt"
	pb "github.com/liyork/gopiano/com/wolf/piano/test/grpc/upandrunning/bidirectionalstreamingrpc/ecommerce"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"testing"
	"time"
)

type server struct{}

type OrderMap map[string]pb.Order

var (
	combinedShipmentMap = map[string]pb.CombinedShipment{}
	batchMarker         = 0
	orderBatchSize      = 5
)

func (s *server) ProcessOrders(
	stream pb.OrderManagement_ProcessOrdersServer) error {
	for {
		orderReply, err := stream.Recv()
		if err == io.EOF {
			for i := 0; i < 10; i++ {
				shipment := pb.CombinedShipment{}
				stream.Send(&shipment)
				time.Sleep(2 * time.Second)
			}
			// end of the stream
			return nil
		}
		if err != nil {
			return err
		}

		// Logic to organize orders into shipments,
		// based on the destination.
		//
		fmt.Println("orderId:", orderReply)

		if batchMarker == orderBatchSize {
			// Stream combined orders to the client in batches
			for _, comb := range combinedShipmentMap {
				// Send combined shipment to the client
				stream.Send(&comb)
			}
			batchMarker = 0
			combinedShipmentMap = make(
				map[string]pb.CombinedShipment)
		} else {
			batchMarker++
		}
	}
}

func ServerCheckTimeout(stream pb.OrderManagement_ProcessOrdersServer) {
	if stream.Context().Err() == context.Canceled {

	}
}

func TestServerBase(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:12345")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOrderManagementServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	time.Sleep(99999 * time.Second)
}
