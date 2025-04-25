package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// HealthReport 定義每一筆健檢報告的結構
// TestResults 用 JSON 字串統一保存所有健檢數據
// 這樣新增欄位時，不需要改 Chaincode

type HealthReport struct {
	ReportID     string `json:"reportID"`      // 報告唯一編號
	PatientHash  string `json:"patientHash"`   // 病患健保卡號 Hash
	AccessList   []string `json:"accessList"`   // 有權限讀取此報告的 CN 名單
	Timestamp    string `json:"timestamp"`      // 上傳時間
	TestResults  string `json:"testResults"`    // 🆕 存所有健檢項目的 JSON 字串
}

// HealthCheckContract 實作 Chaincode

type HealthCheckContract struct {
	contractapi.Contract
}

// UploadReport 上傳一份健檢報告
func (h *HealthCheckContract) UploadReport(ctx contractapi.TransactionContextInterface, reportID, patientHash, testResultsJson string) error {
	exists, err := ctx.GetStub().GetState(reportID)
	if err != nil {
		return err
	}
	if exists != nil {
		return fmt.Errorf("report %s already exists", reportID)
	}

	// 檢查 testResultsJson 是否是合法 JSON（可選）
	var test map[string]interface{}
	err = json.Unmarshal([]byte(testResultsJson), &test)
	if err != nil {
		return fmt.Errorf("invalid testResults JSON: %v", err)
	}

	report := HealthReport{
		ReportID:      reportID,
		PatientHash:   patientHash,
		AccessList:    []string{},
		Timestamp:     time.Now().Format(time.RFC3339),
		TestResults:   testResultsJson,
	}

	reportBytes, err := json.Marshal(report)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(reportID, reportBytes)
}

// ClaimReport 使用者認領自己的報告，將自己加入 AccessList
func (h *HealthCheckContract) ClaimReport(ctx contractapi.TransactionContextInterface, reportID string) error {
	reportBytes, err := ctx.GetStub().GetState(reportID)
	if err != nil || reportBytes == nil {
		return fmt.Errorf("report %s not found", reportID)
	}

	var report HealthReport
	_ = json.Unmarshal(reportBytes, &report)

	cert, err := ctx.GetClientIdentity().GetX509Certificate()
	if err != nil {
		return fmt.Errorf("failed to get client certificate: %v", err)
	}
	callerCN := cert.Subject.CommonName

	for _, cn := range report.AccessList {
		if cn == callerCN {
			return fmt.Errorf("identity %s already authorized", callerCN)
		}
	}

	report.AccessList = append(report.AccessList, callerCN)

	updatedBytes, err := json.Marshal(report)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(reportID, updatedBytes)
}

// ReadReport 查詢報告內容（需要在 AccessList 中）
func (h *HealthCheckContract) ReadReport(ctx contractapi.TransactionContextInterface, reportID string) (*HealthReport, error) {
	reportBytes, err := ctx.GetStub().GetState(reportID)
	if err != nil || reportBytes == nil {
		return nil, fmt.Errorf("report %s not found", reportID)
	}

	var report HealthReport
	_ = json.Unmarshal(reportBytes, &report)

	cert, err := ctx.GetClientIdentity().GetX509Certificate()
	if err != nil {
		return nil, fmt.Errorf("failed to get client certificate: %v", err)
	}
	callerCN := cert.Subject.CommonName

	authorized := false
	for _, cn := range report.AccessList {
		if cn == callerCN {
			authorized = true
			break
		}
	}
	if !authorized {
		return nil, fmt.Errorf("access denied for %s", callerCN)
	}

	return &report, nil
}

// GrantAccess 授權其他人查閱報告
func (h *HealthCheckContract) GrantAccess(ctx contractapi.TransactionContextInterface, reportID, targetCN string) error {
	reportBytes, err := ctx.GetStub().GetState(reportID)
	if err != nil || reportBytes == nil {
		return fmt.Errorf("report %s not found", reportID)
	}

	var report HealthReport
	_ = json.Unmarshal(reportBytes, &report)

	cert, err := ctx.GetClientIdentity().GetX509Certificate()
	if err != nil {
		return fmt.Errorf("failed to get caller certificate: %v", err)
	}
	callerCN := cert.Subject.CommonName

	authorized := false
	for _, cn := range report.AccessList {
		if cn == callerCN {
			authorized = true
			break
		}
	}
	if !authorized {
		return fmt.Errorf("access denied for %s", callerCN)
	}

	for _, cn := range report.AccessList {
		if cn == targetCN {
			return fmt.Errorf("%s already has access", targetCN)
		}
	}

	report.AccessList = append(report.AccessList, targetCN)

	updatedBytes, err := json.Marshal(report)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(reportID, updatedBytes)
}

// RevokeAccess 收回他人查閱權限
func (h *HealthCheckContract) RevokeAccess(ctx contractapi.TransactionContextInterface, reportID, targetCN string) error {
	reportBytes, err := ctx.GetStub().GetState(reportID)
	if err != nil || reportBytes == nil {
		return fmt.Errorf("report %s not found", reportID)
	}

	var report HealthReport
	_ = json.Unmarshal(reportBytes, &report)

	cert, err := ctx.GetClientIdentity().GetX509Certificate()
	if err != nil {
		return fmt.Errorf("failed to get caller certificate: %v", err)
	}
	callerCN := cert.Subject.CommonName

	authorized := false
	for _, cn := range report.AccessList {
		if cn == callerCN {
			authorized = true
			break
		}
	}
	if !authorized {
		return fmt.Errorf("access denied for %s", callerCN)
	}

	// 移除目標 CN
	filtered := []string{}
	for _, cn := range report.AccessList {
		if cn != targetCN {
			filtered = append(filtered, cn)
		}
	}
	report.AccessList = filtered

	updatedBytes, err := json.Marshal(report)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(reportID, updatedBytes)
}

// main 啟動鏈碼
func main() {
	chaincode, err := contractapi.NewChaincode(&HealthCheckContract{})
	if err != nil {
		panic(fmt.Sprintf("Error creating chaincode: %v", err))
	}

	if err := chaincode.Start(); err != nil {
		panic(fmt.Sprintf("Error starting chaincode: %v", err))
	}
}
