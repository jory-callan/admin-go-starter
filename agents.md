# Go Web 脚手架项目开发规范

## 1. 技术栈

### 核心框架
- **配置管理**: `github.com/spf13/viper`
- **CLI 框架**: `github.com/spf13/cobra`
- **Web 框架**: `github.com/labstack/echo/v4`
- **ORM**: `gorm.io/gorm` (配合对应数据库驱动)
- **Redis**: `github.com/redis/go-redis/v9`
- **日志**: `slog` 需要支持默认输出控制台，支持同时输出文件和控制台持久化，支持日志轮转
- **UUID**: `github.com/gofrs/uuid/v5` (用于 UUIDv7 生成)

### 必要依赖
```go
// 数据库驱动根据实际选择
gorm.io/driver/mysql      // MySQL
gorm.io/driver/postgres   // PostgreSQL
gorm.io/driver/sqlite     // SQLite
```

---

## 2. 项目目录结构

```
project-root/
├── cmd/                      # 应用入口
│   └── main.go              # 主程序入口
├── config/                   # 配置文件目录
│   └── config.yml           # 默认配置文件
├── internal/                 # 内部业务逻辑（不对外暴露）
│   ├── handler/             # HTTP 处理器层
│   │   ├── user_handler.go
│   │   └── common.go        # 统一响应结构
│   ├── service/             # 业务逻辑层
│   │   └── user_service.go
│   ├── repo/                # 数据访问层
│   │   ├── user_repo.go
│   │   └── base_repo.go     # 泛型基础 Repo
│   ├── model/               # 数据模型层
│   │   ├── user.go
│   │   └── base.go          # 基础模型（审计字段）
│   ├── router/              # 路由注册与依赖装配
│   │   └── router.go
│   └── middleware/          # 中间件
│       └── auth.go
├── pkg/                      # 基础设施封装（可复用）
│   ├── app/                 # App 结构体定义
│   │   └── app.go
│   ├── config/              # 配置加载
│   │   └── config.go
│   ├── db/                  # 数据库初始化
│   │   └── db.go
│   ├── redis/               # Redis 初始化
│   │   └── redis.go
│   └── logger/              # 日志初始化
│       └── logger.go
├── go.mod
└── go.sum
```

---

## 3. 核心设计原则

### 3.1 架构模式
- **分层架构**: handler → service → repo → model
- **结构体优先**: 每层直接采用结构体，**禁止使用接口**（除非有明确的多实现需求）
- **功能模块化**: 按业务功能组织代码，每个模块包含完整的 handler/service/repo/model

### 3.2 依赖管理（组合根模式）
- **禁止使用依赖注入框架**（wire、fx 等）
- 在 `internal/router/router.go` 中手动装配所有依赖
- `main.go` 保持简洁，只负责初始化和启动

### 3.3 基础设施统一结构体

```go
// pkg/app/app.go
package app

import (
    "gorm.io/gorm"
    "github.com/redis/go-redis/v9"
    "github.com/spf13/viper"
    "go.uber.org/zap"
)

// App 只包含与业务逻辑相关的资源，不包含任何 HTTP 框架依赖
type App struct {
    DB    *gorm.DB
    Redis *redis.Client
    Conf  *viper.Viper
}
```

---

## 4. 配置文件规范

### 4.1 配置文件查找优先级
1. 通过 `-c` 参数指定的配置文件路径
2. `./config.yml`
3. `./config/config.yml`
4. `./conf/config.yml`

---

## 5. 数据库设计规范

### 5.1 基础模型（审计字段）

```go
// internal/model/base.go
package model

import (
    "time"
    "gorm.io/gorm"
)

// 所有业务表的公共审计字段 必须完整写入，禁止结构体引入
    ID        string         `gorm:"primaryKey;type:varchar(36);comment:主键ID" json:"id"`
    CreatedAt time.Time      `gorm:"comment:创建时间" json:"created_at"`
    CreatedBy string         `gorm:"type:varchar(36);comment:创建人ID" json:"created_by"`
    UpdatedAt time.Time      `gorm:"comment:更新时间" json:"updated_at"`
    UpdatedBy string         `gorm:"type:varchar(36);comment:更新人ID" json:"updated_by"`
    DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"deleted_at"`
    DeletedBy string         `gorm:"type:varchar(36);comment:删除人ID" json:"deleted_by"`

