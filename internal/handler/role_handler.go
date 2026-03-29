package handler

import (
	"aicode/internal/model"
	"aicode/internal/service"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// RoleHandler 角色处理器
type RoleHandler struct {
	roleService *service.RoleService
}

// NewRoleHandler 创建角色处理器
func NewRoleHandler(db *gorm.DB) *RoleHandler {
	return &RoleHandler{
		roleService: service.NewRoleService(db),
	}
}

// List 角色列表
func (h *RoleHandler) List(c echo.Context) error {
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

	result, err := h.roleService.List(c.Request().Context(), &query)
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

// Get 获取角色详情
func (h *RoleHandler) Get(c echo.Context) error {
	id := c.Param("id")

	role, err := h.roleService.GetByID(c.Request().Context(), id)
	if err != nil {
		return Fail(c, CodeNotFound, "角色不存在")
	}

	return Success(c, role)
}

// Create 创建角色
func (h *RoleHandler) Create(c echo.Context) error {
	var req service.CreateRoleRequest
	if err := c.Bind(&req); err != nil {
		return Fail(c, CodeBadRequest, "参数错误")
	}

	operatorID := c.Get("user_id").(string)
	if err := h.roleService.Create(c.Request().Context(), &req, operatorID); err != nil {
		return Fail(c, CodeInternalError, err.Error())
	}

	return Success(c, nil)
}

// Update 更新角色
func (h *RoleHandler) Update(c echo.Context) error {
	id := c.Param("id")

	var req service.UpdateRoleRequest
	if err := c.Bind(&req); err != nil {
		return Fail(c, CodeBadRequest, "参数错误")
	}

	operatorID := c.Get("user_id").(string)
	if err := h.roleService.Update(c.Request().Context(), id, &req, operatorID); err != nil {
		return Fail(c, CodeInternalError, err.Error())
	}

	return Success(c, nil)
}

// Delete 删除角色
func (h *RoleHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	operatorID := c.Get("user_id").(string)

	if err := h.roleService.Delete(c.Request().Context(), id, operatorID); err != nil {
		return Fail(c, CodeInternalError, err.Error())
	}

	return Success(c, nil)
}

// AssignPermissions 为角色分配权限
func (h *RoleHandler) AssignPermissions(c echo.Context) error {
	id := c.Param("id")

	var req struct {
		PermissionIDs []string `json:"permission_ids"`
	}
	if err := c.Bind(&req); err != nil {
		return Fail(c, CodeBadRequest, "参数错误")
	}

	if err := h.roleService.AssignPermissions(c.Request().Context(), id, req.PermissionIDs); err != nil {
		return Fail(c, CodeInternalError, err.Error())
	}

	return Success(c, nil)
}
