package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zero-study/service/user/model"
	"go-zero-study/service/user/rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config

	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(conn, c.Cache),
	}
}
