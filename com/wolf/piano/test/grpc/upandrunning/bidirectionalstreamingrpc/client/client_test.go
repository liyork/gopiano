package client

import (
	"context"
	"fmt"
	pb "github.com/liyork/gopiano/com/wolf/piano/test/grpc/upandrunning/bidirectionalstreamingrpc/ecommerce"
	"google.golang.org/grpc"
	"io"
	"log"
	"testing"
	"time"
)

// “the client and server can read and write in any order—the streams operate completely independently”

func TestClientBase(t *testing.T) {
	//ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	conn, _ := grpc.Dial("localhost:12345", grpc.WithInsecure())
	client := pb.NewOrderManagementClient(conn)
	// Process Order
	streamProcOrder, _ := client.ProcessOrders(ctx)
	if err := streamProcOrder.Send(
		&pb.OrderReply{Message: "102"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", client, "102", err)
	}
	if err := streamProcOrder.Send(
		&pb.OrderReply{Message: "103"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", client, "103", err)
	}

	if err := streamProcOrder.Send(
		&pb.OrderReply{Message: "104"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", client, "104", err)
	}

	channel := make(chan struct{})
	go asncClientBidirectionalRPC(streamProcOrder, channel)
	time.Sleep(time.Millisecond * 1000)

	go func(conn *grpc.ClientConn) {
		time.Sleep(2 * time.Second)
		fmt.Println("close conn")
		//conn.Close()
		fmt.Println("close conn")
		cancel()
	}(conn)

	if err := streamProcOrder.Send(
		&pb.OrderReply{Message: "101"}); err != nil {
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
		fmt.Println("before streamProcOrder.Recv return")
		_, errProcOrder := streamProcOrder.Recv()
		fmt.Println("after streamProcOrder.Recv return")
		if errProcOrder == io.EOF { // “detect the end of the stream”
			fmt.Println("after streamProcOrder.Recv EOF return")
			break
		}
		if errProcOrder != nil {
			fmt.Println("Error Receiving messages ", errProcOrder)
		}
		//fmt.Println("Combined shipment : ", combinedShipment.OrdersList)
	}
	<-c
}

func TestClientCancel(t *testing.T) {
	conn, _ := grpc.Dial("localshot:xxxxxx", grpc.WithInsecure())
	client := pb.NewOrderManagementClient(conn)

	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	ctx, cancel := context.WithCancel(context.Background())

	streamProcOrder, _ := client.ProcessOrders(ctx)
	_ = streamProcOrder.Send(&pb.OrderReply{Message: "102"})
	_ = streamProcOrder.Send(&pb.OrderReply{Message: "103"})
	_ = streamProcOrder.Send(&pb.OrderReply{Message: "104"})

	channel := make(chan struct{})

	go asncClientBidirectionalRPC(streamProcOrder, channel)
	time.Sleep(time.Millisecond * 1000)

	// Canceling the RPC
	cancel()
	log.Printf("RPC Status : %s", ctx.Err())

	_ = streamProcOrder.Send(&pb.OrderReply{Message: "101"})
	_ = streamProcOrder.CloseSend()

	<-channel
}
