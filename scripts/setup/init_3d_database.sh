#!/bin/bash

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== 3D Charts Database Initialization ===${NC}"
echo

# 获取脚本所在目录的上级目录（项目根目录）
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"

# 检查SQLite数据库文件是否存在
DB_FILE="gobi.db"
if [ ! -f "$DB_FILE" ]; then
    echo -e "${YELLOW}Warning: Database file $DB_FILE not found.${NC}"
    echo "Please start the Gobi server first to create the database."
    echo "Run: go run cmd/server/main.go"
    exit 1
fi

echo "1. Checking database connection..."
if ! sqlite3 "$DB_FILE" "SELECT 1;" > /dev/null 2>&1; then
    echo -e "${RED}Error: Cannot connect to database $DB_FILE${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Database connection successful${NC}"

echo "2. Creating 3D charts sample tables..."

# 执行SQL脚本
if sqlite3 "$DB_FILE" < "$PROJECT_ROOT/scripts/data/generate_3d_sample_data.sql"; then
    echo -e "${GREEN}✓ 3D charts sample data created successfully${NC}"
else
    echo -e "${RED}✗ Failed to create sample data${NC}"
    exit 1
fi

echo "3. Verifying data creation..."
echo "Table record counts:"
sqlite3 "$DB_FILE" "
SELECT 'sales_3d' as table_name, COUNT(*) as record_count FROM sales_3d
UNION ALL
SELECT 'products_3d' as table_name, COUNT(*) as record_count FROM products_3d
UNION ALL
SELECT 'terrain_3d' as table_name, COUNT(*) as record_count FROM terrain_3d
UNION ALL
SELECT 'cities_3d' as table_name, COUNT(*) as record_count FROM cities_3d;
"

echo
echo -e "${GREEN}=== Database Initialization Complete ===${NC}"
echo "You can now test 3D charts with the sample data."
echo
echo "Sample queries you can use:"
echo "1. 3D Bar Chart: SELECT category as x, region as y, SUM(amount) as z FROM sales_3d GROUP BY category, region"
echo "2. 3D Scatter: SELECT performance_score as x, price as y, customer_rating as z, product_category as category FROM products_3d"
echo "3. 3D Surface: SELECT longitude as x, latitude as y, elevation as z FROM terrain_3d"
echo "4. 3D Bubble: SELECT gdp as x, population as y, area as z, city_name as category FROM cities_3d" 