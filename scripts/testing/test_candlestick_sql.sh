#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=== Candlestick Chart SQL Data Test ==="
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
if sqlite3 "$DB_FILE" < "$PROJECT_ROOT/scripts/data/generate_candlestick_sample_data.sql"; then
    echo -e "${GREEN}   ✓ SQL script executed successfully${NC}"
else
    echo -e "${RED}   ✗ Failed to execute SQL script${NC}"
    exit 1
fi

# 测试数据表创建
echo "4. Testing table creation..."
TABLES=("stock_prices" "crypto_prices")

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

echo "   Testing stock_prices table..."
stock_count=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM stock_prices;")
if [ "$stock_count" -gt 0 ]; then
    echo -e "   ${GREEN}✓${NC} stock_prices has $stock_count records"
else
    echo -e "   ${RED}✗${NC} No stock_prices data found"
    exit 1
fi

echo "   Testing crypto_prices table..."
crypto_count=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM crypto_prices;")
if [ "$crypto_count" -gt 0 ]; then
    echo -e "   ${GREEN}✓${NC} crypto_prices has $crypto_count records"
else
    echo -e "   ${RED}✗${NC} No crypto_prices data found"
    exit 1
fi

# 测试K线图数据查询
echo "6. Testing candlestick data queries..."

echo "   Testing stock prices by symbol..."
echo "   STOCK_A price summary:"
sqlite3 "$DB_FILE" "SELECT symbol, MIN(low_price) as min_price, MAX(high_price) as max_price, AVG(close_price) as avg_close, COUNT(*) as days FROM stock_prices WHERE symbol='STOCK_A' GROUP BY symbol;"

echo "   STOCK_B price summary:"
sqlite3 "$DB_FILE" "SELECT symbol, MIN(low_price) as min_price, MAX(high_price) as max_price, AVG(close_price) as avg_close, COUNT(*) as days FROM stock_prices WHERE symbol='STOCK_B' GROUP BY symbol;"

echo "   STOCK_C price summary:"
sqlite3 "$DB_FILE" "SELECT symbol, MIN(low_price) as min_price, MAX(high_price) as max_price, AVG(close_price) as avg_close, COUNT(*) as days FROM stock_prices WHERE symbol='STOCK_C' GROUP BY symbol;"

echo "   Testing crypto prices by symbol..."
echo "   BTC price summary:"
sqlite3 "$DB_FILE" "SELECT symbol, MIN(low_price) as min_price, MAX(high_price) as max_price, AVG(close_price) as avg_close, COUNT(*) as days FROM crypto_prices WHERE symbol='BTC' GROUP BY symbol;"

echo "   ETH price summary:"
sqlite3 "$DB_FILE" "SELECT symbol, MIN(low_price) as min_price, MAX(high_price) as max_price, AVG(close_price) as avg_close, COUNT(*) as days FROM crypto_prices WHERE symbol='ETH' GROUP BY symbol;"

echo "   Testing recent price trends..."
echo "   Recent STOCK_A prices (last 5 days):"
sqlite3 "$DB_FILE" "SELECT date, open_price, high_price, low_price, close_price, volume FROM stock_prices WHERE symbol='STOCK_A' ORDER BY date DESC LIMIT 5;"

echo "   Recent BTC prices (last 5 days):"
sqlite3 "$DB_FILE" "SELECT date, open_price, high_price, low_price, close_price, volume FROM crypto_prices WHERE symbol='BTC' ORDER BY date DESC LIMIT 5;"

echo "   Testing price volatility..."
echo "   Most volatile stock (by price range):"
sqlite3 "$DB_FILE" "SELECT symbol, AVG(high_price - low_price) as avg_range, AVG((high_price - low_price) / low_price * 100) as avg_volatility_pct FROM stock_prices GROUP BY symbol ORDER BY avg_volatility_pct DESC LIMIT 1;"

echo "   Most volatile crypto (by price range):"
sqlite3 "$DB_FILE" "SELECT symbol, AVG(high_price - low_price) as avg_range, AVG((high_price - low_price) / low_price * 100) as avg_volatility_pct FROM crypto_prices GROUP BY symbol ORDER BY avg_volatility_pct DESC LIMIT 1;"

echo
echo -e "${GREEN}=== Candlestick Data Test Complete ===${NC}"
echo "Sample queries for candlestick charts:"
echo "1. Stock prices: SELECT date, open_price, high_price, low_price, close_price, volume FROM stock_prices WHERE symbol='STOCK_A' ORDER BY date"
echo "2. Crypto prices: SELECT date, open_price, high_price, low_price, close_price, volume FROM crypto_prices WHERE symbol='BTC' ORDER BY date"
echo "3. Multiple stocks: SELECT date, symbol, open_price, high_price, low_price, close_price, volume FROM stock_prices ORDER BY date, symbol"
echo "4. Price comparison: SELECT date, symbol, close_price FROM stock_prices WHERE symbol IN ('STOCK_A', 'STOCK_B') ORDER BY date, symbol"
echo "5. Volume analysis: SELECT date, symbol, volume, (high_price - low_price) as price_range FROM stock_prices ORDER BY volume DESC LIMIT 10"
echo
echo "Candlestick chart data is now ready for testing in Gobi BI Engine!" 