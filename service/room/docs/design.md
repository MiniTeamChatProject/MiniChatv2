# 房间服务设计文档

## 1. 服务概述

房间服务（Room Service）是 MiniChat 系统的核心服务之一，负责管理聊天房间的生命周期、成员关系和状态同步。

### 1.1 职责边界

| 负责 | 不负责 |
|------|--------|
| 房间 CRUD 操作 | 聊天消息存储（由消息服务负责） |
| 成员管理（加入、退出、权限） | 消息推送（连接由网关/连接服务管理） |
| 房间状态同步 | 用户认证（由认证服务负责） |
| 房间元数据持久化 | - |

### 1.2 技术栈

| 组件 | 选型 |
|------|------|
| 语言 | Go 1.22+ |
| 框架 | Kratos v2 |
| ORM | GORM v2 |
| 数据库 | PostgreSQL (Docker) |
| 协议 | gRPC + HTTP |
| API 文档 | Swagger/OpenAPI 3.0 |
| 状态同步 | WebSocket (通过网关) |

---

## 2. 数据模型设计

### 2.1 ER 图

```
┌─────────────┐       ┌──────────────┐       ┌─────────────┐
│   rooms     │──────<│ room_members │>──────│   users     │
│             │       │              │       │(外部服务)   │
│ - id        │       │ - room_id    │       │ - id        │
│ - name      │       │ - user_id    │       │ - name      │
│ - owner_id  │       │ - role       │       │ - avatar    │
│ - type      │       │ - status     │       │             │
│ - max_count │       │ - joined_at  │       │             │
│ - created_at│       │ - updated_at │       │             │
└─────────────┘       └──────────────┘       └─────────────┘
```

### 2.2 表结构

#### rooms - 房间表

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | bigint | 房间 ID | PRIMARY KEY |
| name | varchar(100) | 房间名称 | NOT NULL |
| owner_id | bigint | 房主用户 ID | NOT NULL |
| type | smallint | 房间类型 | NOT NULL, DEFAULT 1 |
| max_count | int | 最大成员数 | DEFAULT 500 |
| avatar | varchar(255) | 房间头像 | |
| description | text | 房间描述 | |
| tags | varchar(500) | 标签（逗号分隔） | |
| is_public | boolean | 是否公开 | DEFAULT true |
| status | smallint | 房间状态 | DEFAULT 1 (1:正常 2:禁言 3:封禁) |
| created_at | timestamptz | 创建时间 | DEFAULT NOW() |
| updated_at | timestamptz | 更新时间 | DEFAULT NOW() |

**房间类型枚举**：
- 1 = 普通群聊
- 2 = 语音房
- 3 = 视频房
- 4 = 直播间

#### room_members - 房间成员表

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | bigint | 主键 ID | PRIMARY KEY |
| room_id | bigint | 房间 ID | NOT NULL, FK -> rooms(id) |
| user_id | bigint | 用户 ID | NOT NULL |
| role | smallint | 成员角色 | DEFAULT 1 |
| status | smallint | 成员状态 | DEFAULT 1 |
| mute_until | timestamptz | 禁言到期时间 | |
| joined_at | timestamptz | 加入时间 | DEFAULT NOW() |
| updated_at | timestamptz | 更新时间 | DEFAULT NOW() |

**唯一索引**：`(room_id, user_id)`

**角色枚举**：
- 1 = 普通成员
- 2 = 管理员
- 3 = 房主

**成员状态枚举**：
- 1 = 正常
- 2 = 禁言中
- 3 = 已离开（软删除记录）

### 2.3 GORM 模型定义

