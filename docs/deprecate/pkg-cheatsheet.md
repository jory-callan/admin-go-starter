# pkg 层快速参考

## 组件总览

| 组件 | 包路径 | 主要功能 | 状态 |
|------|--------|---------|------|
| 📦 Database | `pkg/database` | GORM 数据库连接 | ✅ |
| 🔴 Redis | `pkg/redis` | Redis 客户端 | ✅ |
| 📝 Logger | `pkg/logger` | 结构化日志 | ✅ |
| 🔐 JWT | `pkg/jwt` | Token 认证 | ✅ |
| 📨 Kafka | `pkg/kafka` | 消息队列 | ✅ |
| 🐰 RabbitMQ | `pkg/rabbitmq` | 消息队列 | ✅ |
| 📧 Email | `pkg/email` | 邮件发送 | ✅ |
| 🌐 HTTP | `pkg/http` | HTTP 服务器 | ✅ |
| 🔍 Discovery | `pkg/discovery` | 服务发现 | ✅ |
| 🔗 Tracing | `pkg/tracing` | 链路追踪 | ✅ |

## 快速使用模板

### 1. 基本模式

```go
import "aicode/pkg/<component>"

cfg := <component>.DefaultConfig()
// 修改配置...

client, err := <component>.NewClient(cfg)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// 使用 client...
```

### 2. 配置示例（YAML）

```yaml
# 所有组件的配置示例
database:
  driver: "mysql"
  dsn: "root:pass@tcp(localhost:3306)/db"

redis:
  addr: "localhost:6379"

email:
  smtp_host: "smtp.example.com"
  smtp_port: 587
  username: "user@example.com"
  password: "secret"

kafka:
  brokers: ["localhost:9092"]
  topic: "my-topic"

rabbitmq:
  url: "amqp://guest:guest@localhost:5672/"
```

## API 速查

### Database (GORM)

```go
db, err := database.Open(cfg)
db.AutoMigrate(&User{})
db.Create(&user)
db.Where("id = ?", id).First(&user)
```

### Redis

```go
client, err := redis.Open(cfg)
client.Set(ctx, "key", "value", 0).Err()
client.Get(ctx, "key").Result()
```

### Logger

```go
log := logger.New(cfg, "app-name")
log.Info("message", "key", value)
log.Error("error", "err", err)
subLog := logger.C("component")
```

### JWT

```go
jwt.Init(cfg)
token, _ := jwt.GenerateToken(id, name, roles, perms)
claims, _ := jwt.ParseToken(token)
```

### Kafka Producer

```go
client, _ := kafka.NewClient(cfg)
producer := kafka.NewProducer(client)
producer.SendJSON(ctx, "topic", "key", data)
```

### Kafka Consumer

```go
consumer := kafka.NewConsumer(client, kafka.ConsumerConfig{
    Topic:   "topic",
    GroupID: "group",
    Handler: handlerFunc,
})
consumer.Start()
defer consumer.Stop()
```

### RabbitMQ Producer

```go
client, _ := rabbitmq.NewClient(cfg)
producer := rabbitmq.NewProducer(client, exchange)
producer.SendJSON(ctx, "routing.key", data)
```

### RabbitMQ Consumer

```go
consumer := rabbitmq.NewConsumer(client, rabbitmq.ConsumerConfig{
    Queue:   "queue",
    Handler: handlerFunc,
})
consumer.Start()
defer consumer.Stop()
```

### Email

```go
client, _ := email.NewClient(cfg)
client.SendSimple([]string{"to@example.com"}, "subj", "body")
client.SendHTML([]string{"to@example.com"}, "subj", "<html>...</html>")
```

### HTTP Server

```go
server, _ := http.NewServer(cfg)
server.Use(middleware...)
server.GET("/path", handler)
server.POST("/path", handler)
server.Start()
defer server.Shutdown(ctx)
```

### Discovery

```go
client, _ := discovery.NewClient(cfg)
client.Register(instance)
instances, _ := client.Discover("service-name")
```

