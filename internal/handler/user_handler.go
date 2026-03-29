package handler

import (
	"aicode/internal/model"
	"aicode/internal/service"
	"aicode/pkg/goutils/echoutil"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(db),
	}
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
	query.NeedCount = true

	result, err := h.userService.List(c.Request().Context(), &query)
	if err != nil {
		return Fail(c, CodeInternalError, err.Error())
	}

	// 转换为 PageResult[interface{}]
	items := make([]interface{}, len(result.Items))
	for i, item := range result.Items {
		items[i] = item
	}

	return Success(c, model.PageResult[interface{}]{
		Items:   items,
		Total:   result.Total,
		Page:    result.Page,
		Size:    result.Size,
		HasMore: result.HasMore,
	})
}

// Get 获取用户详情
func (h *UserHandler) Get(c echo.Context) error {
	id := c.Param("id")

	user, err := h.userService.GetByID(c.Request().Context(), id)
	if err != nil {
		return Fail(c, CodeNotFound, "用户不存在")
	}

	return Success(c, user)
}

// Create 创建用户
func (h *UserHandler) Create(c echo.Context) error {
	var req service.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return Fail(c, CodeBadRequest, "参数错误")
	}

	operatorID := echoutil.GetUserID(c)

	if err := h.userService.Create(c.Request().Context(), &req, operatorID); err != nil {
		return Fail(c, CodeInternalError, err.Error())
	}

	return Success(c, nil)
}

// Update 更新用户
func (h *UserHandler) Update(c echo.Context) error {
	id := c.Param("id")

	var req service.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return Fail(c, CodeBadRequest, "参数错误")
	}

	operatorID := c.Get("user_id").(string)
	if err := h.userService.Update(c.Request().Context(), id, &req, operatorID); err != nil {
		return Fail(c, CodeInternalError, err.Error())
	}

	return Success(c, nil)
}

// Delete 删除用户
func (h *UserHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	operatorID := c.Get("user_id").(string)

	if err := h.userService.Delete(c.Request().Context(), id, operatorID); err != nil {
		return Fail(c, CodeInternalError, err.Error())
	}

	return Success(c, nil)
}

// AssignRoles 为用户分配角色
func (h *UserHandler) AssignRoles(c echo.Context) error {
	id := c.Param("id")

	var req struct {
		RoleIDs []string `json:"role_ids"`
	}
	if err := c.Bind(&req); err != nil {
		return Fail(c, CodeBadRequest, "参数错误")
	}

	if err := h.userService.AssignRoles(c.Request().Context(), id, req.RoleIDs); err != nil {
		return Fail(c, CodeInternalError, err.Error())
	}

	return Success(c, nil)
}
