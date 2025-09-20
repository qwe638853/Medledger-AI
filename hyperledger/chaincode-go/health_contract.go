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
	docAccessRequest = "AccessRequest"
	keyReportNS = "REPORT"
	keyAuthNS   = "AUTH"
	keyAccessRequestNS = "ACCESS_REQUEST"
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

type AccessRequest struct {
	DocType      string `json:"docType"`
	RequestID    string `json:"requestId"`
	ReportID     string `json:"reportId"`
	PatientHash  string `json:"patientHash"`
	RequesterHash string `json:"requesterHash"`
	Reason       string `json:"reason"`
	RequestedAt  int64  `json:"requestedAt"`
	Expiry       int64  `json:"expiry"`
	Status       string `json:"status"`
}

// 只包含 metadata
type ReportMeta struct {
	ReportID  string `json:"reportId"`
	ClinicID  string `json:"clinicId"`
	CreatedAt int64  `json:"createdAt"`
}

type HealthCheckContract struct {
	contractapi.Contract
}

/**
 * @notice 將身分證字號轉換為SHA256雜湊值
 * @dev 使用SHA256演算法對身分證字號進行單向雜湊，保護隱私
 * @param id 身分證字號或任何需要雜湊的字串
 * @return string 16進位編碼的雜湊值
 */
func hashID(id string) string {
	hash := sha256.Sum256([]byte(id))
	return hex.EncodeToString(hash[:])
}

/**
 * @notice 取得調用者身分資訊
 * @dev 從Fabric憑證中提取用戶ID和角色，用於權限驗證
 * @param ctx Fabric合約上下文
 * @return userID 用戶ID, role 角色, err 錯誤訊息
 */
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

/**
 * @notice 取得診所ID
 * @dev 從Fabric憑證中提取診所ID屬性
 * @param ctx Fabric合約上下文
 * @return string 診所ID
 */
func getClinicID(ctx contractapi.TransactionContextInterface) string {
	id, _ := cid.New(ctx.GetStub())
	clinic, _, _ := id.GetAttributeValue("clinicId")
	return clinic
}

/**
 * @notice 從原始位元組資料中提取病患雜湊值
 * @dev 解析JSON格式的位元組資料，提取patientHash欄位
 * @param raw 原始位元組資料
 * @return string 病患雜湊值
 */
func recPatientHash(raw []byte) string {
	var t struct {
		PatientHash string `json:"patientHash"`
	}
	_ = json.Unmarshal(raw, &t)
	return t.PatientHash
}

/**
 * @notice 取得當前Unix時間戳（秒）
 * @dev 取得目前時間的Unix時間戳，用於記錄時間
 * @return int64 Unix時間戳（秒）
 */
func nowSec() int64 {
	return time.Now().Unix()
}

/**
 * @notice 上傳健檢報告到區塊鏈
 * @dev 只允許 clinic 身份上傳，會驗證reportID唯一性並記錄時間戳
 * @param ctx Fabric合約上下文
 * @param reportID 報告唯一識別ID
 * @param patientHash 病患身份雜湊值
 * @param resultJSON 健檢結果JSON字串
 * @return error 上傳失敗或權限錯誤
 */
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

