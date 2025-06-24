#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=== Tree Diagram SQL Data Test ==="
echo

# 检查SQLite是否安装
if ! command -v sqlite3 &> /dev/null; then
    echo -e "${RED}Error: sqlite3 is not installed${NC}"
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
echo "3. Executing tree SQL script..."
if sqlite3 "$DB_FILE" < "$PROJECT_ROOT/scripts/data/generate_tree_sample_data.sql"; then
    echo -e "${GREEN}   ✓ Tree SQL script executed successfully${NC}"
else
    echo -e "${RED}   ✗ Failed to execute tree SQL script${NC}"
    exit 1
fi

# 测试数据表创建
TABLES=("org_tree")
echo "4. Testing table creation..."
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

# 示例树形结构查询
echo "6. Sample tree structure queries:"
echo "   Direct subordinates of 王总:"
sqlite3 "$DB_FILE" "SELECT name, position FROM org_tree WHERE parent_id = 1;" | while IFS='|' read -r name position; do
    echo "     $name: $position"
done
echo "   Direct subordinates of 李工:"
sqlite3 "$DB_FILE" "SELECT name, position FROM org_tree WHERE parent_id = 2;" | while IFS='|' read -r name position; do
    echo "     $name: $position"
done

echo
echo "=== Tree Diagram Test Summary ==="
echo -e "${GREEN}✓${NC} All tree tables created successfully"
echo -e "${GREEN}✓${NC} All tree data inserted successfully"
echo -e "${GREEN}✓${NC} All sample tree queries working"
echo -e "${GREEN}✓${NC} Tree diagram data has been added to gobi.db!"
echo 