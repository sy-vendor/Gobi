# Gobi BI API 示例文档

本文档提供了Gobi BI系统中所有图表类型的完整API请求实例和返回示例。

## 📋 目录

- [认证](#认证)
- [面积图 (Area Chart)](#面积图-area-chart)
- [柱状图 (Bar Chart)](#柱状图-bar-chart)
- [折线图 (Line Chart)](#折线图-line-chart)
- [饼图 (Pie Chart)](#饼图-pie-chart)
- [散点图 (Scatter Chart)](#散点图-scatter-chart)
- [3D柱状图 (3D Bar Chart)](#3d柱状图-3d-bar-chart)
- [3D散点图 (3D Scatter Chart)](#3d散点图-3d-scatter-chart)
- [3D表面图 (3D Surface Chart)](#3d表面图-3d-surface-chart)
- [3D气泡图 (3D Bubble Chart)](#3d气泡图-3d-bubble-chart)
- [矩形树状图（TreeMap）](#矩形树状图treemap)
- [旭日图 (Sunburst)](#旭日图-sunburst)
- [树形图 (Tree Diagram)](#树形图-tree-diagram)
- [箱线图 (Box Plot)](#箱线图-box-plot)
- [K线图/蜡烛图 (Candlestick Chart)](#k线图蜡烛图-candlestick-chart)
- [词云图 (Word Cloud)](#词云图-word-cloud)
- [关系图/力导向图 (Graph/Network/Force-directed)](#关系图力导向图-graphnetworkforce-directed)
- [瀑布图 (Waterfall Chart)](#瀑布图-waterfall-chart)

## 🔐 认证

### 用户登录获取JWT Token

**请求**:
```bash
curl -X POST "http://localhost:8080/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

**返回**:
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "ID": 1,
      "username": "admin",
      "email": "admin@gobi.com",
      "role": "admin"
    }
  },
  "message": "Login successful"
}
```

### API Key认证

**请求**:
```bash
curl -X GET "http://localhost:8080/api/queries" \
  -H "Authorization: ApiKey YOUR_API_KEY"
```

---

## 📊 面积图 (Area Chart)

### 1. 创建数据源

**请求**:
```bash
curl -X POST "http://localhost:8080/api/datasources" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "销售数据源",
    "type": "sqlite",
    "database": "gobi.db",
    "description": "包含销售趋势数据的SQLite数据源"
  }'
```

**返回**:
```json
{
  "success": true,
  "data": {
    "ID": 1,
    "name": "销售数据源",
    "type": "sqlite",
    "database": "gobi.db",
    "description": "包含销售趋势数据的SQLite数据源",
    "userID": 1,
    "createdAt": "2025-06-24T11:50:00Z",
    "updatedAt": "2025-06-24T11:50:00Z"
  },
  "message": "Data source created successfully"
}
```

### 2. 创建查询

**请求**:
```bash
curl -X POST "http://localhost:8080/api/queries" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "销售趋势查询",
    "dataSourceId": 1,
    "sql": "SELECT month, product_category, sales_amount FROM sales_trend ORDER BY month, product_category",
    "description": "查询不同产品类别的月度销售数据"
  }'
```

**返回**:
```json
{
  "success": true,
  "data": {
    "ID": 1,
    "name": "销售趋势查询",
    "dataSourceId": 1,
    "sql": "SELECT month, product_category, sales_amount FROM sales_trend ORDER BY month, product_category",
    "description": "查询不同产品类别的月度销售数据",
    "userID": 1,
    "createdAt": "2025-06-24T11:50:00Z",
    "updatedAt": "2025-06-24T11:50:00Z"
  },
  "message": "Query created successfully"
}
```

### 3. 创建面积图

**请求**:
```bash
curl -X POST "http://localhost:8080/api/charts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "销售趋势面积图",
    "queryId": 1,
    "type": "area",
    "config": "{
      \"xField\": \"month\",
      \"yField\": \"sales_amount\",
      \"seriesField\": \"product_category\",
      \"title\": \"月度销售趋势\",
      \"legend\": true,
      \"color\": [\"#1890ff\", \"#2fc25b\", \"#facc14\"],
      \"tooltip\": true,
      \"smooth\": true,
      \"fillOpacity\": 0.6
    }",
    "description": "展示不同产品类别的月度销售趋势"
  }'
```

**返回**:
```json
{
  "success": true,
  "data": {
    "ID": 1,
    "name": "销售趋势面积图",
    "queryId": 1,
    "type": "area",
    "config": "{\"xField\":\"month\",\"yField\":\"sales_amount\",\"seriesField\":\"product_category\",\"title\":\"月度销售趋势\",\"legend\":true,\"color\":[\"#1890ff\",\"#2fc25b\",\"#facc14\"],\"tooltip\":true,\"smooth\":true,\"fillOpacity\":0.6}",
    "description": "展示不同产品类别的月度销售趋势",
    "userID": 1,
    "createdAt": "2025-06-24T11:50:00Z",
    "updatedAt": "2025-06-24T11:50:00Z"
  },
  "message": "Chart created successfully"
}
```

### 4. 获取面积图数据

**请求**:
```bash
curl -X GET "http://localhost:8080/api/charts/1/data" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**返回**:
```json
{
  "success": true,
  "data": {
    "chart": {
      "ID": 1,
      "name": "销售趋势面积图",
      "type": "area",
      "config": {
        "xField": "month",
        "yField": "sales_amount",
        "seriesField": "product_category",
        "title": "月度销售趋势",
        "legend": true,
        "color": ["#1890ff", "#2fc25b", "#facc14"],
        "tooltip": true,
        "smooth": true,
        "fillOpacity": 0.6
      }
    },
    "data": [
      {
        "month": "2024-01",
        "product_category": "Electronics",
        "sales_amount": 125000
      },
      {
        "month": "2024-01",
        "product_category": "Clothing",
        "sales_amount": 89000
      },
      {
        "month": "2024-02",
        "product_category": "Electronics",
        "sales_amount": 138000
      },
      {
        "month": "2024-02",
        "product_category": "Clothing",
        "sales_amount": 92000
      }
    ]
  },
  "message": "Chart data retrieved successfully"
}
```

---

## 📈 柱状图 (Bar Chart)

### 1. 创建数据源
```bash
curl -X POST "http://localhost:8080/api/datasources" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "销售数据源",
    "type": "sqlite",
    "database": "gobi.db",
    "description": "包含销售数据的SQLite数据源"
  }'
```

### 2. 创建查询
```bash
curl -X POST "http://localhost:8080/api/queries" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "月度销售查询",
    "dataSourceId": 1,
    "sql": "SELECT month, product_category, SUM(sales_amount) as total_sales FROM sales_trend GROUP BY month, product_category ORDER BY month, product_category",
    "description": "查询各产品类别的月度销售总额"
  }'
```

### 3. 创建柱状图

**请求**:
```bash
curl -X POST "http://localhost:8080/api/charts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "月度销售柱状图",
    "queryId": 2,
    "type": "bar",
    "config": "{
      \"xField\": \"month\",
      \"yField\": \"total_sales\",
      \"seriesField\": \"product_category\",
      \"title\": \"月度销售对比\",
      \"legend\": true,
      \"color\": [\"#1890ff\", \"#2fc25b\", \"#facc14\"],
      \"tooltip\": true,
      \"barWidth\": 20,
      \"barGap\": 0.1
    }",
    "description": "展示不同产品类别的月度销售对比"
  }'
```

**返回**:
```json
{
  "success": true,
  "data": {
    "ID": 2,
    "name": "月度销售柱状图",
    "queryId": 2,
    "type": "bar",
    "config": "{\"xField\":\"month\",\"yField\":\"total_sales\",\"seriesField\":\"product_category\",\"title\":\"月度销售对比\",\"legend\":true,\"color\":[\"#1890ff\",\"#2fc25b\",\"#facc14\"],\"tooltip\":true,\"barWidth\":20,\"barGap\":0.1}",
    "description": "展示不同产品类别的月度销售对比",
    "userID": 1,
    "createdAt": "2025-06-24T11:50:00Z",
    "updatedAt": "2025-06-24T11:50:00Z"
  },
  "message": "Chart created successfully"
}
```

---

## 📉 折线图 (Line Chart)

### 1. 创建数据源
```bash
curl -X POST "http://localhost:8080/api/datasources" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "用户数据源",
    "type": "sqlite",
    "database": "gobi.db",
    "description": "包含用户增长数据的SQLite数据源"
  }'
```

### 2. 创建查询
```bash
curl -X POST "http://localhost:8080/api/queries" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "用户增长查询",
    "dataSourceId": 2,
    "sql": "SELECT date, user_type, new_users FROM user_growth ORDER BY date, user_type",
    "description": "查询不同类型用户的增长数据"
  }'
```

### 3. 创建折线图

**请求**:
```bash
curl -X POST "http://localhost:8080/api/charts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "用户增长折线图",
    "queryId": 3,
    "type": "line",
    "config": "{
      \"xField\": \"date\",
      \"yField\": \"new_users\",
      \"seriesField\": \"user_type\",
      \"title\": \"用户增长趋势\",
      \"legend\": true,
      \"color\": [\"#722ed1\", \"#13c2c2\", \"#eb2f96\"],
      \"tooltip\": true,
      \"smooth\": true,
      \"pointSize\": 4
    }",
    "description": "展示不同类型用户的增长趋势"
  }'
```

**返回**:
```json
{
  "success": true,
  "data": {
    "ID": 3,
    "name": "用户增长折线图",
    "queryId": 3,
    "type": "line",
    "config": "{\"xField\":\"date\",\"yField\":\"new_users\",\"seriesField\":\"user_type\",\"title\":\"用户增长趋势\",\"legend\":true,\"color\":[\"#722ed1\",\"#13c2c2\",\"#eb2f96\"],\"tooltip\":true,\"smooth\":true,\"pointSize\":4}",
    "description": "展示不同类型用户的增长趋势",
    "userID": 1,
    "createdAt": "2025-06-24T11:50:00Z",
    "updatedAt": "2025-06-24T11:50:00Z"
  },
  "message": "Chart created successfully"
}
```

---

## 🥧 饼图 (Pie Chart)

### 1. 创建数据源
```bash
curl -X POST "http://localhost:8080/api/datasources" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "市场数据源",
    "type": "sqlite",
    "database": "gobi.db",
    "description": "包含市场份额数据的SQLite数据源"
  }'
```

### 2. 创建查询
```bash
curl -X POST "http://localhost:8080/api/queries" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "市场份额查询",
    "dataSourceId": 3,
    "sql": "SELECT company, AVG(market_share) as avg_market_share FROM market_share GROUP BY company ORDER BY avg_market_share DESC",
    "description": "查询各公司的平均市场份额"
  }'
```

### 3. 创建饼图

**请求**:
```bash
curl -X POST "http://localhost:8080/api/charts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "市场份额饼图",
    "queryId": 4,
    "type": "pie",
    "config": "{
      \"angleField\": \"avg_market_share\",
      \"colorField\": \"company\",
      \"title\": \"市场份额分布\",
      \"legend\": true,
      \"color\": [\"#fa541c\", \"#a0d911\", \"#2f54eb\", \"#722ed1\"],
      \"tooltip\": true,
      \"radius\": 0.8,
      \"innerRadius\": 0.4
    }",
    "description": "展示各公司市场份额分布"
  }'
```

**返回**:
```json
{
  "success": true,
  "data": {
    "ID": 4,
    "name": "市场份额饼图",
    "queryId": 4,
    "type": "pie",
    "config": "{\"angleField\":\"avg_market_share\",\"colorField\":\"company\",\"title\":\"市场份额分布\",\"legend\":true,\"color\":[\"#fa541c\",\"#a0d911\",\"#2f54eb\",\"#722ed1\"],\"tooltip\":true,\"radius\":0.8,\"innerRadius\":0.4}",
    "description": "展示各公司市场份额分布",
    "userID": 1,
    "createdAt": "2025-06-24T11:50:00Z",
    "updatedAt": "2025-06-24T11:50:00Z"
  },
  "message": "Chart created successfully"
}
```

---

## 🔵 散点图 (Scatter Chart)

### 1. 创建数据源
```bash
curl -X POST "http://localhost:8080/api/datasources" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "产品数据源",
    "type": "sqlite",
    "database": "gobi.db",
    "description": "包含产品性能数据的SQLite数据源"
  }'
```

### 2. 创建查询
```bash
curl -X POST "http://localhost:8080/api/queries" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "产品性能查询",
    "dataSourceId": 4,
    "sql": "SELECT price, performance_score, product_category, sales_volume FROM products_3d WHERE performance_score IS NOT NULL AND price IS NOT NULL",
    "description": "查询产品价格、性能和销售数据"
  }'
```

### 3. 创建散点图

**请求**:
```bash
curl -X POST "http://localhost:8080/api/charts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "产品性能散点图",
    "queryId": 5,
    "type": "scatter",
    "config": "{
      \"xField\": \"price\",
      \"yField\": \"performance_score\",
      \"colorField\": \"product_category\",
      \"sizeField\": \"sales_volume\",
      \"title\": \"产品价格与性能关系\",
      \"legend\": true,
      \"color\": [\"#1890ff\", \"#2fc25b\", \"#facc14\"],
      \"tooltip\": true,
      \"pointSize\": 8
    }",
    "description": "展示产品价格与性能的关系"
  }'
```

**返回**:
```json
{
  "success": true,
  "data": {
    "ID": 5,
    "name": "产品性能散点图",
    "queryId": 5,
    "type": "scatter",
    "config": "{\"xField\":\"price\",\"yField\":\"performance_score\",\"colorField\":\"product_category\",\"sizeField\":\"sales_volume\",\"title\":\"产品价格与性能关系\",\"legend\":true,\"color\":[\"#1890ff\",\"#2fc25b\",\"#facc14\"],\"tooltip\":true,\"pointSize\":8}",
    "description": "展示产品价格与性能的关系",
    "userID": 1,
    "createdAt": "2025-06-24T11:50:00Z",
    "updatedAt": "2025-06-24T11:50:00Z"
  },
  "message": "Chart created successfully"
}
```

---

## 🏗️ 3D柱状图 (3D Bar Chart)

### 1. 创建数据源
```bash
curl -X POST "http://localhost:8080/api/datasources" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "3D销售数据源",
    "type": "sqlite",
    "database": "gobi.db",
    "description": "包含3D销售数据的SQLite数据源"
  }'
```

### 2. 创建查询
```bash
curl -X POST "http://localhost:8080/api/queries" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "3D销售数据查询",
    "dataSourceId": 5,
    "sql": "SELECT category as x, region as y, SUM(amount) as z FROM sales_3d GROUP BY category, region ORDER BY category, region",
    "description": "查询3D销售数据，按类别和地区分组"
  }'
