# Gobi - BI Engine MVP

A minimal viable product (MVP) for a Business Intelligence engine built with Go.

# Gobi - BI 引擎最小可行产品

一个使用 Go 语言构建的商业智能引擎最小可行产品。

## Features | 功能特性

- SQL query management and execution | SQL 查询管理和执行
- Interactive chart visualization using go-echarts | 使用 go-echarts 进行交互式图表可视化
- Excel template management and export | Excel 模板管理和导出
- User authentication and authorization | 用户认证和授权
- Data isolation between users | 用户数据隔离
- Dashboard statistics and analytics | 仪表盘统计和分析
- Query execution tracking | 查询执行追踪
- Support for multiple chart types | 支持多种图表类型
  - Bar charts | 柱状图
  - Line charts | 折线图
  - Pie charts | 饼图
  - Scatter plots | 散点图
  - Radar charts | 雷达图
  - Heat maps | 热力图
  - Gauge charts | 仪表盘
  - Funnel charts | 漏斗图
- Data source configuration | 数据源配置
  - Support for multiple database types | 支持多种数据库类型
  - Secure credential storage with AES-256 encryption | 使用 AES-256 加密的安全凭证存储
  - Public/private data source sharing | 公开/私有数据源共享

## Prerequisites | 环境要求

- Go 1.21 or later | Go 1.21 或更高版本
- SQLite (for development) | SQLite（用于开发）
- MySQL/PostgreSQL (optional) | MySQL/PostgreSQL（可选）

## Installation | 安装步骤

1. Clone the repository | 克隆仓库:
```bash
git clone https://github.com/yourusername/gobi.git
cd gobi
```

2. Install dependencies | 安装依赖:
```bash
go mod download
```

3. Set environment variables (optional) | 设置环境变量（可选）:
```bash
export DATABASE_ENCRYPTION_KEY="your-32-byte-encryption-key"
export JWT_SECRET="your-jwt-secret"
```

4. Run the application | 运行应用:
```bash
go run cmd/server/main.go
```

The server will start on port 8080 by default. | 服务器默认在 8080 端口启动。

## Quick Start | 快速开始

```bash
git clone https://github.com/yourusername/gobi.git
cd gobi
go mod download
go run cmd/server/main.go
```

## API Endpoints | API 接口

### Authentication | 认证
- POST /api/auth/register - Register a new user | 注册新用户
- POST /api/auth/login - Login and get JWT token | 登录并获取 JWT 令牌

### Dashboard | 仪表盘
- GET /api/dashboard/stats - Get dashboard statistics | 获取仪表盘统计信息

### Users | 用户管理
- GET /api/users - List all users | 列出所有用户
- PUT /api/users/:id - Update user information | 更新用户信息
- POST /api/users/:id/reset-password - Reset user password | 重置用户密码

### Data Sources | 数据源
- POST /api/datasources - Create a new data source | 创建新数据源
- GET /api/datasources - List all data sources | 列出所有数据源
- GET /api/datasources/:id - Get a specific data source | 获取特定数据源
- PUT /api/datasources/:id - Update a data source | 更新数据源
- DELETE /api/datasources/:id - Delete a data source | 删除数据源

### Queries | 查询
- POST /api/queries - Create a new query | 创建新查询
- GET /api/queries - List all queries | 列出所有查询
- GET /api/queries/:id - Get a specific query | 获取特定查询
- PUT /api/queries/:id - Update a query | 更新查询
- DELETE /api/queries/:id - Delete a query | 删除查询
- POST /api/queries/:id/execute - Execute a query | 执行查询

### Charts | 图表
- POST /api/charts - Create a new chart | 创建新图表
- GET /api/charts - List all charts | 列出所有图表
- GET /api/charts/:id - Get a specific chart | 获取特定图表
- PUT /api/charts/:id - Update a chart | 更新图表
- DELETE /api/charts/:id - Delete a chart | 删除图表

### Excel Templates | Excel 模板
- POST /api/templates - Upload a new template | 上传新模板
- GET /api/templates - List all templates | 列出所有模板
- GET /api/templates/:id - Get a specific template | 获取特定模板
- PUT /api/templates/:id - Update a template | 更新模板
- DELETE /api/templates/:id - Delete a template | 删除模板

## Chart Types and Configuration | 图表类型和配置

### Supported Chart Types | 支持的图表类型

#### 1. Bar Chart (bar) | 柱状图
```json
{
  "type": "bar",
  "config": {
    "xField": "category",
    "yField": "value",
    "seriesField": "series",
    "title": "Sales by Category",
    "legend": true,
    "color": ["#1890ff", "#2fc25b", "#facc14"],
    "tooltip": true,
    "label": true,
    "stack": false
  }
}
```

