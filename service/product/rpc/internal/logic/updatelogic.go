package logic

import (
	"context"
	"go-zero-study/service/product/model"
	"google.golang.org/grpc/status"
	"net/http"

	"go-zero-study/service/product/rpc/internal/svc"
	"go-zero-study/service/product/rpc/types/product"

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

func (l *UpdateLogic) Update(in *product.UpdateRequest) (*product.UpdateResponse, error) {
	res, err := l.svcCtx.ProductModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Errorf(http.StatusContinue, err.Error())
		}
		return nil, status.Errorf(http.StatusInternalServerError, err.Error())
	}

	if in.Name != "" {
		res.Name = in.Name
	}
	if in.Desc != "" {
		res.Desc = in.Desc
	}
	if in.Stock != 0 {
		res.Stock = in.Stock
	}
	if in.Amount != 0 {
		res.Amount = in.Amount
	}
	if in.Status != 0 {
		res.Status = in.Status
	}

	err = l.svcCtx.ProductModel.Update(l.ctx, res)
	if err != nil {
		return nil, status.Errorf(http.StatusInternalServerError, err.Error())
	}
	return &product.UpdateResponse{}, nil
}
