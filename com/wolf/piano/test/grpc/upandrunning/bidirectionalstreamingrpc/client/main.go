package client

import (
	"context"
	"github.com/golang/protobuf/ptypes/wrappers"
	pb "github.com/liyork/gopiano/com/wolf/piano/test/grpc/upandrunning/bidirectionalstreamingrpc/ecommerce"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

// “the client and server can read and write in any order—the streams operate completely independently”

func TestClientBase() {
	conn, _ := grpc.Dial("localshot:xxxxxx", grpc.WithInsecure())
	client := pb.NewOrderManagementClient(conn)
	// Process Order
	streamProcOrder, _ := client.ProcessOrders(ctx)
	if err := streamProcOrder.Send(
		&wrappers.StringValue{Value: "102"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", client, "102", err)
	}
	if err := streamProcOrder.Send(
		&wrappers.StringValue{Value: "103"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", client, "103", err)
	}

	if err := streamProcOrder.Send(
		&wrappers.StringValue{Value: "104"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", client, "104", err)
	}

	channel := make(chan struct{})
	go asncClientBidirectionalRPC(streamProcOrder, channel)
	time.Sleep(time.Millisecond * 1000)

	if err := streamProcOrder.Send(
		&wrappers.StringValue{Value: "101"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", client, "101", err)
	}

	// Mark the end of stream
	if err := streamProcOrder.CloseSend(); err != nil {
		log.Fatal(err)
	}

	<-channel
}

func asncClientBidirectionalRPC(
	streamProcOrder pb.OrderManagement_ProcessOrdersClient,
	c chan struct{}) {
	for {
		combinedShipment, errProcOrder := streamProcOrder.Recv()
		if errProcOrder == io.EOF { // “detect the end of the stream”
			break
		}
		if errProcOrder != nil {
			log.Printf("Error Receiving messages %v", errProcOrder)
		}
		log.Printf("Combined shipment : ", combinedShipment.OrdersList)
	}
	<-c
}

func TestClientCancel() {
	conn, _ := grpc.Dial("localshot:xxxxxx", grpc.WithInsecure())
	client := pb.NewOrderManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	streamProcOrder, _ := client.ProcessOrders(ctx)
	_ = streamProcOrder.Send(&wrappers.StringValue{Value: "102"})
	_ = streamProcOrder.Send(&wrappers.StringValue{Value: "103"})
	_ = streamProcOrder.Send(&wrappers.StringValue{Value: "104"})

	channel := make(chan struct{})

	go asncClientBidirectionalRPC(streamProcOrder, channel)
	time.Sleep(time.Millisecond * 1000)

	// Canceling the RPC
	cancel()
	log.Printf("RPC Status : %s", ctx.Err())

	_ = streamProcOrder.Send(&wrappers.StringValue{Value: "101"})
	_ = streamProcOrder.CloseSend()

	<-channel
}