#### 2. Line Chart (line) | 折线图
```json
{
  "type": "line",
  "config": {
    "xField": "date",
    "yField": "value",
    "seriesField": "series",
    "title": "Trend Analysis",
    "legend": true,
    "color": ["#1890ff", "#2fc25b"],
    "tooltip": true,
    "label": false
  }
}
```

#### 3. Pie Chart (pie) | 饼图
```json
{
  "type": "pie",
  "config": {
    "angleField": "value",
    "colorField": "category",
    "title": "Market Share",
    "legend": true,
    "color": ["#1890ff", "#2fc25b", "#facc14", "#f5222d"],
    "tooltip": true,
    "label": true,
    "radius": 0.8
  }
}
```

#### 4. Scatter Plot (scatter) | 散点图
```json
{
  "type": "scatter",
  "config": {
    "xField": "x",
    "yField": "y",
    "colorField": "category",
    "title": "Correlation Analysis",
    "legend": true,
    "color": ["#1890ff", "#2fc25b", "#facc14"],
    "tooltip": true
  }
}
```

#### 5. Radar Chart (radar) | 雷达图
```json
{
  "type": "radar",
  "config": {
    "angleField": "dimension",
    "valueField": "value",
    "seriesField": "series",
    "title": "Performance Metrics",
    "legend": true,
    "color": ["#1890ff", "#2fc25b"],
    "tooltip": true
  }
}
```

#### 6. Heat Map (heatmap) | 热力图
```json
{
  "type": "heatmap",
  "config": {
    "xField": "x",
    "yField": "y",
    "colorField": "value",
    "title": "Correlation Heatmap",
    "legend": true,
    "color": ["#313695", "#4575b4", "#74add1", "#abd9e9", "#e0f3f8", "#ffffcc", "#fee090", "#fdae61", "#f46d43", "#d73027", "#a50026"],
    "tooltip": true
  }
}
```

#### 7. Gauge Chart (gauge) | 仪表盘
```json
{
  "type": "gauge",
  "config": {
    "valueField": "value",
    "title": "Progress Indicator",
    "color": ["#1890ff", "#2fc25b", "#facc14"],
    "min": 0,
    "max": 100
  }
}
```

#### 8. Funnel Chart (funnel) | 漏斗图
```json
{
  "type": "funnel",
  "config": {
    "angleField": "value",
    "colorField": "stage",
    "title": "Conversion Funnel",
    "legend": true,
    "color": ["#1890ff", "#2fc25b", "#facc14", "#f5222d"],
    "tooltip": true,
    "label": true
  }
}
```

### Config Field Descriptions | 配置字段说明

| Field | Type | Description | 说明 |
|-------|------|-------------|------|
| xField | string | X-axis field name | X轴字段名 |
| yField | string | Y-axis field name | Y轴字段名 |
| seriesField | string | Series grouping field | 系列分组字段 |
| angleField | string | Angle field for pie/funnel | 饼图/漏斗图角度字段 |
| valueField | string | Value field for gauge/radar | 仪表盘/雷达图数值字段 |
| colorField | string | Color grouping field | 颜色分组字段 |
| title | string | Chart title | 图表标题 |
| legend | boolean | Show legend | 显示图例 |
| color | array | Color palette | 颜色配置 |
| tooltip | boolean | Show tooltip | 显示提示框 |
| label | boolean | Show data labels | 显示数据标签 |
| stack | boolean | Stack series | 堆叠系列 |
| radius | number | Pie chart radius (0-1) | 饼图半径 |
| min | number | Gauge minimum value | 仪表盘最小值 |
| max | number | Gauge maximum value | 仪表盘最大值 |

## SQL Examples | SQL 示例

### Sample Database Schema | 示例数据库结构

```sql
-- Sales table
CREATE TABLE sales (
    id INTEGER PRIMARY KEY,
    date DATE,
    category VARCHAR(50),
    product VARCHAR(100),
    amount DECIMAL(10,2),
    region VARCHAR(50),
    salesperson VARCHAR(100)
);

-- Insert sample data
INSERT INTO sales VALUES
(1, '2024-01-01', 'Electronics', 'Laptop', 1200.00, 'North', 'Alice'),
(2, '2024-01-02', 'Electronics', 'Phone', 800.00, 'South', 'Bob'),
(3, '2024-01-03', 'Clothing', 'Shirt', 50.00, 'East', 'Charlie'),
(4, '2024-01-04', 'Clothing', 'Pants', 80.00, 'West', 'David'),
(5, '2024-01-05', 'Electronics', 'Tablet', 600.00, 'North', 'Alice');

-- Funnel table
CREATE TABLE funnel (
    id INTEGER PRIMARY KEY,
    stage VARCHAR(50),
    count INTEGER,
    conversion_rate DECIMAL(5,2)
);

INSERT INTO funnel VALUES
(1, 'Visits', 1000, 100.00),
(2, 'Signups', 200, 20.00),
(3, 'Purchases', 50, 5.00),
(4, 'Repeat', 10, 1.00);
```

