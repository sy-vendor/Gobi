# Gobi - BI 引擎最小可行产品

[English](./README.md)

一个使用 Go 语言构建的商业智能引擎最小可行产品。

## 功能特性

- SQL 查询管理和执行
- 交互式图表可视化
- Excel 模板管理和导出
- 用户认证和授权
- 用户数据隔离
- 仪表盘统计和分析
- **定时报告生成**
- **增强的JWT配置**
- **改进的错误处理**

## 环境要求

- Go 1.21 或更高版本
- SQLite（用于开发）
- MySQL/PostgreSQL（可选）

## 快速开始

```bash
git clone https://github.com/sy-vendor/gobi.git
cd gobi
go mod download
go run cmd/server/main.go
```

服务器默认在 8080 端口启动。

## 配置

### 配置文件

应用程序使用 `config/config.yaml` 进行配置管理。

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

### JWT配置

- `jwt.secret`: JWT签名密钥
- `jwt.expiration_hours`: Token过期时间（小时）
  - 168小时 = 7天
  - 720小时 = 30天
  - 2160小时 = 90天

## API 接口

### 认证
- POST /api/auth/register - 注册新用户
- POST /api/auth/login - 登录并获取 JWT 令牌

### 仪表盘
- GET /api/dashboard/stats - 获取仪表盘统计信息

### 数据源
- POST /api/datasources - 创建新数据源
- GET /api/datasources - 列出所有数据源
- GET /api/datasources/:id - 获取特定数据源
- PUT /api/datasources/:id - 更新数据源
- DELETE /api/datasources/:id - 删除数据源

### 查询
- POST /api/queries - 创建新查询
- GET /api/queries - 列出所有查询
- GET /api/queries/:id - 获取特定查询
- PUT /api/queries/:id - 更新查询
- DELETE /api/queries/:id - 删除查询
- POST /api/queries/:id/execute - 执行查询

### 图表
- POST /api/charts - 创建新图表
- GET /api/charts - 列出所有图表
- GET /api/charts/:id - 获取特定图表
- PUT /api/charts/:id - 更新图表
- DELETE /api/charts/:id - 删除图表

### Excel 模板
- POST /api/templates - 上传新模板
- GET /api/templates - 列出所有模板
- GET /api/templates/:id/download - 下载模板

### 定时报告
- POST /api/reports/schedules - 创建新的定时报告
- GET /api/reports/schedules - 列出所有定时报告
- GET /api/reports/schedules/:id - 获取特定定时报告
- PUT /api/reports/schedules/:id - 更新定时报告
- DELETE /api/reports/schedules/:id - 删除定时报告

### 报告
- GET /api/reports - 列出所有生成的报告
- GET /api/reports/:id/download - 下载特定报告

## 图表类型

支持的图表类型：
- 柱状图
- 折线图
- 饼图
- 散点图
- 雷达图
- 热力图
- 仪表盘
- 漏斗图
- **3D柱状图**
- **3D散点图**
- **3D曲面图**
- **3D气泡图**

## Cron表达式指南

### 基本格式
```
* * * * *
│ │ │ │ │
│ │ │ │ └── 星期几 (0-7)
│ │ │ └──── 月份 (1-12)
│ │ └────── 日期 (1-31)
│ └──────── 小时 (0-23)
└────────── 分钟 (0-59)
```

### 常见示例
- `0 9 * * *` - 每天上午9点
- `0 0 * * 1` - 每周一午夜
- `35 16 * * *` - 每天下午4点35分

## API 使用示例

### 登录
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

### 创建定时报告
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

## 错误处理

所有 API 错误响应均为 JSON 格式：

```json
{
  "code": 401,
  "message": "Token expired",
  "error": "Token expired: token is expired"
}
```

### 常见Token错误
- `Authorization header is required` - 缺少认证头
- `Invalid token` - Token无效
- `Token expired` - Token已过期
- `Token missing required claims` - Token缺少必要信息

## 安全特性

- 所有接口都需要 JWT 认证
- 使用 bcrypt 加密密码
- 用户数据隔离 