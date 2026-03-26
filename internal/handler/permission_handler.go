package handler

import (
	"aicode/internal/model"
	"aicode/internal/service"
	"github.com/labstack/echo/v4"
)

// PermissionHandler 权限处理器
type PermissionHandler struct {
	permService *service.PermissionService
}

// NewPermissionHandler 创建权限处理器
func NewPermissionHandler(permService *service.PermissionService) *PermissionHandler {
	return &PermissionHandler{
		permService: permService,
	}
}

// List 权限列表
func (h *PermissionHandler) List(c echo.Context) error {
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

	result, err := h.permService.List(c.Request().Context(), &query)
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

// GetTree 获取权限树
func (h *PermissionHandler) GetTree(c echo.Context) error {
	tree, err := h.permService.GetTree(c.Request().Context())
	if err != nil {
		return Fail(c, CodeInternalError, err.Error())
	}

	return Success(c, tree)
}

// Get 获取权限详情
func (h *PermissionHandler) Get(c echo.Context) error {
	id := c.Param("id")

	perm, err := h.permService.GetByID(c.Request().Context(), id)
	if err != nil {
		return Fail(c, CodeNotFound, "权限不存在")
	}

	return Success(c, perm)
}

// Create 创建权限
func (h *PermissionHandler) Create(c echo.Context) error {
	var req service.CreatePermissionRequest
	if err := c.Bind(&req); err != nil {
		return Fail(c, CodeBadRequest, "参数错误")
	}

	operatorID := c.Get("user_id").(string)
	if err := h.permService.Create(c.Request().Context(), &req, operatorID); err != nil {
		return Fail(c, CodeInternalError, err.Error())
	}

	return Success(c, nil)
}

// Update 更新权限
func (h *PermissionHandler) Update(c echo.Context) error {
	id := c.Param("id")

	var req service.UpdatePermissionRequest
	if err := c.Bind(&req); err != nil {
		return Fail(c, CodeBadRequest, "参数错误")
	}

	operatorID := c.Get("user_id").(string)
	if err := h.permService.Update(c.Request().Context(), id, &req, operatorID); err != nil {
		return Fail(c, CodeInternalError, err.Error())
	}

	return Success(c, nil)
}

// Delete 删除权限
func (h *PermissionHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	operatorID := c.Get("user_id").(string)

	if err := h.permService.Delete(c.Request().Context(), id, operatorID); err != nil {
		return Fail(c, CodeInternalError, err.Error())
	}

	return Success(c, nil)
}
