# 构建阶段
FROM golang:1.24-alpine AS builder

# 安装构建依赖
RUN apk add --no-cache git

# 设置工作目录
WORKDIR /app

# 首先仅复制依赖文件以利用缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o server ./cmd/server

# 最终阶段
FROM alpine:3.15

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN adduser -D -g '' appuser

# 复制编译好的应用
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/configs ./configs

# 创建logs目录并设置权限
RUN mkdir -p /app/logs && chown -R appuser:appuser /app/logs && chmod 755 /app/logs

# 使用非root用户运行
USER appuser

# 声明端口
EXPOSE 18080

# 启动命令
CMD ["./server"]
