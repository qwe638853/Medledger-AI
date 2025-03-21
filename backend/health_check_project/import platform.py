import pymssql
from langchain_ollama import OllamaLLM
from langchain.prompts import PromptTemplate
from langchain_core.memory import ConversationBufferMemory
from fastapi import HTTPException, UploadFile, File
from pydantic import BaseModel
from typing import Optional, Dict, Any
from passlib.context import CryptContext
import PyPDF2
from docx import Document
from io import BytesIO
from datetime import datetime
import logging
import pyodbc
import platform

# 設定日誌
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# 初始化密碼加密工具
pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")

# 定義請求體結構（用於註冊）
class RegisterRequest(BaseModel):
    full_name: str
    gender: str
    birth_date: str
    id_number: str
    password: str
    phone_number: str
    email: str

# 定義請求體結構（用於登入）
class LoginRequest(BaseModel):
    id_number: str
    password: str

# 定義請求體結構（用於忘記密碼）
class ForgotPasswordRequest(BaseModel):
    id_number: str
    email: str

# 定義請求體結構（用於互動模式）
class InteractiveRequest(BaseModel):
    query: str

# 連線配置
db_config = {
    'server': 'healthdbserver123.database.windows.net',
    'port': 1433,
    'user': 'bojay@healthdbserver123',
    'password': '!Aa1085121',
    'database': 'health_db'
}

# 連接到 Azure SQL Database 的通用函數
def connect_to_db():
    try:
        logger.info(f"嘗試連線到 {db_config['server']}，用戶: {db_config['user']}")
        drivers = pyodbc.drivers()
        logger.info(f"可用的 ODBC 驅動程式: {drivers}")
        
        # 檢查 Python 位元
        logger.info(f"Python 位元: {platform.architecture()[0]}")
        
        # 優先使用 ODBC Driver 18 for SQL Server
        driver = "ODBC Driver 18 for SQL Server"
        if driver not in drivers:
            logger.warning(f"未找到 {driver}，嘗試使用 ODBC Driver 17 for SQL Server")
            driver = "ODBC Driver 17 for SQL Server"
            if driver not in drivers:
                logger.warning(f"未找到 {driver}，嘗試使用其他驅動程式")
                for d in drivers:
                    if "SQL Server" in d:
                        driver = d
                        logger.info(f"改用驅動程式: {driver}")
                        break
                else:
                    raise HTTPException(status_code=500, detail="未找到合適的 ODBC 驅動程式，請確保已安裝 ODBC Driver 17 或 18 for SQL Server")

        conn = pyodbc.connect(
            f"DRIVER={{{driver}}};"
            f"SERVER=tcp:{db_config['server']},1433;"
            f"DATABASE={db_config['database']};"
            f"UID={db_config['user']};"
            f"PWD={db_config['password']};"
            f"Encrypt=yes;"
            f"TrustServerCertificate=yes;"
            f"Connection Timeout=30;"
        )
        logger.info("資料庫連線成功 (pyodbc)")
        return conn
    except Exception as e:
        logger.error(f"資料庫連線失敗 (pyodbc): {e}")
        try:
            conn = pymssql.connect(
                server=db_config['server'],
                port=1433,
                user=db_config['user'],
                password=db_config['password'],
                database=db_config['database'],
                login_timeout=30,
                charset='UTF-8',
                tds_version='7.4'
            )
            logger.info("資料庫連線成功 (pymssql)")
            return conn
        except Exception as e2:
            logger.error(f"資料庫連線失敗 (pymssql): {e2}")
            raise HTTPException(status_code=500, detail=f"連線 Azure SQL Database 時發生錯誤：{e2}")

