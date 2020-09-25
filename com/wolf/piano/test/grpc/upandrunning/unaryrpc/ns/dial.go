package ns

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

// Dial 封装 `grpc.Dial()` 方法以供业务方代码初始化 *grpc.ClientConn。
// 业务方可使用此 Dial 方法基于主调方服务名、被调方服务名等参数构造 *grpc.ClientConn 实例，
// 随后可在业务代码中使用 *grpc.ClientConn 实例构造桩代码中生成的 grpcServiceClient 并发起 RPC 调用。
func Dial(callerService, calleeService string, dialOpts ...grpc.DialOption) (*grpc.ClientConn, error) {
	// 根据 callerService 和 calleeService 构造对应的 URI
	URI := URI(callerService, calleeService)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithInsecure(),
	}
	dialOpts = append(dialOpts, dialOpts...)

	conn, err := grpc.DialContext(
		ctx,
		URI,
		opts...,
	)
	if err != nil {
		println("did not connect", URI, err)
		return nil, err
	}
	return conn, err
}
