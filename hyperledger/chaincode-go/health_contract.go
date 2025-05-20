package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const (
	docHealth = "HealthRecord"
	docAuth   = "AuthTicket"

	keyReportNS = "REPORT"
	keyAuthNS   = "AUTH"
)

type HealthReport struct {
	DocType     string `json:"docType"`
	ReportID    string `json:"reportId"`
	PatientHash string `json:"patientHash"`
	ClinicID    string `json:"clinicId"`
	ResultJSON  string `json:"resultJson"`
	CreatedAt   int64  `json:"createdAt"`
}

type AuthTicket struct {
	DocType      string `json:"docType"`
	PatientHash  string `json:"patientHash"`
	TargetHash   string `json:"targetHash"`
	ReportID     string `json:"reportId"`
	GrantedAt    int64  `json:"grantedAt"`
	Expiry       int64  `json:"expiry"`
}

type HealthCheckContract struct {
	contractapi.Contract
}

func hashID(id string) string {
	hash := sha256.Sum256([]byte(id))
	return hex.EncodeToString(hash[:])
}

func getCaller(ctx contractapi.TransactionContextInterface) (userID, role string, err error) {
	id, err := cid.New(ctx.GetStub())
	if err != nil {
		return "", "", fmt.Errorf("cannot create client ID: %v", err)
	}

	userID, ok1, _ := id.GetAttributeValue("hf.EnrollmentID")
	role, ok2, _ := id.GetAttributeValue("role")

	if !ok1 || !ok2 {
		err = fmt.Errorf("missing hf.EnrollmentID or role attribute in cert")
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
		PatientHash string `json:"patientHash"`
	}
	_ = json.Unmarshal(raw, &t)
	return t.PatientHash
}

func nowSec() int64 {
	return time.Now().Unix()
}

func (h *HealthCheckContract) UploadReport(ctx contractapi.TransactionContextInterface, reportID, patientHash, resultJSON string) error {
	id, err := cid.New(ctx.GetStub())
	if err != nil {
		return fmt.Errorf("failed to get client identity")
	}
	role, ok, err := id.GetAttributeValue("role")
	if err != nil || !ok || role != "clinic" {
		return fmt.Errorf("only clinic can upload report")
	}

	repKey, _ := ctx.GetStub().CreateCompositeKey(keyReportNS, []string{reportID})
	b, _ := ctx.GetStub().GetState(repKey)
	if b != nil {
		return fmt.Errorf("reportID already exists")
	}

	rec := HealthReport{
		DocType:     docHealth,
		ReportID:    reportID,
		PatientHash: patientHash,
		ClinicID:    getClinicID(ctx),
		ResultJSON:  resultJSON,
		CreatedAt:   nowSec(),
	}
	bytes, _ := json.Marshal(rec)
	return ctx.GetStub().PutState(repKey, bytes)
}

func (h *HealthCheckContract) AuthorizeAccess(ctx contractapi.TransactionContextInterface, reportID, patientHash, targetHash, expiryStr string) error {
	userID, role, err := getCaller(ctx)
	if err != nil || role != "patient" {
		return fmt.Errorf("only patient can grant access")
	}
	expiry, errExp := strconv.ParseInt(expiryStr, 10, 64)
	if errExp != nil || expiry <= nowSec() {
		return fmt.Errorf("invalid expiry")
	}

	// 驗證用戶身份
	callerHash := hashID(userID)
	if callerHash != patientHash {
		return fmt.Errorf("caller identity does not match patientHash")
	}

	repKey, _ := ctx.GetStub().CreateCompositeKey(keyReportNS, []string{reportID})
	rb, _ := ctx.GetStub().GetState(repKey)
	if rb == nil || recPatientHash(rb) != patientHash {
		return fmt.Errorf("not your report")
	}
	ticketKey, _ := ctx.GetStub().CreateCompositeKey(keyAuthNS, []string{patientHash, targetHash, reportID})
	tk := AuthTicket{
		DocType:     docAuth, 
		PatientHash: patientHash, 
		TargetHash:  targetHash, 
		ReportID:    reportID, 
		GrantedAt:   nowSec(), 
		Expiry:      expiry,
	}
	tbytes, _ := json.Marshal(tk)
	return ctx.GetStub().PutState(ticketKey, tbytes)
}

func (h *HealthCheckContract) RevokeAccess(ctx contractapi.TransactionContextInterface, targetHash, reportID string) error {
	userID, role, err := getCaller(ctx)
	if err != nil || role != "patient" {
		return fmt.Errorf("only patient can revoke")
	}

	patientHash := hashID(userID)

	ticketKey, _ := ctx.GetStub().CreateCompositeKey(keyAuthNS, []string{patientHash, targetHash, reportID})
	return ctx.GetStub().DelState(ticketKey)
}

