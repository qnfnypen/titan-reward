package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/qnfnypen/titan-reward/cmd/reward/api/internal/types"
)

// HeaderMiddleware 请求头处理全局中间键
type HeaderMiddleware struct {
}

// NewHeaderMiddleware 新建 请求头处理全局中间键
func NewHeaderMiddleware() *HeaderMiddleware {
	return &HeaderMiddleware{}
}

// Handle 实现 请求头处理全局中间键
func (m *HeaderMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取语言信息
		ctx := r.Context()
		language := r.Header.Get("lang")
		ctx = context.WithValue(ctx, types.LangKey, strings.TrimSpace(language))
		r = r.WithContext(ctx)
		next(w, r)
	}
}
