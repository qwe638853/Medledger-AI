package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// HealthReport 定義每一筆健檢報告的結構
// TestResults 用 JSON 字串統一保存所有健檢數據
// 這樣新增欄位時，不需要改 Chaincode

// 全域變數
const (
	docHealth = "HealthRecord" // 報告檔案標籤
	docAuth   = "AuthTicket"   // 授權票標籤

	keyReportNS = "REPORT" // CompositeKey namespace
	keyAuthNS   = "AUTH"
)

type HealthReport struct {
	DocType    string `json:"docType"`     // 固定 "HealthRecord"
	ReportID   string `json:"reportId"`    // 主鍵，由診所自訂
	PatientH   string `json:"patientHash"` // 患者 username 雜湊
	ClinicID   string `json:"clinicId"`    // 產生報告的診所 ID
	ResultJSON string `json:"resultJson"`  // 測試結果
	CreatedAt  int64  `json:"createdAt"`   // Unix 秒
}

type AuthTicket struct {
	DocType   string `json:"docType"` // 固定 "AuthTicket"
	PatientH  string `json:"patientHash"`
	TargetH   string `json:"TargetHash"`
	ReportID  string `json:"reportId"`
	GrantedAt int64  `json:"grantedAt"`
	Expiry    int64  `json:"expiry"` // Unix 秒
}

// HealthCheckContract 實作 Chaincode

type HealthCheckContract struct {
	contractapi.Contract
}

func getCaller(ctx contractapi.TransactionContextInterface) (unameHash, role string, err error) {
	id, _ := cid.New(ctx.GetStub())
	unameHash, ok1, _ := id.GetAttributeValue("username")
	role, ok2, _ := id.GetAttributeValue("role")
	if !ok1 || !ok2 {
		err = fmt.Errorf("missing username or role attribute in cert")
	}
	return
}

func getClinicID(ctx contractapi.TransactionContextInterface) string {
	id, _ := cid.New(ctx.GetStub())
	clinic, _, _ := id.GetAttributeValue("clinicId")
	return clinic
}

func recPatientHash(raw []byte) string {
	var t struct {
		PatientH string `json:"patientHash"`
	}
	_ = json.Unmarshal(raw, &t)
	return t.PatientH
}

func nowSec() int64 { return time.Now().Unix() }

// UploadReport 上傳一份健檢報告
func (h *HealthCheckContract) UploadReport(ctx contractapi.TransactionContextInterface, reportID, patientHash, resultJSON string) peer.Response {

	// (A) 僅診所（role=clinic）可呼叫
	if err := cid.AssertAttributeValue(ctx.GetStub(), "role", "clinic"); err != nil {
		return shim.Error("only clinic can upload report")
	}

	// (B) 不允許重複報告 ID
	repKey, _ := ctx.GetStub().CreateCompositeKey(keyReportNS, []string{reportID})
	if b, _ := ctx.GetStub().GetState(repKey); b != nil {
		return shim.Error("reportID already exists")
	}

	// (C) 組裝並寫入
	rec := HealthReport{
		DocType: docHealth, ReportID: reportID, PatientH: patientHash,
		ClinicID: getClinicID(ctx), ResultJSON: resultJSON, CreatedAt: nowSec(),
	}
	bytes, _ := json.Marshal(rec)
	if err := ctx.GetStub().PutState(repKey, bytes); err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)

}

// ReadReport 查詢報告內容（需要在 AccessList 中）
func (h *HealthCheckContract) ReadReport(ctx contractapi.TransactionContextInterface,
	reportID string) peer.Response {

	unameHash, role, err := getCaller(ctx)
	if err != nil {
		return shim.Error(err.Error())
	}

	repKey, _ := ctx.GetStub().CreateCompositeKey(keyReportNS, []string{reportID})
	rb, err := ctx.GetStub().GetState(repKey)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to get report state: %v", err))
	}
	if rb == nil {
		return shim.Error(fmt.Sprintf("report with ID '%s' not found", reportID))
	}

	//  將報告數據 Unmarshal 到 HealthReport 結構體
	var report HealthReport
	if err := json.Unmarshal(rb, &report); err != nil {
		return shim.Error(fmt.Sprintf("failed to unmarshal report data for ID '%s': %v", reportID, err))
	}

	switch role {
	case "clinic":
		return shim.Success(rb)

	case "patient":
		if report.PatientH == unameHash {
			return shim.Success(rb)
		}
		return shim.Error("unauthorized")

	case "insurer":
		ticketKey, _ := ctx.GetStub().CreateCompositeKey(keyAuthNS,
			[]string{report.PatientH, unameHash, reportID})
		tb, _ := ctx.GetStub().GetState(ticketKey)
		if tb == nil {
			return shim.Error("no authorization")
		}
		var tk AuthTicket
		_ = json.Unmarshal(tb, &tk)
		if nowSec() > tk.Expiry {
			return shim.Error("authorization expired")
		}
		return shim.Success(rb)

	default:
		return shim.Error("role not allowed")
	}
}

