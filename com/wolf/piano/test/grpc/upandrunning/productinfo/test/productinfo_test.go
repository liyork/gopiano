package test

import (
	"context"
	pb "github.com/liyork/gopiano/com/wolf/piano/test/grpc/upandrunning/productinfo/ecommerce"
	"google.golang.org/grpc"
	"log"
	"testing"
	"time"
)

const address = "localhost:22222"

func TestServer_AddProduct(t *testing.T) {
	grpcServer := initGRPCServerHTTP2()
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {

		grpcServer.Stop()
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProductInfoClient(conn)

	name := "Sumsung S10"
	description := "Samsung Galaxy S10 is the latest smart phone, launched in February 2019"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.AddProduct(ctx, &pb.Product{Name: name,
		Description: description})
	if err != nil {
		t.Fatalf("Could not add product: %v", err)
	}

	if r.Value == "" {
		t.Errorf("Invalid Product ID %s", r.Value)
	}
	log.Printf("Res %s", r.Value)
	grpcServer.Stop()
}

func initGRPCServerHTTP2() grpc.Server {
	return grpc.Server{}
}
