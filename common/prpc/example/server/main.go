package main

import (
	"context"
	"runtime"
	"strings"

	"github.com/zhangx1n/plato/common/config"
	"github.com/zhangx1n/plato/common/prpc"
	"github.com/zhangx1n/plato/common/prpc/example/helloservice"
	ptrace "github.com/zhangx1n/plato/common/prpc/trace"
	"google.golang.org/grpc"
)

const (
	testIp   = "127.0.0.1"
	testPort = 8867
)

func main() {
	config.Init(currentFileDir() + "/prpc_server.yaml")

	ptrace.StartAgent()
	defer ptrace.StopAgent()

	s := prpc.NewPServer(prpc.WithServiceName("prpc_server"), prpc.WithIP(testIp), prpc.WithPort(testPort), prpc.WithWeight(100))
	s.RegisterService(func(server *grpc.Server) {
		helloservice.RegisterGreeterServer(server, helloservice.HelloServer{})
	})
	s.Start(context.TODO())
}

func currentFileDir() string {
	// runtime.Caller(1) 函数是从 runtime 包提供的一个函数，用于获取关于当前函数调用堆栈的信息。
	// 参数 1 指明要获取当前函数的调用者的信息（即调用 currentFileDir 函数的那个函数的信息）。
	_, file, _, ok := runtime.Caller(1)
	// file: /home/user/project/src/main.go
	// parts: [home user project src main.go]
	parts := strings.Split(file, "/")

	if !ok {
		return ""
	}

	dir := ""
	for i := 0; i < len(parts)-1; i++ {
		dir += "/" + parts[i]
	}
	// dir: /home/user/project/src
	return dir[1:]
}
