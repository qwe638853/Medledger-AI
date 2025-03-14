import pymssql
from langchain_ollama import OllamaLLM
from langchain.prompts import PromptTemplate
from langchain.memory import ConversationBufferMemory
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import Optional, Dict, Any

# 初始化 FastAPI 應用
app = FastAPI()

# 定義請求體結構（用於互動模式）
class InteractiveRequest(BaseModel):
    query: str

# 從 Azure SQL Database 提取特定個人的健檢資料
def extract_health_data(db_config: Dict[str, Any], user_id: int) -> tuple[Optional[str], Optional[str]]:
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
        cursor = conn.cursor()

        # 查詢 health_checks 表
        query = """
        SELECT u.full_name, u.gender, u.birth_date, hc.check_date, hc.data, hc.notes
        FROM health_checks hc
        JOIN users u ON hc.user_id = u.id
        WHERE hc.user_id = %s
        """
        cursor.execute(query, (user_id,))
        record = cursor.fetchone()
        conn.close()

        if not record:
            return None, f"找不到 user_id={user_id} 的健檢資料"

        # 格式化資料
        health_data = f"使用者姓名: {record[0]}\n"
        health_data += f"性別: {record[1]}\n"
        health_data += f"出生日期: {record[2]}\n"
        health_data += f"健檢日期: {record[3]}\n"
        health_data += f"健檢資料: {record[4]}\n"
        health_data += f"建議: {record[5] if record[5] else '無'}\n"

        return health_data, None
    except Exception as e:
        return None, f"連線 Azure SQL Database 時發生錯誤：{e}"

# 分析健檢資料
def analyze_health_data(db_config: Dict[str, Any], user_id: int) -> tuple[Optional[str], Optional[OllamaLLM], Optional[ConversationBufferMemory], Optional[PromptTemplate]]:
    memory = ConversationBufferMemory()

    # 初始化 Ollama 模型
    llm = OllamaLLM(model="llama3:8b", base_url="http://localhost:11434", temperature=0.3)

    # 定義分析的 Prompt 模板
    analysis_prompt_template = PromptTemplate(
        input_variables=["data"],
        template="你是一個專業的健康分析專家，請嚴格以繁體中文回答，不得使用英文或其他語言，所有醫學名詞必須使用繁體中文（例如使用「收縮壓」代替「systolic blood pressure」，使用「舒張壓」代替「diastolic blood pressure」）。請分析以下健檢資料並提供詳細建議：\n\n{data}\n\n請提供清晰的分析，包括每個指標是否正常（明確說明正常範圍，例如血壓正常範圍為90-120/60-80 mmHg，心率正常範圍為60-100 bpm），並給出至少三項具體的健康建議（例如飲食調整、運動建議、醫療檢查）和至少一項潛在疾病風險（例如高血壓可能導致心血管疾病）。所有回答必須是繁體中文。"
    )

    # 定義互動模式的 Prompt 模板（僅 other 角色使用）
    interactive_prompt_template = PromptTemplate(
        input_variables=["query"],
        template="你是一個專業的健康分析專家，請嚴格以繁體中文回答，不得使用英文或其他語言，所有醫學名詞必須使用繁體中文（例如使用「收縮壓」代替「systolic blood pressure」，使用「舒張壓」代替「diastolic blood pressure」）。基於之前的健檢資料和上下文，回答以下問題並提供新的具體建議：\n\n{query}\n\n請避免重複之前的回應，確保建議清晰實用，並提供至少三項具體行動建議和至少一項潛在疾病風險。所有回答必須是繁體中文。"
    )

    # 提取特定個人的健檢資料
    health_data, error = extract_health_data(db_config, user_id)
    if error:
        raise HTTPException(status_code=500, detail=error)

    # 儲存健檢資料到記憶中（供互動模式使用）
    memory.save_context({"input": "健檢資料"}, {"output": health_data})

    # 使用 llm.invoke 進行分析
    try:
        analysis_prompt = analysis_prompt_template.format(data=health_data)
        result = llm.invoke(analysis_prompt)
        return {
            "health_data": health_data,
            "analysis_result": result
        }, llm, memory, interactive_prompt_template
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"分析健檢資料時發生錯誤：{e}")

# Azure SQL Database 連線配置
db_config = {
    'server': 'healthdbserver123.database.windows.net',
    'port': 1433,
    'user': 'bojay',
    'password': '!Aa1085121',
    'database': 'health_db'
}

# API 端點：user 角色，獲取健檢資料和分析結果
@app.get("/health-check/user/{user_id}")
async def get_user_health_check(user_id: int):
    result, _, _, _ = analyze_health_data(db_config, user_id)
    return result

# API 端點：other 角色，獲取健檢資料和分析結果，並支援互動模式
# 這裡我們先返回分析結果，互動模式可以通過另一個端點實現
@app.get("/health-check/other/{user_id}")
async def get_other_health_check(user_id: int):
    result, llm, memory, interactive_prompt_template = analyze_health_data(db_config, user_id)
    # 儲存 llm 和 memory 到全局變數（簡單實現，實際應用中應使用更穩健的方式）
    global global_llm, global_memory, global_interactive_prompt_template
    global_llm = llm
    global_memory = memory
    global_interactive_prompt_template = interactive_prompt_template
    return result

# API 端點：other 角色的互動模式
@app.post("/health-check/other/interact")
async def interact_health_check(request: InteractiveRequest):
    if not hasattr(globals(), 'global_llm') or not hasattr(globals(), 'global_memory'):
        raise HTTPException(status_code=400, detail="請先呼叫 /health-check/other/{user_id} 來初始化互動模式")
    
    try:
        interactive_prompt = global_interactive_prompt_template.format(query=request.query)
        response = global_llm.invoke(interactive_prompt)
        return {"response": response}
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"處理互動問題時發生錯誤：{e}")

# 啟動伺服器
if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)