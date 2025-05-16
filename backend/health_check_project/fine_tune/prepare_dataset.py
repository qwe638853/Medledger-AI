from datasets import Dataset
import json
import logging

# 設定日誌
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# 載入 medical_dialogue_traditional.json
medical_dialogue_path = "D:/gg/WOW/medical_dialogue_traditional.json"
with open(medical_dialogue_path, "r", encoding="utf-8") as f:
    medical_dialogue_data = json.load(f)
logger.info(f"載入 medical_dialogue_traditional.json，包含 {len(medical_dialogue_data)} 筆數據")

# 載入 Huatuo 檔案（假設路徑為 huatuo_medical_qa.json）
huatuo_path = "D:/gg/WOW/huatuo_medical_qa_sharegpt_traditional.json"  # 替換為實際路徑
with open(huatuo_path, "r", encoding="utf-8") as f:
    huatuo_data = json.load(f)
logger.info(f"載入 Huatuo 檔案，包含 {len(huatuo_data)} 筆數據")

# 合併數據
combined_data = []
skipped_entries = 0

# 處理 medical_dialogue_traditional.json
for item in medical_dialogue_data:
    question = item.get("question", "")
    answer = item.get("answer", "")
    if not question or not answer:
        skipped_entries += 1
        continue
    combined_data.append({
        "input": question,
        "output": answer,
        "instruction": "請根據問題提供專業的醫療建議，用繁體中文回答。"
    })

# 處理 Huatuo 檔案
for item in huatuo_data:
    question = item.get("question", "")
    answer = item.get("answer", "")
    if not question or not answer:
        skipped_entries += 1
        continue
    combined_data.append({
        "input": question,
        "output": answer,
        "instruction": "請根據問題提供專業的醫療建議，用繁體中文回答。"
    })

logger.info(f"合併後數據量：{len(combined_data)} 筆，跳過 {skipped_entries} 筆無效數據")

# 轉換為 Hugging Face Dataset 格式
dataset = Dataset.from_list(combined_data)

# 分割訓練和驗證集（可選）
train_test_split = dataset.train_test_split(test_size=0.1)
train_dataset = train_test_split["train"]
test_dataset = train_test_split["test"]

# 保存數據集
output_dir = "D:/gg/WOW/backend/health_check_project/fine_tune/combined_medical_dataset"
train_dataset.save_to_disk(f"{output_dir}/train")
test_dataset.save_to_disk(f"{output_dir}/test")
logger.info(f"數據集已保存到 {output_dir}")