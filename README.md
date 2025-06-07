# Healthcare Records Management System on Hyperledger Fabric

A decentralized healthcare records management system built on Hyperledger Fabric blockchain technology, providing secure, transparent, and privacy-preserving medical data management.

## 🏗️ Architecture Overview

This project consists of three main components:

- **Hyperledger Fabric Network**: Blockchain infrastructure with peers, orderers, and CouchDB
- **Smart Contracts (Chaincode)**: Go-based chaincode for health record management
- **Frontend Application**: Vue.js + Vuetify web interface
- **Backend Services**: Node.js/Express API services (if applicable)

## 🚀 Features

- **Secure Health Record Storage**: Store medical records on blockchain with cryptographic hashing
- **Access Control**: Patient-controlled access authorization system
- **Privacy Protection**: Patient identity protection through hash-based anonymization  
- **Multi-role Support**: Support for patients, clinics, and healthcare providers
- **Audit Trail**: Immutable transaction history for compliance and transparency
- **Access Request Management**: Workflow for requesting and approving access to medical records

## 📋 Prerequisites

Before running this project, ensure you have:

- Docker & Docker Compose
- Node.js (v16+ recommended)  
- Go (v1.19+ for chaincode development)
- Hyperledger Fabric binaries and Docker images

## 🛠️ Installation & Setup

### 1. Environment Setup

First, set up the required environment variables:

```bash
export HYPERLEDGER_HOME=$HOME/hyperledger
```

Or create a `.env` file in the `hyperledger/` directory:

```bash
echo "HYPERLEDGER_HOME=$HOME/hyperledger" > hyperledger/.env
```

### 2. Start Hyperledger Fabric Network

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

### 3. Deploy Smart Contracts

Deploy the health records chaincode:

```bash
# Install and instantiate chaincode (specific commands depend on your deployment scripts)
# Refer to your chaincode deployment scripts in hyperledger/sdk_server or similar
```

### 4. Start Frontend Application

```bash
cd frontend
npm install
npm run dev
```

The frontend will be available at `http://localhost:5173` (or the port specified by Vite).

### 5. Start Backend Services (if applicable)

```bash
cd backend
npm install
npm start
```

## 🏥 Usage

### For Patients
- View your medical records
- Authorize access to specific records
- Manage access requests from healthcare providers
- Monitor who has accessed your data

### For Healthcare Providers  
- Upload new medical reports
- Request access to patient records
- View authorized patient data
- Maintain audit trails

### For Clinics
- Upload patient health reports
- Manage patient data securely
- Comply with healthcare regulations

## 🔧 Configuration

### Network Configuration
- Modify `hyperledger/docker-compose.yaml` for network topology changes
- Update `hyperledger/configtx.yaml` for channel and organization configurations

### Smart Contract Configuration
- Health record structure can be modified in `hyperledger/chaincode-go/health_contract.go`
- Access control policies are defined in the chaincode

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

## 📁 Project Structure

```
.
├── hyperledger/              # Blockchain network configuration
│   ├── docker-compose.yaml   # Docker services definition
│   ├── configtx.yaml        # Channel configuration
│   ├── chaincode-go/         # Smart contracts
│   ├── orgs/                 # Organization certificates
│   └── sdk_server/           # SDK and deployment scripts
├── frontend/                 # Vue.js frontend application
│   ├── src/                  # Source code
│   ├── package.json          # Dependencies
│   └── vite.config.js        # Vite configuration
├── backend/                  # Backend services (if applicable)
└── README.md                 # This file
```

## 🔐 Security Features

- **Identity Management**: Certificate-based authentication
- **Data Privacy**: Patient identity hashing and anonymization
- **Access Control**: Role-based permissions (patient, clinic, provider)
- **Audit Trail**: Immutable transaction logging
- **TLS Encryption**: Secure communication between components

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Troubleshooting

### Common Issues

1. **Docker containers not starting**: Ensure Docker daemon is running and ports are available
2. **Certificate errors**: Check if certificates are properly generated in `hyperledger/orgs/`
3. **Chaincode deployment fails**: Verify Go version and dependencies in `chaincode-go/`
4. **Frontend connection errors**: Check if backend services are running and accessible

### Support

For support and questions:
- Check the troubleshooting section above
- Review Hyperledger Fabric documentation
- Create an issue in this repository

## 🙏 Acknowledgments

- [Hyperledger Fabric](https://hyperledger-fabric.readthedocs.io/) for the blockchain platform
- [Vue.js](https://vuejs.org/) for the frontend framework
- [Vuetify](https://vuetifyjs.com/) for the UI components
