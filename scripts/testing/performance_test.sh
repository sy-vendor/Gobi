#!/bin/bash

# 性能测试脚本
# 使用方法: ./scripts/testing/performance_test.sh

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
BASE_URL="http://localhost:8080/api"
TEST_USER="perf_test_user"
TEST_PASS="test123456"
TEST_EMAIL="perf@test.com"

# 计数器
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 打印状态函数
print_status() {
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓${NC} $2"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}✗${NC} $2"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
}

# 打印标题
print_title() {
    echo -e "\n${BLUE}=== $1 ===${NC}"
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
setup_auth() {
    print_title "设置认证"
    
    # 注册用户
    echo -n "注册测试用户... "
    REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
        -H "Content-Type: application/json" \
        -d "{
            \"username\": \"$TEST_USER\",
            \"password\": \"$TEST_PASS\",
            \"email\": \"$TEST_EMAIL\"
        }")
    
    if echo "$REGISTER_RESPONSE" | grep -q "id"; then
        print_status 0 "用户注册成功"
    else
        print_status 1 "用户注册失败" "$REGISTER_RESPONSE"
    fi
    
    # 用户登录
    echo -n "用户登录... "
    LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
        -H "Content-Type: application/json" \
        -d "{
            \"username\": \"$TEST_USER\",
            \"password\": \"$TEST_PASS\"
        }")
    
    JWT_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    
    if [ -n "$JWT_TOKEN" ]; then
        print_status 0 "用户登录成功"
        echo "  Token: ${JWT_TOKEN:0:20}..."
    else
        print_status 1 "用户登录失败" "$LOGIN_RESPONSE"
        exit 1
    fi
}