```

### 3. 创建3D柱状图

**请求**:
```bash
curl -X POST "http://localhost:8080/api/charts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "3D销售柱状图",
    "queryId": 6,
    "type": "3d-bar",
    "config": "{
      \"xField\": \"x\",
      \"yField\": \"y\",
      \"zField\": \"z\",
      \"title\": \"3D销售数据\",
      \"legend\": true,
      \"color\": [\"#1890ff\", \"#2fc25b\", \"#facc14\"],
      \"tooltip\": true,
      \"grid3D\": {
        \"boxWidth\": 100,
        \"boxHeight\": 100,
        \"boxDepth\": 100,
        \"viewControl\": {
          \"alpha\": 20,
          \"beta\": 40,
          \"distance\": 200
        }
      }
    }",
    "description": "3D展示销售数据"
  }'
```

**返回**:
```json
{
  "success": true,
  "data": {
    "ID": 6,
    "name": "3D销售柱状图",
    "queryId": 6,
    "type": "3d-bar",
    "config": "{\"xField\":\"x\",\"yField\":\"y\",\"zField\":\"z\",\"title\":\"3D销售数据\",\"legend\":true,\"color\":[\"#1890ff\",\"#2fc25b\",\"#facc14\"],\"tooltip\":true,\"grid3D\":{\"boxWidth\":100,\"boxHeight\":100,\"boxDepth\":100,\"viewControl\":{\"alpha\":20,\"beta\":40,\"distance\":200}}}",
    "description": "3D展示销售数据",
    "userID": 1,
    "createdAt": "2025-06-24T11:50:00Z",
    "updatedAt": "2025-06-24T11:50:00Z"
  },
  "message": "Chart created successfully"
}
```

---

## 🌐 3D散点图 (3D Scatter Chart)

### 1. 创建数据源
```bash
curl -X POST "http://localhost:8080/api/datasources" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "3D产品数据源",
    "type": "sqlite",
    "database": "gobi.db",
    "description": "包含3D产品数据的SQLite数据源"
  }'
