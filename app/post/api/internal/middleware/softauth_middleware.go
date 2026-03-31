package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// SoftAuthMiddleware 弱认证中间件：尝试解析 JWT 并注入 Context，但不强制阻拦游客
type SoftAuthMiddleware struct {
	Secret string
}

func NewSoftAuthMiddleware(secret string) *SoftAuthMiddleware {
	return &SoftAuthMiddleware{Secret: secret}
}

func (m *SoftAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. 尝试获取 Authorization 头
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// 没带 Token，当做游客直接放行
			next(w, r)
			return
		}

		// 2. 验证 Bearer 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			next(w, r)
			return
		}

		tokenString := parts[1]

		// 3. 解析并校验 Token
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.Secret), nil
		})

		// 4. 如果 Token 合法，将其携带的信息 (如 userId) 注入当前请求的 Context
		if token != nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				ctx := r.Context()
				for k, v := range claims {
					ctx = context.WithValue(ctx, k, v)
				}
				// 替换 Request 的 Context
				r = r.WithContext(ctx)
			}
		}

		// 5. 放行请求
		next(w, r)
	}
}
