# pkg 层基础设施封装

本目录包含项目的所有基础设施组件封装，每个组件都遵循统一的设计模式。

## 组件列表

| 组件 | 状态 | 说明 |
|------|------|------|
| database | ✅ 完成 | GORM 数据库连接管理（支持 MySQL/PostgreSQL/SQLite） |
| redis | ✅ 完成 | Redis 客户端封装 |
| logger | ✅ 完成 | 结构化日志（基于 slog） |
| jwt | ✅ 完成 | JWT Token 生成和解析 |
| kafka | ✅ 完成 | Kafka 生产者和消费者封装 |
| rabbitmq | ✅ 完成 | RabbitMQ 发布/订阅封装 |
| email | ✅ 完成 | SMTP 邮件发送封装 |
| http | ✅ 完成 | HTTP 服务器（基于 Gin） |
| discovery | ✅ 完成 | 服务发现（支持 Consul） |
| tracing | ✅ 完成 | 链路追踪（OpenTelemetry） |

## 设计原则

### 1. 统一配置模式

每个组件都有 `config.go` 文件提供配置结构体和默认值：

```go
// config.go
type Config struct {
    // 配置字段...
}

func DefaultConfig() Config {
    return Config{
        // 默认值...
    }
}
```

### 2. 分离关注点

每个组件按功能拆分为多个文件：

- `config.go` - 配置定义
- `client.go` / `service.go` - 核心客户端逻辑
- `producer.go` / `consumer.go` - 特定功能封装（如消息队列）

### 3. 资源管理

所有客户端都提供 `Close()` 方法用于资源清理：

```go
client, err := pkg.NewClient(cfg)
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

### 4. 健康检查

每个客户端都提供 `HealthCheck()` 方法：

```go
if err := client.HealthCheck(); err != nil {
    log.Error("unhealthy", "error", err)
}
```

### 5. 错误处理

所有方法返回详细的错误信息：

```go
return fmt.Errorf("operation failed: %w", err)
```

## 快速开始

### 1. 安装依赖

**Linux/macOS:**
```bash
./scripts/install-deps.sh
```

**Windows (PowerShell):**
```powershell
.\scripts\install-deps.ps1
```

**手动安装:**
```bash
go get github.com/gin-gonic/gin
go get golang.org/x/time/rate
go get github.com/hashicorp/consul/api
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/exporters/jaeger
go get go.opentelemetry.io/otel/sdk
go mod tidy
```

### 2. 配置

在 `config/config.yaml` 中添加相应配置：

```yaml
# 示例：Email 配置
email:
  smtp_host: "smtp.example.com"
  smtp_port: 587
  username: "noreply@example.com"
  password: "your-password"
  from_name: "My App"
  use_tls: true
  timeout: 10
```

### 3. 使用

```go
import "aicode/pkg/email"

cfg := email.DefaultConfig()
// 从配置文件加载实际值...

client, err := email.NewClient(cfg)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

err = client.SendSimple([]string{"user@example.com"}, "Subject", "Body")
```

## 详细文档

查看 [pkg-usage.md](../docs/pkg-usage.md) 获取每个组件的详细使用示例。

## 添加新组件

如需添加新的基础设施组件，请遵循以下模板：

### 1. 创建目录结构

```
pkg/newcomponent/
├── config.go      # 配置定义
└── client.go      # 客户端实现
```

### 2. 实现 config.go

```go
package newcomponent

// Config 配置结构
type Config struct {
    // 字段定义...
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
    return Config{
        // 默认值...
    }
}
```

### 3. 实现 client.go

```go
package newcomponent

// Client 客户端
type Client struct {
    cfg Config
}

// NewClient 创建客户端
func NewClient(cfg Config) (*Client, error) {
    // 初始化逻辑...
}

// Close 关闭资源
func (c *Client) Close() error {
    // 清理逻辑...
}

// HealthCheck 健康检查
func (c *Client) HealthCheck() error {
    // 检查逻辑...
}
```

### 4. 更新全局配置

编辑 `config/config.go`:

```go
import pkgnewcomponent "aicode/pkg/newcomponent"

type AppConfig struct {
    // ...
    NewComponent pkgnewcomponent.Config `mapstructure:"new_component" yaml:"new_component"`
}
```

### 5. 更新配置加载器

编辑 `config/loader.go`:

```go
// NewComponent
v.SetDefault("new_component.some_field", cfg.NewComponent.SomeField)
```

## 测试建议

为每个组件编写单元测试：

```go
// client_test.go
func TestClient_Send(t *testing.T) {
    cfg := DefaultConfig()
    // 设置测试配置...
    
    client, err := NewClient(cfg)
    if err != nil {
        t.Fatal(err)
    }
    defer client.Close()
    
    // 测试逻辑...
}
```

## 注意事项

1. **并发安全**: 所有客户端都是并发安全的，内部使用 `sync.RWMutex` 保护共享状态
2. **超时控制**: 所有外部调用都应设置合理的超时时间
3. **重试机制**: 对于网络请求，建议在业务层实现重试逻辑
4. **日志记录**: 关键操作应记录日志，便于问题排查
5. **优雅关闭**: 应用退出时应调用所有客户端的 `Close()` 方法

## 故障排查

### 常见问题

**Q: 依赖安装失败**
A: 检查网络连接，尝试使用国内镜像：
```bash
export GOPROXY=https://goproxy.cn,direct
go mod tidy
```

**Q: 连接被拒绝**
A: 检查服务是否启动、端口是否正确、防火墙设置

**Q: 超时错误**
A: 增加配置中的超时时间，检查网络延迟

## 贡献指南

欢迎提交 Issue 和 Pull Request 来改进这些封装！

---

**维护者**: Your Team  
**最后更新**: 2026-03-27
