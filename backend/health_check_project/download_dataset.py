import pandas as pd
import json
import os
import glob
from convert_to_traditional import convert_to_traditional

# 獲取當前腳本所在的目錄
script_dir = os.path.dirname(os.path.abspath(__file__))
print(f"腳本所在目錄: {script_dir}")

# 定義數據目錄（向上兩級到 d:\gg\WOW，然後進入 chinese-medical-dialogue-data/Data_数据）
data_dir = os.path.join(script_dir, "..", "..", "Chinese-medical-dialogue-data", "Data_数据")
data_dir = os.path.abspath(data_dir).replace("Chinese-medical-dialogue-data", "chinese-medical-dialogue-data")
print(f"數據目錄: {data_dir}")

# 列出數據目錄下的所有子目錄，檢查是否存在
parent_dir = os.path.dirname(data_dir)
print(f"父目錄: {parent_dir}")
if os.path.exists(parent_dir):
    print(f"父目錄下的子目錄: {os.listdir(parent_dir)}")
else:
    print(f"父目錄 {parent_dir} 不存在！")

# 確認數據目錄是否存在
if not os.path.exists(data_dir):
    print(f"數據目錄 {data_dir} 不存在，無法處理。")
    exit()

# 尋找所有 CSV 檔案（包括子目錄）
csv_pattern = os.path.join(data_dir, "**/*.csv")
print(f"搜尋模式: {csv_pattern}")
csv_files = glob.glob(csv_pattern, recursive=True)
print(f"找到的 CSV 檔案: {csv_files}")

if not csv_files:
    print(f"在 {data_dir} 中找不到任何 CSV 檔案，無法處理。")
    exit()

# 儲存所有對話數據
all_dialogue_data = []

# 處理每個 CSV 檔案
for csv_file in csv_files:
    try:
        print(f"開始處理 {csv_file}...")
        # 嘗試多種編碼，優先使用 gb18030
        encodings = ['gb18030', 'GB2312', 'utf-8-sig', 'gbk', 'utf-16', 'utf-8']
        df = None
        for encoding in encodings:
            try:
                df = pd.read_csv(csv_file, encoding=encoding)
                print(f"成功使用 {encoding} 編碼讀取 {csv_file}")
                break
            except Exception as e:
                print(f"嘗試 {encoding} 編碼失敗: {str(e)}")
                continue
        if df is None:
            print(f"無法讀取 {csv_file}，跳過此檔案")
            continue

        # 獲取檔案所在的子目錄名稱作為來源（例如 IM、Pediatric）
        source = os.path.basename(os.path.dirname(csv_file))

        # 轉換為問答格式並轉為繁體中文
        for idx, row in df.iterrows():
            question = convert_to_traditional(row["ask"])
            answer = convert_to_traditional(row["answer"])
            all_dialogue_data.append({
                "question": question,
                "answer": answer,
                "source": source  # 添加來源標記
            })

        print(f"完成處理 {csv_file}")
    except Exception as e:
        print(f"處理 {csv_file} 時發生錯誤: {str(e)}")
        continue

# 儲存為 JSON
if all_dialogue_data:
    with open("medical_dialogue_traditional.json", "w", encoding="utf-8") as f:
        json.dump(all_dialogue_data, f, ensure_ascii=False, indent=2)
    print("簡繁轉換完成：medical_dialogue_traditional.json")
else:
    print("沒有成功處理任何數據，無法生成 JSON 檔案。")