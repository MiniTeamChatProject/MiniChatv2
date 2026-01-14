package server

import (
	"context"
	"errors"
	"strings"

	// 改回 registration
	v1 "registration/api/helloworld/v1"
	"registration/internal/conf"
	"registration/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/golang-jwt/jwt/v5"
)

// 定义密钥 (必须和你 Biz 层里写的一样)
var jwtSecret = []byte("MySecretKey_0721")

// 1. 定义 JWT 认证中间件
func AuthMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			operation := tr.Operation()
			if isWhiteList(operation) {
				return handler(ctx, req)
			}

			authHeader := tr.RequestHeader().Get("Authorization")
			if authHeader == "" {
				return nil, errors.New("未提供 Token")
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return nil, errors.New("Token 格式错误")
			}
			tokenString := parts[1]

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				return nil, errors.New("Token 无效或已过期")
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if userIdFloat, ok := claims["user_id"].(float64); ok {
					ctx = context.WithValue(ctx, "user_id", int64(userIdFloat))
				}
			}

			return handler(ctx, req)
		}
	}
}

func isWhiteList(operation string) bool {
	if strings.Contains(operation, "Register") || strings.Contains(operation, "Login") {
		return true
	}
	return false
}

// 2. HTTP Server 初始化
func NewHTTPServer(c *conf.Server, greeter *service.RegistrationService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			AuthMiddleware(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterRegistrationHTTPServer(srv, greeter)
	return srv
}
