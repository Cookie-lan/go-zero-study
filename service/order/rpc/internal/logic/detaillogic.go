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

type DetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DetailLogic) Detail(in *order.DetailRequest) (*order.DetailResponse, error) {
	od, err := l.svcCtx.OrderModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(http.StatusContinue, "订单不存在")
		}
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &order.DetailResponse{
		Id:     od.Id,
		Uid:    od.Uid,
		Pid:    od.Pid,
		Amount: od.Amount,
		Status: od.Status,
	}, nil
}
