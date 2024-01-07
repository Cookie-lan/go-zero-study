package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero-study/service/order/api/internal/config"
	"go-zero-study/service/order/rpc/orderclient"
)

type ServiceContext struct {
	Config config.Config

	OrderRpc orderclient.Order
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		OrderRpc: orderclient.NewOrder(zrpc.MustNewClient(c.OrderRpc)),
	}
}