func (h *HealthCheckContract) UpdateReport(ctx contractapi.TransactionContextInterface, reportID, newResultJSON string) error {
    userID, role, err := getCaller(ctx)
    if err != nil { return fmt.Errorf("failed to get caller identity") }
    if role != "patient" { return fmt.Errorf("only patient can update report") }

    repKey, _ := ctx.GetStub().CreateCompositeKey(keyReportNS, []string{reportID})
    b, _ := ctx.GetStub().GetState(repKey)
    if b == nil { return fmt.Errorf("report not found") }

    var rec HealthReport
    if err := json.Unmarshal(b, &rec); err != nil { return fmt.Errorf("failed to unmarshal report: %v", err) }
    // 歸屬檢查
    if rec.PatientHash != hashID(userID) { return fmt.Errorf("not owner of this report") }

    // 解析舊、新 envelope
    var oldEnv, newEnv struct {
        Ciphertext  string            `json:"ciphertext"`
        Nonce       string            `json:"nonce"`
        WrappedKeys map[string]string `json:"wrappedKeys"`
        Enc         string            `json:"enc"`
        KDF         string            `json:"kdf"`
        Curve       string            `json:"curve"`
    }
    if err := json.Unmarshal([]byte(rec.ResultJSON), &oldEnv); err != nil {
        return fmt.Errorf("bad old envelope: %v", err)
    }
    if err := json.Unmarshal([]byte(newResultJSON), &newEnv); err != nil {
        return fmt.Errorf("bad new envelope: %v", err)
    }

    // 嚴格一致性檢查：不允許改密文/nonce/演算法/版本/AAD
    if oldEnv.Ciphertext != newEnv.Ciphertext ||
       oldEnv.Nonce != newEnv.Nonce ||
       oldEnv.Enc != newEnv.Enc ||
       oldEnv.KDF != newEnv.KDF ||
       oldEnv.Curve != newEnv.Curve || {
        return fmt.Errorf("only wrappedKeys can be changed")
    }

    // 若需要：可以改成「merge」模式，防止把既有 key 覆蓋掉
    if oldEnv.WrappedKeys == nil { oldEnv.WrappedKeys = map[string]string{} }
    for k, v := range newEnv.WrappedKeys {
        oldEnv.WrappedKeys[k] = v
    }

    // 回寫
    merged, _ := json.Marshal(oldEnv)
    rec.ResultJSON = string(merged)
    out, _ := json.Marshal(rec)
    return ctx.GetStub().PutState(repKey, out)
}


/**
 * @notice 病患批准並授權訪問請求
 * @dev 只允許 patient 身份，會產生授權票據並發送事件
 * @param ctx Fabric合約上下文
 * @param requestID 訪問請求ID
 * @return error 批准失敗或權限錯誤
 */
func (h *HealthCheckContract) ApproveAndAuthorizeAccess(ctx contractapi.TransactionContextInterface, requestID string) error {
    userID, role, err := getCaller(ctx)
    if err != nil || role != "patient" {
        return fmt.Errorf("only patient can approve")
    }

    // 1. 取得請求並驗證
    reqKey, _ := ctx.GetStub().CreateCompositeKey(keyAccessRequestNS, []string{requestID})
    reqBytes, err := ctx.GetStub().GetState(reqKey)
    if err != nil || reqBytes == nil {
        return fmt.Errorf("request not found")
    }
    var req AccessRequest
    if err := json.Unmarshal(reqBytes, &req); err != nil {
        return fmt.Errorf("failed to unmarshal request: %v", err)
    }
    patientHash := hashID(userID)
    if req.PatientHash != patientHash {
        return fmt.Errorf("not authorized to approve this request")
    }
    if req.Status != "PENDING" {
        return fmt.Errorf("request already handled")
    }

    // 2. 更新狀態
    req.Status = "APPROVED"
    newReqBytes, _ := json.Marshal(req)
    if err := ctx.GetStub().PutState(reqKey, newReqBytes); err != nil {
        return fmt.Errorf("failed to update request status")
    }

    // 3. 產生授權票據
    ticketKey, _ := ctx.GetStub().CreateCompositeKey(keyAuthNS, []string{req.PatientHash, req.RequesterHash, req.ReportID})
    tk := AuthTicket{
        DocType:     docAuth,
        PatientHash: req.PatientHash,
        TargetHash:  req.RequesterHash,
        ReportID:    req.ReportID,
        GrantedAt:   nowSec(),
        Expiry:      req.Expiry,
    }
    tbytes, _ := json.Marshal(tk)
	if err := ctx.GetStub().PutState(ticketKey, tbytes); err != nil {
        return fmt.Errorf("failed to store auth ticket")
    }

	eventPayload, _ := json.Marshal(map[string]interface{}{
        "requestId":    req.RequestID,
        "reportId":     req.ReportID,
        "patientHash":  req.PatientHash,
        "requesterHash": req.RequesterHash,
        "status":       "APPROVED",
        "grantedAt":    nowSec(),
        "expiry":       req.Expiry,
    })
    if err := ctx.GetStub().SetEvent("AccessApproved", eventPayload); err != nil {
        return fmt.Errorf("failed to set event")
    }

    return nil
}