# 创建测试数据
create_test_data() {
    print_title "创建测试数据"
    
    # 创建数据源
    echo -n "创建测试数据源... "
    DS_RESPONSE=$(curl -s -X POST "$BASE_URL/datasources" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $JWT_TOKEN" \
        -d '{
            "name": "性能测试数据源",
            "type": "sqlite",
            "database": "test.db",
            "description": "用于性能测试的数据源"
        }')
    
    if echo "$DS_RESPONSE" | grep -q "id"; then
        print_status 0 "数据源创建成功"
    else
        print_status 1 "数据源创建失败" "$DS_RESPONSE"
    fi
    
    # 创建多个查询
    echo -n "创建测试查询... "
    for i in {1..5}; do
        QUERY_RESPONSE=$(curl -s -X POST "$BASE_URL/queries" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $JWT_TOKEN" \
            -d "{
                \"name\": \"性能测试查询 $i\",
                \"data_source_id\": 1,
                \"sql\": \"SELECT $i as test_value, 'query_$i' as query_name\",
                \"description\": \"性能测试查询 $i\"
            }")
        
        if echo "$QUERY_RESPONSE" | grep -q "id"; then
            print_status 0 "查询 $i 创建成功"
        else
            print_status 1 "查询 $i 创建失败" "$QUERY_RESPONSE"
        fi
    done
}

# 性能测试
performance_test() {
    print_title "性能测试"
    
    # 测试查询执行性能
    echo -n "测试查询执行性能... "
    start_time=$(date +%s.%N)
    
    for i in {1..10}; do
        curl -s -X POST "$BASE_URL/queries/1/execute" \
            -H "Authorization: Bearer $JWT_TOKEN" > /dev/null
    done
    
    end_time=$(date +%s.%N)
    execution_time=$(echo "$end_time - $start_time" | bc)
    avg_time=$(echo "scale=3; $execution_time / 10" | bc)
    
    print_status 0 "查询执行测试完成"
    echo "  总执行时间: ${execution_time}s"
    echo "  平均执行时间: ${avg_time}s"
    
    # 测试缓存性能
    echo -n "测试缓存性能... "
    start_time=$(date +%s.%N)
    
    for i in {1..20}; do
        curl -s -X POST "$BASE_URL/queries/1/execute" \
            -H "Authorization: Bearer $JWT_TOKEN" > /dev/null
    done
    
    end_time=$(date +%s.%N)
    execution_time=$(echo "$end_time - $start_time" | bc)
    avg_time=$(echo "scale=3; $execution_time / 20" | bc)
    
    print_status 0 "缓存性能测试完成"
    echo "  总执行时间: ${execution_time}s"
    echo "  平均执行时间: ${avg_time}s"
    
    # 测试并发性能
    echo -n "测试并发性能... "
    start_time=$(date +%s.%N)
    
    # 并发执行10个请求
    for i in {1..10}; do
        curl -s -X POST "$BASE_URL/queries/1/execute" \
            -H "Authorization: Bearer $JWT_TOKEN" > /dev/null &
    done
    wait
    
    end_time=$(date +%s.%N)
    execution_time=$(echo "$end_time - $start_time" | bc)
    
    print_status 0 "并发性能测试完成"
    echo "  并发执行时间: ${execution_time}s"
}

# 系统监控测试
system_monitoring_test() {
    print_title "系统监控测试"
    
    # 获取系统统计信息
    echo -n "获取系统统计信息... "
    STATS_RESPONSE=$(curl -s -X GET "$BASE_URL/system/stats" \
        -H "Authorization: Bearer $JWT_TOKEN")
    
    if echo "$STATS_RESPONSE" | grep -q "database"; then
        print_status 0 "系统统计获取成功"
        echo "  数据库连接池信息: $(echo "$STATS_RESPONSE" | jq -r '.database.total_pools') 个连接池"
        echo "  缓存状态: $(echo "$STATS_RESPONSE" | jq -r '.cache.enabled')"
    else
        print_status 1 "系统统计获取失败" "$STATS_RESPONSE"
    fi
    
    # 获取仪表板统计信息
    echo -n "获取仪表板统计信息... "
    DASHBOARD_RESPONSE=$(curl -s -X GET "$BASE_URL/dashboard/stats" \
        -H "Authorization: Bearer $JWT_TOKEN")
    
    if echo "$DASHBOARD_RESPONSE" | grep -q "totalQueries"; then
        print_status 0 "仪表板统计获取成功"
        echo "  总查询数: $(echo "$DASHBOARD_RESPONSE" | jq -r '.totalQueries')"
        echo "  今日查询数: $(echo "$DASHBOARD_RESPONSE" | jq -r '.todayQueries')"
    else
        print_status 1 "仪表板统计获取失败" "$DASHBOARD_RESPONSE"
    fi
}

# 压力测试
stress_test() {
    print_title "压力测试"
    
    # 模拟高并发查询
    echo -n "执行压力测试 (50个并发请求)... "
    start_time=$(date +%s.%N)
    
    for i in {1..50}; do
        curl -s -X POST "$BASE_URL/queries/1/execute" \
            -H "Authorization: Bearer $JWT_TOKEN" > /dev/null &
    done
    wait
    
    end_time=$(date +%s.%N)
    execution_time=$(echo "$end_time - $start_time" | bc)
    rps=$(echo "scale=2; 50 / $execution_time" | bc)
    
    print_status 0 "压力测试完成"
    echo "  总执行时间: ${execution_time}s"
    echo "  请求处理速率: ${rps} RPS"
}

# 清理测试数据
cleanup_test_data() {
    print_title "清理测试数据"
    
    echo -n "清理测试查询... "
    for i in {1..5}; do
        curl -s -X DELETE "$BASE_URL/queries/$i" \
            -H "Authorization: Bearer $JWT_TOKEN" > /dev/null
    done
    print_status 0 "测试查询清理完成"
    
    echo -n "清理测试数据源... "
    curl -s -X DELETE "$BASE_URL/datasources/1" \
        -H "Authorization: Bearer $JWT_TOKEN" > /dev/null
    print_status 0 "测试数据源清理完成"
}

# 打印测试结果
print_results() {
    print_title "测试结果"
    
    echo "总测试数: $TOTAL_TESTS"
    echo -e "通过: ${GREEN}$PASSED_TESTS${NC}"
    echo -e "失败: ${RED}$FAILED_TESTS${NC}"
    
    if [ $FAILED_TESTS -eq 0 ]; then
        echo -e "\n${GREEN}所有测试通过！${NC}"
    else
        echo -e "\n${RED}有 $FAILED_TESTS 个测试失败${NC}"
    fi
}

# 主函数
main() {
    echo -e "${BLUE}Gobi 性能测试开始${NC}"
    
    check_server
    setup_auth
    create_test_data
    performance_test
    system_monitoring_test
    stress_test
    cleanup_test_data
    print_results
    
    echo -e "\n${BLUE}性能测试完成${NC}"
}

# 运行主函数
main 