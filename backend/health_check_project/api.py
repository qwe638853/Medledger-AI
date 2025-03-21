from fastapi import FastAPI, HTTPException, UploadFile, File, Request, Depends, Query
from fastapi.middleware.cors import CORSMiddleware
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm
from fastapi.responses import HTMLResponse
from fastapi.openapi.utils import get_openapi
from fastapi.staticfiles import StaticFiles
import uvicorn
import logging
import os
import jwt
from datetime import datetime, timedelta

# 從 main.py 導入所需的類和函數
from main import (
    CustomOAuth2PasswordRequestForm, RegisterRequest, LoginRequest, ForgotPasswordRequest, InteractiveRequest,
    validate_id_number, register_user, login_user, forgot_password, upload_health_check,
    analyze_health_data, get_user_info, ROLE_USER, ROLE_OTHER, ROLE_HEALTH_CENTER, VALID_ROLES
)

# 定義可用的 scopes
SCOPES = {
    "role:user": "以 user 角色登入",
    "role:health_center": "以 health_center 角色登入",
    "role:other": "以 other 角色登入"
}

# 初始化 FastAPI 應用
app = FastAPI(
    title="Health Check API",
    description="健康檢查系統 API",
    version="1.0.0"
)

# 掛載靜態文件目錄（使用 /api-static 路徑，避免與前端衝突）
app.mount("/api-static", StaticFiles(directory="D:/gg/WOW/backend/health_check_project/static"), name="api-static")

# 設定允許的來源（根據環境變數動態配置）
ALLOWED_ORIGINS = os.getenv("ALLOWED_ORIGINS", "http://localhost:8000").split(",")
if ALLOWED_ORIGINS == ["*"]:
    allow_credentials = False
else:
    allow_credentials = True

# 添加 CORS 中間件以允許前端跨域請求
app.add_middleware(
    CORSMiddleware,
    allow_origins=ALLOWED_ORIGINS,
    allow_credentials=allow_credentials,
    allow_methods=["GET", "POST", "PUT", "DELETE", "OPTIONS"],
    allow_headers=["*"],
)

# 設定日誌
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Azure SQL Database 連線配置
db_config = {
    'server': 'healthdbserver123.database.windows.net',
    'port': 1433,
    'user': 'bojay@healthdbserver123',
    'password': '!Aa1085121',
    'database': 'health_db'
}

# JWT 設定
SECRET_KEY = "your-secret-key"  # 請替換為安全的密鑰
ALGORITHM = "HS256"
ACCESS_TOKEN_EXPIRE_MINUTES = 30

# 定義 OAuth2PasswordBearer 並指定 scopes
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="login", scopes=SCOPES)

# 自定義 OpenAPI 規範（用於簡化表單和設置 scope 下拉選單）
def custom_openapi():
    if app.openapi_schema:
        return app.openapi_schema

    # 使用 get_openapi 函數生成原始的 OpenAPI 規範
    openapi_schema = get_openapi(
        title=app.title,
        version=app.version,
        description=app.description,
        routes=app.routes,
    )

    # 修改 /login 端點的 requestBody
    try:
        for path, path_item in openapi_schema.get("paths", {}).items():
            if path == "/login" and "post" in path_item:
                request_body = path_item.get("post", {}).get("requestBody", {})
                if not request_body:
                    continue
                content = request_body.get("content", {}).get("application/x-www-form-urlencoded", {})
                if not content:
                    continue
                schema = content.get("schema", {})
                if not schema:
                    continue
                # 移除 grant_type、client_id 和 client_secret
                properties = schema.get("properties", {})
                for field in ["grant_type", "client_id", "client_secret"]:
                    if field in properties:
                        del properties[field]
                # 修改 scope 欄位，添加 enum 和顯示名稱
                if "scope" in properties:
                    properties["scope"] = {
                        "type": "string",
                        "enum": ["role:user", "role:health_center", "role:other"],
                        "description": (
                            "選擇角色，可選值：\n"
                            "- 用戶 (role:user)\n"
                            "- 健康中心 (role:health_center)\n"
                            "- 其他 (role:other)"
                        ),
                        "default": "role:user"
                    }
                # 更新 required 字段，移除不需要的字段
                required = schema.get("required", [])
                updated_required = [field for field in required if field not in ["grant_type", "client_id", "client_secret"]]
                if "scope" not in updated_required:
                    updated_required.append("scope")
                schema["required"] = updated_required
    except Exception as e:
        logger.error(f"生成 OpenAPI 規範時發生錯誤: {str(e)}")
        raise

    app.openapi_schema = openapi_schema
    return app.openapi_schema

