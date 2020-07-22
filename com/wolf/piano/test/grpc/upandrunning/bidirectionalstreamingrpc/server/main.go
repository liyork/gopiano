package server

import (
	"context"
	"fmt"
	pb "github.com/liyork/gopiano/com/wolf/piano/test/grpc/upandrunning/bidirectionalstreamingrpc/ecommerce"
	"io"
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
		orderId, err := stream.Recv()
		if err == io.EOF {
			for _, comb := range combinedShipmentMap {
				stream.Send(&comb)
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
		fmt.Println("orderId:", orderId)

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

func TestServerCheckTimeout(stream pb.OrderManagement_ProcessOrdersServer) {
	if stream.Context().Err() == context.Canceled {

	}
}
