// pkg/response/response.go
package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

const (
	CodeSuccess = 200
	CodeError   = 500
)

func Success(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  "success",
		Data: data,
	})
}

func SuccessWithMsg(c echo.Context, msg string, data any) error {
	return c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  msg,
		Data: data,
	})
}

func Error(c echo.Context, code int, msg string) error {
	return c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

func ErrorWithHTTPStatus(c echo.Context, httpStatus, code int, msg string) error {
	return c.JSON(httpStatus, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