```

### 2. 创建查询
```bash
curl -X POST "http://localhost:8080/api/queries" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "3D产品数据查询",
    "dataSourceId": 6,
    "sql": "SELECT performance_score as x, price as y, customer_rating as z, product_category as category, sales_volume as size FROM products_3d WHERE performance_score IS NOT NULL AND price IS NOT NULL AND customer_rating IS NOT NULL",
    "description": "查询3D产品性能、价格和评分数据"
  }'
```

### 3. 创建3D散点图

**请求**:
```bash
curl -X POST "http://localhost:8080/api/charts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "3D产品散点图",
    "queryId": 7,
    "type": "3d-scatter",
    "config": "{
      \"xField\": \"x\",
      \"yField\": \"y\",
      \"zField\": \"z\",
      \"colorField\": \"category\",
      \"sizeField\": \"size\",
      \"title\": \"3D产品性能分析\",
      \"legend\": true,
      \"color\": [\"#1890ff\", \"#2fc25b\", \"#facc14\"],
      \"tooltip\": true,
      \"symbolSize\": 10,
      \"grid3D\": {
        \"boxWidth\": 100,
        \"boxHeight\": 100,
        \"boxDepth\": 100,
        \"viewControl\": {
          \"alpha\": 20,
          \"beta\": 40,
          \"distance\": 200
        }
      }
    }",
    "description": "3D展示产品性能、价格和评分关系"
  }'
```

**返回**:
```json
{
  "success": true,
  "data": {
    "ID": 7,
    "name": "3D产品散点图",
    "queryId": 7,
    "type": "3d-scatter",
    "config": "{\"xField\":\"x\",\"yField\":\"y\",\"zField\":\"z\",\"colorField\":\"category\",\"sizeField\":\"size\",\"title\":\"3D产品性能分析\",\"legend\":true,\"color\":[\"#1890ff\",\"#2fc25b\",\"#facc14\"],\"tooltip\":true,\"symbolSize\":10,\"grid3D\":{\"boxWidth\":100,\"boxHeight\":100,\"boxDepth\":100,\"viewControl\":{\"alpha\":20,\"beta\":40,\"distance\":200}}}",
    "description": "3D展示产品性能、价格和评分关系",
    "userID": 1,
    "createdAt": "2025-06-24T11:50:00Z",
    "updatedAt": "2025-06-24T11:50:00Z"
  },
  "message": "Chart created successfully"
}
```

---

## 🏔️ 3D表面图 (3D Surface Chart)

### 1. 创建数据源
```bash
curl -X POST "http://localhost:8080/api/datasources" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "地形数据源",
    "type": "sqlite",
    "database": "gobi.db",
    "description": "包含地形数据的SQLite数据源"
  }'
```

### 2. 创建查询
```bash
curl -X POST "http://localhost:8080/api/queries" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "地形数据查询",
    "dataSourceId": 7,
    "sql": "SELECT longitude as x, latitude as y, elevation as z FROM terrain_3d ORDER BY longitude, latitude",
    "description": "查询地形经纬度和海拔数据"
  }'
```

### 3. 创建3D表面图

**请求**:
```bash
curl -X POST "http://localhost:8080/api/charts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "3D地形表面图",
    "queryId": 8,
    "type": "3d-surface",
    "config": "{
      \"xField\": \"x\",
      \"yField\": \"y\",
      \"zField\": \"z\",
      \"title\": \"3D地形图\",
      \"color\": [\"#313695\", \"#4575b4\", \"#74add1\", \"#abd9e9\", \"#e0f3f8\", \"#ffffcc\", \"#fee090\", \"#fdae61\", \"#f46d43\", \"#d73027\", \"#a50026\"],
      \"tooltip\": true,
      \"shading\": \"realistic\",
      \"grid3D\": {
        \"boxWidth\": 100,
        \"boxHeight\": 100,
        \"boxDepth\": 100,
        \"viewControl\": {
          \"alpha\": 20,
          \"beta\": 40,
          \"distance\": 200
        }
      }
    }",
    "description": "3D地形表面展示"
  }'
```

**返回**:
```json
{
  "success": true,
  "data": {
    "ID": 8,
    "name": "3D地形表面图",
    "queryId": 8,
    "type": "3d-surface",
    "config": "{\"xField\":\"x\",\"yField\":\"y\",\"zField\":\"z\",\"title\":\"3D地形图\",\"color\":[\"#313695\",\"#4575b4\",\"#74add1\",\"#abd9e9\",\"#e0f3f8\",\"#ffffcc\",\"#fee090\",\"#fdae61\",\"#f46d43\",\"#d73027\",\"#a50026\"],\"tooltip\":true,\"shading\":\"realistic\",\"grid3D\":{\"boxWidth\":100,\"boxHeight\":100,\"boxDepth\":100,\"viewControl\":{\"alpha\":20,\"beta\":40,\"distance\":200}}}",
    "description": "3D地形表面展示",
    "userID": 1,
    "createdAt": "2025-06-24T11:50:00Z",
    "updatedAt": "2025-06-24T11:50:00Z"
  },
  "message": "Chart created successfully"
}
```

---

## 🫧 3D气泡图 (3D Bubble Chart)

### 1. 创建数据源
```bash
curl -X POST "http://localhost:8080/api/datasources" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "城市数据源",
    "type": "sqlite",
    "database": "gobi.db",
    "description": "包含城市数据的SQLite数据源"
  }'
```

### 2. 创建查询
```bash
curl -X POST "http://localhost:8080/api/queries" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "城市数据查询",
    "dataSourceId": 8,
    "sql": "SELECT gdp as x, population as y, area as z, city_name as category, population as size FROM cities_3d ORDER BY gdp DESC",
    "description": "查询城市GDP、人口和面积数据"
  }'
```

### 3. 创建3D气泡图

**请求**:
```bash
curl -X POST "http://localhost:8080/api/charts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "3D城市气泡图",
    "queryId": 9,
    "type": "3d-bubble",
    "config": "{
      \"xField\": \"x\",
      \"yField\": \"y\",
      \"zField\": \"z\",
      \"sizeField\": \"size\",
      \"colorField\": \"category\",
      \"title\": \"3D城市数据\",
      \"legend\": true,
      \"color\": [\"#1890ff\", \"#2fc25b\", \"#facc14\", \"#f5222d\"],
      \"tooltip\": true,
      \"grid3D\": {
        \"boxWidth\": 100,
        \"boxHeight\": 100,
        \"boxDepth\": 100,
        \"viewControl\": {
          \"alpha\": 20,
          \"beta\": 40,
          \"distance\": 200
        }
      }
    }",
    "description": "3D展示城市GDP、人口和面积关系"
  }'
