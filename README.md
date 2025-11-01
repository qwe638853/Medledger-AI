# Healthcare Records Management System on Hyperledger Fabric

A decentralized healthcare records management system built on Hyperledger Fabric blockchain technology, providing secure, transparent, and privacy-preserving medical data management with AI-powered health analysis.

## 🏗️ Architecture Overview

This project consists of four main components:

- **Hyperledger Fabric Network**: Blockchain infrastructure with peers, orderers, and CouchDB
- **Smart Contracts (Chaincode)**: Go-based chaincode for health record management
- **Backend Services**: Python gRPC server with AI-powered health analysis using LangChain and Ollama
- **Frontend Application**: Vue.js + Vuetify web interface

## 🚀 Features

- **Secure Health Record Storage**: Store medical records on blockchain with cryptographic hashing
- **Access Control**: Patient-controlled access authorization system
- **Privacy Protection**: Patient identity protection through hash-based anonymization  
- **Multi-role Support**: Support for patients, clinics, and healthcare providers
- **Audit Trail**: Immutable transaction history for compliance and transparency
- **Access Request Management**: Workflow for requesting and approving access to medical records
- **AI Health Analysis**: Intelligent health report analysis for both users and insurers
- **Multi-language Support**: Traditional Chinese language support with medical terminology translation
- **Risk Assessment**: AI-powered disease risk evaluation and personalized recommendations
- **Insurance Integration**: Automated policy recommendations based on health metrics

## 📋 Prerequisites

Before running this project, ensure you have:

- Docker & Docker Compose
- Node.js (v16+ recommended)  
- Go (v1.19+ for chaincode development)
- Python 3.8+ with pip
- Ollama with Llama3 model (for AI analysis)
- Hyperledger Fabric binaries and Docker images

## 🛠️ Installation & Setup


###  Start Hyperledger Fabric Network

Navigate to the hyperledger directory and start the network:

```bash
cd hyperledger
docker-compose up -d
```

This will start:
- Certificate Authorities (Root CA, TLS CA)
- Peer nodes (peer1-org1, peer2-org1)
- Orderer service (orderer1-org1)
- CouchDB instances for state database

###  Deploy Smart Contracts

部署健康記錄鏈碼的完整步驟：

#### 前置準備

確保你已經啟動 Hyperledger Fabric 網路（見上方「Start Hyperledger Fabric Network」）。

假設你在 `fabric-samples` 目錄下，設定環境變數：

```bash
# 在 fabric-samples 目錄執行
export PATH=$PATH:$PWD/bin
export FABRIC_CFG_PATH=$PWD/config

# 設定組織環境變數（假設使用 test-network）
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_MSPCONFIGPATH=$PWD/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
```

#### 步驟 1: 打包鏈碼

```bash
# 從 fabric-samples 目錄，打包鏈碼
peer lifecycle chaincode package health.tar.gz \
  --path ../Medledger-AI/hyperledger/chaincode-go \
  --lang golang \
  --label health_1.0
```

**說明**：
- `health.tar.gz`：打包後的檔案名稱
- `--path`：指向你的鏈碼目錄（相對於 fabric-samples 的路徑）
- `--lang golang`：指定使用 Go 語言
- `--label health_1.0`：鏈碼標籤，用於識別版本

#### 步驟 2: 安裝鏈碼到所有 Peer

**安裝到 Peer1-Org1：**

```bash
# 設定 Peer1 環境變數
export CORE_PEER_ADDRESS=localhost:7051

# 安裝鏈碼
peer lifecycle chaincode install health.tar.gz
```

**安裝到 Peer2-Org1（如果需要）：**

```bash
# 設定 Peer2 環境變數
export CORE_PEER_ADDRESS=localhost:7053

# 安裝鏈碼
peer lifecycle chaincode install health.tar.gz
```

**確認安裝成功：**

```bash
peer lifecycle chaincode queryinstalled
```

這會顯示已安裝的鏈碼及其 package ID（記下這個 ID，後續會用到）。

#### 步驟 3: 批准鏈碼定義

```bash
# 設定 Peer1 環境變數
export CORE_PEER_ADDRESS=localhost:7051

# 批准鏈碼定義（使用步驟 2 獲取的 package ID）
peer lifecycle chaincode approveformyorg \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --channelID mychannel \
  --name health \
  --version 1.0 \
  --package-id <PACKAGE_ID> \
  --sequence 1 \
  --tls \
  --cafile $PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```

**說明**：
- `<PACKAGE_ID>`：替換為步驟 2 獲取的 package ID
- `--name health`：鏈碼名稱
- `--version 1.0`：鏈碼版本
- `--sequence 1`：升級序列號（首次部署為 1）
- `--channelID mychannel`：目標頻道名稱

#### 步驟 4: 提交鏈碼定義

```bash
# 提交鏈碼定義（提交後鏈碼即可使用）
peer lifecycle chaincode commit \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --channelID mychannel \
  --name health \
  --version 1.0 \
  --sequence 1 \
  --tls \
  --cafile $PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```

#### 步驟 5: 驗證部署

**查詢鏈碼資訊：**

```bash
peer lifecycle chaincode querycommitted --channelID mychannel --name health
```

**測試調用鏈碼方法：**

```bash
# 設定 Peer1 環境變數
export CORE_PEER_ADDRESS=localhost:7051

# 測試調用（例如查詢病患的報告 meta）
peer chaincode query \
  -C mychannel \
  -n health \
  -c '{"function":"ListMyReportMeta","Args":[]}'
```

