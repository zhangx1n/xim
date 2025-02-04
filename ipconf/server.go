package ipconf

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/zhangx1n/xim/common/config"
	"github.com/zhangx1n/xim/ipconf/domain"
	"github.com/zhangx1n/xim/ipconf/source"
)

func RunMain(path string) {
	config.Init(path)
	source.Init()
	domain.Init()
	s := server.Default(server.WithHostPorts(":6969"))
	s.GET("/ip/list", GetIpInfoList)
	s.Spin()
}
