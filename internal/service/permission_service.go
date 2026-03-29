package service

import (
	"aicode/internal/model"
	"aicode/internal/repo"
	"context"

	"gorm.io/gorm"
)

// PermissionService 权限服务
type PermissionService struct {
	permRepo *repo.PermissionRepo
}

// NewPermissionService 创建权限服务
func NewPermissionService(db *gorm.DB) *PermissionService {
	return &PermissionService{
		permRepo: repo.NewPermissionRepo(db),
	}
}

// List 权限列表
func (s *PermissionService) List(ctx context.Context, query *model.PageQuery) (*model.PageResult[model.Permission], error) {
	return s.permRepo.ListWithPagination(ctx, query)
}

// GetTree 获取权限树
func (s *PermissionService) GetTree(ctx context.Context) ([]model.Permission, error) {
	return s.permRepo.GetTree(ctx)
}

// GetByID 获取权限详情
func (s *PermissionService) GetByID(ctx context.Context, id string) (*model.Permission, error) {
	return s.permRepo.GetByID(ctx, id)
}

// Create 创建权限
func (s *PermissionService) Create(ctx context.Context, req *CreatePermissionRequest, operatorID string) error {
	perm := &model.Permission{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Type:        req.Type,
		Sort:        req.Sort,
		Status:      req.Status,
		ParentID:    req.ParentID,
		CreatedBy:   operatorID,
		UpdatedBy:   operatorID,
	}

	return s.permRepo.Create(ctx, perm)
}

// Update 更新权限
func (s *PermissionService) Update(ctx context.Context, id string, req *UpdatePermissionRequest, operatorID string) error {
	perm, err := s.permRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 更新字段
	if req.Name != "" {
		perm.Name = req.Name
	}
	if req.Description != "" {
		perm.Description = req.Description
	}
	if req.Type != 0 {
		perm.Type = req.Type
	}
	if req.Sort != 0 {
		perm.Sort = req.Sort
	}
	if req.Status != 0 {
		perm.Status = req.Status
	}
	perm.UpdatedBy = operatorID

	return s.permRepo.Update(ctx, perm)
}

// Delete 删除权限
func (s *PermissionService) Delete(ctx context.Context, id, operatorID string) error {
	perm, err := s.permRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	perm.DeletedBy = operatorID
	return s.permRepo.Delete(ctx, id)
}

// CreatePermissionRequest 创建权限请求
type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	Type        int    `json:"type"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status"`
	ParentID    string `json:"parent_id"`
}

// UpdatePermissionRequest 更新权限请求
type UpdatePermissionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        int    `json:"type"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status"`
}
