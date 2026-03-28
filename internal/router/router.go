package router

import (
	app "aicode/internal/app/core"
	"aicode/internal/handler"
	"aicode/internal/middleware"
	"aicode/internal/repo"
	"aicode/internal/service"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes 注册所有路由并装配依赖
func RegisterRoutes(e *echo.Echo, application *app.App) {
	// ============ 依赖装配（组合根） ============

	// Repo 层
	userRepo := repo.NewUserRepo(application.DB)
	roleRepo := repo.NewRoleRepo(application.DB)
	permRepo := repo.NewPermissionRepo(application.DB)

	// Service 层
	authService := service.NewAuthService(userRepo, roleRepo)
	userService := service.NewUserService(userRepo, roleRepo)
	roleService := service.NewRoleService(roleRepo)
	permService := service.NewPermissionService(permRepo)

	// Handler 层
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	roleHandler := handler.NewRoleHandler(roleService)
	permHandler := handler.NewPermissionHandler(permService)

	// ============ 路由注册 ============

	// 健康检查
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	// API v1
	v1 := e.Group("/api/v1")

	// ========== 公开接口 ==========

	// 登录认证
	auth := v1.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
	}

	// ========== 需要认证的接口 ==========

	// JWT 认证中间件
	v1.Use(middleware.JWTAuth())

	// 当前用户信息
	v1.GET("/user/info", authHandler.GetUserInfo)

	// ========== 用户管理 ==========
	users := v1.Group("/users")
	{
		users.GET("", userHandler.List)                                                                      // 用户列表
		users.GET("/:id", userHandler.Get)                                                                   // 用户详情
		users.POST("", userHandler.Create, middleware.RequirePermission("system:user:create"))               // 创建用户
		users.PUT("/:id", userHandler.Update, middleware.RequirePermission("system:user:update"))            // 更新用户
		users.DELETE("/:id", userHandler.Delete, middleware.RequirePermission("system:user:delete"))         // 删除用户
		users.PUT("/:id/roles", userHandler.AssignRoles, middleware.RequirePermission("system:user:assign")) // 分配角色
	}

	// ========== 角色管理 ==========
	roles := v1.Group("/roles")
	{
		roles.GET("", roleHandler.List)                                                                                  // 角色列表
		roles.GET("/:id", roleHandler.Get)                                                                               // 角色详情
		roles.POST("", roleHandler.Create, middleware.RequirePermission("system:role:create"))                           // 创建角色
		roles.PUT("/:id", roleHandler.Update, middleware.RequirePermission("system:role:update"))                        // 更新角色
		roles.DELETE("/:id", roleHandler.Delete, middleware.RequirePermission("system:role:delete"))                     // 删除角色
		roles.PUT("/:id/permissions", roleHandler.AssignPermissions, middleware.RequirePermission("system:role:assign")) // 分配权限
	}

	// ========== 权限管理 ==========
	permissions := v1.Group("/permissions")
	{
		permissions.GET("", permHandler.List)                                                                    // 权限列表
		permissions.GET("/tree", permHandler.GetTree)                                                            // 权限树
		permissions.GET("/:id", permHandler.Get)                                                                 // 权限详情
		permissions.POST("", permHandler.Create, middleware.RequirePermission("system:permission:create"))       // 创建权限
		permissions.PUT("/:id", permHandler.Update, middleware.RequirePermission("system:permission:update"))    // 更新权限
		permissions.DELETE("/:id", permHandler.Delete, middleware.RequirePermission("system:permission:delete")) // 删除权限
	}
}
