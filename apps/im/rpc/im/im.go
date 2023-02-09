// Code generated by goctl. DO NOT EDIT!
// Source: im.proto

package im

import (
	"context"

	"liyue/apps/im/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	LoginReq  = pb.LoginReq
	LoginResp = pb.LoginResp

	Im interface {
		Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginResp, error)
	}

	defaultIm struct {
		cli zrpc.Client
	}
)

func NewIm(cli zrpc.Client) Im {
	return &defaultIm{
		cli: cli,
	}
}

func (m *defaultIm) Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginResp, error) {
	client := pb.NewImClient(m.cli.Conn())
	return client.Login(ctx, in, opts...)
}
