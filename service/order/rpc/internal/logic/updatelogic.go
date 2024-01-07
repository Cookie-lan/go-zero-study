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

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLogic) Update(in *order.UpdateRequest) (*order.UpdateResponse, error) {
	od, err := l.svcCtx.OrderModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(http.StatusContinue, "订单不存在")
		}
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if in.Uid > 0 {
		od.Uid = in.Uid
	}
	if in.Pid > 0 {
		od.Pid = in.Pid
	}
	if in.Amount > 0 {
		od.Amount = in.Amount
	}
	if in.Status != 0 {
		od.Status = in.Status
	}
	if err = l.svcCtx.OrderModel.Update(l.ctx, od); err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}
	return &order.UpdateResponse{}, nil
}
