package signature

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
)

// Signer Ed25519 签名器，用于验证和生成签名
type Signer struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

// NewSigner 从 AppSecret 创建签名器
// 根据 QQ 官方算法，将 AppSecret 转换为 32 字节 seed，然后生成密钥对
func NewSigner(appSecret string) (*Signer, error) {
	if appSecret == "" {
		return nil, fmt.Errorf("appSecret cannot be empty")
	}

	// 从 AppSecret 派生 32 字节 seed
	seed, err := deriveKeyFromSecret(appSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key from secret: %w", err)
	}

	// 从 seed 生成 Ed25519 密钥对
	privateKey := ed25519.NewKeyFromSeed(seed)
	publicKey := privateKey.Public().(ed25519.PublicKey)

	return &Signer{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

// deriveKeyFromSecret 从 AppSecret 派生 32 字节 seed
// 根据 QQ 官方算法：重复填充 secret 字符串直到达到 32 字节
func deriveKeyFromSecret(secret string) ([]byte, error) {
	if secret == "" {
		return nil, fmt.Errorf("secret cannot be empty")
	}

	seed := make([]byte, ed25519.SeedSize) // ed25519.SeedSize = 32
	secretBytes := []byte(secret)

	// 重复填充 secret 到 32 字节
	// 例如：长度 16 的 secret 会重复 2 次
	for i := 0; i < ed25519.SeedSize; i++ {
		seed[i] = secretBytes[i%len(secretBytes)]
	}

	return seed, nil
}

// VerifySignature 验证 webhook 签名
// 验证消息格式：timestamp + body
func (s *Signer) VerifySignature(timestamp string, body []byte, signatureHex string) error {
	// 解码 hex 格式的签名
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return fmt.Errorf("invalid signature format: %w", err)
	}

	// Ed25519 签名长度固定为 64 字节
	if len(signature) != ed25519.SignatureSize {
		return fmt.Errorf("invalid signature length: expected %d, got %d",
			ed25519.SignatureSize, len(signature))
	}

	// 构建待验证消息：timestamp + body
	message := append([]byte(timestamp), body...)

	// 使用 Ed25519 验证签名
	if !ed25519.Verify(s.publicKey, message, signature) {
		return fmt.Errorf("signature verification failed")
	}

	return nil
}

// SignVerification 生成验证响应签名（用于响应 op=13 验证请求）
// 签名消息格式：event_ts + plain_token
func (s *Signer) SignVerification(eventTs, plainToken string) (string, error) {
	// 构建待签名消息：event_ts + plain_token
	message := []byte(eventTs + plainToken)

	// 使用私钥签名
	signature := ed25519.Sign(s.privateKey, message)

	// 转换为 hex 字符串返回
	return hex.EncodeToString(signature), nil
}
