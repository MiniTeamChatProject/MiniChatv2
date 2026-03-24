package middleware

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// Auth 身份验证中间件
// 从请求头中提取用户信息并注入到 context 中
// 支持白名单模式，公开端点不需要认证
func Auth(whitelist ...string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 检查是否在白名单中
			if isInWhitelist(ctx, whitelist) {
				// 白名单端点，尝试提取用户信息但不强制要求
				user, _ := NewUserContextFromRequest(ctx)
				if user != nil {
					ctx = NewContextWithUser(ctx, user)
				}
				return handler(ctx, req)
			}

			// 非白名单端点，必须认证
			user, err := NewUserContextFromRequest(ctx)
			if err != nil {
				return nil, err
			}

			// 将用户信息注入到 context 中
			ctx = NewContextWithUser(ctx, user)

			// 继续处理请求
			return handler(ctx, req)
		}
	}
}

// isInWhitelist 检查当前请求是否在白名单中
func isInWhitelist(ctx context.Context, whitelist []string) bool {
	if len(whitelist) == 0 {
		return false
	}

	tr, ok := transport.FromServerContext(ctx)
	if !ok {
		return false
	}

	operation := tr.Operation()
	for _, path := range whitelist {
		if strings.Contains(operation, path) {
			return true
		}
	}

	return false
}

// OptionalAuth 可选的身份验证中间件
// 如果有用户信息则注入，没有则继续处理
func OptionalAuth() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			user, err := NewUserContextFromRequest(ctx)
			if err == nil && user != nil {
				ctx = NewContextWithUser(ctx, user)
			}
			return handler(ctx, req)
		}
	}
}
