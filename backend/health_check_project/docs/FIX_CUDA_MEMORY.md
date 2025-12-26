# 修復 Ollama CUDA 內存錯誤

## 問題描述

當遇到以下錯誤時：
```
Ollama 錯誤: llama runner process has terminated: error loading model: unable to allocate CUDA0 buffer
```

### 為什麼 GPU 內存顯示為 0 卻還報錯？

這是 **WSL2 環境的常見問題**，可能原因：

1. **WSL2 GPU 內存映射問題**：WSL2 的 GPU 內存管理與原生 Linux 不同，Ollama 可能無法正確檢測或分配 GPU 內存
2. **內存碎片化**：Ollama 嘗試一次性分配大塊連續內存（例如 4-8GB），但 WSL2 無法提供足夠的連續內存塊
3. **CUDA 驅動兼容性**：WSL2 的 CUDA 支持可能需要特定配置
4. **模型載入策略**：Ollama 默認嘗試將整個模型載入 GPU，但 WSL2 環境下這種方式可能失敗

**即使 GPU 監控顯示內存為 0，Ollama 仍會嘗試分配 GPU 內存，導致分配失敗。**

## 解決方案（不改模型）

### 方法 0: WSL2 專用修復（推薦，針對您的問題）

如果您在 WSL2 環境中遇到「GPU 內存為 0 但分配失敗」的問題：

```bash
cd backend/health_check_project/scripts
./fix_wsl2_gpu.sh
```

這個腳本會：
1. 自動檢測 WSL2 環境
2. 使用**分層 GPU 載入**（避免一次性分配大量內存）
3. 設置較少的 GPU 層數（15 層，對 8GB GPU 較安全）
4. 優化 WSL2 特定的配置

**原理**：通過限制 GPU 層數，讓 Ollama 分批次載入模型，而不是一次性分配所有內存。

### 方法 1: 使用 CPU 模式（最簡單，但較慢）

執行修復腳本：
```bash
cd backend/health_check_project/scripts
./fix_ollama_cuda.sh
```

這個腳本會：
1. 停止現有的 Ollama 服務
2. 設置環境變數強制使用 CPU
3. 重新啟動 Ollama 服務

### 方法 2: 手動設置環境變數

1. **停止 Ollama 服務**：
```bash
pkill -f ollama
```

2. **設置環境變數並啟動**：
```bash
export OLLAMA_NUM_GPU=0
export OLLAMA_GPU_LAYERS=0
ollama serve
```

或者使用 nohup 在背景運行：
```bash
OLLAMA_NUM_GPU=0 OLLAMA_GPU_LAYERS=0 nohup ollama serve > /tmp/ollama.log 2>&1 &
```

### 方法 3: 部分使用 GPU（WSL2 推薦）

**這是解決 WSL2 GPU 內存分配問題的最佳方法**：使用分層載入，而不是一次性分配。

```bash
# 停止 Ollama
pkill -f ollama

# WSL2 環境建議：使用較少的 GPU 層數
export OLLAMA_NUM_GPU=1
export OLLAMA_GPU_LAYERS=15  # WSL2 建議從 15 層開始
ollama serve
```

**為什麼這樣有效？**
- Ollama 會分批次載入模型層到 GPU，而不是一次性分配所有內存
- 減少內存碎片化問題
- 即使 GPU 內存顯示為 0，也能正常工作

**GPU 層數建議（WSL2 環境）**：
- **4GB GPU 內存**：10-12 層
- **6GB GPU 內存**：12-15 層
- **8GB GPU 內存**：15-20 層（您的 RTX 4060 建議從 15 開始）
- **16GB+ GPU 內存**：20-30 層

**如果 15 層仍有問題，逐步減少**：
```bash
export OLLAMA_GPU_LAYERS=10  # 再試試 10 層
# 或
export OLLAMA_GPU_LAYERS=5   # 最後嘗試 5 層
```

### 方法 4: 使用 systemd 服務（永久設置）

如果您使用 systemd 管理 Ollama，可以創建環境文件：

1. **創建環境文件**：
```bash
sudo nano /etc/systemd/system/ollama.service.d/override.conf
```

2. **添加以下內容**：
```ini
[Service]
Environment="OLLAMA_NUM_GPU=0"
Environment="OLLAMA_GPU_LAYERS=0"
```

3. **重新載入並重啟**：
```bash
sudo systemctl daemon-reload
sudo systemctl restart ollama
```

## 驗證修復

執行以下命令驗證 Ollama 是否正常運行：

```bash
# 檢查服務狀態
curl http://localhost:11434/api/tags

# 測試模型
ollama run qwen2.5:14b "測試"
```

## 性能說明

### CPU 模式
- ✅ 不會有 CUDA 內存錯誤
- ✅ 穩定可靠
- ⚠️ 速度較慢（但對於健康分析通常足夠）

### 部分 GPU 模式
- ✅ 比 CPU 快
- ✅ 減少內存使用
- ⚠️ 需要根據 GPU 內存調整層數

## 環境變數說明

| 變數 | 說明 | 建議值 |
|------|------|--------|
| `OLLAMA_NUM_GPU` | 使用的 GPU 數量 | `0` = CPU, `1` = 單 GPU |
| `OLLAMA_GPU_LAYERS` | 在 GPU 上運行的層數 | `0` = CPU, `10-30` = 部分 GPU |

## 故障排除

### 如果修復後仍然有問題

1. **檢查日誌**：
```bash
tail -f /tmp/ollama.log
```

2. **檢查是否有其他進程佔用 GPU**：
```bash
nvidia-smi  # 如果有 NVIDIA GPU
```

3. **完全清理並重啟**：
```bash
pkill -9 ollama
rm -rf ~/.ollama/models/*/  # 注意：這會刪除模型，需要重新下載
./fix_ollama_cuda.sh
```

## WSL2 特殊說明

### 為什麼 WSL2 會有這個問題？

1. **GPU 內存映射**：WSL2 使用虛擬化技術，GPU 內存映射與原生 Linux 不同
2. **內存分配策略**：Ollama 默認嘗試一次性分配大塊內存，但 WSL2 可能無法提供足夠的連續內存
3. **CUDA 驅動**：WSL2 的 CUDA 驅動是通過 Windows 驅動轉發的，可能有兼容性問題

### WSL2 最佳實踐

1. **使用分層載入**：設置 `OLLAMA_GPU_LAYERS` 而不是完全禁用 GPU
2. **從較少層數開始**：先試 15 層，如果成功再逐步增加
3. **監控內存使用**：使用 `nvidia-smi` 監控實際 GPU 內存使用情況

### 診斷命令

```bash
# 檢查是否在 WSL2
grep -i microsoft /proc/version

# 檢查 CUDA 是否可用
nvidia-smi

# 檢查 Ollama 日誌
tail -f /tmp/ollama.log

# 測試模型載入
ollama run qwen2.5:14b "測試"
```

## 相關文件

- `scripts/fix_wsl2_gpu.sh` - **WSL2 專用修復腳本（推薦）**
- `scripts/fix_ollama_cuda.sh` - CPU 模式修復腳本
- `scripts/start_ollama.sh` - 啟動腳本（已更新支持環境變數）

