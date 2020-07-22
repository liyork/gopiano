package test

import (
	"context"
	pb "github.com/liyork/gopiano/com/wolf/piano/test/grpc/upandrunning/productinfo/ecommerce"
	"testing"
	"time"
)

//“mockgen github.com/grpc-up-and-running/samples/ch07/grpc-docker/go/proto-gen \
//ProductInfoClient > mock_prodinfo/prodinfo_mock.go”

func TestAddProduct(t *testing.T) {
	//ctrl := gomock.NewController(t)
	//defer ctrl.Finish()
	//mocklProdInfoClient := NewMockProductInfoClient(ctrl)
	//...
	//req := &pb.Product{Name: name, Description: description, Price: price}
	//
	//mocklProdInfoClient.
	//	EXPECT().AddProduct(gomock.Any(), &rpcMsg{msg: req},).
	//	Return(&wrapper.StringValue{Value: "ABC123" + name}, nil)

	//testAddProduct(t, mocklProdInfoClient)
}

func testAddProduct(t *testing.T, client pb.ProductInfoClient) {
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	//...
	//
	//r, err := client.AddProduct(ctx, &pb.Product{Name: name,
	//	Description: description, Price: price})
	//
	//// test and verify response.
}
