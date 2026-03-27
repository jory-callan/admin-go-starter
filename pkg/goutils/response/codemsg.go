package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type CodeMsg struct {
	Code int
	Msg  string
}

var (
	CodeMsgOK       = CodeMsg{Code: 200, Msg: "success"}
	CodeMsgErr      = CodeMsg{Code: 500, Msg: "server error"}
	CodeMsgErrParam = CodeMsg{Code: 400, Msg: "parameter error"}
)

func (cm CodeMsg) String() string {
	return cm.Msg
}

func (cm CodeMsg) GetCode() int {
	return cm.Code
}

func (cm CodeMsg) GetMsg() string {
	return cm.Msg
}

// 业务成功，返回200状态码，带数据
func SuccessWithCodeMsg[T any](c echo.Context, cm CodeMsg) error {
	res := ApiResponse[T]{
		Code: cm.Code,
		Msg:  cm.Msg,
	}
	return c.JSON(http.StatusOK, res)
}

// 业务成功，返回200状态码，带数据
func SuccessWithCodeMsgWithData[T any](c echo.Context, cm CodeMsg, data T) error {
	res := ApiResponse[T]{
		Code: cm.Code,
		Msg:  cm.Msg,
		Data: data,
	}
	return c.JSON(http.StatusOK, res)
}

// 业务报错，返回200状态码，不带数据
func ErrorWithCodeMsg[T any](c echo.Context, cm CodeMsg) error {
	res := ApiResponse[T]{
		Code: cm.Code,
		Msg:  cm.Msg,
	}
	return c.JSON(http.StatusOK, res)
}

// 系统报错,500状态码
func SystemErrorWithCodeMsg[T any](c echo.Context) error {
	return c.JSON(http.StatusInternalServerError, CodeMsgErr)
}
