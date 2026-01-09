package signature

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// SignatureMiddleware 返回签名验证中间件
// 仅对 /webhook 路由进行签名验证
func SignatureMiddleware(signer *Signer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 仅验证 webhook 路由
			if c.Request().URL.Path != "/webhook" {
				return next(c)
			}

			// 读取签名相关 Headers
			timestamp := c.Request().Header.Get("X-Signature-Timestamp")
			signature := c.Request().Header.Get("X-Signature-Ed25519")

			// 检查必要的 headers
			if timestamp == "" || signature == "" {
				slog.Warn("missing signature headers")
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "missing signature headers",
				})
			}

			// 验证时间戳（防重放攻击）
			if err := validateTimestamp(timestamp); err != nil {
				slog.Warn("invalid timestamp", "error", err, "timestamp", timestamp)
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid or expired timestamp",
				})
			}

			// 读取请求体
			body, err := io.ReadAll(c.Request().Body)
			if err != nil {
				slog.Error("failed to read request body", "error", err)
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "failed to read request body",
				})
			}

			// 重建 Body 以供后续处理使用
			c.Request().Body = io.NopCloser(bytes.NewReader(body))

			// 验证签名
			if err := signer.VerifySignature(timestamp, body, signature); err != nil {
				slog.Warn("signature verification failed", "error", err)
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "signature verification failed",
				})
			}

			slog.Info("webhook signature verified successfully")
			return next(c)
		}
	}
}

// validateTimestamp 验证时间戳是否在有效期内（5 分钟）
func validateTimestamp(timestampStr string) error {
	// 解析时间戳（假设是 Unix 时间戳，秒级）
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return err
	}

	requestTime := time.Unix(timestamp, 0)
	now := time.Now()

	// 检查时间戳是否在 5 分钟内
	diff := now.Sub(requestTime)
	if diff < 0 {
		diff = -diff
	}

	if diff > 5*time.Minute {
		return echo.NewHTTPError(http.StatusUnauthorized, "timestamp expired")
	}

	return nil
}
