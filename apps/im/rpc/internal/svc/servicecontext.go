package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"liyue/apps/im/rpc/internal/config"
	"liyue/apps/im/rpc/model"
)

type ServiceContext struct {
	Config config.Config
	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		UserModel: model.NewUserModel(sqlx.NewMysql(c.DataSource), c.Cache),
	}
}
