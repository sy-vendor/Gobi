#!/bin/bash

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

echo "开始测试 Gobi BI 引擎..."

# 1. 注册用户
echo -e "\n${GREEN}1. 测试用户注册${NC}"
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass"}'

# 2. 用户登录
echo -e "\n\n${GREEN}2. 测试用户登录${NC}"
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass"}' | jq -r '.token')

echo "获取到的 Token: $TOKEN"

# 3. 创建 SQL 查询
echo -e "\n${GREEN}3. 测试创建 SQL 查询${NC}"
QUERY_RESPONSE=$(curl -s -X POST http://localhost:8080/api/queries \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "测试查询",
    "sql": "SELECT * FROM users",
    "description": "测试查询描述",
    "isPublic": true
  }')

echo "查询响应: $QUERY_RESPONSE"
QUERY_ID=$(echo $QUERY_RESPONSE | jq -r '.ID')

# 4. 创建图表
echo -e "\n${GREEN}4. 测试创建图表${NC}"
CHART_RESPONSE=$(curl -s -X POST http://localhost:8080/api/charts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "测试图表",
    "type": "bar",
    "queryId": '$QUERY_ID',
    "config": "{\"title\":\"测试图表\",\"xAxis\":{\"type\":\"category\"},\"yAxis\":{\"type\":\"value\"}}",
    "data": "{\"categories\":[\"A\",\"B\",\"C\"],\"values\":[10,20,30]}"
  }')

echo "图表响应: $CHART_RESPONSE"
CHART_ID=$(echo $CHART_RESPONSE | jq -r '.ID')

# 5. 上传 Excel 模板
echo -e "\n${GREEN}5. 测试上传 Excel 模板${NC}"
# 创建一个简单的 Excel 文件
echo "创建测试 Excel 文件..."
cat > test_template.xlsx << EOL
PK
EOL

TEMPLATE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/templates \
  -H "Authorization: Bearer $TOKEN" \
  -F "template=@test_template.xlsx")

echo "模板响应: $TEMPLATE_RESPONSE"
TEMPLATE_ID=$(echo $TEMPLATE_RESPONSE | jq -r '.ID')

# 6. 测试数据隔离
echo -e "\n${GREEN}6. 测试数据隔离${NC}"
# 创建另一个用户
curl -s -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser2","password":"testpass2"}'

# 新用户登录
TOKEN2=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser2","password":"testpass2"}' | jq -r '.token')

# 尝试访问第一个用户的查询
echo "尝试访问其他用户的查询..."
curl -s -X GET http://localhost:8080/api/queries/$QUERY_ID \
  -H "Authorization: Bearer $TOKEN2"

# 7. 测试图表类型
echo -e "\n${GREEN}7. 测试不同图表类型${NC}"
# 创建柱状图
curl -s -X POST http://localhost:8080/api/charts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "柱状图测试",
    "type": "bar",
    "queryId": '$QUERY_ID',
    "config": "{\"title\":\"柱状图测试\",\"xAxis\":{\"type\":\"category\"},\"yAxis\":{\"type\":\"value\"}}",
    "data": "{\"categories\":[\"A\",\"B\",\"C\"],\"values\":[10,20,30]}"
  }'

# 创建折线图
curl -s -X POST http://localhost:8080/api/charts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "折线图测试",
    "type": "line",
    "queryId": '$QUERY_ID',
    "config": "{\"title\":\"折线图测试\",\"xAxis\":{\"type\":\"category\"},\"yAxis\":{\"type\":\"value\"}}",
    "data": "{\"categories\":[\"A\",\"B\",\"C\"],\"values\":[10,20,30]}"
  }'

# 创建饼图
curl -s -X POST http://localhost:8080/api/charts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "饼图测试",
    "type": "pie",
    "queryId": '$QUERY_ID',
    "config": "{\"title\":\"饼图测试\"}",
    "data": "{\"items\":[{\"name\":\"A\",\"value\":30},{\"name\":\"B\",\"value\":40},{\"name\":\"C\",\"value\":30}]}"
  }'

# 创建散点图
curl -s -X POST http://localhost:8080/api/charts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "散点图测试",
    "type": "scatter",
    "queryId": '$QUERY_ID',
    "config": "{\"title\":\"散点图测试\",\"xAxis\":{\"type\":\"value\"},\"yAxis\":{\"type\":\"value\"}}",
    "data": "{\"points\":[{\"x\":10,\"y\":20},{\"x\":15,\"y\":25},{\"x\":20,\"y\":30}]}"
  }'

# 创建雷达图
curl -s -X POST http://localhost:8080/api/charts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "雷达图测试",
    "type": "radar",
    "queryId": '$QUERY_ID',
    "config": "{\"title\":\"雷达图测试\",\"radar\":{\"indicator\":[{\"name\":\"指标1\",\"max\":100},{\"name\":\"指标2\",\"max\":100},{\"name\":\"指标3\",\"max\":100}]}}",
    "data": "{\"values\":[80,70,90]}"
  }'

# 创建热力图
curl -s -X POST http://localhost:8080/api/charts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "热力图测试",
    "type": "heatmap",
    "queryId": '$QUERY_ID',
    "config": "{\"title\":\"热力图测试\",\"xAxis\":{\"type\":\"category\",\"data\":[\"A\",\"B\",\"C\"]},\"yAxis\":{\"type\":\"category\",\"data\":[\"X\",\"Y\",\"Z\"]}}",
    "data": "{\"data\":[[10,20,30],[40,50,60],[70,80,90]]}"
  }'

# 创建仪表盘
curl -s -X POST http://localhost:8080/api/charts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "仪表盘测试",
    "type": "gauge",
    "queryId": '$QUERY_ID',
    "config": "{\"title\":\"仪表盘测试\",\"min\":0,\"max\":100}",
    "data": "{\"value\":75}"
  }'

# 创建漏斗图
curl -s -X POST http://localhost:8080/api/charts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "漏斗图测试",
    "type": "funnel",
    "queryId": '$QUERY_ID',
    "config": "{\"title\":\"漏斗图测试\"}",
    "data": "{\"items\":[{\"name\":\"访问\",\"value\":100},{\"name\":\"注册\",\"value\":80},{\"name\":\"下单\",\"value\":60},{\"name\":\"支付\",\"value\":40},{\"name\":\"完成\",\"value\":20}]}"
  }'

# 清理测试文件
rm test_template.xlsx

echo -e "\n${GREEN}测试完成！${NC}" 