app.openapi = custom_openapi

# 自定義 Swagger UI HTML，使用本地資源
SWAGGER_UI_HTML = """
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="/api-static/swagger-ui.css" />
    <style>
        html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin: 0; background: #fafafa; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="/api-static/swagger-ui-bundle.js" charset="UTF-8"></script>
    <script src="/api-static/swagger-ui-standalone-preset.js" charset="UTF-8"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: "/openapi.json",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout",
                onComplete: function() {
                    // 強制將 scope 欄位渲染為下拉選單
                    const scopeInputs = document.querySelectorAll('input[name="scope"]');
                    scopeInputs.forEach(input => {
                        const parent = input.parentElement;
                        const select = document.createElement('select');
                        select.name = input.name;
                        select.className = input.className;
                        const options = [
                            { value: "role:user", text: "用戶 (role:user)" },
                            { value: "role:health_center", text: "健康中心 (role:health_center)" },
                            { value: "role:other", text: "其他 (role:other)" }
                        ];
                        options.forEach(opt => {
                            const option = document.createElement('option');
                            option.value = opt.value;
                            option.text = opt.text;
                            if (opt.value === "role:user") {
                                option.selected = true;
                            }
                            select.appendChild(option);
                        });
                        parent.replaceChild(select, input);
                    });
                }
            });
            window.ui = ui;
        };
    </script>
</body>
</html>
"""

@app.get("/docs", include_in_schema=False)
async def custom_swagger_ui_html():
    return HTMLResponse(SWAGGER_UI_HTML)

# 生成 JWT token
def create_access_token(data: dict, expires_delta: timedelta = None):
    to_encode = data.copy()
    if expires_delta:
        expire = datetime.utcnow() + expires_delta
    else:
        expire = datetime.utcnow() + timedelta(minutes=15)
    to_encode.update({"exp": expire})
    encoded_jwt = jwt.encode(to_encode, SECRET_KEY, algorithm=ALGORITHM)
    return encoded_jwt

# 驗證 JWT token
def decode_access_token(token: str):
    try:
        payload = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
        id_number: str = payload.get("id_number")
        role: str = payload.get("role")
        if id_number is None or role is None:
            raise HTTPException(status_code=401, detail="無效的 token")
        return {"id_number": id_number, "role": role}
    except jwt.ExpiredSignatureError:
        raise HTTPException(status_code=401, detail="token 已過期")
    except jwt.InvalidTokenError:
        raise HTTPException(status_code=401, detail="無效的 token")

# 獲取當前用戶
async def get_current_user(token: str = Depends(oauth2_scheme)):
    logger.info(f"開始驗證 token: {token}")
    try:
        user_data = decode_access_token(token)
        id_number = user_data["id_number"]
        role = user_data["role"]
        
        logger.info(f"查詢用戶: id_number={id_number}")
        user = get_user_info(db_config, id_number)
        if not user:
            logger.error(f"用戶 {id_number} 不存在")
            raise HTTPException(status_code=401, detail=f"用戶 {id_number} 不存在")
        
        if user["role"] != role:
            logger.error(f"角色不匹配：token 中的角色為 {role}，但用戶角色為 {user['role']}")
            raise HTTPException(status_code=401, detail=f"角色不匹配：token 中的角色為 {role}，但用戶角色為 {user['role']}")
        
        logger.info(f"驗證成功: id_number={id_number}, role={role}")
        return {"id_number": id_number, "role": role}
    except HTTPException as e:
        raise e
    except Exception as e:
        logger.error(f"驗證失敗: {str(e)}")
        raise HTTPException(status_code=401, detail=f"驗證失敗: {str(e)}")

# API 端點：註冊新用戶
@app.post("/register")
async def register(request: RegisterRequest):
    logger.info(f"收到註冊請求: id_number={request.id_number}, role={request.role}")
    try:
        result = register_user(db_config, request)
        logger.info(f"註冊成功: id_number={request.id_number}, role={request.role}")
        return result
    except Exception as e:
        logger.error(f"註冊失敗: {str(e)}")
        raise HTTPException(status_code=500, detail=f"註冊時發生錯誤: {str(e)}")

