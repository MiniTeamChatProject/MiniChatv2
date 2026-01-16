# ---------- 构建阶段 ----------
FROM docker.m.daocloud.io/library/golang:1.24-alpine AS builder

ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn

WORKDIR /src

# 先复制依赖文件（利用缓存）
COPY go.mod go.sum ./
RUN go mod download

# 再复制全部源码
COPY . .

# 构建 Kratos 服务
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o user-service ./cmd/registration


# ---------- 运行阶段 ----------
FROM docker.m.daocloud.io/library/alpine:3.19

WORKDIR /app

# 拷贝二进制
COPY --from=builder /src/user-service ./user-service

# 拷贝配置目录（关键）
COPY configs ./configs

EXPOSE 8000 9000

# 明确指定配置路径（绝对路径，100%稳定）
CMD ["./user-service", "-conf", "/app/configs"]

