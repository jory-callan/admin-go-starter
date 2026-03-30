package handler

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	"aicode/internal/app/core"
	"aicode/internal/model"
	"aicode/internal/repo"
	"aicode/pkg/goutils/echoutil"
	"aicode/pkg/goutils/response"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// QueryHandler SQL查询处理器
type QueryHandler struct {
	core         *core.App
	instanceRepo *repo.InstanceRepo
	historyRepo  *repo.QueryHistoryRepo
}

// NewQueryHandler 创建查询处理器
func NewQueryHandler(core *core.App) *QueryHandler {
	return &QueryHandler{
		core:         core,
		instanceRepo: repo.NewInstanceRepo(core.DB),
		historyRepo:  repo.NewQueryHistoryRepo(core.DB),
	}
}

// QueryRequest 查询请求
type QueryRequest struct {
	InstanceID string `json:"instance_id"`
	DBName     string `json:"db_name"`
	SQL        string `json:"sql"`
}

// 白名单允许的前缀（查询类SQL）
var allowedQueryPrefixes = []string{
	"SELECT", "SHOW", "DESCRIBE", "DESC", "EXPLAIN", "USE",
}

// 危险SQL模式（工单执行时才允许，工单外拒绝）
var dangerousPatterns = []string{
	"UPDATE", "DELETE", "ALTER", "DROP", "TRUNCATE", "INSERT", "CREATE",
	"REPLACE", "LOAD", "GRANT", "REVOKE", "RENAME", "TRUNCATE",
}

// isQueryAllowed 检查查询SQL是否在白名单内
func isQueryAllowed(sql string) bool {
	sql = strings.TrimSpace(sql)
	sqlUpper := strings.ToUpper(sql)

	// 检查是否以白名单前缀开头
	for _, prefix := range allowedQueryPrefixes {
		if strings.HasPrefix(sqlUpper, prefix) {
			return true
		}
	}

	return false
}

// isDangerousSQL 检查是否包含危险SQL关键词
func isDangerousSQL(sql string) bool {

	// 使用正则表达式匹配完整的单词
	for _, pattern := range dangerousPatterns {
		// 匹配单词边界，防止误匹配（如 UPDATE 匹配到 UPDATED）
		regex := regexp.MustCompile(`(?i)\b` + pattern + `\b`)
		if regex.MatchString(sql) {
			return true
		}
	}

	return false
}

// Query 执行查询SQL（所有人可执行，仅限查询类SQL）
func (h *QueryHandler) Query(c echo.Context) error {
	var req QueryRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, 400, "参数错误")
	}

	if req.InstanceID == "" || req.DBName == "" || req.SQL == "" {
		return response.Error(c, 400, "实例ID、数据库名、SQL语句不能为空")
	}

	// 验证SQL白名单
	if !isQueryAllowed(req.SQL) {
		return response.Error(c, 403, "只允许执行 SELECT, SHOW, DESCRIBE, EXPLAIN, USE 类型的SQL")
	}

	// 检查危险SQL
	if isDangerousSQL(req.SQL) {
		return response.Error(c, 403, "检测到危险SQL操作，已拒绝执行")
	}

	ctx := context.Background()
	userID := echoutil.GetUserID(c)

	// 获取实例信息
	instance, err := h.instanceRepo.GetByID(ctx, req.InstanceID)
	if err != nil {
		return response.Error(c, 404, "实例不存在")
	}

	// 执行查询
	startTime := time.Now()
	result, err := h.executeQuery(instance, req.DBName, req.SQL)
	duration := time.Since(startTime).Milliseconds()

	// 记录查询历史
	history := &model.QueryHistory{
		UserID:     userID,
		InstanceID: req.InstanceID,
		DBName:     req.DBName,
		SQLContent: req.SQL,
		Duration:   duration,
	}
	if result.Error != nil {
		history.ErrorMsg = result.Error.Error()
	} else {
		history.RowsAffected = int64(len(result.Rows))
	}

	h.historyRepo.Create(ctx, history)

	if result.Error != nil {
		return response.Error(c, 500, "查询失败: "+result.Error.Error())
	}

	return response.Success(c, map[string]any{
		"columns":  result.Columns,
		"rows":     result.Rows,
		"rows_num": len(result.Rows),
	})
}

