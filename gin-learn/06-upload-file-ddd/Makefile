.PHONY: run build clean fmt lint test air install-deps clean-uploads help docker docker-dev docker-stop docker-clean

# 定义变量
APP_NAME = file-upload
MAIN_FILE = ./cmd/main.go
BUILD_DIR = ./build
UPLOAD_DIR = ./uploads
TEMP_CHUNK_DIR = ./temp_chunks
DOCKER_IMAGE = file-upload
DOCKER_CONTAINER = file-upload-service
DOCKER_DEV_CONTAINER = file-upload-dev

# 默认目标
all: build

# 帮助信息
help:
	@echo "文件上传系统 Makefile 命令:"
	@echo "make run              - 运行程序"
	@echo "make build            - 构建可执行文件"
	@echo "make clean            - 清理构建和临时文件"
	@echo "make fmt              - 格式化代码"
	@echo "make lint             - 运行静态代码检查"
	@echo "make test             - 运行测试"
	@echo "make air              - 启动 air 进行热重载"
	@echo "make docker           - 构建和运行 Docker 容器"
	@echo "make docker-dev       - 构建和运行开发环境 Docker 容器"
	@echo "make docker-stop      - 停止 Docker 容器"
	@echo "make docker-clean     - 清理 Docker 容器和镜像"
	@echo "make install-deps     - 安装依赖工具"
	@echo "make clean-uploads    - 清理上传的文件"

# 运行程序
run:
	@echo "正在启动文件上传服务器..."
	@mkdir -p $(UPLOAD_DIR) $(TEMP_CHUNK_DIR)
	@go run $(MAIN_FILE)

# 构建可执行文件
build:
	@echo "正在构建$(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "构建完成: $(BUILD_DIR)/$(APP_NAME)"

# 清理构建文件
clean:
	@echo "清理构建文件..."
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "清理完成"

# 格式化代码
fmt:
	@echo "正在格式化代码..."
	@go fmt ./...
	@echo "格式化完成"

# 代码检查
lint:
	@echo "正在进行代码检查..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint 未安装, 请先运行 make install-deps"; \
		exit 1; \
	fi

# 运行测试
test:
	@echo "正在运行测试..."
	@go test -v ./...

# 使用 air 进行热重载
air:
	@echo "正在使用 air 启动服务（热重载）..."
	@if command -v air >/dev/null 2>&1; then \
		mkdir -p $(UPLOAD_DIR) $(TEMP_CHUNK_DIR); \
		air; \
	else \
		echo "air 未安装, 请先运行 make install-deps"; \
		exit 1; \
	fi

# 构建和运行Docker容器
docker:
	@echo "构建并运行 Docker 容器..."
	@docker-compose up -d upload-service
	@echo "容器已启动，访问 http://localhost:8080"

# 构建和运行开发环境Docker容器
docker-dev:
	@echo "构建并运行开发环境 Docker 容器..."
	@docker-compose up -d upload-dev
	@echo "开发容器已启动，访问 http://localhost:8081"

# 停止Docker容器
docker-stop:
	@echo "停止 Docker 容器..."
	@docker-compose down
	@echo "容器已停止"

# 清理Docker容器和镜像
docker-clean:
	@echo "清理 Docker 容器和镜像..."
	@docker-compose down -v
	@docker rmi -f $(DOCKER_IMAGE) || true
	@echo "Docker 清理完成"

# 安装依赖工具
install-deps:
	@echo "正在安装开发依赖..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/cosmtrek/air@latest
	@go mod tidy
	@echo "依赖安装完成"

# 清理上传的文件
clean-uploads:
	@echo "正在清理上传的文件..."
	@rm -rf $(UPLOAD_DIR)/* $(TEMP_CHUNK_DIR)/*
	@echo "清理完成"
