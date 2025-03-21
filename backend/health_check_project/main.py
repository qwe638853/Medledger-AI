import pymssql
from langchain_ollama import OllamaLLM
from langchain.prompts import PromptTemplate
from langchain.memory import ConversationBufferMemory
from pydantic import BaseModel
from typing import Optional, Dict, Any, List, Tuple
from passlib.context import CryptContext
import PyPDF2
from docx import Document
from io import BytesIO
from datetime import datetime
import logging

# 初始化密碼加密工具
pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")

# 設定日誌
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# 定義角色常量
ROLE_USER = "user"
ROLE_OTHER = "other"
ROLE_HEALTH_CENTER = "health_center"
VALID_ROLES = [ROLE_USER, ROLE_OTHER, ROLE_HEALTH_CENTER]

# 自定義表單類，用於簡化 OAuth2 登入表單
class CustomOAuth2PasswordRequestForm(BaseModel):
    username: str
    password: str
    scope: Optional[str] = None  # 將 scope 設為可選

    class Config:
        json_schema_extra = {
            "example": {
                "username": "A123456789",
                "password": "testpassword",
                "scope": "role:user"
            }
        }

# 定義請求體結構（用於註冊）
class RegisterRequest(BaseModel):
    full_name: str
    gender: str
    birth_date: str  # 格式：YYYY-MM-DD
    id_number: str
    password: str
    phone_number: str
    email: str
    role: str

# 定義請求體結構（用於登入）
class LoginRequest(BaseModel):
    id_number: str
    password: str
    role: str

# 定義請求體結構（用於忘記密碼）
class ForgotPasswordRequest(BaseModel):
    id_number: str
    email: str

# 定義請求體結構（用於互動模式）
class InteractiveRequest(BaseModel):
    query: str

# 連接到 Azure SQL Database 的通用函數
def connect_to_db(db_config: Dict[str, Any]):
    try:
        conn = pymssql.connect(
            server=db_config['server'],
            port=db_config['port'],
            user=db_config['user'],
            password=db_config['password'],
            database=db_config['database'],
            login_timeout=30,
            charset='UTF-8'
        )
        logger.info("成功連接到 Azure SQL Database")
        return conn
    except Exception as e:
        logger.error(f"連線 Azure SQL Database 時發生錯誤: {str(e)}")
        raise Exception(f"連線 Azure SQL Database 時發生錯誤: {str(e)}")

# 檢查身分證字號是否已存在
def check_id_number_exists(db_config: Dict[str, Any], id_number: str) -> bool:
    try:
        conn = connect_to_db(db_config)
        cursor = conn.cursor()
        query = "SELECT COUNT(*) FROM users WHERE id_number = %s AND id_number != ''"
        cursor.execute(query, (id_number,))
        count = cursor.fetchone()[0]
        conn.close()
        logger.info(f"檢查身分證字號是否存在: id_number={id_number}, 存在={count > 0}")
        return count > 0
    except Exception as e:
        logger.error(f"檢查身分證字號是否存在時發生錯誤: id_number={id_number}, 錯誤: {str(e)}")
        raise Exception(f"檢查身分證字號是否存在時發生錯誤: {str(e)}")

