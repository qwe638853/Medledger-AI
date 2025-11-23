import os
import sys
# 確保 Python 能找到專案根目錄的套件 (如果有的話)
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from langchain_chroma import Chroma
from langchain_huggingface import HuggingFaceEmbeddings
from langchain_community.document_loaders import TextLoader
from langchain_text_splitters import MarkdownHeaderTextSplitter

# --- 路徑設定 ---
BASE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
KNOWLEDGE_FILE = os.path.join(BASE_DIR, "data", "health_rules.txt")
CHROMA_DIR = os.path.join(BASE_DIR, "chroma_db")
EMBEDDING_MODEL = "intfloat/multilingual-e5-large"

def add_knowledge():
    print(f"🔍 除錯資訊：")
    print(f"  - 專案根目錄 (BASE_DIR): {BASE_DIR}")
    print(f"  - 規則文件路徑 (KNOWLEDGE_FILE): {KNOWLEDGE_FILE}")
    print(f"  - 資料庫路徑 (CHROMA_DIR): {CHROMA_DIR}")

    # 1. 檢查目錄
    if not os.path.exists(CHROMA_DIR):
        print(f"📂 資料庫目錄不存在，正在建立: {CHROMA_DIR}")
        os.makedirs(CHROMA_DIR)
    else:
        print(f"📂 資料庫目錄已存在: {CHROMA_DIR}")
        
    # 2. 檢查文件
    if not os.path.exists(KNOWLEDGE_FILE):
        print(f"❌ 錯誤：找不到知識文件！路徑: {KNOWLEDGE_FILE}")
        print(f"請確認您是否已執行 'mv health_rules.txt data/'")
        return

    # 3. 載入模型 (這步最容易卡住)
    print(f"🚀 正在初始化 RAG 嵌入模型 ({EMBEDDING_MODEL})... (如果是第一次執行，下載可能需要幾分鐘，請耐心等待)")
    try:
        embeddings = HuggingFaceEmbeddings(model_name=EMBEDDING_MODEL)
        print("✅ 模型載入成功！")
    except Exception as e:
        print(f"❌ 模型載入失敗: {e}")
        return
    
    # 4. 讀取文件
    print(f"📄 正在讀取文件內容...")
    try:
        with open(KNOWLEDGE_FILE, 'r', encoding='utf-8') as f:
            full_text = f.read()
        print(f"✅ 文件讀取成功，大小: {len(full_text)} 字元")
    except Exception as e:
        print(f"❌ 讀取文件失敗: {e}")
        return

    # 5. 切割文件
    print("✂️ 正在切割文件...")
    headers_to_split_on = [("###", "TopicHeader")]
    markdown_splitter = MarkdownHeaderTextSplitter(
        headers_to_split_on=headers_to_split_on,
        strip_headers=False
    )
    docs = markdown_splitter.split_text(full_text)
    print(f"📚 切割完成，共 {len(docs)} 個知識塊")
    
    # 6. 建立索引
    print(f"�� 正在寫入 ChromaDB 資料庫...")
    try:
        # 清除舊資料
        try:
            Chroma(persist_directory=CHROMA_DIR, embedding_function=embeddings, collection_name="health_knowledge").delete_collection()
            print("  - 已清除舊索引")
        except:
            print("  - 無舊索引需清除")

        # 寫入新資料
        Chroma.from_documents(
            docs, 
            embeddings, 
            persist_directory=CHROMA_DIR, 
            collection_name="health_knowledge"
        )
        print(f"✅ 知識索引全部完成！資料庫位於: {CHROMA_DIR}")
    except Exception as e:
        print(f"❌ 資料庫寫入失敗: {e}")

if __name__ == "__main__":
    add_knowledge()

