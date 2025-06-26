# Go CI/CD 示例项目

一个使用 Go、Gin 构建的 RESTful API 服务，集成了 MySQL、Redis，并实现了基于 GitHub Actions 的 CI/CD 流程。

## 项目架构

```
├── api/                        # API 文档相关
│   └── openapi/                # OpenAPI/Swagger 文档
├── cmd/                        # 主要应用程序入口
│   └── server/                 # 服务器启动
│       └── main.go             # 主函数
├── configs/                    # 配置文件
│   ├── config.yaml             # 应用配置
│   └── db_init.sql             # 数据库初始化脚本
├── deployments/                # 部署文件
│   ├── docker/                 # Docker 部署文件
│   │   ├── dev/                # 开发环境Docker配置
│   │   │   ├── Dockerfile      # 开发环境Dockerfile
│   │   │   └── docker-compose.yml # 开发环境Docker Compose配置
│   │   └── prod/               # 生产环境Docker配置
│   │       ├── Dockerfile      # 生产环境Dockerfile
│   │       └── docker-compose.yml # 生产环境Docker Compose配置
│   └── shell/                  # 部署脚本
├── docs/                       # 项目文档
├── internal/                   # 私有应用程序代码
│   ├── api/                    # API层
│   │   └── user/               # 用户相关 API 接口
│   ├── config/                 # 配置模型与加载
│   ├── middleware/             # HTTP 中间件
│   ├── model/                  # 数据模型
│   ├── pkg/                    # 内部通用包
│   │   ├── db/                 # 数据库连接管理
│   │   ├── known/              # 常量定义
│   │   └── log/                # 日志工具
│   ├── repository/             # 数据访问层
│   └── service/                # 业务逻辑层
├── scripts/                    # 工具脚本
├── .github/                    # GitHub 相关配置
│   └── workflows/              # GitHub Actions 工作流
│       └── deploy.yml          # CI/CD 工作流
├── docker-compose.dev-env.yml  # 开发环境容器配置
├── .air.toml                   # Air 热重载配置
├── .editorconfig               # 编辑器配置
└── Makefile                    # 项目管理命令
```

## 技术栈

- **后端框架**: [Gin](https://github.com/gin-gonic/gin)
- **数据库**:
  - MySQL: 用户数据持久化存储
  - Redis: 缓存和会话管理
- **项目管理**: Makefile
- **文档**: Swagger/OpenAPI
- **配置管理**: Viper
- **日志**: Zap
- **CI/CD**: GitHub Actions
- **容器化**: Docker & Docker Compose
- **热重载**: [Air](https://github.com/cosmtrek/air)

## 功能特性

- RESTful API设计
- 用户认证与授权
- 基于MySQL的数据持久化
- Redis缓存支持
- 请求ID追踪
- 结构化日志
- API文档自动生成
- 容器化部署
- CI/CD流水线

## 开始使用

### 前提条件

- Go 1.24 或更高版本
- Docker 和 Docker Compose
- Make

### 本地开发环境

1. 克隆仓库

   ```bash
   git clone https://github.com/clin211/go-cicd-github-actions-docker.git
   cd go-cicd-github-actions-docker
   ```

2. 启动开发环境数据库

   ```bash
   docker-compose -f docker-compose.dev-env.yml up -d
   ```

3. 运行应用

   ```bash
   make run
   ```

   或使用热重载:

   ```bash
   make dev
   ```

4. API 将在 <http://localhost:18080> 可用

### Makefile 命令

```bash
# 编译项目
make build

# 运行项目
make run

# 使用热重载运行开发环境
make dev

# 运行测试
make test

# 清理生成文件
make clean

# 生成 Swagger 文档
make swagger

# 构建开发环境Docker镜像
make docker-dev

# 构建生产环境Docker镜像
make docker-prod

# 使用Docker运行开发环境
make docker-compose-dev

# 使用Docker运行生产环境
make docker-compose-prod

# 初始化数据库
make init-db
```

## API 端点

### 用户管理

- `GET /api/v1/users` - 获取用户列表
- `GET /api/v1/users/:id` - 获取单个用户
- `POST /api/v1/users` - 创建用户
- `PUT /api/v1/users/:id` - 更新用户
- `DELETE /api/v1/users/:id` - 删除用户

### 认证

- `POST /api/v1/auth/signin` - 用户登录
- `POST /api/v1/auth/signup` - 用户注册

### 系统

- `GET /api/v1/health` - 健康检查

## 数据库设计

项目使用MySQL数据库，主要表结构:

- `user` - 用户信息
- `role` - 角色定义
- `user_role` - 用户角色关联
- `access_token` - 访问令牌

数据库初始化脚本位于 `configs/db_init.sql`

## CI/CD 流程

项目使用GitHub Actions实现CI/CD自动化，工作流程:

1. 代码推送触发工作流
2. 运行测试
3. 构建Docker镜像
4. 推送镜像到Docker Hub
5. 部署到目标环境

配置文件: `.github/workflows/deploy.yml`

### GitHub Actions Secrets 配置

部署过程需要配置以下secrets变量：

| 变量名 | 描述 |
|--------|------|
| `ALIYUN_USERNAME` | 阿里云容器镜像服务的用户名 |
| `ALIYUN_PASSWORD` | 阿里云容器镜像服务的密码 |
| `SERVER_HOST` | 部署服务器的主机地址 |
| `SERVER_USER` | 部署服务器的用户名 |
| `SERVER_SSH_KEY` | 部署服务器的 SSH 私钥 |

这些变量需要在 GitHub 仓库的 Settings -> Secrets and variables -> Actions 中进行配置，以确保 CI/CD 流程能够正常执行。

## 项目结构说明

- **分层架构**: 采用仓库模式和服务层模式，实现关注点分离
- **依赖注入**: 使用构造函数注入各层依赖
- **配置管理**: 使用Viper管理配置，支持环境变量覆盖
- **错误处理**: 统一的错误处理机制
- **中间件**: 请求日志、跨域、安全头等
- **数据库连接池**: 支持主从读写分离

## 贡献指南

1. Fork 仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 许可证

本项目基于 MIT 许可证开源
