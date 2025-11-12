# Go Server 超時配置說明

## Python Backend 調用超時設置

Go Server 在調用 Python Backend 進行健康分析時，可以通過環境變數配置超時時間。

### 配置方式

#### 方法 1: 使用環境變數
```bash
export PYTHON_BACKEND_TIMEOUT_SECONDS=600  # 10 分鐘
go run main.go
```

#### 方法 2: 使用 .env 文件
在 `hyperledger/go_server/.env` 文件中添加：
```bash
PYTHON_BACKEND_GRPC_ADDR=localhost:50052
PYTHON_BACKEND_TIMEOUT_SECONDS=600  # 10 分鐘（秒）
```

#### 方法 3: 使用默認值
如果不設置，默認超時時間為 **300 秒（5 分鐘）**

### 超時時間建議

根據 Ollama 分析的速度，建議設置：

- **快速模型（llama3:8b）**: 300-600 秒（5-10 分鐘）
- **大型模型**: 600-1200 秒（10-20 分鐘）
- **極慢的模型或複雜分析**: 1200+ 秒（20 分鐘以上）

### 工作原理

1. Go Server 接收到分析請求
2. 創建一個新的 context，設置超時時間（從環境變數讀取）
3. 使用這個 context 調用 Python Backend 的 gRPC 服務
4. 如果 Python Backend 在超時時間內未完成，會返回 `DeadlineExceeded` 錯誤

### 錯誤處理

如果超時，會返回以下錯誤：
```
Python Backend 分析超時（超過 600 秒）
```

### 日誌輸出

啟動時會顯示當前配置的超時時間：
```
[callPythonAnalysisService] 連接到 Python Backend: localhost:50052, 超時時間: 600 秒
```

### 注意事項

1. **超時時間應該大於 Python Backend 的實際分析時間**
2. **如果經常超時，考慮增加超時時間或優化 Python Backend 的性能**
3. **超時時間過長可能會導致前端請求等待時間過長**

