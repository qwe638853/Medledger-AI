package fabric

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/pkg/errors"
)

type RegisterRequest struct {
	ID          string `json:"id"`
	Secret      string `json:"secret"`
	Type        string `json:"type"`
	Affiliation string `json:"affiliation"`
}

func B64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func uriPath(full string) string {
	u, err := url.Parse(full)
	if err != nil {
		return full
	}
	return u.RequestURI()
}

// ✨ BCCSP-based token generation
func GenECDSAToken(csp bccsp.BCCSP, cert []byte, key bccsp.Key, method, uri string, body []byte) (string, error) {
	b64body := B64Encode(body)
	b64cert := B64Encode(cert)
	b64uri := B64Encode([]byte(uriPath(uri)))
	payload := method + "." + b64uri + "." + b64body + "." + b64cert

	return genECDSAToken(csp, key, b64cert, payload)
}

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

// 📤 RegisterUser using BCCSP
func RegisterUser(caURL, certPath, keyPath string, req RegisterRequest) error {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("❌ Failed to marshal request: %w", err)
	}

	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return fmt.Errorf("❌ Failed to read cert: %w", err)
	}
	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return fmt.Errorf("❌ Failed to read key: %w", err)
	}

	// Init BCCSP
	factory.InitFactories(nil)
	csp := factory.GetDefault()

	keyBlock, _ := pem.Decode(keyPEM)
	key, err := csp.KeyImport(keyBlock.Bytes, &bccsp.ECDSAPrivateKeyImportOpts{Temporary: true})
	if err != nil {
		return fmt.Errorf("❌ Failed to import key: %w", err)
	}

	token, err := GenECDSAToken(csp, certPEM, key, "POST", caURL+"/api/v1/register", bodyBytes)
	if err != nil {
		return fmt.Errorf("❌ Failed to generate token: %w", err)
	}

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

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return fmt.Errorf("❌ Register failed (%d): %s", resp.StatusCode, respBody)
	}

	fmt.Printf("✅ Register success: %s\n", respBody)
	return nil
}