```

**返回**:
```json
{
  "success": true,
  "data": {
    "ID": 9,
    "name": "3D城市气泡图",
    "queryId": 9,
    "type": "3d-bubble",
    "config": "{\"xField\":\"x\",\"yField\":\"y\",\"zField\":\"z\",\"sizeField\":\"size\",\"colorField\":\"category\",\"title\":\"3D城市数据\",\"legend\":true,\"color\":[\"#1890ff\",\"#2fc25b\",\"#facc14\",\"#f5222d\"],\"tooltip\":true,\"grid3D\":{\"boxWidth\":100,\"boxHeight\":100,\"boxDepth\":100,\"viewControl\":{\"alpha\":20,\"beta\":40,\"distance\":200}}}",
    "description": "3D展示城市GDP、人口和面积关系",
    "userID": 1,
    "createdAt": "2025-06-24T11:50:00Z",
    "updatedAt": "2025-06-24T11:50:00Z"
  },
  "message": "Chart created successfully"
}
```

---

## 🟦 矩形树状图（TreeMap）

### 1. 创建数据源
```bash
curl -X POST "http://localhost:8080/api/datasources" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "层级数据源",
    "type": "sqlite",
    "database": "gobi.db",
    "description": "包含层级结构数据的SQLite数据源"
  }'
```

### 2. 创建查询
```bash
curl -X POST "http://localhost:8080/api/queries" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "部门层级查询",
    "dataSourceId": 10,
    "sql": "SELECT id, parent_id, name, value, category FROM department_hierarchy",
    "description": "查询公司部门层级结构"
  }'
```

### 3. 创建树状图

**请求**:
```bash
curl -X POST "http://localhost:8080/api/charts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "公司部门树状图",
    "queryId": 11,
    "type": "treemap",
    "config": "{\"dataField\":\"name\",\"valueField\":\"value\",\"colorField\":\"category\",\"title\":\"公司部门分布\",\"legend\":true,\"tooltip\":true}",
    "description": "展示公司各部门及其层级结构"
  }'
```

**返回**:
```json
{
  "success": true,
  "data": {
    "ID": 11,
    "name": "公司部门树状图",
    "queryId": 11,
    "type": "treemap",
    "config": "{\"dataField\":\"name\",\"valueField\":\"value\",\"colorField\":\"category\",\"title\":\"公司部门分布\",\"legend\":true,\"tooltip\":true}",
    "description": "展示公司各部门及其层级结构",
    "userID": 1,
    "createdAt": "2025-06-24T11:50:00Z",
    "updatedAt": "2025-06-24T11:50:00Z"
  },
  "message": "Chart created successfully"
}
```

---

## ☀️ 旭日图 (Sunburst)

### 1. 创建数据源
```bash
curl -X POST "http://localhost:8080/api/datasources" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "产品层级数据源",
    "type": "sqlite",
    "database": "gobi.db",
    "description": "包含产品层级数据的SQLite数据源"
  }'
```

### 2. 创建查询
```bash
curl -X POST "http://localhost:8080/api/queries" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "产品分类层级查询",
    "dataSourceId": 12,
    "sql": "SELECT id, parent_id, name, value, category FROM product_hierarchy",
    "description": "查询产品分类层级结构"
  }'
```

### 3. 创建旭日图

**请求**:
```bash
curl -X POST "http://localhost:8080/api/charts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "产品分类旭日图",
    "queryId": 13,
    "type": "sunburst",
    "config": "{\"dataField\":\"name\",\"valueField\":\"value\",\"colorField\":\"category\",\"title\":\"产品分类层级分布\",\"legend\":true,\"tooltip\":true}",
    "description": "展示产品分类的层级结构"
  }'
```

**返回**:
```json
{
  "success": true,
  "data": {
    "ID": 13,
    "name": "产品分类旭日图",
    "queryId": 13,
    "type": "sunburst",
    "config": "{\"dataField\":\"name\",\"valueField\":\"value\",\"colorField\":\"category\",\"title\":\"产品分类层级分布\",\"legend\":true,\"tooltip\":true}",
    "description": "展示产品分类的层级结构",
    "userID": 1,
    "createdAt": "2025-06-24T11:50:00Z",
    "updatedAt": "2025-06-24T11:50:00Z"
  },
  "message": "Chart created successfully"
}
```

---

## 🌲 树形图 (Tree Diagram)

### 1. 创建数据源
```bash
curl -X POST "http://localhost:8080/api/datasources" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "组织架构数据源",
    "type": "sqlite",
    "database": "gobi.db",
    "description": "包含组织架构树的SQLite数据源"
  }'
```

### 2. 创建查询
```bash
curl -X POST "http://localhost:8080/api/queries" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "组织架构树查询",
    "dataSourceId": 20,
    "sql": "SELECT id, parent_id, name, position FROM org_tree",
    "description": "查询公司组织架构树"
  }'
```

### 3. 创建树形图

**请求**:
```bash
curl -X POST "http://localhost:8080/api/charts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "公司组织架构树形图",
    "queryId": 21,
    "type": "tree",
    "config": "{\"idField\":\"id\",\"parentField\":\"parent_id\",\"nameField\":\"name\",\"valueField\":\"position\",\"title\":\"公司组织架构\",\"legend\":true,\"tooltip\":true}",
    "description": "展示公司组织架构的分支结构树"
  }'
```

**返回**:
```json
{
  "success": true,
  "data": {
    "ID": 21,
    "name": "公司组织架构树形图",
    "queryId": 21,
    "type": "tree",
    "config": "{\"idField\":\"id\",\"parentField\":\"parent_id\",\"nameField\":\"name\",\"valueField\":\"position\",\"title\":\"公司组织架构\",\"legend\":true,\"tooltip\":true}",
    "description": "展示公司组织架构的分支结构树",
    "userID": 1,
    "createdAt": "2025-06-24T11:50:00Z",
    "updatedAt": "2025-06-24T11:50:00Z"
  },
  "message": "Chart created successfully"
}
```

---

## 📦 箱线图 (Box Plot)

### 1. 创建数据源
```bash
curl -X POST "http://localhost:8080/api/datasources" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "成绩数据源",
    "type": "sqlite",
    "database": "gobi.db",
    "description": "包含学生成绩分布数据的SQLite数据源"
  }'
```

### 2. 创建查询
```bash
curl -X POST "http://localhost:8080/api/queries" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "成绩分布查询",
    "dataSourceId": 30,
    "sql": "SELECT class, subject, score FROM student_scores ORDER BY class, subject",
    "description": "查询不同班级不同科目的成绩分布"
  }'
```

### 3. 创建箱线图

**请求**:
```bash
curl -X POST "http://localhost:8080/api/charts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "班级成绩箱线图",
    "queryId": 30,
    "type": "boxplot",
    "config": "{
      \"xField\": \"class\",
      \"yField\": \"score\",
      \"seriesField\": \"subject\",
      \"title\": \"各班级各科目成绩分布\",
      \"legend\": true,
      \"color\": [\"#1890ff\", \"#2fc25b\", \"#facc14\"],
      \"tooltip\": true,
      \"boxStyle\": {
        \"stroke\": \"#545454\",
        \"fill\": \"#f6f6f6\"
      },
      \"outlierStyle\": {
        \"fill\": \"#f5222d\",
        \"stroke\": \"#f5222d\"
      }
    }",
    "description": "展示不同班级各科目成绩的分布情况"
  }'
