package middleware

import (
	"fmt"
	"net/http"

	"github.com/qnfnypen/gzocomm/response"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/types"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// AuthMiddleware 鉴权中间件
type AuthMiddleware struct {
	rcli *redis.Redis
}

// NewAuthMiddleware 初始化鉴权中间件
func NewAuthMiddleware(cli *redis.Redis) *AuthMiddleware {
	return &AuthMiddleware{
		rcli: cli,
	}
}

// Handle 实现鉴权中间件
func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuid, _ := r.Context().Value("uuid").(string)
		key := fmt.Sprintf("%s_%s", types.TokenPre, uuid)

		// 判断用户是否已退出登陆
		value, _ := m.rcli.GetCtx(r.Context(), key)
		if value == "" {
			response.UnauthorizedResp(w, "Login time out, please login again")
			return
		}

		// Passthrough to next handler if need
		next(w, r)
	}
}
