# Go Server 連接測試說明

## 連接配置

Go Server 會在啟動時測試以下連接：

### 1. Gateway 連線測試
- 測試 Hyperledger Fabric Gateway 是否可用
- 如果失敗：會顯示警告，但不阻止服務啟動

### 2. Python Backend gRPC 連線測試（新增）
- 測試 Python Backend gRPC 服務是否可用
- 默認地址：`localhost:50052`
- 可通過環境變數 `PYTHON_BACKEND_GRPC_ADDR` 配置
- 如果失敗：會顯示警告和啟動提示，但不阻止服務啟動

## 配置方式

### 方法 1: 使用環境變數
```bash
export PYTHON_BACKEND_GRPC_ADDR=localhost:50052
go run main.go
```

### 方法 2: 使用 .env 文件
```bash
# 在 hyperledger/go_server/.env 文件中
PYTHON_BACKEND_GRPC_ADDR=localhost:50052
```

### 方法 3: 使用默認值
如果不設置，默認使用 `localhost:50052`

## 連接流程

```
Go Server 啟動
  ↓
1. 測試 Gateway 連線
  ↓
2. 測試 Python Backend gRPC 連線
  ↓
3. 啟動 Go Server gRPC (:50051)
  ↓
4. 啟動 HTTP Gateway (:8080)
```

## 實際調用時的連接

當前端調用 `/v1/health/analyze` 時：

```
前端 → HTTP Gateway (:8080)
  ↓
GetHealthAnalysis (gRPC :50051)
  ↓
HandleGetHealthAnalysis (service)
  ↓
callPythonAnalysisService
  ↓
連接到 Python Backend (:50052)
  ↓
調用 AnalyzeHealthReportForUser/ForInsurer
  ↓
返回分析結果
```

## 注意事項

1. **啟動順序**：
   - 先啟動 Python Backend（端口 50052）
   - 再啟動 Go Server（端口 50051）

2. **連接測試**：
   - 啟動時的測試是**非阻塞**的
   - 即使測試失敗，服務仍會啟動
   - 實際調用時如果連接失敗會返回錯誤

3. **錯誤處理**：
   - 如果 Python Backend 未運行，會返回友好的錯誤信息
   - 日誌會記錄詳細的連接錯誤
