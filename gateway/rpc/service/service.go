package service

import (
	context "context"
)

const (
	DelConnCmd = 1 // DelConn
	PushCmd    = 2 // push
)

type CmdContext struct {
	Ctx     *context.Context
	Cmd     int32
	FD      int
	Payload []byte
}

type Service struct {
	UnimplementedGatewayServer
	CmdChannel chan *CmdContext
}

func (s *Service) DelConn(ctx context.Context, gr *GatewayRequest) (*GatewayResponse, error) {
	c := context.TODO()
	s.CmdChannel <- &CmdContext{
		Ctx: &c,
		Cmd: DelConnCmd,
		FD:  int(gr.GetFd()),
	}
	return &GatewayResponse{
		Code: 0,
		Msg:  "success",
	}, nil
}

func (s *Service) Push(ctx context.Context, gr *GatewayRequest) (*GatewayResponse, error) {
	c := context.TODO()
	s.CmdChannel <- &CmdContext{
		Ctx:     &c,
		Cmd:     PushCmd,
		FD:      int(gr.GetFd()),
		Payload: gr.GetData(),
	}
	return &GatewayResponse{
		Code: 0,
		Msg:  "success",
	}, nil
}
