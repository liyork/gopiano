package unaryrpc

import (
	"context"
	"google.golang.org/grpc"
	"log"
)

// server
func orderUnaryServerInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	// Preprocessing logic
	// Gets info about the current RPC call by examining the args passed in
	log.Println("======= [Server Interceptor] ", info.FullMethod)

	// Invoking the handler to complete the normal execution of a unary RPC.
	m, err := handler(ctx, req)

	// Post processing logic
	log.Printf(" Post Proc Message : %s", m)
	return m, err
}

// client
func orderUnaryClientInterceptor(
	ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// Preprocessor phase
	log.Println("Method : " + method)

	// Invoking the remote method
	err := invoker(ctx, method, req, reply, cc, opts...)

	// Postprocessor phase
	log.Println(reply)

	return err
}
