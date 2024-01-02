package logic

import (
	"context"
	"encoding/json"
	"go-zero-study/service/user/rpc/userclient"

	"go-zero-study/service/user/api/internal/svc"
	"go-zero-study/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo() (resp *types.UserInfoResponse, err error) {
	// 获取用户信息
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	u, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &userclient.UserInfoRequest{
		Id: uid,
	})
	if err != nil {
		return
	}

	return &types.UserInfoResponse{
		Id:     u.Id,
		Name:   u.Name,
		Gender: u.Gender,
		Mobile: u.Mobile,
	}, nil
}
