package main

import (
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
	DocType    string `json:"docType"`
	ReportID   string `json:"reportId"`
	PatientH   string `json:"patientHash"`
	ClinicID   string `json:"clinicId"`
	ResultJSON string `json:"resultJson"`
	CreatedAt  int64  `json:"createdAt"`
}

type AuthTicket struct {
	DocType   string `json:"docType"`
	PatientH  string `json:"patientHash"`
	TargetH   string `json:"TargetHash"`
	ReportID  string `json:"reportId"`
	GrantedAt int64  `json:"grantedAt"`
	Expiry    int64  `json:"expiry"`
}

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
		DocType: docHealth,
		ReportID: reportID,
		PatientH: patientHash,
		ClinicID: getClinicID(ctx),
		ResultJSON: resultJSON,
		CreatedAt: nowSec(),
	}
	bytes, _ := json.Marshal(rec)
	return ctx.GetStub().PutState(repKey, bytes)
}

func (h *HealthCheckContract) GrantAccess(ctx contractapi.TransactionContextInterface, targetHash, reportID, expiryStr string) error {
	patientHash, role, err := getCaller(ctx)
	if err != nil || role != "patient" {
		return fmt.Errorf("only patient can grant access")
	}
	expiry, errExp := strconv.ParseInt(expiryStr, 10, 64)
	if errExp != nil || expiry <= nowSec() {
		return fmt.Errorf("invalid expiry")
	}

	repKey, _ := ctx.GetStub().CreateCompositeKey(keyReportNS, []string{reportID})
	rb, _ := ctx.GetStub().GetState(repKey)
	if rb == nil || recPatientHash(rb) != patientHash {
		return fmt.Errorf("not your report")
	}
	ticketKey, _ := ctx.GetStub().CreateCompositeKey(keyAuthNS, []string{patientHash, targetHash, reportID})
	tk := AuthTicket{DocType: docAuth, PatientH: patientHash, TargetH: targetHash, ReportID: reportID, GrantedAt: nowSec(), Expiry: expiry}
	tbytes, _ := json.Marshal(tk)
	return ctx.GetStub().PutState(ticketKey, tbytes)
}

func (h *HealthCheckContract) RevokeAccess(ctx contractapi.TransactionContextInterface, targetHash, reportID string) error {
	patientHash, role, err := getCaller(ctx)
	if err != nil || role != "patient" {
		return fmt.Errorf("only patient can revoke")
	}
	ticketKey, _ := ctx.GetStub().CreateCompositeKey(keyAuthNS, []string{patientHash, targetHash, reportID})
	return ctx.GetStub().DelState(ticketKey)
}

func (h *HealthCheckContract) ListMyReports(ctx contractapi.TransactionContextInterface) ([]HealthReport, error) {
	unameHash, role, err := getCaller(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get caller identity: %v", err)
	}
	if role != "patient" {
		return nil, fmt.Errorf("only patient can list their reports")
	}

	query := fmt.Sprintf(`{
		"selector": {
			"docType": "%s",
			"patientHash": "%s"
		}
	}`, docHealth, unameHash)

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
	targetHash, role, err := getCaller(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get caller identity: %v", err)
	}
	if role != "insurer" {
		return nil, fmt.Errorf("only insurer can list authorized reports")
	}

	query := fmt.Sprintf(`{
		"selector": {
			"docType": "%s",
			"TargetHash": "%s"
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



func main() {
	chaincode, err := contractapi.NewChaincode(&HealthCheckContract{})
	if err != nil {
		panic(fmt.Sprintf("Error creating chaincode: %v", err))
	}
	if err := chaincode.Start(); err != nil {
		panic(fmt.Sprintf("Error starting chaincode: %v", err))
	}
}
