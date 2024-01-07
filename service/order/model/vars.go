package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var (
	ErrNotFound = sqlx.ErrNotFound

	OrderStatusNormal int64 = 0
	// OrderStatusInvalid 订单失效
	OrderStatusInvalid int64 = 9
)
