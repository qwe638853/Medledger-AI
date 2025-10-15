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

/**
 * @notice 產生 JWT Token
 * @dev 使用 HS256 簽章，payload 含 user_id 與 24h 過期時間
 * @param userID 使用者識別
 * @return string 簽章後的 JWT, error 失敗
 */
// 產生 JWT Token
func GenerateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // 設定過期時間（1天）
	})
	return token.SignedString(jwtSecret)
}

/**
 * @notice 驗證 JWT 並回傳 userID
 * @dev 以預設密鑰驗簽；解析 payload 中的 user_id
 * @param tokenStr JWT 字串
 * @return string userID, error 驗證或解析失敗
 */
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

/**
 * @notice 從 gRPC Context 的 metadata 解析出 userID
 * @dev 讀取 authorization 欄位作為 JWT 來源，並呼叫 ValidateJWT
 * @param ctx gRPC 請求上下文
 * @return string userID, error 未授權或驗證失敗
 */
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
