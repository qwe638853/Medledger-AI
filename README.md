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

Deploy the health records chaincode:

```bash
# Install and instantiate chaincode (specific commands depend on your deployment scripts)
# Refer to your chaincode deployment scripts in hyperledger/sdk_server or similar
```

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


