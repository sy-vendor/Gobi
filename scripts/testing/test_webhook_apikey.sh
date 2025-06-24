#!/bin/bash

# 综合测试脚本：Webhook 和 API Key 功能测试
# 使用方法: ./scripts/test_webhook_apikey.sh

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 配置
BASE_URL="http://localhost:8080/api"
TEST_USER="webhook_test_user"
TEST_PASS="test123456"
TEST_EMAIL="webhook@test.com"

# 计数器
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 测试结果存储
JWT_TOKEN=""
API_KEY=""
WEBHOOK_ID=""

# 打印状态函数
print_status() {
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓${NC} $2"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}✗${NC} $2"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        if [ -n "$3" ]; then
            echo -e "${RED}   Error: $3${NC}"
        fi
    fi
}

# 打印标题
print_title() {
    echo -e "\n${BLUE}=== $1 ===${NC}"
}

# 打印子标题
print_subtitle() {
    echo -e "\n${CYAN}--- $1 ---${NC}"
}

# 检查服务器是否运行
check_server() {
    print_title "检查服务器状态"
    if curl -s "$BASE_URL/healthz" > /dev/null; then
        print_status 0 "服务器运行正常"
    else
        print_status 1 "服务器未运行" "请先启动服务器: go run cmd/server/main.go"
        exit 1
    fi
}

# 用户注册和登录
setup_user() {
    print_title "用户设置"
    
    # 注册用户
    echo -n "注册测试用户... "
    REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
        -H "Content-Type: application/json" \
        -d "{
            \"username\": \"$TEST_USER\",
            \"password\": \"$TEST_PASS\",
            \"email\": \"$TEST_EMAIL\"
        }")
    
    if echo "$REGISTER_RESPONSE" | grep -q "already exists"; then
        print_status 0 "用户已存在"
    elif echo "$REGISTER_RESPONSE" | grep -q "id"; then
        print_status 0 "用户注册成功"
    else
        print_status 1 "用户注册失败" "$REGISTER_RESPONSE"
    fi
    
    # 登录获取 JWT Token
    echo -n "用户登录... "
    LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
        -H "Content-Type: application/json" \
        -d "{
            \"username\": \"$TEST_USER\",
            \"password\": \"$TEST_PASS\"
        }")
    
    JWT_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    
    if [ -n "$JWT_TOKEN" ]; then
        print_status 0 "登录成功，获取 JWT Token"
        echo "  Token: ${JWT_TOKEN:0:20}..."
    else
        print_status 1 "登录失败" "$LOGIN_RESPONSE"
        exit 1
    fi
}

# API Key 测试
test_api_keys() {
    print_title "API Key 功能测试"
    
    # 创建 API Key
    print_subtitle "创建 API Key"
    echo -n "创建 API Key... "
    CREATE_KEY_RESPONSE=$(curl -s -X POST "$BASE_URL/apikeys" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $JWT_TOKEN" \
        -d '{
            "name": "测试 API Key",
            "expires_at": "2025-12-31T23:59:59Z"
        }')
    
    API_KEY=$(echo "$CREATE_KEY_RESPONSE" | grep -o '"api_key":"[^"]*"' | cut -d'"' -f4)
    
    if [ -n "$API_KEY" ]; then
        print_status 0 "API Key 创建成功"
        echo "  Key: ${API_KEY:0:20}..."
    else
        print_status 1 "API Key 创建失败" "$CREATE_KEY_RESPONSE"
        return 1
    fi
    
    # 使用 API Key 访问接口
    print_subtitle "API Key 认证测试"
    echo -n "使用 API Key 访问查询接口... "
    API_KEY_RESPONSE=$(curl -s -X GET "$BASE_URL/queries" \
        -H "Authorization: ApiKey $API_KEY")
    
    if echo "$API_KEY_RESPONSE" | grep -q "\[\]"; then
        print_status 0 "API Key 认证成功"
    else
        print_status 1 "API Key 认证失败" "$API_KEY_RESPONSE"
    fi
    
    # 列出 API Keys
    print_subtitle "API Key 管理"
    echo -n "列出 API Keys... "
    LIST_KEYS_RESPONSE=$(curl -s -X GET "$BASE_URL/apikeys" \
        -H "Authorization: Bearer $JWT_TOKEN")
    
    if echo "$LIST_KEYS_RESPONSE" | grep -q "测试 API Key"; then
        print_status 0 "API Key 列表获取成功"
    else
        print_status 1 "API Key 列表获取失败" "$LIST_KEYS_RESPONSE"
    fi
}

