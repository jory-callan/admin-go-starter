package app

import (
	"aicode/pkg/jwt"
	"log/slog"
)

// initJWT 初始化 JWT 全局参数
func (a *App) initJWT() {
	jwt.Init(a.Config.JWT)
	slog.Info("initialized", "issuer", a.Config.JWT.Issuer, "expires", a.Config.JWT.Expires)
}
