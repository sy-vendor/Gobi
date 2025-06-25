#!/bin/bash

# 关系图/力导向图示例数据插入脚本
# 用于向数据库中插入graph/network/force-directed图表所需的示例数据

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

DB_FILE="gobi.db"
SQL_FILE="scripts/data/generate_graph_sample_data.sql"

echo -e "${BLUE}=== Gobi BI 关系图/力导向图示例数据插入 ===${NC}"

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

echo -e "${YELLOW}开始插入关系图/力导向图示例数据...${NC}"

# 备份原始数据库
echo -e "${YELLOW}备份原始数据库...${NC}"
cp "$DB_FILE" "${DB_FILE}.backup.$(date +%Y%m%d_%H%M%S)"
echo -e "${GREEN}数据库备份完成${NC}"

# 执行SQL文件插入数据
echo -e "${YELLOW}执行SQL文件插入数据...${NC}"
if sqlite3 "$DB_FILE" < "$SQL_FILE"; then
    echo -e "${GREEN}关系图/力导向图示例数据插入成功！${NC}"
else
    echo -e "${RED}数据插入失败${NC}"
    echo -e "${YELLOW}正在恢复数据库备份...${NC}"
    cp "${DB_FILE}.backup.$(date +%Y%m%d_%H%M%S)" "$DB_FILE"
    exit 1
fi

# 验证数据插入结果
echo -e "\n${YELLOW}验证数据插入结果...${NC}"

# 检查节点数据
NODE_COUNT=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM graph_nodes;")
echo -e "${GREEN}节点数据: $NODE_COUNT 条记录${NC}"

# 检查边数据
EDGE_COUNT=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM graph_edges;")
echo -e "${GREEN}边数据: $EDGE_COUNT 条记录${NC}"

# 显示数据统计
echo -e "\n${BLUE}=== 数据统计 ===${NC}"
sqlite3 "$DB_FILE" "
SELECT 'graph_nodes' as table_name, COUNT(*) as record_count FROM graph_nodes
UNION ALL
SELECT 'graph_edges' as table_name, COUNT(*) as record_count FROM graph_edges;
"

# 显示节点示例
echo -e "\n${BLUE}=== 节点示例（前10条）===${NC}"
sqlite3 "$DB_FILE" "
SELECT id, name, group_id, value, category 
FROM graph_nodes 
ORDER BY id ASC 
LIMIT 10;
"

# 显示边示例
echo -e "\n${BLUE}=== 边示例（前10条）===${NC}"
sqlite3 "$DB_FILE" "
SELECT source, target, weight, relation 
FROM graph_edges 
ORDER BY id ASC 
LIMIT 10;
"

echo -e "\n${GREEN}=== 关系图/力导向图示例数据插入完成 ===${NC}"
echo -e "${YELLOW}已插入的数据表:${NC}"
echo "1. graph_nodes - 关系图节点"
echo "2. graph_edges - 关系图边"

echo -e "\n${BLUE}数据说明:${NC}"
echo "- graph_nodes: 包含用户、设备、组织、地点、网络等多类型节点"
echo "- graph_edges: 包含多种关系类型和权重的边"

echo -e "\n${YELLOW}接下来您可以:${NC}"
echo "1. 使用API创建关系图/力导向图表"
echo "2. 查看API示例文档: scripts/docs/api_examples.md" 