package router

import (
	"aicode/internal/handler"
)

// registerProtected 注册需要认证的接口
func (r *Router) registerProtected() {
	// 用户相关接口
	userHandler := handler.NewUserHandler(r.core)
	roleHandler := handler.NewRoleHandler(r.core)
	instanceHandler := handler.NewInstanceHandler(r.core)
	ticketHandler := handler.NewTicketHandler(r.core)
	queryHandler := handler.NewQueryHandler(r.core)

	// 需要 JWT 认证的路由
	auth := r.group.Group("/auth", r.mw.JWTAuth())
	auth.POST("/logout", userHandler.Logout)
	auth.GET("/current_user", userHandler.GetCurrentUser)
	auth.POST("/change_password", userHandler.ChangePassword)

	// 用户管理
	users := r.group.Group("/users", r.mw.JWTAuth())
	users.POST("", userHandler.Create)
	users.GET("", userHandler.List)
	users.GET("/:id", userHandler.GetByID)
	users.PUT("/:id", userHandler.Update)
	users.DELETE("/:id", userHandler.Delete)
	users.POST("/:id/roles", userHandler.AssignRoles)

	// 角色管理
	roles := r.group.Group("/roles", r.mw.JWTAuth())
	roles.POST("", roleHandler.Create)
	roles.GET("", roleHandler.List)
	roles.GET("/all", roleHandler.GetAll)
	roles.GET("/:id", roleHandler.GetByID)
	roles.PUT("/:id", roleHandler.Update)
	roles.DELETE("/:id", roleHandler.Delete)

	// 数据库实例管理（需要 admin 或 db 角色）
	instances := r.group.Group("/instances", r.mw.JWTAuth())
	instances.POST("", instanceHandler.Create)
	instances.GET("", instanceHandler.List)
	instances.GET("/:id", instanceHandler.GetByID)
	instances.PUT("/:id", instanceHandler.Update)
	instances.DELETE("/:id", instanceHandler.Delete)
	// 实例下的数据库和表结构查询
	instances.GET("/:id/databases", instanceHandler.GetDatabases)
	instances.GET("/:id/tables", instanceHandler.GetTables)
	instances.GET("/:id/columns", instanceHandler.GetColumns)

	// 工单管理
	tickets := r.group.Group("/tickets", r.mw.JWTAuth())
	tickets.POST("", ticketHandler.Create)
	tickets.GET("", ticketHandler.List)
	tickets.GET("/:id", ticketHandler.GetByID)
	tickets.PUT("/:id", ticketHandler.Update)
	tickets.DELETE("/:id", ticketHandler.Delete)

	// SQL查询（所有人可执行查询，exec角色可执行工单）
	queryHandlerGroup := r.group.Group("/query", r.mw.JWTAuth())
	queryHandlerGroup.POST("/execute", queryHandler.Query)                     // 执行查询SQL
	queryHandlerGroup.POST("/tickets/:id/execute", queryHandler.ExecuteTicket) // 执行工单
	queryHandlerGroup.GET("/history", queryHandler.GetQueryHistory)            // 查询历史
}
