package logic

import (
	"context"
	"liyue/apps/im/rpc/im"

	"github.com/zeromicro/go-zero/core/logx"
	"liyue/apps/im/gateway/internal/svc"
)

type GatewayLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGatewayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayLogic {
	return &GatewayLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (g *GatewayLogic) Login(uid int64, password string) (User, error) {
	res, err := g.svcCtx.ImRpc.Login(g.ctx, &im.LoginReq{
		Uid:      uid,
		Password: password,
	})
	if err != nil {
		return User{}, err
	}
	return User{
		Uid:      uid,
		UserName: res.Username,
	}, nil
}

func (g *GatewayLogic) SendMessage(from, to int64, content string) error {
	return nil
}
