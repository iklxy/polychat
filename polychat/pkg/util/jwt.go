package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret 密钥
var jwtSecret = []byte("polychat_secret_key_2026")

// Claims 自定义载荷结构体
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 Token，有效期 24 小时
func GenerateToken(userID uint) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)), // 过期时间 24 小时
			IssuedAt:  jwt.NewNumericDate(now),                     // 签发时间
			Issuer:    "polychat",                                  // 签发人
		},
	}

	// 使用 HS256 签名算法
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 生成签名后的 Token 字符串
	return tokenClaims.SignedString(jwtSecret)
}

// ParseToken 解析 Token
func ParseToken(tokenString string) (*Claims, error) {
	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法是否为 HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// 验证 Token 是否有效并提取 Claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}
