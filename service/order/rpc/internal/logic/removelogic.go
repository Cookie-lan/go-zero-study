package logic

import (
	"context"
	"go-zero-study/service/order/model"
	"google.golang.org/grpc/status"
	"net/http"

	"go-zero-study/service/order/rpc/internal/svc"
	"go-zero-study/service/order/rpc/types/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveLogic {
	return &RemoveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveLogic) Remove(in *order.RemoveRequest) (*order.RemoveResponse, error) {
	od, err := l.svcCtx.OrderModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(http.StatusContinue, "订单不存在")
		}
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if err = l.svcCtx.OrderModel.Delete(l.ctx, od.Id); err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}
	return &order.RemoveResponse{}, nil
}