```go
// internal/data/model/room.go
package model

import "time"

// RoomType 房间类型
type RoomType int16

const (
    RoomTypeGroup   RoomType = 1 // 普通群聊
    RoomTypeVoice   RoomType = 2 // 语音房
    RoomTypeVideo   RoomType = 3 // 视频房
    RoomTypeLive    RoomType = 4 // 直播间
)

// RoomStatus 房间状态
type RoomStatus int16

const (
    RoomStatusNormal   RoomStatus = 1 // 正常
    RoomStatusMuted    RoomStatus = 2 // 禁言
    RoomStatusBanned   RoomStatus = 3 // 封禁
)

// Room 房间模型
type Room struct {
    ID          int64       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
    Name        string      `gorm:"column:name;type:varchar(100);not null" json:"name"`
    OwnerID     int64       `gorm:"column:owner_id;not null;index" json:"owner_id"`
    Type        RoomType    `gorm:"column:type;type:smallint;not null;default:1" json:"type"`
    MaxCount    int         `gorm:"column:max_count;type:int;default:500" json:"max_count"`
    Avatar      string      `gorm:"column:avatar;type:varchar(255)" json:"avatar,omitempty"`
    Description string      `gorm:"column:description;type:text" json:"description,omitempty"`
    Tags        string      `gorm:"column:tags;type:varchar(500)" json:"tags,omitempty"`
    IsPublic    bool        `gorm:"column:is_public;type:boolean;default:true" json:"is_public"`
    Status      RoomStatus  `gorm:"column:status;type:smallint;default:1" json:"status"`
    CreatedAt   time.Time   `gorm:"column:created_at;type:timestamptz;default:now()" json:"created_at"`
    UpdatedAt   time.Time   `gorm:"column:updated_at;type:timestamptz;default:now()" json:"updated_at"`

    // 关联
    Members     []RoomMember `gorm:"foreignKey:RoomID" json:"members,omitempty"`
}

// TableName 指定表名
func (Room) TableName() string {
    return "rooms"
}

// MemberRole 成员角色
type MemberRole int16

const (
    RoleMember    MemberRole = 1 // 普通成员
    RoleAdmin     MemberRole = 2 // 管理员
    RoleOwner     MemberRole = 3 // 房主
)

// MemberStatus 成员状态
type MemberStatus int16

const (
    MemberStatusNormal   MemberStatus = 1 // 正常
    MemberStatusMuted    MemberStatus = 2 // 禁言中
    MemberStatusLeft     MemberStatus = 3 // 已离开
)

// RoomMember 房间成员模型
type RoomMember struct {
    ID         int64         `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
    RoomID     int64         `gorm:"column:room_id;not null;uniqueIndex:idx_room_user" json:"room_id"`
    UserID     int64         `gorm:"column:user_id;not null;uniqueIndex:idx_room_user;index" json:"user_id"`
    Role       MemberRole    `gorm:"column:role;type:smallint;default:1" json:"role"`
    Status     MemberStatus  `gorm:"column:status;type:smallint;default:1" json:"status"`
    MuteUntil  *time.Time    `gorm:"column:mute_until;type:timestamptz" json:"mute_until,omitempty"`
    JoinedAt   time.Time     `gorm:"column:joined_at;type:timestamptz;default:now()" json:"joined_at"`
    UpdatedAt  time.Time     `gorm:"column:updated_at;type:timestamptz;default:now()" json:"updated_at"`

    // 关联
    Room       *Room         `gorm:"foreignKey:RoomID" json:"room,omitempty"`
}

// TableName 指定表名
func (RoomMember) TableName() string {
    return "room_members"
}
```

### 2.4 数据库初始化

```go
// internal/data/data.go
package data

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "room/internal/data/model"
)

type Data struct {
    db *gorm.DB
}

// NewData 初始化数据层
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
    // 配置 GORM
    db, err := gorm.Open(postgres.Open(c.Database.Source), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NamingStrategy: schema.NamingStrategy{
            SingularTable: false, // 使用复数表名
        },
    })
    if err != nil {
        return nil, nil, err
    }

    // 连接池配置
    sqlDB, _ := db.DB()
    sqlDB.SetMaxOpenConns(int(c.Database.MaxOpenConns))
    sqlDB.SetMaxIdleConns(int(c.Database.MaxIdleConns))
    sqlDB.SetConnMaxLifetime(c.Database.ConnMaxLifetime.AsDuration())

    // 自动迁移（开发环境）
    if c.Database.EnableAutoMigrate {
        db.AutoMigrate(&model.Room{}, &model.RoomMember{})
    }

    cleanup := func() {
        sqlDB.Close()
    }

    return &Data{db: db}, cleanup, nil
}
```

---

## 3. API 设计

### 3.1 服务定义（Proto）

```protobuf
syntax = "proto3";

package room.v1;

option go_package = "room/api/room/v1;v1";

