package middleware

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
)

// Auth 身份验证中间件
// 从请求头中提取用户信息并注入到 context 中
func Auth() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 从请求中提取用户信息
			user, err := NewUserContextFromRequest(ctx)
			if err != nil {
				// 如果没有用户信息，返回错误或使用默认值
				// 根据业务需求，可以选择：
				// 1. 返回未授权错误
				// 2. 使用默认用户（开发环境）
				// 3. 继续处理（允许匿名访问的接口）
				return nil, err
			}

			// 将用户信息注入到 context 中
			ctx = NewContextWithUser(ctx, user)

			// 继续处理请求
			return handler(ctx, req)
		}
	}
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