### Chart-Specific SQL Examples | 图表专用 SQL 示例

#### Bar Chart SQL | 柱状图 SQL
```sql
-- Sales by category
SELECT category, SUM(amount) as total_sales
FROM sales 
GROUP BY category
ORDER BY total_sales DESC;

-- Sales by region and category
SELECT region, category, SUM(amount) as sales
FROM sales 
GROUP BY region, category
ORDER BY region, sales DESC;
```

#### Line Chart SQL | 折线图 SQL
```sql
-- Sales trend over time
SELECT date, SUM(amount) as daily_sales
FROM sales 
GROUP BY date
ORDER BY date;

-- Sales by category over time
SELECT date, category, SUM(amount) as sales
FROM sales 
GROUP BY date, category
ORDER BY date, category;
```

#### Pie Chart SQL | 饼图 SQL
```sql
-- Market share by category
SELECT category, SUM(amount) as total_sales
FROM sales 
GROUP BY category
ORDER BY total_sales DESC;
```

#### Radar Chart SQL | 雷达图 SQL
```sql
-- Performance metrics by salesperson
SELECT 
    salesperson,
    COUNT(*) as total_sales,
    AVG(amount) as avg_amount,
    SUM(amount) as total_amount,
    COUNT(DISTINCT category) as categories_sold
FROM sales 
GROUP BY salesperson;
```

#### Heat Map SQL | 热力图 SQL
```sql
-- Sales correlation matrix
SELECT 
    s1.category as category1,
    s2.category as category2,
    COUNT(*) as correlation
FROM sales s1
JOIN sales s2 ON s1.region = s2.region
WHERE s1.category != s2.category
GROUP BY s1.category, s2.category;
```

#### Funnel Chart SQL | 漏斗图 SQL
```sql
-- Conversion funnel
SELECT stage, count as value
FROM funnel
ORDER BY id;
```

## API Usage Examples | API 使用示例

### Authentication | 认证

#### Register | 注册
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "test123",
    "email": "test@example.com"
  }'
```

#### Login | 登录
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "test123"
  }'
```

### Dashboard | 仪表盘

#### Get Dashboard Statistics | 获取仪表盘统计
```bash
curl -X GET http://localhost:8080/api/dashboard/stats \
  -H "Authorization: Bearer <your_token>"
```

### Users | 用户管理

#### List Users | 列出用户
```bash
curl -X GET http://localhost:8080/api/users \
  -H "Authorization: Bearer <your_token>"
```

#### Update User | 更新用户
```bash
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{
    "username": "updateduser",
    "email": "updated@example.com"
  }'
```

#### Reset Password | 重置密码
```bash
curl -X POST http://localhost:8080/api/users/1/reset-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{
    "newPassword": "newpassword123"
  }'
```

### Data Sources | 数据源

#### Create Data Source | 创建数据源
```bash
curl -X POST http://localhost:8080/api/datasources \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{
    "name": "MySQL Database",
    "type": "mysql",
    "host": "localhost",
    "port": 3306,
    "database": "test_db",
    "username": "test_user",
    "password": "test_pass",
    "description": "Test database for development",
    "isPublic": true
  }'
```

#### List Data Sources | 列出数据源
```bash
curl -X GET http://localhost:8080/api/datasources \
  -H "Authorization: Bearer <your_token>"
```

#### Get Data Source | 获取数据源
```bash
curl -X GET http://localhost:8080/api/datasources/1 \
  -H "Authorization: Bearer <your_token>"
```

#### Update Data Source | 更新数据源
```bash
curl -X PUT http://localhost:8080/api/datasources/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{
    "name": "Updated MySQL Database",
    "description": "Updated description",
    "isPublic": false
  }'
```

### Queries | 查询

#### Create Query | 创建查询
```bash
curl -X POST http://localhost:8080/api/queries \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{
    "name": "Sales by Category",
    "dataSourceId": 1,
    "sql": "SELECT category, SUM(amount) as total_sales FROM sales GROUP BY category ORDER BY total_sales DESC",
    "description": "Query to analyze sales by product category"
  }'
```

