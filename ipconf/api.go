package ipconf

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/zhangx1n/plato/ipconf/domain"
)

type Response struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

func GetIpInfoList(c context.Context, ctx *app.RequestContext) {
	defer func() {
		if err := recover(); err != nil {
			ctx.JSON(consts.StatusBadRequest, utils.H{"err": err})
		}
	}()

	// step0: 构建用户请求信息
	ipConfCtx := domain.BuildIpConfContext(&c, ctx)
	// step1: 进行 ip 调度
	eds := domain.Dispatch(ipConfCtx)
	// step2: 根据得分返回 top5 的 ip 信息
	ipConfCtx.AppCtx.JSON(consts.StatusOK, packRes(top5Endports(eds)))
}
