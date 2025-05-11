package main

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"sdk_test/database"
	fc "sdk_test/fabric"
	pb "sdk_test/proto"
	sc "sdk_test/service"
	"sdk_test/wallet"
	wl "sdk_test/wallet"
)

// 🔧 MockWallet 模擬 wallet.Exists 與 PutFile 的行為
type MockWallet struct{}

func (w *MockWallet) Exists(userID string) bool {
	return true // 模擬錢包已存在
}

func (w *MockWallet) PutFile(userID, certPath, keyPath, mspID string) error {
	return nil // 模擬成功寫入錢包
}

func (m *MockWallet) Get(userID string) (*wallet.Entry, bool) {
	return nil, true // 或你想模擬的行為
}

// 🧪 測試 UploadReport（假設 HandleUploadReport 成功）
func TestUploadReport(t *testing.T) {
	s := &server{
		Wallet:  &wl.Wallet{},
		Builder: fc.GWBuilder{},
	}
	req := &pb.UploadReportRequest{
		ReportId:        "r001",
		PatientHash:     "hash001",
		TestResultsJson: `{"LDL": "123"}`,
	}

	resp, err := s.UploadReport(context.Background(), req)

	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

// 🧪 測試 ClaimReport 是否正常回傳成功訊息
func TestClaimReport(t *testing.T) {
	s := &server{}

	req := &pb.ClaimReportRequest{
		ReportId: "r001",
	}

	resp, err := s.ClaimReport(context.Background(), req)

	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Report Claimed Successfully", resp.Message)
}

// 🧪 測試 ReadReport 是否回傳假資料內容
func TestReadReport(t *testing.T) {
	s := &server{}

	req := &pb.ReadReportRequest{
		ReportId: "r001",
	}

	resp, err := s.ReadReport(context.Background(), req)

	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Fake report content...", resp.ReportContent)
}

// 🧪 測試 Register（模擬成功流程）
func TestRegister(t *testing.T) {
	// ✅ 初始化 SQLite 測試資料庫（避免 nil pointer）
	err := database.InitDB("test_user_data.sqlite")
	assert.NoError(t, err)
	defer os.Remove("test_user_data.sqlite") // 測試後刪除

	mockWallet := &MockWallet{}

	req := &pb.RegisterRequest{
		UserId:   "test_user",
		Password: "test_pass",
		Name:     "測試用戶",
		Date:     "2024-05-01",
		Email:    "test@example.com",
		Phone:    "0987654321",
	}

	resp, err := sc.HandleRegister(context.Background(), req, mockWallet)

	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "註冊成功", resp.Message)
}

// 🧪 測試 Login（模擬成功登入）
func TestLogin(t *testing.T) {
	mockWallet := &MockWallet{}

	req := &pb.LoginRequest{
		UserId:   "test_user",
		Password: "test_pass",
	}

	resp, err := sc.HandleLogin(context.Background(), req, mockWallet)

	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "登入成功", resp.Message)
	assert.NotEmpty(t, resp.Token)
}
