// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"go-zero-study/service/order/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/orders/create",
				Handler: CreateHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/orders/update",
				Handler: UpdateHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/api/orders/:id",
				Handler: RemoveHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/orders/:id",
				Handler: DetailHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/orders",
				Handler: ListHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)
}
