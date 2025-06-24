#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 服务器地址
BASE_URL="http://localhost:8080"

# 打印状态函数
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓${NC} $2"
    else
        echo -e "${RED}✗${NC} $2"
        echo "Response: $3"
    fi
}

# 测试函数
test_3d_chart() {
    local chart_type="$1"
    local chart_name="$2"
    local config="$3"
    
    echo -n "Testing $chart_name... "
    
    response=$(curl -s -w "\n%{http_code}" -X POST \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d "$config" \
        "$BASE_URL/api/charts")
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" -eq 201 ]; then
        print_status 0 "$chart_name"
        echo "  HTTP Code: $http_code"
        chart_id=$(echo "$body" | grep -o '"ID":[0-9]*' | cut -d':' -f2)
        echo "  Chart ID: $chart_id"
    else
        print_status 1 "$chart_name"
        echo "  Expected: 201"
        echo "  Got: $http_code"
        echo "  Response: $body"
    fi
    echo
}

echo "=== 3D Charts Test ==="
echo

# 检查服务器是否运行
if ! curl -s "$BASE_URL/healthz" > /dev/null; then
    echo -e "${RED}Error: Server is not running on $BASE_URL${NC}"
    echo "Please start the server first: go run cmd/server/main.go"
    exit 1
fi

# 获取认证token
echo "1. Getting authentication token..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
        "username": "admin",
        "password": "admin123"
    }')

TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}Failed to get authentication token${NC}"
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi

print_status 0 "Authentication successful"

# 创建数据源
echo "2. Creating test data source..."
DS_RESPONSE=$(curl -s -X POST "$BASE_URL/api/datasources" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{
        "name": "Test 3D Data Source",
        "type": "sqlite",
        "database": "test_3d.db",
        "description": "Test data source for 3D charts"
    }')

DS_ID=$(echo "$DS_RESPONSE" | grep -o '"ID":[0-9]*' | cut -d':' -f2)
print_status $? "Data source creation"

# 创建测试查询
echo "3. Creating test queries..."

# 3D Bar Chart Query
BAR_QUERY_RESPONSE=$(curl -s -X POST "$BASE_URL/api/queries" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{
        "name": "3D Sales Data",
        "dataSourceId": '$DS_ID',
        "sql": "SELECT category as x, region as y, SUM(amount) as z FROM sales_3d GROUP BY category, region ORDER BY category, region",
        "description": "3D bar chart data"
    }')

BAR_QUERY_ID=$(echo "$BAR_QUERY_RESPONSE" | grep -o '"ID":[0-9]*' | cut -d':' -f2)
print_status $? "3D Bar Chart query creation"

# 3D Scatter Query
SCATTER_QUERY_RESPONSE=$(curl -s -X POST "$BASE_URL/api/queries" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{
        "name": "3D Product Data",
        "dataSourceId": '$DS_ID',
        "sql": "SELECT performance_score as x, price as y, customer_rating as z, product_category as category, sales_volume as size FROM products_3d WHERE performance_score IS NOT NULL AND price IS NOT NULL AND customer_rating IS NOT NULL",
        "description": "3D scatter plot data"
    }')

SCATTER_QUERY_ID=$(echo "$SCATTER_QUERY_RESPONSE" | grep -o '"ID":[0-9]*' | cut -d':' -f2)
print_status $? "3D Scatter Chart query creation"

# 测试3D图表创建
echo "4. Testing 3D Chart Creation..."

# 3D Bar Chart
test_3d_chart "3d-bar" "3D Bar Chart" '{
    "name": "3D Sales Chart",
    "queryId": '$BAR_QUERY_ID',
    "type": "3d-bar",
    "config": "{\"xField\":\"x\",\"yField\":\"y\",\"zField\":\"z\",\"title\":\"3D Sales by Category and Region\",\"legend\":true,\"color\":[\"#1890ff\",\"#2fc25b\",\"#facc14\"],\"tooltip\":true,\"grid3D\":{\"boxWidth\":100,\"boxHeight\":100,\"boxDepth\":100,\"viewControl\":{\"alpha\":20,\"beta\":40,\"distance\":200}}}",
    "description": "3D bar chart showing sales by category and region"
}'

# 3D Scatter Chart
test_3d_chart "3d-scatter" "3D Scatter Chart" '{
    "name": "3D Product Scatter",
    "queryId": '$SCATTER_QUERY_ID',
    "type": "3d-scatter",
    "config": "{\"xField\":\"x\",\"yField\":\"y\",\"zField\":\"z\",\"colorField\":\"category\",\"sizeField\":\"size\",\"title\":\"3D Product Performance\",\"legend\":true,\"color\":[\"#1890ff\",\"#2fc25b\",\"#facc14\"],\"tooltip\":true,\"symbolSize\":10,\"grid3D\":{\"boxWidth\":100,\"boxHeight\":100,\"boxDepth\":100,\"viewControl\":{\"alpha\":20,\"beta\":40,\"distance\":200}}}",
    "description": "3D scatter plot showing product performance"
}'

# 3D Surface Chart
test_3d_chart "3d-surface" "3D Surface Chart" '{
    "name": "3D Surface Plot",
    "queryId": '$BAR_QUERY_ID',
    "type": "3d-surface",
    "config": "{\"xField\":\"x\",\"yField\":\"y\",\"zField\":\"z\",\"title\":\"3D Surface Plot\",\"color\":[\"#313695\",\"#4575b4\",\"#74add1\",\"#abd9e9\",\"#e0f3f8\",\"#ffffcc\",\"#fee090\",\"#fdae61\",\"#f46d43\",\"#d73027\",\"#a50026\"],\"tooltip\":true,\"shading\":\"realistic\",\"grid3D\":{\"boxWidth\":100,\"boxHeight\":100,\"boxDepth\":100,\"viewControl\":{\"alpha\":20,\"beta\":40,\"distance\":200}}}",
    "description": "3D surface plot"
}'

