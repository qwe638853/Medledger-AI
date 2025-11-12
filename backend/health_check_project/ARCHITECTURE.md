# Health Check Project 專案架構解析

## 📋 專案概述

這是一個基於 AI 的健康檢查數據分析系統，使用 Python 實現，提供 gRPC 服務接口，整合 Ollama LLM 進行智能健康分析。

---

## 🏗️ 專案結構

```
health_check_project/
├── test.py                    # 主 gRPC 服務器（RAG 4.0/4.1 版本）
├── client.py                  # gRPC 客戶端測試工具
├── proto/                     # Protocol Buffers 定義
│   ├── data.proto            # gRPC 服務和消息定義
│   └── data_pb2.py           # 生成的 Python proto 文件
│   └── data_pb2_grpc.py      # 生成的 gRPC 服務文件
├── google/                    # Google Protobuf 依賴
│   ├── api/                  # gRPC-Gateway 註解
│   └── protobuf/             # Protobuf 標準定義
├── first/                     # 舊版 FastAPI 服務器（備用）
│   ├── main.py               # FastAPI 主服務器
│   ├── api.py                # API 路由定義
│   ├── analyze_health_data.py # 健康數據分析邏輯
│   └── ...
├── chroma_db/                 # ChromaDB 向量數據庫（用於 first/ 版本）
├── health_rules.txt          # 醫學知識庫（RAG 4.0 使用）
├── requirements.txt          # Python 依賴
├── wsl_setup.sh              # WSL 環境初始化腳本
├── start_grpc_server.sh      # 啟動 gRPC 服務器腳本
├── start_ollama.sh           # 啟動 Ollama 服務腳本
└── check_model.sh            # 檢查 Ollama 模型腳本
```

---

## 🔧 核心組件

### 1. **主服務器：test.py**

#### 功能
- 提供 gRPC 服務接口，實現 `HealthService` 協議
- 使用 **RAG 4.0/4.1** 進行健康數據分析
- 整合 Ollama LLM 進行智能分析

#### 主要類別

##### `HealthAnalysisServicer`
- 實現 `data_pb2_grpc.HealthServiceServicer`
- 提供兩個主要 RPC 方法：
  - `AnalyzeHealthReportForUser`: 用戶健康分析
  - `AnalyzeHealthReportForInsurer`: 保險業者風險評估

#### RAG 4.0 架構

```
┌─────────────────────────────────────────┐
│  健康檢查數據 (test_results)            │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│  _get_key_to_topic_map()                │
│  映射指標 → 知識主題                     │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│  _filter_context_for_rag()               │
│  動態篩選相關知識塊                       │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│  health_rules.txt                        │
│  醫學知識庫（按主題分割）                 │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│  Ollama LLM (llama3:8b)                  │
│  生成結構化 JSON 分析                     │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│  Pydantic 驗證                           │
│  HealthAnalysis / InsurerAnalysis        │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│  gRPC Response                           │
│  UserHealthAnalysisResponse /            │
│  InsurerHealthAnalysisResponse           │
└─────────────────────────────────────────┘
```

#### 關鍵方法

1. **`_load_health_rules()`**
   - 載入 `health_rules.txt` 醫學知識庫
   - 按 `### 主題` 分割為知識塊
   - 返回完整文本和知識塊字典

2. **`_get_key_to_topic_map()`**
   - 定義健檢指標到知識主題的映射
   - 例如：`Glu-AC`, `HbA1c` → `血糖指標`

3. **`_filter_context_for_rag()`**
   - 根據輸入的健檢指標，動態篩選相關知識塊
   - 只包含與輸入指標相關的醫學知識
   - 提高上下文相關性和效率

4. **`_call_ollama_json()`**
   - 調用 Ollama LLM 生成結構化 JSON
   - 使用 Pydantic schema 驗證輸出
   - 支持 JSON Schema 格式約束

5. **`_check_ollama_connection()`**
   - 檢查 Ollama 服務可用性
   - 支持重試機制

---

### 2. **Protocol Buffers：proto/data.proto**

#### 服務定義
```protobuf
service HealthService {
  // Python Backend 提供
  rpc AnalyzeHealthReportForUser(...) returns (...);
  rpc AnalyzeHealthReportForInsurer(...) returns (...);
  
  // Go Server 提供（通過 HTTP Gateway）
  rpc GetHealthAnalysis(...) returns (...);
  rpc UploadReport(...) returns (...);
  // ... 其他方法
}
```

