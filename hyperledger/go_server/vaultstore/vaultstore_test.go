package vaultstore

import (
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "os"
    "strings"
    "sync"
    "testing"
)

// fakeKVServer 建立一個最小的 KV v2 相容測試伺服器
// - 支援 Write: POST /v1/{mount}/data/{path}
// - 支援 Read:  GET  /v1/{mount}/data/{path}
func fakeKVServer(t *testing.T, mount string) *httptest.Server {
    t.Helper()
    store := struct {
        mu sync.Mutex
        m  map[string]map[string]string // key: path, val: data(csr/key/cert)
    }{m: make(map[string]map[string]string)}

    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 期望路徑前綴 /v1/{mount}/data/
        prefix := "/v1/" + mount + "/data/"
        if !strings.HasPrefix(r.URL.Path, prefix) {
            http.Error(w, "not found", http.StatusNotFound)
            return
        }
        path := strings.TrimPrefix(r.URL.Path, prefix)

        switch r.Method {
        case http.MethodPost:
            // 寫入：請求 JSON 應為 {"data": {"csr":"...","key":"...","cert":"..."}}
            var body struct{
                Data map[string]string `json:"data"`
            }
            if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }
            store.mu.Lock()
            store.m[path] = body.Data
            store.mu.Unlock()
            // 回傳空成功（Vault 會帶 metadata，但對本測試非必要）
            _ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"metadata": map[string]any{}}})

        case http.MethodGet:
            store.mu.Lock()
            data, ok := store.m[path]
            store.mu.Unlock()
            if !ok {
                http.Error(w, "not found", http.StatusNotFound)
                return
            }
            // 回傳 KV v2 結構：{"data":{"data":{...},"metadata":{}}}
            resp := map[string]any{
                "data": map[string]any{
                    "data":     data,
                    "metadata": map[string]any{},
                },
            }
            _ = json.NewEncoder(w).Encode(resp)

        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    return httptest.NewServer(handler)
}

// TestWriteAndReadUserMaterial
// 目的：驗證 WriteUserMaterial / ReadUserMaterial 的基本讀寫流程（KV v2）
func TestWriteAndReadUserMaterial(t *testing.T) {
    srv := fakeKVServer(t, "kv")
    defer srv.Close()

    // 設定 Vault 連線環境（僅測試用途）
    t.Setenv("VAULT_ADDR", srv.URL)
    t.Setenv("VAULT_TOKEN", "dev-token")
    t.Setenv("VAULT_MOUNT", "kv")

    st, err := NewFromEnv()
    if err != nil { t.Fatalf("new store: %v", err) }

    ctx := context.Background()
    userID := "alice"
    csr := []byte("CSR-PEM")
    key := []byte("KEY-PEM")
    cert := []byte("CERT-PEM")

    if err := st.WriteUserMaterial(ctx, userID, csr, key, cert); err != nil {
        t.Fatalf("write user: %v", err)
    }

    rc, rk, rcert, err := st.ReadUserMaterial(ctx, userID)
    if err != nil { t.Fatalf("read user: %v", err) }
    if string(rc) != string(csr) || string(rk) != string(key) || string(rcert) != string(cert) {
        t.Fatalf("mismatch: got (%s,%s,%s)", rc, rk, rcert)
    }
}

// TestReadPath
// 目的：驗證 ReadPath 可用相對路徑讀取資料（users/{id}）
func TestReadPath(t *testing.T) {
    srv := fakeKVServer(t, "kv")
    defer srv.Close()

    t.Setenv("VAULT_ADDR", srv.URL)
    t.Setenv("VAULT_TOKEN", "dev-token")
    t.Setenv("VAULT_MOUNT", "kv")

    st, err := NewFromEnv()
    if err != nil { t.Fatalf("new store: %v", err) }

    ctx := context.Background()
    userID := "bob"
    csr := []byte("CSR-PEM")
    key := []byte("KEY-PEM")
    cert := []byte("CERT-PEM")

    if err := st.WriteUserMaterial(ctx, userID, csr, key, cert); err != nil {
        t.Fatalf("write user: %v", err)
    }

    // 使用相對路徑讀取
    rc, rk, rcert, err := st.ReadPath(ctx, "users/"+userID)
    if err != nil { t.Fatalf("read path: %v", err) }
    if string(rc) != string(csr) || string(rk) != string(key) || string(rcert) != string(cert) {
        t.Fatalf("mismatch: got (%s,%s,%s)", rc, rk, rcert)
    }
}

// TestWriteAndReadOthers
// 目的：驗證 insurers/clinics/platform 讀寫 API
func TestWriteAndReadOthers(t *testing.T) {
    srv := fakeKVServer(t, "kv")
    defer srv.Close()

    t.Setenv("VAULT_ADDR", srv.URL)
    t.Setenv("VAULT_TOKEN", "dev-token")
    t.Setenv("VAULT_MOUNT", "kv")

    st, err := NewFromEnv()
    if err != nil { t.Fatalf("new store: %v", err) }
    ctx := context.Background()

    // insurers
    if err := st.WriteInsurerMaterial(ctx, "ins1", []byte("C"), []byte("K"), []byte("T")); err != nil {
        t.Fatalf("write insurer: %v", err)
    }
    _, _, _, err = st.ReadInsurerMaterial(ctx, "ins1")
    if err != nil { t.Fatalf("read insurer: %v", err) }

    // clinics
    if err := st.WriteClinicMaterial(ctx, "cli1", []byte("C"), []byte("K"), []byte("T")); err != nil {
        t.Fatalf("write clinic: %v", err)
    }
    _, _, _, err = st.ReadClinicMaterial(ctx, "cli1")
    if err != nil { t.Fatalf("read clinic: %v", err) }

    // platform
    if err := st.WritePlatformMaterial(ctx, []byte("C"), []byte("K"), []byte("T")); err != nil {
        t.Fatalf("write platform: %v", err)
    }
    _, _, _, err = st.ReadPlatformMaterial(ctx)
    if err != nil { t.Fatalf("read platform: %v", err) }
}


