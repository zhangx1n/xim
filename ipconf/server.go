package ipconf

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/zhangx1n/plato/ipconf/domain"
	"github.com/zhangx1n/plato/ipconf/source"
)

func RunMain() {
	source.Init()
	domain.Init()
	s := server.Default(server.WithHostPorts(":6969"))
	s.GET("/ip/list", GetIpInfoList)
	s.Spin()
}
