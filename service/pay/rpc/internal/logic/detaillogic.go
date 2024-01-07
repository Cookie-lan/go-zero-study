package logic

import (
	"context"
	"go-zero-study/service/pay/model"
	"google.golang.org/grpc/status"
	"net/http"

	"go-zero-study/service/pay/rpc/internal/svc"
	"go-zero-study/service/pay/rpc/types/pay"

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

func (l *DetailLogic) Detail(in *pay.DetailRequest) (*pay.DetailResponse, error) {
	// 查询支付是否存在
	res, err := l.svcCtx.PayModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(http.StatusContinue, "支付不存在")
		}
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &pay.DetailResponse{
		Id:     res.Id,
		Uid:    res.Uid,
		Oid:    res.Oid,
		Amount: res.Amount,
		Source: res.Source,
		Status: res.Status,
	}, nil
}
