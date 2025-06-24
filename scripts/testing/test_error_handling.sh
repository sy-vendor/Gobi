#!/bin/bash
# 自动化测试主要接口和错误处理
set -e
API_URL="http://localhost:8080/api"
echo "==== 注册用户 ===="
curl -s -X POST $API_URL/auth/register -H "Content-Type: application/json" -d '{"username":"testuser","email":"testuser@example.com","password":"test123"}'
echo -e "\n==== 重复注册 ===="
curl -s -X POST $API_URL/auth/register -H "Content-Type: application/json" -d '{"username":"testuser","email":"testuser@example.com","password":"test123"}'
echo -e "\n==== 注册缺少参数 ===="
curl -s -X POST $API_URL/auth/register -H "Content-Type: application/json" -d '{"username":""}'
echo -e "\n==== 登录 ===="
LOGIN_RESP=$(curl -s -X POST $API_URL/auth/login -H "Content-Type: application/json" -d '{"username":"testuser","password":"test123"}')
TOKEN=$(echo $LOGIN_RESP | grep -o '"token":"[^"]*"' | cut -d '"' -f4)
echo "Token: $TOKEN"
echo -e "\n==== 登录密码错误 ===="
curl -s -X POST $API_URL/auth/login -H "Content-Type: application/json" -d '{"username":"testuser","password":"wrong"}'
echo -e "\n==== 创建数据源 ===="
curl -s -X POST $API_URL/datasources -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d '{"name":"Test DS","type":"sqlite","database":"/tmp/test.db","username":"","password":"","host":"","port":0,"description":"Test SQLite"}'
echo -e "\n==== 创建数据源缺少参数 ===="
curl -s -X POST $API_URL/datasources -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d '{"name":""}'
echo -e "\n==== 创建查询 ===="
DS_ID=$(curl -s -X GET $API_URL/datasources -H "Authorization: Bearer $TOKEN" | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')
curl -s -X POST $API_URL/queries -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d '{"name":"Test Query","data_source_id":'$DS_ID',"sql":"SELECT 1"}'
echo -e "\n==== 查询列表 ===="
curl -s -X GET $API_URL/queries -H "Authorization: Bearer $TOKEN"
echo -e "\n==== 未登录访问 ===="
curl -s -X GET $API_URL/queries
echo -e "\n==== 权限校验 ===="
curl -s -X GET $API_URL/users -H "Authorization: Bearer $TOKEN"
echo -e "\n==== 错误响应格式 ===="
curl -s -X POST $API_URL/queries -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d '{"name":""}'
echo -e "\n==== 测试完成 ====" 