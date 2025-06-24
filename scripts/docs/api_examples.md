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

*最后更新：2025年6月* 