from opencc import OpenCC

# 初始化簡繁轉換器（簡體轉繁體）
cc = OpenCC('s2t')

def convert_to_traditional(text):
    """
    將簡體中文轉換為繁體中文
    Args:
        text (str): 簡體中文文本
    Returns:
        str: 繁體中文文本
    """
    if not isinstance(text, str):
        return str(text)  # 如果輸入不是字串，轉為字串
    return cc.convert(text)