// 创建前自动生成 UUIDv7 ID
func (m *xxxModel) BeforeCreate(tx *gorm.DB) error {
    if m.ID == "" {
        m.ID = generateUUIDv7()
    }
    return nil
}
```

### 5.2 UUIDv7 生成

```go
// pkg/uuid/uuid.go
// 全局复用同一个生成器，避免每次 new。
var defaultGen = uuid.NewGen()

// UUID 返回标准 UUIDv7 字符串。
func UUIDv7() string {
	return uuid.Must(defaultGen.NewV7()).String()
}

// ShortUUIDv7 返回去掉连字符的 UUIDv7。
func ShortUUIDv7() string {
	return strings.ReplaceAll(UUIDv7(), "-", "")
}

```

### 5.3 业务模型示例

```go
// internal/model/user.go
package model

// User 用户表
type User struct {
    // 每个表必备审计字段，禁止结构体引入

    // 其他业务字段
    Username string `gorm:"type:varchar(50);uniqueIndex;comment:用户名" json:"username"`
    Email    string `gorm:"type:varchar(100);uniqueIndex;comment:邮箱" json:"email"`
    Password string `gorm:"type:varchar(255);comment:密码" json:"-"`
    Nickname string `gorm:"type:varchar(50);comment:昵称" json:"nickname"`
    Status   int    `gorm:"type:tinyint;default:1;comment:状态(1:正常 2:禁用)" json:"status"`
}

func (User) TableName() string {
    return "users"
}
```

**关键要求**：
- ✅ 所有表的每个字段**必须**有 `comment` 标签
- ✅ 所有业务表**必须**嵌入 `BaseModel`
- ✅ ID 字段使用字符串类型的 UUIDv7
- ✅ 敏感字段（如 password）JSON 标签设为 `-`

---

## 6. 泛型设计规范

### 6.1 分页查询结构

```go
// internal/model/page.go
package model

// PageQuery 分页查询参数
type PageQuery struct {
    Page      int    `json:"page" query:"page"`           // 页码，从 1 开始
    Size      int    `json:"size" query:"size"`           // 每页数量
    NeedCount bool   `json:"need_count" query:"need_count"` // 是否需要总数
    Order     string `json:"order" query:"order"`         // 排序字段，如 "created_at desc"
}

// PageResult 分页返回结果
type PageResult[T any] struct {
    Items   []T   `json:"items"`    // 数据列表
    Total   int64 `json:"total"`    // 总记录数
    Page    int   `json:"page"`     // 当前页码
    Size    int   `json:"size"`     // 每页数量
    HasMore bool  `json:"hasMore"`  // 是否有下一页
}
```

### 6.2 泛型基础 Repository

```go
// internal/repo/base_repo.go
package repo

import (
    "context"
    "gorm.io/gorm"
    "your-project/internal/model"
)

// BaseRepo 泛型基础 Repository，提供通用 CRUD 操作
type BaseRepo[T any] struct {
    DB *gorm.DB
}

// NewBaseRepo 创建基础 Repo
func NewBaseRepo[T any](db *gorm.DB) *BaseRepo[T] {
    return &BaseRepo[T]{DB: db}
}

// Create 创建记录
func (r *BaseRepo[T]) Create(ctx context.Context, entity *T) error {
    return r.DB.WithContext(ctx).Create(entity).Error
}

// Update 更新记录
func (r *BaseRepo[T]) Update(ctx context.Context, entity *T) error {
    return r.DB.WithContext(ctx).Save(entity).Error
}

// Delete 软删除记录
func (r *BaseRepo[T]) Delete(ctx context.Context, id string) error {
    var entity T
    return r.DB.WithContext(ctx).Delete(&entity, "id = ?", id).Error
}

// GetByID 根据 ID 查询
func (r *BaseRepo[T]) GetByID(ctx context.Context, id string) (*T, error) {
    var entity T
    err := r.DB.WithContext(ctx).First(&entity, "id = ?", id).Error
    if err != nil {
        return nil, err
    }
    return &entity, nil
}

