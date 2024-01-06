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

func (l *DetailLogic) Detail(in *product.DetailRequest) (*product.DetailResponse, error) {
	produce, err := l.svcCtx.ProductModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Errorf(http.StatusContinue, err.Error())
		}
		return nil, status.Errorf(http.StatusInternalServerError, err.Error())
	}

	return &product.DetailResponse{
		Id:     produce.Id,
		Name:   produce.Name,
		Desc:   produce.Desc,
		Stock:  produce.Stock,
		Amount: produce.Amount,
		Status: produce.Status,
	}, nil
}
