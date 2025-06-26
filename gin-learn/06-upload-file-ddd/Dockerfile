# 使用多阶段构建
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装依赖
RUN apk add --no-cache git

# 复制go.mod和go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制所有源文件
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o file-upload .

# 第二阶段，使用轻量级基础镜像
FROM alpine:latest

# 安装ca-certificates以支持HTTPS
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从builder阶段复制二进制文件
COPY --from=builder /app/file-upload .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/web ./web

# 创建上传目录
RUN mkdir -p uploads temp_chunks

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./file-upload"]
