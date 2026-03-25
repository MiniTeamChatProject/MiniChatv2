# Room Service 部署文档

## 目录

- [环境准备](#环境准备)
- [本地开发部署](#本地开发部署)
- [Docker 部署](#docker-部署)
- [生产环境部署](#生产环境部署)
- [配置管理](#配置管理)
- [监控与日志](#监控与日志)

---

## 环境准备

### 系统要求

- **操作系统**：Linux (推荐 Ubuntu 20.04+ / CentOS 8+)
- **Go 版本**：1.24+
- **Docker 版本**：20.10+
- **Docker Compose 版本**：2.0+
- **PostgreSQL 版本**：16+

### 依赖服务

- PostgreSQL 16（必需）
- User Service（可选，用于用户认证）

---

## 本地开发部署

### 1. 安装依赖工具

```bash
# 安装 Go 1.24+
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# 验证安装
go version

# 安装 Make
sudo apt-get install build-essential

# 安装 Protoc 编译器
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
go install github.com/google/wire/cmd/wire@latest
```

### 2. 克隆代码

```bash
git clone https://github.com/your-org/MiniChatv2.git
cd MiniChatv2/service/room
```

### 3. 配置数据库

```bash
# 启动 PostgreSQL
docker-compose up -d room-db

# 等待数据库启动
sleep 5

# 验证连接
docker exec room_db psql -U minichat -d room_service -c "SELECT version();"
```

### 4. 配置文件

```bash
# 复制配置文件模板
cp configs/config.yaml.example configs/config.yaml

# 编辑配置
vim configs/config.yaml
```

**配置示例**：

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
    host: localhost
    port: 5432
    user: minichat
    password: minichat123
    dbname: room_service
    max_open_conns: 100
    max_idle_conns: 10
    conn_max_lifetime: 300s
    enable_auto_migrate: true

user_service:
  base_url: http://localhost:8000
  timeout: 10s
```

### 5. 构建和运行

```bash
# 安装依赖
go mod tidy

# 生成代码
make all

# 构建
make build

# 运行
./bin/room -conf ./configs
```

### 6. 验证服务

```bash
# 检查 HTTP 服务
curl http://localhost:8000/v1/rooms

# 检查 gRPC 服务
grpcurl -plaintext localhost:9000 list

# 查看日志
tail -f logs/room.log
```

---

## Docker 部署

### 1. 构建 Docker 镜像

```bash
# 使用 Docker Compose 构建
docker-compose build room-service

# 或手动构建
docker build -t minichat-room-service:latest .
```

### 2. 启动服务

```bash
# 启动所有服务
docker-compose up -d

# 查看运行状态
docker-compose ps

# 查看日志
docker-compose logs -f room-service
```

### 3. Docker Compose 配置

```yaml
version: "3.9"

services:
  room-service:
    build:
      context: ./service/room
      dockerfile: Dockerfile
    container_name: room_service
    ports:
      - "8001:8000"   # HTTP
      - "9001:9000"   # gRPC
    depends_on:
      - room-db
    environment:
      - DB_HOST=room-db
      - DB_PORT=5432
      - DB_USER=minichat
      - DB_PASSWORD=minichat123
      - DB_NAME=room_service
    networks:
      - minichat-network
    restart: unless-stopped

  room-db:
    image: postgres:16-alpine
    container_name: room_db
    environment:
      POSTGRES_USER: minichat
      POSTGRES_PASSWORD: minichat123
      POSTGRES_DB: room_service
    ports:
      - "5432:5432"
    volumes:
      - room-postgres-data:/var/lib/postgresql/data
    networks:
      - minichat-network
    restart: unless-stopped

volumes:
  room-postgres-data:

networks:
  minichat-network:
    driver: bridge
```

### 4. 健康检查

```bash
# 检查服务健康状态
curl http://localhost:8001/v1/rooms

# 检查容器状态
docker ps | grep room

# 检查数据库连接
docker exec room_db psql -U minichat -d room_service -c "SELECT 1"
```

---

## 生产环境部署

### 1. 环境变量配置

创建 `.env` 文件：

```bash
# 服务配置
SERVICE_NAME=room-service
SERVICE_VERSION=v1.0.0
ENVIRONMENT=production

# HTTP 服务
HTTP_ADDR=0.0.0.0:8000
HTTP_TIMEOUT=1s

# gRPC 服务
GRPC_ADDR=0.0.0.0:9000
GRPC_TIMEOUT=1s

# 数据库配置
DB_DRIVER=postgres
DB_HOST=your-db-host
DB_PORT=5432
DB_USER=minichat
DB_PASSWORD=your-secure-password
DB_NAME=room_service
DB_MAX_OPEN_CONNS=100
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=300s

# 自动迁移（生产环境建议设为 false）
ENABLE_AUTO_MIGRATE=false

# 用户服务配置
USER_SERVICE_BASE_URL=http://user-service:8000
USER_SERVICE_TIMEOUT=10s

# 日志配置
LOG_LEVEL=info
LOG_FORMAT=json
LOG_OUTPUT=stdout
```

### 2. 配置文件

生产环境配置 `configs/config.prod.yaml`：

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
    host: ${DB_HOST}
    port: ${DB_PORT}
    user: ${DB_USER}
    password: ${DB_PASSWORD}
    dbname: ${DB_NAME}
    max_open_conns: ${DB_MAX_OPEN_CONNS}
    max_idle_conns: ${DB_MAX_IDLE_CONNS}
    conn_max_lifetime: ${DB_CONN_MAX_LIFETIME}
    enable_auto_migrate: ${ENABLE_AUTO_MIGRATE}

user_service:
  base_url: ${USER_SERVICE_BASE_URL}
  timeout: ${USER_SERVICE_TIMEOUT}

log:
  level: ${LOG_LEVEL}
  format: ${LOG_FORMAT}
  output: ${LOG_OUTPUT}
```

### 3. Kubernetes 部署

**Deployment 配置**：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: room-service
  labels:
    app: room-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: room-service
  template:
    metadata:
      labels:
        app: room-service
    spec:
      containers:
      - name: room-service
        image: minichat-room-service:v1.0.0
        ports:
        - containerPort: 8000
          name: http
        - containerPort: 9000
          name: grpc
        env:
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: room-config
              key: db-host
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: room-secret
              key: db-password
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /v1/rooms
            port: 8000
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /v1/rooms
            port: 8000
          initialDelaySeconds: 5
          periodSeconds: 10

---
apiVersion: v1
kind: Service
metadata:
  name: room-service
spec:
  selector:
    app: room-service
  ports:
  - name: http
    port: 8000
    targetPort: 8000
  - name: grpc
    port: 9000
    targetPort: 9000
  type: LoadBalancer

---
apiVersion: v1
kind: ConfigMap
metadata:
  name: room-config
data:
  db-host: "room-db-postgresql"
  enable-auto-migrate: "false"

---
apiVersion: v1
kind: Secret
metadata:
  name: room-secret
type: Opaque
data:
  db-password: c2VjdXJlLXBhc3N3b3Jk  # base64 encoded
```

**部署命令**：

```bash
# 创建 ConfigMap 和 Secret
kubectl apply -f k8s/config.yaml
kubectl apply -f k8s/secret.yaml

# 部署服务
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml

# 查看部署状态
kubectl get pods -l app=room-service
kubectl logs -f deployment/room-service

# 扩容
kubectl scale deployment room-service --replicas=5
```

---

## 配置管理

### 开发环境配置

```yaml
server:
  http:
    addr: 0.0.0.0:8000
  grpc:
    addr: 0.0.0.0:9000

data:
  database:
    driver: postgres
    host: localhost
    port: 5432
    user: minichat
    password: minichat123
    dbname: room_service
    enable_auto_migrate: true  # 开发环境启用自动迁移

log:
  level: debug
  format: console
```

### 生产环境配置

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
    host: room-db-prod.postgresql
    port: 5432
    user: minichat_prod
    password: ${DB_PASSWORD}
    dbname: room_service_prod
    max_open_conns: 100
    max_idle_conns: 10
    conn_max_lifetime: 300s
    enable_auto_migrate: false  # 生产环境禁用自动迁移

log:
  level: info
  format: json
  output: stdout
```

---

## 监控与日志

### 1. 日志配置

服务支持三种日志级别：
- `debug` - 开发调试
- `info` - 正常运行信息
- `error` - 错误信息
- `warn` - 警告信息

### 2. 日志收集

**Docker 环境**：

```bash
# 查看实时日志
docker-compose logs -f room-service

# 查看最近 100 行
docker-compose logs --tail 100 room-service

# 查看特定时间的日志
docker-compose logs --since 1h room-service
```

**Kubernetes 环境**：

```bash
# 查看所有 Pod 日志
kubectl logs -l app=room-service --all-containers=true

# 查看特定 Pod 日志
kubectl logs room-service-xxx-xxx

# 流式查看日志
kubectl logs -f room-service-xxx-xxx
```

### 3. 性能监控

建议集成以下监控工具：

- **Prometheus** - 指标收集
- **Grafana** - 可视化监控面板
- **Jaeger** - 分布式追踪

**配置 Prometheus 监控**：

```yaml
# 在 main.go 中添加 Prometheus 支持
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// 启动 Prometheus 指标端点
http.Handle("/metrics", promhttp.Handler())
```

### 4. 健康检查

```bash
# HTTP 健康检查
curl http://localhost:8001/v1/rooms

# 检查数据库连接
docker exec room_db psql -U minichat -d room_service -c "SELECT 1"

# 检查 WebSocket 连接
wscat -c "ws://localhost:8001/ws/room/1?user_id=1&username=test"
```

---

## 故障恢复

### 数据库备份

```bash
# 备份数据库
docker exec room_db pg_dump -U minichat room_service > backup.sql

# 恢复数据库
docker exec -i room_db psql -U minichat room_service < backup.sql
```

### 服务重启

```bash
# Docker Compose
docker-compose restart room-service

# Kubernetes
kubectl rollout restart deployment room-service

# 手动重启（本地）
pkill room
./bin/room -conf ./configs &
```

### 数据库迁移

生产环境建议手动执行 SQL 迁移脚本：

```bash
# 生成迁移脚本
go run cmd/migrate/main.go

# 执行迁移
docker exec room_db psql -U minichat -d room_service < migrations/001_init.up.sql
```

---

## 安全配置

### 1. 数据库安全

```yaml
# 使用强密码
DB_PASSWORD: $(openssl rand -base64 32)

# 限制数据库访问
# 在 postgresql.conf 中设置:
listen_addresses = 'localhost'
# 或使用防火墙规则
```

### 2. TLS/SSL 配置

生产环境建议启用 HTTPS：

```yaml
server:
  http:
    addr: 0.0.0.0:443
    tls:
      cert: /path/to/cert.pem
      key: /path/to/key.pem
```

### 3. 认证中间件

建议集成 API Gateway 进行统一认证：

```go
// 在 middleware/auth.go 中添加 JWT 验证
func JWTAuth() middleware.Middleware {
    return func(handler middleware.Handler) middleware.Handler {
        return func(ctx context.Context, req interface{}) (interface{}, error) {
            token := req.Header.Get("Authorization")
            if !validateToken(token) {
                return nil, errors.New("unauthorized")
            }
            return handler(ctx, req)
        }
    }
}
```

---

## 性能优化

### 1. 数据库连接池

```yaml
data:
  database:
    max_open_conns: 100      # 最大打开连接数
    max_idle_conns: 10       # 最大空闲连接数
    conn_max_lifetime: 300s  # 连接最大生命周期
```

### 2. WebSocket 连接管理

```go
// 在 websocket.go 中配置
upgrader: websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    HandshakeTimeout: 10 * time.Second,
}
```

### 3. 缓存策略

建议添加 Redis 缓存：

```go
// 缓存房间信息
func (uc *RoomUsecase) Get(ctx context.Context, id int64) (*Room, error) {
    // 先查缓存
    if cached, err := uc.cache.Get(ctx, fmt.Sprintf("room:%d", id)); err == nil {
        return cached.(*Room), nil
    }
    // 查数据库
    room, err := uc.repo.Get(ctx, id)
    if err == nil {
        uc.cache.Set(ctx, fmt.Sprintf("room:%d", id), room, time.Minute*5)
    }
    return room, err
}
```

---

## 常见问题

### 1. 端口冲突

```bash
# 检查端口占用
lsof -i :8001
lsof -i :9001

# 杀死占用进程
kill -9 <PID>
```

### 2. 数据库连接失败

```bash
# 检查数据库状态
docker ps | grep room_db

# 查看数据库日志
docker logs room-db

# 测试连接
docker exec room_db psql -U minichat -d room_service
```

### 3. 内存泄漏

```bash
# 查看内存使用
docker stats room-service

# 分析内存
go tool pprof http://localhost:8001/debug/pprof/heap
```

---

## 维护建议

1. **定期备份**：每天备份数据库
2. **日志轮转**：配置日志轮转策略
3. **监控告警**：配置关键指标监控和告警
4. **安全更新**：定期更新依赖包
5. **容量规划**：根据负载调整实例数量
