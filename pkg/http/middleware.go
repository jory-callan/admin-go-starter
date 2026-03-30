package http

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"aicode/pkg/goutils/idutil"
	"aicode/pkg/goutils/response"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"golang.org/x/time/rate"
)

func Recover() echo.MiddlewareFunc {
	return echoMiddleware.RecoverWithConfig(echoMiddleware.RecoverConfig{
		StackSize: 2 << 10,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			errMsg := fmt.Sprintf("http panic error --> errmsg: %s , stack: %s", err.Error(), string(stack))
			slog.Error(errMsg)
			_ = response.Error(c, http.StatusInternalServerError, "unknown error")
			return nil
		},
	})
}

func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			start := time.Now()
			err := next(c)
			latency := time.Since(start)

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}
			if id == "" {
				id = idutil.ShortUUIDv7()
			}
			res.Header().Set(echo.HeaderXRequestID, id)

			status := res.Status
			if err != nil {
				var he *echo.HTTPError
				if errors.As(err, &he) {
					status = he.Code
				}
			}

			args := []any{
				"id", id,
				"request_id", id,
				"remote_ip", c.RealIP(),
				"host", req.Host,
				"proto", req.Proto,
				"method", req.Method,
				"uri", req.RequestURI,
				"path", req.URL.Path,
				"route", c.Path(),
				"user_agent", req.UserAgent(),
				"referer", req.Referer(),
				"status", status,
				"latency", latency,
				"latency_human", latency.String(),
				"bytes_in", req.ContentLength,
				"bytes_out", res.Size,
			}

			if qs := req.URL.RawQuery; qs != "" {
				args = append(args, "query_string", qs)
			}
			if contentType := req.Header.Get("Content-Type"); contentType != "" {
				args = append(args, "content_type", contentType)
			}

			if err != nil {
				args = append(args, "error", err.Error())
			}

			switch {
			case status >= 500:
				slog.Error("http request 500 error --> ", args...)
			case status >= 400:
				slog.Warn("http request 400 error --> ", args...)
			default:
				slog.Info("http request --> ", args...)
			}

			return err
		}
	}
}

func ErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}
		var he *echo.HTTPError
		if c.Request().Method == http.MethodHead {
			return
		}

		if errors.As(err, &he) {
			_ = c.JSON(he.Code, map[string]any{
				"code": he.Code,
				"msg":  he.Message,
				"data": nil,
			})
			return
		}

		expose := c.Request().Header.Get("Debug") == "true" || c.Request().Header.Get("From") == "in"

		msg := "internal error"
		if expose && err != nil {
			msg = err.Error()
			requestID := c.Response().Header().Get(echo.HeaderXRequestID)
			slog.Error(msg, "request_id", requestID)
		}

		_ = c.JSON(http.StatusOK, map[string]any{
			"code": "-1",
			"msg":  msg,
			"data": nil,
		})
	}
}

func CORS(cfg CORSConfig) echo.MiddlewareFunc {
	if len(cfg.AllowOrigins) == 0 {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}
	}
	return echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     cfg.AllowMethods,
		AllowHeaders:     cfg.AllowHeaders,
		ExposeHeaders:    cfg.ExposeHeaders,
		MaxAge:           cfg.MaxAge,
		AllowCredentials: true,
	})
}

func RateLimit(cfg RateLimitConfig) echo.MiddlewareFunc {
	if !cfg.Enabled {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}
	}
	store := echoMiddleware.NewRateLimiterMemoryStore(rate.Limit(cfg.RequestsPerSecond))
	return echoMiddleware.RateLimiter(store)
}