# 驗證身分證字號格式（改進版，提供詳細錯誤訊息）
def validate_id_number(id_number: str) -> Tuple[bool, Optional[str]]:
    if not id_number:
        return False, "請輸入身分證字號"

    if len(id_number) != 10:
        return False, "身分證字號長度必須為 10 個字元"

    if not id_number[0].isalpha() or not id_number[0].isupper():
        return False, "身分證字號第一個字元必須為大寫字母"

    if id_number[1] not in ['1', '2']:
        return False, "身分證字號第二個字元必須為 1 或 2"

    if not id_number[2:].isdigit():
        return False, "身分證字號後 8 個字元必須為數字"

    letter_map = {
        'A': 10, 'B': 11, 'C': 12, 'D': 13, 'E': 14, 'F': 15, 'G': 16, 'H': 17,
        'I': 34, 'J': 18, 'K': 19, 'L': 20, 'M': 21, 'N': 22, 'O': 35, 'P': 23,
        'Q': 24, 'R': 25, 'S': 26, 'T': 27, 'U': 28, 'V': 29, 'W': 32, 'X': 30,
        'Y': 31, 'Z': 33
    }

    first_char_value = letter_map.get(id_number[0])
    if first_char_value is None:
        return False, "身分證字號第一個字元無效"

    total = (first_char_value // 10) + (first_char_value % 10) * 9
    weights = [8, 7, 6, 5, 4, 3, 2, 1]
    for i in range(8):
        total += int(id_number[i + 1]) * weights[i]

    check_digit = (10 - (total % 10)) % 10
    if check_digit != int(id_number[9]):
        return False, "身分證字號檢查碼錯誤"

    return True, None

# 驗證手機號碼格式
def validate_phone_number(phone_number: str) -> bool:
    return len(phone_number) == 10 and phone_number.isdigit()

# 驗證電子郵件格式
def validate_email(email: str) -> bool:
    return "@" in email and "." in email and len(email) <= 100

# 驗證角色
def validate_role(role: str) -> bool:
    return role in VALID_ROLES

# 提取 PDF 文件的文本
def extract_text_from_pdf(file: BytesIO) -> str:
    try:
        reader = PyPDF2.PdfReader(file)
        text = ""
        for page in reader.pages:
            page_text = page.extract_text()
            text += page_text or ""
        logger.info(f"成功從 PDF 文件提取文本，長度: {len(text)}")
        return text
    except Exception as e:
        logger.error(f"無法提取 PDF 文件內容: {str(e)}")
        raise Exception(f"無法提取 PDF 文件內容: {str(e)}")

# 提取 Word 文件的文本
def extract_text_from_docx(file: BytesIO) -> str:
    try:
        doc = Document(file)
        text = ""
        for paragraph in doc.paragraphs:
            text += paragraph.text + "\n"
        logger.info(f"成功從 Word 文件提取文本，長度: {len(text)}")
        return text
    except Exception as e:
        logger.error(f"無法提取 Word 文件內容: {str(e)}")
        raise Exception(f"無法提取 Word 文件內容: {str(e)}")

# 根據 id_number 查找用戶資訊
def get_user_info(db_config: Dict[str, Any], id_number: str) -> Optional[Dict[str, Any]]:
    is_valid, error_message = validate_id_number(id_number)
    if not is_valid:
        logger.error(f"身分證字號格式不正確: id_number={id_number}, 錯誤: {error_message}")
        raise Exception(f"身分證字號格式不正確: {error_message}")
    try:
        conn = connect_to_db(db_config)
        cursor = conn.cursor()
        query = """
        SELECT id, full_name, gender, birth_date, id_number, password, 
               phone_number, email, role, created_at 
        FROM users 
        WHERE id_number = %s AND id_number != ''
        """
        cursor.execute(query, (id_number,))
        user = cursor.fetchone()
        conn.close()
        if user:
            user_info = {
                "id": user[0],
                "full_name": user[1],
                "gender": user[2],
                "birth_date": user[3],
                "id_number": user[4],
                "password": user[5],
                "phone_number": user[6],
                "email": user[7],
                "role": user[8],
                "created_at": user[9]
            }
            logger.info(f"成功查詢用戶: id_number={id_number}, role={user_info['role']}")
            return user_info
        logger.warning(f"用戶不存在: id_number={id_number}")
        return None
    except Exception as e:
        logger.error(f"查詢用戶資訊時發生錯誤: id_number={id_number}, 錯誤: {str(e)}")
        raise Exception(f"查詢用戶資訊時發生錯誤: {str(e)}")

# 註冊新用戶
def register_user(db_config: Dict[str, Any], user_data: RegisterRequest):
    is_valid, error_message = validate_id_number(user_data.id_number)
    if not is_valid:
        raise Exception(f"身分證字號格式不正確: {error_message}")
    if check_id_number_exists(db_config, user_data.id_number):
        raise Exception("此身分證字號已註冊，請確認是否輸入錯誤")
    if len(user_data.password) < 6:
        raise Exception("密碼長度必須至少為 6 個字元")
    if len(user_data.password) > 100:
        raise Exception("密碼長度不得超過 100 個字元")
    if not validate_phone_number(user_data.phone_number):
        raise Exception("手機號碼格式不正確，應為 10 位數字，例如 0912345678")
    if not validate_email(user_data.email):
        raise Exception("電子郵件格式不正確，例如 user@example.com")
    if not user_data.full_name or len(user_data.full_name) > 100:
        raise Exception("姓名必須提供且長度不得超過 100 個字元")
    if user_data.gender not in ["M", "F", "Other"]:
        raise Exception("性別必須是 'M'、'F' 或 'Other'")
    if not validate_role(user_data.role):
        raise Exception(f"角色無效，必須是 {', '.join(VALID_ROLES)} 之一")
    try:
        datetime.strptime(user_data.birth_date, "%Y-%m-%d")
    except ValueError:
        raise Exception("出生日期格式不正確，應為 YYYY-MM-DD")

    hashed_password = pwd_context.hash(user_data.password)

    try:
        conn = connect_to_db(db_config)
        cursor = conn.cursor()
        query = """
        INSERT INTO users (full_name, gender, birth_date, id_number, password, phone_number, email, role, created_at)
        VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s)
        """
        cursor.execute(query, (
            user_data.full_name,
            user_data.gender,
            user_data.birth_date,
            user_data.id_number,
            hashed_password,
            user_data.phone_number,
            user_data.email,
            user_data.role,
            datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        ))
        conn.commit()
        conn.close()
        logger.info(f"用戶註冊成功: id_number={user_data.id_number}, role={user_data.role}")
        return {"message": "註冊成功"}
    except Exception as e:
        logger.error(f"註冊過程中發生錯誤: id_number={user_data.id_number}, 錯誤: {str(e)}")
        raise Exception(f"註冊過程中發生錯誤: {str(e)}")

# 登入驗證
def login_user(db_config: Dict[str, Any], login_data: LoginRequest):
    is_valid, error_message = validate_id_number(login_data.id_number)
    if not is_valid:
        raise Exception(f"身分證字號格式不正確: {error_message}")
    if not login_data.password or len(login_data.password) > 100:
        raise Exception("密碼無效")
    if not validate_role(login_data.role):
        raise Exception(f"角色無效，必須是 {', '.join(VALID_ROLES)} 之一")

    try:
        user = get_user_info(db_config, login_data.id_number)
        if not user:
            raise Exception("身分證字號不存在")

        stored_password = user["password"]
        if not stored_password or not stored_password.startswith('$2b$'):
            raise Exception(f"資料庫中的密碼格式無效，需聯繫管理員修復 (id_number: {login_data.id_number})")

        if not pwd_context.verify(login_data.password, stored_password):
            raise Exception("密碼錯誤")

        if user["role"] != login_data.role:
            raise Exception(f"角色不匹配，該用戶的角色為 {user['role']}，但嘗試以 {login_data.role} 登入")

        logger.info(f"用戶登入成功: id_number={login_data.id_number}, role={user['role']}")
        return {"message": "登入成功", "id_number": login_data.id_number, "role": user["role"]}
    except Exception as e:
        logger.error(f"登入過程中發生錯誤: id_number={login_data.id_number}, 錯誤: {str(e)}")
        raise Exception(f"登入過程中發生錯誤: {str(e)}")

# 忘記密碼功能
def forgot_password(db_config: Dict[str, Any], forgot_data: ForgotPasswordRequest):
    is_valid, error_message = validate_id_number(forgot_data.id_number)
    if not is_valid:
        raise Exception(f"身分證字號格式不正確: {error_message}")
    if not validate_email(forgot_data.email):
        raise Exception("電子郵件格式不正確，例如 user@example.com")

    try:
        conn = connect_to_db(db_config)
        cursor = conn.cursor()
        query = "SELECT id FROM users WHERE id_number = %s AND email = %s"
        cursor.execute(query, (forgot_data.id_number, forgot_data.email))
        user = cursor.fetchone()

        if not user:
            conn.close()
            raise Exception("身分證字號與電子郵件不匹配")

        temp_password = "temp123456"
        hashed_temp_password = pwd_context.hash(temp_password)

        query = "UPDATE users SET password = %s WHERE id = %s"
        cursor.execute(query, (hashed_temp_password, user[0]))
        conn.commit()
        conn.close()

        logger.info(f"密碼重置成功: id_number={forgot_data.id_number}")
        return {"message": "密碼已重置，請使用臨時密碼登入並更改密碼", "temp_password": temp_password}
    except Exception as e:
        logger.error(f"重置密碼時發生錯誤: id_number={forgot_data.id_number}, 錯誤: {str(e)}")
        raise Exception(f"重置密碼時發生錯誤: {str(e)}")

# 上傳健檢資料文件（僅限健檢中心）
async def upload_health_check(db_config: Dict[str, Any], id_number: str, file):
    if file.content_type not in ["application/pdf", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"]:
        raise Exception("僅支援 PDF 和 Word (.docx) 文件")

    file_data = await file.read()
    file_extension = file.filename.split('.')[-1].lower()
    file_stream = BytesIO(file_data)
    if file_extension == "pdf":
        extracted_text = extract_text_from_pdf(file_stream)
    elif file_extension == "docx":
        extracted_text = extract_text_from_docx(file_stream)
    else:
        raise Exception("不支援的文件格式")

    if not extracted_text.strip():
        raise Exception("文件內容為空，無法提取有效文本")

    try:
        conn = connect_to_db(db_config)
        cursor = conn.cursor()

        query = """
        INSERT INTO health_checks (id_number, check_date, data, file_data, extracted_text, created_at)
        VALUES (%s, %s, %s, %s, %s, %s)
        """
        cursor.execute(query, (
            id_number,
            datetime.now().strftime("%Y-%m-%d"),
            extracted_text,
            file_data,
            extracted_text,
            datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        ))

        conn.commit()
        conn.close()
        logger.info(f"健檢資料上傳成功: id_number={id_number}")
        return {"message": "健檢資料上傳成功", "extracted_text": extracted_text}
    except Exception as e:
        logger.error(f"上傳健檢資料時發生錯誤: id_number={id_number}, 錯誤: {str(e)}")
        raise Exception(f"上傳健檢資料時發生錯誤: {str(e)}")

# 從 Azure SQL Database 提取特定個人的健檢資料（包含歷史資料，支持日期範圍篩選）
def extract_health_data(db_config: Dict[str, Any], id_number: str, start_date: str = None, end_date: str = None) -> tuple[Optional[List[str]], Optional[str]]:
    try:
        conn = connect_to_db(db_config)
        cursor = conn.cursor()

        query = """
        SELECT u.full_name, u.gender, u.birth_date, hc.check_date, hc.extracted_text, hc.notes, u.created_at, u.id_number, u.phone_number, u.email, hc.created_at
        FROM health_checks hc
        JOIN users u ON hc.id_number = u.id_number
        WHERE hc.id_number = %s
        """
        params = [id_number]

        if start_date:
            query += " AND hc.check_date >= %s"
            params.append(start_date)
        if end_date:
            query += " AND hc.check_date <= %s"
            params.append(end_date)

        query += " ORDER BY hc.check_date DESC"
        cursor.execute(query, tuple(params))
        records = cursor.fetchall()
        conn.close()

        if not records:
            logger.warning(f"找不到 id_number={id_number} 的健檢資料")
            return None, f"找不到 id_number={id_number} 的健檢資料"

        health_data_list = []
        for record in records:
            health_data = f"使用者姓名: {record[0]}\n"
            health_data += f"性別: {record[1] if record[1] else '未提供'}\n"
            health_data += f"出生日期: {record[2]}\n"
            health_data += f"健檢日期: {record[3]}\n"
            health_data += f"健檢資料: {record[4] if record[4] else '無'}\n"
            health_data += f"建議: {record[5] if record[5] else '無'}\n"
            health_data += f"創建時間: {record[6] if record[6] else '未提供'}\n"
            health_data += f"身分證字號: {record[7]}\n"
            health_data += f"手機號碼: {record[8]}\n"
            health_data += f"電子郵件: {record[9]}\n"
            health_data += f"健檢記錄創建時間: {record[10]}\n"
            health_data += "-" * 50 + "\n"
            health_data_list.append(health_data)

        logger.info(f"成功提取健檢資料: id_number={id_number}, 記錄數={len(health_data_list)}")
        return health_data_list, None
    except Exception as e:
        logger.error(f"提取健檢資料時發生錯誤: id_number={id_number}, 錯誤: {str(e)}")
        return None, f"連線 Azure SQL Database 時發生錯誤: {str(e)}"

# 分析健檢資料
def analyze_health_data(db_config: Dict[str, Any], id_number: str, start_date: str = None, end_date: str = None) -> tuple[Dict[str, Any], Optional[OllamaLLM], Optional[ConversationBufferMemory], Optional[PromptTemplate]]:
    memory = ConversationBufferMemory()

    llm = OllamaLLM(model="llama3:8b", base_url="http://localhost:11434", temperature=0.3)

    analysis_prompt_template = PromptTemplate(
        input_variables=["data"],
        template="你是一個專業的健康分析專家，請嚴格以繁體中文回答，不得使用英文或其他語言，所有醫學名詞必須使用繁體中文（例如使用「收縮壓」代替「systolic blood pressure」，使用「舒張壓」代替「diastolic blood pressure」）。請分析以下健檢資料並提供詳細建議：\n\n{data}\n\n請提供清晰的分析，包括每個指標是否正常（明確說明正常範圍，例如血壓正常範圍為90-120/60-80 mmHg，心率正常範圍為60-100 bpm），並給出至少三項具體的健康建議（例如飲食調整、運動建議、醫療檢查）和至少一項潛在疾病風險（例如高血壓可能導致心血管疾病）。所有回答必須是繁體中文。"
    )

    interactive_prompt_template = PromptTemplate(
        input_variables=["query"],
        template="你是一個專業的健康分析專家，請嚴格以繁體中文回答，不得使用英文或其他語言，所有醫學名詞必須使用繁體中文（例如使用「收縮壓」代替「systolic blood pressure」，使用「舒張壓」代替「diastolic blood pressure」）。基於之前的健檢資料和上下文，回答以下問題並提供新的具體建議：\n\n{query}\n\n請避免重複之前的回應，確保建議清晰實用，並提供至少三項具體行動建議和至少一項潛在疾病風險。所有回答必須是繁體中文。"
    )

    health_data_list, error = extract_health_data(db_config, id_number, start_date, end_date)
    if error:
        raise Exception(error)

    health_data = "\n".join(health_data_list)
    memory.save_context({"input": "健檢資料"}, {"output": health_data})

    try:
        analysis_prompt = analysis_prompt_template.format(data=health_data)
        result = llm.invoke(analysis_prompt)
        logger.info(f"健檢資料分析成功: id_number={id_number}")
        return {
            "health_data": health_data_list,
            "analysis_result": result
        }, llm, memory, interactive_prompt_template
    except Exception as e:
        logger.error(f"分析健檢資料時發生錯誤: id_number={id_number}, 錯誤: {str(e)}")
        raise Exception(f"分析健檢資料時發生錯誤: {str(e)}")