# 檢查身分證字號是否已存在
def check_id_number_exists(id_number: str) -> bool:
    try:
        conn = connect_to_db()
        cursor = conn.cursor()
        query = "SELECT COUNT(*) FROM users WHERE id_number = ? AND id_number != ''"
        cursor.execute(query, (id_number,))
        count = cursor.fetchone()[0]
        conn.close()
        return count > 0
    except Exception as e:
        logger.error(f"檢查身分證字號是否存在時發生錯誤: {e}")
        raise HTTPException(status_code=500, detail=f"檢查身分證字號是否存在時發生錯誤：{e}")

# 驗證身分證字號格式
def validate_id_number(id_number: str) -> bool:
    if len(id_number) != 10:
        return False
    if not id_number[0].isalpha() or not id_number[1:].isdigit():
        return False
    return True

# 驗證手機號碼格式
def validate_phone_number(phone_number: str) -> bool:
    return len(phone_number) == 10 and phone_number.isdigit()

# 驗證電子郵件格式
def validate_email(email: str) -> bool:
    return "@" in email and "." in email and len(email) <= 100

# 提取 PDF 文件的文本
def extract_text_from_pdf(file: BytesIO) -> str:
    try:
        reader = PyPDF2.PdfReader(file)
        text = ""
        for page in reader.pages:
            text += page.extract_text() or ""
        return text
    except Exception as e:
        raise HTTPException(status_code=400, detail=f"無法提取 PDF 文件內容：{e}")

# 提取 Word 文件的文本
def extract_text_from_docx(file: BytesIO) -> str:
    try:
        doc = Document(file)
        text = ""
        for paragraph in doc.paragraphs:
            text += paragraph.text + "\n"
        return text
    except Exception as e:
        raise HTTPException(status_code=400, detail=f"無法提取 Word 文件內容：{e}")

# 根據 id_number 查找 id_number
def get_user_id_by_id_number(id_number: str) -> Optional[str]:
    if not validate_id_number(id_number):
        raise HTTPException(status_code=400, detail="身分證字號格式不正確，應為 1 個字母 + 9 個數字，例如 A123456789")
    try:
        conn = connect_to_db()
        cursor = conn.cursor()
        query = "SELECT id_number FROM users WHERE id_number = ? AND id_number != ''"
        cursor.execute(query, (id_number,))
        user = cursor.fetchone()
        conn.close()
        return user[0] if user else None
    except Exception as e:
        logger.error(f"查詢身分證字號時發生錯誤: {e}")
        raise HTTPException(status_code=500, detail=f"查詢身分證字號時發生錯誤：{e}")

# 提取健檢資料
def extract_health_data(id_number: str):
    try:
        conn = connect_to_db()
        cursor = conn.cursor()
        cursor.execute("SELECT COUNT(*) FROM users WHERE id_number = ?", (id_number,))
        if cursor.fetchone()[0] == 0:
            conn.close()
            return {"health_data": None, "error": f"id_number={id_number} 未在 users 表中註冊"}
        query = """
        SELECT hc.check_date, hc.extracted_text, hc.notes
        FROM health_checks hc
        WHERE hc.id_number = ?
        ORDER BY hc.check_date DESC
        """
        cursor.execute(query, (id_number,))
        record = cursor.fetchone()
        conn.close()
        if not record:
            return {"health_data": None, "error": f"找不到 id_number={id_number} 的健檢資料"}
        health_data = f"健檢日期: {record[0]}\n"
        health_data += f"健檢資料: {record[1] if record[1] else '無'}\n"
        health_data += f"建議: {record[2] if record[2] else '無'}\n"
        health_data += f"身分證字號: {id_number}\n"
        return {"health_data": health_data, "error": None}
    except Exception as e:
        logger.error(f"提取健檢資料時發生錯誤: {e}")
        return {"health_data": None, "error": f"連線錯誤：{e}"}

