#!/bin/bash

# 上传Excel模板到API
# 使用方法: ./upload_template.sh [template_file]

# 默认模板文件
TEMPLATE_FILE=${1:-"daily_report_template_20250619_161225.xlsx"}

# API配置
API_BASE="http://localhost:8080/api"
LOGIN_URL="$API_BASE/auth/login"
TEMPLATE_URL="$API_BASE/templates"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}开始上传Excel模板...${NC}"

# 检查模板文件是否存在
if [ ! -f "$TEMPLATE_FILE" ]; then
    echo -e "${RED}错误: 模板文件 $TEMPLATE_FILE 不存在${NC}"
    echo "请先运行: go run scripts/create_template.go"
    exit 1
fi

echo -e "${GREEN}找到模板文件: $TEMPLATE_FILE${NC}"

# 登录获取token
echo -e "${YELLOW}正在登录...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$LOGIN_URL" \
    -H "Content-Type: application/json" \
    -d '{
        "username": "admin",
        "password": "admin123"
    }')

# 提取token
TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}登录失败，请检查用户名和密码${NC}"
    echo "响应: $LOGIN_RESPONSE"
    exit 1
fi

echo -e "${GREEN}登录成功，获取到token${NC}"

# 上传模板
echo -e "${YELLOW}正在上传模板...${NC}"
UPLOAD_RESPONSE=$(curl -s -X POST "$TEMPLATE_URL" \
    -H "Authorization: Bearer $TOKEN" \
    -F "template=@$TEMPLATE_FILE" \
    -F "description=标准日报模板，包含销售数据、用户数据和系统性能指标")

# 检查上传结果
if echo "$UPLOAD_RESPONSE" | grep -q '"id"'; then
    TEMPLATE_ID=$(echo "$UPLOAD_RESPONSE" | grep -o '"id":[0-9]*' | cut -d':' -f2)
    echo -e "${GREEN}模板上传成功！${NC}"
    echo -e "${GREEN}模板ID: $TEMPLATE_ID${NC}"
    echo -e "${YELLOW}你可以在创建报告计划时使用这个template_ids: [$TEMPLATE_ID]${NC}"
    
    # 显示模板详情
    echo -e "\n${YELLOW}模板详情:${NC}"
    echo "$UPLOAD_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$UPLOAD_RESPONSE"
    
else
    echo -e "${RED}模板上传失败${NC}"
    echo "响应: $UPLOAD_RESPONSE"
    exit 1
fi

echo -e "\n${GREEN}上传完成！${NC}"
echo -e "${YELLOW}现在你可以使用template_ids: [$TEMPLATE_ID] 来创建报告计划${NC}" 