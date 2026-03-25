# MiniChatv2

> 🚀 基于 Kratos 框架的微服务实时聊天平台 - 团队协作聊天系统

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Kratos Framework](https://img.shields.io/badge/Kratos-v2.0-blue)](https://github.com/go-kratos/kratos)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-blue)](https://www.postgresql.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

---

## 📋 目录

- [项目概述](#项目概述)
- [系统架构](#系统架构)
- [功能特性](#功能特性)
- [技术栈](#技术栈)
- [快速开始](#快速开始)
- [部署指南](#部署指南)
- [服务端口](#服务端口)
- [开发文档](#开发文档)
- [团队信息](#团队信息)

---

## 项目概述

MiniChatv2 是一个基于微服务架构的实时团队协作聊天平台，采用 **Go + Kratos 框架** 构建，支持房间管理、成员权限控制、实时消息推送等功能。

### 微服务架构

```
┌─────────────────────────────────────────────────────────────┐
│                    MiniChatv2 微服务架构                        │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐        ┌──────────────┐                   │
│  │  用户服务     │        │  房间服务     │                   │
│  │  :8000       │   ◄──► │  :8001       │                   │
│  │              │  gRPC  │              │                   │
│  │  注册/登录    │        │  房间/聊天     │                   │
│  │  JWT Token   │        │  WebSocket   │                   │
│  └──────┬───────┘        └──────┬───────┘                   │
│         │                        │                           │
│         ▼                        ▼                           │
│  ┌──────────────┐        ┌──────────────┐                   │
│  │  用户数据库   │        │  房间数据库   │                   │
│  │  benutzer_db │        │ room_service │                   │
│  │  :15432      │        │   :5432      │                   │
│  └──────────────┘        └──────────────┘                   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 系统架构

### 服务列表

| 服务 | 端口映射 | 数据库 | 说明 |
|------|----------|--------|------|
| User Service | 8000 (HTTP)<br>9000 (gRPC) | benutzer_db:15432 | 用户注册、登录、JWT 认证 |
| Room Service | 8001 (HTTP)<br>9001 (gRPC) | room_service:5432 | 房间管理、实时消息、WebSocket |

### 前端项目

- **web-vue**: Vue 3 测试前端（见 `web-vue/README_TEST_ONLY.md`）

---

## 功能特性

### 用户服务 (User Service)
- ✅ 用户注册与登录
- ✅ JWT Token 认证
- ✅ 用户资料管理
- ✅ gRPC/HTTP 双协议支持

### 房间服务 (Room Service)
- ✅ 房间创建与管理
- ✅ 成员权限控制（房主、管理员、普通成员）
- ✅ 实时消息推送 (WebSocket)
- ✅ 消息历史查询
- ✅ 成员踢出、禁言功能

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

## 技术栈

### 后端技术
- **框架**: [Kratos v2](https://github.com/go-kratos/kratos) - Go 微服务框架
- **ORM**: GORM v1.25.12 (Go 1.24 兼容)
- **数据库**: PostgreSQL 16
- **实时通信**: gorilla/websocket
- **协议**: Protobuf v3 (gRPC + HTTP)
- **依赖注入**: Wire

### 前端技术 (web-vue)
- **框架**: Vue 3 + TypeScript
- **构建工具**: Vite
- **样式**: TailwindCSS
- **HTTP**: Axios

---

## 快速开始

### 环境要求

- **Go**: 1.24+ (Room Service), 1.22+ (User Service)
- **Docker**: 20.10+
- **Docker Compose**: 2.0+
- **PostgreSQL**: 16+
- **Node.js**: 18+ (web-vue 前端，运行在端口 3000)

### 一键启动（推荐）

最简单的启动方式是使用 Docker Compose：

```bash
# 1. 克隆项目
git clone https://github.com/MiniTeamChatProject/MiniChatv2.git
cd MiniChatv2

# 2. 启动所有服务（包括数据库）
docker-compose up -d

# 3. 查看服务状态
docker-compose ps

# 4. 查看日志
docker-compose logs -f

# 5. 停止服务
docker-compose down

# 6. 停止服务并删除数据卷（清空数据库）
docker-compose down -v
```

启动后服务将在以下端口运行：
- **User Service HTTP**: http://localhost:8000
- **User Service gRPC**: localhost:9000
- **Room Service HTTP**: http://localhost:8001
- **Room Service gRPC**: localhost:9001
- **User Database**: localhost:15432
- **Room Database**: localhost:5432

---

## 部署指南

### 开发环境部署

#### 方式一：使用 Docker Compose（推荐）

这是最简单的部署方式，适合快速开发和测试。

```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看特定服务日志
docker-compose logs -f user-service
docker-compose logs -f room-service

# 重启服务
docker-compose restart user-service
docker-compose restart room-service

# 停止所有服务
docker-compose down
```

#### 方式二：手动部署

##### 1. 启动数据库

```bash
# 启动 User Database
docker run -d \
  --name user_db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=07210721 \
  -e POSTGRES_DB=benutzer_db \
  -p 15432:5432 \
  docker.m.daocloud.io/library/postgres:16-alpine

# 启动 Room Database
docker run -d \
  --name room_db \
  -e POSTGRES_USER=minichat \
  -e POSTGRES_PASSWORD=minichat123 \
  -e POSTGRES_DB=room_service \
  -p 5432:5432 \
  docker.m.daocloud.io/library/postgres:16-alpine
```

##### 2. 启动 User Service

```bash
cd service/user

# 安装依赖
go mod tidy

# 配置数据库连接
export DB_HOST=localhost
export DB_PORT=15432
export DB_USER=postgres
export DB_PASSWORD=07210721
export DB_NAME=benutzer_db

# 运行服务
go run ./cmd/registration/main.go -conf ./configs
```

##### 3. 启动 Room Service

```bash
cd service/room

# 安装依赖
go mod tidy

# 配置数据库连接
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=minichat
export DB_PASSWORD=minichat123
export DB_NAME=room_service

# 运行服务
go run ./cmd/room/main.go -conf ./configs
```

### 生产环境部署

#### 方式一：Docker Compose (生产环境)

创建生产环境配置文件 `docker-compose.prod.yml`：

```yaml
version: "3.9"

services:
  user-service:
    build:
      context: ./service/user
      dockerfile: Dockerfile
    container_name: user_service_prod
    ports:
      - "8000:8000"
      - "9000:9000"
    depends_on:
      - user-db
    environment:
      - DB_HOST=user-db
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=benutzer_db
    networks:
      - minichat-network
    restart: always
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

  room-service:
    build:
      context: ./service/room
      dockerfile: Dockerfile
    container_name: room_service_prod
    ports:
      - "8001:8000"
      - "9001:9000"
    depends_on:
      - room-db
    environment:
      - DB_HOST=room-db
      - DB_PORT=5432
      - DB_USER=minichat
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=room_service
    networks:
      - minichat-network
    restart: always
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

  user-db:
    image: postgres:16-alpine
    container_name: user_db_prod
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: benutzer_db
    volumes:
      - user-postgres-data:/var/lib/postgresql/data
    networks:
      - minichat-network
    restart: always

  room-db:
    image: postgres:16-alpine
    container_name: room_db_prod
    environment:
      POSTGRES_USER: minichat
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: room_service
    volumes:
      - room-postgres-data:/var/lib/postgresql/data
    networks:
      - minichat-network
    restart: always

volumes:
  user-postgres-data:
  room-postgres-data:

networks:
  minichat-network:
    driver: bridge
```

**部署命令**：

```bash
# 设置环境变量
export DB_PASSWORD=your_secure_password

# 启动生产环境服务
docker-compose -f docker-compose.prod.yml up -d

# 查看日志
docker-compose -f docker-compose.prod.yml logs -f
```

#### 方式二：Kubernetes 部署

详见各服务的部署文档：
- [User Service 部署文档](./service/user/README.md)
- [Room Service 部署文档](./service/room/docs/deployment.md)

### 环境变量配置

生产环境需要配置的关键环境变量：

```bash
# 数据库密码（请使用强密码）
DB_PASSWORD=your_secure_password_here

# 服务端口（可根据需要调整）
USER_SERVICE_HTTP_PORT=8000
USER_SERVICE_GRPC_PORT=9000
ROOM_SERVICE_HTTP_PORT=8001
ROOM_SERVICE_GRPC_PORT=9001
```

### 健康检查

部署后验证服务状态：

```bash
# 检查 User Service
curl http://localhost:8000/q/

# 检查 Room Service
curl http://localhost:8001/q/

# 检查数据库连接
docker exec user_db psql -U postgres -d benutzer_db -c "SELECT 1"
docker exec room_db psql -U minichat -d room_service -c "SELECT 1"

# 测试用户注册
curl -X POST http://localhost:8000/register \
  -H "Content-Type: application/json" \
  -d '{"username": "test", "password": "pass123", "nickname": "Test User"}'

# 测试登录获取 Token
curl -X POST http://localhost:8000/login \
  -H "Content-Type: application/json" \
  -d '{"username": "test", "password": "pass123"}'

# 测试创建房间
curl -X POST http://localhost:8001/v1/rooms \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{"name": "测试房间", "max_count": 100}'
```

---

## 服务端口

| 服务 | HTTP | gRPC | 数据库 | DB 端口 |
|------|------|------|--------|--------|
| User Service | 8000 | 9000 | benutzer_db | 15432 |
| Room Service | 8001 | 9001 | room_service | 5432 |

### Swagger UI

- **User Service**: http://localhost:8000/q/
- **Room Service**: http://localhost:8001/q/

---

## 开发文档

### 详细文档

- [User Service README](./service/user/README.md) - 用户服务详细文档
- [Room Service README](./service/room/README.md) - 房间服务详细文档
- [Room Service API 文档](./service/room/docs/api.md) - API 接口文档
- [Room Service 部署文档](./service/room/docs/deployment.md) - 部署指南
- [PRD.md](./PRD.md) - 产品需求文档

### API 文档

#### 用户服务 API

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/register` | 用户注册 |
| POST | `/login` | 用户登录（返回 JWT） |
| GET | `/user/profile` | 获取用户资料 |
| PUT | `/user/profile` | 更新用户资料 |
| DELETE | `/user/{id}` | 删除用户 |

#### 房间服务 API

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/v1/rooms` | 创建房间 |
| GET | `/v1/rooms/{id}` | 获取房间信息 |
| POST | `/v1/rooms/{id}/join` | 加入房间 |
| POST | `/v1/rooms/{id}/leave` | 离开房间（状态变 Left）|
| POST | `/v1/rooms/{id}/quit` | 退出群（删除记录）|
| GET | `/v1/rooms/{id}/messages` | 获取消息历史 |
| **WebSocket** | `ws://localhost:8001/ws/room/{id}` | 实时聊天 |

### 开发工作流

```bash
# 1. 启动开发环境
docker-compose up -d

# 2. 查看服务日志
docker-compose logs -f user-service
docker-compose logs -f room-service

# 3. 重启某个服务
docker-compose restart room-service

# 4. 重置数据库（清空所有数据）
docker-compose down -v
docker-compose up -d

# 5. 停止所有服务
docker-compose down
```

### 数据库重置

```bash
# 完全重置（删除所有数据）
docker-compose down -v
docker-compose up -d

# 重置特定数据库
# User Service Database
docker exec user_db psql -U postgres -d benutzer_db -c "DELETE FROM users;"

# Room Service Database
docker exec room_db psql -U minichat -d room_service -c "DELETE FROM rooms; DELETE FROM room_members; DELETE FROM messages;"
```

---

## 前端项目

### web-vue（测试版本）

⚠️ **注意**：这是一个用于测试的最简陋版本前端，详见 [web-vue/README_TEST_ONLY.md](./web-vue/README_TEST_ONLY.md)

**快速启动**：

```bash
cd web-vue
npm install
npm run dev
```

访问 http://localhost:3000 即可测试。

---

## 常见问题

### 1. 端口被占用

```bash
# 检查端口占用
lsof -i :8000
lsof -i :8001
lsof -i :15432
lsof -i :5432

# 杀死占用进程
kill -9 <PID>
```

### 2. Docker 容器启动失败

```bash
# 查看容器日志
docker logs user-service
docker logs room-service
docker logs user-db
docker logs room-db

# 检查容器状态
docker ps -a

# 重启容器
docker-compose restart <service-name>
```

### 3. 数据库连接失败

```bash
# 检查数据库是否运行
docker ps | grep postgres

# 测试数据库连接
docker exec user_db psql -U postgres -c "SELECT 1"
docker exec room_db psql -U minichat -c "SELECT 1"

# 检查网络配置
docker network ls
docker network inspect minichat-network
```

### 4. "user already in room" 错误

这是正常行为，当用户状态为 `Left` 时重新加入会复用记录。详情见 [Room Service 文档](./service/room/CLAUDE.md)。

### 5. 服务间通信失败

```bash
# 测试 User Service 内部接口
curl http://localhost:8000/internal/verify/1

# 检查服务是否在同一网络
docker network inspect minichat-network
```

---

## 技术架构详解

### 分层架构

每个服务都遵循 Kratos 的分层架构：

```
┌─────────────────────────────────┐
│         API Layer (proto)         │
│  Protobuf 定义 → Go 代码生成      │
└──────────────┬──────────────────┘
               │
┌──────────────▼──────────────────┐
│       Service Layer               │
│  HTTP/gRPC Handler               │
│  业务逻辑编排                    │
└──────────────┬──────────────────┘
               │
┌──────────────▼──────────────────┐
│       Business Logic (biz)        │
│  用例 (Usecase)                   │
│  领域模型 (Entity)                 │
│  仓储接口 (Repository)             │
└──────────────┬──────────────────┘
               │
┌──────────────▼──────────────────┐
│         Data Layer                │
│  Repository 实现                   │
│  GORM 数据访问                    │
└──────────────┬──────────────────┘
               │
┌──────────────▼──────────────────┐
│       PostgreSQL                  │
│  数据持久化                       │
└─────────────────────────────────┘
```

### 依赖注入

使用 [Wire](https://github.com/google/wire) 进行编译时依赖注入：

```go
//go:generate wire
//go:build !debug
// +build wire

// Wire 代码生成
go run github.com/google/wire/cmd/wire
```

---

## 团队信息

**项目名称**: MiniChatv2

**团队成员**：
- NameIsNotName (郭子凡) - Team Lead
- Wang Haokun - 后端开发
- Chen Xinqi - 后端开发
- Jonas Peng - 前端开发

**项目类型**: 武汉理工大学课程设计项目

---

## 许可证

MIT License

---

## 联系方式

- **GitHub**: https://github.com/MiniTeamChatProject/MiniChatv2

---

## 更新日志

### [Unreleased]

- 完善房间服务功能和文档
- 修复 Join/Leave 房间逻辑问题
- 修复 WebSocket Context 问题
- 新增测试前端 web-vue
- 完善 API 和部署文档
