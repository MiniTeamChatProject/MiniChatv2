package data

import (
	"context"
	"fmt"
	"time"

	userv1 "room/api/user/v1"
	"room/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

// UserClient 用户服务客户端
type UserClient struct {
	client userv1.UserServiceClient
	conn   *grpc.ClientConn
	log    *log.Helper
}

// NewUserClient 创建用户服务客户端
func NewUserClient(c *conf.Bootstrap, logger log.Logger) (*UserClient, func(), error) {
	// 默认配置
	addr := "localhost:9000"
	if c.UserClient != nil && c.UserClient.Addr != "" {
		addr = c.UserClient.Addr
	}

	// 配置 gRPC 连接
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				MaxDelay: 3 * time.Second,
			},
		}),
	}

	// 建立连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to user service at %s: %w", addr, err)
	}

	client := &UserClient{
		client: userv1.NewUserServiceClient(conn),
		conn:   conn,
		log:    log.NewHelper(logger),
	}

	cleanup := func() {
		log.Info("closing user service gRPC connection")
		client.Close()
	}

	return client, cleanup, nil
}

// VerifyUser 验证用户是否存在
func (c *UserClient) VerifyUser(ctx context.Context, userID int64) (*userv1.VerifyUserResponse, error) {
	c.log.Infof("Verifying user: %d", userID)

	resp, err := c.client.VerifyUser(ctx, &userv1.VerifyUserRequest{
		UserId: userID,
	})
	if err != nil {
		c.log.Errorf("Failed to verify user %d: %v", userID, err)
		return nil, err
	}

	if !resp.Valid {
		return nil, fmt.Errorf("user %d not found or invalid", userID)
	}

	return resp, nil
}

// GetUser 获取用户信息
func (c *UserClient) GetUser(ctx context.Context, userID int64) (*userv1.GetUserResponse, error) {
	c.log.Infof("Getting user: %d", userID)

	resp, err := c.client.GetUser(ctx, &userv1.GetUserRequest{
		Id: userID,
	})
	if err != nil {
		c.log.Errorf("Failed to get user %d: %v", userID, err)
		return nil, err
	}

	return resp, nil
}

// Close 关闭连接
func (c *UserClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
