import pymssql
from passlib.context import CryptContext

# 資料庫配置
db_config = {
    'server': 'healthdbserver123.database.windows.net',
    'port': 1433,
    'user': 'bojay',
    'password': '!Aa1085121',
    'database': 'health_db'
}

# 初始化密碼加密工具
pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")

# 連接到資料庫
def connect_to_db():
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
        return conn
    except Exception as e:
        print(f"連線資料庫時發生錯誤：{e}")
        return None

# 修復密碼
def fix_passwords():
    passwords_to_fix = {
        'A123450003': 'fuckme',
        'A123450001': 'default_password',
        'A123450002': 'default_password'
    }

    conn = connect_to_db()
    if not conn:
        return

    cursor = conn.cursor()

    try:
        for id_number, plain_password in passwords_to_fix.items():
            hashed_password = pwd_context.hash(plain_password)
            print(f"Generated hash for {id_number} (password: {plain_password}): {hashed_password}")

            if id_number == 'A123450003':
                query = """
                UPDATE users
                SET password = %s
                WHERE id_number = %s
                """
                cursor.execute(query, (hashed_password, id_number))
            else:
                query = """
                UPDATE users
                SET password = %s
                WHERE id_number = %s
                """
                cursor.execute(query, (hashed_password, id_number))

            print(f"Updated password for id_number {id_number}")

        conn.commit()

    except Exception as e:
        print(f"更新密碼時發生錯誤：{e}")
        conn.rollback()

    finally:
        cursor.close()
        conn.close()

if __name__ == "__main__":
    print("開始修復密碼...")
    fix_passwords()
    print("修復完成！")