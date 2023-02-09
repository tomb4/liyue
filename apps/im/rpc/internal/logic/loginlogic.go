package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/hash"

	"liyue/apps/im/rpc/internal/svc"
	"liyue/apps/im/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *pb.LoginReq) (*pb.LoginResp, error) {
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Uid)
	if err != nil {
		l.Error("FindOne err:", err)
		return nil, err
	}

	if user.Password != hash.Md5Hex([]byte(in.Password)) {
		return nil, errors.New("passwords do not match")
	}

	return &pb.LoginResp{
		Username: user.Username,
	}, nil
}
