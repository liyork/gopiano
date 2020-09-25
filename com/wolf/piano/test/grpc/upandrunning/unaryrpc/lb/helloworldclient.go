package main

import (
	"flag"
	"fmt"
	pb "github.com/liyork/gopiano/com/wolf/piano/test/grpc/upandrunning/unaryrpc/ecommerce"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"strconv"
	"time"
)

var (
	serverName = flag.String("service", "ecommerce/OrderManagement", "service name")
	reg        = flag.String("reg", "http://127.0.0.1:2379", "register etcd address")
)

// lb关联resolver，resolver关联Watcher，Watcher会被lb定期获取然后更新到他的地址中用于后续负载以及更新连接
func main() {
	flag.Parse()
	r := NewResolver(*serverName)
	b := grpc.RoundRobin(r)

	ctx, _ := context.WithTimeout(context.Background(), 110*time.Second)
	conn, err := grpc.DialContext(ctx, *reg, grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(1 * time.Second)
	for t := range ticker.C {
		client := pb.NewOrderManagementClient(conn)
		resp, err := client.GetOrder(context.Background(), &pb.OrderRequest{Value: "world " + strconv.Itoa(t.Second())})
		if err == nil {
			fmt.Printf("%v: Reply is %v\n", t, resp)
		}
	}
}