// GrantAccess 授權其他人查閱報告
func (h *HealthCheckContract) GrantAccess(ctx contractapi.TransactionContextInterface, targetHash, reportID, expiryStr string) peer.Response {

	// (A) 僅患者可呼叫
	patientHash, role, err := getCaller(ctx)
	if err != nil || role != "patient" {
		return shim.Error("only patient can grant access")
	}

	// (B) 檢查到期時間有效
	expiry, errExp := strconv.ParseInt(expiryStr, 10, 64)
	if errExp != nil || expiry <= nowSec() {
		return shim.Error("invalid expiry")
	}

	// (C) 確認報告屬於患者
	repKey, _ := ctx.GetStub().CreateCompositeKey(keyReportNS, []string{reportID})
	rb, _ := ctx.GetStub().GetState(repKey)
	if rb == nil {
		return shim.Error("report not found")
	}
	if recPatientHash(rb) != patientHash {
		return shim.Error("not your report")
	}

	// (D) 寫入授權票
	ticketKey, _ := ctx.GetStub().CreateCompositeKey(keyAuthNS,
		[]string{patientHash, targetHash, reportID})
	tk := AuthTicket{DocType: docAuth,
		PatientH:  patientHash,
		TargetH:   targetHash,
		ReportID:  reportID,
		GrantedAt: nowSec(),
		Expiry:    expiry,
	}

	tbytes, _ := json.Marshal(tk)
	if err := ctx.GetStub().PutState(ticketKey, tbytes); err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)

}

func (h *HealthCheckContract) ListGrantedAccesses(ctx contractapi.TransactionContextInterface) peer.Response {
	targetHash, role, err := getCaller(ctx)
	if err != nil {
		return shim.Error("failed to get caller identity: " + err.Error())
	}
	if role != "insurer" {
		return shim.Error("only insurer can list granted accesses")
	}
	query := fmt.Sprintf(`{
		"selector": {
			"docType": "%s",
			"TargetHash": "%s"
		}`,
		docAuth, targetHash)

	iter, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return shim.Error("query execution failed: " + err.Error())
	}
	defer iter.Close()

	now := nowSec()

	var results []map[string]interface{}
	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			continue // 跳過損壞的記錄
		}

		var tk AuthTicket
		if err := json.Unmarshal(kv.Value, &tk); err != nil {
			continue
		}

		if now > tk.Expiry {
			continue // 跳過過期的授權票
		}

		results = append(results, map[string]interface{}{
			"patientHash": tk.PatientH,
			"reportId":    tk.ReportID,
			"grantedAt":   tk.GrantedAt,
			"expiry":      tk.Expiry,
		})
	}
	bytes, _ := json.Marshal(results)
	return shim.Success(bytes)
}

// ListMyReports 根據使用者身份查詢自己可見的報告清單（不含 ResultJSON）
// 優化點：使用 rich query 篩選 + 僅回傳摘要 + 僅允許 patient/clinic 查詢
func (h *HealthCheckContract) ListMyReports(ctx contractapi.TransactionContextInterface) peer.Response {
	unameHash, role, err := getCaller(ctx)
	if err != nil {
		return shim.Error("failed to get caller identity: " + err.Error())
	}

	var query string

	switch role {
	case "patient":
		// 查詢自己是病患的報告
		query = fmt.Sprintf(`{
			"selector": {
				"docType": "%s",
				"patientHash": "%s"
			}
		}`, docHealth, unameHash)

	case "clinic":
		// 查詢自己診所上傳的報告
		clinicID := unameHash // 診所 ID 目前假設綁在 username（依你的設計）
		query = fmt.Sprintf(`{
			"selector": {
				"docType": "%s",
				"clinicId": "%s"
			}
		}`, docHealth, clinicID)

	default:
		return shim.Error("unauthorized role: only patient or clinic can list their reports")
	}

	iter, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return shim.Error("query execution failed: " + err.Error())
	}
	defer iter.Close()

	var results []map[string]interface{}

	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			continue // 跳過損壞的記錄
		}

		var rep HealthReport
		if err := json.Unmarshal(kv.Value, &rep); err != nil {
			continue
		}

		// 僅回傳摘要，不含 ResultJSON（防洩漏）
		results = append(results, map[string]interface{}{
			"reportId":  rep.ReportID,
			"clinicId":  rep.ClinicID,
			"createdAt": rep.CreatedAt,
		})
	}

	bytes, _ := json.Marshal(results)
	return shim.Success(bytes)
}

// RevokeAccess 收回他人查閱權限
func (h *HealthCheckContract) RevokeAccess(ctx contractapi.TransactionContextInterface, targetHash, reportID string) peer.Response {
	patientHash, role, err := getCaller(ctx)
	if err != nil || role != "patient" {
		return shim.Error("only patient can revoke")
	}
	ticketKey, _ := ctx.GetStub().CreateCompositeKey(keyAuthNS,
		[]string{patientHash, targetHash, reportID})
	if err := ctx.GetStub().DelState(ticketKey); err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
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
