package logic

import (
	"context"
	"go-zero-study/service/order/rpc/orderclient"

	"go-zero-study/service/order/api/internal/svc"
	"go-zero-study/service/order/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.DetailRequest) (resp *types.DetailResponse, err error) {
	od, err := l.svcCtx.OrderRpc.Detail(l.ctx, &orderclient.DetailRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	return &types.DetailResponse{
		Id:     od.Id,
		Uid:    od.Uid,
		Pid:    od.Pid,
		Amount: od.Amount,
		Status: od.Status,
	}, nil
}