/**
 * @notice 病患拒絕授權請求
 * @dev 只允許 patient 身份，更新請求狀態為拒絕並發送事件
 * @param ctx Fabric合約上下文
 * @param requestID 請求ID
 * @return error 拒絕失敗或權限錯誤
 */
func (h *HealthCheckContract) RejectAccessRequest(ctx contractapi.TransactionContextInterface, requestID string) error {
    userID, role, err := getCaller(ctx)
    if err != nil || role != "patient" {
        return fmt.Errorf("only patient can reject access request")
    }

    // 1. 取得請求並驗證
    reqKey, _ := ctx.GetStub().CreateCompositeKey(keyAccessRequestNS, []string{requestID})
    reqBytes, err := ctx.GetStub().GetState(reqKey)
    if err != nil || reqBytes == nil {
        return fmt.Errorf("request not found")
    }
    
    var req AccessRequest
    if err := json.Unmarshal(reqBytes, &req); err != nil {
        return fmt.Errorf("failed to unmarshal request: %v", err)
    }
    
    patientHash := hashID(userID)
    if req.PatientHash != patientHash {
        return fmt.Errorf("not authorized to reject this request")
    }
    
    if req.Status != "PENDING" {
        return fmt.Errorf("request already handled")
    }

    // 2. 更新狀態為拒絕
    req.Status = "REJECTED"
    newReqBytes, _ := json.Marshal(req)
    if err := ctx.GetStub().PutState(reqKey, newReqBytes); err != nil {
        return fmt.Errorf("failed to update request status")
    }

    // 3. 發送拒絕事件
    eventPayload, _ := json.Marshal(map[string]interface{}{
        "requestId":    req.RequestID,
        "reportId":     req.ReportID,
        "patientHash":  req.PatientHash,
        "requesterHash": req.RequesterHash,
        "status":       "REJECTED",
        "rejectedAt":   nowSec(),
    })
    if err := ctx.GetStub().SetEvent("AccessRejected", eventPayload); err != nil {
        return fmt.Errorf("failed to set event")
    }

    return nil
}


/**
 * @notice 查詢自己的所有健檢報告，只回傳 metadata，不含內容
 * @dev 只允許 patient 身份，回傳該 patient 的所有報告 meta 資訊
 * @param ctx Fabric合約上下文
 * @return []ReportMeta 報告meta陣列, error 查詢失敗或權限錯誤
 */
