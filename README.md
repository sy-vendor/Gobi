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
- Support for multiple chart types | 支持多种图表类型
  - Bar charts | 柱状图
  - Line charts | 折线图
  - Pie charts | 饼图
  - Scatter plots | 散点图
  - Radar charts | 雷达图
  - Heat maps | 热力图
  - Gauge charts | 仪表盘
  - Funnel charts | 漏斗图

## Prerequisites | 环境要求

- Go 1.21 or later | Go 1.21 或更高版本
- SQLite (for development) | SQLite（用于开发）

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

3. Run the application | 运行应用:
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

访问 [http://localhost:8080](http://localhost:8080) 查看服务是否启动成功。

## API Endpoints | API 接口

### Authentication | 认证
- POST /api/auth/register - Register a new user | 注册新用户
- POST /api/auth/login - Login and get JWT token | 登录并获取 JWT 令牌

### Queries | 查询
- POST /api/queries - Create a new query | 创建新查询
- GET /api/queries - List all queries | 列出所有查询
- GET /api/queries/:id - Get a specific query | 获取特定查询
- PUT /api/queries/:id - Update a query | 更新查询
- DELETE /api/queries/:id - Delete a query | 删除查询

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

## Example API Usage | API 使用示例

### Register | 注册
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"test123"}'
```

### Login | 登录
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"test123"}'
```

### Create Query | 创建查询
```bash
curl -X POST http://localhost:8080/api/queries \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{"name":"Test Query","sql":"SELECT * FROM users"}'
```

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
│   │   └── auth.go         # Authentication middleware | 认证中间件
│   └── models/
│       └── models.go       # Data models | 数据模型
├── pkg/
│   ├── database/
│   │   └── database.go     # Database connection | 数据库连接
│   └── utils/              # Utility functions | 工具函数
├── web/
│   ├── static/            # Static assets | 静态资源
│   └── templates/         # HTML templates | HTML 模板
├── go.mod                 # Go module file | Go 模块文件
└── README.md             # Project documentation | 项目文档
```

## Testing | 测试

可使用 `scripts/test_error_handling.sh` 脚本自动化测试主要接口和错误处理：

```bash
chmod +x scripts/test_error_handling.sh
./scripts/test_error_handling.sh
```
