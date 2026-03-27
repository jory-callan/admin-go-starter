# pkg 层封装说明

本文档介绍 pkg 目录下各基础设施组件的封装和使用方法。

## 目录结构

```
pkg/
├── database/          # 数据库封装（GORM）
├── discovery/         # 服务发现封装（Consul）
├── email/            # 邮件服务封装（SMTP）
├── http/             # HTTP 服务器封装（Gin）
├── jwt/              # JWT 认证封装
├── kafka/            # Kafka 消息队列封装
├── logger/           # 日志封装（slog）
├── rabbitmq/         # RabbitMQ 消息队列封装
├── redis/            # Redis 封装
└── tracing/          # 链路追踪封装（OpenTelemetry）
```

## 1. Database - 数据库封装

### 文件结构
- `config.go` - 数据库配置
- `database.go` - 数据库连接管理

### 使用示例

```go
import "aicode/pkg/database"

// 加载配置
cfg := database.DefaultConfig()
cfg.Driver = "mysql"
cfg.DSN = "root:password@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=True"

// 打开连接
db, err := database.Open(cfg)
if err != nil {
    log.Fatal(err)
}
defer db.Close()
```

### 支持的数据库类型
- MySQL
- PostgreSQL
- SQLite

## 2. Redis - Redis 客户端封装

### 文件结构
- `config.go` - Redis 配置
- `redis.go` - Redis 连接管理

### 使用示例

```go
import "aicode/pkg/redis"

cfg := redis.DefaultConfig()
cfg.Addr = "127.0.0.1:6379"
cfg.Password = "" // 如果有密码
cfg.DB = 0

client, err := redis.Open(cfg)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// 使用 Redis 客户端
err = client.Set(ctx, "key", "value", 0).Err()
```

## 3. Logger - 结构化日志封装

### 文件结构
- `config.go` - 日志配置
- `logger.go` - 日志实例管理

### 使用示例

```go
import "aicode/pkg/logger"

cfg := logger.DefaultConfig()
cfg.Level = "info"
cfg.Format = "json"
cfg.FilePath = "/var/log/app.log"

log := logger.New(cfg, "my-app")

// 使用全局默认 logger
log.Info("application started")
log.Error("something failed", "error", err)

// 创建子 logger（带 component 标记）
dbLog := logger.C("database")
dbLog.Info("connected to database")
```

### 特性
- 支持 JSON 和文本格式
- 支持日志轮转（lumberjack）
- 支持按级别过滤

## 4. JWT - JWT 认证封装

### 文件结构
- `config.go` - JWT 配置
- `jwt.go` - Token 生成和解析

### 使用示例

```go
import "aicode/pkg/jwt"

// 初始化 JWT（应用启动时调用一次）
cfg := jwt.DefaultConfig()
cfg.Secret = "your-secret-key"
cfg.Expires = "24h"
jwt.Init(cfg)

// 生成 Token
token, err := jwt.GenerateToken(
    "user-123",
    "username",
    []string{"admin", "user"},
    []string{"read", "write"},
)

// 解析 Token
claims, err := jwt.ParseToken(tokenString)
if err != nil {
    log.Fatal(err)
}
fmt.Println(claims.UserID)
```

## 5. Kafka - Kafka 消息队列封装

### 文件结构
- `config.go` - Kafka 配置
- `client.go` - Kafka 客户端（管理生产者和消费者）
- `producer.go` - 生产者封装
- `consumer.go` - 消费者封装

### 使用示例

#### 生产者

```go
import "aicode/pkg/kafka"

cfg := kafka.DefaultConfig()
cfg.Brokers = []string{"127.0.0.1:9092"}
cfg.Topic = "app_events"

client, err := kafka.NewClient(cfg)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

producer := kafka.NewProducer(client)

// 发送简单消息
err = producer.Send(ctx, "my-topic", kafka.Message{
    Key:   "user-event",
    Value: []byte("user created"),
})

// 发送 JSON 消息
err = producer.SendJSON(ctx, "my-topic", "user-123", map[string]interface{}{
    "event": "user_created",
    "data": map[string]interface{}{
        "id": 123,
        "name": "John",
    },
})

// 批量发送
messages := []kafka.Message{...}
err = producer.SendBatch(ctx, "my-topic", messages)
```

