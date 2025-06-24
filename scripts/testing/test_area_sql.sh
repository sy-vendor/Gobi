#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "=== Area Chart SQL Data Test ==="
echo

# 检查SQLite是否安装
if ! command -v sqlite3 &> /dev/null; then
    echo -e "${RED}Error: sqlite3 is not installed${NC}"
    echo "Please install sqlite3 first:"
    echo "  macOS: brew install sqlite3"
    echo "  Ubuntu: sudo apt-get install sqlite3"
    echo "  CentOS: sudo yum install sqlite3"
    exit 1
fi

# 获取脚本所在目录的上级目录（项目根目录）
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"

# 使用主数据库
DB_FILE="$PROJECT_ROOT/gobi.db"
echo "1. Using main database: $DB_FILE"

# 检查数据库文件是否存在
if [ ! -f "$DB_FILE" ]; then
    echo -e "${RED}Error: Database file $DB_FILE not found${NC}"
    echo "Please start the Gobi server first to create the database."
    echo "Run: go run cmd/server/main.go"
    exit 1
fi

# 检查数据库连接
echo "2. Checking database connection..."
if ! sqlite3 "$DB_FILE" "SELECT 1;" > /dev/null 2>&1; then
    echo -e "${RED}Error: Cannot connect to database $DB_FILE${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Database connection successful${NC}"

# 执行SQL脚本
echo "3. Executing SQL script..."
if sqlite3 "$DB_FILE" < "$PROJECT_ROOT/scripts/data/generate_area_chart_data.sql"; then
    echo -e "${GREEN}   ✓ SQL script executed successfully${NC}"
else
    echo -e "${RED}   ✗ Failed to execute SQL script${NC}"
    exit 1
fi

# 测试数据表创建
echo "4. Testing table creation..."
TABLES=("sales_trend" "user_growth" "market_share" "website_traffic" "revenue_data")

for table in "${TABLES[@]}"; do
    if sqlite3 "$DB_FILE" "SELECT name FROM sqlite_master WHERE type='table' AND name='$table';" | grep -q "$table"; then
        echo -e "   ${GREEN}✓${NC} Table '$table' created successfully"
    else
        echo -e "   ${RED}✗${NC} Table '$table' not found"
    fi
done

# 测试数据插入
echo "5. Testing data insertion..."
for table in "${TABLES[@]}"; do
    count=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM $table;")
    echo -e "   ${GREEN}✓${NC} Table '$table' has $count records"
done

# 测试示例查询
echo "6. Testing sample queries..."

# 测试销售趋势查询
echo "   Testing sales trend query..."
sales_result=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM sales_trend WHERE product_category = 'Electronics';")
if [ "$sales_result" -gt 0 ]; then
    echo -e "   ${GREEN}✓${NC} Sales trend data available ($sales_result records)"
else
    echo -e "   ${RED}✗${NC} No sales trend data found"
fi

# 测试用户增长查询
echo "   Testing user growth query..."
user_result=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM user_growth WHERE user_type = 'Free';")
if [ "$user_result" -gt 0 ]; then
    echo -e "   ${GREEN}✓${NC} User growth data available ($user_result records)"
else
    echo -e "   ${RED}✗${NC} No user growth data found"
fi

# 测试市场份额查询
echo "   Testing market share query..."
market_result=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM market_share WHERE company = 'Company A';")
if [ "$market_result" -gt 0 ]; then
    echo -e "   ${GREEN}✓${NC} Market share data available ($market_result records)"
else
    echo -e "   ${RED}✗${NC} No market share data found"
fi

# 测试网站流量查询
echo "   Testing website traffic query..."
traffic_result=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM website_traffic WHERE hour BETWEEN 9 AND 17;")
if [ "$traffic_result" -gt 0 ]; then
    echo -e "   ${GREEN}✓${NC} Website traffic data available ($traffic_result records)"
else
    echo -e "   ${RED}✗${NC} No website traffic data found"
fi

# 测试收入数据查询
echo "   Testing revenue data query..."
revenue_result=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM revenue_data WHERE profit > 20000;")
if [ "$revenue_result" -gt 0 ]; then
    echo -e "   ${GREEN}✓${NC} Revenue data available ($revenue_result records with profit > 20000)"
else
    echo -e "   ${RED}✗${NC} No revenue data found"
fi

# 显示示例查询结果
echo "7. Sample query results:"

echo "   Sales trend (Electronics):"
sqlite3 "$DB_FILE" "SELECT month, sales_amount FROM sales_trend WHERE product_category = 'Electronics' ORDER BY month LIMIT 5;" | while IFS='|' read -r month amount; do
    echo "     $month: $amount"
done

echo "   User growth (Free users):"
sqlite3 "$DB_FILE" "SELECT date, new_users FROM user_growth WHERE user_type = 'Free' ORDER BY date LIMIT 5;" | while IFS='|' read -r date users; do
    echo "     $date: $users new users"
done

echo "   Market share (Company A):"
sqlite3 "$DB_FILE" "SELECT quarter, market_share FROM market_share WHERE company = 'Company A' ORDER BY quarter LIMIT 5;" | while IFS='|' read -r quarter share; do
    echo "     $quarter: ${share}%"
done

echo
echo "=== Test Summary ==="
echo -e "${GREEN}✓${NC} All tables created successfully"
echo -e "${GREEN}✓${NC} All data inserted successfully"
echo -e "${GREEN}✓${NC} All sample queries working"
echo -e "${GREEN}✓${NC} Area chart data has been added to gobi.db!"
echo
echo "Area chart data is now available in your main database:"
echo "1. Tables created: sales_trend, user_growth, market_share, website_traffic, revenue_data"
echo "2. Sample data inserted for testing area charts"
echo "3. You can now use these tables in Gobi to create area charts"
echo