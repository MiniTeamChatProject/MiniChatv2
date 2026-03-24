# 用户系统集成指南

## 架构概述

```
┌─────────┐      ┌──────────┐      ┌──────────┐      ┌──────────┐
│ 客户端   │ ──> │ API 网关  │ ──> │ 用户服务  │      │ 房间服务  │
│(Web/App)│      │(JWT验证) │      │          │      │          │
└─────────┘      └──────────┘      └──────────┘      └──────────┘
                       │                                  ^
                       │ 添加用户信息头                     │
                       │ X-User-ID, X-Username            │
                       └──────────────────────────────────┘
```

## 方案一：通过 API 网关（推荐）

### 1. API 网关配置

API 网关负责验证 JWT Token，并将用户信息通过请求头转发给后端服务：

```yaml
# API Gateway 配置示例 (Nginx/Kong/Traefik)
location /api/rooms/ {
    # JWT 验证配置
    auth_jwt /path/to/jwt/key;

    # 转发用户信息到后端
    proxy_set_header X-User-ID $jwt_claim_user_id;
    proxy_set_header X-Username $jwt_claim_username;
    proxy_set_header X-User-Role $jwt_claim_role;

    proxy_pass http://room-service:8000/;
}
```

### 2. 房间服务配置

在 `internal/server/http.go` 中添加 Auth 中间件：

```go
import (
    "room/internal/middleware"
    "github.com/go-kratos/kratos/v2/middleware/auth"
)

func NewHTTPServer(c *conf.Server, roomSvc *service.RoomService, wsSvc *service.WebSocketService, logger log.Logger) *http.Server {
    srv := http.NewServer(opts...)

    // 添加 Auth 中间件
    srv.Route("/v1").Use(middleware.Auth())

    v1.RegisterRoomServiceHTTPServer(srv, roomSvc)
    // ...
}
```

### 3. 请求流程

```
1. 客户端请求: POST /v1/rooms
   Header: Authorization: Bearer <JWT_TOKEN>

2. API 网关验证 JWT，解析用户信息
   user_id: 123, username: "alice"

3. 网关转发请求到房间服务
   Header: X-User-ID: 123
           X-Username: alice

4. 房间服务从 Header 提取用户信息，注入到 context

5. 业务逻辑从 context 获取 user_id
   userID := middleware.GetUserIDFromContext(ctx)
```

## 方案二：直接验证 JWT

如果房间服务需要独立验证 JWT，可以实现完整的 JWT 中间件：

```go
// internal/middleware/jwt.go
package middleware

import (
    "github.com/golang-jwt/jwt/v4"
)

type JWTConfig struct {
    Secret string
    Issuer string
}

func JWTAuth(config *JWTConfig) middleware.Middleware {
    return func(handler middleware.Handler) middleware.Handler {
        return func(ctx context.Context, req interface{}) (interface{}, error) {
            // 获取 Authorization header
            token := extractToken(ctx)

            // 验证 JWT
            claims, err := validateJWT(token, config)
            if err != nil {
                return nil, err
            }

            // 注入用户信息到 context
            user := &UserContext{
                UserID:   claims.UserID,
                Username: claims.Username,
            }
            ctx = NewContextWithUser(ctx, user)

            return handler(ctx, req)
        }
    }
}
```

## 方案三：服务间通信（gRPC）

如果需要验证用户是否真实存在，可以调用用户服务：

```go
// internal/biz/user_client.go
package biz

type UserClient interface {
    ValidateUser(ctx context.Context, userID int64) (*User, error)
}

type userClient struct {
    conn *grpc.ClientConn
}

func (c *userClient) ValidateUser(ctx context.Context, userID int64) (*User, error) {
    // 调用用户服务的 gRPC 接口
    client := pb.NewUserServiceClient(c.conn)
    return client.GetUser(ctx, &pb.GetUserRequest{Id: userID})
}
```

## 测试

### 使用 curl 测试（带用户信息头）

```bash
# 创建房间（模拟 API 网关转发）
curl -X POST http://localhost:8000/v1/rooms \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 123" \
  -H "X-Username: alice" \
  -d '{"name": "Alice的聊天室"}'
```

### 使用 wscat 测试 WebSocket

```bash
# WebSocket 连接需要在 URL 参数或子协议中传递用户信息
wscat -c "ws://localhost:8000/ws/room/1?user_id=123&username=alice"
```

**注意**：当前 WebSocket 实现使用硬编码 user_id，需要更新以支持从连接参数获取用户信息：

```go
// internal/service/websocket.go
func (s *WebSocketService) HandleWebSocket(ctx http.Context) error {
    // 从 URL 参数获取用户信息
    userIDStr := ctx.Query().Get("user_id")
    userID := parseInt64(userIDStr)

    // 或从 Header 获取（用于 WebSocket 握手）
    userID = parseInt64(ctx.Request().Header.Get("X-User-ID"))
    // ...
}
```

## 生产环境注意事项

1. **HTTPS**: 生产环境必须使用 HTTPS
2. **JWT Secret**: 使用环境变量或密钥管理服务存储
3. **Token 过期**: 设置合理的过期时间
4. **刷新 Token**: 实现 Token 刷新机制
5. **日志**: 记录用户操作用于审计
6. **限流**: 防止恶意请求
