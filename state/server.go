package state

import (
	"context"
	"fmt"
	"github.com/zhangx1n/plato/common/config"
	idl "github.com/zhangx1n/plato/common/idl/state"
	"github.com/zhangx1n/plato/common/prpc"
	"github.com/zhangx1n/plato/state/rpc/client"
	"github.com/zhangx1n/plato/state/rpc/service"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"sync"
	"time"
)

// RunMain 启动网关服务
func RunMain(path string) {
	config.Init(path)
	cmdChannel = make(chan *service.CmdContext, config.GetSateCmdChannelNum())
	connToStateTable = sync.Map{}
	s := prpc.NewPServer(
		prpc.WithServiceName(config.GetStateServiceName()),
		prpc.WithIP(config.GetSateServiceAddr()),
		prpc.WithPort(config.GetSateServerPort()), prpc.WithWeight(config.GetSateRPCWeight()))

	s.RegisterService(func(server *grpc.Server) {
		service.RegisterStateServer(server, &service.Service{CmdChannel: cmdChannel})
	})
	// 初始化RPC 客户端
	client.Init()
	// 启动时间轮
	InitTimer()
	// 启动 命令处理写协程
	go cmdHandler()
	// 启动 rpc server
	s.Start(context.TODO())
}

func cmdHandler() {
	for cmdCtx := range cmdChannel {
		switch cmdCtx.Cmd {
		case service.CancelConnCmd:
			fmt.Printf("cancelconn endpoint:%s, fd:%d, data:%+v", cmdCtx.Endpoint, cmdCtx.ConnID, cmdCtx.Payload)
		case service.SendMsgCmd:
			fmt.Println("cmdHandler", string(cmdCtx.Payload))
			msgCmd := &idl.MsgCmd{}
			err := proto.Unmarshal(cmdCtx.Payload, msgCmd)
			if err != nil {
				fmt.Printf("SendMsgCmd:err=%s\n", err.Error())
			}
			msgCmdHandler(cmdCtx, msgCmd)
		}
	}
}

func msgCmdHandler(cmdCtx *service.CmdContext, msgCmd *idl.MsgCmd) {
	switch msgCmd.Type {
	case idl.CmdType_Login:
		loginMsgHandler(cmdCtx, msgCmd)
	case idl.CmdType_Heartbeat:
		hearbeatMsgHandler(cmdCtx, msgCmd)
	case idl.CmdType_ReConn:
		reConnMsgHandler(cmdCtx, msgCmd)
	}
}

func reConnMsgHandler(cmdCtx *service.CmdContext, msgCmd *idl.MsgCmd) {
	reConnMsg := &idl.ReConnMsg{}
	err := proto.Unmarshal(msgCmd.Payload, reConnMsg)
	if err != nil {
		fmt.Printf("reConnMsgHandler:err=%s\n", err.Error())
		return
	}
	// 重连的消息头中的connID才是上一次断开连接的connID
	if data, ok := connToStateTable.Load(reConnMsg.Head.ConnID); ok {
		state, _ := data.(*connState)
		state.Lock()
		defer state.Unlock()
		// 先停止定时任务的回调
		if state.reConnTimer != nil {
			state.reConnTimer.Stop()
			state.reConnTimer = nil // 重连定时器被清除
		}
		// 从索引中删除 旧的connID
		connToStateTable.Delete(reConnMsg.Head.ConnID)
		// 变更connID, cmdCtx中的connID才是 gateway重连的新连接
		state.connID = cmdCtx.ConnID
		connToStateTable.Store(cmdCtx.ConnID, state)
		sendACKMsg(cmdCtx.ConnID, 0, "reconn ok")
	} else {
		sendACKMsg(cmdCtx.ConnID, 1, "reconn feailed")
	}
}

func hearbeatMsgHandler(cmdCtx *service.CmdContext, msgCmd *idl.MsgCmd) {
	heartMsg := &idl.HeartbeatMsg{}
	err := proto.Unmarshal(msgCmd.Payload, heartMsg)
	if err != nil {
		fmt.Printf("hearbeatMsgHandler:err=%s\n", err.Error())
		return
	}
	if data, ok := connToStateTable.Load(cmdCtx.ConnID); ok {
		sate, _ := data.(*connState)
		sate.reSetHeartTimer()
	}
	// 未减少通信量，可以暂时不回复心跳的ack
}

func loginMsgHandler(cmdCtx *service.CmdContext, msgCmd *idl.MsgCmd) {
	loginMsg := &idl.LoginMsg{}
	err := proto.Unmarshal(msgCmd.Payload, loginMsg)
	if err != nil {
		fmt.Printf("loginMsgHandler:err=%s\n", err.Error())
		return
	}
	if loginMsg.Head != nil {
		// 这里会把 login msg 传送给业务层做处理
		fmt.Println("loginMsgHandler", loginMsg.Head.DeviceID)
	}
	// 创建定时器
	t := AfterFunc(300*time.Second, func() {
		clearState(cmdCtx.ConnID)
	})
	// 初始化连接的状态
	connToStateTable.Store(cmdCtx.ConnID, &connState{heartTimer: t, connID: cmdCtx.ConnID})
}
func sendACKMsg(connID uint64, code uint32, msg string) {
	ackMsg := &idl.ACKMsg{}
	ackMsg.Code = code
	ackMsg.Msg = msg
	ctx := context.TODO()
	downLoad, err := proto.Marshal(ackMsg)
	if err != nil {
		fmt.Println("sendACKMsg", err)
	}
	client.Push(&ctx, connID, downLoad)
}
