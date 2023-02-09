package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"liyue/apps/im/gateway/internal/config"
	"liyue/apps/im/rpc/im"
)

type ServiceContext struct {
	Config config.Config
	ImRpc  im.Im
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		ImRpc:  im.NewIm(zrpc.MustNewClient(c.ImRpc)),
	}
}
