# ---------- 构建阶段 ----------
FROM docker.m.daocloud.io/library/golang:1.24-alpine AS builder

# 国内网络必备
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn

WORKDIR /src

# 先复制 mod 文件，利用 Docker cache
COPY go.mod go.sum ./
RUN go mod download

# 再复制全部源码
COPY . .
COPY cmd ./cmd

# 构建
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o user-service ./cmd/registration

# ---------- 运行阶段 ----------
FROM docker.m.daocloud.io/library/alpine:3.19

WORKDIR /app

COPY --from=builder /src/user-service .

EXPOSE 8000 9000

CMD ["./user-service"]