# 分析健檢資料
def analyze_health_data(id_number: str):
    memory = ConversationBufferMemory()
    llm = OllamaLLM(model="llama3:8b", base_url="http://localhost:11434", temperature=0.3)
    analysis_prompt_template = PromptTemplate(
        input_variables=["data"],
        template="你是一個專業的健康分析專家，請以繁體中文分析以下健檢資料並提供建議：\n\n{data}\n\n分析每個指標是否正常並給出至少三項建議和一項潛在風險。"
    )
    interactive_prompt_template = PromptTemplate(
        input_variables=["query"],
        template="你是一個專業的健康分析專家，根據之前的健檢資料，回答以下問題：\n\n{query}"
    )
    health_data_response = extract_health_data(id_number)
    health_data = health_data_response["health_data"]
    if health_data_response["error"]:
        return None, llm, memory, interactive_prompt_template
    memory.save_context({"input": "健檢資料"}, {"output": health_data})
    try:
        analysis_prompt = analysis_prompt_template.format(data=health_data)
        result = llm.invoke(analysis_prompt)
        return {"health_data": health_data, "analysis_result": result, "id_number": id_number}, llm, memory, interactive_prompt_template
    except Exception as e:
        logger.error(f"分析健檢資料時發生錯誤: {e}")
        return None, llm, memory, interactive_prompt_template

# 註冊新用戶
def register_user(user_data: RegisterRequest):
    if not validate_id_number(user_data.id_number):
        raise HTTPException(status_code=400, detail="身分證字號格式不正確")
    if check_id_number_exists(user_data.id_number):
        raise HTTPException(status_code=400, detail="此身分證字號已註冊")
    if len(user_data.password) < 6 or len(user_data.password) > 100:
        raise HTTPException(status_code=400, detail="密碼長度無效")
    if not validate_phone_number(user_data.phone_number):
        raise HTTPException(status_code=400, detail="手機號碼格式不正確")
    if not validate_email(user_data.email):
        raise HTTPException(status_code=400, detail="電子郵件格式不正確")
    if not user_data.full_name or len(user_data.full_name) > 100:
        raise HTTPException(status_code=400, detail="姓名無效")
    if user_data.gender not in ["M", "F", "Other"]:
        raise HTTPException(status_code=400, detail="性別無效")
    try:
        datetime.strptime(user_data.birth_date, "%Y-%m-%d")
    except ValueError:
        raise HTTPException(status_code=400, detail="出生日期格式不正確")

    hashed_password = pwd_context.hash(user_data.password)

    try:
        conn = connect_to_db()
        cursor = conn.cursor()
        query = """
        INSERT INTO users (full_name, gender, birth_date, id_number, password, phone_number, email)
        VALUES (?, ?, ?, ?, ?, ?, ?)
        """
        cursor.execute(query, (
            user_data.full_name, user_data.gender, user_data.birth_date, user_data.id_number,
            hashed_password, user_data.phone_number, user_data.email
        ))
        conn.commit()
        conn.close()
        return {"message": "註冊成功"}
    except Exception as e:
        logger.error(f"註冊用戶時發生錯誤: {e}")
        raise HTTPException(status_code=500, detail=f"註冊錯誤：{e}")

# 登入驗證
def login_user(login_data: LoginRequest):
    if not validate_id_number(login_data.id_number):
        raise HTTPException(status_code=400, detail="身分證字號格式不正確")
    if not login_data.password or len(login_data.password) > 100:
        raise HTTPException(status_code=400, detail="密碼無效")

    try:
        conn = connect_to_db()
        cursor = conn.cursor()
        query = "SELECT id, password FROM users WHERE id_number = ?"
        cursor.execute(query, (login_data.id_number,))
        user = cursor.fetchone()
        conn.close()

        if not user:
            raise HTTPException(status_code=401, detail="用戶不存在")
        if not pwd_context.verify(login_data.password, user[1]):
            raise HTTPException(status_code=401, detail="密碼錯誤")

        return {"message": "登入成功", "id_number": login_data.id_number}
    except Exception as e:
        logger.error(f"登入時發生錯誤: {e}")
        raise HTTPException(status_code=500, detail=f"登入錯誤：{e}")

