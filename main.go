package main

import (
	"errors"
	"log"
	"log/slog"
	"net/http"

	"go-qq-bot/internal/config"
	"go-qq-bot/internal/handler"
	"go-qq-bot/internal/service"
	"go-qq-bot/internal/signature"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	slog.Info("config loaded", "app_id", cfg.AppID, "port", cfg.Port)

	// 创建签名器
	signer, err := signature.NewSigner(cfg.AppSecret)
	if err != nil {
		log.Fatalf("Failed to create signer: %v", err)
	}
	slog.Info("ed25519 signer initialized")

	// 创建服务
	eventService := service.NewEventService()
	webhookHandler := handler.NewWebhookHandler(signer, eventService)
	slog.Info("services initialized")

	// 创建 Echo 实例
	e := echo.New()

	// 全局中间件
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// 签名验证中间件（仅对 /webhook 路由生效）
	e.Use(signature.SignatureMiddleware(signer))

	// 路由
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "QQ Bot Webhook Server")
	})
	e.POST("/webhook", webhookHandler.HandleWebhook)

	// 启动服务器
	slog.Info("starting server", "port", cfg.Port)
	if err := e.Start(cfg.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}
