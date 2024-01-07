package logic

import (
	"context"
	"database/sql"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"

	"go-zero-study/service/product/rpc/internal/svc"
	"go-zero-study/service/product/rpc/types/product"

	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DecrStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDecrStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DecrStockLogic {
	return &DecrStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DecrStockLogic) DecrStock(in *product.DecrStockRequest) (*product.DecrStockResponse, error) {
	// 获取 RawDB
	db, err := sqlx.NewMysql(l.svcCtx.Config.Mysql.DataSource).RawDB()
	if err != nil {
		return nil, status.Errorf(http.StatusInternalServerError, err.Error())
	}

	// 获取子事务屏障对象
	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		return nil, status.Errorf(http.StatusInternalServerError, err.Error())
	}

	// 开启自事务屏障
	err = barrier.CallWithDB(db, func(tx *sql.Tx) error {
		// 更新后产品库存
		result, err := l.svcCtx.ProductModel.TxAdjustStock(l.ctx, tx, in.Id, -1)
		if err != nil {
			return err
		}
		affected, err := result.RowsAffected()
		if err == nil && affected == 0 { // 库存不足，返回子事务失败
			return dtmcli.ErrFailure
		}

		return err
	})
	if err == dtmcli.ErrFailure { // 这种情况是库存不足，不再重试，走回滚
		return nil, status.Errorf(codes.Aborted, dtmcli.ResultFailure)
	}

	if err != nil { // 其他错误，返回错误，dtm 会重试
		return nil, status.Errorf(http.StatusInternalServerError, err.Error())
	}

	return &product.DecrStockResponse{}, nil
}
