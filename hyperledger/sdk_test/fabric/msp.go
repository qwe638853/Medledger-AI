package fabric

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/hyperledger/fabric-ca/api"
	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/pkg/errors"
)

type EnrollRequest struct {
	Profile             string `json:"profile,omitempty"`
	Certificate_request string `json:"certificate_request,omitempty"`
}

// 將 byte 資料轉換為 Base64 編碼字串，用於 token payload 組合
func B64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// 從完整 URL 中取出 URI 路徑（不含主機與查詢字串）
func uriPath(full string) string {
	u, err := url.Parse(full)
	if err != nil {
		return full
	}
	return u.RequestURI()
}

// 建立 ECDSA token，用來在與 Fabric CA 溝通時進行授權驗證（Authorization header）
func GenECDSAToken(csp bccsp.BCCSP, cert []byte, key bccsp.Key, method, uri string, body []byte) (string, error) {
	b64body := B64Encode(body)
	b64cert := B64Encode(cert)
	b64uri := B64Encode([]byte(uriPath(uri)))
	payload := method + "." + b64uri + "." + b64body + "." + b64cert
	return genECDSAToken(csp, key, b64cert, payload)
}

// 實際進行雜湊與簽章，用於支援 GenECDSAToken
func genECDSAToken(csp bccsp.BCCSP, key bccsp.Key, b64cert, payload string) (string, error) {
	digest, err := csp.Hash([]byte(payload), &bccsp.SHAOpts{})
	if err != nil {
		return "", errors.WithMessage(err, fmt.Sprintf("Hash failed on '%s'", payload))
	}
	sig, err := csp.Sign(key, digest, nil)
	if err != nil {
		return "", errors.WithMessage(err, "BCCSP signature generation failure")
	}
	if len(sig) == 0 {
		return "", errors.New("BCCSP signature creation failed. Signature must be different than nil")
	}
	b64sig := B64Encode(sig)
	return b64cert + "." + b64sig, nil
}

// 使用 Fabric 官方 api.RegistrationRequest 結構向 Fabric CA 註冊身份
func RegisterUser(caURL, certPath, keyPath string, req api.RegistrationRequest) error {
	// 將註冊資料轉為 JSON
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("❌ Failed to marshal request: %w", err)
	}

	// 讀取 admin 的證書與私鑰，用於產生 token
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return fmt.Errorf("❌ Failed to read cert: %w", err)
	}
	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return fmt.Errorf("❌ Failed to read key: %w", err)
	}

	// 初始化 BCCSP（密鑰與簽章模組）
	factory.InitFactories(nil)
	csp := factory.GetDefault()

	// 匯入私鑰以供 BCCSP 使用
	keyBlock, _ := pem.Decode(keyPEM)
	key, err := csp.KeyImport(keyBlock.Bytes, &bccsp.ECDSAPrivateKeyImportOpts{Temporary: true})
	if err != nil {
		return fmt.Errorf("❌ Failed to import key: %w", err)
	}

	// 產生 token
	token, err := GenECDSAToken(csp, certPEM, key, "POST", caURL+"/api/v1/register", bodyBytes)
	if err != nil {
		return fmt.Errorf("❌ Failed to generate token: %w", err)
	}

	// 建立 HTTP request 發送註冊資料
	httpReq, err := http.NewRequest("POST", caURL+"/api/v1/register", bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("❌ Failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", token)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("❌ Failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 檢查回應是否成功
	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return fmt.Errorf("❌ Register failed (%d): %s", resp.StatusCode, respBody)
	}

	fmt.Printf("✅ Register success: %s\n", respBody)
	return nil
}

func EnrollUser(caURL, enrollID, enrollSecret string, enrollRequest EnrollRequest) ([]byte, error) {

	var enrollResp struct {
		Result struct {
			Cert string `json:"Cert"`
		} `json:"result"`
	}
	bodyBytes, err := json.Marshal(enrollRequest)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to marshal empty body: %w", err)
	}
	// Basic Auth header: base64("id:secret")
	authStr := base64.StdEncoding.EncodeToString([]byte(enrollID + ":" + enrollSecret))

	httpReq, err := http.NewRequest("POST", caURL+"/api/v1/enroll", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to create enroll request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Basic "+authStr)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to send enroll request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, fmt.Errorf("❌ Enroll failed (%d): %s", resp.StatusCode, respBody)
	}

	fmt.Printf("✅ Enroll success: %s\n", respBody)

	json.Unmarshal(respBody, &enrollResp)

	certPEM, err := base64.StdEncoding.DecodeString(enrollResp.Result.Cert)
	if err != nil {
		return nil, fmt.Errorf("無法解碼憑證: %v", err)
	}

	return certPEM, nil
}

// 建立使用者專屬 CSR（含 CommonName）與私鑰
func GenerateCSR(commonName string) (*ecdsa.PrivateKey, []byte, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	template := x509.CertificateRequest{
		Subject: pkix.Name{CommonName: commonName},
	}
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &template, priv)
	if err != nil {
		return nil, nil, err
	}
	csrPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrDER})
	return priv, csrPEM, nil
}

// 儲存私鑰為 PEM 格式至指定路徑
func SavePrivateKeyToFile(key *ecdsa.PrivateKey, filename string) error {
	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER})
	return ioutil.WriteFile(filename, keyPEM, 0600)
}

// 將 CSR PEM 資料寫入檔案
func SaveCSRToFile(csrPEM []byte, filename string) error {
	return ioutil.WriteFile(filename, csrPEM, 0600)
}

func SaveCertToFile(certPEM []byte, filename string) error {
	return ioutil.WriteFile(filename, certPEM, 0600)
}
