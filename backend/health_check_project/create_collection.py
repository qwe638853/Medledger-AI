import chromadb
client = chromadb.PersistentClient(path="D:\\gg\\WOW\\chroma_db")
collections = client.list_collections()
print([col.name for col in collections])