#### List Queries | 列出查询
```bash
curl -X GET http://localhost:8080/api/queries \
  -H "Authorization: Bearer <your_token>"
```

#### Execute Query | 执行查询
```bash
curl -X POST http://localhost:8080/api/queries/1/execute \
  -H "Authorization: Bearer <your_token>"
```

#### Update Query | 更新查询
```bash
curl -X PUT http://localhost:8080/api/queries/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{
    "name": "Updated Sales Query",
    "sql": "SELECT category, SUM(amount) as total_sales FROM sales WHERE date >= '2024-01-01' GROUP BY category",
    "description": "Updated query with date filter"
  }'
```

### Charts | 图表

#### Create Bar Chart | 创建柱状图
```bash
curl -X POST http://localhost:8080/api/charts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{
    "name": "Sales by Category Chart",
    "queryId": 1,
    "type": "bar",
    "config": {
      "xField": "category",
      "yField": "total_sales",
      "title": "Sales by Category",
      "legend": true,
      "color": ["#1890ff", "#2fc25b", "#facc14"],
      "tooltip": true,
      "label": true
    },
    "description": "Bar chart showing sales by product category"
  }'
```

#### Create Line Chart | 创建折线图
```bash
curl -X POST http://localhost:8080/api/charts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{
    "name": "Sales Trend Chart",
    "queryId": 2,
    "type": "line",
    "config": {
      "xField": "date",
      "yField": "daily_sales",
      "title": "Daily Sales Trend",
      "legend": true,
      "color": ["#1890ff"],
      "tooltip": true,
      "label": false
    },
    "description": "Line chart showing daily sales trend"
  }'
```

#### Create Pie Chart | 创建饼图
```bash
curl -X POST http://localhost:8080/api/charts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{
    "name": "Market Share Chart",
    "queryId": 1,
    "type": "pie",
    "config": {
      "angleField": "total_sales",
      "colorField": "category",
      "title": "Market Share by Category",
      "legend": true,
      "color": ["#1890ff", "#2fc25b", "#facc14", "#f5222d"],
      "tooltip": true,
      "label": true,
      "radius": 0.8
    },
    "description": "Pie chart showing market share distribution"
  }'
```

#### List Charts | 列出图表
```bash
curl -X GET http://localhost:8080/api/charts \
  -H "Authorization: Bearer <your_token>"
```

#### Get Chart | 获取图表
```bash
curl -X GET http://localhost:8080/api/charts/1 \
  -H "Authorization: Bearer <your_token>"
```

### Templates | 模板

#### Upload Template | 上传模板
```bash
curl -X POST http://localhost:8080/api/templates \
  -H "Authorization: Bearer <your_token>" \
  -F "file=@/path/to/template.xlsx" \
  -F "name=Sales Report Template" \
  -F "description=Template for monthly sales reports"
```

#### List Templates | 列出模板
```bash
curl -X GET http://localhost:8080/api/templates \
  -H "Authorization: Bearer <your_token>"
```

## Permission Control | 权限控制

### User Roles | 用户角色

1. **Admin User** | 管理员用户
   - Can access all data and operations | 可以访问所有数据和操作
   - Can manage all users | 可以管理所有用户
   - Can view all data sources, queries, and charts | 可以查看所有数据源、查询和图表

2. **Regular User** | 普通用户
   - Can only access own data | 只能访问自己的数据
   - Can access public data sources | 可以访问公开的数据源
   - Cannot modify other users' data | 不能修改其他用户的数据

### Data Isolation | 数据隔离

- Users can only see their own queries and charts by default | 用户默认只能看到自己的查询和图表
- Data sources can be marked as public for sharing | 数据源可以标记为公开以共享
- Admin users bypass all permission restrictions | 管理员用户绕过所有权限限制

## Dashboard Statistics | 仪表盘统计

The `/api/dashboard/stats` endpoint provides comprehensive statistics:

- **Total Queries**: Number of queries in the system | 系统中的查询总数
- **Total Charts**: Number of charts created | 创建的图表总数  
- **Total Users**: Number of registered users | 注册用户总数
- **Today's Queries**: Queries executed today | 今日执行的查询数
- **Query Trends**: Query execution trends over the last 7 days | 最近7天的查询执行趋势
- **Popular Queries**: Most frequently executed queries | 最频繁执行的查询

## Error Handling | 错误处理

所有 API 错误响应均为 JSON 格式，例如：

