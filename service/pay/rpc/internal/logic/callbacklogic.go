package logic

import (
	"context"
	"go-zero-study/service/order/rpc/types/order"
	"go-zero-study/service/pay/model"
	"go-zero-study/service/user/rpc/types/user"
	"google.golang.org/grpc/status"
	"net/http"

	"go-zero-study/service/pay/rpc/internal/svc"
	"go-zero-study/service/pay/rpc/types/pay"

	"github.com/zeromicro/go-zero/core/logx"
)

type CallbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CallbackLogic {
	return &CallbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CallbackLogic) Callback(in *pay.CallbackRequest) (*pay.CallbackResponse, error) {
	// 查询用户是否存在
	_, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user.UserInfoRequest{
		Id: in.Uid,
	})
	if err != nil {
		return nil, err
	}

	// 查询订单是否存在
	_, err = l.svcCtx.OrderRpc.Detail(l.ctx, &order.DetailRequest{
		Id: in.Oid,
	})
	if err != nil {
		return nil, err
	}

	// 查询支付是否存在
	res, err := l.svcCtx.PayModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(http.StatusContinue, "支付不存在")
		}
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	// 支付金额与订单金额不符
	if res.Amount != in.Amount {
		return nil, status.Error(http.StatusContinue, "支付金额与订单金额不符")
	}

	res.Source = in.Source
	res.Status = in.Status
	if err = l.svcCtx.PayModel.Update(l.ctx, res); err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	// 更新订单支付状态
	if _, err = l.svcCtx.OrderRpc.Paid(l.ctx, &order.PaidRequest{
		Id: res.Oid,
	}); err != nil {
		return nil, err
	}
	return &pay.CallbackResponse{}, nil
}
