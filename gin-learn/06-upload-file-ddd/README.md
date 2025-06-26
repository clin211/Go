# Gin文件上传系统 (DDD架构)

这个项目是一个基于Gin框架的文件上传系统，使用领域驱动设计(DDD)架构实现，支持多种上传方式：

- 单文件上传
- 多文件批量上传
- 文件夹上传（保持目录结构）
- 大文件分片上传

## 项目结构

```
file-upload/
│
├── cmd/                            # 应用入口
│   └── main.go
│
├── configs/                        # 配置文件
│   └── config.yaml
│
├── internal/
│   ├── upload/                     # [限界上下文：文件上传]
│   │   ├── domain/                 # 领域层
│   │   │   ├── model.go            # 实体（File、ChunkFile、Folder）
│   │   │   ├── value_object.go     # 值对象（FileID、Path 等）
│   │   │   └── service.go          # 领域服务（合并分片等）
│   │   │
│   │   ├── application/            # 应用层
│   │   │   └── upload_service.go   # 应用服务，处理用例流程
│   │   │
│   │   ├── infrastructure/         # 基础设施层
│   │   │   └── storage/
│   │   │       └── local.go        # 本地存储实现
│   │   │
│   │   └── interfaces/             # 接口层
│   │       └── http/
│   │           ├── handler.go      # 上传请求处理
│   │           └── routes.go       # 上传路由注册
│   │
│   └── shared/                     # 通用工具
│       ├── logger/                 # 日志
│       ├── response/               # 响应封装
│       └── errors/                 # 错误处理
│
├── pkg/                            # 通用工具库
│   └── utils/
│       └── file.go                 # 文件操作工具
│
├── web/                            # 前端页面
│   └── index.html                  # 测试上传的示例页面
│
├── Dockerfile                      # 生产环境Docker配置
├── Dockerfile.dev                  # 开发环境Docker配置
├── docker-compose.yml              # Docker Compose配置
├── Makefile                        # 项目管理脚本
├── .air.toml                       # Air热重载配置
├── go.mod
└── README.md
```

## DDD架构分层

| 层级       | 职责说明 |
|------------|----------|
| **领域层** | 定义核心模型与业务规则：文件对象、文件夹结构、分片合并逻辑等 |
| **应用层** | 协调上传流程：协调各种上传用例的执行 |
| **基础设施层** | 实现对存储系统的访问：本地文件系统或对象存储 |
| **接口层** | 接收请求并返回响应：处理HTTP请求，调用应用服务 |

## 功能特性

- **多种上传模式**: 单文件、多文件、文件夹、分片上传
- **多样化存储引擎**: 本地文件系统存储、阿里云OSS对象存储
- **完整的错误处理**: 标准化的错误响应和日志记录
- **可配置**: 通过YAML配置文件进行系统配置

## Makefile 使用指南

项目包含一个Makefile，提供了多种开发和管理命令：

```bash
# 运行程序（正常模式）
make run

# 使用air进行热重载开发
make air

# 构建可执行文件
make build

# 格式化代码
make fmt

# 运行代码检查
make lint

# 运行测试
make test

# 使用Docker部署
make docker

# 使用Docker部署开发环境
make docker-dev

# 停止Docker容器
make docker-stop

# 清理Docker容器和镜像
make docker-clean

# 安装开发依赖工具(golangci-lint, air)
make install-deps

# 清理上传的文件
make clean-uploads

# 清理构建文件
make clean

# 显示帮助信息
make help
```

## Docker 部署

项目提供了Docker部署支持，有两种模式：

### 生产模式

```bash
# 使用Docker Compose启动生产模式
make docker

# 或者直接使用docker-compose
docker-compose up -d upload-service
```

生产模式访问地址：<http://localhost:8080>

### 开发模式

开发模式支持热重载，当代码变更时会自动重新编译：

```bash
# 使用Docker Compose启动开发模式
make docker-dev

# 或者直接使用docker-compose
docker-compose up -d upload-dev
```

开发模式访问地址：<http://localhost:8081>

### 停止容器

```bash
make docker-stop
```

### 清理Docker环境

```bash
make docker-clean
```

## 快速开始

1. 克隆代码库
2. 配置 `configs/config.yaml`
3. 安装开发依赖

```bash
make install-deps
```

4. 开始开发（使用热重载）

```bash
make air
```

5. 或者直接运行

```bash
make run
```

6. 访问 <http://localhost:8080> 使用Web界面测试上传功能

## API文档

### 单文件上传

```
POST /upload/single
```

### 多文件上传

```
POST /upload/multiple
```

### 文件夹上传

```
POST /upload/folder
```

### 分片上传

初始化上传:

```
POST /upload/chunk/init
```

上传分片:

```
POST /upload/chunk
```

合并分片:

```
POST /upload/chunk/merge
```
