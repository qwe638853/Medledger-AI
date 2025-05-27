from langchain_chroma import Chroma
from langchain_huggingface import HuggingFaceEmbeddings

embedding_model = HuggingFaceEmbeddings(model_name="sentence-transformers/all-MiniLM-L6-v2")
vectorstore = Chroma(
    embedding_function=embedding_model,
    collection_name="health_knowledge",
    persist_directory="./chroma_db"
)
medical_docs = vectorstore.get(where={"source": "medicalgpt"})
print(f"MedicalGPT 數據數量: {len(medical_docs.get('documents', []))}")