```

**返回**:
```json
{
  "success": true,
  "data": {
    "ID": 30,
    "name": "班级成绩箱线图",
    "queryId": 30,
    "type": "boxplot",
    "config": "{\"xField\":\"class\",\"yField\":\"score\",\"seriesField\":\"subject\",\"title\":\"各班级各科目成绩分布\",\"legend\":true,\"color\":[\"#1890ff\",\"#2fc25b\",\"#facc14\"],\"tooltip\":true,\"boxStyle\":{\"stroke\":\"#545454\",\"fill\":\"#f6f6f6\"},\"outlierStyle\":{\"fill\":\"#f5222d\",\"stroke\":\"#f5222d\"}}",
    "description": "展示不同班级各科目成绩的分布情况",
    "userID": 1,
    "createdAt": "2025-06-24T11:50:00Z",
    "updatedAt": "2025-06-24T11:50:00Z"
  },
  "message": "Chart created successfully"
}
```

### 4. 获取箱线图数据

**请求**:
```bash
curl -X GET "http://localhost:8080/api/charts/30/data" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**返回**:
```json
{
  "success": true,
  "data": {
    "chart": {
      "ID": 30,
      "name": "班级成绩箱线图",
      "type": "boxplot",
      "config": {
        "xField": "class",
        "yField": "score",
        "seriesField": "subject",
        "title": "各班级各科目成绩分布",
        "legend": true,
        "color": ["#1890ff", "#2fc25b", "#facc14"],
        "tooltip": true,
        "boxStyle": {
          "stroke": "#545454",
          "fill": "#f6f6f6"
        },
        "outlierStyle": {
          "fill": "#f5222d",
          "stroke": "#f5222d"
        }
      }
    },
    "data": [
      {
        "class": "A班",
        "subject": "数学",
        "score": 85.5
      },
      {
        "class": "A班",
        "subject": "数学",
        "score": 92.3
      },
      {
        "class": "A班",
        "subject": "英语",
        "score": 88.2
      },
      {
        "class": "B班",
        "subject": "数学",
        "score": 72.8
      },
      {
        "class": "B班",
        "subject": "英语",
        "score": 82.5
      }
    ]
  },
  "message": "Chart data retrieved successfully"
}
```

---

## 📦 产品性能箱线图 (Product Performance Box Plot)

### 1. 创建数据源
```bash
curl -X POST "http://localhost:8080/api/datasources" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "产品性能数据源",
    "type": "sqlite",
    "database": "gobi.db",
    "description": "包含产品性能测试数据的SQLite数据源"
  }'
```

### 2. 创建查询
```bash
curl -X POST "http://localhost:8080/api/queries" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "产品性能查询",
    "dataSourceId": 31,
    "sql": "SELECT product_type, test_metric, value FROM product_performance ORDER BY product_type, test_metric",
    "description": "查询不同产品类型的性能测试数据"
  }'
```

### 3. 创建产品性能箱线图

**请求**:
```bash
curl -X POST "http://localhost:8080/api/charts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "产品性能箱线图",
    "queryId": 31,
    "type": "boxplot",
    "config": "{
      \"xField\": \"product_type\",
      \"yField\": \"value\",
      \"seriesField\": \"test_metric\",
      \"title\": \"产品性能测试分布\",
      \"legend\": true,
      \"color\": [\"#722ed1\", \"#13c2c2\"],
      \"tooltip\": true,
      \"boxStyle\": {
        \"stroke\": \"#545454\",
        \"fill\": \"#f6f6f6\"
      },
      \"outlierStyle\": {
        \"fill\": \"#f5222d\",
        \"stroke\": \"#f5222d\"
      }
    }",
    "description": "展示不同产品类型的性能测试分布"
  }'
```

**返回**:
```json
{
  "success": true,
  "data": {
    "ID": 31,
    "name": "产品性能箱线图",
    "queryId": 31,
    "type": "boxplot",
    "config": "{\"xField\":\"product_type\",\"yField\":\"value\",\"seriesField\":\"test_metric\",\"title\":\"产品性能测试分布\",\"legend\":true,\"color\":[\"#722ed1\",\"#13c2c2\"],\"tooltip\":true,\"boxStyle\":{\"stroke\":\"#545454\",\"fill\":\"#f6f6f6\"},\"outlierStyle\":{\"fill\":\"#f5222d\",\"stroke\":\"#f5222d\"}}",
    "description": "展示不同产品类型的性能测试分布",
    "userID": 1,
    "createdAt": "2025-06-24T11:50:00Z",
    "updatedAt": "2025-06-24T11:50:00Z"
  },
  "message": "Chart created successfully"
}
```

---

## 🔧 通用操作

### 更新图表

**请求**:
```bash
curl -X PUT "http://localhost:8080/api/charts/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "更新后的图表名称",
    "config": "{\"xField\":\"month\",\"yField\":\"sales_amount\",\"title\":\"更新后的标题\"}",
    "description": "更新后的描述"
  }'
```

### 删除图表

**请求**:
```bash
curl -X DELETE "http://localhost:8080/api/charts/1" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 获取所有图表

**请求**:
```bash
curl -X GET "http://localhost:8080/api/charts" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 按类型获取图表

**请求**:
```bash
curl -X GET "http://localhost:8080/api/charts?type=area" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## 📝 配置参数说明

### 通用配置
- `title`: 图表标题
- `legend`: 是否显示图例（true/false）
- `color`: 颜色数组
- `tooltip`: 是否显示提示框（true/false）

### 2D图表配置
- `xField`: X轴字段名
- `yField`: Y轴字段名
- `seriesField`: 系列字段名
- `angleField`: 角度字段名（饼图）
- `colorField`: 颜色字段名
- `sizeField`: 大小字段名

### 3D图表配置
- `xField`: X轴字段名
- `yField`: Y轴字段名
- `zField`: Z轴字段名
- `colorField`: 颜色字段名
- `sizeField`: 大小字段名
- `grid3D`: 3D网格配置
  - `boxWidth`: 盒子宽度
  - `boxHeight`: 盒子高度
  - `boxDepth`: 盒子深度
  - `viewControl`: 视角控制
    - `alpha`: 水平旋转角度
    - `beta`: 垂直旋转角度
    - `distance`: 距离

### 样式配置
- `smooth`: 是否平滑曲线（true/false）
- `fillOpacity`: 填充透明度（0-1）
- `barWidth`: 柱状图宽度
- `barGap`: 柱状图间距
- `pointSize`: 点大小
- `radius`: 饼图半径
- `innerRadius`: 饼图内半径
- `symbolSize`: 符号大小
- `shading`: 3D表面着色方式

---

## 🚀 快速开始

1. **启动服务器**:
   ```bash
   go run cmd/server/main.go
   ```

2. **登录获取Token**:
   ```bash
   curl -X POST "http://localhost:8080/api/auth/login" \
     -H "Content-Type: application/json" \
     -d '{"username": "admin", "password": "admin123"}'
   ```

3. **创建数据源**:
   ```bash
   curl -X POST "http://localhost:8080/api/datasources" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -d '{
       "name": "测试数据源",
       "type": "sqlite",
       "database": "gobi.db"
     }'
   ```

4. **创建查询**:
   ```bash
   curl -X POST "http://localhost:8080/api/queries" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -d '{
       "name": "销售数据查询",
       "dataSourceId": 1,
       "sql": "SELECT month, product_category, sales_amount FROM sales_trend"
     }'
   ```

5. **创建图表**:
   使用上面的示例创建你需要的图表类型。

---

## 📈 K线图/蜡烛图 (Candlestick Chart)

### 1. 创建数据源
```bash
curl -X POST "http://localhost:8080/api/datasources" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "股票数据源",
    "type": "sqlite",
    "database": "gobi.db",
    "description": "包含股票价格数据的SQLite数据源"
  }'
