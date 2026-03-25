# Room Service API 文档

## 概述

Room Service 提供完整的房间管理和实时消息功能，支持 HTTP REST API 和 WebSocket 双向通信。

**基础 URL**: `http://localhost:8001`

**认证方式**:
- HTTP Header: `X-User-ID: {user_id}`
- (可选) `Authorization: Bearer {jwt_token}`

---

## 房间管理 API

### 1. 创建房间

创建新的聊天房间。

**请求**
```http
POST /v1/rooms
Content-Type: application/json
X-User-ID: 1

{
  "name": "技术交流群",
  "type": 1,
  "max_count": 500,
  "avatar": "https://example.com/avatar.png",
  "description": "讨论技术问题的群组",
  "tags": ["技术", "交流"],
  "is_public": true
}
```

**参数说明**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 房间名称（最大长度 100） |
| type | RoomType | 否 | 房间类型，默认 1 |
| max_count | int | 否 | 最大成员数，默认 500 |
| avatar | string | 否 | 头像 URL |
| description | string | 否 | 房间描述 |
| tags | string[] | 否 | 标签列表 |
| is_public | bool | 否 | 是否公开，默认 true |

**响应**
```json
{
  "room": {
    "id": 1,
    "name": "技术交流群",
    "ownerId": "1",
    "type": "ROOM_TYPE_GROUP",
    "maxCount": 500,
    "currentCount": 1,
    "avatar": "https://example.com/avatar.png",
    "description": "讨论技术问题的群组",
    "tags": ["技术", "交流"],
    "isPublic": true,
    "status": "ROOM_STATUS_NORMAL",
    "createdAt": "1774408929",
    "updatedAt": "1774408929"
  }
}
```

### 2. 获取房间信息

获取指定房间的详细信息。

**请求**
```http
GET /v1/rooms/{room_id}
X-User-ID: 1
```

**响应**
```json
{
  "room": {
    "id": 1,
    "name": "技术交流群",
    "ownerId": "1",
    "currentCount": 5,
    ...
  }
}
```

### 3. 更新房间信息

更新房间的基本信息（仅房主和管理员可操作）。

**请求**
```http
PUT /v1/rooms/{room_id}
Content-Type: application/json
X-User-ID: 1

{
  "name": "新房间名称",
  "description": "更新后的描述",
  "avatar": "https://example.com/new-avatar.png"
}
```

**响应**
```json
{
  "room": {
    "id": 1,
    "name": "新房间名称",
    ...
  }
}
```

### 4. 删除房间

删除指定房间（仅房主可操作）。

**请求**
```http
DELETE /v1/rooms/{room_id}
X-User-ID: 1
```

**响应**
```json
{
  "success": true
}
```

### 5. 获取所有房间列表

获取所有公开房间列表，支持分页。

**请求**
```http
GET /v1/rooms?page=1&page_size=20
X-User-ID: 1
```

**参数说明**

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量（最大 100） |

**响应**
```json
{
  "rooms": [
    {
      "id": 1,
      "name": "技术交流群",
      "currentCount": 5,
      ...
    }
  ],
  "total": 100
}
```

---

## 成员管理 API

### 1. 加入房间

加入指定房间，成为房间成员。

**请求**
```http
POST /v1/rooms/{room_id}/join
Content-Type: application/json
X-User-ID: 2

{}
```

**响应**
```json
{
  "member": {
    "id": 10,
    "roomId": "1",
    "userId": "2",
    "role": "ROLE_MEMBER",
    "status": "MEMBER_STATUS_NORMAL",
    "muteUntil": "0",
    "joinedAt": "1774408929",
    "updatedAt": "1774408929"
  },
  "room": {
    "id": 1,
    "name": "技术交流群",
    ...
  }
}
```

**状态说明**
- 如果用户之前离开过房间（status=Left），此接口会将状态恢复为 Normal
- 如果用户从未加入过，会创建新的成员记录

### 2. 离开房间

离开房间但不删除成员记录（状态变更为 Left）。

**请求**
```http
POST /v1/rooms/{room_id}/leave
Content-Type: application/json
X-User-ID: 2

{}
```

**响应**
```json
{
  "success": true
}
```

**注意**：此操作不会删除成员记录，只是将状态变更为 Left。用户可以重新加入房间。

### 3. 退出群

完全退出群组，删除成员记录。

**请求**
```http
POST /v1/rooms/{room_id}/quit
Content-Type: application/json
X-User-ID: 2

{}
```

**响应**
```json
{
  "success": true
}
```

**注意**：此操作会删除成员记录，退出后需要重新加入才能再次进入房间。

### 4. 获取成员列表

获取房间成员列表，支持分页。