func (h *HealthCheckContract) ListMyReportMeta(ctx contractapi.TransactionContextInterface) ([]ReportMeta, error) {
	userID, role, err := getCaller(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get caller identity: %v", err)
	}
	if role != "patient" {
		return nil, fmt.Errorf("only patient can list their report meta")
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

/**
 * @notice 病患讀取自己的報告詳細內容（需先取得 meta 再查詢內容）
 * @dev 只允許 patient 身份，只能讀自己擁有的報告
 * @param ctx Fabric合約上下文
 * @param reportID 報告ID
 * @return string 報告內容 JSON, error 查詢失敗或無權限
 */
func (h *HealthCheckContract) ReadMyReport(ctx contractapi.TransactionContextInterface, reportID string) (string, error) {
	userID, role, err := getCaller(ctx)
	if err != nil || role != "patient" {
		return "", fmt.Errorf("only patient can read own report")
	}
	patientHash := hashID(userID)
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
	if rep.PatientHash != patientHash {
		return "", fmt.Errorf("not authorized to read this report")
	}
	return rep.ResultJSON, nil
}


/**
 * @notice 查詢特定病患的所有健檢報告 meta，只回傳 meta，不含內容
 * @dev 只允許 insurer 身份，需傳入病患ID
 * @param ctx Fabric合約上下文
 * @param patientID 病患身份證或ID
 * @return []ReportMeta 報告meta陣列, error 查詢失敗或權限錯誤
 */
func (h *HealthCheckContract) ListReportMetaByPatientID(ctx contractapi.TransactionContextInterface, patientID string) ([]ReportMeta, error) {
	_, role, err := getCaller(ctx)
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

/**
 * @notice 保險業者讀取已授權的健檢報告內容
 * @dev 只允許 insurer 身份，且必須獲得授權票據，且票據未過期
 * @param ctx Fabric合約上下文
 * @param patientHash 病患hash
 * @param reportID 報告ID
 * @return string 報告內容 JSON, error 查詢失敗或無權限
 */
func (h *HealthCheckContract) ReadAuthorizedReport(ctx contractapi.TransactionContextInterface, patientHash, reportID string) (string, error) {
	userID, role, err := getCaller(ctx)
	if err != nil || role != "insurer" {
		return "", fmt.Errorf("only insurer can read report")
	}
	targetHash := hashID(userID)

	// 查詢授權票
	iter, err := ctx.GetStub().GetStateByPartialCompositeKey(keyAuthNS, []string{patientHash, targetHash, reportID})
	if err != nil {
		return "", fmt.Errorf("error accessing authorization ticket")
	}
	defer iter.Close()

	found := false
	expired := false
	var tk AuthTicket
	now := nowSec()
	for iter.HasNext() {
		kv, _ := iter.Next()
		if err := json.Unmarshal(kv.Value, &tk); err != nil {
			continue
		}
		if tk.ReportID == reportID {
			if now <= tk.Expiry {
				found = true
				break
			} else {
				expired = true
			}
		}
	}
	if !found {
		if expired {
			return "", fmt.Errorf("access expired for report %s", reportID)
		}
		return "", fmt.Errorf("access denied for report %s", reportID)
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


/**
 * @notice 病患查看自己已授權的票據列表
 * @dev 只允許 patient 身份，回傳該病患發出的所有授權票據
 * @param ctx Fabric合約上下文
 * @return []AuthTicket 授權票據陣列, error 查詢失敗或權限錯誤
 */
func (h *HealthCheckContract) ListMyAuthorizedTickets(ctx contractapi.TransactionContextInterface) ([]AuthTicket, error) {
	userID, role, err := getCaller(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get caller identity: %v", err)
	}
	if role != "patient" {
		return nil, fmt.Errorf("only patient can list their authorized reports")
	}

	patientHash := hashID(userID)

	query := fmt.Sprintf(`{
		"selector": {
			"docType": "%s",
			"patientHash": "%s"
		}
	}`, docAuth, patientHash)
	
	iter, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %v", err)
	}
	defer iter.Close()

	var results []AuthTicket
	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			continue
		}
		var tk AuthTicket
		if err := json.Unmarshal(kv.Value, &tk); err != nil {
			continue
		}
		results = append(results, tk)
	}
	return results, nil
}
	
/**
 * @notice 保險業者查看所有已授權的報告
 * @dev 只允許 insurer 身份，回傳所有授權給該保險業者且未過期的報告
 * @param ctx Fabric合約上下文
 * @return []map[string]interface{} 已授權報告資訊陣列, error 查詢失敗或權限錯誤
 */
func (h *HealthCheckContract) ListAuthorizedReports(ctx contractapi.TransactionContextInterface) ([]map[string]interface{}, error) {
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
	var results []map[string]interface{}

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

		// 組合報告和授權信息
		result := map[string]interface{}{
			"reportId":    rep.ReportID,
			"clinicId":    rep.ClinicID,
			"patientHash": rep.PatientHash,
			"resultJson":  rep.ResultJSON,
			"createdAt":   rep.CreatedAt,
			"expiry":      tk.Expiry,
		}
		results = append(results, result)
	}

	return results, nil
}



/**
 * @notice 保險業者請求訪問健檢報告
 * @dev 只允許 insurer 身份，會建立待處理的訪問請求
 * @param ctx Fabric合約上下文
 * @param reportID 報告ID
 * @param patientHash 病患雜湊值
 * @param reason 請求原因
 * @param expiryStr 過期時間戳字串
 * @return error 請求失敗或權限錯誤
 */
func (h *HealthCheckContract) RequestAccess(ctx contractapi.TransactionContextInterface, reportID, patientHash, reason, expiryStr string) error {
	userID, role, err := getCaller(ctx)
	if err != nil || role != "insurer" {
		return fmt.Errorf("only insurer can request access")
	}

	requesterHash := hashID(userID)
	requestID := fmt.Sprintf("req_%d", nowSec())

	// 檢查報告是否存在
	repKey, _ := ctx.GetStub().CreateCompositeKey(keyReportNS, []string{reportID})
	rb, _ := ctx.GetStub().GetState(repKey)
	if rb == nil {
		return fmt.Errorf("report not found")
	}

	expiry, errExp := strconv.ParseInt(expiryStr, 10, 64)
	if errExp != nil || expiry <= nowSec() {
		return fmt.Errorf("invalid expiry")
	}

	reqKey, _ := ctx.GetStub().CreateCompositeKey(keyAccessRequestNS, []string{requestID})
	req := AccessRequest{
		DocType:       docAccessRequest,
		RequestID:     requestID,
		ReportID:      reportID,
		PatientHash:   patientHash,
		RequesterHash: requesterHash,
		Reason:        reason,
		RequestedAt:   nowSec(),
		Expiry:        expiry,
		Status:        "PENDING",
	}

	reqBytes, _ := json.Marshal(req)
	return ctx.GetStub().PutState(reqKey, reqBytes)
}

/**
 * @notice 病患列出待處理的授權請求
 * @dev 只允許 patient 身份，回傳所有狀態為PENDING的請求
 * @param ctx Fabric合約上下文
 * @return []AccessRequest 待處理請求陣列, error 查詢失敗或權限錯誤
 */
func (h *HealthCheckContract) ListPendingAccessRequests(ctx contractapi.TransactionContextInterface) ([]AccessRequest, error) {
	userID, role, err := getCaller(ctx)
	if err != nil || role != "patient" {
		return nil, fmt.Errorf("only patient can list pending requests")
	}

	patientHash := hashID(userID)
	query := fmt.Sprintf(`{
		"selector": {
			"docType": "%s",
			"patientHash": "%s",
			"status": "PENDING"
		}
	}`, docAccessRequest, patientHash)

	iter, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %v", err)
	}
	defer iter.Close()

	var results []AccessRequest
	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			continue
		}

		var req AccessRequest
		if err := json.Unmarshal(kv.Value, &req); err != nil {
			continue
		}
		results = append(results, req)
	}

	return results, nil
}

/**
 * @notice 保險業者列出自己發出的授權請求
 * @dev 只允許 insurer 身份，回傳該保險業者發出的所有請求
 * @param ctx Fabric合約上下文
 * @return []AccessRequest 請求陣列, error 查詢失敗或權限錯誤
 */
func (h *HealthCheckContract) ListMyAccessRequests(ctx contractapi.TransactionContextInterface) ([]AccessRequest, error) {
    userID, role, err := getCaller(ctx)
    if err != nil || role != "insurer" {
        return nil, fmt.Errorf("only insurer can list their access requests")
    }
    requesterHash := hashID(userID)

    query := fmt.Sprintf(`{
        "selector": {
            "docType": "%s",
            "requesterHash": "%s"
        }
    }`, docAccessRequest, requesterHash)

    iter, err := ctx.GetStub().GetQueryResult(query)
    if err != nil {
        return nil, fmt.Errorf("query execution failed: %v", err)
    }
    defer iter.Close()

    var results []AccessRequest
    for iter.HasNext() {
        kv, err := iter.Next()
        if err != nil {
            continue
        }
        var req AccessRequest
        if err := json.Unmarshal(kv.Value, &req); err != nil {
            continue
        }
        results = append(results, req)
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
