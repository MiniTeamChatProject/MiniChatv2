# MiniChat Room Service

> 基于 [Kratos v2](https://github.com/go-kratos/kratos) 框架的微服务房间管理系统，提供房间管理、成员权限控制、实时消息推送等功能。

## 📋 目录

- [功能特性](#功能特性)
- [架构设计](#架构设计)
- [快速开始](#快速开始)
- [API 文档](#api-文档)
- [WebSocket 接口](#websocket-接口)
- [数据模型](#数据模型)
- [开发指南](#开发指南)
- [部署指南](#部署指南)
- [故障排查](#故障排查)

---

## 功能特性

### 核心功能
- ✅ **房间管理**：创建、查询、更新、删除房间
- ✅ **成员管理**：加入、离开、踢出成员，角色权限控制
- ✅ **实时消息**：WebSocket 双向通信，消息广播
- ✅ **消息历史**：分页查询历史消息
- ✅ **用户集成**：与 User Service 的 HTTP/gRPC 通信

### 权限模型

| 操作 | 房主 | 管理员 | 普通成员 |
|------|------|--------|----------|
| 更新房间 | ✓ | ✓ | ✗ |
| 删除房间 | ✓ | ✗ | ✗ |
| 踢出成员 | ✓ | ✓（不能踢管理员/房主） | ✗ |
| 更新角色 | ✓ | ✗ | ✗ |
| 禁言/解禁 | ✓ | ✓ | ✗ |
| 发送消息 | ✓ | ✓ | ✓ |

---

## 架构设计

### 分层架构

```
┌─────────────────────────────────────────────────────────────┐
│                    Room Service                             │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐        ┌──────────────┐                   │
│  │ API Layer    │        │              │                   │
│  │ (proto/api)  │        │              │                   │
│  │              │        │              │                   │
│  └──────┬───────┘        │              │                   │
│         │                │              │                   │
│  ┌──────▼──────────────┐ │              │                   │
│  │ Service Layer       │ │              │                   │
│  │ (internal/service)  │ │              │                   │
│  │                     │ │              │                   │
│  │ - RoomService       │ │              │                   │
│  │ - WebSocketService  │ │              │                   │
│  └──────┬──────────────┘ │              │                   │
│         │                │              │                   │
│  ┌──────▼──────────────┐ │              │                   │
│  │ Business Logic      │ │              │                   │
│  │ (internal/biz)      │ │              │                   │
│  │                     │ │              │                   │
│  │ - RoomUsecase       │ │              │                   │
│  │ - RoomMemberUsecase │ │              │                   │
│  │ - MessageUsecase    │ │              │                   │
│  └──────┬──────────────┘ │              │                   │
│         │                │              │                   │
│  ┌──────▼──────────────┐ │              │                   │
│  │ Data Layer          │ │              │                   │
│  │ (internal/data)     │ │              │                   │
│  │                     │ │              │                   │
│  │ - RoomRepository    │ │              │                   │
│  │ - RoomMemberRepo    │ │              │                   │
│  │ - MessageRepository │ │              │                   │
│  └──────┬──────────────┘ │              │                   │
│         │                │              │                   │
│  ┌──────▼──────────────┐ │              │                   │
│  │ PostgreSQL          │ │              │                   │
│  └─────────────────────┘ │              │                   │
│                           │              │                   │
│  User Service ◄──────────┘              │                   │
│  (HTTP/gRPC)                              │                   │
└───────────────────────────────────────────┘                   │
```

### 目录结构

```
service/room/
├── api/                    # Protobuf API 定义
│   └── room/v1/
│       ├── room.proto      # 房间服务 API
│       └── room.pb.go      # 生成的 Go 代码
├── cmd/room/              # 应用入口
│   ├── main.go            # 主函数
│   └── wire.go            # Wire 依赖注入配置
├── configs/               # 配置文件
│   └── config.yaml        # 服务配置
├── internal/              # 内部代码
│   ├── biz/               # 业务逻辑层
│   │   ├── room.go        # 领域模型定义
│   │   ├── room_usecase.go
│   │   ├── room_member_usecase.go
│   │   └── message_usecase.go
│   ├── data/              # 数据访问层
│   │   ├── data.go        # Data 层初始化
│   │   ├── model/         # GORM 数据模型
│   │   ├── room_repo.go   # 房间仓储实现
│   │   ├── room_member_repo.go
│   │   ├── message_repo.go
│   │   └── user_client.go # 用户服务 HTTP 客户端
│   ├── service/           # 服务层
│   │   ├── room_service.go
│   │   └── websocket.go   # WebSocket 处理
│   ├── server/           # 服务器配置
│   │   ├── http.go       # HTTP 服务器
│   │   └── grpc.go       # gRPC 服务器
│   └── middleware/       # 中间件
│       └── auth.go       # 认证中间件
├── Dockerfile            # Docker 构建文件
├── Makefile             # 构建脚本
├── go.mod
└── go.sum
```

---

## 快速开始

### 环境要求

- Go 1.24+
- Docker & Docker Compose
- PostgreSQL 16
- Make（可选）

### 本地开发

```bash
# 1. 克隆项目
cd service/room

# 2. 安装依赖
go mod tidy

# 3. 启动数据库
docker-compose up -d room-db

# 4. 生成代码（修改 proto 后需要执行）
make api

# 5. 生成 Wire 依赖注入代码
make generate

# 6. 配置文件
cp configs/config.yaml.example configs/config.yaml
# 编辑 configs/config.yaml，配置数据库连接

# 7. 运行服务
make build
./bin/room -conf ./configs
```

### Docker 部署

```bash
# 启动所有服务（包括数据库）
docker-compose up -d

# 查看日志
docker-compose logs -f room-service

# 停止服务
docker-compose down
```

---

## API 文档

### HTTP 接口

| 方法 | 路径 | 说明 | 请求体 | 响应 |
|------|------|------|--------|------|
| POST | `/v1/rooms` | 创建房间 | `{name, type?, max_count?, ...}` | `{room}` |
| GET | `/v1/rooms/{id}` | 获取房间信息 | - | `{room}` |
| PUT | `/v1/rooms/{id}` | 更新房间 | `{name?, avatar?, ...}` | `{room}` |
| DELETE | `/v1/rooms/{id}` | 删除房间 | - | `{success}` |
| GET | `/v1/rooms` | 获取所有房间 | `page?, page_size?` | `{rooms, total}` |
| POST | `/v1/rooms/{id}/join` | 加入房间 | `{}` | `{member, room}` |
| POST | `/v1/rooms/{id}/leave` | 离开房间（状态变 Left） | `{}` | `{success}` |
| POST | `/v1/rooms/{id}/quit` | 退出群（删除记录） | `{}` | `{success}` |
| GET | `/v1/rooms/{id}/members` | 获取成员列表 | `page?, page_size?` | `{members, total}` |
| DELETE | `/v1/rooms/{id}/members/{user_id}` | 踢出成员 | - | `{success}` |
| PUT | `/v1/rooms/{id}/members/{user_id}/role` | 更新角色 | `{role}` | `{success}` |
| PUT | `/v1/rooms/{id}/members/{user_id}/mute` | 禁言/解禁 | `{mute_until}` | `{success}` |
| GET | `/v1/users/{user_id}/rooms` | 获取用户房间列表 | `page?, page_size?` | `{rooms, total}` |
| POST | `/v1/rooms/{id}/messages` | 发送消息 | `{content, type?}` | `{message}` |
| GET | `/v1/rooms/{id}/messages` | 获取消息历史 | `page?, page_size?` | `{messages, total}` |

### 认证方式

所有 API 请求需要在请求头中包含用户信息：

```http
X-User-ID: 123
Authorization: Bearer <jwt_token>
```

### 示例请求

```bash
# 创建房间
curl -X POST http://localhost:8001/v1/rooms \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{
    "name": "测试房间",
    "max_count": 100,
    "description": "这是一个测试房间"
  }'

# 加入房间
curl -X POST http://localhost:8001/v1/rooms/1/join \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 2"

# 发送消息
curl -X POST http://localhost:8001/v1/rooms/1/messages \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 2" \
  -d '{"content": "Hello!"}'

# 获取消息历史
curl http://localhost:8001/v1/rooms/1/messages \
  -H "X-User-ID: 2"
```

---

## WebSocket 接口

### 连接地址

```
ws://localhost:8001/ws/room/{room_id}?user_id={user_id}&username={username}
```

### 连接参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| room_id | number | 是 | 房间 ID（路径参数） |
| user_id | number | 是 | 用户 ID（URL 参数） |
| username | string | 是 | 用户名（URL 参数） |

### 消息格式

**客户端 → 服务端**

```json
{
  "type": "message",
  "data": {
    "content": "消息内容"
  }
}
```

**服务端 → 客户端**

```json
{
  "type": "message",
  "data": {
    "id": 123,
    "room_id": 1,
    "user_id": 2,
    "username": "张三",
    "content": "消息内容",
    "type": 1,
    "created_at": 1774408929
  }
}
```

**错误消息**

```json
{
  "type": "error",
  "error": "错误信息"
}
```

### 消息类型

| Type | 说明 |
|------|------|
| `message` | 聊天消息 |
| `system` | 系统消息 |
| `error` | 错误消息 |
| `ping` | 心跳请求 |
| `pong` | 心跳响应 |

---

## 数据模型

### Room（房间）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 房间 ID |
| name | string | 房间名称 |
| owner_id | int64 | 房主 ID |
| type | RoomType | 房间类型 |
| max_count | int | 最大成员数 |
| current_count | int | 当前成员数（非持久化） |
| avatar | string | 头像 URL |
| description | string | 房间描述 |
| tags | string[] | 标签列表 |
| is_public | bool | 是否公开 |
| status | RoomStatus | 房间状态 |
| created_at | int64 | 创建时间（Unix 时间戳） |
| updated_at | int64 | 更新时间（Unix 时间戳） |

**RoomType 枚举**
- `ROOM_TYPE_GROUP` (1) - 普通群聊
- `ROOM_TYPE_VOICE` (2) - 语音房
- `ROOM_TYPE_VIDEO` (3) - 视频房
- `ROOM_TYPE_LIVE` (4) - 直播间

**RoomStatus 枚举**
- `ROOM_STATUS_NORMAL` (1) - 正常
- `ROOM_STATUS_MUTED` (2) - 禁言
- `ROOM_STATUS_BANNED` (3) - 封禁

### RoomMember（房间成员）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 成员记录 ID |
| room_id | int64 | 房间 ID |
| user_id | int64 | 用户 ID |
| role | MemberRole | 成员角色 |
| status | MemberStatus | 成员状态 |
| mute_until | int64 | 禁言截止时间（0 表示未禁言） |
| joined_at | int64 | 加入时间 |
| updated_at | int64 | 更新时间 |

**MemberRole 枚举**
- `ROLE_MEMBER` (1) - 普通成员
- `ROLE_ADMIN` (2) - 管理员
- `ROLE_OWNER` (3) - 房主

**MemberStatus 枚举**
- `MEMBER_STATUS_NORMAL` (1) - 正常
- `MEMBER_STATUS_MUTED` (2) - 禁言中
- `MEMBER_STATUS_LEFT` (3) - 已离开

### Message（消息）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 消息 ID |
| room_id | int64 | 房间 ID |
| user_id | int64 | 发送者 ID |
| content | string | 消息内容 |
| type | MessageType | 消息类型 |
| created_at | int64 | 创建时间 |

**MessageType 枚举**
- `MESSAGE_TYPE_TEXT` (1) - 文本消息
- `MESSAGE_TYPE_IMAGE` (2) - 图片消息（预留）
- `MESSAGE_TYPE_VOICE` (3) - 语音消息（预留）
- `MESSAGE_TYPE_VIDEO` (4) - 视频消息（预留）

---

## 开发指南

### 添加新功能

1. **定义 API**：在 `api/room/v1/room.proto` 中添加服务定义
2. **生成代码**：运行 `make api` 生成 Go 代码
3. **实现仓储**：在 `internal/data/` 中实现数据访问
4. **实现用例**：在 `internal/biz/` 中实现业务逻辑
5. **实现服务**：在 `internal/service/` 中实现服务层
6. **注册服务**：更新 `ProviderSet` 和 `wire.go`
7. **重新生成 Wire**：运行 `make generate`

### 数据库迁移

```bash
# 重置数据库（删除所有数据）
docker-compose down -v
docker-compose up -d

# Auto-migration 会在服务启动时自动执行
```

### 代码生成

```bash
# 生成 Proto 代码
make api

# 生成 Wire 依赖注入
make generate

# 生成所有
make all
```

### 测试

```bash
# 运行所有测试
go test ./...

# 运行测试并显示覆盖率
go test -cover ./...

# 运行特定包的测试
go test ./internal/biz/...
```

---

## 部署指南

### 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| DB_HOST | 数据库主机 | room-db |
| DB_PORT | 数据库端口 | 5432 |
| DB_USER | 数据库用户 | minichat |
| DB_PASSWORD | 数据库密码 | minichat123 |
| DB_NAME | 数据库名称 | room_service |
| HTTP_PORT | HTTP 服务端口 | 8000 |
| GRPC_PORT | gRPC 服务端口 | 9000 |

### 配置文件

`configs/config.yaml`：

```yaml
server:
  http:
    addr: 0.0.0.0:8000
    timeout: 1s
  grpc:
    addr: 0.0.0.0:9000
    timeout: 1s

data:
  database:
    driver: postgres
    host: room-db
    port: 5432
    user: minichat
    password: minichat123
    dbname: room_service
    max_open_conns: 100
    max_idle_conns: 10
    conn_max_lifetime: 300s
    enable_auto_migrate: true

# 用户服务配置
user_service:
  base_url: http://user-service:8000
  timeout: 10s
```

### Docker 部署

```bash
# 构建镜像
docker build -t minichat-room-service .

# 运行容器
docker run -d \
  --name room-service \
  -p 8001:8000 \
  -p 9001:9000 \
  -e DB_HOST=room-db \
  -e DB_PASSWORD=your_password \
  minichat-room-service
```

---

## 故障排查

### 常见问题

#### 1. 数据库连接失败

```
failed to connect to `host=room-db user=minichat database=room_service`:
connection refused
```

**解决方法**：
- 检查 PostgreSQL 容器是否运行：`docker ps | grep room_db`
- 检查网络配置：`docker network ls`
- 验证数据库密码是否正确

#### 2. WebSocket 连接断开

**可能原因**：
- 用户认证失败
- 房间成员状态异常（不是 Normal）
- Context 被取消

**解决方法**：
- 检查 `X-User-ID` 请求头是否正确
- 验证用户是否在房间中且状态为 Normal
- 查看服务日志确认具体错误

#### 3. "User already in room" 错误

**原因**：用户已存在 `Left` 状态的记录

**解决方法**：
- 调用 `/v1/rooms/{id}/join` 接口会自动处理重新加入
- 系统会将状态从 `Left` 更新为 `Normal`

#### 4. 消息显示 "Unknown" 或 "Invalid Date"

**原因**：
- 后端返回驼峰式字段名（`roomId`），前端期望蛇形式（`room_id`）

**解决方法**：
- 前端已实现自动转换，确保使用最新版本代码

### 日志查看

```bash
# Docker 容器日志
docker logs room-service -f

# 查看最近 100 行
docker logs room-service --tail 100

# 查看特定时间的日志
docker logs room-service --since 1h
```

### 健康检查

```bash
# 检查服务状态
curl http://localhost:8001/v1/rooms

# 检查数据库连接
docker exec room_db psql -U minichat -d room_service -c "SELECT 1"
```

---

## 技术栈

- **框架**：[Kratos v2](https://github.com/go-kratos/kratos) - Go 微服务框架
- **ORM**：[GORM v1.25.12](https://github.com/go-gorm/gorm) - Go ORM 库
- **数据库**：PostgreSQL 16
- **WebSocket**：[gorilla/websocket](https://github.com/gorilla/websocket)
- **API 协议**：gRPC + HTTP (Protobuf)
- **依赖注入**：[Wire](https://github.com/google/wire)
- **配置管理**：Kratos Config

---

## 许可证

MIT License

---

## 联系方式

- **团队**：Wuhan University of Technology Students
- **成员**：NameIsNotName (Guo Zifan), Wang Haokun, Chen Xinqi, Jonas Peng
- **项目**：[MiniChatv2](https://github.com/your-org/MiniChatv2)