```

### 2. 创建查询
```bash
curl -X POST "http://localhost:8080/api/queries" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "股票价格查询",
    "dataSourceId": 40,
    "sql": "SELECT date, open_price, high_price, low_price, close_price, volume FROM stock_prices WHERE symbol='STOCK_A' ORDER BY date",
    "description": "查询股票A的价格数据"
  }'
```

### 3. 创建K线图

**请求**:
```bash
curl -X POST "http://localhost:8080/api/charts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "股票A K线图",
    "queryId": 40,
    "type": "candlestick",
    "config": "{
      \"xField\": \"date\",
      \"openField\": \"open_price\",
      \"highField\": \"high_price\",
      \"lowField\": \"low_price\",
      \"closeField\": \"close_price\",
      \"volumeField\": \"volume\",
      \"title\": \"股票A价格走势\",
      \"legend\": true,
      \"color\": [\"#f5222d\", \"#52c41a\"],
      \"tooltip\": true,
      \"candlestickStyle\": {
        \"stroke\": \"#000000\",
        \"lineWidth\": 1
      },
      \"volumeStyle\": {
        \"fill\": \"#1890ff\",
        \"opacity\": 0.6
      }
    }",
    "description": "展示股票A的价格走势和成交量"
  }'
```

**返回**:
```json
{
  "success": true,
  "data": {
    "ID": 40,
    "name": "股票A K线图",
    "queryId": 40,
    "type": "candlestick",
    "config": "{\"xField\":\"date\",\"openField\":\"open_price\",\"highField\":\"high_price\",\"lowField\":\"low_price\",\"closeField\":\"close_price\",\"volumeField\":\"volume\",\"title\":\"股票A价格走势\",\"legend\":true,\"color\":[\"#f5222d\",\"#52c41a\"],\"tooltip\":true,\"candlestickStyle\":{\"stroke\":\"#000000\",\"lineWidth\":1},\"volumeStyle\":{\"fill\":\"#1890ff\",\"opacity\":0.6}}",
    "description": "展示股票A的价格走势和成交量",
    "userID": 1,
    "createdAt": "2025-06-24T11:50:00Z",
    "updatedAt": "2025-06-24T11:50:00Z"
  },
  "message": "Chart created successfully"
}
```

### 4. 获取K线图数据

**请求**:
```bash
curl -X GET "http://localhost:8080/api/charts/40/data" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**返回**:
```json
{
  "success": true,
  "data": {
    "chart": {
      "ID": 40,
      "name": "股票A K线图",
      "type": "candlestick",
      "config": {
        "xField": "date",
        "openField": "open_price",
        "highField": "high_price",
        "lowField": "low_price",
        "closeField": "close_price",
        "volumeField": "volume",
        "title": "股票A价格走势",
        "legend": true,
        "color": ["#f5222d", "#52c41a"],
        "tooltip": true,
        "candlestickStyle": {
          "stroke": "#000000",
          "lineWidth": 1
        },
        "volumeStyle": {
          "fill": "#1890ff",
          "opacity": 0.6
        }
      }
    },
    "data": [
      {
        "date": "2024-01-02",
        "open_price": 100.50,
        "high_price": 102.30,
        "low_price": 99.80,
        "close_price": 101.20,
        "volume": 1500000
      },
      {
        "date": "2024-01-03",
        "open_price": 101.20,
        "high_price": 103.50,
        "low_price": 100.90,
        "close_price": 102.80,
        "volume": 1800000
      },
      {
        "date": "2024-01-04",
        "open_price": 102.80,
        "high_price": 104.20,
        "low_price": 101.50,
        "close_price": 103.90,
        "volume": 2200000
      }
    ]
  },
  "message": "Chart data retrieved successfully"
}
```

---

## 📈 加密货币K线图 (Crypto Candlestick Chart)

### 1. 创建数据源
```bash
curl -X POST "http://localhost:8080/api/datasources" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "加密货币数据源",
    "type": "sqlite",
    "database": "gobi.db",
    "description": "包含加密货币价格数据的SQLite数据源"
  }'
```

### 2. 创建查询
```bash
curl -X POST "http://localhost:8080/api/queries" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "比特币价格查询",
    "dataSourceId": 41,
    "sql": "SELECT date, open_price, high_price, low_price, close_price, volume FROM crypto_prices WHERE symbol='BTC' ORDER BY date",
    "description": "查询比特币的价格数据"
  }'
```

### 3. 创建加密货币K线图

**请求**:
```bash
curl -X POST "http://localhost:8080/api/charts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "比特币K线图",
    "queryId": 41,
    "type": "candlestick",
    "config": "{
      \"xField\": \"date\",
      \"openField\": \"open_price\",
      \"highField\": \"high_price\",
      \"lowField\": \"low_price\",
      \"closeField\": \"close_price\",
      \"volumeField\": \"volume\",
      \"title\": \"比特币价格走势\",
      \"legend\": true,
      \"color\": [\"#f5222d\", \"#52c41a\"],
      \"tooltip\": true,
      \"candlestickStyle\": {
        \"stroke\": \"#000000\",
        \"lineWidth\": 1
      },
      \"volumeStyle\": {
        \"fill\": \"#722ed1\",
        \"opacity\": 0.6
      }
    }",
    "description": "展示比特币的价格走势和成交量"
  }'
```

**返回**:
```json
{
  "success": true,
  "data": {
    "ID": 41,
    "name": "比特币K线图",
    "queryId": 41,
    "type": "candlestick",
    "config": "{\"xField\":\"date\",\"openField\":\"open_price\",\"highField\":\"high_price\",\"lowField\":\"low_price\",\"closeField\":\"close_price\",\"volumeField\":\"volume\",\"title\":\"比特币价格走势\",\"legend\":true,\"color\":[\"#f5222d\",\"#52c41a\"],\"tooltip\":true,\"candlestickStyle\":{\"stroke\":\"#000000\",\"lineWidth\":1},\"volumeStyle\":{\"fill\":\"#722ed1\",\"opacity\":0.6}}",
    "description": "展示比特币的价格走势和成交量",
    "userID": 1,
    "createdAt": "2025-06-24T11:50:00Z",
    "updatedAt": "2025-06-24T11:50:00Z"
  },
  "message": "Chart created successfully"
}
```

---

## 📈 多股票对比K线图 (Multi-Stock Candlestick Chart)

### 1. 创建查询
```bash
curl -X POST "http://localhost:8080/api/queries" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "多股票价格查询",
    "dataSourceId": 40,
    "sql": "SELECT date, symbol, open_price, high_price, low_price, close_price, volume FROM stock_prices WHERE symbol IN ('STOCK_A', 'STOCK_B') ORDER BY date, symbol",
    "description": "查询多只股票的价格数据"
  }'
```

### 3. 创建多股票对比K线图

**请求**:
```bash
curl -X POST "http://localhost:8080/api/charts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "多股票对比K线图",
    "queryId": 42,
    "type": "candlestick",
    "config": "{
      \"xField\": \"date\",
      \"openField\": \"open_price\",
      \"highField\": \"high_price\",
      \"lowField\": \"low_price\",
      \"closeField\": \"close_price\",
      \"volumeField\": \"volume\",
      \"seriesField\": \"symbol\",
      \"title\": \"多股票价格对比\",
      \"legend\": true,
      \"color\": [\"#f5222d\", \"#52c41a\", \"#1890ff\"],
      \"tooltip\": true,
      \"candlestickStyle\": {
        \"stroke\": \"#000000\",
        \"lineWidth\": 1
      },
      \"volumeStyle\": {
        \"fill\": \"#722ed1\",
        \"opacity\": 0.6
      }
    }",
    "description": "展示多只股票的价格对比"
  }'
```

