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

func (l *RemoveLogic) Remove(in *product.RemoveRequest) (*product.RemoveResponse, error) {
	res, err := l.svcCtx.ProductModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Errorf(http.StatusContinue, err.Error())
		}
		return nil, status.Errorf(http.StatusInternalServerError, err.Error())
	}

	err = l.svcCtx.ProductModel.Delete(l.ctx, res.Id)
	if err != nil {
		return nil, status.Errorf(http.StatusInternalServerError, err.Error())
	}
	return &product.RemoveResponse{}, nil
}
