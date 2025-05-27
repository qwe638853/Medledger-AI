import chromadb
import json

client = chromadb.PersistentClient(path="D:\gg\WOW\chroma_db")
collection = client.get_or_create_collection("medical_knowledge")

# 清空現有數據以避免重複
collection.delete(ids=collection.get()["ids"])

# 健康數據
health_data = {
    "Glu-AC": "89 mg/dL",
    "HbA1c": "4.1 %",
    "Glu-PC": "124 mg/dL",
    "Alb": "4.5 g/dL",
    "TP": "6.5 g/dL",
    "AST（GOT）": "27 U/L",
    "ALT（GPT）": "10 U/L",
    "D-Bil": "0.03 mg/dL",
    "ALP": "74 U/L",
    "T-Bil": "0.7 mg/dL",
    "UN": "23 mg/dL",
    "CRE": "1.2 mg/dL",
    "U.A": "4.9 mg/dL",
    "T-CHO": "164 mg/dL",
    "LDL-C": "128 mg/dL",
    "HDL-C": "54 mg/dL",
    "TG": "143 mg/dL",
    "Hb": "13.3 g/dL",
    "Hct": "49.7 %",
    "PLT": "286 x10^3/uL",
    "WBC": "4.04 x10^3/uL",
    "RBC": "4.56 x10^6/uL",
    "hsCRP": "0.38 mg/dL",
    "AFP": "14 ng/mL",
    "CEA": "2.8 ng/mL",
    "CA-125": "28 U/mL",
    "CA19-9": "29 U/mL",
    "BP": "127/61 mmHg",
    "MCV": "96.2 fL",
    "MCH": "26.3 pg",
    "MCHC": "33.8 g/dL",
    "PT": "10.4 sec",
    "aPTT": "27.8 sec",
    "ESR": "5 mm/hr",
    "RDW-CV": "12.5 %",
    "Specific Gravity": "1.033",
    "PH": "6.0",
    "Protein (Dipstick)": "-",
    "Glucose (Dipstick)": "-",
    "Bilirubin (Dipstick)": "-",
    "Urobilinogen (Dipstick)": "0.4 mg/dL",
    "RBC (Urine)": "0 /HPF",
    "WBC (Urine)": "5 /HPF",
    "Epithelial Cells": "5 /HPF",
    "Casts": "2 /LPF",
    "Ketone": "-",
    "Crystal": "None",
    "Bacteria": "-",
    "Albumin (Dipstick)": "10 mg/L",
    "Creatinine (Dipstick)": "N/A",
    "Alb/CRE Ratio": "10",
    "Nitrite": "-",
    "Occult Blood": "-",
    "WBC Esterase": "-"
}

# 生成文檔和元數據
documents = []
metadatas = []

for key, value in health_data.items():
    # 僅儲存原始健康數據
    document = json.dumps({
        "health_metric": {
            "name": key,
            "value": value
        }
    })
    metadata = {"query": f"{key} {value}"}

    documents.append(document)
    metadatas.append(metadata)

# 將數據添加到 ChromaDB
collection.add(
    ids=[str(i) for i in range(len(documents))],
    documents=documents,
    metadatas=metadatas
)

print("Raw health data added to 'medical_knowledge' collection for LLM analysis!")