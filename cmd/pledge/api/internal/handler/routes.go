// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	auth "github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/handler/auth"
	user "github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/handler/user"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.HeaderMiddleware},
			[]rest.Route{
				{
					// 钱包的随机码
					Method:  http.MethodGet,
					Path:    "/code",
					Handler: auth.GetCodeHandler(serverCtx),
				},
				{
					// 用户登陆
					Method:  http.MethodPost,
					Path:    "/login",
					Handler: auth.UserLoginHandler(serverCtx),
				},
			}...,
		),
		rest.WithPrefix("/api/pledge"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.AuthMiddleware, serverCtx.HeaderMiddleware},
			[]rest.Route{
				{
					// 质押token
					Method:  http.MethodPost,
					Path:    "/delegate",
					Handler: user.DelegateHandler(serverCtx),
				},
				{
					// 获取进行中的解除质押
					Method:  http.MethodGet,
					Path:    "/delegate/unbinding",
					Handler: user.GetUnbindingDelegateHandler(serverCtx),
				},
				{
					// 获取用户信息详情
					Method:  http.MethodGet,
					Path:    "/info",
					Handler: user.InfoHandler(serverCtx),
				},
				{
					// 用户登出
					Method:  http.MethodPut,
					Path:    "/login_out",
					Handler: user.LoginoutHandler(serverCtx),
				},
				{
					// 质押转移
					Method:  http.MethodPut,
					Path:    "/redelegate",
					Handler: user.ReDelegateHandler(serverCtx),
				},
				{
					// 提取收益
					Method:  http.MethodPost,
					Path:    "/rewards/withdraw",
					Handler: user.WithdrawRewardsHandler(serverCtx),
				},
				{
					// 解除质押
					Method:  http.MethodPut,
					Path:    "/undelegate",
					Handler: user.UnDelegateHandler(serverCtx),
				},
				{
					// 取消解除质押
					Method:  http.MethodPut,
					Path:    "/undelegate/cancel",
					Handler: user.CancelUnDelegateHandler(serverCtx),
				},
				{
					// 获取验证者信息
					Method:  http.MethodGet,
					Path:    "/validators",
					Handler: user.ValidatorsHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/api/pledge/user"),
	)
}