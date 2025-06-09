package utils

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var jwtSecret = []byte("my_super_secret") // 可放環境變數

// 產生 JWT Token
func GenerateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // 設定過期時間（1天）
	})
	return token.SignedString(jwtSecret)
}

// 驗證 token 並解析 userID（內部用）
func ValidateJWT(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", err
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", err
	}

	return userID, nil
}

// ✅ 封裝版本：從 context 取得 userID
func ExtractUserIDFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "沒有 metadata")
	}

	authHeader := md["authorization"]
	if len(authHeader) == 0 {
		return "", status.Error(codes.Unauthenticated, "請提供 JWT token")
	}

	token := authHeader[0]
	userID, err := ValidateJWT(token)
	if err != nil {
		return "", status.Error(codes.Unauthenticated, "JWT 驗證失敗")
	}
	return userID, nil
}
