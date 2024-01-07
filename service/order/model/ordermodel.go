package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderModel = (*customOrderModel)(nil)

type (
	// OrderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderModel.
	OrderModel interface {
		orderModel
	}

	customOrderModel struct {
		*defaultOrderModel
	}
)

func (m *customOrderModel) TxInsert(ctx context.Context, tx *sql.Tx, data *Order) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, orderRowsExpectAutoSet)
	return tx.ExecContext(ctx, query, data.Uid, data.Pid, data.Amount, data.Status)
}

func (m *customOrderModel) TxUpdate(ctx context.Context, tx *sql.Tx, data *Order) error {
	key := fmt.Sprintf("%s%v", cacheOrderIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (sql.Result, error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, orderRowsWithPlaceHolder)
		return tx.ExecContext(ctx, query, data.Uid, data.Pid, data.Amount, data.Status, data.Id)
	}, key)
	return err
}

func (m *customOrderModel) FindOneByUid(ctx context.Context, uid int64) (*Order, error) {
	var out Order
	query := fmt.Sprintf("select %s from %s where `uid` = ? order by create_time desc limit 1", orderRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &out, query, uid)

	switch err {
	case nil:
		return &out, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customOrderModel) FindAllByUID(ctx context.Context, uid int64) ([]*Order, error) {
	var out []*Order
	query := fmt.Sprintf("select %s from %s where `uid` = ?", orderRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &out, query, uid)

	switch err {
	case nil:
		return out, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// NewOrderModel returns a model for the database table.
func NewOrderModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OrderModel {
	return &customOrderModel{
		defaultOrderModel: newOrderModel(conn, c, opts...),
	}
}