**返回**:
```json
{
  "success": true,
  "data": {
    "ID": 42,
    "name": "多股票对比K线图",
    "queryId": 42,
    "type": "candlestick",
    "config": "{\"xField\":\"date\",\"openField\":\"open_price\",\"highField\":\"high_price\",\"lowField\":\"low_price\",\"closeField\":\"close_price\",\"volumeField\":\"volume\",\"seriesField\":\"symbol\",\"title\":\"多股票价格对比\",\"legend\":true,\"color\":[\"#f5222d\",\"#52c41a\",\"#1890ff\"],\"tooltip\":true,\"candlestickStyle\":{\"stroke\":\"#000000\",\"lineWidth\":1},\"volumeStyle\":{\"fill\":\"#722ed1\",\"opacity\":0.6}}",
    "description": "展示多只股票的价格对比",
    "userID": 1,
    "createdAt": "2025-06-24T11:50:00Z",
    "updatedAt": "2025-06-24T11:50:00Z"
  },
  "message": "Chart created successfully"
}
```

---

## ☁️ 词云图 (Word Cloud)

### 1. 创建数据源

**请求**:
```bash
curl -X POST "http://localhost:8080/api/datasources" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "词云数据源",
    "type": "sqlite",
    "database": "gobi.db",
    "description": "包含社交媒体话题、新闻关键词和产品评论关键词的SQLite数据源"
  }'
```

**返回**:
```json
{
  "success": true,
  "data": {
    "ID": 50,
    "name": "词云数据源",
    "type": "sqlite",
    "database": "gobi.db",
    "description": "包含社交媒体话题、新闻关键词和产品评论关键词的SQLite数据源",
    "userID": 1,
    "createdAt": "2025-06-24T11:50:00Z",
    "updatedAt": "2025-06-24T11:50:00Z"
  },
  "message": "Data source created successfully"
}
```

### 2. 创建社交媒体话题查询

**请求**:
```bash
curl -X POST "http://localhost:8080/api/queries" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "社交媒体热门话题",
    "dataSourceId": 50,
    "sql": "SELECT topic as word, frequency as value, category, sentiment FROM social_media_topics ORDER BY frequency DESC LIMIT 30",
    "description": "查询社交媒体热门话题及其频率"
  }'
```

**返回**:
```json
{
  "success": true,
  "data": {
    "ID": 50,
    "name": "社交媒体热门话题",
    "dataSourceId": 50,
    "sql": "SELECT topic as word, frequency as value, category, sentiment FROM social_media_topics ORDER BY frequency DESC LIMIT 30",
    "description": "查询社交媒体热门话题及其频率",
    "userID": 1,
    "createdAt": "2025-06-24T11:50:00Z",
    "updatedAt": "2025-06-24T11:50:00Z"
  },
  "message": "Query created successfully"
}
```

### 3. 创建社交媒体话题词云

**请求**:
```bash
curl -X POST "http://localhost:8080/api/charts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "社交媒体热门话题词云",
    "queryId": 50,
    "type": "wordcloud",
    "config": "{
      \"wordField\": \"word\",
      \"weightField\": \"value\",
      \"colorField\": \"category\",
      \"title\": \"社交媒体热门话题词云\",
      \"subtitle\": \"基于话题频率和分类\",
      \"color\": [\"#1890ff\", \"#2fc25b\", \"#facc14\", \"#f5222d\", \"#722ed1\"],
      \"fontSize\": [12, 60],
      \"rotation\": [-90, 90],
      \"spiral\": \"archimedean\",
      \"shape\": \"circle\",
      \"tooltip\": true,
      \"legend\": true
    }",
    "description": "展示社交媒体热门话题的词云图，字体大小表示话题热度"
  }'
```

**返回**:
```json
{
  "success": true,
  "data": {
    "ID": 50,
    "name": "社交媒体热门话题词云",
    "queryId": 50,
    "type": "wordcloud",
    "config": "{\"wordField\":\"word\",\"weightField\":\"value\",\"colorField\":\"category\",\"title\":\"社交媒体热门话题词云\",\"subtitle\":\"基于话题频率和分类\",\"color\":[\"#1890ff\",\"#2fc25b\",\"#facc14\",\"#f5222d\",\"#722ed1\"],\"fontSize\":[12,60],\"rotation\":[-90,90],\"spiral\":\"archimedean\",\"shape\":\"circle\",\"tooltip\":true,\"legend\":true}",
    "description": "展示社交媒体热门话题的词云图，字体大小表示话题热度",
    "userID": 1,
    "createdAt": "2025-06-24T11:50:00Z",
    "updatedAt": "2025-06-24T11:50:00Z"
  },
  "message": "Chart created successfully"
}
```

### 4. 创建新闻关键词词云

**请求**:
```bash
curl -X POST "http://localhost:8080/api/queries" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "新闻关键词统计",
    "dataSourceId": 50,
    "sql": "SELECT keyword as word, frequency as value, source, date FROM news_keywords ORDER BY frequency DESC LIMIT 25",
    "description": "查询新闻关键词及其出现频率"
  }'
```

**创建新闻关键词词云**:
```bash
curl -X POST "http://localhost:8080/api/charts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "新闻关键词词云",
    "queryId": 51,
    "type": "wordcloud",
    "config": "{
      \"wordField\": \"word\",
      \"weightField\": \"value\",
      \"colorField\": \"source\",
      \"title\": \"新闻关键词词云\",
      \"subtitle\": \"基于关键词出现频率\",
      \"color\": [\"#1890ff\", \"#2fc25b\", \"#facc14\"],
      \"fontSize\": [14, 50],
      \"rotation\": [-45, 45],
      \"spiral\": \"rectangular\",
      \"shape\": \"diamond\",
      \"tooltip\": true,
      \"legend\": true
    }",
    "description": "展示新闻关键词的词云图，字体大小表示关键词重要性"
  }'
```

### 5. 创建产品评论词云

**请求**:
```bash
curl -X POST "http://localhost:8080/api/queries" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "产品正面评价关键词",
    "dataSourceId": 50,
    "sql": "SELECT keyword as word, frequency as value, product_category FROM product_review_keywords WHERE sentiment = \"positive\" ORDER BY frequency DESC LIMIT 20",
    "description": "查询产品正面评价关键词"
  }'
```

**创建产品正面评价词云**:
```bash
curl -X POST "http://localhost:8080/api/charts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "产品正面评价词云",
    "queryId": 52,
    "type": "wordcloud",
    "config": "{
      \"wordField\": \"word\",
      \"weightField\": \"value\",
      \"colorField\": \"product_category\",
      \"title\": \"产品正面评价词云\",
      \"subtitle\": \"基于评价关键词频率\",
      \"color\": [\"#52c41a\", \"#1890ff\", \"#722ed1\"],
      \"fontSize\": [16, 48],
      \"rotation\": [0, 0],
      \"spiral\": \"archimedean\",
      \"shape\": \"circle\",
      \"tooltip\": true,
      \"legend\": true
    }",
    "description": "展示产品正面评价关键词的词云图"
  }'
```

### 6. 获取词云图数据

