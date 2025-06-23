# Gobi - BI Engine MVP

A minimal viable product (MVP) for a Business Intelligence engine built with Go.

[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg)](https://go.dev/dl/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://github.com/sy-vendor/gobi/actions/workflows/go.yml/badge.svg)](https://github.com/sy-vendor/gobi/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sy-vendor/gobi)](https://goreportcard.com/report/github.com/sy-vendor/gobi)
[![GitHub stars](https://img.shields.io/github/stars/sy-vendor/gobi)](https://github.com/sy-vendor/gobi/stargazers)

# Gobi - BI 引擎最小可行产品

一个使用 Go 语言构建的商业智能引擎最小可行产品。

## Features | 功能特性

- SQL query management and execution | SQL 查询管理和执行
- Interactive chart visualization | 交互式图表可视化
- Excel template management and export | Excel 模板管理和导出
- User authentication and authorization | 用户认证和授权
- Data isolation between users | 用户数据隔离
- Dashboard statistics and analytics | 仪表盘统计和分析
- **Scheduled Report Generation | 定时报告生成**
- **Enhanced JWT Configuration | 增强的JWT配置**
- **Improved Error Handling | 改进的错误处理**

## Prerequisites | 环境要求

- Go 1.21 or later | Go 1.21 或更高版本
- SQLite (for development) | SQLite（用于开发）
- MySQL/PostgreSQL (optional) | MySQL/PostgreSQL（可选）

## Quick Start | 快速开始

```bash
git clone https://github.com/sy-vendor/gobi.git
cd gobi
go mod download
go run cmd/server/main.go
```

The server will start on port 8080 by default. | 服务器默认在 8080 端口启动。

## Configuration | 配置

### Configuration File | 配置文件

The application uses `config/config.yaml` for configuration management. | 应用程序使用 `config/config.yaml` 进行配置管理。

```yaml
default:
  server:
    port: "8080"
  jwt:
    secret: "default_jwt_secret"
    expiration_hours: 168  # 7天
  database:
    type: "sqlite"
    dsn: "gobi.db"
```

### JWT Configuration | JWT配置

- `jwt.secret`: JWT签名密钥
- `jwt.expiration_hours`: Token过期时间（小时）
  - 168小时 = 7天
  - 720小时 = 30天
  - 2160小时 = 90天

## API Endpoints | API 接口

### Authentication | 认证
- POST /api/auth/register - Register a new user | 注册新用户
- POST /api/auth/login - Login and get JWT token | 登录并获取 JWT 令牌

### Dashboard | 仪表盘
- GET /api/dashboard/stats - Get dashboard statistics | 获取仪表盘统计信息

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
- GET /api/templates/:id/download - Download a template | 下载模板

### Report Schedules | 定时报告
- POST /api/reports/schedules - Create a new report schedule | 创建新的定时报告
- GET /api/reports/schedules - List all report schedules | 列出所有定时报告
- GET /api/reports/schedules/:id - Get a specific report schedule | 获取特定定时报告
- PUT /api/reports/schedules/:id - Update a report schedule | 更新定时报告
- DELETE /api/reports/schedules/:id - Delete a report schedule | 删除定时报告

### Reports | 报告
- GET /api/reports - List all generated reports | 列出所有生成的报告
- GET /api/reports/:id/download - Download a specific report | 下载特定报告

## Chart Types | 图表类型

Supported chart types: | 支持的图表类型：
- Bar charts | 柱状图
- Line charts | 折线图
- Pie charts | 饼图
- Scatter plots | 散点图
- Radar charts | 雷达图
- Heat maps | 热力图
- Gauge charts | 仪表盘
- Funnel charts | 漏斗图
- **3D Bar charts | 3D柱状图**
- **3D Scatter plots | 3D散点图**
- **3D Surface charts | 3D曲面图**
- **3D Bubble charts | 3D气泡图**

## Cron Expression Guide | Cron表达式指南

### Basic Format | 基本格式
```
* * * * *
│ │ │ │ │
│ │ │ │ └── 星期几 (0-7)
│ │ │ └──── 月份 (1-12)
│ │ └────── 日期 (1-31)
│ └──────── 小时 (0-23)
└────────── 分钟 (0-59)
```

### Common Examples | 常见示例
- `0 9 * * *` - 每天上午9点
- `0 0 * * 1` - 每周一午夜
- `35 16 * * *` - 每天下午4点35分

## API Usage Examples | API 使用示例

### Login | 登录
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

### Create Report Schedule | 创建定时报告
```bash
curl -X POST http://localhost:8080/api/reports/schedules \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_jwt_token>" \
  -d '{
    "name": "每日销售报告",
    "type": "daily",
    "query_ids": [1, 2, 3],
    "chart_ids": [1, 2],
    "template_ids": [1],
    "cron_pattern": "35 16 * * *"
  }'
```

## Error Handling | 错误处理

所有 API 错误响应均为 JSON 格式：

```json
{
  "code": 401,
  "message": "Token expired",
  "error": "Token expired: token is expired"
}
```

### Common Token Errors | 常见Token错误
- `Authorization header is required` - 缺少认证头
- `Invalid token` - Token无效
- `Token expired` - Token已过期
- `Token missing required claims` - Token缺少必要信息

## Security | 安全特性

- JWT authentication for all endpoints | 所有接口都需要 JWT 认证
- Password hashing with bcrypt | 使用 bcrypt 加密密码
- User data isolation | 用户数据隔离
- Configurable JWT token expiration | 可配置的JWT token过期时间
- Database credentials encryption | 数据库凭证加密

## Docker Deployment | Docker 部署

### Quick Start | 快速开始
```bash
docker-compose up --build
```

### Manual Build | 手动构建
```bash
docker build -t gobi .
docker run -p 8080:8080 gobi
```

## Project Structure | 项目结构

```
gobi/
├── cmd/server/main.go      # Main entry point | 主程序入口
├── config/                 # Configuration | 配置
├── internal/               # Internal packages | 内部包
│   ├── handlers/          # API handlers | API 处理函数
│   ├── middleware/        # Middleware | 中间件
│   └── models/           # Data models | 数据模型
├── pkg/                   # Public packages | 公共包
│   ├── database/         # Database | 数据库
│   ├── errors/           # Error handling | 错误处理
│   └── utils/            # Utilities | 工具
├── scripts/              # Scripts | 脚本
├── migrations/           # Database migrations | 数据库迁移
└── README.md            # Documentation | 文档
```

## Testing | 测试

Run automated tests: | 运行自动化测试：
```bash
chmod +x scripts/test_error_handling.sh
./scripts/test_error_handling.sh
```

## Contributing | 贡献

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## Support | 支持

For issues and questions, please create an issue on GitHub.
