// pkg/response/response.go
package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ApiResponse[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

const (
	CodeSuccess = 200
	CodeError   = 500
)

func Success(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, ApiResponse[any]{
		Code: CodeSuccess,
		Msg:  "success",
		Data: data,
	})
}

func SuccessWithMsg[T any](c echo.Context, msg string, data T) error {
	return c.JSON(http.StatusOK, ApiResponse[T]{
		Code: CodeSuccess,
		Msg:  msg,
		Data: data,
	})
}

func Error(c echo.Context, code int, msg string) error {
	return c.JSON(http.StatusOK, ApiResponse[any]{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

func ErrorWithHTTPStatus(c echo.Context, httpStatus, code int, msg string) error {
	return c.JSON(httpStatus, ApiResponse[any]{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
