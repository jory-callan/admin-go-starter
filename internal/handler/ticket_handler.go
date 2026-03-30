package handler

import (
	"context"
	"strconv"

	"aicode/internal/app/core"
	"aicode/internal/model"
	"aicode/internal/repo"
	"aicode/pkg/goutils/echoutil"
	"aicode/pkg/goutils/response"

	"github.com/labstack/echo/v4"
)

// TicketHandler 工单处理器
type TicketHandler struct {
	core       *core.App
	ticketRepo *repo.TicketRepo
}

// NewTicketHandler 创建工单处理器
func NewTicketHandler(core *core.App) *TicketHandler {
	return &TicketHandler{
		core:       core,
		ticketRepo: repo.NewTicketRepo(core.DB),
	}
}

// Create 创建工单
func (h *TicketHandler) Create(c echo.Context) error {
	var req struct {
		Title      string `json:"title"`
		SQLContent string `json:"sql_content"`
		InstanceID string `json:"instance_id"`
		DBName     string `json:"db_name"`
	}
	if err := c.Bind(&req); err != nil {
		return response.Error(c, 400, "参数错误")
	}

	if req.Title == "" || req.SQLContent == "" || req.InstanceID == "" || req.DBName == "" {
		return response.Error(c, 400, "标题、SQL内容、实例ID、数据库名不能为空")
	}

	ctx := context.Background()
	creatorID := echoutil.GetUserID(c)

	ticket := &model.Ticket{
		Title:      req.Title,
		SQLContent: req.SQLContent,
		InstanceID: req.InstanceID,
		DBName:     req.DBName,
		Status:     model.TicketStatusPending,
		CreatorID:  creatorID,
	}

	if err := h.ticketRepo.Create(ctx, ticket); err != nil {
		return response.Error(c, 500, "创建工单失败")
	}

	return response.Success(c, ticket)
}

// GetByID 获取工单详情
func (h *TicketHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.Error(c, 400, "无效的工单ID")
	}

	ctx := context.Background()
	ticket, err := h.ticketRepo.GetByID(ctx, id)
	if err != nil {
		return response.Error(c, 404, "工单不存在")
	}

	return response.Success(c, ticket)
}

// List 查询工单列表
// 普通用户只能看到自己的工单，admin 和 exec 角色可以看到所有工单
func (h *TicketHandler) List(c echo.Context) error {
	var pq response.PageQuery
	if err := c.Bind(&pq); err != nil {
		pq = response.DefaultPageQuery()
	}

	ctx := context.Background()
	userID, _ := strconv.ParseInt(c.Get("user_id").(string), 10, 64)

	// 检查用户角色
	hasAllPermission := false
	if roles, ok := c.Get("roles").([]string); ok {
		for _, role := range roles {
			if role == "admin" || role == "exec" {
				hasAllPermission = true
				break
			}
		}
	}

	var result *response.PageResult
	var err error

	if hasAllPermission {
		// admin 和 exec 角色可以看到所有工单
		result, err = h.ticketRepo.ListAll(ctx, &pq)
	} else {
		// 普通用户只能看到自己的工单
		result, err = h.ticketRepo.ListByCreator(ctx, userID, &pq)
	}

	if err != nil {
		return response.Error(c, 500, "查询工单列表失败")
	}

	return response.SuccessWithPage(c, *result)
}

// Update 更新工单（只能更新自己的 PENDING 状态的工单）
func (h *TicketHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.Error(c, 400, "无效的工单ID")
	}

	var req struct {
		Title      string `json:"title"`
		SQLContent string `json:"sql_content"`
	}
	if err := c.Bind(&req); err != nil {
		return response.Error(c, 400, "参数错误")
	}

	ctx := context.Background()
	userID := echoutil.GetUserID(c)

	// 获取原工单
	ticket, err := h.ticketRepo.GetByID(ctx, id)
	if err != nil {
		return response.Error(c, 404, "工单不存在")
	}

	// 检查是否是创建者且状态是 PENDING
	if ticket.CreatorID != userID || ticket.Status != model.TicketStatusPending {
		return response.Error(c, 403, "只能更新自己创建的 PENDING 状态的工单")
	}

	// 更新字段
	if req.Title != "" {
		ticket.Title = req.Title
	}
	if req.SQLContent != "" {
		ticket.SQLContent = req.SQLContent
	}

	if err := h.ticketRepo.Update(ctx, ticket); err != nil {
		return response.Error(c, 500, "更新工单失败")
	}

	return response.Success(c, nil)
}

// Delete 删除工单（只能删除自己的 PENDING 状态的工单）
func (h *TicketHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.Error(c, 400, "无效的工单ID")
	}

	ctx := context.Background()
	userID := echoutil.GetUserID(c)

	// 获取原工单
	ticket, err := h.ticketRepo.GetByID(ctx, id)
	if err != nil {
		return response.Error(c, 404, "工单不存在")
	}

	// 检查是否是创建者且状态是 PENDING
	if ticket.CreatorID != userID || ticket.Status != model.TicketStatusPending {
		return response.Error(c, 403, "只能删除自己创建的 PENDING 状态的工单")
	}

	if err := h.ticketRepo.Delete(ctx, idStr, c.Get("user_id").(string)); err != nil {
		return response.Error(c, 500, "删除工单失败")
	}

	return response.Success(c, nil)
}
