package logic

import (
	"context"
	"go-zero-study/common/cryptx"
	"go-zero-study/service/user/model"
	"google.golang.org/grpc/status"
	"net/http"

	"go-zero-study/service/user/rpc/internal/svc"
	"go-zero-study/service/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginRequest) (*user.LoginResponse, error) {
	// 查询用户是否存在
	u, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, in.Mobile)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(http.StatusContinue, "该用户不存在")
		}
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	// 判断密码是否正确
	password := cryptx.PasswordEncrypt(l.svcCtx.Config.Salt, in.Password)
	if password != u.Password {
		return nil, status.Error(http.StatusContinue, "密码错误")
	}

	return &user.LoginResponse{
		Id:     u.Id,
		Name:   u.Name,
		Mobile: u.Mobile,
		Gender: u.Gender,
	}, nil
}
