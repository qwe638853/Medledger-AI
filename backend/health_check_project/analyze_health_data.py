import pymssql
from langchain_community.llms import Ollama
from langchain.prompts import PromptTemplate
from langchain.chains import LLMChain
from langchain.memory import ConversationBufferMemory

# 步驟 1：從 Azure SQL Database 提取特定個人的健檢資料
def extract_health_data(db_config, user_id):
    try:
        conn = pymssql.connect(
            server=db_config['server'],
            port=db_config['port'],
            user=db_config['user'],
            password=db_config['password'],
            database=db_config['database'],
            login_timeout=30,
            charset='UTF-8'  # 明確指定 UTF-8 編碼
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

        print(f"\n從 Azure SQL Database 提取的健檢資料（user_id={user_id}）：\n{health_data}")
        return health_data, None
    except Exception as e:
        return None, f"連線 Azure SQL Database 時發生錯誤：{e}"

# 步驟 2：分析健檢資料並支援互動模式
def analyze_health_data(db_config, user_id):
    all_results = []
    memory = ConversationBufferMemory()

    # 初始化 Ollama 模型
    llm = Ollama(model="llama3:8b", base_url="http://localhost:11434", temperature=0.3)
    print("Ollama 模型初始化成功，模型名稱：", llm.model)

    # 定義分析的 Prompt 模板
    analysis_prompt_template = PromptTemplate(
        input_variables=["data"],
        template="你是一個專業的健康分析專家，請嚴格以繁體中文回答，不得使用英文或其他語言。請分析以下健檢資料並提供詳細建議：\n\n{data}\n\n請提供清晰的分析，包括每個指標是否正常，並給出具體的健康建議。所有回答必須是繁體中文。"
    )

    # 定義互動模式的 Prompt 模板
    interactive_prompt_template = PromptTemplate(
        input_variables=["query"],
        template="你是一個專業的健康分析專家，請嚴格以繁體中文回答，不得使用英文或其他語言。基於之前的健檢資料和上下文，回答以下問題並提供新的具體建議：\n\n{query}\n\n請避免重複之前的回應，確保建議清晰實用，並提供至少一項具體行動建議。所有回答必須是繁體中文。"
    )

    # 創建 LLMChain（用於分析）
    analysis_chain = LLMChain(llm=llm, prompt=analysis_prompt_template, memory=memory)
    print("LangChain 分析鏈接創建成功")

    # 創建 LLMChain（用於互動模式）
    interactive_chain = LLMChain(llm=llm, prompt=interactive_prompt_template, memory=memory)
    print("LangChain 互動鏈接創建成功")

    # 提取特定個人的健檢資料
    health_data, error = extract_health_data(db_config, user_id)
    if error:
        print(error)
        return [], interactive_chain, memory

    # 使用 Chain 進行分析
    try:
        result = analysis_chain.run(data=health_data)
        all_results.append(result)
        print(f"\n分析結果：\n{result}")
    except Exception as e:
        print(f"分析健檢資料時發生錯誤：{e}")

    return all_results, interactive_chain, memory

# 主程式：支援動態查詢與互動模式
if __name__ == "__main__":
    # Azure SQL Database 連線配置（與 SQLTools 一致）
    db_config = {
        'server': 'healthdbserver123.database.windows.net',  # 替換為你的 Azure SQL 伺服器地址
        'port': 1433,
        'user': 'adminuser',  # 替換為你的管理員名稱
        'password': 'your_secure_password',  # 替換為你的密碼
        'database': 'health_db'
    }

    try:
        # 動態輸入 user_id
        user_id = int(input("請輸入 user_id（例如 1）："))
        results, interactive_chain, memory = analyze_health_data(db_config, user_id=user_id)
        
        print("\n所有分析結果總結：")
        for i, result in enumerate(results, 1):
            print(f"結果 {i}:")
            print(result)

        # 進入互動模式
        print("\n=== 進入互動模式 ===")
        print("你可以繼續提問，例如：'我需要更多健康建議' 或 '我的血壓有問題怎麼辦？'（輸入 'exit' 退出，或 'clear' 清除記憶）")
        while True:
            user_input = input("請輸入你的問題：")
            if user_input.lower() == "exit":
                print("退出互動模式")
                break
            elif user_input.lower() == "clear":
                memory.clear()
                print("記憶已清除")
                continue
            try:
                # 使用 LLMChain 處理問題
                response = interactive_chain.run(query=user_input)
                print("\n回應：")
                print(response)
            except Exception as e:
                print(f"處理問題時發生錯誤：{e}")

    except Exception as e:
        print(f"主程式發生錯誤：{e}")