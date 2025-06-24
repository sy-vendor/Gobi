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
  -d '{
    "username": "admin",
    "password": "adminpass",
    "email": "admin@example.com",
    "is_admin": true
  }'

# 2. 用户登录
echo -e "\n\n${GREEN}2. 测试用户登录${NC}"
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass"}' | jq -r '.token')

echo "获取到的 Token: $TOKEN"

echo "2.1 测试 admin 用户登录"
ADMIN_TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"adminpass"}' | jq -r '.token')
echo "获取到的 Admin Token: $ADMIN_TOKEN"

# 3. 创建数据源
echo -e "\n${GREEN}3. 测试创建数据源${NC}"
curl -X POST http://localhost:8080/api/datasources \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "name": "Test DB",
    "type": "mysql",
    "host": "localhost",
    "port": 3306,
    "database": "test",
    "username": "root",
    "password": "password"
  }'

# 4. 创建 SQL 查询
echo -e "\n${GREEN}4. 测试创建 SQL 查询${NC}"
curl -X POST http://localhost:8080/api/queries \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "name": "Test Query",
    "sql": "SELECT * FROM users",
    "description": "Test query description",
    "data_source_id": 1
  }'

# 5. 创建图表
echo -e "\n${GREEN}5. 测试创建图表${NC}"
curl -X POST http://localhost:8080/api/charts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "name": "Test Chart",
    "type": "line",
    "query_id": 1,
    "config": "{\"x_axis\":\"created_at\",\"y_axis\":\"count\"}"
  }'

# 6. 上传 Excel 模板
echo -e "\n${GREEN}6. 测试上传 Excel 模板${NC}"
curl -X POST http://localhost:8080/api/templates \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -F "template=@test_template.xlsx" \
  -F "name=Test Template" \
  -F "description=Test template description"

# 7. 测试数据隔离
echo -e "\n${GREEN}7. 测试数据隔离${NC}"
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user2",
    "password": "password2",
    "email": "user2@example.com"
  }'

# 8. 测试图表类型
echo -e "\n${GREEN}8. 测试不同图表类型${NC}"
for type in "bar" "pie" "line" "scatter" "area" "radar" "gauge" "funnel" "3d-bar" "3d-scatter" "3d-surface" "3d-bubble"; do
  echo "Creating $type chart..."
  curl -X POST http://localhost:8080/api/charts \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -d "{
      \"name\": \"$type Chart\",
      \"type\": \"$type\",
      \"query_id\": 1,
      \"config\": \"{\\\"x_axis\\\":\\\"created_at\\\",\\\"y_axis\\\":\\\"count\\\"}\"
    }"
  echo -e "\n"
done

# 9. 测试手动清理缓存接口
echo "Testing manual cache clearing..."
curl -X POST http://localhost:8080/api/cache/clear \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "type": "all"
  }'

echo -e "\n${GREEN}测试完成！${NC}" 