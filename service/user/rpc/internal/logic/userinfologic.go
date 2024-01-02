package logic

import (
	"context"
	"go-zero-study/service/user/model"
	"google.golang.org/grpc/status"
	"net/http"

	"go-zero-study/service/user/rpc/internal/svc"
	"go-zero-study/service/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoLogic) UserInfo(in *user.UserInfoRequest) (*user.UserInfoResponse, error) {
	// 查询用户是否存在
	u, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(http.StatusContinue, "该用户不存在")
		}
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &user.UserInfoResponse{
		Id:     u.Id,
		Name:   u.Name,
		Mobile: u.Mobile,
		Gender: u.Gender,
	}, nil
}