#### 消息類型
- `AnalyzeHealthReportRequest`: 分析請求
- `UserHealthAnalysisResponse`: 用戶分析結果
- `InsurerHealthAnalysisResponse`: 保險業者分析結果
- `HealthAnalysisResponse`: 統一分析響應（包含用戶或保險業者結果）

---

### 3. **醫學知識庫：health_rules.txt**

#### 結構
- 按主題分割（使用 `### 主題` 標記）
- 包含：
  - 核心評分準則
  - 血糖指標規則
  - 肝功能指標規則
  - 腎功能指標規則
  - 血脂指標規則
  - 全血球計數規則
  - 凝血功能規則
  - 發炎指標規則
  - 癌症指標規則
  - 尿液常規檢查規則

#### 使用方式
- RAG 4.0 根據輸入指標動態篩選相關主題
- 只將相關知識塊傳遞給 LLM，減少 token 消耗

---

### 4. **舊版服務器：first/main.py**

#### 功能
- FastAPI REST API 服務器
- 使用 ChromaDB 向量數據庫進行 RAG
- 整合 Azure SQL Database
- 支持文件上傳（PDF/DOCX）

#### 與 test.py 的差異

| 特性 | test.py (RAG 4.0) | first/main.py (舊版) |
|------|------------------|---------------------|
| 接口類型 | gRPC | REST API (FastAPI) |
| 知識庫 | 文本文件 (health_rules.txt) | ChromaDB 向量數據庫 |
| RAG 方式 | 規則映射 + 主題篩選 | 向量相似度搜索 |
| LLM | Ollama (llama3:8b) | Ollama (llama3:8b) |
| 數據庫 | 無（通過 Go Server） | Azure SQL Database |

---

## 🔌 服務通信架構

```
┌─────────────┐
│  前端 (Vue) │
└──────┬──────┘
       │ HTTP REST
       ▼
┌─────────────────────────────────┐
│  Go Server (Hyperledger Fabric) │
│  - gRPC: :50051                  │
│  - HTTP Gateway: :8080           │
└──────┬──────────────────────────┘
       │ gRPC
       ▼
┌─────────────────────────────────┐
│  Python Backend (test.py)       │
│  - gRPC: :50052                  │
│  - Ollama: localhost:11434       │
└─────────────────────────────────┘
```

### 請求流程

1. **前端請求分析**
   ```
   前端 → Go Server HTTP Gateway (:8080)
        → POST /v1/health/analyze
   ```

2. **Go Server 處理**
   ```
   Go Server → 從 Hyperledger Fabric 讀取報告
            → 解密數據
            → 調用 Python Backend gRPC (:50052)
   ```

3. **Python Backend 分析**
   ```
   Python Backend → RAG 4.0 篩選知識
                 → Ollama LLM 分析
                 → 返回結構化結果
   ```

4. **返回結果**
   ```
   Python Backend → Go Server
                 → HTTP Gateway
                 → 前端
   ```

---

## 📦 依賴項

### Python 依賴 (`requirements.txt`)
- `grpcio==1.66.0`: gRPC 框架
- `grpcio-tools==1.66.0`: Protocol Buffers 編譯工具
- `pydantic==2.9.2`: 數據驗證
- `ollama==0.6.0`: Ollama LLM 客戶端
- `langchain-chroma==0.1.2`: ChromaDB 集成（舊版使用）
- `langchain-huggingface==0.0.3`: HuggingFace 嵌入模型（舊版使用）
- `chromadb==0.5.3`: 向量數據庫（舊版使用）
- `sentence-transformers==3.0.1`: 句子嵌入模型（舊版使用）

---

## 🚀 啟動流程

### 1. 環境初始化
```bash
./wsl_setup.sh
```
- 安裝 Python 依賴
- 編譯 Protocol Buffers
- 設置虛擬環境

### 2. 啟動 Ollama 服務
```bash
./start_ollama.sh
# 或手動：
ollama serve
```

### 3. 下載 Ollama 模型
```bash
ollama pull llama3:8b
```

