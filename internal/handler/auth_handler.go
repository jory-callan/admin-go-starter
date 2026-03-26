package handler

import (
	"aicode/internal/model"
	"aicode/internal/service"
	"github.com/labstack/echo/v4"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login 用户登录
func (h *AuthHandler) Login(c echo.Context) error {
	var req model.LoginRequest
	if err := c.Bind(&req); err != nil {
		return Fail(c, CodeBadRequest, "参数错误")
	}

	result, err := h.authService.Login(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		if err == service.ErrUserDisabled {
			return Fail(c, CodeBadRequest, "用户已被禁用")
		}
		return Fail(c, CodeBadRequest, "用户名或密码错误")
	}

	return Success(c, result)
}

// GetUserInfo 获取当前用户信息
func (h *AuthHandler) GetUserInfo(c echo.Context) error {
	userID := c.Get("user_id").(string)

	userInfo, err := h.authService.GetUserInfo(c.Request().Context(), userID)
	if err != nil {
		return Fail(c, CodeNotFound, "用户不存在")
	}

	return Success(c, userInfo)
}
