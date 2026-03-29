package router

import (
	"aicode/internal/handler"
)

// registerPublic 注册公开接口（无需认证）
func (r *Router) registerPublic() {
	auth := handler.NewAuthHandler(r.core.DB)

	public := r.group.Group("/auth")
	{
		public.POST("/login", auth.Login)
	}
}