**请求**:
```bash
curl -X GET "http://localhost:8080/api/charts/50/data" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**返回**:
```json
{
  "success": true,
  "data": {
    "chart": {
      "ID": 50,
      "name": "社交媒体热门话题词云",
      "type": "wordcloud",
      "config": {
        "wordField": "word",
        "weightField": "value",
        "colorField": "category",
        "title": "社交媒体热门话题词云",
        "subtitle": "基于话题频率和分类",
        "color": ["#1890ff", "#2fc25b", "#facc14", "#f5222d", "#722ed1"],
        "fontSize": [12, 60],
        "rotation": [-90, 90],
        "spiral": "archimedean",
        "shape": "circle",
        "tooltip": true,
        "legend": true
      }
    },
    "data": [
      {
        "word": "人工智能",
        "value": 1250,
        "category": "科技",
        "sentiment": "positive"
      },
      {
        "word": "大数据",
        "value": 1100,
        "category": "科技",
        "sentiment": "positive"
      },
      {
        "word": "心理健康",
        "value": 1200,
        "category": "健康",
        "sentiment": "positive"
      },
      {
        "word": "在线教育",
        "value": 1100,
        "category": "教育",
        "sentiment": "positive"
      },
      {
        "word": "数字化转型",
        "value": 890,
        "category": "商业",
        "sentiment": "positive"
      }
    ]
  },
  "message": "Chart data retrieved successfully"
}
```

---

## 📝 词云图配置参数说明

### 基本配置
- `wordField`: 词语字段名（必填）
- `weightField`: 权重字段名，决定字体大小（必填）
- `colorField`: 颜色字段名，用于分类着色
- `title`: 图表标题
- `subtitle`: 图表副标题

### 样式配置
- `fontSize`: 字体大小范围 [最小值, 最大值]
- `rotation`: 旋转角度范围 [最小值, 最大值]
- `spiral`: 螺旋排列方式
  - `archimedean`: 阿基米德螺旋（圆形）
  - `rectangular`: 矩形螺旋
- `shape`: 词云形状
  - `circle`: 圆形
  - `diamond`: 菱形
  - `triangle`: 三角形
  - `star`: 星形
- `color`: 颜色数组，用于不同分类的着色

### 交互配置
- `tooltip`: 是否显示提示框（true/false）
- `legend`: 是否显示图例（true/false）

### 数据格式要求
词云图数据需要包含以下字段：
- 词语字段：包含要显示的词语
- 权重字段：数值类型，决定字体大小
- 颜色字段：可选，用于分类着色

*最后更新：2025年6月* 

---

## 🔗 关系图/力导向图 (Graph/Network/Force-directed)

### 1. 创建数据源

**请求**:
```bash
curl -X POST "http://localhost:8080/api/datasources" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "关系图数据源",
    "type": "sqlite",
    "database": "gobi.db",
    "description": "包含节点和边的关系图数据"
  }'
```

### 2. 创建节点查询

**请求**:
```bash
curl -X POST "http://localhost:8080/api/queries" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Graph Nodes",
    "dataSourceId": 1,
    "sql": "SELECT id, name, group_id, value, category FROM graph_nodes",
    "description": "关系图节点数据"
  }'
```

### 3. 创建边查询

**请求**:
```bash
curl -X POST "http://localhost:8080/api/queries" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Graph Edges",
    "dataSourceId": 1,
    "sql": "SELECT source, target, weight, relation FROM graph_edges",
    "description": "关系图边数据"
  }'
```

### 4. 创建关系图表

**请求**:
```bash
curl -X POST "http://localhost:8080/api/charts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "关系图示例",
    "queryId": 10,
    "type": "graph",
    "config": "{\"nodesQueryId\":10,\"edgesQueryId\":11,\"sourceField\":\"source\",\"targetField\":\"target\",\"nodeIdField\":\"id\",\"nodeNameField\":\"name\",\"groupField\":\"group_id\",\"categoryField\":\"category\",\"valueField\":\"value\",\"relationField\":\"relation\",\"weightField\":\"weight\",\"title\":\"关系图示例\",\"legend\":true,\"color\":[\"#1890ff\",\"#2fc25b\",\"#facc14\",\"#f5222d\"],\"tooltip\":true,\"repulsion\":200,\"gravity\":0.1,\"edgeSymbol\":[\"circle\",\"arrow\"],\"layout\":\"force\"}",
    "description": "展示节点和边的关系图/力导向图"
  }'
```

**配置参数说明**:
- `nodesQueryId`: 节点查询ID
- `edgesQueryId`: 边查询ID
- `sourceField`/`targetField`: 边的起止字段
- `nodeIdField`/`nodeNameField`: 节点ID/名称字段
- `groupField`/`categoryField`: 分组/分类字段
- `valueField`: 节点权重
- `relationField`: 边关系类型
- `weightField`: 边权重
- `title`/`legend`/`color`/`tooltip`/`repulsion`/`gravity`/`edgeSymbol`/`layout` 等

**数据格式要求**:
- 节点表需包含 id、name、group_id、value、category 等字段
- 边表需包含 source、target、weight、relation 等字段

---

## 🪜 瀑布图 (Waterfall Chart)

### 1. 创建数据源
```bash
curl -X POST "http://localhost:8080/api/datasources" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "财务数据源",
    "type": "sqlite",
    "database": "gobi.db",
    "description": "包含利润拆解数据的SQLite数据源"
  }'
```

### 2. 创建查询
```bash
curl -X POST "http://localhost:8080/api/queries" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "利润拆解查询",
    "dataSourceId": 1,
    "sql": "SELECT step, amount, type, description FROM waterfall_demo ORDER BY id",
    "description": "查询利润拆解的各步骤数据"
  }'
```

### 3. 创建瀑布图
```bash
curl -X POST "http://localhost:8080/api/charts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "年度利润瀑布图",
    "queryId": 1,
    "type": "waterfall",
    "config": "{\n      \"xField\": \"step\",\n      \"yField\": \"amount\",\n      \"typeField\": \"type\",\n      \"descriptionField\": \"description\",\n      \"title\": \"年度利润拆解\",\n      \"legend\": true,\n      \"color\": [\"#1890ff\", \"#f5222d\", \"#2fc25b\"],\n      \"tooltip\": true\n    }",
    "description": "展示年度利润的各项增减变化"
  }'
```

**返回**:
```json
{
  "success": true,
  "data": {
    "ID": 1,
    "name": "年度利润瀑布图",
    "queryId": 1,
    "type": "waterfall",
    "config": "{\"xField\":\"step\",\"yField\":\"amount\",\"typeField\":\"type\",\"descriptionField\":\"description\",\"title\":\"年度利润拆解\",\"legend\":true,\"color\":[\"#1890ff\",\"#f5222d\",\"#2fc25b\"],\"tooltip\":true}",
    "description": "展示年度利润的各项增减变化",
    "userID": 1,
    "createdAt": "2025-06-24T11:50:00Z",
    "updatedAt": "2025-06-24T11:50:00Z"
  },
  "message": "Chart created successfully"
}
```

### 4. 获取瀑布图数据
```bash
curl -X GET "http://localhost:8080/api/charts/1/data" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**返回**:
```json
{
  "success": true,
  "data": {
    "chart": {
      "ID": 1,
      "name": "年度利润瀑布图",
      "type": "waterfall",
      "config": {
        "xField": "step",
        "yField": "amount",
        "typeField": "type",
        "descriptionField": "description",
        "title": "年度利润拆解",
        "legend": true,
        "color": ["#1890ff", "#f5222d", "#2fc25b"],
        "tooltip": true
      }
    },
    "data": [
      { "step": "期初余额", "amount": 1000, "type": "base", "description": "年初资金" },
      { "step": "主营业务收入", "amount": 2000, "type": "increase", "description": "主营业务带来的收入" },
      { "step": "其他收入", "amount": 500, "type": "increase", "description": "其他来源收入" },
      { "step": "运营成本", "amount": -1200, "type": "decrease", "description": "日常运营支出" },
      { "step": "税费", "amount": -300, "type": "decrease", "description": "税收及附加" },
      { "step": "净利润", "amount": 2000, "type": "base", "description": "年末净利润" }
    ]
  },
  "message": "Chart data retrieved successfully"
}
```