package handler

import (
	"context"
	"log/slog"

	"aicode/internal/app/core"
	"aicode/internal/model"
	"aicode/internal/repo"
	"aicode/pkg/goutils/echoutil"
	"aicode/pkg/goutils/idutil"
	"aicode/pkg/goutils/response"
	"aicode/pkg/jwt"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// UserHandler 用户处理器
type UserHandler struct {
	userRepo     *repo.UserRepo
	roleRepo     *repo.RoleRepo
	userRoleRepo *repo.UserRoleRepo
}

// NewUserHandler 创建用户处理器
func NewUserHandler(core *core.App) *UserHandler {
	return &UserHandler{
		userRepo:     repo.NewUserRepo(core.DB),
		roleRepo:     repo.NewRoleRepo(core.DB),
		userRoleRepo: repo.NewUserRoleRepo(core.DB),
	}
}

// Create 创建用户
func (h *UserHandler) Create(c echo.Context) error {
	var req model.User
	if err := c.Bind(&req); err != nil {
		return response.Error(c, 400, "参数错误")
	}

	// 生成 ID
	req.ID = idutil.UUIDv7()

	// 加密密码
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			slog.Error("密码加密失败", "error", err)
			return response.Error(c, 500, err.Error())
		}
		req.Password = string(hash)
	}

	ctx := context.Background()
	if err := h.userRepo.Create(ctx, &req); err != nil {
		return response.Error(c, 500, err.Error())
	}

	return response.Success(c, req)
}

// GetByID 根据ID获取用户
func (h *UserHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.Error(c, 400, "ID不能为空")
	}

	ctx := context.Background()
	user, err := h.userRepo.GetByID(ctx, id)
	if err != nil {
		return response.Error(c, 404, "用户不存在")
	}
	// 去掉密码
	user.Password = ""

	// 获取用户角色
	roles, _ := h.roleRepo.ListByUserID(ctx, id)
	return response.Success(c, map[string]any{
		"user":  user,
		"roles": roles,
	})
}

// Update 更新用户
func (h *UserHandler) Update(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.Error(c, 400, "ID不能为空")
	}

	var req model.User
	if err := c.Bind(&req); err != nil {
		return response.Error(c, 400, "参数错误")
	}

	ctx := context.Background()

	// 获取原用户
	oldUser, err := h.userRepo.GetByID(ctx, id)
	if err != nil {
		return response.Error(c, 404, "用户不存在")
	}

	// 如果更新了密码，需要重新加密
	if req.Password != "" && req.Password != oldUser.Password {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return response.Error(c, 500, "密码加密失败")
		}
		req.Password = string(hash)
	}

	req.ID = id
	if err := h.userRepo.UpdateByID(ctx, &req, id); err != nil {
		return response.Error(c, 500, "更新用户失败")
	}

	return response.Success(c, nil)
}

// Delete 删除用户
func (h *UserHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.Error(c, 400, "ID不能为空")
	}

	ctx := context.Background()
	operatorID := c.Get("user_id").(string)

	if err := h.userRepo.Delete(ctx, id, operatorID); err != nil {
		return response.Error(c, 500, "删除用户失败")
	}

	return response.Success(c, nil)
}

// List 分页查询用户列表
func (h *UserHandler) List(c echo.Context) error {
	var pq response.PageQuery
	if err := c.Bind(&pq); err != nil {
		pq = response.DefaultPageQuery()
	}

	ctx := context.Background()
	result, err := h.userRepo.Pagination(ctx, &pq, nil)
	if err != nil {
		return response.Error(c, 500, "查询用户列表失败")
	}

	// 将password字段置空
	for _, user := range result.Items.([]*model.User) {
		user.Password = ""
	}

	return response.SuccessWithPage(c, *result)
}

