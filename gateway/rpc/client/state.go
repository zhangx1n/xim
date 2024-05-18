package client

import (
	"context"
	"github.com/zhangx1n/plato/common/config"
	"github.com/zhangx1n/plato/common/prpc"
	"github.com/zhangx1n/plato/state/rpc/service"
	"time"
)

var stateClient service.StateClient

func initStateClient() {
	pCli, err := prpc.NewPClient(config.GetStateServiceName())
	if err != nil {
		panic(err)
	}
	stateClient = service.NewStateClient(pCli.Conn())
}

func CancelConn(ctx *context.Context, endpoint string, connID uint64, payLoad []byte) error {
	rpcCtx, _ := context.WithTimeout(*ctx, 100*time.Millisecond)
	stateClient.CancelConn(rpcCtx, &service.StateRequest{
		Endpoint: endpoint,
		ConnID:   connID,
		Data:     payLoad,
	})
	return nil
}

func SendMsg(ctx *context.Context, endpoint string, connID uint64, payLoad []byte) error {
	rpcCtx, _ := context.WithTimeout(*ctx, 100*time.Millisecond)
	_, err := stateClient.SendMsg(rpcCtx, &service.StateRequest{
		Endpoint: endpoint,
		ConnID:   connID,
		Data:     payLoad,
	})
	if err != nil {
		panic(err)
	}
	return nil
}
