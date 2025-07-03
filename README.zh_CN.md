# Super Web Server

一个基于 Go 语言构建的高性能 Web 服务器，具有清晰的架构和现代 Web 应用所需的完整功能。

## 🚀 特性

- **RESTful API**: 清晰且结构良好的 API 端点
- **JWT 认证**: 安全的用户身份验证和授权
- **基于角色的访问控制**: 细粒度的权限管理
- **数据库集成**: MySQL 配合 GORM ORM 和 Redis 缓存
- **中间件支持**: CORS、日志记录、异常恢复和自定义中间件
- **多环境配置**: 支持开发、生产、测试和本地模式
- **结构化日志**: 完整的日志记录，支持轮转和压缩
- **优雅关闭**: 正确的服务器关闭处理
- **数据验证**: 请求验证，支持自定义错误消息
- **唯一ID生成**: 使用雪花算法生成分布式唯一ID
- **国际化支持**: 多语言支持

## 🏗️ 架构

```
├── cmd/server/          # 应用程序入口点
├── configs/             # 配置文件
├── internal/
│   ├── api/v1/         # API 版本 1 路由
│   ├── app/            # 应用程序核心
│   ├── config/         # 配置管理
│   ├── controller/     # HTTP 处理器
│   ├── dto/            # 数据传输对象
│   ├── middleware/     # HTTP 中间件
│   ├── model/          # 数据库模型
│   ├── repo/           # 数据访问层
│   ├── service/        # 业务逻辑层
│   └── validator/      # 请求验证
├── pkg/                # 共享包
│   ├── database/       # 数据库工具
│   ├── jwt/            # JWT 工具
│   ├── logger/         # 日志工具
│   ├── redis/          # Redis 工具
│   └── utils/          # 通用工具
└── static/             # 静态文件
```

## 🛠️ 技术栈

- **语言**: Go 1.24.2
- **Web 框架**: Gin
- **数据库**: MySQL 配合 GORM
- **缓存**: Redis
- **认证**: JWT
- **日志**: Uber Zap
- **配置**: Viper
- **验证**: go-playground/validator
- **ID 生成**: Snowflake

## 📋 前置要求

- Go 1.24.2 或更高版本
- MySQL 8.0 或更高版本
- Redis 6.0 或更高版本

## 🚀 快速开始

### 1. 克隆仓库

```bash
git clone https://github.com/yourusername/super-web-server.git
cd super-web-server
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 设置配置

复制示例配置文件并根据您的环境进行修改：

```bash
cp configs/config.dev.yml configs/config.local.yml
```

### 4. 设置数据库

创建 MySQL 数据库并更新配置文件中的数据库配置：

```yaml
database:
  host: localhost
  port: 3306
  username: your_username
  password: your_password
  database: super_db
  charset: utf8mb4
  timezone: UTC
```

### 5. 设置 Redis

确保 Redis 正在运行并更新 Redis 配置：

```yaml
redis:
  host: localhost
  port: 6379
  password: your_redis_password
  db: 0
```

### 6. 运行服务器

```bash
# 开发模式
go run cmd/server/server.go -mode dev

# 生产模式
go run cmd/server/server.go -mode prod

# 或者构建后运行
go build -o bin/server cmd/server/server.go
./bin/server -mode dev
```

服务器默认将在 `http://localhost:8080` 启动。

## 📚 API 文档

### 认证

#### 邮箱登录
```bash
POST /api/v1/user/login-by-email
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "your_password"
}
```

#### 获取用户信息（需要认证）
```bash
GET /api/v1/user/info
Authorization: Bearer <your_jwt_token>
```

### 健康检查

```bash
GET /api/v1/hello
```

## ⚙️ 配置

应用程序支持多种环境配置：

- `dev`: 开发环境
- `prod`: 生产环境
- `test`: 测试环境
- `local`: 本地开发环境

配置文件位于 `configs/` 目录中，遵循 `config.{mode}.yml` 的命名约定。

### 配置选项

```yaml
server:
  port: 8080
  readTimeout: 30s
  writeTimeout: 30s
  maxHeaderBytes: 1048576
  snowflakeNode: 1

log:
  level: info
  path: ./logs
  filename: app.log
  maxSize: 100
  maxBackups: 30
  maxAge: 7
  compress: true
  stdout: true

database:
  host: localhost
  port: 3306
  username: root
  password: root
  database: super_db
  charset: utf8mb4
  timezone: UTC
  maxIdleConns: 10
  maxOpenConns: 100
  connMaxLifetime: 10s
  logLevel: info
  parseTime: true

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: your-jwt-secret-key
  expire: 24h
  issuer: super-web-server
```

## 🔧 开发

### 运行测试

```bash
go test ./...
```

### 代码格式化

```bash
go fmt ./...
```

### 代码检查

```bash
golangci-lint run
```

### 数据库迁移

应用程序包含数据库迁移和种子数据功能。模型会在启动时自动迁移。

## 🐳 Docker 支持

创建 `Dockerfile` 进行容器化：

```dockerfile
FROM golang:1.24.2-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/server cmd/server/server.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/bin/server .
COPY --from=builder /app/configs ./configs

EXPOSE 8080
CMD ["./server", "-mode", "prod"]
```

## 🤝 贡献

1. Fork 仓库
2. 创建您的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开一个 Pull Request

## 📝 许可证

本项目采用 MIT 许可证 - 请查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [Gin](https://github.com/gin-gonic/gin) - HTTP Web 框架
- [GORM](https://gorm.io/) - ORM 库
- [Viper](https://github.com/spf13/viper) - 配置管理
- [Zap](https://github.com/uber-go/zap) - 结构化日志
- [JWT](https://github.com/golang-jwt/jwt) - JSON Web 令牌

## 📞 支持

如需支持，请在 GitHub 仓库中创建 issue 或联系维护人员。