// ExecuteTicket 执行工单SQL（仅 exec 角色可执行）
func (h *QueryHandler) ExecuteTicket(c echo.Context) error {
	ticketIDStr := c.Param("id")
	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil {
		return response.Error(c, 400, "无效的工单ID")
	}

	ctx := context.Background()
	executorID, _ := strconv.ParseInt(c.Get("user_id").(string), 10, 64)

	// 获取工单
	ticketRepo := repo.NewTicketRepo(h.core.DB)
	ticket, err := ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return response.Error(c, 404, "工单不存在")
	}

	// 检查工单状态
	if ticket.Status != model.TicketStatusPending {
		return response.Error(c, 400, "工单状态不是 PENDING，无法执行")
	}

	// 获取实例信息
	instance, err := h.instanceRepo.GetByID(ctx, ticket.InstanceID)
	if err != nil {
		return response.Error(c, 404, "关联的实例不存在")
	}

	// 执行工单SQL
	startTime := time.Now()
	execResult, err := h.executeQuery(instance, ticket.DBName, ticket.SQLContent)
	duration := time.Since(startTime).Milliseconds()

	// 更新工单状态
	var resultMsg string
	var newStatus model.TicketStatus
	if err != nil {
		newStatus = model.TicketStatusFailed
		resultMsg = err.Error()
	} else {
		newStatus = model.TicketStatusExecuted
		resultMsg = "执行成功，影响行数: " + strconv.FormatInt(int64(len(execResult.Rows)), 10)
	}

	// 更新工单状态
	if err := ticketRepo.UpdateStatus(ctx, ticketID, newStatus, &executorID, &resultMsg); err != nil {
		return response.Error(c, 500, "更新工单状态失败")
	}

	return response.Success(c, map[string]any{
		"status":   newStatus,
		"result":   resultMsg,
		"rows":     execResult.Rows,
		"duration": duration,
	})
}

// QueryResult 查询结果
type QueryResult struct {
	Columns []string
	Rows    [][]any
	Error   error
}

// executeQuery 执行SQL查询
func (h *QueryHandler) executeQuery(instance *model.Instance, dbName, sql string) (*QueryResult, error) {
	// TODO: 加密存储的密码在这里解密
	// adminPass := decryptPassword(instance.AdminPass)
	adminPass := instance.AdminPass

	dsn := instance.AdminUser + ":" + adminPass + "@tcp(" + instance.Host + ":" + strconv.Itoa(instance.Port) + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	defer func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()

	// 判断是查询还是执行
	sqlUpper := strings.ToUpper(strings.TrimSpace(sql))
	isQuery := strings.HasPrefix(sqlUpper, "SELECT") || strings.HasPrefix(sqlUpper, "SHOW") ||
		strings.HasPrefix(sqlUpper, "DESCRIBE") || strings.HasPrefix(sqlUpper, "DESC") ||
		strings.HasPrefix(sqlUpper, "EXPLAIN") || strings.HasPrefix(sqlUpper, "USE")

	if isQuery {
		// 查询语句
		rows, err := db.Raw(sql).Rows()
		if err != nil {
			return &QueryResult{Error: err}, err
		}
		defer rows.Close()

		columns, _ := rows.Columns()
		var results [][]any

		for rows.Next() {
			values := make([]any, len(columns))
			valuePtrs := make([]any, len(columns))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			if err := rows.Scan(valuePtrs...); err != nil {
				continue
			}

			row := make([]any, len(columns))
			for i, v := range values {
				row[i] = v
			}
			results = append(results, row)
		}

		return &QueryResult{
			Columns: columns,
			Rows:    results,
		}, nil
	} else {
		// 执行语句（UPDATE, DELETE 等）
		result := db.Exec(sql)
		if result.Error != nil {
			return &QueryResult{Error: result.Error}, result.Error
		}

		affected := result.RowsAffected
		return &QueryResult{
			Columns: []string{"Rows_Affected"},
			Rows:    [][]any{{affected}},
		}, nil
	}
}

// GetQueryHistory 获取查询历史
func (h *QueryHandler) GetQueryHistory(c echo.Context) error {
	var pq response.PageQuery
	if err := c.Bind(&pq); err != nil {
		pq = response.DefaultPageQuery()
	}

	ctx := context.Background()
	userID, _ := strconv.ParseInt(c.Get("user_id").(string), 10, 64)

	result, err := h.historyRepo.ListByUser(ctx, userID, &pq)
	if err != nil {
		return response.Error(c, 500, "获取查询历史失败")
	}

	return response.SuccessWithPage(c, *result)
}

// decryptPassword 解密密码（预留加密逻辑）
// func decryptPassword(encryptedPass string) string {
// 	// TODO: 实现 AES-256 解密
// 	return encryptedPass
// }
