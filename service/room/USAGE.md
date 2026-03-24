# 房间服务使用文档

## 目录

- [快速开始](#快速开始)
- [启动服务](#启动服务)
- [API 文档 (Swagger UI)](#api-文档-swagger-ui)
- [API 使用示例](#api-使用示例)
- [数据库管理](#数据库管理)
- [开发指南](#开发指南)

---

## 快速开始

### 前置要求

- Go 1.22+
- Docker & Docker Compose
- Make 工具

### 1. 安装依赖

```bash
# 安装 Go 工具链
make init

# 安装项目依赖
go mod tidy
```

### 2. 启动数据库

```bash
# 启动 PostgreSQL (Docker)
docker-compose up -d

# 查看数据库状态
docker-compose ps

# 停止数据库
docker-compose down
```

### 3. 生成代码

```bash
# 生成 API 代码
make api

# 生成配置代码
make config

# 生成 Wire 依赖注入
make generate

# 或一次性生成所有
make all
```

### 4. 构建服务

```bash
make build
```

### 5. 启动服务

```bash
# 方式一：使用构建的二进制文件
./bin/room -conf ./configs

# 方式二：直接运行
go run ./cmd/room -conf ./configs
```

服务启动后：
- **HTTP 服务**: http://localhost:8000
- **gRPC 服务**: localhost:9000

---

## 启动服务

### 完整启动流程

```bash
# 1. 启动数据库
docker-compose up -d

# 2. 生成并构建
make all
make build

# 3. 启动服务
./bin/room -conf ./configs
```

### 验证服务状态

```bash
# 检查 HTTP 服务
curl http://localhost:8000

# 检查服务列表 API
curl http://localhost:8000/q/services
```

---

## API 文档 (Swagger UI)

### 方式一：在线 Swagger Editor（推荐）

1. **获取 OpenAPI 文档**：

项目根目录已生成 `openapi.yaml` 文件，包含完整的 API 定义。

2. **访问在线 Swagger Editor**：

   打开浏览器访问：https://editor.swagger.io/

3. **导入文档**：

   - 方式 A：复制 `openapi.yaml` 的内容并粘贴到编辑器
   - 方式 B：直接拖拽 `openapi.yaml` 文件到编辑器页面

4. **浏览和测试 API**：

   - 左侧查看所有 API 端点
   - 点击任意端点查看详细信息
   - 点击 "Try it out" 按钮在线测试

### 方式二：本地 Swagger UI

使用 Docker 运行 Swagger UI：

```bash
# 在项目根目录运行
docker run -p 8080:8080 \
  -e SWAGGER_JSON=/openapi.yaml \
  -v $(pwd)/openapi.yaml:/openapi.yaml \
  swaggerapi/swagger-ui
```

然后访问：http://localhost:8080

### 方式三：使用 Postman/Insomnia

1. 下载项目中的 `openapi.yaml` 文件
2. 打开 Postman/Insomnia
3. 选择 Import → 选择 `openapi.yaml`

### 方式四：服务发现 API

服务提供了服务发现端点（仅 JSON 数据，无 UI）：

```bash
# 获取所有服务列表
curl http://localhost:8000/q/services

# 返回示例：
# {
#   "services": ["room.v1.RoomService", ...],
#   "methods": ["/room.v1.RoomService/CreateRoom", ...]
# }
```

---

## API 使用示例

### 1. 创建房间

```bash
curl -X POST http://localhost:8000/v1/rooms \
  -H "Content-Type: application/json" \
  -d '{
    "name": "聊天室",
    "type": 1,
    "max_count": 100,
    "avatar": "https://example.com/avatar.png",
    "description": "这是一个测试房间",
    "tags": ["游戏", "聊天"],
    "is_public": true
  }'
```

**响应示例:**
```json
{
  "room": {
    "id": "1",
    "name": "聊天室",
    "ownerId": "1",
    "type": "ROOM_TYPE_GROUP",
    "maxCount": 100,
    "currentCount": 1,
    "avatar": "https://example.com/avatar.png",
    "description": "这是一个测试房间",
    "tags": ["游戏", "聊天"],
    "isPublic": true,
    "status": "ROOM_STATUS_NORMAL",
    "createdAt": "1711234567",
    "updatedAt": "1711234567"
  }
}
```

### 2. 获取房间信息

```bash
curl http://localhost:8000/v1/rooms/1
```

### 3. 更新房间信息

```bash
curl -X PUT http://localhost:8000/v1/rooms/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "更新后的房间名",
    "description": "更新后的描述",
    "is_public": false
  }'
```

### 4. 删除房间

```bash
curl -X DELETE http://localhost:8000/v1/rooms/1
```

### 5. 加入房间

```bash
curl -X POST http://localhost:8000/v1/rooms/1/join \
  -H "Content-Type: application/json"
```

### 6. 退出房间

```bash
curl -X POST http://localhost:8000/v1/rooms/1/leave \
  -H "Content-Type: application/json"
```

### 7. 获取房间成员列表

```bash
curl "http://localhost:8000/v1/rooms/1/members?page=1&page_size=50"
```

### 8. 踢出成员

```bash
curl -X DELETE http://localhost:8000/v1/rooms/1/members/2
```

### 9. 更新成员角色

```bash
curl -X PUT http://localhost:8000/v1/rooms/1/members/2/role \
  -H "Content-Type: application/json" \
  -d '{
    "role": 2
  }'
```

**角色值:**
- `1` = 普通成员 (ROLE_MEMBER)
- `2` = 管理员 (ROLE_ADMIN)
- `3` = 房主 (ROLE_OWNER)

### 10. 禁言/解禁成员

```bash
# 禁言 1 小时
curl -X PUT http://localhost:8000/v1/rooms/1/members/2/mute \
  -H "Content-Type: application/json" \
  -d '{
    "mute_until": '$(date -v+1H +%s)'
  }'

# 解禁
curl -X PUT http://localhost:8000/v1/rooms/1/members/2/mute \
  -H "Content-Type: application/json" \
  -d '{
    "mute_until": 0
  }'
```

### 11. 获取用户房间列表

```bash
curl "http://localhost:8000/v1/users/1/rooms?page=1&page_size=20"
```

---

## 数据模型

### 房间类型 (RoomType)

| 值 | 名称 | 说明 |
|----|------|------|
| 1 | ROOM_TYPE_GROUP | 普通群聊 |
| 2 | ROOM_TYPE_VOICE | 语音房 |
| 3 | ROOM_TYPE_VIDEO | 视频房 |
| 4 | ROOM_TYPE_LIVE | 直播间 |

### 房间状态 (RoomStatus)

| 值 | 名称 | 说明 |
|----|------|------|
| 1 | ROOM_STATUS_NORMAL | 正常 |
| 2 | ROOM_STATUS_MUTED | 禁言 |
| 3 | ROOM_STATUS_BANNED | 封禁 |

### 成员角色 (MemberRole)

| 值 | 名称 | 说明 |
|----|------|------|
| 1 | ROLE_MEMBER | 普通成员 |
| 2 | ROLE_ADMIN | 管理员 |
| 3 | ROLE_OWNER | 房主 |

### 成员状态 (MemberStatus)

| 值 | 名称 | 说明 |
|----|------|------|
| 1 | MEMBER_STATUS_NORMAL | 正常 |
| 2 | MEMBER_STATUS_MUTED | 禁言中 |
| 3 | MEMBER_STATUS_LEFT | 已离开 |

---

## 数据库管理

### 连接数据库

```bash
# 使用 psql 连接
docker exec -it minichat-room-db psql -U minichat -d room_service
```

### 常用 SQL 命令

```sql
-- 查看所有表
\dt

-- 查看房间表结构
\d rooms

-- 查看成员表结构
\d room_members

-- 查询所有房间
SELECT * FROM rooms;

-- 查询房间成员
SELECT * FROM room_members WHERE room_id = 1;

-- 统计房间成员数
SELECT room_id, COUNT(*) as member_count
FROM room_members
WHERE status = 1
GROUP BY room_id;
```

### 重置数据库

```bash
# 停止并删除数据卷
docker-compose down -v

# 重新启动
docker-compose up -d
```

---

## 开发指南

### 目录结构

```
service/room/
├── api/              # API 定义 (Proto)
│   └── room/v1/      # 房间服务 API
├── cmd/              # 命令行工具
│   └── room/         # 主程序入口
├── configs/          # 配置文件
├── docs/             # 文档
├── internal/         # 内部代码
│   ├── biz/          # 业务逻辑层
│   ├── data/         # 数据访问层
│   ├── server/       # 服务器 (HTTP/gRPC)
│   └── service/      # 服务实现
├── bin/              # 编译输出
├── docker-compose.yml
├── Makefile
└── openapi.yaml      # OpenAPI 文档
```

### 添加新功能

1. **定义 API**: 在 `api/room/v1/*.proto` 中定义接口
2. **生成代码**: 运行 `make api`
3. **实现 Biz**: 在 `internal/biz/` 添加业务逻辑
4. **实现 Data**: 在 `internal/data/` 添加数据访问
5. **实现 Service**: 在 `internal/service/` 添加服务
6. **注册服务**: 在 `internal/server/` 注册新服务
7. **更新 Wire**: 运行 `make generate`

### 配置说明

`configs/config.yaml`:

```yaml
server:
  http:
    addr: 0.0.0.0:8000      # HTTP 地址
    timeout: 10s            # 超时时间
    metadata: true          # 启用 OpenAPI
  grpc:
    addr: 0.0.0.0:9000      # gRPC 地址
    timeout: 10s

data:
  database:
    driver: postgres
    source: postgres://minichat:minichat123@localhost:5432/room_service?sslmode=disable
    max_open_conns: 100
    max_idle_conns: 10
    conn_max_lifetime: 300s
    enable_auto_migrate: true  # 开发环境启用自动迁移
```

---

## 常见问题

### 1. 服务启动失败

**问题**: `failed to connect to database`

**解决**:
```bash
# 检查数据库是否运行
docker-compose ps

# 启动数据库
docker-compose up -d
```

### 2. 端口被占用

**问题**: `bind: address already in use`

**解决**:
```bash
# 查找占用端口的进程
lsof -i :8000
lsof -i :9000

# 杀死进程
kill -9 <PID>

# 或修改 configs/config.yaml 中的端口
```

### 3. API 返回 404

**问题**: 接口路径不正确

**解决**: 检查请求路径，确保以 `/v1/` 开头

### 4. 无法访问 Swagger UI

**问题**: Swagger UI 未配置

**解决**:
- 使用在线 Swagger Editor: https://editor.swagger.io/
- 导入项目根目录的 `openapi.yaml` 文件

---

## 许可证

MIT License
