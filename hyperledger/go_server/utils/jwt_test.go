package utils

import (
    "context"
    "testing"
    "time"

    "google.golang.org/grpc/metadata"
)

// TestGenerateAndValidateJWT
// 目的：驗證 JWT 產生與驗證行為
func TestGenerateAndValidateJWT(t *testing.T) {
    tok, err := GenerateJWT("alice")
    if err != nil || tok == "" {
        t.Fatalf("GenerateJWT: %v", err)
    }
    user, err := ValidateJWT(tok)
    if err != nil || user != "alice" {
        t.Fatalf("ValidateJWT mismatch: %v user=%s", err, user)
    }
}

// TestExtractUserIDFromContext
// 目的：驗證從 gRPC metadata 取出 token 並解析 userID
func TestExtractUserIDFromContext(t *testing.T) {
    tok, err := GenerateJWT("bob")
    if err != nil { t.Fatalf("GenerateJWT: %v", err) }

    // 模擬 gRPC incoming metadata
    md := metadata.New(map[string]string{"authorization": tok})
    ctx := metadata.NewIncomingContext(context.Background(), md)

    user, err := ExtractUserIDFromContext(ctx)
    if err != nil || user != "bob" {
        t.Fatalf("ExtractUserIDFromContext: %v user=%s", err, user)
    }

    // 過期測試（簡易）：產生極短期 token 並等待過期（此處僅示範，不做實際 sleep 測試）
    _ = time.Now() // 保留時間相關引用（未實作過期快速測）
}