# 忘記密碼功能
def forgot_password(forgot_data: ForgotPasswordRequest):
    if not validate_id_number(forgot_data.id_number):
        raise HTTPException(status_code=400, detail="身分證字號格式不正確")
    if not validate_email(forgot_data.email):
        raise HTTPException(status_code=400, detail="電子郵件格式不正確")

    try:
        conn = connect_to_db()
        cursor = conn.cursor()
        query = "SELECT id FROM users WHERE id_number = ? AND email = ?"
        cursor.execute(query, (forgot_data.id_number, forgot_data.email))
        user = cursor.fetchone()
        if not user:
            conn.close()
            raise HTTPException(status_code=404, detail="用戶不存在")

        temp_password = "temp123456"
        hashed_temp_password = pwd_context.hash(temp_password)
        query = "UPDATE users SET password = ? WHERE id = ?"
        cursor.execute(query, (hashed_temp_password, user[0]))
        conn.commit()
        conn.close()
        return {"message": "密碼重置成功", "temp_password": temp_password}
    except Exception as e:
        logger.error(f"重置密碼時發生錯誤: {e}")
        raise HTTPException(status_code=500, detail=f"重置密碼錯誤：{e}")

# 上傳健檢資料文件
async def upload_health_check(id_number: str, file: UploadFile):
    if file.content_type not in ["application/pdf", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"]:
        raise HTTPException(status_code=400, detail="僅支援 PDF 和 Word 文件")
    file_data = await file.read()
    file_extension = file.filename.split('.')[-1].lower()
    file_stream = BytesIO(file_data)
    if file_extension == "pdf":
        extracted_text = extract_text_from_pdf(file_stream)
    elif file_extension == "docx":
        extracted_text = extract_text_from_docx(file_stream)
    else:
        raise HTTPException(status_code=400, detail="不支援的文件格式")
    if not extracted_text.strip():
        raise HTTPException(status_code=400, detail="文件內容為空")

    try:
        conn = connect_to_db()
        cursor = conn.cursor()
        query_check = "SELECT COUNT(*) FROM health_checks WHERE id_number = ?"
        cursor.execute(query_check, (id_number,))
        exists = cursor.fetchone()[0] > 0
        if exists:
            query = """
            UPDATE health_checks SET check_date = ?, data = ?, file_data = ?, extracted_text = ?
            WHERE id_number = ?
            """
            cursor.execute(query, (datetime.now().strftime("%Y-%m-%d"), extracted_text, file_data, extracted_text, id_number))
            message = "健檢資料更新成功"
        else:
            query = """
            INSERT INTO health_checks (id_number, check_date, data, file_data, extracted_text)
            VALUES (?, ?, ?, ?, ?)
            """
            cursor.execute(query, (id_number, datetime.now().strftime("%Y-%m-%d"), extracted_text, file_data, extracted_text))
            message = "健檢資料上傳成功"
        conn.commit()
        conn.close()
        return {"message": message, "extracted_text": extracted_text}
    except Exception as e:
        logger.error(f"上傳健檢資料時發生錯誤: {e}")
        raise HTTPException(status_code=500, detail=f"上傳錯誤：{e}")

# 測試資料庫連線
def test_db_connection():
    try:
        conn = connect_to_db()
        cursor = conn.cursor()
        cursor.execute("SELECT 1")
        result = cursor.fetchone()
        conn.close()
        logger.info(f"資料庫連線測試成功: {result}")
        return True
    except Exception as e:
        logger.error(f"資料庫連線測試失敗: {e}")
        return False

# 測試 Ollama 連線
def test_ollama_connection():
    try:
        llm = OllamaLLM(model="llama3:8b", base_url="http://localhost:11434", temperature=0.3)
        result = llm.invoke("你好，這是一個測試。")
        logger.info(f"Ollama 連線測試成功: {result}")
        return True
    except Exception as e:
        logger.error(f"Ollama 連線測試失敗: {e}")
        return False

if __name__ == "__main__":
    test_db_connection()
    test_ollama_connection()