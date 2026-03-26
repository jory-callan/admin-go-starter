---
name: Go Web 脚手架 Agents.md 设计
overview: 设计一个 agents.md 规范文档，用于指导 AI 助手根据用户需求生成符合规范的 Go Web 脚手架项目代码。文档包含技术栈、架构设计、代码规范、目录结构等核心内容。
todos:
  - id: analyze-requirements
    content: 分析用户需求，梳理脚手架规范要点
    status: completed
  - id: create-agents-md
    content: 使用 [skill:skill-creator] 创建完整的 agents.md 规范文档
    status: completed
    dependencies:
      - analyze-requirements
---

## 产品概述

创建一份 agents.md 规范文档，用于指导 AI 自动生成符合用户技术栈和架构要求的 Go 脚手架项目

## 核心特性

- 定义技术栈规范：viper + cobra + echo + gorm + goredis(v9)
- 定义启动参数规范：-c 指定配置文件，默认查找顺序 ./config.yml > ./config/config.yml > ./conf/config.yml
- 定义架构规范：整洁架构 + 按功能模块拆分，handler -> service -> repo -> model 四层结构
- 定义依赖管理规范：组合根模式，统一 App 结构体，在 router 中装配依赖
- 定义数据库规范：审计字段、UUIDv7、字段注释
- 定义泛型 CRUD 规范：类似 mybatis-plus 的快速开发结构
- 定义统一响应结构：HTTP 返回和分页返回格式
- 定义目录结构规范：cmd、config、internal、pkg

## 技术栈选择

- **配置管理**: viper - Go 生态最流行的配置库
- **CLI 框架**: cobra - 强大的命令行应用框架
- **Web 框架**: echo - 高性能、简洁的 HTTP 框架
- **ORM**: gorm - Go 最成熟的 ORM 库
- **Redis 客户端**: goredis v9 - 官方推荐的 Redis 客户端

## 实现方案

### 文档结构设计

agents.md 将作为 AI 生成脚手架项目的"宪法"，包含以下核心章节：

1. **项目概述** - 技术栈、设计理念
2. **目录结构规范** - 完整的文件夹组织和职责说明
3. **基础设施层规范** - pkg 目录下各组件的初始化模式
4. **业务层规范** - internal 目录下四层架构的实现规则
5. **泛型 CRUD 规范** - BaseRepo、BaseService、PageQuery 等泛型结构
6. **统一响应规范** - Response、PageResult 等结构定义
7. **数据库规范** - 审计字段、UUIDv7、字段注释规则
8. **配置规范** - 配置文件结构和读取规则
9. **启动流程规范** - main.go 和 cmd 目录的组织方式
10. **代码示例** - 完整的用户模块示例代码

### 核心设计要点

#### 组合根模式

```
// pkg/app/app.go
type App struct {
    DB    *gorm.DB
    Redis *redis.Client
    Conf  *config.Config
    Log   *logger.Logger
}
```

#### 泛型 CRUD 结构

```
// internal/base/repo.go
type BaseRepo[T any] struct {
    db *gorm.DB
}

func (r *BaseRepo[T]) Create(entity *T) error
func (r *BaseRepo[T]) GetByID(id string) (*T, error)
func (r *BaseRepo[T]) List(query *PageQuery) (*PageResult[T], error)
// ... 更多通用方法
```

#### 统一响应结构

```
// pkg/response/response.go
type Response struct {
    Code int         `json:"code"`
    Msg  string      `json:"msg"`
    Data interface{} `json:"data"`
}

type PageResult[T any] struct {
    Items   []T   `json:"items"`
    Total   int64 `json:"total"`
    Page    int   `json:"page"`
    Size    int   `json:"size"`
    HasMore bool  `json:"hasMore"`
}
```

## 实施要点

- 文档需足够详细，让 AI 能够准确理解每个规范
- 提供完整代码示例，确保生成的代码风格一致
- 强调"最少代码实现最多功能"的原则
- 避免过度设计，适合单人开发场景

## Agent Extensions

### Skill

- **skill-creator**
- Purpose: 指导创建有效的 agents.md 规范文档
- Expected outcome: 生成符合最佳实践的 agents.md 文档结构和内容