# API 端點：用戶登入（支援 OAuth2 標準）
@app.post("/login")
async def login(form_data: OAuth2PasswordRequestForm = Depends()):
    logger.info(f"收到登入請求: username={form_data.username}, scopes={form_data.scopes}")
    
    # 從 scopes 中提取 role
    role = None
    for scope in form_data.scopes:
        if scope.startswith("role:"):
            role = scope.split(":")[1]
            break
    
    if not role:
        logger.error("未提供角色，請在 scopes 中指定 role，例如 role:user")
        raise HTTPException(status_code=400, detail="未提供角色，請在 scopes 中指定 role，例如 role:user")

    try:
        request = LoginRequest(
            id_number=form_data.username,
            password=form_data.password,
            role=role
        )
        result = login_user(db_config, request)
        # 生成 JWT token
        access_token_expires = timedelta(minutes=ACCESS_TOKEN_EXPIRE_MINUTES)
        access_token = create_access_token(
            data={"id_number": result["id_number"], "role": result["role"]},
            expires_delta=access_token_expires
        )
        logger.info(f"登入成功: id_number={form_data.username}, role={role}")
        return {
            "access_token": access_token,
            "token_type": "bearer"
        }
    except Exception as e:
        logger.error(f"登入失敗: {str(e)}")
        raise HTTPException(status_code=500, detail=f"登入時發生錯誤: {str(e)}")

# API 端點：忘記密碼
@app.post("/forgot-password")
async def forgot_password_endpoint(request: ForgotPasswordRequest):
    logger.info(f"收到忘記密碼請求: id_number={request.id_number}")
    try:
        result = forgot_password(db_config, request)
        logger.info(f"密碼重置成功: id_number={request.id_number}")
        return result
    except Exception as e:
        logger.error(f"密碼重置失敗: {str(e)}")
        raise HTTPException(status_code=500, detail=f"密碼重置時發生錯誤: {str(e)}")

# API 端點：上傳健檢資料文件（僅限健檢中心）
@app.post("/health-check/upload/{id_number}")
async def upload_health_check_endpoint(
    id_number: str,
    file: UploadFile = File(...),
    current_user: dict = Depends(get_current_user)
):
    logger.info(f"收到上傳健檢資料請求: id_number='{id_number}', 當前用戶={current_user['id_number']}, 角色={current_user['role']}")
    
    if not id_number:
        logger.error("id_number 為空")
        raise HTTPException(status_code=400, detail="請提供身分證字號")

    if current_user["role"] != ROLE_HEALTH_CENTER:
        logger.error(f"權限不足: 當前用戶角色為 {current_user['role']}，僅健檢中心可上傳")
        raise HTTPException(status_code=403, detail="權限不足，僅健康中心可以上傳健檢資料")

    is_valid, error_message = validate_id_number(id_number)
    if not is_valid:
        logger.error(f"身分證字號格式不正確: id_number='{id_number}', 錯誤: {error_message}")
        raise HTTPException(status_code=400, detail=f"身分證字號格式不正確: {error_message}")
    
    user = get_user_info(db_config, id_number)
    if user is None:
        logger.error(f"身分證字號不存在: id_number='{id_number}'")
        raise HTTPException(status_code=404, detail="身分證字號不存在")
    
    try:
        result = await upload_health_check(db_config, id_number, file)
        logger.info(f"上傳健檢資料成功: id_number='{id_number}'")
        return result
    except Exception as e:
        logger.error(f"上傳健檢資料失敗: {str(e)}")
        raise HTTPException(status_code=500, detail=f"上傳健檢資料時發生錯誤: {str(e)}")