#### 消费者

```go
// 创建消费者
consumerCfg := kafka.ConsumerConfig{
    Topic:   "app_events",
    GroupID: "app_consumer",
    Workers: 3, // 并发 worker 数量
    Handler: func(ctx context.Context, msg kafka.Message) error {
        // 处理消息
        fmt.Printf("Received: %s\n", string(msg.Value))
        
        // 解析 JSON 消息
        var data map[string]interface{}
        if err := msg.UnmarshalValue(&data); err != nil {
            return err
        }
        
        return nil
    },
}

consumer := kafka.NewConsumer(client, consumerCfg)

// 启动消费者
if err := consumer.Start(); err != nil {
    log.Fatal(err)
}

// 优雅关闭
defer consumer.Stop()
```

## 6. RabbitMQ - RabbitMQ 消息队列封装

### 文件结构
- `config.go` - RabbitMQ 配置
- `client.go` - RabbitMQ 客户端
- `producer.go` - 生产者封装
- `consumer.go` - 消费者封装

### 使用示例

#### 生产者

```go
import "aicode/pkg/rabbitmq"

cfg := rabbitmq.DefaultConfig()
cfg.URL = "amqp://guest:guest@localhost:5672/"
cfg.Exchange = "app_events"

client, err := rabbitmq.NewClient(cfg)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

producer := rabbitmq.NewProducer(client, cfg.Exchange)

// 发送消息
err = producer.SendJSON(ctx, "routing.key", map[string]interface{}{
    "event": "order_created",
    "data": map[string]interface{}{
        "order_id": 456,
    },
})

// 发送消息并等待确认
err = producer.SendWithConfirm(ctx, "routing.key", []byte(`{"event":"test"}`))
```

#### 消费者

```go
consumerCfg := rabbitmq.ConsumerConfig{
    Queue:         "app_queue",
    Consumer:      "app_consumer",
    AutoAck:       false,
    PrefetchCount: 10,
    Workers:       3,
    Handler: func(ctx context.Context, msg *rabbitmq.Message) error {
        // 处理消息
        fmt.Printf("Received: %s\n", string(msg.Body))
        
        // 解析 JSON
        var data map[string]interface{}
        if err := msg.UnmarshalMessage(&data); err != nil {
            return err
        }
        
        return nil
    },
}

consumer := rabbitmq.NewConsumer(client, consumerCfg)

if err := consumer.Start(); err != nil {
    log.Fatal(err)
}

defer consumer.Stop()
```

## 7. Email - 邮件服务封装

### 文件结构
- `config.go` - 邮件服务配置
- `client.go` - 邮件客户端

### 使用示例

```go
import "aicode/pkg/email"

cfg := email.DefaultConfig()
cfg.SMTPHost = "smtp.example.com"
cfg.SMTPPort = 587
cfg.Username = "noreply@example.com"
cfg.Password = "your-password"
cfg.FromName = "My App"
cfg.UseTLS = true

client, err := email.NewClient(cfg)
if err != nil {
    log.Fatal(err)
}

// 发送简单邮件
err = client.SendSimple(
    []string{"user@example.com"},
    "Welcome",
    "Thank you for registering!",
)

// 发送 HTML 邮件
htmlBody := `
<html>
<body>
<h1>Welcome!</h1>
<p>Thank you for registering.</p>
</body>
</html>
`
err = client.SendHTML(
    []string{"user@example.com"},
    "Welcome",
    htmlBody,
)

// 发送带附件的邮件（高级用法）
err = client.Send(&email.EmailMessage{
    To:      []string{"user@example.com"},
    Subject: "Report",
    Body:    "Please find the attached report.",
    ContentType: "text/plain",
    Attachments: []*email.EmailAttachment{
        {
            Filename: "report.pdf",
            Data:     pdfData,
        },
    },
})
```

## 8. HTTP - HTTP 服务器封装（基于 Gin）

### 文件结构
- `config.go` - HTTP 服务配置
- `server.go` - HTTP 服务器
- `middleware.go` - 中间件（限流、超时、恢复）
- `cors.go` - CORS 跨域中间件

### 使用示例

