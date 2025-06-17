#!/bin/bash

# 设置基础URL
BASE_URL="http://localhost:8080"

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# 登录获取 token
login_and_get_token() {
    local username=$1
    local password=$2
    TOKEN=$(curl -s -X POST -H "Content-Type: application/json" -d '{"username": "'$username'", "password": "'$password'"}' "$BASE_URL/api/auth/login" | jq -r '.token')
}

# 测试函数
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local expected_status=$4
    local test_name=$5
    local auth_header=$6

    echo "Testing: $test_name"
    
    if [ -z "$data" ]; then
        if [ -z "$auth_header" ]; then
            response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint")
        else
            response=$(curl -s -w "\n%{http_code}" -X $method -H "$auth_header" "$BASE_URL$endpoint")
        fi
    else
        if [ -z "$auth_header" ]; then
            response=$(curl -s -w "\n%{http_code}" -X $method -H "Content-Type: application/json" -d "$data" "$BASE_URL$endpoint")
        else
            response=$(curl -s -w "\n%{http_code}" -X $method -H "Content-Type: application/json" -H "$auth_header" -d "$data" "$BASE_URL$endpoint")
        fi
    fi

    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')

    if [ "$status_code" -eq "$expected_status" ]; then
        echo -e "${GREEN}✓ Test passed${NC}"
    else
        echo -e "${RED}✗ Test failed${NC}"
        echo "Expected status: $expected_status"
        echo "Got status: $status_code"
        echo "Response body: $body"
    fi
    echo "----------------------------------------"
}

# 注册用户
USERNAME="testuser_$(date +%s)"
PASSWORD="test123"
test_endpoint "POST" "/api/auth/register" '{"username": "'$USERNAME'", "password": "'$PASSWORD'"}' 201 "Valid Registration"

# 登录获取 token
login_and_get_token "$USERNAME" "$PASSWORD"

# 测试未找到的资源
test_endpoint "GET" "/api/queries/999" "" 404 "Not Found Error" "Authorization: Bearer $TOKEN"

# 测试无效的JSON数据
test_endpoint "POST" "/api/queries" "invalid json" 400 "Invalid JSON Error" "Authorization: Bearer $TOKEN"

# 测试未授权访问
test_endpoint "GET" "/api/queries/1" "" 401 "Unauthorized Access"

# 测试用户注册 - 无效数据
test_endpoint "POST" "/api/auth/register" '{"username": "test"}' 400 "Invalid Registration Data"

# 测试用户注册 - 重复用户名
test_endpoint "POST" "/api/auth/register" '{"username": "$USERNAME", "password": "$PASSWORD"}' 409 "Duplicate Username"

# 测试登录 - 无效凭证
test_endpoint "POST" "/api/auth/login" '{"username": "$USERNAME", "password": "wrong"}' 401 "Invalid Credentials"

# 测试登录 - 成功
test_endpoint "POST" "/api/auth/login" '{"username": "$USERNAME", "password": "$PASSWORD"}' 200 "Valid Login"

# 测试无效的token
test_endpoint "GET" "/api/queries/1" "" 401 "Invalid Token" "Authorization: Bearer invalid_token"

# 测试文件上传 - 无文件
test_endpoint "POST" "/api/templates" "" 400 "No File Uploaded" "Authorization: Bearer $TOKEN"

echo "All tests completed!" 