# 使用官方 Go 镜像作为基础（包含完整的 Go 编译器）
# Go mirror is in need to configure with compelete compiler 
#FROM golang:1.24
#国内也下载不了GO 
FROM docker.m.daocloud.io/library/golang:1.24
# 1. 设置环境变量
# 开启 Module 模式，并配置国内代理（加速下载）
# 1, Set up enviroment path; then open Module mode, this confiuration has been set for China, please chage it if you are out of China 
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \ 
    CGO_ENABLED=0

# 2. 安装系统级依赖
# 主要是 Protobuf 编译器，Kratos 生成代码需要它
#2, install system dependencies, mainly compiler of Protobuf wihch Kratos need to generate codes
RUN apt-get update && apt-get install -y --no-install-recommends \
    protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*

# 3. 安装 Go 开发工具链 (Kratos, Wire, Protoc Plugins)
# 这样组员不需要在本地安装这些，容器里都有
#3, install developing chains for Go, such as Kratos, Wire and Plugins, which woulb be set up fine in the container

RUN go install github.com/go-kratos/kratos/cmd/kratos/v2@latest && \
    go install github.com/google/wire/cmd/wire@latest && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest && \
    go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest

# 4. 设置工作目录
#4, set working path
WORKDIR /app

# 5. 这一步很关键：
# 我们不在这里 COPY 代码，而是通过 Docker Compose 的 "挂载" (Volume) 功能
# 让容器直接读取你硬盘上的代码。这样你修改本地代码，容器里也会跟着变。
#5, Attention
# Docker Compose will read codes through a function of volume, instead of copying codes directly. Once you change codes in the local, that in the container would change along with it  


# 6. 容器启动时默认执行的命令
# 默认启动服务，但也可以被覆盖用来执行 wire 或 protoc 命令
#6, commands to be executed as default when the container starts. it can also be used to cover wire or protoc commands 
CMD ["kratos", "run"]
