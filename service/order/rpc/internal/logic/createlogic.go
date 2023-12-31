package logic

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dtm-labs/client/dtmgrpc"
	"net/http"

	"go-zero-study/service/order/model"
	"go-zero-study/service/order/rpc/internal/svc"
	"go-zero-study/service/order/rpc/types/order"
	"go-zero-study/service/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/status"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateLogic) Create(in *order.CreateRequest) (*order.CreateResponse, error) {
	// 获取 RawDB
	db, err := sqlx.NewMysql(l.svcCtx.Config.Mysql.DataSource).RawDB()
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	// 获取子事务屏障对象
	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	newOrder := model.Order{
		Uid:    in.Uid,
		Pid:    in.Pid,
		Amount: in.Amount,
		Status: model.OrderStatusNormal,
	}
	// 开启自事务屏障
	if err = barrier.CallWithDB(db, func(tx *sql.Tx) error {
		// 查询用户是否存在
		_, err = l.svcCtx.UserRpc.UserInfo(l.ctx, &user.UserInfoRequest{
			Id: in.Uid,
		})
		if err != nil {
			return fmt.Errorf("用户不存在")
		}

		// 创建订单
		res, err := l.svcCtx.OrderModel.TxInsert(l.ctx, tx, &newOrder)
		if err != nil {
			return fmt.Errorf("订单创建失败")
		}

		newOrder.Id, err = res.LastInsertId()
		return err
	}); err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &order.CreateResponse{
		Id: newOrder.Id,
	}, nil
}
