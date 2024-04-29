package source

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/zhangx1n/plato/common/discovery"
)

func testServiceRegister(ctx *context.Context, port, node string) {
	// 模拟服务发现
	go func() {
		ed := discovery.EndpointInfo{
			IP:   "127.0.0.1",
			Port: port,
			MetaData: map[string]interface{}{
				"connect_num":   float64(rand.Int63n(12312321231231131)),
				"message_bytes": float64(rand.Int63n(1231232131556)),
			},
		}
		sr, err := discovery.NewServiceRegister(ctx, fmt.Sprintf("/plato/ip_dispatcher/%s", node), &ed, time.Now().Unix())
		if err != nil {
			panic(err)
		}
		go sr.ListenLeaseRespChan()
		for {
			ed = discovery.EndpointInfo{
				IP:   "127.0.0.1",
				Port: port,
				MetaData: map[string]interface{}{
					"connect_num":   float64(rand.Int63n(12312321231231131)),
					"message_bytes": float64(rand.Int63n(1231232131556)),
				},
			}
			err := sr.UpdateValue(&ed)
			if err != nil {
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()
}