# Webhook 测试
test_webhooks() {
    print_title "Webhook 功能测试"
    
    # 创建 Webhook
    print_subtitle "创建 Webhook"
    echo -n "创建 Webhook... "
    CREATE_WEBHOOK_RESPONSE=$(curl -s -X POST "$BASE_URL/webhooks" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $JWT_TOKEN" \
        -d '{
            "name": "测试 Webhook",
            "url": "https://httpbin.org/post",
            "events": ["report.generated", "report.failed", "webhook.test"],
            "headers": {
                "X-Test-Header": "test-value"
            }
        }')
    
    WEBHOOK_ID=$(echo "$CREATE_WEBHOOK_RESPONSE" | grep -o '"ID":[0-9]*' | cut -d':' -f2)
    
    if [ -n "$WEBHOOK_ID" ]; then
        print_status 0 "Webhook 创建成功"
        echo "  Webhook ID: $WEBHOOK_ID"
    else
        print_status 1 "Webhook 创建失败" "$CREATE_WEBHOOK_RESPONSE"
        return 1
    fi
    
    # 测试 Webhook
    print_subtitle "测试 Webhook"
    echo -n "发送测试 Webhook... "
    TEST_WEBHOOK_RESPONSE=$(curl -s -X POST "$BASE_URL/webhooks/$WEBHOOK_ID/test" \
        -H "Authorization: Bearer $JWT_TOKEN")
    
    if echo "$TEST_WEBHOOK_RESPONSE" | grep -q "sent successfully"; then
        print_status 0 "Webhook 测试发送成功"
    else
        print_status 1 "Webhook 测试发送失败" "$TEST_WEBHOOK_RESPONSE"
    fi
    
    # 查看 Webhook 发送记录
    print_subtitle "Webhook 发送记录"
    echo -n "获取发送记录... "
    DELIVERIES_RESPONSE=$(curl -s -X GET "$BASE_URL/webhooks/$WEBHOOK_ID/deliveries" \
        -H "Authorization: Bearer $JWT_TOKEN")
    
    if echo "$DELIVERIES_RESPONSE" | grep -q "webhook.test"; then
        print_status 0 "Webhook 发送记录获取成功"
    else
        print_status 1 "Webhook 发送记录获取失败" "$DELIVERIES_RESPONSE"
    fi
    
    # 列出 Webhooks
    print_subtitle "Webhook 管理"
    echo -n "列出 Webhooks... "
    LIST_WEBHOOKS_RESPONSE=$(curl -s -X GET "$BASE_URL/webhooks" \
        -H "Authorization: Bearer $JWT_TOKEN")
    
    if echo "$LIST_WEBHOOKS_RESPONSE" | grep -q "测试 Webhook"; then
        print_status 0 "Webhook 列表获取成功"
    else
        print_status 1 "Webhook 列表获取失败" "$LIST_WEBHOOKS_RESPONSE"
    fi
}

# 创建测试数据源和报告
create_test_data() {
    print_title "创建测试数据"
    
    # 创建数据源
    print_subtitle "创建测试数据源"
    echo -n "创建 SQLite 数据源... "
    DS_RESPONSE=$(curl -s -X POST "$BASE_URL/datasources" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $JWT_TOKEN" \
        -d '{
            "name": "测试数据源",
            "type": "sqlite",
            "database": "test.db",
            "description": "用于测试的数据源"
        }')
    
    if echo "$DS_RESPONSE" | grep -q "id"; then
        print_status 0 "数据源创建成功"
    else
        print_status 1 "数据源创建失败" "$DS_RESPONSE"
    fi
    
    # 创建查询
    print_subtitle "创建测试查询"
    echo -n "创建测试查询... "
    QUERY_RESPONSE=$(curl -s -X POST "$BASE_URL/queries" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $JWT_TOKEN" \
        -d '{
            "name": "测试查询",
            "data_source_id": 1,
            "sql": "SELECT 1 as test_value",
            "description": "用于测试的查询"
        }')
    
    if echo "$QUERY_RESPONSE" | grep -q "id"; then
        print_status 0 "查询创建成功"
    else
        print_status 1 "查询创建失败" "$QUERY_RESPONSE"
    fi
}

