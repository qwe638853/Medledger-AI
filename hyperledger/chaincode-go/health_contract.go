package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type HealthReport struct {
	ReportID      string            `json:"reportID"`
	PatientID     string            `json:"patientID"`
	EncryptedData string            `json:"encryptedData"`
	Flags         map[string]int    `json:"flags"`
	SummaryHash   string            `json:"summaryHash"`
	AccessList    []string          `json:"accessList"`
	Timestamp     string            `json:"timestamp"`
}

type HealthCheckContract struct {
	contractapi.Contract
}

func (h *HealthCheckContract) UploadReport(ctx contractapi.TransactionContextInterface, reportID, patientID, encryptedData, flagsJson, summaryHash string) error {
	exists, err := ctx.GetStub().GetState(reportID)
	if err != nil {
		return err
	}
	if exists != nil {
		return fmt.Errorf("Report %s already exists", reportID)
	}

	var flags map[string]int
	err = json.Unmarshal([]byte(flagsJson), &flags)
	if err != nil {
		return err
	}

	report := HealthReport{
		ReportID:      reportID,
		PatientID:     patientID,
		EncryptedData: encryptedData,
		Flags:         flags,
		SummaryHash:   summaryHash,
		AccessList:    []string{},
		Timestamp:     time.Now().Format(time.RFC3339),
	}

	reportBytes, err := json.Marshal(report)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(reportID, reportBytes)
}

func (h *HealthCheckContract) GrantAccess(ctx contractapi.TransactionContextInterface, reportID, accessorID string) error {
	reportBytes, err := ctx.GetStub().GetState(reportID)
	if err != nil {
		return err
	}
	if reportBytes == nil {
		return fmt.Errorf("Report %s not found", reportID)
	}

	var report HealthReport
	err = json.Unmarshal(reportBytes, &report)
	if err != nil {
		return err
	}

	for _, id := range report.AccessList {
		if id == accessorID {
			return nil // already granted
		}
	}
	report.AccessList = append(report.AccessList, accessorID)

	reportBytes, err = json.Marshal(report)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(reportID, reportBytes)
}

func (h *HealthCheckContract) RevokeAccess(ctx contractapi.TransactionContextInterface, reportID, accessorID string) error {
	reportBytes, err := ctx.GetStub().GetState(reportID)
	if err != nil {
		return err
	}
	if reportBytes == nil {
		return fmt.Errorf("Report %s not found", reportID)
	}

	var report HealthReport
	err = json.Unmarshal(reportBytes, &report)
	if err != nil {
		return err
	}

	newList := []string{}
	for _, id := range report.AccessList {
		if id != accessorID {
			newList = append(newList, id)
		}
	}
	report.AccessList = newList

	reportBytes, err = json.Marshal(report)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(reportID, reportBytes)
}

func (h *HealthCheckContract) ReadReport(ctx contractapi.TransactionContextInterface, reportID, requesterID string) (*HealthReport, error) {
	reportBytes, err := ctx.GetStub().GetState(reportID)
	if err != nil {
		return nil, err
	}
	if reportBytes == nil {
		return nil, fmt.Errorf("Report %s not found", reportID)
	}

	var report HealthReport
	err = json.Unmarshal(reportBytes, &report)
	if err != nil {
		return nil, err
	}

	authorized := false
	for _, id := range report.AccessList {
		if id == requesterID {
			authorized = true
			break
		}
	}

	if !authorized {
		return nil, fmt.Errorf("Access denied for %s", requesterID)
	}

	return &report, nil
}

func (h *HealthCheckContract) GetSummaryHash(ctx contractapi.TransactionContextInterface, reportID string) (string, error) {
	reportBytes, err := ctx.GetStub().GetState(reportID)
	if err != nil {
		return "", err
	}
	if reportBytes == nil {
		return "", fmt.Errorf("Report %s not found", reportID)
	}

	var report HealthReport
	err = json.Unmarshal(reportBytes, &report)
	if err != nil {
		return "", err
	}

	return report.SummaryHash, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(HealthCheckContract))
	if err != nil {
		panic(fmt.Sprintf("Error creating chaincode: %v", err))
	}

	if err := chaincode.Start(); err != nil {
		panic(fmt.Sprintf("Error starting chaincode: %v", err))
	}
}
