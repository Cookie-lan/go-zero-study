package logic

import (
	"context"
	"go-zero-study/service/order/model"
	"go-zero-study/service/product/rpc/types/product"
	"go-zero-study/service/user/rpc/types/user"
	"google.golang.org/grpc/status"
	"net/http"

	"go-zero-study/service/order/rpc/internal/svc"
	"go-zero-study/service/order/rpc/types/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateLogic) Create(in *order.CreateRequest) (*order.CreateResponse, error) {
	// 查询用户是否存在
	_, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user.UserInfoRequest{
		Id: in.Uid,
	})
	if err != nil {
		return nil, err
	}
	// 查询产品是否存在
	pd, err := l.svcCtx.ProductRpc.Detail(l.ctx, &product.DetailRequest{
		Id: in.Pid,
	})
	if err != nil {
		return nil, err
	}
	// 判断产品库存是否充足
	if pd.Stock <= 0 {
		return nil, status.Error(http.StatusInternalServerError, "库存不足")
	}

	newOrder := &model.Order{
		Uid:    in.Uid,
		Pid:    in.Pid,
		Amount: in.Amount,
		Status: 0,
	}
	// 创建订单
	od, err := l.svcCtx.OrderModel.Insert(l.ctx, newOrder)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	newOrder.Id, err = od.LastInsertId()
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	// 扣减库存
	_, err = l.svcCtx.ProductRpc.Update(l.ctx, &product.UpdateRequest{
		Id:    in.Pid,
		Stock: pd.Stock - 1,
	})
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &order.CreateResponse{
		Id: newOrder.Id,
	}, nil
}
