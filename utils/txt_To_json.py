import json

# 读取文件内容
with open('floor.txt', 'r', encoding='utf-8') as file:
    content = file.read()

# 解析JSON内容
data = json.loads(content)

# 提取message字段的内容
message_data = data['message']

# 将message字段的内容转换为标准JSON格式
standard_json = json.dumps(message_data, ensure_ascii=False, indent=4)

# 输出标准JSON
print(standard_json)

# 如果需要将标准JSON保存到文件
with open('标准分类.json', 'w', encoding='utf-8') as output_file:
    output_file.write(standard_json)