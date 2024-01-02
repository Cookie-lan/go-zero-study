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

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterRequest) (*user.RegisterResponse, error) {
	// 判断手机号是否已经注册
	_, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, in.Mobile)
	if err == nil {
		return nil, status.Error(http.StatusContinue, "该用户已存在")
	}

	if err != model.ErrNotFound {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	newUser := model.User{
		Name:     in.Name,
		Gender:   in.Gender,
		Mobile:   in.Mobile,
		Password: cryptx.PasswordEncrypt(l.svcCtx.Config.Salt, in.Password),
	}

	res, err := l.svcCtx.UserModel.Insert(l.ctx, &newUser)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	newUser.Id, err = res.LastInsertId()
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &user.RegisterResponse{
		Id:     newUser.Id,
		Name:   newUser.Name,
		Gender: newUser.Gender,
		Mobile: newUser.Mobile,
	}, nil
}