# API 端點：user 角色，獲取健檢資料和分析結果
@app.get("/health-check/user/{id_number}")
async def get_user_health_check(
    id_number: str,
    start_date: str = Query(None, description="開始日期 (格式: YYYY-MM-DD)"),
    end_date: str = Query(None, description="結束日期 (格式: YYYY-MM-DD)"),
    current_user: dict = Depends(get_current_user)
):
    logger.info(f"收到獲取健檢資料請求 (user 角色): id_number={id_number}, 當前用戶={current_user['id_number']}, 角色={current_user['role']}, start_date={start_date}, end_date={end_date}")
    if current_user["role"] != ROLE_USER:
        logger.error(f"權限不足: 當前用戶角色為 {current_user['role']}，此端點僅限 user 角色")
        raise HTTPException(status_code=403, detail="權限不足，此端點僅限 user 角色")
    
    is_valid, error_message = validate_id_number(id_number)
    if not is_valid:
        logger.error(f"身分證字號格式不正確: id_number={id_number}, 錯誤: {error_message}")
        raise HTTPException(status_code=400, detail=f"身分證字號格式不正確: {error_message}")
    
    user = get_user_info(db_config, id_number)
    if user is None:
        logger.error(f"身分證字號不存在: id_number={id_number}")
        raise HTTPException(status_code=404, detail="身分證字號不存在")
    
    if current_user["id_number"] != id_number:
        logger.error(f"權限不足: 用戶 {current_user['id_number']} 無權查看 {id_number} 的資料")
        raise HTTPException(status_code=403, detail="權限不足，只能查看自己的健檢資料")
    
    try:
        result, _, _, _ = analyze_health_data(db_config, id_number, start_date, end_date)
        logger.info(f"獲取健檢資料成功 (user 角色): id_number={id_number}")
        return result
    except Exception as e:
        logger.error(f"獲取健檢資料失敗: {str(e)}")
        raise HTTPException(status_code=500, detail=f"獲取健檢資料時發生錯誤: {str(e)}")

# API 端點：other 和 health_center 角色，獲取健檢資料和分析結果，並支援互動模式
@app.get("/health-check/other/{id_number}")
async def get_other_health_check(
    id_number: str,
    request: Request,
    start_date: str = Query(None, description="開始日期 (格式: YYYY-MM-DD)"),
    end_date: str = Query(None, description="結束日期 (格式: YYYY-MM-DD)"),
    current_user: dict = Depends(get_current_user)
):
    logger.info(f"收到獲取健檢資料請求 (other/health_center 角色): id_number={id_number}, 當前用戶={current_user['id_number']}, 角色={current_user['role']}, start_date={start_date}, end_date={end_date}")
    if current_user["role"] not in [ROLE_OTHER, ROLE_HEALTH_CENTER]:
        logger.error(f"權限不足: 當前用戶角色為 {current_user['role']}，此端點僅限 other 和 health_center 角色")
        raise HTTPException(status_code=403, detail="權限不足，此端點僅限 other 和 health_center 角色")
    
    is_valid, error_message = validate_id_number(id_number)
    if not is_valid:
        logger.error(f"身分證字號格式不正確: id_number={id_number}, 錯誤: {error_message}")
        raise HTTPException(status_code=400, detail=f"身分證字號格式不正確: {error_message}")
    
    user = get_user_info(db_config, id_number)
    if user is None:
        logger.error(f"身分證字號不存在: id_number={id_number}")
        raise HTTPException(status_code=404, detail="身分證字號不存在")
    
    try:
        result, llm, memory, interactive_prompt_template = analyze_health_data(db_config, id_number, start_date, end_date)
        request.state.llm = llm
        request.state.memory = memory
        request.state.interactive_prompt_template = interactive_prompt_template
        logger.info(f"獲取健檢資料成功 (other/health_center 角色): id_number={id_number}")
        return result
    except Exception as e:
        logger.error(f"獲取健檢資料失敗: {str(e)}")
        raise HTTPException(status_code=500, detail=f"獲取健檢資料時發生錯誤: {str(e)}")

# API 端點：other 和 health_center 角色的互動模式
@app.post("/health-check/other/interact")
async def interact_health_check(request: Request, data: InteractiveRequest, current_user: dict = Depends(get_current_user)):
    logger.info(f"收到互動模式請求: 當前用戶={current_user['id_number']}, 角色={current_user['role']}")
    if current_user["role"] not in [ROLE_OTHER, ROLE_HEALTH_CENTER]:
        logger.error(f"權限不足: 當前用戶角色為 {current_user['role']}，此端點僅限 other 和 health_center 角色")
        raise HTTPException(status_code=403, detail="權限不足，此端點僅限 other 和 health_center 角色")
    
    if not hasattr(request.state, 'llm') or not hasattr(request.state, 'memory'):
        logger.error("互動模式未初始化")
        raise HTTPException(status_code=400, detail="請先呼叫 /health-check/other/{id_number} 來初始化互動模式")
    
    try:
        interactive_prompt = request.state.interactive_prompt_template.format(query=data.query)
        response = request.state.llm.invoke(interactive_prompt)
        logger.info("互動模式請求成功")
        return {"response": response}
    except Exception as e:
        logger.error(f"處理互動問題時發生錯誤: {str(e)}")
        raise HTTPException(status_code=500, detail=f"處理互動問題時發生錯誤: {str(e)}")

# 啟動伺服器
if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)