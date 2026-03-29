package service

import (
	"aicode/internal/model"
	"aicode/internal/repo"
	"context"
	"gorm.io/gorm"
)

// RoleService 角色服务
type RoleService struct {
	roleRepo *repo.RoleRepo
}

// NewRoleService 创建角色服务
func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{
		roleRepo: repo.NewRoleRepo(db),
	}
}

// List 角色列表
func (s *RoleService) List(ctx context.Context, query *model.PageQuery) (*model.PageResult[model.Role], error) {
	return s.roleRepo.ListWithPermissions(ctx, query)
}

// GetByID 获取角色详情
func (s *RoleService) GetByID(ctx context.Context, id string) (*model.Role, error) {
	return s.roleRepo.GetByIDWithPermissions(ctx, id)
}

// Create 创建角色
func (s *RoleService) Create(ctx context.Context, req *CreateRoleRequest, operatorID string) error {
	role := &model.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Sort:        req.Sort,
		Status:      req.Status,
		CreatedBy:   operatorID,
		UpdatedBy:   operatorID,
	}

	if err := s.roleRepo.Create(ctx, role); err != nil {
		return err
	}

	// 分配权限
	if len(req.PermissionIDs) > 0 {
		return s.roleRepo.AssignPermissions(ctx, role.ID, req.PermissionIDs)
	}

	return nil
}

// Update 更新角色
func (s *RoleService) Update(ctx context.Context, id string, req *UpdateRoleRequest, operatorID string) error {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 更新字段
	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.Sort != 0 {
		role.Sort = req.Sort
	}
	if req.Status != 0 {
		role.Status = req.Status
	}
	role.UpdatedBy = operatorID

	if err := s.roleRepo.Update(ctx, role); err != nil {
		return err
	}

	// 更新权限
	if req.PermissionIDs != nil {
		if err := s.roleRepo.AssignPermissions(ctx, id, req.PermissionIDs); err != nil {
			return err
		}
	}

	return nil
}

// Delete 删除角色
func (s *RoleService) Delete(ctx context.Context, id, operatorID string) error {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	role.DeletedBy = operatorID
	return s.roleRepo.Delete(ctx, id)
}

// AssignPermissions 为角色分配权限
func (s *RoleService) AssignPermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	return s.roleRepo.AssignPermissions(ctx, roleID, permissionIDs)
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name          string   `json:"name" binding:"required"`
	Code          string   `json:"code" binding:"required"`
	Description   string   `json:"description"`
	Sort          int      `json:"sort"`
	Status        int      `json:"status"`
	PermissionIDs []string `json:"permission_ids"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Sort          int      `json:"sort"`
	Status        int      `json:"status"`
	PermissionIDs []string `json:"permission_ids"`
}
