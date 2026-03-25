# Room Service 变更日志

本文档记录 Room Service 的所有重要变更。

---

## [Unreleased]

### 新增
- 完善的 API 文档和部署文档

### 修复
- 修复 WebSocket 连接使用已取消的 context 导致的 "user is not a member of this room" 错误
- 修复前端消息显示问题，添加字段格式转换（camelCase ↔ snake_case）

### 优化
- 用户名显示优先使用 nickname
- 添加用户名缓存机制

---

## [1.0.0] - 2026-03-25

### 新增
- **房间管理功能**
  - 创建、查询、更新、删除房间
  - 支持公开/私有房间
  - 支持多种房间类型（群聊、语音房、视频房、直播间）

- **成员管理功能**
  - 加入/离开/退出群组
  - 成员角色权限控制（房主、管理员、普通成员）
  - 踢出成员、更新角色、禁言/解禁

- **实时消息功能**
  - WebSocket 双向通信
  - 消息广播到房间所有成员
  - 心跳保活机制
  - 自动重连机制

- **消息历史功能**
  - 分页查询历史消息
  - 消息按时间倒序排列

- **用户服务集成**
  - HTTP 方式调用用户服务
  - 验证用户存在性
  - 获取用户详细信息

### 技术栈
- Go 1.24+
- Kratos v2 框架
- GORM v1.25.12
- PostgreSQL 16
- gorilla/websocket
- Wire 依赖注入

### API 端点
- POST `/v1/rooms` - 创建房间
- GET `/v1/rooms/{id}` - 获取房间信息
- PUT `/v1/rooms/{id}` - 更新房间信息
- DELETE `/v1/rooms/{id}` - 删除房间
- GET `/v1/rooms` - 获取所有房间列表
- POST `/v1/rooms/{id}/join` - 加入房间
- POST `/v1/rooms/{id}/leave` - 离开房间（状态变 Left）
- POST `/v1/rooms/{id}/quit` - 退出群（删除记录）
- GET `/v1/rooms/{id}/members` - 获取成员列表
- DELETE `/v1/rooms/{id}/members/{user_id}` - 踢出成员
- PUT `/v1/rooms/{id}/members/{user_id}/role` - 更新成员角色
- PUT `/v1/rooms/{id}/members/{user_id}/mute` - 禁言/解禁
- GET `/v1/users/{user_id}/rooms` - 获取用户房间列表
- POST `/v1/rooms/{id}/messages` - 发送消息
- GET `/v1/rooms/{id}/messages` - 获取消息历史

### WebSocket 端点
- `ws://localhost:8001/ws/room/{room_id}?user_id={user_id}&username={username}`

### 数据模型
- Room - 房间实体
- RoomMember - 房间成员实体
- Message - 消息实体

### 权限模型
| 操作 | 房主 | 管理员 | 普通成员 |
|------|------|--------|----------|
| 更新房间 | ✓ | ✓ | ✗ |
| 删除房间 | ✓ | ✗ | ✗ |
| 踢出成员 | ✓ | ✓（不能踢管理员/房主）| ✗ |
| 更新角色 | ✓ | ✗ | ✗ |
| 禁言/解禁 | ✓ | ✓ | ✗ |
| 发送消息 | ✓ | ✓ | ✓ |

---

## [0.9.0] - 2026-03-20

### 新增
- 基础房间 CRUD 功能
- 成员加入/离开功能
- 基础 WebSocket 支持

---

## 贡献指南

### 变更类型

- `Added` - 新增功能
- `Changed` - 功能变更
- `Deprecated` - 即将废弃的功能
- `Removed` - 已删除的功能
- `Fixed` - 问题修复
- `Security` - 安全相关修复

### 变更日志格式

```markdown
## [版本号] - 日期

### 新增
- 功能描述

### 修复
- 问题描述

### 变更
- 功能变更说明

### 废弃
- 即将废弃的功能

### 移除
- 已删除的功能

### 安全
- 安全相关修复
```

### 版本号规范

遵循 [语义化版本 2.0.0](https://semver.org/lang/zh-CN)

- MAJOR：不兼容的 API 修改
- MINOR：向下兼容的功能性新增
- PATCH：向下兼容的问题修正
