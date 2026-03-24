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
| 数据库 | PostgreSQL (Docker) |
| 协议 | gRPC + HTTP |
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
- [ ] 数据库表创建
- [ ] Data 层 Repository 实现
- [ ] 数据库连接池配置

### Phase 3: 业务层
- [ ] Biz 层用例实现
- [ ] 权限校验逻辑
- [ ] 事件发布机制

### Phase 4: 服务层
- [ ] Proto 定义生成
- [ ] Service 层实现
- [ ] HTTP/gRPC 服务注册

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

# 事件发布配置（待定）
event:
  type: kafka  # or rabbitmq, redis
  # ...
```