// List 列表查询
func (r *BaseRepo[T]) List(ctx context.Context, query *model.PageQuery) (*model.PageResult[T], error) {
    var items []T
    var total int64

    db := r.DB.WithContext(ctx).Model(new(T))

    // 排序
    if query.Order != "" {
        db = db.Order(query.Order)
    }

    // 分页
    offset := (query.Page - 1) * query.Size
    if err := db.Offset(offset).Limit(query.Size).Find(&items).Error; err != nil {
        return nil, err
    }

    // 是否统计总数
    if query.NeedCount {
        if err := db.Count(&total).Error; err != nil {
            return nil, err
        }
    }

    // 计算是否有下一页
    hasMore := false
    if query.NeedCount {
        hasMore = int64(query.Page*query.Size) < total
    }

    return &model.PageResult[T]{
        Items:   items,
        Total:   total,
        Page:    query.Page,
        Size:    query.Size,
        HasMore: hasMore,
    }, nil
}

// Where 条件查询（类似 mybatis-plus）
func (r *BaseRepo[T]) Where(conditions ...interface{}) *gorm.DB {
    return r.DB.Where(conditions[0], conditions[1:]...)
}
```

### 6.3 业务 Repository 示例

```go
// internal/repo/user_repo.go
package repo

import (
    "context"
    "gorm.io/gorm"
    "your-project/internal/model"
)

// UserRepo 用户数据访问层
type UserRepo struct {
    *BaseRepo[model.User]
}

// NewUserRepo 创建用户 Repo
func NewUserRepo(db *gorm.DB) *UserRepo {
    return &UserRepo{
        BaseRepo: NewBaseRepo[model.User](db),
    }
}

// GetByUsername 根据用户名查询
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
    var user model.User
    err := r.DB.WithContext(ctx).Where("username = ?", username).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}
```

---

## 7. 统一响应结构

### 7.1 HTTP 响应结构

```go
// internal/handler/common.go
package handler

import (
    "github.com/labstack/echo/v4"
    "net/http"
)

// Response 统一 HTTP 响应结构
type Response struct {
    Code int         `json:"code"`       // 业务状态码
    Msg  string      `json:"msg"`        // 提示信息
    Data interface{} `json:"data"`       // 业务数据
}

// 成功响应
func Success(c echo.Context, data interface{}) error {
    return c.JSON(http.StatusOK, &Response{
        Code: 0,
        Msg:  "success",
        Data: data,
    })
}

// 失败响应
func Fail(c echo.Context, code int, msg string) error {
    return c.JSON(http.StatusOK, &Response{
        Code: code,
        Msg:  msg,
        Data: nil,
    })
}

// 分页响应
func SuccessPage(c echo.Context, items interface{}, total int64, page, size int, hasMore bool) error {
    return Success(c, map[string]interface{}{
        "items":   items,
        "total":   total,
        "page":    page,
        "size":    size,
        "hasMore": hasMore,
    })
}
```

### 7.2 业务状态码规范

```go
// internal/handler/code.go
package handler

const (
    CodeSuccess       = 0      // 成功
    CodeBadRequest    = 400    // 请求参数错误
    CodeUnauthorized  = 401    // 未授权
    CodeForbidden     = 403    // 禁止访问
    CodeNotFound      = 404    // 资源不存在
    CodeInternalError = 500    // 内部错误
)
```

---

## 8. 路由与依赖装配

### 8.1 路由注册示例

```go
// internal/router/router.go
package router

import (
    "github.com/labstack/echo/v4"
    "your-project/internal/handler"
    "your-project/internal/middleware"
    "your-project/internal/repo"
    "your-project/internal/service"
    "your-project/pkg/app"
)

// RegisterRoutes 注册所有路由并装配依赖
func RegisterRoutes(e *echo.Echo, app *app.App) {
    // ============ 依赖装配（组合根） ============
    
    // Repo 层
    userRepo := repo.NewUserRepo(app.DB)
    
    // Service 层
    userService := service.NewUserService(userRepo, app.Log)
    
    // Handler 层
    userHandler := handler.NewUserHandler(userService)
    
    // ============ 路由注册 ============
    
    // 健康检查
    e.GET("/health", func(c echo.Context) error {
        return c.String(200, "OK")
    })
    
    // API v1
    v1 := e.Group("/api/v1")
    
    // 用户模块
    userGroup := v1.Group("/users")
    {
        userGroup.GET("", userHandler.List)       // 用户列表
        userGroup.GET("/:id", userHandler.Get)    // 用户详情
        userGroup.POST("", userHandler.Create)    // 创建用户
        userGroup.PUT("/:id", userHandler.Update) // 更新用户
        userGroup.DELETE("/:id", userHandler.Delete) // 删除用户
    }
    
    // 中间件
    e.Use(middleware.Logger(app.Log))
    e.Use(middleware.Recover())
}
```

---

## 9. 应用启动流程

### 9.1 主程序入口

```go
// cmd/main.go
package main

