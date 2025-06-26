#!/bin/bash

# 甘特图示例数据插入脚本
# 用于向数据库中插入gantt chart所需的示例数据

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

DB_FILE="gobi.db"
SQL_FILE="scripts/data/generate_gantt_sample_data.sql"

echo -e "${BLUE}=== Gobi BI 甘特图示例数据插入 ===${NC}"

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

echo -e "${YELLOW}开始插入甘特图示例数据...${NC}"

# 备份原始数据库
BACKUP_FILE="${DB_FILE}.backup.$(date +%Y%m%d_%H%M%S)"
echo -e "${YELLOW}备份原始数据库...${NC}"
cp "$DB_FILE" "$BACKUP_FILE"
echo -e "${GREEN}数据库备份完成${NC}"

# 执行SQL文件插入数据
if sqlite3 "$DB_FILE" < "$SQL_FILE"; then
    echo -e "${GREEN}甘特图示例数据插入成功！${NC}"
else
    echo -e "${RED}数据插入失败${NC}"
    echo -e "${YELLOW}正在恢复数据库备份...${NC}"
    cp "$BACKUP_FILE" "$DB_FILE"
    exit 1
fi

# 验证数据插入结果
echo -e "\n${YELLOW}验证数据插入结果...${NC}"

# 检查数据条数
ROW_COUNT=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM gantt_demo;")
echo -e "${GREEN}gantt_demo 表数据: $ROW_COUNT 条记录${NC}"

# 显示数据统计
echo -e "\n${BLUE}=== 数据统计 ===${NC}"
sqlite3 "$DB_FILE" "SELECT 'gantt_demo' as table_name, COUNT(*) as record_count FROM gantt_demo;"

# 显示示例数据
echo -e "\n${BLUE}=== 甘特图示例数据（前10条）===${NC}"
sqlite3 "$DB_FILE" "SELECT task_id, task_name, start_date, end_date, duration, progress, status, assignee, project FROM gantt_demo ORDER BY id ASC LIMIT 10;"

echo -e "\n${GREEN}=== 甘特图示例数据插入完成 ===${NC}"
echo -e "${YELLOW}已插入的数据表:${NC}"
echo "1. gantt_demo - 甘特图数据表"

echo -e "\n${BLUE}数据说明:${NC}"
echo "- gantt_demo: 包含软件开发、建筑项目等甘特图场景的任务、时间、进度、状态、负责人等字段"

echo -e "\n${YELLOW}接下来您可以:${NC}"
echo "1. 使用API创建甘特图表（参考 scripts/docs/api_examples.md）"
echo "2. 通过SQL自定义更多甘特图场景数据" 