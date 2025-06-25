#!/bin/bash

# 词云图表示例数据插入脚本
# 用于向数据库中插入词云图表所需的示例数据

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
DB_FILE="gobi.db"
SQL_FILE="scripts/data/generate_wordcloud_sample_data.sql"

echo -e "${BLUE}=== Gobi BI 词云图表示例数据插入 ===${NC}"

# 检查数据库文件是否存在
if [ ! -f "$DB_FILE" ]; then
    echo -e "${RED}错误: 数据库文件 $DB_FILE 不存在${NC}"
    echo -e "${YELLOW}请先运行服务器以创建数据库文件${NC}"
    exit 1
fi

# 检查SQL文件是否存在
if [ ! -f "$SQL_FILE" ]; then
    echo -e "${RED}错误: SQL文件 $SQL_FILE 不存在${NC}"
    exit 1
fi

echo -e "${YELLOW}开始插入词云图表示例数据...${NC}"

# 备份原始数据库
echo -e "${YELLOW}备份原始数据库...${NC}"
cp "$DB_FILE" "${DB_FILE}.backup.$(date +%Y%m%d_%H%M%S)"
echo -e "${GREEN}数据库备份完成${NC}"

# 执行SQL文件插入数据
echo -e "${YELLOW}执行SQL文件插入数据...${NC}"
if sqlite3 "$DB_FILE" < "$SQL_FILE"; then
    echo -e "${GREEN}词云图表示例数据插入成功！${NC}"
else
    echo -e "${RED}数据插入失败${NC}"
    echo -e "${YELLOW}正在恢复数据库备份...${NC}"
    cp "${DB_FILE}.backup.$(date +%Y%m%d_%H%M%S)" "$DB_FILE"
    exit 1
fi

# 验证数据插入结果
echo -e "\n${YELLOW}验证数据插入结果...${NC}"

# 检查社交媒体话题数据
SOCIAL_COUNT=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM social_media_topics;")
echo -e "${GREEN}社交媒体话题数据: $SOCIAL_COUNT 条记录${NC}"

# 检查新闻关键词数据
NEWS_COUNT=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM news_keywords;")
echo -e "${GREEN}新闻关键词数据: $NEWS_COUNT 条记录${NC}"

# 检查产品评论关键词数据
REVIEW_COUNT=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM product_review_keywords;")
echo -e "${GREEN}产品评论关键词数据: $REVIEW_COUNT 条记录${NC}"

# 显示数据统计
echo -e "\n${BLUE}=== 数据统计 ===${NC}"
sqlite3 "$DB_FILE" "
SELECT 'social_media_topics' as table_name, COUNT(*) as record_count FROM social_media_topics
UNION ALL
SELECT 'news_keywords' as table_name, COUNT(*) as record_count FROM news_keywords
UNION ALL
SELECT 'product_review_keywords' as table_name, COUNT(*) as record_count FROM product_review_keywords;
"

# 显示热门话题示例
echo -e "\n${BLUE}=== 热门话题示例（前10条）===${NC}"
sqlite3 "$DB_FILE" "
SELECT topic, frequency, category, sentiment 
FROM social_media_topics 
ORDER BY frequency DESC 
LIMIT 10;
"

# 显示新闻关键词示例
echo -e "\n${BLUE}=== 新闻关键词示例（前10条）===${NC}"
sqlite3 "$DB_FILE" "
SELECT keyword, frequency, source, date 
FROM news_keywords 
ORDER BY frequency DESC 
LIMIT 10;
"

# 显示产品评论关键词示例
echo -e "\n${BLUE}=== 产品评论关键词示例（前10条）===${NC}"
sqlite3 "$DB_FILE" "
SELECT keyword, frequency, product_category, sentiment 
FROM product_review_keywords 
ORDER BY frequency DESC 
LIMIT 10;
"

echo -e "\n${GREEN}=== 词云图表示例数据插入完成 ===${NC}"
echo -e "${YELLOW}已插入的数据表:${NC}"
echo "1. social_media_topics - 社交媒体热门话题"
echo "2. news_keywords - 新闻关键词"
echo "3. product_review_keywords - 产品评论关键词"

echo -e "\n${BLUE}数据说明:${NC}"
echo "- social_media_topics: 包含科技、商业、健康、教育、娱乐等类别的热门话题"
echo "- news_keywords: 包含财经、科技、社会等新闻关键词"
echo "- product_review_keywords: 包含电子产品和服装的正面/负面评价关键词"

echo -e "\n${YELLOW}接下来您可以:${NC}"
echo "1. 运行测试脚本: ./scripts/testing/test_wordcloud_sql.sh"
echo "2. 使用API创建词云图表"
echo "3. 查看API示例文档: scripts/docs/api_examples.md" 