// 房间服务
service RoomService {
  // 创建房间
  rpc CreateRoom (CreateRoomRequest) returns (CreateRoomReply);

  // 获取房间信息
  rpc GetRoom (GetRoomRequest) returns (GetRoomReply);

  // 更新房间信息
  rpc UpdateRoom (UpdateRoomRequest) returns (UpdateRoomReply);

  // 删除房间
  rpc DeleteRoom (DeleteRoomRequest) returns (DeleteRoomReply);

  // 加入房间
  rpc JoinRoom (JoinRoomRequest) returns (JoinRoomReply);

  // 退出房间
  rpc LeaveRoom (LeaveRoomRequest) returns (LeaveRoomReply);

  // 获取房间成员列表
  rpc ListMembers (ListMembersRequest) returns (ListMembersReply);

  // 踢出成员
  rpc KickMember (KickMemberRequest) returns (KickMemberReply);

  // 更新成员角色
  rpc UpdateMemberRole (UpdateMemberRoleRequest) returns (UpdateMemberRoleReply);

  // 禁言/解禁成员
  rpc MuteMember (MuteMemberRequest) returns (MuteMemberReply);

  // 获取用户加入的房间列表
  rpc ListUserRooms (ListUserRoomsRequest) returns (ListUserRoomsReply);
}
```

### 3.2 HTTP 端点（示例）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /v1/rooms | 创建房间 |
| GET | /v1/rooms/{id} | 获取房间信息 |
| PUT | /v1/rooms/{id} | 更新房间信息 |
| DELETE | /v1/rooms/{id} | 删除房间 |
| POST | /v1/rooms/{id}/join | 加入房间 |
| POST | /v1/rooms/{id}/leave | 退出房间 |
| GET | /v1/rooms/{id}/members | 获取成员列表 |
| DELETE | /v1/rooms/{id}/members/{user_id} | 踢出成员 |
| PUT | /v1/rooms/{id}/members/{user_id}/role | 更新成员角色 |
| PUT | /v1/rooms/{id}/members/{user_id}/mute | 禁言/解禁 |
| GET | /v1/users/{user_id}/rooms | 获取用户房间列表 |

### 3.3 Proto 定义与 Swagger 注释

```protobuf
syntax = "proto3";

package room.v1;

import "google/api/annotations.proto";

option go_package = "room/api/room/v1;v1";

// 创建房间请求
message CreateRoomRequest {
  string name = 1;
  int32 type = 2;        // 1:群聊 2:语音房 3:视频房 4:直播间
  int32 max_count = 3;   // 最大成员数
  string avatar = 4;
  string description = 5;
  repeated string tags = 6;
  bool is_public = 7;
}

// 创建房间响应
message CreateRoomReply {
  int64 room_id = 1;
  string name = 2;
  int64 owner_id = 3;
}

// 房间信息
message Room {
  int64 id = 1;
  string name = 2;
  int64 owner_id = 3;
  int32 type = 4;
  int32 max_count = 5;
  int32 current_count = 6;
  string avatar = 7;
  string description = 8;
  repeated string tags = 9;
  bool is_public = 10;
  int32 status = 11;
  int64 created_at = 12;
}

service RoomService {
  // 创建房间
  rpc CreateRoom (CreateRoomRequest) returns (CreateRoomReply) {
    option (google.api.http) = {
      post: "/v1/rooms"
      body: "*"
    };
  }

  // 获取房间信息
  rpc GetRoom (GetRoomRequest) returns (GetRoomReply) {
    option (google.api.http) = {
      get: "/v1/rooms/{id}"
    };
  }
}
```

运行 `make api` 后，Kratos 会自动生成 `openapi.yaml` 文件。

### 3.4 Swagger UI 访问

启动服务后，可通过以下方式访问 API 文档：

- **OpenAPI Spec**: http://localhost:8000/q/openapi.yaml
- **Swagger UI**: 通过 Kratos 的 swagger 中间件访问

---

## 4. 业务逻辑设计

### 4.1 权限模型

| 操作 | 房主 | 管理员 | 普通成员 |
|------|------|--------|----------|
| 更新房间信息 | ✓ | ✓ | ✗ |
| 删除房间 | ✓ | ✗ | ✗ |
| 踢出成员 | ✓ | ✓ (非管理员) | ✗ |
| 更新成员角色 | ✓ | ✗ | ✗ |
| 禁言成员 | ✓ | ✓ | ✗ |
| 解禁成员 | ✓ | ✓ | ✗ |

### 4.2 状态同步事件

房间服务通过消息队列发布状态变更事件，供网关/连接服务消费并推送：

| 事件类型 | Payload | 触发时机 |
|----------|---------|----------|
| `room.created` | room_id, creator_id | 创建房间 |
| `room.updated` | room_id, changes | 更新房间信息 |
| `room.deleted` | room_id | 删除房间 |
| `member.joined` | room_id, user_id | 成员加入 |
| `member.left` | room_id, user_id | 成员退出 |
| `member.kicked` | room_id, user_id, operator_id | 踢出成员 |
| `member.role_changed` | room_id, user_id, new_role | 角色变更 |
| `member.muted` | room_id, user_id, until | 禁言 |
| `member.unmuted` | room_id, user_id | 解禁 |

### 4.3 核心流程

#### 加入房间流程
```
1. 验证房间存在且未满
2. 检查用户是否已在房间中
3. 创建成员记录
4. 发布 member.joined 事件
5. 返回房间信息
```

#### 踢出成员流程
```
1. 验证操作者权限（房主/管理员）
2. 验证目标成员角色（不能踢房主/同级）
3. 删除/软删除成员记录
4. 发布 member.kicked 事件
5. 返回成功
```

---

## 5. 架构设计

### 5.1 分层架构

```
┌─────────────────────────────────────────────────────────┐
│                    API Layer (proto)                    │
└─────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────┐
│              Service Layer (internal/service)           │
│  - 请求参数验证                                          │
│  - 调用 Biz 层                                           │
│  - 响应转换                                              │
└─────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────┐
│                Biz Layer (internal/biz)                 │
│  - 业务逻辑编排                                          │
│  - 权限校验                                              │
│  - 事件发布                                              │
└─────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────┐
│                 Data Layer (internal/data)              │
│  - 数据库操作 (PostgreSQL)                              │
│  - 缓存操作 (Redis, 可选)                               │
│  - 外部服务调用 (用户服务)                               │
└─────────────────────────────────────────────────────────┘
```

### 5.2 依赖关系

```
room-service
    ├── PostgreSQL (房间元数据存储)
    ├── Redis (可选: 在线状态缓存)
    ├── Kafka/RabbitMQ (状态事件发布)
    └── user-service (用户信息查询)
