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

type PaidLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPaidLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaidLogic {
	return &PaidLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PaidLogic) Paid(in *order.PaidRequest) (*order.PaidResponse, error) {
	// 查询订单是否存在
	od, err := l.svcCtx.OrderModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(http.StatusContinue, "订单不存在")
		}
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	od.Status = 1

	if err = l.svcCtx.OrderModel.Update(l.ctx, od); err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}
	return &order.PaidResponse{}, nil
}
