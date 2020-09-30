package echomw

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func ZapLogger(log *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			if err != nil {
				c.Error(err)
			}
			req := c.Request()
			res := c.Response()

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}
			fields := []zapcore.Field{
				zap.String("time", time.Now().Format(time.RFC3339)),
				zap.Int("status", res.Status),
				zap.String("id", id),
				zap.String("latency", time.Since(start).String()),
				zap.String("method", req.Method),
				zap.String("host", req.Host),
				zap.String("remote_ip", c.RealIP()),
				zap.String("uri", req.RequestURI),
				zap.String("user_agent", req.UserAgent()),
			}
			n := res.Status
			switch {
			case n >= 500:
				log.Error("Server error", fields...)
			case n >= 400:
				log.Warn("Client error", fields...)
			case n >= 300:
				log.Info("Redirection", fields...)
			default:
				log.Info("Success", fields...)
			}
			return nil
		}
	}
}