// Login 用户登录
func (h *UserHandler) Login(c echo.Context) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return response.Error(c, 400, "参数错误")
	}

	if req.Username == "" || req.Password == "" {
		return response.Error(c, 400, "用户名密码不能为空")
	}

	ctx := context.Background()

	// 根据用户名查询用户
	user, err := h.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return response.Error(c, 401, "用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		return response.Error(c, 401, "用户已被禁用")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return response.Error(c, 401, "用户名或密码错误")
	}

	// 获取用户角色
	roles, err := h.roleRepo.ListByUserID(ctx, user.ID)
	if err != nil {
		roles = []model.Role{}
	}

	// 生成角色编码列表
	roleCodes := make([]string, len(roles))
	for i, role := range roles {
		roleCodes[i] = role.Code
	}

	// 生成 JWT token
	token, err := jwt.GenerateToken(user.ID, user.Username, roleCodes, nil)
	if err != nil {
		return response.Error(c, 500, "生成token失败")
	}

	return response.Success(c, map[string]any{
		"token": token,
		"user":  user,
		"role":  roleCodes,
	})
}

// Logout 用户登出
func (h *UserHandler) Logout(c echo.Context) error {
	// JWT 无状态登出，只需要告知客户端清除 token
	// 如果使用 Redis 黑名单机制，可以在这里添加
	return response.SuccessWithMsg(c, "登出成功", nil)
}

// Register 注册新用户
func (h *UserHandler) Register(c echo.Context) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
	}
	if err := c.Bind(&req); err != nil {
		return response.Error(c, 400, "参数错误")
	}

	if req.Username == "" || req.Password == "" {
		return response.Error(c, 400, "用户名密码不能为空")
	}

	ctx := context.Background()

	// 检查用户名是否已存在
	_, err := h.userRepo.GetByUsername(ctx, req.Username)
	if err == nil {
		return response.Error(c, 400, "用户名已存在")
	}

	// 创建用户
	user := &model.User{
		ID:       idutil.UUIDv7(),
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   1,
	}

	// 加密密码
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return response.Error(c, 500, "密码加密失败")
	}
	user.Password = string(hash)

	if err := h.userRepo.Create(ctx, user); err != nil {
		return response.Error(c, 500, "注册用户失败")
	}

	return response.Success(c, user)
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(c echo.Context) error {
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.Bind(&req); err != nil {
		return response.Error(c, 400, "参数错误")
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		return response.Error(c, 400, "旧密码和新密码不能为空")
	}

	userID := c.Get("user_id").(string)
	ctx := context.Background()

	// 获取当前用户
	user, err := h.userRepo.GetByID(ctx, userID)
	if err != nil {
		return response.Error(c, 404, "用户不存在")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return response.Error(c, 400, "旧密码错误")
	}

	// 加密新密码
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return response.Error(c, 500, "密码加密失败")
	}

	// 更新密码
	user.Password = string(hash)
	if err := h.userRepo.Update(ctx, user); err != nil {
		return response.Error(c, 500, "修改密码失败")
	}

	return response.SuccessWithMsg(c, "密码修改成功", nil)
}

// GetCurrentUser 获取当前用户信息
func (h *UserHandler) GetCurrentUser(c echo.Context) error {
	// userID := c.Get("user_id").(string)
	userID := echoutil.GetUserID(c)
	if userID == "" {
		return response.Error(c, 401, "未登录")
	}

	ctx := context.Background()

	user, err := h.userRepo.GetByID(ctx, userID)
	if err != nil {
		return response.Error(c, 404, "用户不存在")
	}

	// 获取用户角色
	roles, _ := h.roleRepo.ListByUserID(ctx, userID)

	return response.Success(c, map[string]any{
		"user":  user,
		"roles": roles,
	})
}

// AssignRoles 给用户分配角色
func (h *UserHandler) AssignRoles(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return response.Error(c, 400, "用户ID不能为空")
	}

	var req struct {
		RoleIDs []string `json:"role_ids"`
	}
	if err := c.Bind(&req); err != nil {
		return response.Error(c, 400, "参数错误")
	}

	ctx := context.Background()

	// 检查用户是否存在
	_, err := h.userRepo.GetByID(ctx, userID)
	if err != nil {
		return response.Error(c, 404, "用户不存在")
	}

	if err := h.userRoleRepo.SetRolesForUser(ctx, userID, req.RoleIDs); err != nil {
		return response.Error(c, 500, "分配角色失败")
	}

	return response.SuccessWithMsg(c, "角色分配成功", nil)
}
