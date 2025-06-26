.PHONY: all build run clean test lint swagger docker-dev docker-prod docker-compose-dev docker-compose-prod help dev deps install-tools init-db check-db

# 全局变量
APP_NAME=go-cicd-github-actions-docker
ROOT_DIR := $(shell pwd)
BUILD_DIR=$(ROOT_DIR)/build
MAIN_FILE=$(ROOT_DIR)/cmd/server/main.go
DOCKER_IMAGE_NAME=clin211/go-cicd-github-actions-docker
COVERAGE_FILE=$(ROOT_DIR)/coverage.out
GO_VERSION=1.24

# 默认目标
all: clean build

# 帮助命令
help:
	@echo "用户服务项目管理命令"
	@echo ""
	@echo "使用方法:"
	@echo "    make [target]"
	@echo ""
	@echo "可用的目标:"
	@echo "    build              - 编译项目"
	@echo "    run                - 运行项目"
	@echo "    dev                - 运行开发环境(自动重载)"
	@echo "    deps               - 更新依赖"
	@echo "    clean              - 清理编译文件"
	@echo "    test               - 运行测试"
	@echo "    cover              - 运行测试覆盖率"
	@echo "    lint               - 运行代码检查"
	@echo "    swagger            - 生成Swagger文档"
	@echo "    install-tools      - 安装开发工具"
	@echo "    docker-dev         - 构建开发环境Docker镜像"
	@echo "    docker-prod        - 构建生产环境Docker镜像"
	@echo "    docker-compose-dev - 运行开发环境Docker Compose"
	@echo "    docker-compose-prod- 运行生产环境Docker Compose"
	@echo "    help               - 显示帮助信息"
	@echo "    init-db            - 初始化数据库"
	@echo "    check-db           - 检查数据库状态"
	@echo ""

# 编译项目
build:
	@echo "编译项目..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "编译完成: $(BUILD_DIR)/$(APP_NAME)"

# 运行项目
run:
	@echo "运行项目..."
	@go run $(MAIN_FILE)

# 开发环境(自动重载)
dev:
	@echo "启动开发环境..."
	@mkdir -p tmp
	@which air > /dev/null || (echo "安装 air..." && go install github.com/cosmtrek/air@latest)
	@echo "启动热重载服务，按 Ctrl+C 停止..."
	@GOFLAGS="-mod=mod" air -c .air.toml

# 更新依赖
deps:
	@echo "更新依赖..."
	@go mod tidy
	@go mod vendor
	@echo "依赖更新完成"

# 清理编译文件
clean:
	@echo "清理编译文件..."
	@rm -rf $(BUILD_DIR)
	@rm -rf tmp
	@rm -f $(COVERAGE_FILE)
	@echo "清理完成"

# 运行测试
test:
	@echo "运行测试..."
	@go test -v ./...

# 运行测试覆盖率
cover:
	@echo "运行测试覆盖率分析..."
	@go test -cover -coverprofile=$(COVERAGE_FILE) ./...
	@go tool cover -html=$(COVERAGE_FILE)
	@echo "测试覆盖率分析完成"

# 运行代码检查
lint:
	@echo "检查代码格式..."
	@! gofmt -d ./... | grep '^'
	@echo "运行golint..."
	@which golint > /dev/null || (echo "安装 golint..." && go install golang.org/x/lint/golint@latest)
	@golint ./...
	@echo "运行go vet..."
	@go vet ./...

# 生成Swagger文档
swagger:
	@echo "生成Swagger文档..."
	@which swag > /dev/null || (echo "安装 swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
	@if [ ! -d "$(ROOT_DIR)/api/openapi/docs" ]; then mkdir -p $(ROOT_DIR)/api/openapi/docs; fi
	@cd $(ROOT_DIR) && swag init -g cmd/server/main.go -o api/openapi/docs
	@which swagger > /dev/null || (echo "安装 swagger..." && go install github.com/go-swagger/go-swagger/cmd/swagger@latest)
	@swagger serve -F=swagger --no-open --port 65534 $(ROOT_DIR)/api/openapi/docs/swagger.yaml
	@echo "Swagger文档生成完成"

# 安装开发工具
install-tools:
	@echo "安装开发工具..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/cosmtrek/air@latest
	@go install golang.org/x/lint/golint@latest
	@echo "开发工具安装完成"

# 构建开发环境Docker镜像
docker-dev:
	@echo "构建开发环境Docker镜像..."
	@docker build -t $(DOCKER_IMAGE_NAME):dev -f deployments/docker/dev/Dockerfile .
	@echo "开发环境Docker镜像构建完成: $(DOCKER_IMAGE_NAME):dev"

# 构建生产环境Docker镜像
docker-prod:
	@echo "构建生产环境Docker镜像..."
	@docker build -t $(DOCKER_IMAGE_NAME):prod -f deployments/docker/prod/Dockerfile .
	@echo "生产环境Docker镜像构建完成: $(DOCKER_IMAGE_NAME):prod"

# 运行开发环境Docker Compose
docker-compose-dev:
	@echo "运行开发环境Docker Compose..."
	@docker-compose -f deployments/docker/dev/docker-compose.yml up -d
	@echo "开发环境Docker Compose启动完成"

# 运行生产环境Docker Compose
docker-compose-prod:
	@echo "运行生产环境Docker Compose..."
	@docker-compose -f deployments/docker/prod/docker-compose.yml up -d
	@echo "生产环境Docker Compose启动完成"

# 初始化数据库
init-db:
	@echo "初始化数据库..."
	@chmod +x ./scripts/init-db.sh
	@if ! docker ps | grep -q "mysql"; then \
		echo "MySQL容器未运行，正在启动容器..."; \
		docker-compose -f docker-compose.dev-env.yml up -d; \
		echo "等待容器启动完成..."; \
		sleep 10; \
	fi
	@./scripts/init-db.sh
	@echo "数据库初始化完成"

# 检查数据库状态
check-db:
	@echo "检查数据库状态..."
	@chmod +x ./scripts/check-db.sh
	@./scripts/check-db.sh