### 4. 啟動 Python gRPC 服務器
```bash
./start_grpc_server.sh
# 或手動：
source venv/bin/activate
python test.py
```

### 5. 環境變數配置

#### Python Backend (.env)
```bash
OLLAMA_HOST=http://localhost:11434
OLLAMA_MODEL=llama3:8b
PYTHON_BACKEND_GRPC_PORT=50052
PYTHON_BACKEND_GRPC_HOST=[::]
```

#### Go Server (.env)
```bash
PYTHON_BACKEND_GRPC_ADDR=localhost:50052
```

---

## 🔍 RAG 4.0 vs RAG 4.1

### RAG 4.0
- 基於規則映射的主題篩選
- 使用 `_get_key_to_topic_map()` 映射指標到主題
- 動態篩選相關知識塊
- 不依賴向量數據庫

### RAG 4.1
- 在 RAG 4.0 基礎上增強
- 恢復列表結構（`disease_risks` 為列表而非單一對象）
- 強化 Prompt 指令
- 更嚴格的 JSON 格式驗證

---

## 📊 數據流

### 用戶分析流程

```
1. 接收 AnalyzeHealthReportRequest
   ├─ report_id
   ├─ patient_id
   └─ test_results_json

2. 解析 test_results_json
   └─ 提取健檢指標

3. RAG 4.0 知識篩選
   ├─ 映射指標 → 主題
   ├─ 篩選相關知識塊
   └─ 構建上下文

4. 構建 Prompt
   ├─ 健康檢查數據
   ├─ 相關醫學知識
   └─ JSON Schema 約束

5. 調用 Ollama LLM
   └─ 生成結構化 JSON

6. Pydantic 驗證
   └─ HealthAnalysis 模型

7. 轉換為 Proto 消息
   └─ UserHealthAnalysisResponse

8. 返回 gRPC 響應
```

---

## 🛠️ 工具腳本

### `wsl_setup.sh`
- WSL 環境初始化
- 安裝 Python 依賴
- 編譯 Protocol Buffers

### `start_grpc_server.sh`
- 啟動 Python gRPC 服務器
- 設置環境變數
- 檢查依賴

### `start_ollama.sh`
- 啟動 Ollama 服務

### `check_model.sh`
- 檢查 Ollama 模型是否已安裝

### `client.py`
- gRPC 客戶端測試工具
- 用於測試 Python Backend 服務

---

## 🔐 安全與配置

### 端口配置
- **Go Server gRPC**: `:50051`
- **Python Backend gRPC**: `:50052` (可配置)
- **Go Server HTTP Gateway**: `:8080`
- **Ollama**: `:11434`

### 環境變數
- 所有配置通過環境變數管理
- 支持 `.env` 文件（需手動實現）

---

## 📝 注意事項

1. **模型依賴**
   - 默認使用 `llama3:8b`
   - 可通過 `OLLAMA_MODEL` 環境變數配置
   - 確保 Ollama 服務運行且模型已下載

2. **知識庫文件**
   - `health_rules.txt` 必須存在於項目根目錄
   - 文件編碼必須為 UTF-8

3. **Proto 文件同步**
   - `proto/data.proto` 必須與 Go Server 的 `data.proto` 同步
   - 修改後需重新編譯：`python -m grpc_tools.protoc ...`

4. **虛擬環境**
   - 建議使用虛擬環境 (`venv`)
   - 啟動前需激活：`source venv/bin/activate`

---

## 🔄 版本演進

- **v1.0 (first/)**: FastAPI + ChromaDB 向量數據庫
- **v2.0 (test.py)**: gRPC + RAG 4.0 規則映射
- **v2.1 (test.py)**: RAG 4.1 增強版（恢復列表結構）

---

## 📚 相關文檔

- `README_OLLAMA.md`: Ollama 安裝與配置指南
- `PORT_CONFIG.md`: 端口配置說明
- `.env.example`: 環境變數示例

---

## 🎯 未來改進方向

1. **性能優化**
   - 知識庫緩存機制
   - LLM 響應緩存
   - 並發處理優化

2. **功能增強**
   - 支持更多健檢指標
   - 多語言支持
   - 歷史數據對比分析

3. **架構優化**
   - 微服務化
   - 容器化部署
   - 監控與日誌系統

