package logic

import (
	"context"
	"go-zero-study/service/order/rpc/orderclient"

	"go-zero-study/service/order/api/internal/svc"
	"go-zero-study/service/order/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveLogic {
	return &RemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveLogic) Remove(req *types.RemoveRequest) (resp *types.RemoveResponse, err error) {
	_, err = l.svcCtx.OrderRpc.Remove(l.ctx, &orderclient.RemoveRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	return &types.RemoveResponse{}, nil
}
