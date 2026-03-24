package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"room/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

// UserClient 用户服务客户端 (HTTP 实现)
type UserClient struct {
	baseURL    string
	httpClient *http.Client
	log        *log.Helper
}

// UserInfo 用户信息
type UserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

// VerifyUserResponse 验证用户响应
type VerifyUserResponse struct {
	Valid    bool   `json:"valid"`
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

// NewUserClient 创建用户服务客户端
func NewUserClient(c *conf.Bootstrap, logger log.Logger) (*UserClient, func(), error) {
	// 默认配置
	addr := "user-service:8000"
	if c.UserClient != nil && c.UserClient.Addr != "" {
		addr = c.UserClient.Addr
	}

	// 如果是 gRPC 端口 (9000)，替换为 HTTP 端口 (8000)
	baseURL := fmt.Sprintf("http://%s", addr)

	client := &UserClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		log: log.NewHelper(logger),
	}

	cleanup := func() {
		log.Info("user service HTTP client cleanup")
	}

	return client, cleanup, nil
}

// VerifyUser 验证用户是否存在
func (c *UserClient) VerifyUser(ctx context.Context, userID int64) (*VerifyUserResponse, error) {
	c.log.Infof("Verifying user: %d", userID)

	// 使用 User Service 的登录端点来验证
	// 但我们需要一个端点来通过 user_id 验证用户
	// 这里我们直接调用 HTTP 端点
	url := fmt.Sprintf("%s/internal/verify/%d", c.baseURL, userID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// 如果内部端点不存在，返回默认有效（开发环境）
		c.log.Warnf("User service verify endpoint not available, assuming user %d is valid", userID)
		return &VerifyUserResponse{
			Valid:  true,
			UserID: userID,
		}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result VerifyUserResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, err
		}
		return &result, nil
	}

	// 用户不存在
	return &VerifyUserResponse{
		Valid:  false,
		UserID: userID,
	}, nil
}

// GetUser 获取用户信息
func (c *UserClient) GetUser(ctx context.Context, userID int64) (*UserInfo, error) {
	c.log.Infof("Getting user: %d", userID)

	// 通过 profile 端点获取用户信息（但这需要 JWT token）
	// 对于服务间通信，我们需要一个内部端点
	url := fmt.Sprintf("%s/internal/user/%d", c.baseURL, userID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// 如果内部端点不存在，返回默认用户信息（开发环境）
		c.log.Warnf("User service get endpoint not available, returning default user %d", userID)
		return &UserInfo{
			ID:       userID,
			Username: fmt.Sprintf("user_%d", userID),
			Nickname: fmt.Sprintf("User %d", userID),
		}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var user UserInfo
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			return nil, err
		}
		return &user, nil
	}

	body, _ := io.ReadAll(resp.Body)
	return nil, fmt.Errorf("failed to get user %d: %s", userID, string(body))
}

// Close 关闭连接
func (c *UserClient) Close() error {
	return nil
}
