# ==========================================
# 第一阶段：构建环境 (Builder)
# 使用 Go 1.24 解决版本依赖问题
# ==========================================
FROM docker.m.daocloud.io/library/golang:1.24 AS builder

# 1. 设置环境变量 (开启 Module，配置国内七牛云代理加速)
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0

# 2. 安装系统级依赖 (Protoc 编译器)
RUN apt-get update && apt-get install -y --no-install-recommends \
    protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*

# 3. 安装 Go 开发工具链 (方便队友进入容器维护)
RUN go install github.com/go-kratos/kratos/cmd/kratos/v2@latest && \
    go install github.com/google/wire/cmd/wire@latest && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest && \
    go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest

# 4. 设置工作目录
WORKDIR /src

# 5. 预下载依赖 (利用 Docker 缓存层)
COPY go.mod go.sum ./
RUN go mod download

# 6. 复制代码并编译
COPY . .
RUN go build -o /bin/server ./cmd/user

# ==========================================
# 第二阶段：运行环境 (Runner)
# 极简镜像
# ==========================================
FROM docker.m.daocloud.io/library/debian:stable-slim

# 安装基础证书
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates netbase && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# 1. 复制编译好的二进制程序
COPY --from=builder /bin/server ./

# 2. 复制配置文件 (确保路径正确)
COPY --from=builder /src/configs ./configs

# 3. 复制 API 定义 (保留目录结构，防止找不到文件)
COPY --from=builder /src/api ./api

# 4. 暴露端口
EXPOSE 8000
EXPOSE 9000

# 5. 启动命令
CMD ["./server", "-conf", "/app/configs"]
