#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=== Box Plot SQL Data Test ==="
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
echo "3. Executing SQL script..."
if sqlite3 "$DB_FILE" < "$PROJECT_ROOT/scripts/data/generate_boxplot_sample_data.sql"; then
    echo -e "${GREEN}   ✓ SQL script executed successfully${NC}"
else
    echo -e "${RED}   ✗ Failed to execute SQL script${NC}"
    exit 1
fi

# 测试数据表创建
echo "4. Testing table creation..."
TABLES=("student_scores" "product_performance")

for table in "${TABLES[@]}"; do
    echo "   Testing $table table..."
    if sqlite3 "$DB_FILE" "SELECT name FROM sqlite_master WHERE type='table' AND name='$table';" | grep -q "$table"; then
        echo -e "   ${GREEN}✓${NC} $table table created successfully"
    else
        echo -e "   ${RED}✗${NC} $table table not found"
        exit 1
    fi
done

# 测试数据插入
echo "5. Testing data insertion..."

echo "   Testing student_scores table..."
score_count=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM student_scores;")
if [ "$score_count" -gt 0 ]; then
    echo -e "   ${GREEN}✓${NC} student_scores has $score_count records"
else
    echo -e "   ${RED}✗${NC} No student_scores data found"
    exit 1
fi

echo "   Testing product_performance table..."
perf_count=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM product_performance;")
if [ "$perf_count" -gt 0 ]; then
    echo -e "   ${GREEN}✓${NC} product_performance has $perf_count records"
else
    echo -e "   ${RED}✗${NC} No product_performance data found"
    exit 1
fi

# 测试箱线图数据查询
echo "6. Testing box plot data queries..."

echo "   Testing student scores by class and subject..."
echo "   A班数学成绩分布:"
sqlite3 "$DB_FILE" "SELECT class, subject, MIN(score) as min, AVG(score) as avg, MAX(score) as max, COUNT(*) as count FROM student_scores WHERE class='A班' AND subject='数学' GROUP BY class, subject;"

echo "   B班英语成绩分布:"
sqlite3 "$DB_FILE" "SELECT class, subject, MIN(score) as min, AVG(score) as avg, MAX(score) as max, COUNT(*) as count FROM student_scores WHERE class='B班' AND subject='英语' GROUP BY class, subject;"

echo "   C班物理成绩分布:"
sqlite3 "$DB_FILE" "SELECT class, subject, MIN(score) as min, AVG(score) as avg, MAX(score) as max, COUNT(*) as count FROM student_scores WHERE class='C班' AND subject='物理' GROUP BY class, subject;"

echo "   Testing product performance by type and metric..."
echo "   手机A响应时间分布:"
sqlite3 "$DB_FILE" "SELECT product_type, test_metric, MIN(value) as min, AVG(value) as avg, MAX(value) as max, COUNT(*) as count FROM product_performance WHERE product_type='手机A' AND test_metric='响应时间' GROUP BY product_type, test_metric;"

echo "   手机B电池续航分布:"
sqlite3 "$DB_FILE" "SELECT product_type, test_metric, MIN(value) as min, AVG(value) as avg, MAX(value) as max, COUNT(*) as count FROM product_performance WHERE product_type='手机B' AND test_metric='电池续航' GROUP BY product_type, test_metric;"

echo "   手机C响应时间分布:"
sqlite3 "$DB_FILE" "SELECT product_type, test_metric, MIN(value) as min, AVG(value) as avg, MAX(value) as max, COUNT(*) as count FROM product_performance WHERE product_type='手机C' AND test_metric='响应时间' GROUP BY product_type, test_metric;"

echo
echo -e "${GREEN}=== Box Plot Data Test Complete ===${NC}"
echo "Sample queries for box plots:"
echo "1. Student scores by class: SELECT class, subject, score FROM student_scores ORDER BY class, subject"
echo "2. Product performance: SELECT product_type, test_metric, value FROM product_performance ORDER BY product_type, test_metric"
echo "3. Math scores comparison: SELECT class, score FROM student_scores WHERE subject='数学' ORDER BY class"
echo "4. Response time comparison: SELECT product_type, value FROM product_performance WHERE test_metric='响应时间' ORDER BY product_type"
echo
echo "Box plot data is now ready for testing in Gobi BI Engine!" 