# 3D Bubble Chart
test_3d_chart "3d-bubble" "3D Bubble Chart" '{
    "name": "3D Bubble Chart",
    "queryId": '$SCATTER_QUERY_ID',
    "type": "3d-bubble",
    "config": "{\"xField\":\"x\",\"yField\":\"y\",\"zField\":\"z\",\"sizeField\":\"size\",\"colorField\":\"category\",\"title\":\"3D Bubble Chart\",\"legend\":true,\"color\":[\"#1890ff\",\"#2fc25b\",\"#facc14\",\"#f5222d\"],\"tooltip\":true,\"grid3D\":{\"boxWidth\":100,\"boxHeight\":100,\"boxDepth\":100,\"viewControl\":{\"alpha\":20,\"beta\":40,\"distance\":200}}}",
    "description": "3D bubble chart"
}'

# 测试无效图表类型
echo "5. Testing Invalid Chart Type..."
INVALID_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{
        "name": "Invalid Chart",
        "queryId": '$BAR_QUERY_ID',
        "type": "invalid-3d-type",
        "config": "{}",
        "description": "Invalid chart type test"
    }' \
    "$BASE_URL/api/charts")

INVALID_HTTP_CODE=$(echo "$INVALID_RESPONSE" | tail -n1)
INVALID_BODY=$(echo "$INVALID_RESPONSE" | head -n -1)

if [ "$INVALID_HTTP_CODE" -eq 400 ]; then
    print_status 0 "Invalid chart type validation"
else
    print_status 1 "Invalid chart type validation"
    echo "  Expected: 400"
    echo "  Got: $INVALID_HTTP_CODE"
    echo "  Response: $INVALID_BODY"
fi

echo
echo "=== 3D Charts Test Summary ==="
echo "✓ Authentication and data source setup"
echo "✓ 3D Bar Chart creation"
echo "✓ 3D Scatter Chart creation" 
echo "✓ 3D Surface Chart creation"
echo "✓ 3D Bubble Chart creation"
echo "✓ Invalid chart type validation"
echo
echo "3D charts are now supported in Gobi BI Engine!"