```

---

## 6. 开发计划

### Phase 1: 基础设施
- [ ] PostgreSQL Docker 环境搭建
- [ ] 数据库迁移工具集成
- [ ] 基础项目结构调整

### Phase 2: 数据层
- [ ] GORM 模型定义 (`internal/data/model/`)
- [ ] 数据库初始化与连接池配置
- [ ] Repository 接口定义 (`internal/biz/`)
- [ ] Repository 实现 (`internal/data/`)
  - [ ] RoomRepository
  - [ ] RoomMemberRepository

### Phase 3: 业务层
- [ ] Biz 层用例实现
- [ ] 权限校验逻辑
- [ ] 事件发布机制

### Phase 4: 服务层
- [ ] Proto 定义编写 (`api/room/v1/`)
- [ ] 运行 `make api` 生成代码
- [ ] Service 层实现 (`internal/service/`)
- [ ] HTTP/gRPC 服务注册
- [ ] Swagger UI 集成与测试

### Phase 5: 测试与优化
- [ ] 单元测试
- [ ] 集成测试
- [ ] 性能优化

---

## 7. 配置参考

### 7.1 Docker Compose

```yaml
services:
  postgres:
    image: postgres:16-alpine
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: minichat
      POSTGRES_PASSWORD: minichat123
      POSTGRES_DB: room_service
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

### 7.2 应用配置 (configs/config.yaml)

```yaml
server:
  http:
    addr: 0.0.0.0:8000
    timeout: 10s
    # Swagger UI 配置
    metadata: true
  grpc:
    addr: 0.0.0.0:9000
    timeout: 10s

data:
  database:
    driver: postgres
    source: postgres://minichat:minichat123@localhost:5432/room_service?sslmode=disable
    max_open_conns: 100
    max_idle_conns: 10
    conn_max_lifetime: 300s
    # 开发环境启用自动迁移（生产环境建议关闭）
    enable_auto_migrate: true

# 事件发布配置（待定）
event:
  type: kafka  # or rabbitmq, redis
  # ...
```

### 7.3 HTTP Server 配置 (internal/server/http.go)

```go
package server

import (
    v1 "room/api/room/v1"
    "room/internal/conf"
    "room/internal/service"

    "github.com/go-kratos/kratos/v2/log"
    "github.com/go-kratos/kratos/v2/middleware/recovery"
    "github.com/go-kratos/kratos/v2/middleware/selector"
    "github.com/go-kratos/kratos/v2/transport/http"
    "github.com/go-kratos/swagger-api" // Swagger UI 中间件
)

// NewHTTPServer 创建 HTTP 服务器
func NewHTTPServer(c *conf.Server, roomSvc *service.RoomService, logger log.Logger) *http.Server {
    var opts = []http.ServerOption{
        http.Middleware(
            recovery.Recovery(),
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

    // 注册 gRPC-Gateway
    v1.RegisterRoomServiceHTTPServer(srv, roomSvc)

    // 启用 Swagger UI
    if c.Http.Metadata {
        srv.Route("/").GET(swagger.UIHandler)
    }

    return srv
}
```

### 7.4 依赖安装

```bash
# GORM
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres

# Swagger UI (Kratos 集成)
go get -u github.com/go-kratos/swagger-api
```

---