```json
{
  "code": 400,
  "message": "Invalid registration request",
  "error": "Invalid registration request: Key: 'Password' Error:Field validation for 'Password' failed on the 'required' tag"
}
```

常见错误码：
- 400 Bad Request：参数错误或无效 JSON
- 401 Unauthorized：未认证或 token 无效
- 403 Forbidden：无权限
- 404 Not Found：资源不存在
- 409 Conflict：资源冲突（如用户名已存在）
- 500 Internal Server Error：服务器内部错误

## Security | 安全特性

- All endpoints except registration and login require JWT authentication | 除注册和登录外的所有接口都需要 JWT 认证
- Passwords are hashed using bcrypt | 使用 bcrypt 加密密码
- User data is isolated by default | 默认进行用户数据隔离
- Admin users can access all data | 管理员可以访问所有数据
- Database credentials are encrypted using AES-256 | 数据库凭证使用 AES-256 加密存储
- JWT tokens expire after 24 hours | JWT 令牌24小时后过期

## Testing | 测试

### Automated Testing | 自动化测试

Run the automated test script to verify all functionality:

```bash
chmod +x scripts/test_error_handling.sh
./scripts/test_error_handling.sh
```

The test script covers:
- User registration and authentication | 用户注册和认证
- Data source management | 数据源管理
- Query creation and execution | 查询创建和执行
- Chart creation and configuration | 图表创建和配置
- Permission control | 权限控制
- Error handling | 错误处理

### Manual Testing | 手动测试

Use the provided curl examples in the API Usage section to test individual endpoints.

## Project Structure | 项目结构

```
gobi/
├── cmd/
│   └── server/
│       └── main.go          # Main entry point | 主程序入口
├── config/
│   └── config.go           # Configuration | 配置文件
├── internal/
│   ├── handlers/
│   │   └── handlers.go     # API handlers | API 处理函数
│   ├── middleware/
│   │   ├── auth.go         # Authentication middleware | 认证中间件
│   │   └── error.go        # Error handling middleware | 错误处理中间件
│   └── models/
│       └── models.go       # Data models | 数据模型
├── pkg/
│   ├── database/
│   │   └── database.go     # Database connection | 数据库连接
│   ├── errors/
│   │   └── errors.go       # Error definitions | 错误定义
│   └── utils/
│       └── cache.go        # Cache utilities | 缓存工具
├── scripts/
│   └── test_error_handling.sh  # Automated test script | 自动化测试脚本
├── go.mod                 # Go module file | Go 模块文件
├── go.sum                 # Go module checksums | Go 模块校验和
├── test.sh               # Test runner | 测试运行器
└── README.md             # Project documentation | 项目文档
```

## Development | 开发

### Adding New Chart Types | 添加新图表类型

1. Update the chart type validation in handlers | 在处理器中更新图表类型验证
2. Add corresponding SQL examples | 添加相应的 SQL 示例
3. Update this README with new configuration options | 使用新的配置选项更新此 README

### Database Migrations | 数据库迁移

The application automatically handles database schema changes. For production deployments, consider using a proper migration tool.

### Environment Variables | 环境变量

- `DATABASE_ENCRYPTION_KEY`: 32-byte key for encrypting database credentials
- `JWT_SECRET`: Secret key for JWT token signing
- `PORT`: Server port (default: 8080)
- `DB_PATH`: SQLite database path (default: gobi.db)

## Contributing | 贡献

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Support | 支持

For issues and questions, please create an issue on GitHub or contact the development team.

## Docker & 容器化部署

### 1. 构建镜像
```bash
docker build -t gobi .
```

### 2. 运行容器
```bash
docker run -p 8080:8080 \
  -e DATA_SOURCE_SECRET=12345678901234567890123456789012 \
  -e GOBI_ENV=dev \
  gobi
```

### 3. 使用 docker-compose 一键启动（推荐开发/测试）
```bash
docker-compose up --build
```
- 默认会启动 gobi 服务和 MySQL 数据库
- 可在 `docker-compose.yml` 中自定义数据库密码、端口等

### 4. 挂载配置和迁移目录
- Dockerfile 会自动拷贝 `config/`、`migrations/`、`.env`、`config.yaml` 到容器内
- 如需自定义配置，建议挂载本地目录或修改镜像后重建

### 5. 环境变量说明
- `DATA_SOURCE_SECRET`：数据源加密密钥，必须32字节
- `GOBI_ENV`：配置环境（dev/prod等）

### 6. 访问服务
- 默认端口：`http://localhost:8080`
- 可通过 `docker-compose.yml` 或 `-p` 参数自定义端口映射