import (
    "fmt"
    "github.com/labstack/echo/v4"
    "github.com/spf13/cobra"
    "your-project/internal/router"
    "your-project/pkg/app"
    "your-project/pkg/config"
    "your-project/pkg/db"
    "your-project/pkg/logger"
    "your-project/pkg/redis"
)

var configFile string

func main() {
    var rootCmd = &cobra.Command{
        Use:   "myapp",
        Short: "My Go Web Application",
        Run:   run,
    }
    
    // -c 参数指定配置文件
    rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "config file path")
    
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
    }
}

func run(cmd *cobra.Command, args []string) {
    // 1. 加载配置
    conf := config.Load(configFile)
    
    // 2. 初始化日志
    log := logger.New(conf)
    
    // 3. 初始化数据库
    database := db.New(conf, log)
    
    // 4. 初始化 Redis
    rdb := redis.New(conf, log)
    
    // 5. 创建 App 结构体（基础设施）
    application := &app.App{
        DB:    database,
        Redis: rdb,
        Conf:  conf,
        Log:   log,
    }
    
    // 6. 创建 Echo 实例
    e := echo.New()
    
    // 7. 注册路由和依赖装配
    router.RegisterRoutes(e, application)
    
    // 8. 启动服务
    port := conf.GetString("server.port")
    log.Info("Server starting on :" + port)
    if err := e.Start(":" + port); err != nil {
        log.Fatal("Failed to start server", "error", err)
    }
}
```

---

## 10. Handler/Service/Repo 层示例

### 10.1 Handler 层

```go
// internal/handler/user_handler.go
package handler

import (
    "github.com/labstack/echo/v4"
    "net/http"
    "your-project/internal/model"
    "your-project/internal/service"
)

type UserHandler struct {
    svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
    return &UserHandler{svc: svc}
}

// List 用户列表
func (h *UserHandler) List(c echo.Context) error {
    var query model.PageQuery
    if err := c.Bind(&query); err != nil {
        return Fail(c, CodeBadRequest, "参数错误")
    }
    
    // 默认值
    if query.Page == 0 {
        query.Page = 1
    }
    if query.Size == 0 {
        query.Size = 20
    }
    
    result, err := h.svc.List(c.Request().Context(), &query)
    if err != nil {
        return Fail(c, CodeInternalError, err.Error())
    }
    
    return SuccessPage(c, result.Items, result.Total, result.Page, result.Size, result.HasMore)
}

// Get 用户详情
func (h *UserHandler) Get(c echo.Context) error {
    id := c.Param("id")
    user, err := h.svc.GetByID(c.Request().Context(), id)
    if err != nil {
        return Fail(c, CodeNotFound, "用户不存在")
    }
    return Success(c, user)
}

// Create 创建用户
func (h *UserHandler) Create(c echo.Context) error {
    var req struct {
        Username string `json:"username"`
        Email    string `json:"email"`
        Password string `json:"password"`
        Nickname string `json:"nickname"`
    }
    
    if err := c.Bind(&req); err != nil {
        return Fail(c, CodeBadRequest, "参数错误")
    }
    
    user := &model.User{
        Username: req.Username,
        Email:    req.Email,
        Password: req.Password,
        Nickname: req.Nickname,
    }
    
    if err := h.svc.Create(c.Request().Context(), user); err != nil {
        return Fail(c, CodeInternalError, err.Error())
    }
    
    return Success(c, user)
}

// Update 更新用户
func (h *UserHandler) Update(c echo.Context) error {
    id := c.Param("id")
    
    var req struct {
        Nickname string `json:"nickname"`
        Status   int    `json:"status"`
    }
    
    if err := c.Bind(&req); err != nil {
        return Fail(c, CodeBadRequest, "参数错误")
    }
    
    user, err := h.svc.GetByID(c.Request().Context(), id)
    if err != nil {
        return Fail(c, CodeNotFound, "用户不存在")
    }
    
    user.Nickname = req.Nickname
    user.Status = req.Status
    
    if err := h.svc.Update(c.Request().Context(), user); err != nil {
        return Fail(c, CodeInternalError, err.Error())
    }
    
    return Success(c, user)
}