#### 完整部署腳本範例

```bash
#!/bin/bash
# deploy_chaincode.sh

# 設定環境
export PATH=$PATH:$PWD/bin
export FABRIC_CFG_PATH=$PWD/config
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_MSPCONFIGPATH=$PWD/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp

# 步驟 1: 打包
echo "打包鏈碼..."
peer lifecycle chaincode package health.tar.gz \
  --path ../Medledger-AI/hyperledger/chaincode-go \
  --lang golang \
  --label health_1.0

# 步驟 2: 安裝（Peer1）
echo "安裝鏈碼到 Peer1..."
export CORE_PEER_ADDRESS=localhost:7051
peer lifecycle chaincode install health.tar.gz

# 獲取 Package ID
PACKAGE_ID=$(peer lifecycle chaincode queryinstalled | grep "health_1.0" | cut -d' ' -f3 | cut -d',' -f1)
echo "Package ID: $PACKAGE_ID"

# 步驟 3: 批准
echo "批准鏈碼定義..."
peer lifecycle chaincode approveformyorg \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --channelID mychannel \
  --name health \
  --version 1.0 \
  --package-id $PACKAGE_ID \
  --sequence 1 \
  --tls \
  --cafile $PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

# 步驟 4: 提交
echo "提交鏈碼定義..."
peer lifecycle chaincode commit \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --channelID mychannel \
  --name health \
  --version 1.0 \
  --sequence 1 \
  --tls \
  --cafile $PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

echo "鏈碼部署完成！"
```

#### 注意事項

- 確保網路已啟動且所有服務運行正常
- 路徑需根據實際專案位置調整
- 頻道名稱需與實際建立的頻道一致
- 升級鏈碼時，需增加 `--sequence` 號碼
- 如使用 TLS，需提供正確的 CA 證書路徑

###  Setup AI Backend Services

#### Install Python Dependencies
```bash
cd backend/health_check_project
pip install -r requirements.txt  # Create requirements.txt with necessary packages
```

#### Install Required Python Packages
```bash
pip install grpcio grpcio-tools
pip install langchain langchain-chroma langchain-huggingface langchain-ollama
pip install chromadb
pip install protobuf
```

#### Setup Ollama (for AI Analysis)
```bash
# Install Ollama
curl -fsSL https://ollama.ai/install.sh | sh

# Pull Llama3 model
ollama pull llama3:8b

# Start Ollama service
ollama serve
```

#### Setup ChromaDB Vector Database
```bash
cd backend/health_check_project
python create_collection.py
python add_data.py
```

#### Start Backend gRPC Server
```bash
cd backend/health_check_project
python test.py
```

The backend server will start on `localhost:50051`

###  Start Frontend Application

```bash
cd frontend
npm install
npm run dev
```

The frontend will be available at `http://localhost:5173` (or the port specified by Vite).

## 🏥 Usage

### For Patients
- View your medical records
- Authorize access to specific records
- Manage access requests from healthcare providers
- Monitor who has accessed your data
- **Get AI-powered health analysis and personalized recommendations**

### For Healthcare Providers  
- Upload new medical reports
- Request access to patient records
- View authorized patient data
- Maintain audit trails
- **Access professional health analysis with risk assessments**

### For Clinics
- Upload patient health reports
- Manage patient data securely
- Comply with healthcare regulations

### For Insurance Companies
- **Automated risk assessment based on health metrics**
- **AI-powered policy recommendations**
- **Comprehensive health report analysis for underwriting**

## 🤖 AI Analysis Features

### Health Report Analysis
- **Comprehensive Health Summary**: Detailed analysis of all health metrics
- **Risk Assessment**: AI-powered evaluation of potential health risks
- **Personalized Recommendations**: Customized advice for diet, exercise, and medical monitoring
- **Insurance Policy Suggestions**: Automated recommendations for suitable insurance products

### Technical Implementation
- **LangChain Framework**: Advanced prompt engineering and retrieval-augmented generation (RAG)
- **Vector Database**: ChromaDB for efficient similarity search and context retrieval
- **Multi-Query Retrieval**: Enhanced context gathering through category-based queries
- **HyDE (Hypothetical Document Embeddings)**: Improved retrieval accuracy
- **Medical Terminology Translation**: Automatic English-to-Chinese medical term conversion

## 🔧 Configuration

### Network Configuration
- Modify `hyperledger/docker-compose.yaml` for network topology changes
- Update `hyperledger/configtx.yaml` for channel and organization configurations

### Smart Contract Configuration
- Health record structure can be modified in `hyperledger/chaincode-go/health_contract.go`
- Access control policies are defined in the chaincode

### Backend Configuration
- **gRPC Service Configuration**: Modify `backend/health_check_project/test.py`
- **AI Model Settings**: Configure Ollama model parameters and ChromaDB paths
- **Medical Translations**: Update translation dictionaries for different languages
- **Analysis Prompts**: Customize LangChain prompts for different analysis scenarios

### Frontend Configuration
- Update API endpoints in frontend configuration files
- Modify UI components in `frontend/src/`

## 🐳 Docker Services

The system runs the following Docker containers:

| Service | Purpose | Port |
|---------|---------|------|
| root-ca | Root Certificate Authority | 7054 |
| ca-tls | TLS Certificate Authority | 7052 |
| peer1-org1 | Primary peer node | 7051 |
| peer2-org1 | Secondary peer node | 7053 |
| orderer1-org1 | Ordering service | 7050 |
| couchdb1/couchdb2 | State databases | 5984 |


