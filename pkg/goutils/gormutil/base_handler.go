package gormutil

import (
	"net/http"

	"aicode/pkg/goutils/echoutil"
	"aicode/pkg/goutils/response"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// BaseHandler 极简泛型处理器
// 仅负责：参数绑定 -> 调用 Repo -> 统一响应
// 不包含任何自动审计、钩子或隐式逻辑
type BaseHandler[T any] struct {
	Repo *BaseRepo[T]
}

// NewBaseHandler 构造函数
func NewBaseHandler[T any](repo *BaseRepo[T]) *BaseHandler[T] {
	return &BaseHandler[T]{Repo: repo}
}

// Create 创建单条 - POST /
// 注意：实体中的 created_by 等字段需在上层业务逻辑或 Service 层赋值，此处不做处理
func (h *BaseHandler[T]) Create(c echo.Context) error {
	ctx := c.Request().Context()

	var entity T
	// 1. 绑定参数
	if err := c.Bind(&entity); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request body")
	}

	// 2. 直接调用 Repo (无隐式逻辑)
	if err := h.Repo.Create(ctx, &entity); err != nil {
		return response.Error(c, http.StatusInternalServerError, "create failed: "+err.Error())
	}

	// 3. 返回结果
	return response.Success(c, entity)
}

// GetByID 查询单条 - GET /:id
func (h *BaseHandler[T]) GetByID(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	entity, err := h.Repo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Error(c, http.StatusNotFound, "resource not found")
		}
		return response.Error(c, http.StatusInternalServerError, "query failed: "+err.Error())
	}

	return response.Success(c, entity)
}

// Pagination 分页查询 - GET /?page=1&size=10
func (h *BaseHandler[T]) Pagination(c echo.Context) error {
	ctx := c.Request().Context()

	var pq response.PageQuery
	if err := c.Bind(&pq); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid query params")
	}

	result, err := h.Repo.PaginationWithScopes(ctx, &pq)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "query failed: "+err.Error())
	}

	return response.SuccessWithPage(c, *result)
}

func (h *BaseHandler[T]) UpdateByID(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	// 获取操作人 ID
	// operatorID := echoutil.GetUserID(c)

	var entity T
	if err := c.Bind(&entity); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request body")
	}

	if err := h.Repo.UpdateByID(ctx, &entity, id); err != nil {
		return response.Error(c, http.StatusInternalServerError, "update failed: "+err.Error())
	}

	return response.Success(c, entity)
}

// Delete 删除单条 - DELETE /:id
// 这里需要获取 operatorID，根据您的要求使用 echoutil.GetUserID(c)
func (h *BaseHandler[T]) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	// 获取操作人 ID
	operatorID := echoutil.GetUserID(c)

	if err := h.Repo.Delete(ctx, id, operatorID); err != nil {
		return response.Error(c, http.StatusInternalServerError, "delete failed: "+err.Error())
	}

	return response.Success(c, nil)
}