// Delete 删除用户
func (h *UserHandler) Delete(c echo.Context) error {
    id := c.Param("id")
    
    if err := h.svc.Delete(c.Request().Context(), id); err != nil {
        return Fail(c, CodeInternalError, err.Error())
    }
    
    return Success(c, nil)
}
```

### 10.2 Service 层

```go
// internal/service/user_service.go
package service

import (
    "context"
    "go.uber.org/zap"
    "your-project/internal/model"
    "your-project/internal/repo"
)

type UserService struct {
    repo *repo.UserRepo
    log  *zap.Logger
}

func NewUserService(repo *repo.UserRepo, log *zap.Logger) *UserService {
    return &UserService{
        repo: repo,
        log:  log,
    }
}

func (s *UserService) Create(ctx context.Context, user *model.User) error {
    // TODO: 密码加密、参数校验等业务逻辑
    return s.repo.Create(ctx, user)
}

func (s *UserService) Update(ctx context.Context, user *model.User) error {
    return s.repo.Update(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, id string) error {
    return s.repo.Delete(ctx, id)
}

func (s *UserService) GetByID(ctx context.Context, id string) (*model.User, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *UserService) List(ctx context.Context, query *model.PageQuery) (*model.PageResult[model.User], error) {
    return s.repo.List(ctx, query)
}
```

---

## 11. 代码风格与最佳实践

### 11.1 简洁优先
- ✅ 用最少代码实现最多功能
- ✅ 避免过度设计，保持简单直观
- ✅ 单人开发不需要复杂抽象

### 11.2 错误处理
- ✅ 在 handler 层统一处理错误响应
- ✅ 使用自定义错误码区分业务错误
- ✅ 日志记录关键错误信息

### 11.3 上下文传递
- ✅ 所有方法第一参数为 `context.Context`
- ✅ 使用 `ctx` 统一命名

### 11.4 结构体命名
- ✅ Handler: `XxxHandler`
- ✅ Service: `XxxService`
- ✅ Repo: `XxxRepo`
- ✅ Model: 直接使用领域名称（如 `User`）

---

## 12. 快速开发命令

### 12.1 启动应用
```bash
# 使用默认配置
go run cmd/main.go

# 指定配置文件
go run cmd/main.go -c /path/to/config.yml

# 编译后运行
go build -o myapp cmd/main.go
./myapp -c config/config.yml
```

### 12.2 数据库迁移
```go
// 在 db.New() 中添加自动迁移
func New(conf *viper.Viper, log *zap.Logger) *gorm.DB {
    // ... 连接数据库
    
    // 自动迁移
    db.AutoMigrate(
        &model.User{},
        // 其他模型...
    )
    
    return db
}
```

---

## 13. 扩展指南

### 13.1 新增业务模块
1. 在 `internal/model/` 创建模型
2. 在 `internal/repo/` 创建 Repo（继承 `BaseRepo`）
3. 在 `internal/service/` 创建 Service
4. 在 `internal/handler/` 创建 Handler
5. 在 `internal/router/router.go` 装配依赖并注册路由

### 13.2 添加中间件
```go
// internal/middleware/auth.go
package middleware

import (
    "github.com/labstack/echo/v4"
    "net/http"
)

func Auth() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // TODO: 实现 JWT 验证等逻辑
            token := c.Request().Header.Get("Authorization")
            if token == "" {
                return c.JSON(http.StatusUnauthorized, map[string]string{
                    "code": "401",
                    "msg":  "unauthorized",
                })
            }
            return next(c)
        }
    }
}
```

---

## 14. 总结

本规范旨在实现：
- ✅ **简单直观**: 代码结构清晰，易于理解和维护
- ✅ **快速开发**: 泛型 CRUD + 类似 mybatis-plus 的便捷操作
- ✅ **统一规范**: 响应结构、错误处理、命名规则
- ✅ **易于维护**: 清晰的分层架构，手动管理依赖
- ✅ **适合单人**: 不引入复杂框架，保持代码可控
