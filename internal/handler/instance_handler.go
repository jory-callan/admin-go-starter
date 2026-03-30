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
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InstanceHandler 实例处理器
type InstanceHandler struct {
	instanceRepo *repo.InstanceRepo
}

// NewInstanceHandler 创建实例处理器
func NewInstanceHandler(core *core.App) *InstanceHandler {
	return &InstanceHandler{
		instanceRepo: repo.NewInstanceRepo(core.DB),
	}
}

// Create 创建数据库实例（需要 admin 或 db 角色）
func (h *InstanceHandler) Create(c echo.Context) error {
	var req model.Instance
	if err := c.Bind(&req); err != nil {
		return response.Error(c, 400, "参数错误")
	}

	if req.Name == "" || req.Host == "" || req.Port == 0 || req.AdminUser == "" {
		return response.Error(c, 400, "名称、主机、端口、管理员用户名不能为空")
	}

	ctx := context.Background()
	creatorID := echoutil.GetUserID(c)
	req.CreatedBy = creatorID

	if err := h.instanceRepo.Create(ctx, &req); err != nil {
		return response.Error(c, 500, "创建实例失败")
	}

	return response.Success(c, req)
}

// GetByID 获取实例详情
func (h *InstanceHandler) GetByID(c echo.Context) error {
	id := c.Param("id")

	ctx := context.Background()
	instance, err := h.instanceRepo.GetByID(ctx, id)
	if err != nil {
		return response.Error(c, 404, "实例不存在")
	}

	return response.Success(c, instance)
}

// Update 更新实例
func (h *InstanceHandler) Update(c echo.Context) error {
	id := c.Param("id")

	var req model.Instance
	if err := c.Bind(&req); err != nil {
		return response.Error(c, 400, "参数错误")
	}

	ctx := context.Background()

	// 检查实例是否存在
	_, err := h.instanceRepo.GetByID(ctx, id)
	if err != nil {
		return response.Error(c, 404, "实例不存在")
	}

	req.ID = id
	if err := h.instanceRepo.UpdateByID(ctx, &req, id); err != nil {
		return response.Error(c, 500, "更新实例失败")
	}

	return response.Success(c, nil)
}

// Delete 删除实例
func (h *InstanceHandler) Delete(c echo.Context) error {
	id := c.Param("id")

	operatorID := c.Get("user_id").(string)
	ctx := context.Background()

	if err := h.instanceRepo.Delete(ctx, id, operatorID); err != nil {
		return response.Error(c, 500, "删除实例失败")
	}

	return response.Success(c, nil)
}

// List 获取所有实例
func (h *InstanceHandler) List(c echo.Context) error {
	ctx := context.Background()
	instances, err := h.instanceRepo.ListAll(ctx)
	if err != nil {
		return response.Error(c, 500, "获取实例列表失败")
	}

	return response.Success(c, instances)
}

// GetDatabases 获取实例下的数据库列表
func (h *InstanceHandler) GetDatabases(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.Error(c, 400, "无效的实例ID")
	}

	ctx := context.Background()
	instance, err := h.instanceRepo.GetByID(ctx, id)
	if err != nil {
		return response.Error(c, 404, "实例不存在")
	}

	// 连接数据库获取数据库列表
	databases, err := h.fetchDatabases(instance)
	if err != nil {
		return response.Error(c, 500, "获取数据库列表失败: "+err.Error())
	}

	return response.Success(c, databases)
}

// GetTables 获取指定数据库下的表列表
func (h *InstanceHandler) GetTables(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.Error(c, 400, "无效的实例ID")
	}

	dbName := c.QueryParam("db")
	if dbName == "" {
		return response.Error(c, 400, "数据库名称不能为空")
	}

	ctx := context.Background()
	instance, err := h.instanceRepo.GetByID(ctx, id)
	if err != nil {
		return response.Error(c, 404, "实例不存在")
	}

	tables, err := h.fetchTables(instance, dbName)
	if err != nil {
		return response.Error(c, 500, "获取表列表失败: "+err.Error())
	}

	return response.Success(c, tables)
}

// GetColumns 获取指定表的列信息
func (h *InstanceHandler) GetColumns(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.Error(c, 400, "无效的实例ID")
	}

	dbName := c.QueryParam("db")
	tableName := c.QueryParam("table")
	if dbName == "" || tableName == "" {
		return response.Error(c, 400, "数据库名称和表名不能为空")
	}

	ctx := context.Background()
	instance, err := h.instanceRepo.GetByID(ctx, id)
	if err != nil {
		return response.Error(c, 404, "实例不存在")
	}

	columns, err := h.fetchColumns(instance, dbName, tableName)
	if err != nil {
		return response.Error(c, 500, "获取列信息失败: "+err.Error())
	}

	return response.Success(c, columns)
}

// 连接实例获取数据库列表
func (h *InstanceHandler) fetchDatabases(instance *model.Instance) ([]string, error) {
	// TODO: 加密存储的密码在这里解密
	// adminPass := decryptPassword(instance.AdminPass)
	adminPass := instance.AdminPass

	dsn := instance.AdminUser + ":" + adminPass + "@tcp(" + instance.Host + ":" + strconv.Itoa(instance.Port) + ")/?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gormOpen(dsn)
	if err != nil {
		return nil, err
	}
	defer closeDB(db)

	var databases []string
	rows, err := db.Raw("SHOW DATABASES").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err == nil {
			// 过滤掉系统数据库
			if dbName != "information_schema" && dbName != "mysql" && dbName != "performance_schema" && dbName != "sys" {
				databases = append(databases, dbName)
			}
		}
	}

	return databases, nil
}

// 连接实例获取表列表
func (h *InstanceHandler) fetchTables(instance *model.Instance, dbName string) ([]string, error) {
	adminPass := instance.AdminPass
	dsn := instance.AdminUser + ":" + adminPass + "@tcp(" + instance.Host + ":" + strconv.Itoa(instance.Port) + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gormOpen(dsn)
	if err != nil {
		return nil, err
	}
	defer closeDB(db)

	var tables []string
	rows, err := db.Raw("SHOW TABLES").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err == nil {
			tables = append(tables, table)
		}
	}

	return tables, nil
}

// 连接实例获取列信息
func (h *InstanceHandler) fetchColumns(instance *model.Instance, dbName, tableName string) ([]map[string]interface{}, error) {
	adminPass := instance.AdminPass
	dsn := instance.AdminUser + ":" + adminPass + "@tcp(" + instance.Host + ":" + strconv.Itoa(instance.Port) + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gormOpen(dsn)
	if err != nil {
		return nil, err
	}
	defer closeDB(db)

	var columns []map[string]any
	rows, err := db.Raw("DESCRIBE " + tableName).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var field, colType, null, key string
		var defaultVal, extra any
		if err := rows.Scan(&field, &colType, &null, &key, &defaultVal, &extra); err == nil {
			columns = append(columns, map[string]any{
				"field":   field,
				"type":    colType,
				"null":    null,
				"key":     key,
				"default": defaultVal,
				"extra":   extra,
			})
		}
	}

	return columns, nil
}

// gormOpen 创建 gorm DB 连接
func gormOpen(dsn string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

// closeDB 关闭数据库连接
func closeDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

// decryptPassword 解密密码（预留加密逻辑）
// func decryptPassword(encryptedPass string) string {
// 	// TODO: 实现 AES-256 解密
// 	return encryptedPass
// }