### Tracing

```go
tracer, _ := tracing.NewTracer(cfg)
defer tracer.Close()
tr := tracer.Tracer("operation")
ctx, span := tr.Start(ctx, "span-name")
defer span.End()
```

## 常用配置默认值

### Database
```go
Driver:          "sqlite"
MaxOpenConns:    50
MaxIdleConns:    10
ConnMaxLifetime: 1800
```

### Redis
```go
Addr:         "127.0.0.1:6379"
PoolSize:     100
MinIdleConns: 10
DialTimeout:  5s
```

### Logger
```go
Level:      "info"
Format:     "json"
Output:     "stdout"
MaxSize:    100MB
MaxBackups: 5
MaxAge:     7 days
```

### JWT
```go
Expires: "24h"
Issuer:  "aicode"
```

### HTTP
```go
Host:           "0.0.0.0"
Port:           8080
ReadTimeout:    10s
WriteTimeout:   10s
IdleTimeout:    60s
MaxHeaderBytes: 1MB
```

### Email
```go
SMTPHost:  "smtp.example.com"
SMTPPort:  587
UseTLS:    true
Timeout:   10s
```

## 错误处理最佳实践

```go
// ❌ 不好的做法
if err != nil {
    return err
}

// ✅ 好的做法（添加上下文）
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// ✅ 更好的做法（包含关键参数）
if err != nil {
    return fmt.Errorf("connect database (driver=%s): %w", cfg.Driver, err)
}
```

## 资源管理最佳实践

```go
// 立即关闭资源
client, err := pkg.NewClient(cfg)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// 多个资源时逆序关闭
db, _ := database.Open(dbCfg)
defer db.Close()

redis, _ := redis.Open(redisCfg)
defer redis.Close()

kafka, _ := kafka.NewClient(kafkaCfg)
defer kafka.Close()
```

## 并发安全提示

✅ 所有客户端都是并发安全的  
✅ 可以在多个 goroutine 中共享同一个客户端实例  
✅ 内部使用 sync.RWMutex 保护共享状态  

```go
// 安全的用法
client := pkg.NewClient(cfg)
go useClient(client)
go useClient(client)
go useClient(client)
```

## 健康检查

```go
// 定期检查健康状态
ticker := time.NewTicker(30 * time.Second)
go func() {
    for range ticker.C {
        if err := client.HealthCheck(); err != nil {
            log.Error("unhealthy", "error", err)
            // 触发告警或重启
        }
    }
}()
```

## 依赖安装命令

```bash
# 一键安装所有依赖
go get github.com/gin-gonic/gin
go get golang.org/x/time/rate
go get github.com/hashicorp/consul/api
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/exporters/jaeger
go get go.opentelemetry.io/otel/sdk
go mod tidy
```

## 调试技巧

### 1. 启用 Debug 模式

```go
// HTTP
cfg.EnableDebug = true

// Logger
cfg.Level = "debug"

// Database
cfg.LogLevel = "info"  // 或 "debug"
```

### 2. 查看详细日志

```yaml
log:
  level: "debug"
  format: "text"  # text 格式更易读
  output: "stdout"
```

### 3. 监控连接池

```go
// Redis
stats := client.PoolStats()
fmt.Printf("Active: %d, Idle: %d\n", stats.ActiveConns, stats.IdleConns)

// Database
stats := db.Stats()
fmt.Printf("Open: %d, InUse: %d, Idle: %d\n", stats.OpenConnections, stats.InUse, stats.Idle)
```

## 性能优化建议

1. **连接池大小**: 根据负载调整 PoolSize
2. **批量操作**: Kafka/RabbitMQ 使用批量发送
3. **超时设置**: 合理设置 timeout 避免长时间阻塞
4. **重试机制**: 网络请求添加指数退避重试
5. **缓存策略**: Redis 使用 pipeline 减少 RTT

---

**最后更新**: 2026-03-27  
**维护者**: Your Team