```go
import "aicode/pkg/http"

cfg := http.DefaultConfig()
cfg.Host = "0.0.0.0"
cfg.Port = 8080

server, err := http.NewServer(cfg)
if err != nil {
    log.Fatal(err)
}

// 添加中间件
server.Use(http.RecoveryMiddleware())
server.Use(http.RateLimitMiddleware(100, 20)) // 100 req/s, burst 20
server.Use(http.CORSMiddleware(cfg.CORS))

// 注册路由
server.GET("/health", func(c *gin.Context) {
    c.JSON(200, gin.H{"status": "ok"})
})

server.POST("/api/users", handler.CreateUser)

// 启动服务器
go func() {
    if err := server.Start(); err != nil {
        log.Fatal(err)
    }
}()

// 优雅关闭
defer server.Shutdown(context.Background())
```

## 9. Discovery - 服务发现封装

### 文件结构
- `config.go` - 服务发现配置
- `types.go` - 类型定义和接口
- `client.go` - 统一客户端
- `consul.go` - Consul 驱动实现

### 使用示例

```go
import "aicode/pkg/discovery"

cfg := discovery.DefaultConfig()
cfg.Enabled = true
cfg.Driver = "consul"
cfg.Address = "127.0.0.1:8500"
cfg.ServiceName = "my-service"
cfg.ServicePort = 8080

client, err := discovery.NewClient(cfg)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// 注册服务
instance := discovery.ServiceInstance{
    ID:      "my-service-1",
    Name:    "my-service",
    Address: "192.168.1.100",
    Port:    8080,
    Tags:    []string{"v1", "production"},
}

err = client.Register(instance)
if err != nil {
    log.Fatal(err)
}

// 发现服务
instances, err := client.Discover("other-service")
if err != nil {
    log.Fatal(err)
}

// 健康检查
status, err := client.HealthCheck("my-service-1")

// 监听服务变化
ctx, cancel := context.WithCancel(context.Background())
stream, err := client.WatchService(ctx, "other-service")
for instances := range stream {
    fmt.Printf("Instances: %+v\n", instances)
}
```

## 10. Tracing - 链路追踪封装

### 文件结构
- `config.go` - 链路追踪配置
- `tracer.go` - Tracer 实现
- `propagator.go` - 上下文传播器

### 使用示例

```go
import "aicode/pkg/tracing"

cfg := tracing.DefaultConfig()
cfg.Enabled = true
cfg.Driver = "jaeger"
cfg.Endpoint = "http://jaeger:14268/api/traces"
cfg.ServiceName = "my-service"
cfg.SampleRate = 0.1

tracer, err := tracing.NewTracer(cfg)
if err != nil {
    log.Fatal(err)
}
defer tracer.Close()

// 创建 span
tr := tracer.Tracer("my-operation")
ctx, span := tr.Start(context.Background(), "operation-name")
defer span.End()

// 记录错误
if err := doSomething(); err != nil {
    span.RecordError(err)
}

// 设置属性
span.SetAttribute("custom.key", "value")
```

## 最佳实践

### 1. 配置管理
- 所有配置都通过 `config.go` 中的 `DefaultConfig()` 提供默认值
- 使用 YAML 配置文件时，未指定的字段自动使用默认值
- 敏感信息（密码、密钥）应通过环境变量注入

### 2. 资源管理
- 所有客户端都需要在应用退出时调用 `Close()` 方法
- 使用 `defer` 确保资源被正确释放
- 对于长连接（数据库、Redis），应用启动时初始化一次即可

### 3. 错误处理
- 所有方法都返回详细的错误信息（包含上下文）
- 使用 `fmt.Errorf("...: %w", err)` 包装错误
- 对于可重试的错误，应在业务层处理重试逻辑

### 4. 并发安全
- 所有客户端都是并发安全的
- 内部使用 `sync.RWMutex` 保护共享状态
- 消费者支持多 worker 并发处理

### 5. 健康检查
- 每个客户端都提供 `HealthCheck()` 方法
- 应用应定期检查基础设施健康状态
- 不健康的组件应及时告警或重启

## 依赖安装

运行以下命令安装所有依赖：

```bash
go get github.com/gin-gonic/gin
go get golang.org/x/time/rate
go get github.com/hashicorp/consul/api
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/exporters/jaeger
go get go.opentelemetry.io/otel/sdk
```

或者使用 `go mod tidy` 自动整理依赖。
