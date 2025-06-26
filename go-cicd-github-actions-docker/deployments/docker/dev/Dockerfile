# 基础阶段：包含通用依赖
FROM golang:1.24-alpine AS base

# 安装开发工具和依赖
RUN apk add --no-cache git curl make gcc libc-dev &&
    go install github.com/cosmtrek/air@latest &&
    go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app

# 开发阶段：用于热重载的镜像
FROM base AS dev

# 复制应用代码
COPY . .

# 创建logs目录并设置权限
RUN mkdir -p /app/logs && chmod 755 /app/logs

# 设置环境变量
ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOOS=linux
ENV GOARCH=amd64

# 暴露应用端口和调试端口
EXPOSE 18080

# 使用 air 实现热重载
CMD ["air", "-c", ".air.toml"]