func (h *HealthCheckContract) ListMyReports(ctx contractapi.TransactionContextInterface) ([]HealthReport, error) {
	userID, role, err := getCaller(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get caller identity: %v", err)
	}
	if role != "patient" {
		return nil, fmt.Errorf("only patient can list their reports")
	}

	patientHash := hashID(userID)

	query := fmt.Sprintf(`{
		"selector": {
			"docType": "%s",
			"patientHash": "%s"
		}
	}`, docHealth, patientHash)

	iter, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %v", err)
	}
	defer iter.Close()

	var results []HealthReport

	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			continue
		}

		var rep HealthReport
		if err := json.Unmarshal(kv.Value, &rep); err != nil {
			continue
		}

		results = append(results, rep)
	}

	return results, nil
}

func (h *HealthCheckContract) ListAuthorizedReports(ctx contractapi.TransactionContextInterface) ([]HealthReport, error) {
	userID, role, err := getCaller(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get caller identity: %v", err)
	}
	if role != "insurer" {
		return nil, fmt.Errorf("only insurer can list authorized reports")
	}

	targetHash := hashID(userID)

	query := fmt.Sprintf(`{
		"selector": {
			"docType": "%s",
			"targetHash": "%s"
		}
	}`, docAuth, targetHash)

	iter, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query authorization tickets: %v", err)
	}
	defer iter.Close()

	now := nowSec()
	var results []HealthReport

	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			continue
		}

		var tk AuthTicket
		if err := json.Unmarshal(kv.Value, &tk); err != nil {
			continue
		}
		if now > tk.Expiry {
			continue // 忽略過期授權
		}

		// 查詢對應報告
		repKey, _ := ctx.GetStub().CreateCompositeKey(keyReportNS, []string{tk.ReportID})
		rb, err := ctx.GetStub().GetState(repKey)
		if err != nil || rb == nil {
			continue
		}

		var rep HealthReport
		if err := json.Unmarshal(rb, &rep); err != nil {
			continue
		}

		results = append(results, rep)
	}

	return results, nil
}

// 不含健檢數據的報告元數據
type ReportMeta struct {
	ReportID  string `json:"reportId"`
	ClinicID  string `json:"clinicId"`
	CreatedAt int64  `json:"createdAt"`
}

func (h *HealthCheckContract) ListReportMetaByPatientID(ctx contractapi.TransactionContextInterface, patientID string) ([]ReportMeta, error) {
	_  , role, err := getCaller(ctx)
	if err != nil || role != "insurer" {
		return nil, fmt.Errorf("only insurer can query report meta")
	}

	patientHash := hashID(patientID)

	query := fmt.Sprintf(`{
		"selector": {
			"docType": "%s",
			"patientHash": "%s"
		}
	}`, docHealth, patientHash)

	iter, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %v", err)
	}
	defer iter.Close()

	var results []ReportMeta
	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			continue
		}
		var report HealthReport
		if err := json.Unmarshal(kv.Value, &report); err != nil {
			continue
		}

		results = append(results, ReportMeta{
			ReportID:  report.ReportID,
			ClinicID:  report.ClinicID,
			CreatedAt: report.CreatedAt,
		})
	}

	return results, nil
}

// 保險業者讀取授權報告
func (h *HealthCheckContract) ReadAuthorizedReport(ctx contractapi.TransactionContextInterface, reportID string) (string, error) {
	userID, role, err := getCaller(ctx)
	if err != nil || role != "insurer" {
		return "", fmt.Errorf("only insurer can read report")
	}

	targetHash := hashID(userID)

	// 查詢授權票
	iter, err := ctx.GetStub().GetStateByPartialCompositeKey(keyAuthNS, []string{"", targetHash, reportID})
	if err != nil {
		return "", fmt.Errorf("error accessing authorization ticket")
	}
	defer iter.Close()

	found := false
	var tk AuthTicket
	now := nowSec()
	for iter.HasNext() {
		kv, _ := iter.Next()
		if err := json.Unmarshal(kv.Value, &tk); err != nil {
			continue
		}
		if tk.ReportID == reportID && now <= tk.Expiry {
			found = true
			break
		}
	}

	if !found {
		return "", fmt.Errorf("access denied or expired for report %s", reportID)
	}

	// 查詢報告內容
	repKey, _ := ctx.GetStub().CreateCompositeKey(keyReportNS, []string{reportID})
	data, err := ctx.GetStub().GetState(repKey)
	if err != nil || data == nil {
		return "", fmt.Errorf("report not found")
	}

	var rep HealthReport
	if err := json.Unmarshal(data, &rep); err != nil {
		return "", fmt.Errorf("failed to parse report")
	}

	return rep.ResultJSON, nil
}


func main() {
	chaincode, err := contractapi.NewChaincode(&HealthCheckContract{})
	if err != nil {
		panic(fmt.Sprintf("Error creating chaincode: %v", err))
	}
	if err := chaincode.Start(); err != nil {
		panic(fmt.Sprintf("Error starting chaincode: %v", err))
	}
}
