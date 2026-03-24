package middleware

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kratos/kratos/v2/transport"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

// Context keys for user information
type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UsernameKey contextKey = "username"
)

// UserContext 用户上下文信息
type UserContext struct {
	UserID   int64
	Username string
}

// NewUserContextFromRequest 从请求中提取用户信息
// 支持两种方式：
// 1. 通过 API 网关转发的请求头（X-User-ID, X-Username）
// 2. 直接解析 JWT Token（适用于独立部署）
func NewUserContextFromRequest(ctx context.Context) (*UserContext, error) {
	tr, ok := transport.FromServerContext(ctx)
	if !ok {
		return nil, errors.New("no transport found in context")
	}

	// 方式1: 从 HTTP Header 获取（API 网关转发）
	if httpTr, ok := tr.(kratoshttp.Transporter); ok {
		userID := httpTr.Request().Header.Get("X-User-ID")
		username := httpTr.Request().Header.Get("X-Username")

		if userID != "" {
			return &UserContext{
				UserID:   parseInt64(userID),
				Username: username,
			}, nil
		}
	}

	// 方式2: 从 JWT Token 解析（需要实现 JWT 中间件）
	// return extractFromJWT(ctx)

	return nil, errors.New("user not authenticated")
}

// GetUserIDFromContext 从 context 中获取用户 ID
func GetUserIDFromContext(ctx context.Context) int64 {
	if userID, ok := ctx.Value(UserIDKey).(int64); ok {
		return userID
	}
	return 0
}

// NewContextWithUser 创建带有用户信息的 context
func NewContextWithUser(ctx context.Context, user *UserContext) context.Context {
	ctx = context.WithValue(ctx, UserIDKey, user.UserID)
	ctx = context.WithValue(ctx, UsernameKey, user.Username)
	return ctx
}

func parseInt64(s string) int64 {
	var i int64
	if _, err := fmt.Sscanf(s, "%d", &i); err != nil {
		return 0
	}
	return i
}
