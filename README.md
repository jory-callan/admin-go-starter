# Go Web RBAC 后端管理框架

一个基于 Go + Echo + GORM + Redis 的生产级后端管理框架，内置完整的 RBAC 权限控制系统。

## ✨ 特性

- 🏗️ **清晰的分层架构** - Handler → Service → Repo → Model
- 🔐 **完善的 RBAC 权限控制** - 支持用户、角色、权限三级管理
- 🎯 **灵活的权限匹配** - 支持通配符 `*`，如 `system:user:*`、`system:*`、`*`
- 🚀 **泛型基础库** - 类似 MyBatis-Plus 的快速 CRUD 开发体验
- 🔑 **JWT 认证** - 安全的 Token 认证机制
- 📦 **组合根模式** - 手动依赖装配，简洁直观
- 📊 **统一响应格式** - 标准化的 API 响应结构
- 🗄️ **审计字段** - 完整的创建/更新/删除追踪
- 🆔 **UUIDv7** - 基于时间的分布式 ID 生成

## 📋 技术栈

- **Web 框架**: [Echo](https://github.com/labstack/echo)
- **ORM**: [GORM](https://github.com/go-gorm/gorm)
- **数据库**: MySQL
- **缓存**: Redis
- **配置管理**: [Viper](https://github.com/spf13/viper)
- **CLI 框架**: [Cobra](https://github.com/spf13/cobra)
- **JWT**: [golang-jwt/jwt](https://github.com/golang-jwt/jwt)
- **日志**: slog

## 🏗️ pkg 层基础设施

项目包含完整的基础设施封装，所有组件都遵循统一的设计模式，生产可用、易扩展、易维护。

### 已完成的组件

| 组件 | 功能 | 状态 |
|------|------|------|
| 📦 **Database** | GORM 数据库连接（支持 MySQL/PostgreSQL/SQLite） | ✅ |
| 🔴 **Redis** | Redis 客户端封装 | ✅ |
| 📝 **Logger** | 结构化日志（基于 slog，支持 JSON/T ext 格式） | ✅ |
| 🔐 **JWT** | JWT Token 生成和解析 | ✅ |
| 📨 **Kafka** | Kafka 生产者和消费者封装 | ✅ |
| 🐰 **RabbitMQ** | RabbitMQ 发布/订阅封装 | ✅ |
| 📧 **Email** | SMTP 邮件发送封装 | ✅ |
| 🌐 **HTTP** | HTTP 服务器（基于 Gin，含中间件） | ✅ |
| 🔍 **Discovery** | 服务发现（支持 Consul） | ✅ |
| 🔗 **Tracing** | 链路追踪（OpenTelemetry + Jaeger） | ✅ |

### pkg 层文档

- 📖 [详细使用指南](docs/pkg-usage.md) - 每个组件的完整说明和示例
- 📋 [快速参考](docs/pkg-cheatsheet.md) - API 速查和最佳实践
- 📚 [pkg README](pkg/README.md) - pkg 层总览和贡献指南

## 📁 项目结构

```
.
├── cmd                    # 应用入口
│   └── main.go           # 主程序入口
├── config                 # 配置文件
│   ├── config.go         # 配置结构定义
│   ├── loader.go         # 配置加载逻辑
│   └── config.yml        # 默认配置
├── configs                # 配置示例
│   └── config.yaml.example  # 完整配置示例
├── internal               # 内部业务逻辑
│   ├── handler           # HTTP 处理器层
│   ├── service           # 业务逻辑层
│   ├── repo              # 数据访问层
│   ├── model             # 数据模型层
│   ├── router            # 路由注册与依赖装配
│   └── middleware        # 中间件
├── pkg                    # 基础设施封装（生产级）
│   ├── database          # 数据库封装（GORM）
│   ├── redis             # Redis 封装
│   ├── logger            # 日志封装（slog）
│   ├── jwt               # JWT 工具
│   ├── kafka             # Kafka 封装
│   ├── rabbitmq          # RabbitMQ 封装
│   ├── email             # Email 封装
│   ├── http              # HTTP 服务器（Gin）
│   ├── discovery         # 服务发现（Consul）
│   └── tracing           # 链路追踪（OpenTelemetry）
├── docs                   # 文档
│   ├── pkg-usage.md      # pkg 层使用指南
│   ├── pkg-cheatsheet.md # 快速参考
│   └── COMPLETION_SUMMARY.md  # 完成总结
└── scripts                # 工具脚本
    ├── install-deps.sh   # Linux/macOS 依赖安装
    └── install-deps.ps1  # Windows 依赖安装
```

## 🚀 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 配置数据库

编辑 `config/config.yml`:

```yaml
server:
  port: "8080"

database:
  dsn: "root:password@tcp(localhost:3306)/aicode?charset=utf8mb4&parseTime=True&loc=Local"

redis:
  addr: "localhost:6379"
  password: ""
  db: 0

jwt:
  secret: "your-secret-key-change-in-production"
  expire: 24h
```

### 3. 创建数据库

```sql
CREATE DATABASE aicode CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. 运行项目

```bash
go run cmd/main.go
```

### 5. 初始化超级管理员

首次运行后，需要手动插入超级管理员：

```sql
-- 创建超级管理员角色
INSERT INTO roles (id, name, code, description, sort, status, created_at, updated_at)
VALUES ('00000000-0000-0000-0000-000000000001', '超级管理员', 'super_admin', '拥有所有权限', 1, 1, NOW(), NOW());

-- 创建用户（密码: admin123）
INSERT INTO users (id, username, password, nickname, status, created_at, updated_at)
VALUES ('00000000-0000-0000-0000-000000000001', 'admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '超级管理员', 1, NOW(), NOW());

-- 给用户分配超级管理员角色
INSERT INTO user_roles (user_id, role_id)
VALUES ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001');
```

## 📝 API 接口文档

### 认证接口

#### 登录
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

#### 获取当前用户信息
```http
GET /api/v1/user/info
Authorization: Bearer <token>
```

### 用户管理

#### 用户列表
```http
GET /api/v1/users?page=1&size=20
Authorization: Bearer <token>
```

#### 创建用户
```http
POST /api/v1/users
Authorization: Bearer <token>
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123",
  "nickname": "测试用户",
  "role_ids": ["role-id"]
}
```

### 角色管理

#### 角色列表
```http
GET /api/v1/roles?page=1&size=20
Authorization: Bearer <token>
```

### 权限管理

#### 权限树
```http
GET /api/v1/permissions/tree
Authorization: Bearer <token>
```

## 🔐 权限规则说明

### 权限码设计

权限码采用三级分层设计：`module:resource:action`

示例：
- `system:user:read` - 系统管理用户模块读取权限
- `system:user:write` - 系统管理用户模块写入权限
- `system:*` - 系统管理所有权限
- `*` - 超级管理员（所有权限）

### 通配符规则

- `*` - 匹配所有权限（超级管理员）
- `system:*` - 匹配 system 模块下所有权限
- `system:user:*` - 匹配 system:user 资源下所有权限

## 📄 License

MIT