**请求**
```http
GET /v1/rooms/{room_id}/members?page=1&page_size=50
X-User-ID: 1
```

**响应**
```json
{
  "members": [
    {
      "id": 1,
      "roomId": "1",
      "userId": "1",
      "role": "ROLE_OWNER",
      "status": "MEMBER_STATUS_NORMAL",
      "joinedAt": "1774408929"
    },
    {
      "id": 10,
      "roomId": "1",
      "userId": "2",
      "role": "ROLE_MEMBER",
      "status": "MEMBER_STATUS_NORMAL",
      "joinedAt": "1774408930"
    }
  ],
  "total": 50
}
```

**排序规则**：按角色降序（房主 > 管理员 > 成员），加入时间升序排列。

### 5. 踢出成员

将指定成员踢出房间（仅房主和管理员可操作）。

**请求**
```http
DELETE /v1/rooms/{room_id}/members/{user_id}
X-User-ID: 1
```

**响应**
```json
{
  "success": true
}
```

**权限规则**：
- 房主可以踢出任何成员
- 管理员可以踢出普通成员，不能踢出房主和其他管理员

### 6. 更新成员角色

更新成员的角色（仅房主可操作）。

**请求**
```http
PUT /v1/rooms/{room_id}/members/{user_id}/role
Content-Type: application/json
X-User-ID: 1

{
  "role": 2
}
```

**参数说明**

| Role | 说明 |
|------|------|
| 1 | 普通成员 |
| 2 | 管理员 |
| 3 | 房主（不能通过此接口设置） |

**响应**
```json
{
  "success": true
}
```

### 7. 禁言/解禁成员

禁言或解禁指定成员（仅房主和管理员可操作）。

**请求**
```http
PUT /v1/rooms/{room_id}/members/{user_id}/mute
Content-Type: application/json
X-User-ID: 1

{
  "mute_until": 1774495329
}
```

**参数说明**

| 参数 | 类型 | 说明 |
|------|------|------|
| mute_until | int64 | 禁言截止时间（Unix 时间戳），0 表示解禁 |

**响应**
```json
{
  "success": true
}
```

---

## 消息 API

### 1. 发送消息

向房间发送消息。

**请求**
```http
POST /v1/rooms/{room_id}/messages
Content-Type: application/json
X-User-ID: 2

{
  "content": "Hello, World!",
  "type": 1
}
```

**参数说明**

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| content | string | 是 | - | 消息内容（最大 5000 字符） |
| type | MessageType | 否 | 1 | 消息类型 |

**响应**
```json
{
  "message": {
    "id": 100,
    "roomId": "1",
    "userId": "2",
    "type": "MESSAGE_TYPE_TEXT",
    "content": "Hello, World!",
    "createdAt": "1774408929"
  }
}
```

### 2. 获取消息历史

获取房间消息历史，支持分页。

**请求**
```http
GET /v1/rooms/{room_id}/messages?page=1&page_size=50
X-User-ID: 1
```

**响应**
```json
{
  "messages": [
    {
      "id": 100,
      "roomId": "1",
      "userId": "2",
      "type": "MESSAGE_TYPE_TEXT",
      "content": "Hello, World!",
      "createdAt": "1774408929"
    },
    {
      "id": 99,
      "roomId": "1",
      "userId": "1",
      "type": "MESSAGE_TYPE_TEXT",
      "content": "Hi!",
      "createdAt": "1774408920"
    }
  ],
  "total": 500
}
```

**排序规则**：按创建时间倒序排列（最新消息在前）。

---

## 用户相关 API

### 获取用户房间列表

获取指定用户加入的所有房间。

**请求**
```http
GET /v1/users/{user_id}/rooms?page=1&page_size=20
X-User-ID: 1
```

**响应**
```json
{
  "rooms": [
    {
      "id": 1,
      "name": "技术交流群",
      "currentCount": 5,
      ...
    }
  ],
  "total": 10
}
```

**说明**：只返回用户是 Normal 状态的房间。

---

## WebSocket API

### 连接格式

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

#### 发送消息

**客户端 → 服务端**

```json
{
  "type": "message",
  "data": {
    "content": "Hello, everyone!"
  }
}
```

#### 接收消息

**服务端 → 客户端**

聊天消息：
```json
{
  "type": "message",
  "data": {
    "id": 100,
    "room_id": 1,
    "user_id": 2,
    "username": "张三",
    "content": "Hello, everyone!",
    "type": 1,
    "created_at": 1774408929
  }
}
```

系统消息：
```json
{
  "type": "system",
  "data": {
    "message": "用户 李四 加入了房间"
  }
}
```

错误消息：
```json
{
  "type": "error",
  "error": "user is not a member of this room"
}
```

