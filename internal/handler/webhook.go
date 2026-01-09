package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"go-qq-bot/internal/model"
	"go-qq-bot/internal/service"
	"go-qq-bot/internal/signature"

	"github.com/labstack/echo/v4"
)

// WebhookHandler Webhook 请求处理器
type WebhookHandler struct {
	signer       *signature.Signer
	eventService *service.EventService
}

// NewWebhookHandler 创建 Webhook 处理器
func NewWebhookHandler(signer *signature.Signer, eventService *service.EventService) *WebhookHandler {
	return &WebhookHandler{
		signer:       signer,
		eventService: eventService,
	}
}

// HandleWebhook 处理 webhook POST 请求
func (h *WebhookHandler) HandleWebhook(c echo.Context) error {
	var payload model.WebhookPayload
	if err := c.Bind(&payload); err != nil {
		slog.Error("failed to bind webhook payload", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid json payload",
		})
	}

	slog.Info("webhook payload received", "op", payload.Op, "type", payload.T, "id", payload.ID)

	// 根据 op 字段路由
	switch payload.Op {
	case 13:
		// 验证请求
		return h.handleVerification(c, payload)
	case 0:
		// 事件处理
		return h.handleEvent(c, payload)
	default:
		slog.Warn("unsupported opcode", "op", payload.Op)
		return c.JSON(http.StatusOK, nil)
	}
}

// handleVerification 处理验证请求（op=13）
func (h *WebhookHandler) handleVerification(c echo.Context, payload model.WebhookPayload) error {
	slog.Info("handling verification challenge")

	// 解析验证数据
	var verifyData model.VerificationPayload
	data, err := json.Marshal(payload.D)
	if err != nil {
		slog.Error("failed to marshal verification data", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal error",
		})
	}

	if err := json.Unmarshal(data, &verifyData); err != nil {
		slog.Error("failed to unmarshal verification data", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid verification payload",
		})
	}

	// 生成签名
	msgSig, err := h.signer.SignVerification(verifyData.EventTs, verifyData.PlainToken)
	if err != nil {
		slog.Error("failed to sign verification", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to generate signature",
		})
	}

	// 返回验证响应
	response := model.VerificationResponse{
		PlainToken:   verifyData.PlainToken,
		MsgSig:       msgSig,
		ResponceTime: verifyData.EventTs,
	}

	slog.Info("verification challenge completed successfully")
	return c.JSON(http.StatusOK, response)
}

// handleEvent 处理事件（op=0）
func (h *WebhookHandler) handleEvent(c echo.Context, payload model.WebhookPayload) error {
	slog.Info("handling event", "type", payload.T, "id", payload.ID)

	// 解析事件数据
	var eventPayload model.EventPayload
	data, err := json.Marshal(payload.D)
	if err != nil {
		slog.Error("failed to marshal event data", "error", err)
		// 即使处理失败，也返回 200 避免 QQ 服务器重试
		return c.JSON(http.StatusOK, nil)
	}

	if err := json.Unmarshal(data, &eventPayload); err != nil {
		slog.Warn("failed to unmarshal event data", "error", err, "type", payload.T)
		// 即使处理失败，也返回 200 避免 QQ 服务器重试
		return c.JSON(http.StatusOK, nil)
	}

	// 调用事件处理服务
	if err := h.eventService.Handle(c.Request().Context(), model.EventType(payload.T), &eventPayload); err != nil {
		slog.Error("failed to handle event", "error", err, "type", payload.T)
		// 即使处理失败，也返回 200 避免 QQ 服务器重试
		return c.JSON(http.StatusOK, nil)
	}

	slog.Info("event processed successfully", "type", payload.T)
	return c.JSON(http.StatusOK, nil)
}
