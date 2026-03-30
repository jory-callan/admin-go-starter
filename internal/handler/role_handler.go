package handler

import (
	"context"

	"aicode/internal/app/core"
	"aicode/internal/model"
	"aicode/internal/repo"
	"aicode/pkg/goutils/idutil"
	"aicode/pkg/goutils/response"

	"github.com/labstack/echo/v4"
)

// RoleHandler 角色处理器
type RoleHandler struct {
	core     *core.App
	roleRepo *repo.RoleRepo
}

// NewRoleHandler 创建角色处理器
func NewRoleHandler(core *core.App) *RoleHandler {
	return &RoleHandler{
		core:     core,
		roleRepo: repo.NewRoleRepo(core.DB),
	}
}

// Create 创建角色
func (h *RoleHandler) Create(c echo.Context) error {
	var req model.Role
	if err := c.Bind(&req); err != nil {
		return response.Error(c, 400, "参数错误")
	}

	if req.Code == "" || req.Name == "" {
		return response.Error(c, 400, "角色编码和名称不能为空")
	}

	ctx := context.Background()

	// 检查角色编码是否已存在
	existingRole, err := h.roleRepo.GetByCode(ctx, req.Code)
	if err == nil && existingRole != nil {
		return response.Error(c, 400, "角色编码已存在")
	}

	// 生成 ID
	req.ID = idutil.ShortUUIDv7()

	if err := h.roleRepo.Create(ctx, &req); err != nil {
		return response.Error(c, 500, "创建角色失败")
	}

	return response.Success(c, req)
}

// GetByID 根据ID获取角色
func (h *RoleHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.Error(c, 400, "ID不能为空")
	}

	ctx := context.Background()
	role, err := h.roleRepo.GetByID(ctx, id)
	if err != nil {
		return response.Error(c, 404, "角色不存在")
	}

	// 获取使用该角色的用户列表
	userRepo := repo.NewUserRepo(h.core.DB)
	users, _ := userRepo.ListByRoleID(ctx, id)

	return response.Success(c, map[string]any{
		"role":  role,
		"users": users,
	})
}

// Update 更新角色
func (h *RoleHandler) Update(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.Error(c, 400, "ID不能为空")
	}

	var req model.Role
	if err := c.Bind(&req); err != nil {
		return response.Error(c, 400, "参数错误")
	}

	if req.Code == "" || req.Name == "" {
		return response.Error(c, 400, "角色编码和名称不能为空")
	}

	ctx := context.Background()

	// 检查角色是否存在
	_, err := h.roleRepo.GetByID(ctx, id)
	if err != nil {
		return response.Error(c, 404, "角色不存在")
	}

	// 如果更新了角色编码，检查是否与其他角色冲突
	if req.Code != "" {
		existingRole, err := h.roleRepo.GetByCode(ctx, req.Code)
		if err == nil && existingRole != nil && existingRole.ID != id {
			return response.Error(c, 400, "角色编码已存在")
		}
	}

	req.ID = id
	if err := h.roleRepo.UpdateByID(ctx, &req, id); err != nil {
		return response.Error(c, 500, "更新角色失败")
	}

	return response.Success(c, nil)
}

// Delete 删除角色
func (h *RoleHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.Error(c, 400, "ID不能为空")
	}

	ctx := context.Background()
	operatorID := c.Get("user_id").(string)

	if err := h.roleRepo.Delete(ctx, id, operatorID); err != nil {
		return response.Error(c, 500, "删除角色失败")
	}

	return response.Success(c, nil)
}

// List 分页查询角色列表
func (h *RoleHandler) List(c echo.Context) error {
	var pq response.PageQuery
	if err := c.Bind(&pq); err != nil {
		pq = response.DefaultPageQuery()
	}

	ctx := context.Background()
	result, err := h.roleRepo.Pagination(ctx, &pq, nil)
	if err != nil {
		return response.Error(c, 500, "查询角色列表失败")
	}

	return response.SuccessWithPage(c, *result)
}

// GetAll 获取所有角色
func (h *RoleHandler) GetAll(c echo.Context) error {
	ctx := context.Background()
	db := h.roleRepo.GetDB(ctx)
	var roles []model.Role
	if err := db.Find(&roles).Error; err != nil {
		return response.Error(c, 500, "获取角色列表失败")
	}
	return response.Success(c, roles)
}
