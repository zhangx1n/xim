package helloservice

import "context"

type HelloServer struct {
	UnimplementedGreeterServer
}

func (h HelloServer) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	return &HelloReply{
		Message: "Hello " + req.Name,
	}, nil
}
