package fabric

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RegisterRequest struct {
	ID          string `json:"id"`
	Secret      string `json:"secret"`
	Affiliation string `json:"affiliation"`
	Type        string `json:"type"` // usually "client"
}

type EnrollRequest struct {
	Username string
	Password string
}

func RegisterUser(caURL, adminUser, adminPass string, req RegisterRequest) error {
	body := map[string]interface{}{
		"id":              req.ID,
		"secret":          req.Secret,
		"affiliation":     req.Affiliation,
		"type":            req.Type,
		"max_enrollments": -1,
	}
	b, _ := json.Marshal(body)

	httpReq, _ := http.NewRequest("POST", caURL+"/api/v1/register", bytes.NewReader(b))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Basic "+basicAuth(adminUser, adminPass))

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		data, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("❌ Register failed: %s", string(data))
	}
	return nil
}

func EnrollUser(caURL string, req EnrollRequest) ([]byte, []byte, error) {
	body := map[string]interface{}{
		"certificate_request": "", // 讓 server 自動產生
	}
	b, _ := json.Marshal(body)

	httpReq, _ := http.NewRequest("POST", caURL+"/api/v1/enroll", bytes.NewReader(b))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.SetBasicAuth(req.Username, req.Password)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		data, _ := ioutil.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("❌ Enroll failed: %s", string(data))
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	cert := result["result"].(map[string]interface{})["Cert"].(string)
	key := result["result"].(map[string]interface{})["ServerInfo"].(map[string]interface{})["CAName"].(string)

	return []byte(cert), []byte(key), nil // Enroll 預設只會給 Cert，Key 需自己產生或 CSR 模式
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