# 测试报告生成触发 Webhook
test_report_webhook() {
    print_title "报告生成 Webhook 测试"
    
    # 创建报告计划
    print_subtitle "创建报告计划"
    echo -n "创建报告计划... "
    SCHEDULE_RESPONSE=$(curl -s -X POST "$BASE_URL/reports/schedules" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $JWT_TOKEN" \
        -d '{
            "name": "测试报告计划",
            "type": "daily",
            "query_ids": "[1]",
            "chart_ids": "[]",
            "template_ids": "[]",
            "cron_pattern": "0 12 * * *"
        }')
    
    SCHEDULE_ID=$(echo "$SCHEDULE_RESPONSE" | grep -o '"ID":[0-9]*' | cut -d':' -f2)
    
    if [ -n "$SCHEDULE_ID" ]; then
        print_status 0 "报告计划创建成功"
        echo "  Schedule ID: $SCHEDULE_ID"
    else
        print_status 1 "报告计划创建失败" "$SCHEDULE_RESPONSE"
        return 1
    fi
    
    # 手动生成报告（会触发 webhook）
    print_subtitle "手动生成报告"
    echo -n "生成报告... "
    GENERATE_RESPONSE=$(curl -s -X POST "$BASE_URL/reports/generate/excel" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $JWT_TOKEN" \
        -d '{
            "name": "测试报告",
            "query_ids": [1],
            "chart_ids": [],
            "template_id": null
        }')
    
    if echo "$GENERATE_RESPONSE" | grep -q "success"; then
        print_status 0 "报告生成成功"
    else
        print_status 1 "报告生成失败" "$GENERATE_RESPONSE"
    fi
    
    # 等待一下让 webhook 发送
    sleep 2
    
    # 检查 webhook 发送记录
    print_subtitle "检查报告 Webhook"
    echo -n "检查报告生成 Webhook... "
    REPORT_DELIVERIES=$(curl -s -X GET "$BASE_URL/webhooks/$WEBHOOK_ID/deliveries" \
        -H "Authorization: Bearer $JWT_TOKEN")
    
    if echo "$REPORT_DELIVERIES" | grep -q "report.generated"; then
        print_status 0 "报告生成 Webhook 发送成功"
    else
        print_status 1 "报告生成 Webhook 未发送" "可能需要检查 webhook 配置"
    fi
}

# 清理测试数据
cleanup() {
    print_title "清理测试数据"
    
    if [ -n "$WEBHOOK_ID" ]; then
        echo -n "删除测试 Webhook... "
        DELETE_WEBHOOK_RESPONSE=$(curl -s -X DELETE "$BASE_URL/webhooks/$WEBHOOK_ID" \
            -H "Authorization: Bearer $JWT_TOKEN")
        
        if echo "$DELETE_WEBHOOK_RESPONSE" | grep -q "deleted successfully"; then
            print_status 0 "Webhook 删除成功"
        else
            print_status 1 "Webhook 删除失败" "$DELETE_WEBHOOK_RESPONSE"
        fi
    fi
    
    # 获取 API Key ID 并删除
    if [ -n "$JWT_TOKEN" ]; then
        echo -n "删除测试 API Key... "
        KEYS_RESPONSE=$(curl -s -X GET "$BASE_URL/apikeys" \
            -H "Authorization: Bearer $JWT_TOKEN")
        
        KEY_ID=$(echo "$KEYS_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
        
        if [ -n "$KEY_ID" ]; then
            DELETE_KEY_RESPONSE=$(curl -s -X DELETE "$BASE_URL/apikeys/$KEY_ID" \
                -H "Authorization: Bearer $JWT_TOKEN")
            
            if echo "$DELETE_KEY_RESPONSE" | grep -q "revoked"; then
                print_status 0 "API Key 删除成功"
            else
                print_status 1 "API Key 删除失败" "$DELETE_KEY_RESPONSE"
            fi
        fi
    fi
}

# 显示测试结果
show_results() {
    print_title "测试结果汇总"
    
    echo -e "${YELLOW}总测试数: $TOTAL_TESTS${NC}"
    echo -e "${GREEN}通过: $PASSED_TESTS${NC}"
    echo -e "${RED}失败: $FAILED_TESTS${NC}"
    
    if [ $FAILED_TESTS -eq 0 ]; then
        echo -e "\n${GREEN}🎉 所有测试通过！${NC}"
    else
        echo -e "\n${RED}❌ 有 $FAILED_TESTS 个测试失败${NC}"
    fi
    
    echo -e "\n${BLUE}测试完成时间: $(date)${NC}"
}

# 主函数
main() {
    echo -e "${PURPLE}🚀 Gobi Webhook & API Key 功能测试${NC}"
    echo -e "${PURPLE}================================${NC}"
    
    check_server
    setup_user
    test_api_keys
    test_webhooks
    create_test_data
    test_report_webhook
    cleanup
    show_results
}

# 运行主函数
main "$@" 