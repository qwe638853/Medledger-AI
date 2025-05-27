import requests
import logging
import os
import tkinter as tk
from tkinter import filedialog, messagebox, simpledialog

# 設定日誌
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# 設置請求的頭部
headers = {
    "Authorization": "Bearer C123456789:health_center",  # 使用健康中心的 token
    "Accept": "application/json"
}

# 創建一個簡單的 Tkinter 窗口來輸入 id_number 和選擇檔案
def get_user_input():
    root = tk.Tk()
    root.withdraw()  # 隱藏主窗口

    # 讓使用者輸入 id_number
    logger.info("請輸入目標用戶的身分證字號")
    id_number = simpledialog.askstring(
        title="輸入身分證字號",
        prompt="請輸入目標用戶的身分證字號（例如 A456789123）："
    )
    if not id_number:
        root.destroy()
        logger.error("未輸入身分證字號，程式退出")
        raise ValueError("未輸入身分證字號，程式退出")

    # 驗證身分證字號格式（簡單檢查）
    if len(id_number) != 10 or not id_number[0].isalpha() or not id_number[1:].isdigit():
        root.destroy()
        logger.error(f"身分證字號格式不正確: {id_number}")
        raise ValueError("身分證字號格式不正確，應為 1 個字母 + 9 個數字，例如 A123456789")

    # 設置 URL
    url = f"http://localhost:8000/health-check/upload/{id_number}"

    # 讓使用者選擇檔案
    logger.info("請選擇要上傳的檔案")
    file_path = filedialog.askopenfilename(
        title="選擇健康檢查報告檔案",
        filetypes=[("PDF files", "*.pdf"), ("All files", "*.*")]
    )
    if not file_path:
        root.destroy()
        logger.error("未選擇任何檔案，程式退出")
        raise ValueError("未選擇任何檔案，程式退出")

    root.destroy()
    return url, file_path

# 獲取使用者輸入
try:
    url, file_path = get_user_input()
except ValueError as e:
    logger.error(str(e))
    exit(1)

# 檢查檔案是否存在
if not os.path.exists(file_path):
    logger.error(f"檔案不存在: {file_path}")
    raise FileNotFoundError(f"檔案不存在: {file_path}")

# 獲取檔案名稱
file_name = os.path.basename(file_path)

# 設置文件
files = {
    "file": (file_name, open(file_path, "rb"), "application/pdf")
}

# 發送 POST 請求
try:
    logger.info(f"開始發送上傳請求到: {url}")
    logger.info(f"上傳檔案: {file_path}")
    response = requests.post(url, headers=headers, files=files)
    logger.info(f"狀態碼: {response.status_code}")
    logger.info(f"響應內容: {response.json()}")
    # 顯示成功訊息
    root = tk.Tk()
    root.withdraw()
    messagebox.showinfo("成功", "健康檢查資料上傳成功！")
    root.destroy()
except requests.exceptions.RequestException as e:
    logger.error(f"請求失敗: {str(e)}")
    if e.response is not None:
        logger.error(f"錯誤詳情: {e.response.text}")
    # 顯示錯誤訊息
    root = tk.Tk()
    root.withdraw()
    messagebox.showerror("錯誤", f"上傳失敗: {str(e)}")
    root.destroy()
except Exception as e:
    logger.error(f"發生未知錯誤: {str(e)}")
    # 顯示錯誤訊息
    root = tk.Tk()
    root.withdraw()
    messagebox.showerror("錯誤", f"發生未知錯誤: {str(e)}")
    root.destroy()
finally:
    # 關閉文件
    if "file" in files:
        files["file"][1].close()
        logger.info("文件已關閉")