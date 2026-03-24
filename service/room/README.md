# MiniChat 房间服务

> 基于 Kratos 框架的房间服务，提供聊天房间的管理功能。

## 快速开始

```bash
# 1. 启动数据库
docker-compose up -d

# 2. 构建服务
make build

# 3. 启动服务
./bin/room -conf ./configs
```

服务启动后：
- HTTP: http://localhost:8000
- gRPC: localhost:9000

## 文档

- 📖 **[使用文档](USAGE.md)** - API 使用说明、Swagger UI 配置
- 📐 **[设计文档](docs/design.md)** - 架构设计、数据模型

## API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/v1/rooms` | 创建房间 |
| GET | `/v1/rooms/{id}` | 获取房间信息 |
| PUT | `/v1/rooms/{id}` | 更新房间信息 |
| DELETE | `/v1/rooms/{id}` | 删除房间 |
| POST | `/v1/rooms/{id}/join` | 加入房间 |
| POST | `/v1/rooms/{id}/leave` | 退出房间 |
| GET | `/v1/rooms/{id}/members` | 获取成员列表 |
| DELETE | `/v1/rooms/{id}/members/{user_id}` | 踢出成员 |
| PUT | `/v1/rooms/{id}/members/{user_id}/role` | 更新成员角色 |
| PUT | `/v1/rooms/{id}/members/{user_id}/mute` | 禁言/解禁 |
| GET | `/v1/users/{user_id}/rooms` | 获取用户房间列表 |

## 开发

```bash
# 安装依赖
make init

# 生成代码
make all

# 运行测试
go test ./...

# 生成 Wire
make generate
```

## 技术栈

- **框架**: [Kratos v2](https://github.com/go-kratos/kratos)
- **ORM**: [GORM v2](https://github.com/go-kratos/kratos)
- **数据库**: PostgreSQL 16
- **API**: gRPC + HTTP
- **协议**: Protobuf v3

## 许可证

MIT License
