package config

import (
	"fmt"
	"os"
)

// Config 包含 QQ 机器人的配置信息
type Config struct {
	AppID     string // NITRO_BOT_APP_ID
	Token     string // NITRO_BOT_TOKEN
	AppSecret string // NITRO_BOT_APP_SECRET
	Port      string // 服务器端口
}

// Load 从环境变量加载配置
func Load() (*Config, error) {
	cfg := &Config{
		AppID:     os.Getenv("NITRO_BOT_APP_ID"),
		Token:     os.Getenv("NITRO_BOT_TOKEN"),
		AppSecret: os.Getenv("NITRO_BOT_APP_SECRET"),
		Port:      os.Getenv("SERVER_PORT"),
	}

	// 设置默认端口
	if cfg.Port == "" {
		cfg.Port = ":8080"
	}

	// 验证必要的配置项
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate 验证配置的有效性
func (c *Config) Validate() error {
	if c.AppID == "" {
		return fmt.Errorf("NITRO_BOT_APP_ID is required")
	}
	if c.Token == "" {
		return fmt.Errorf("NITRO_BOT_TOKEN is required")
	}
	if c.AppSecret == "" {
		return fmt.Errorf("NITRO_BOT_APP_SECRET is required")
	}
	return nil
}