#### 心跳机制

**客户端发送心跳**：
```json
{
  "type": "ping"
}
```

**服务端响应**：
```json
{
  "type": "pong"
}
```

### 连接要求

1. **成员验证**：连接时用户必须是房间的 Normal 状态成员
2. **心跳保活**：客户端需定期发送 ping 消息（建议间隔 30 秒）
3. **重连机制**：连接断开后应自动重连（建议最多重试 5 次）

---

## 错误码

| HTTP Code | Error Message | 说明 |
|------------|---------------|------|
| 400 | Invalid room ID | 房间 ID 无效 |
| 403 | Not a member of this room | 不是房间成员 |
| 403 | Member status not normal | 成员状态异常 |
| 403 | Permission denied | 权限不足 |
| 404 | Room not found | 房间不存在 |
| 404 | Member not found | 成员不存在 |
| 500 | user already in room | 用户已在房间中（状态为 Normal） |
| 500 | Internal server error | 服务器内部错误 |

---

## 数据模型

### RoomType（房间类型）

| 值 | 常量 | 说明 |
|----|------|------|
| 1 | ROOM_TYPE_GROUP | 普通群聊 |
| 2 | ROOM_TYPE_VOICE | 语音房 |
| 3 | ROOM_TYPE_VIDEO | 视频房 |
| 4 | ROOM_TYPE_LIVE | 直播间 |

### RoomStatus（房间状态）

| 值 | 常量 | 说明 |
|----|------|------|
| 1 | ROOM_STATUS_NORMAL | 正常 |
| 2 | ROOM_STATUS_MUTED | 禁言 |
| 3 | ROOM_STATUS_BANNED | 封禁 |

### MemberRole（成员角色）

| 值 | 常量 | 说明 |
|----|------|------|
| 1 | ROLE_MEMBER | 普通成员 |
| 2 | ROLE_ADMIN | 管理员 |
| 3 | ROLE_OWNER | 房主 |

### MemberStatus（成员状态）

| 值 | 常量 | 说明 |
|----|------|------|
| 1 | MEMBER_STATUS_NORMAL | 正常 |
| 2 | MEMBER_STATUS_MUTED | 禁言中 |
| 3 | MEMBER_STATUS_LEFT | 已离开 |

### MessageType（消息类型）

| 值 | 常量 | 说明 |
|----|------|------|
| 1 | MESSAGE_TYPE_TEXT | 文本消息 |
| 2 | MESSAGE_TYPE_IMAGE | 图片消息（预留） |
| 3 | MESSAGE_TYPE_VOICE | 语音消息（预留） |
| 4 | MESSAGE_TYPE_VIDEO | 视频消息（预留） |

---

## 前端集成说明

### 字段命名格式

**后端返回格式**：驼峰式（camelCase）
```json
{
  "roomId": 1,
  "userId": 2,
  "createdAt": 1774408929
}
```

**前端期望格式**：蛇形式（snake_case）
```json
{
  "room_id": 1,
  "user_id": 2,
  "created_at": 1774408929
}
```

**解决方案**：前端需要实现字段格式转换，支持两种格式。

### 用户名显示

后端 WebSocket 消息包含 `username` 字段，但历史消息 API 不返回用户名。

**推荐方案**：
1. WebSocket 实时消息：使用后端返回的 `username`
2. 历史消息：前端根据 `user_id` 缓存用户信息
3. 优先显示用户的 `nickname`，其次 `username`

---

## 调试工具

### Swagger UI

访问 http://localhost:8001/q/ 查看自动生成的 API 文档。

### cURL 示例

```bash
# 完整的加入房间流程
# 1. 创建房间
curl -X POST http://localhost:8001/v1/rooms \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{"name": "测试房间", "max_count": 100}'

# 2. 用户2加入房间
curl -X POST http://localhost:8001/v1/rooms/1/join \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 2"

# 3. 发送消息
curl -X POST http://localhost:8001/v1/rooms/1/messages \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 2" \
  -d '{"content": "Hello!"}'

# 4. 获取消息历史
curl http://localhost:8001/v1/rooms/1/messages \
  -H "X-User-ID: 2"

# 5. 离开房间
curl -X POST http://localhost:8001/v1/rooms/1/leave \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 2"
```

### WebSocket 测试

使用 `wscat` 测试 WebSocket 连接：

```bash
# 安装 wscat
npm install -g wscat

# 连接到房间
wscat -c "ws://localhost:8001/ws/room/1?user_id=2&username=TestUser"

# 发送消息
> {"type":"message","data":{"content":"Hello from WebSocket!"}}

# 发送心跳
> {"type":"ping"}
```
