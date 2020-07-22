package client

import (
	"context"
	pb "github.com/liyork/gopiano/com/wolf/piano/test/grpc/upandrunning/clientstreamingrpc/ecommerce"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	address = "localhost:22222"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithStreamInterceptor(clientStreamInterceptor))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrderManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	updateStream, err := client.UpdateOrders(ctx)

	//get metadata
	updateStream.Header()
	// Trailers are used to send status codes and the status message.
	updateStream.Trailer()

	if err != nil {
		log.Fatalf("%v.UpdateOrders(_) = _, %v", client, err)
	}

	// Updating order 1
	if err := updateStream.Send(&pb.Order{}); err != nil {
		log.Fatalf("%v.Send(%v) = %v",
			updateStream, &pb.Order{}, err)
	}

	// Updating order 2
	if err := updateStream.Send(&pb.Order{}); err != nil {
		log.Fatalf("%v.Send(%v) = %v",
			updateStream, &pb.Order{}, err)
	}

	// Updating order 3
	if err := updateStream.Send(&pb.Order{}); err != nil {
		log.Fatalf("%v.Send(%v) = %v",
			updateStream, &pb.Order{}, err)
	}

	updateRes, err := updateStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v",
			updateStream, err, nil)
	}
	log.Printf("Update Orders Res : %s", updateRes)
}
