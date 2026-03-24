FROM golang:1.24-alpine AS builder

COPY . /src
WORKDIR /src

RUN apk add --no-cache make git && \
    GOPROXY=https://goproxy.cn make build

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /src/bin /app
COPY --from=builder /src/configs /app/configs

WORKDIR /app

EXPOSE 8000
EXPOSE 9000

CMD ["./room", "-conf", "./configs"]
