# Super Web Server

A high-performance web server built with Go, featuring a clean architecture and comprehensive functionality for modern web applications.

## 🚀 Features

- **RESTful API**: Clean and well-structured API endpoints
- **JWT Authentication**: Secure user authentication and authorization
- **Role-based Access Control**: Fine-grained permission management
- **Database Integration**: MySQL with GORM ORM and Redis caching
- **Middleware Support**: CORS, logging, recovery, and custom middleware
- **Multi-environment Configuration**: Support for dev, prod, test, and local modes
- **Structured Logging**: Comprehensive logging with rotation and compression
- **Graceful Shutdown**: Proper server shutdown handling
- **Data Validation**: Request validation with custom error messages
- **Unique ID Generation**: Snowflake algorithm for distributed ID generation
- **Internationalization**: Multi-language support

## 🏗️ Architecture

```
├── cmd/server/          # Application entry point
├── configs/             # Configuration files
├── internal/
│   ├── api/v1/         # API version 1 routes
│   ├── app/            # Application core
│   ├── config/         # Configuration management
│   ├── controller/     # HTTP handlers
│   ├── dto/            # Data transfer objects
│   ├── middleware/     # HTTP middleware
│   ├── model/          # Database models
│   ├── repo/           # Data access layer
│   ├── service/        # Business logic layer
│   └── validator/      # Request validation
├── pkg/                # Shared packages
│   ├── database/       # Database utilities
│   ├── jwt/            # JWT utilities
│   ├── logger/         # Logging utilities
│   ├── redis/          # Redis utilities
│   └── utils/          # Common utilities
└── static/             # Static files
```

## 🛠️ Tech Stack

- **Language**: Go 1.24.2
- **Web Framework**: Gin
- **Database**: MySQL with GORM
- **Cache**: Redis
- **Authentication**: JWT
- **Logging**: Uber Zap
- **Configuration**: Viper
- **Validation**: go-playground/validator
- **ID Generation**: Snowflake

## 📋 Prerequisites

- Go 1.24.2 or later
- MySQL 8.0 or later
- Redis 6.0 or later

## 🚀 Quick Start

### 1. Clone the repository

```bash
git clone https://github.com/yourusername/super-web-server.git
cd super-web-server
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Set up configuration

Copy the example configuration file and modify it according to your environment:

```bash
cp configs/config.dev.yml configs/config.local.yml
```

### 4. Set up database

Create a MySQL database and update the database configuration in your config file:

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

### 5. Set up Redis

Ensure Redis is running and update the Redis configuration:

```yaml
redis:
  host: localhost
  port: 6379
  password: your_redis_password
  db: 0
```

### 6. Run the server

```bash
# Development mode
go run cmd/server/server.go -mode dev

# Production mode
go run cmd/server/server.go -mode prod

# Or build and run
go build -o bin/server cmd/server/server.go
./bin/server -mode dev
```

The server will start on `http://localhost:8080` by default.

## 📚 API Documentation

### Authentication

#### Login by Email
```bash
POST /api/v1/user/login-by-email
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "your_password"
}
```

#### Get User Info (Protected)
```bash
GET /api/v1/user/info
Authorization: Bearer <your_jwt_token>
```

### Health Check

```bash
GET /api/v1/hello
```

## ⚙️ Configuration

The application supports multiple environment configurations:

- `dev`: Development environment
- `prod`: Production environment
- `test`: Testing environment
- `local`: Local development environment

Configuration files are located in the `configs/` directory and follow the naming convention `config.{mode}.yml`.

### Configuration Options

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

## 🔧 Development

### Running Tests

```bash
go test ./...
```

### Code Formatting

```bash
go fmt ./...
```

### Linting

```bash
golangci-lint run
```

### Database Migration

The application includes database migration and seeding functionality. Models are automatically migrated on startup.

## 🐳 Docker Support

Create a `Dockerfile` for containerization:

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

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [GORM](https://gorm.io/) - ORM library
- [Viper](https://github.com/spf13/viper) - Configuration management
- [Zap](https://github.com/uber-go/zap) - Structured logging
- [JWT](https://github.com/golang-jwt/jwt) - JSON Web Tokens

## 📞 Support

For support, please create an issue in the GitHub repository or contact the maintainers.


## TODO

- [ ] RabbitMQ
- [ ] Cron