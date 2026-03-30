package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type PageQuery struct {
	Page      int         `query:"page" json:"page"` // 当前页，从1开始
	Size      int         `query:"size" json:"size"`
	NeedCount bool        `query:"need_count" json:"needCount"` // 是否需要总数（性能敏感场景可关闭）
	Order     []OrderItem `query:"order" json:"order"`          // 支持多字段排序
}

// OrderItem 表示单个排序字段和排序方式
type OrderItem struct {
	Field string `json:"field" query:"field"` // 数据库字段名
	Desc  bool   `json:"desc" query:"desc"`   // true: DESC, false: ASC
}

func (q *PageQuery) Validate() {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Size <= 0 {
		q.Size = 10
	}
	if q.Size > 100 {
		q.Size = 100
	}
}

func DefaultPageQuery() PageQuery {
	return PageQuery{
		Page: 1,
		Size: 10,
	}
}

// PageResult 分页结果（响应 data 里的内容）
type PageResult struct {
	Items      any    `json:"items"`
	Total      int64  `json:"total,omitempty"`
	Page       int    `json:"page"`
	Size       int    `json:"size"`
	HasMore    bool   `json:"hasMore,omitempty"`
	NextCursor string `json:"nextCursor,omitempty"` // 预留 cursor 分页
	// 未来扩展位
	// TotalPages int    `json:"totalPages,omitempty"`
	// Stats      any    `json:"stats,omitempty"`
}

func NewPageResult(items []any, page, size int, total int64) PageResult {
	return PageResult{
		Items:   items,
		Page:    page,
		Size:    size,
		Total:   total,
		HasMore: int64(page*size) < total,
	}
}

func SuccessWithPage(c echo.Context, result PageResult) error {
	return c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "success",
		Data: result,
	})
}
