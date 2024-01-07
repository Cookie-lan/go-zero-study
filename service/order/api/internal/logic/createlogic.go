package logic

import (
	"context"
	"net/http"

	"go-zero-study/service/order/api/internal/svc"
	"go-zero-study/service/order/api/internal/types"
	"go-zero-study/service/order/rpc/orderclient"
	"go-zero-study/service/product/rpc/productclient"

	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.CreateRequest) (resp *types.CreateResponse, err error) {
	// 获取 OrderRpc BuildTarget
	_, err = l.svcCtx.Config.OrderRpc.BuildTarget()
	if err != nil {
		return nil, status.Error(http.StatusContinue, "订单创建异常")
	}

	// 获取 ProductRpc BuildTarget
	_, err = l.svcCtx.Config.ProductRpc.BuildTarget()
	if err != nil {
		return nil, status.Error(http.StatusContinue, "订单商品创建异常")
	}

	// 创建一个gid
	gid := dtmgrpc.MustGenGid(l.svcCtx.Config.DtmServer)
	// 创建一个saga协议的事务
	sage := dtmgrpc.NewSagaGrpc(l.svcCtx.Config.DtmServer, gid).
		Add("etcd://etcd-zhoutao:2379/order.rpc/orderclient.Order/Create",
			"etcd://etcd-zhoutao:2379/order.rpc/orderclient.Order/CreateRevert",
			&orderclient.CreateRequest{
				Uid:    req.Uid,
				Pid:    req.Pid,
				Amount: req.Amount,
				Status: req.Status,
			}).
		Add("etcd://etcd-zhoutao:2379/product.rpc/productclient.Product/DecrStock",
			"etcd://etcd-zhoutao:2379/product.rpc/productclient.Product/DecrStockRevert",
			&productclient.DecrStockRequest{
				Id:  req.Pid,
				Num: 1,
			})

	// 提交事务
	err = sage.Submit()
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &types.CreateResponse{}, nil
}
