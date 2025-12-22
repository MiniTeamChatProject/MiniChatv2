# User Service (Kratos重构版)
## 用户登录注册验证微服务

**开发者：Jonas Peng**


> 本项目是用户鉴权微服务的 **Kratos 重构版本**。实现了基于 DDD (领域驱动设计) 的分层架构，支持 HTTP 和 gRPC 双协议。

## 技术栈 

- **Framework**: [Kratos v2](https://github.com/go-kratos/kratos) (Microservice Framework)
- **Database**: PostgreSQL
- **ORM**: GORM (With AutoMigrate)
- **Auth**: JWT (JSON Web Token) + Bcrypt
- **DI**: Google Wire (Dependency Injection)
- **Protocol**: HTTP / gRPC

---

## 快速启动 (Quick Start)

### 1. 环境准备
- Ubuntu 24.04.3 LTS
- Go 1.22+
- PostgreSQL (数据库名: `benutzer_db`)
- [Kratos CLI](https://go-kratos.dev/docs/getting-started/start#install-cli) (推荐安装)

### 2. 配置文件
请检查 `internal/data/data.go` 中的数据库连接字符串，或查看 `configs/config.yaml` (如有)。
> 默认连接: `postgres://postgres:07210721@localhost:5432/benutzer_db`

### 3. 运行服务
 *Linux脚本*
```bash
# 方式一：使用 Kratos CLI (推荐)
kratos run

# 方式二：直接运行 Go
cd cmd/user
go run .