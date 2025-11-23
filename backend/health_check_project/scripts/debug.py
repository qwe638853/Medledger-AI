import logging
from ollama import Client

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

client = Client(host="http://localhost:11434")
prompt = '''
{{
  "name": "test",
  "value": 123
}}
**絕對重要：僅輸出JSON，禁止添加任何說明文字！若無法生成，輸出 {{"success": false}}！**
'''
response = client.chat(model="llama3:8b", messages=[{"role": "user", "content": prompt}], format={"type": "json"})
logger.info(f"回應: {response['message']['content']}")