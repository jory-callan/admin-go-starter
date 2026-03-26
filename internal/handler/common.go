package handler

import (
	"aicode/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

// Success 成功响应
func Success(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, &model.Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

// Fail 失败响应
func Fail(c echo.Context, code int, msg string) error {
	return c.JSON(http.StatusOK, &model.Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

// SuccessPage 分页响应
func SuccessPage(c echo.Context, result *model.PageResult[interface{}]) error {
	return Success(c, map[string]interface{}{
		"items":   result.Items,
		"total":   result.Total,
		"page":    result.Page,
		"size":    result.Size,
		"hasMore": result.HasMore,
	})
}
