# Gobi - 现代 Go 原生商业智能引擎

[English](./README.md)

🚀 **一个轻量级、API优先的商业智能引擎，使用 Go 构建** - 专为需要嵌入式分析、自动化报告和实时数据可视化的现代应用程序而设计。

## ✨ 为什么选择 Gobi？

- **🔧 Go 原生**: 完全使用 Go 构建，性能优异、简单易用、部署便捷
- **🔌 API 优先**: 提供 RESTful API，支持 JWT 和 API Key 认证，无缝集成
- **📊 多图表支持**: 从基础图表到高级 3D 可视化
- **🤖 自动化就绪**: 支持定时报告和 Webhook 通知
- **🔐 企业级安全**: 多用户隔离、API Key 管理、Webhook 签名验证
- **📈 生产就绪**: 服务层架构、全面的错误处理和日志记录

## 🎯 完美适用于

- **SaaS 应用程序** 需要嵌入式分析功能
- **微服务** 需要轻量级 BI 能力
- **内部工具** 用于数据可视化和报告
- **API 优先平台** 需要无头 BI 功能
- **Go 应用程序** 寻找原生 BI 集成

[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg)](https://go.dev/dl/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://github.com/sy-vendor/gobi/actions/workflows/go.yml/badge.svg)](https://github.com/sy-vendor/gobi/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sy-vendor/gobi)](https://goreportcard.com/report/github.com/sy-vendor/gobi)
[![GitHub stars](https://img.shields.io/github/stars/sy-vendor/gobi)](https://github.com/sy-vendor/gobi/stargazers)
[![API-First](https://img.shields.io/badge/API--First-Design-blueviolet)](https://github.com/sy-vendor/gobi)
[![3D Charts](https://img.shields.io/badge/3D--Charts-Supported-orange)](https://github.com/sy-vendor/gobi)

---

## 🚀 核心功能

### 🔌 **API 优先架构**
- 提供全面的 CRUD 操作的 RESTful API
- **API Key 认证** 用于服务间通信
- **Webhook 系统** 带签名验证的实时通知
- 统一的 JSON 响应格式和适当的错误处理

### 📊 **高级可视化**
- **12+ 图表类型**: 柱状图、折线图、饼图、散点图、雷达图、热力图、仪表盘、漏斗图
- **3D 图表**: 3D 柱状图、3D 散点图、3D 曲面图、3D 气泡图
- 交互式图表配置和自定义
- Excel 模板集成，生成专业报告

### 🤖 **自动化和调度**
- **基于 Cron 的调度** 用于自动化报告生成
- **Webhook 通知** 用于报告完成事件
- **重试逻辑** 带指数退避的失败重试
- **发送跟踪** 带详细日志记录

### 🔐 **企业级安全**
- **JWT 认证** 带可配置过期时间
- **API Key 管理** 带安全生成和吊销功能
- **多用户隔离** 确保数据隐私
- **Webhook 签名验证** 用于安全通知
- **基于角色的访问控制** (管理员/用户角色)

### 🏗️ **现代架构**
- **服务层模式** 用于清晰的关注点分离
- **依赖注入** 用于改进的可测试性
- **数据库连接池** 用于最佳性能
- **全面的错误处理** 带详细日志记录
- **配置管理** 支持 YAML

### 📈 **数据管理**
- **多数据库支持** (SQLite、MySQL、PostgreSQL)
- **SQL 查询管理** 带执行跟踪
- **数据源管理** 用于集中连接处理
- **查询缓存** 用于改进性能
- **仪表盘统计** 和分析

---

## 🎯 使用场景

### **嵌入式分析**
```go
// 将 BI 直接集成到您的 Go 应用程序中
client := gobi.NewClient("https://your-gobi-instance.com")
client.SetAPIKey("your-api-key")

// 以编程方式创建图表
chart := &gobi.Chart{
    Name: "销售分析",
    Type: "3d_surface",
    Data: salesData,
}
```

### **自动化报告**
```yaml
# 使用 Webhook 通知调度每日报告
schedule:
  name: "每日销售报告"
  cron: "0 9 * * *"  # 每天上午 9 点
  webhook: "https://your-app.com/webhooks/reports"
```

### **API 优先集成**
```bash
# 服务间认证
curl -H "Authorization: ApiKey your-api-key" \
     https://gobi.example.com/api/charts

# 实时 Webhook 通知
POST /webhooks/reports
{
  "event": "report.generated",
  "data": { "report_id": 123, "status": "success" }
}
```

---

## 🛠️ 技术栈

- **后端**: Go 1.21+ 与 Gin 框架
- **数据库**: SQLite (开发) / MySQL/PostgreSQL (生产)
- **认证**: JWT + API Keys，带 bcrypt 哈希
- **图表**: 自定义 3D 渲染，支持 WebGL
- **调度**: 基于 Cron，支持时区
- **通知**: Webhook 系统，带 HMAC 签名
- **文档**: 支持 OpenAPI/Swagger

---

## 📊 图表展示

| 图表类型 | 2D | 3D | 交互式 |
|------------|----|----|-------------|
| 柱状图 | ✅ | ✅ | ✅ |
| 折线图 | ✅ | ❌ | ✅ |
| 饼图 | ✅ | ❌ | ✅ |
| 散点图 | ✅ | ✅ | ✅ |
| 面积图 | ✅ | ❌ | ✅ |
| 曲面图 | ❌ | ✅ | ✅ |
| 热力图 | ✅ | ❌ | ✅ |
| 仪表盘 | ✅ | ❌ | ✅ |
| 漏斗图 | ✅ | ❌ | ✅ |
| 矩形树状图 | ✅ | ❌ | ✅ |
| 旭日图 | ✅ | ❌ | ✅ |
| 树形图 | ✅ | ❌ | ✅ |
| 箱线图 | ✅ | ❌ | ✅ |
| K线图 | ✅ | ❌ | ✅ |

---

## 🔧 环境要求

- Go 1.21 或更高版本
- SQLite（用于开发）
- MySQL/PostgreSQL（用于生产）

---

## ⚡ 快速开始

```bash
# 克隆并运行
git clone https://github.com/sy-vendor/gobi.git
cd gobi
go mod download
go run cmd/server/main.go

# 服务器在 http://localhost:8080 启动
# 默认管理员: admin/admin123
```

---

## 📋 配置

### 配置文件

应用程序使用 `config/config.yaml` 进行配置管理。

```yaml
default:
  server:
    port: "8080"
  jwt:
    secret: "default_jwt_secret"
    expiration_hours: 168
  database:
    type: "sqlite"
    dsn: "gobi.db"
```

### JWT 配置

- `jwt.secret`: JWT 签名密钥
- `jwt.expiration_hours`: Token 过期时间（小时）
  - 168 = 7 天
  - 720 = 30 天
  - 2160 = 90 天

---

## 🔌 API 接口

### 认证
- `POST /api/auth/register` — 注册新用户
- `POST /api/auth/login` — 登录并获取 JWT 令牌

### API Key 管理
- `POST /api/apikeys` — 创建新的 API Key
- `GET /api/apikeys` — 列出所有 API Key（用户自己的或管理员可查看所有）
- `DELETE /api/apikeys/:id` — 吊销 API Key

### Webhook 管理
- `POST /api/webhooks` — 创建新的 Webhook
- `GET /api/webhooks` — 列出所有 Webhook（用户自己的或管理员可查看所有）
- `GET /api/webhooks/:id` — 获取特定 Webhook
- `PUT /api/webhooks/:id` — 更新 Webhook
- `DELETE /api/webhooks/:id` — 删除 Webhook
- `GET /api/webhooks/:id/deliveries` — 列出 Webhook 发送记录
- `POST /api/webhooks/:id/test` — 测试 Webhook

### 仪表盘
- `GET /api/dashboard/stats` — 获取仪表盘统计信息

### 数据源
- `POST /api/datasources` — 创建新数据源
- `GET /api/datasources` — 列出所有数据源
- `GET /api/datasources/:id` — 获取特定数据源
- `PUT /api/datasources/:id` — 更新数据源
- `DELETE /api/datasources/:id` — 删除数据源

### 查询
- `POST /api/queries` — 创建新查询
- `GET /api/queries` — 列出所有查询
- `GET /api/queries/:id` — 获取特定查询
- `PUT /api/queries/:id` — 更新查询
- `DELETE /api/queries/:id` — 删除查询
- `POST /api/queries/:id/execute` — 执行查询

### 图表
- `POST /api/charts` — 创建新图表
- `GET /api/charts` — 列出所有图表
- `GET /api/charts/:id` — 获取特定图表
- `PUT /api/charts/:id` — 更新图表
- `DELETE /api/charts/:id` — 删除图表

### Excel 模板
- `POST /api/templates` — 上传新模板
- `GET /api/templates` — 列出所有模板
- `GET /api/templates/:id/download` — 下载模板

### 定时报告
- `POST /api/reports/schedules` — 创建新的定时报告
- `GET /api/reports/schedules` — 列出所有定时报告
- `GET /api/reports/schedules/:id` — 获取特定定时报告
- `PUT /api/reports/schedules/:id` — 更新定时报告
- `DELETE /api/reports/schedules/:id` — 删除定时报告

### 报告
- `GET /api/reports` — 列出所有生成的报告
- `GET /api/reports/:id/download` — 下载特定报告

---

## 📊 图表类型

支持的图表类型：
- 柱状图
- 折线图
- 饼图
- 散点图
- 雷达图
- 热力图
- 仪表盘
- 漏斗图
- 面积图
- 3D 柱状图
- 3D 散点图
- 3D 曲面图
- 3D 气泡图
- 矩形树状图（TreeMap）
- 旭日图
- 树形图（Tree Diagram）
- 箱线图（Box Plot）
- K线图/蜡烛图（Candlestick）

---

## 📊 图表类型说明

- 矩形树状图（TreeMap）：用嵌套矩形面积表达层级和数值，适合占比分析。
- 树形图（Tree Diagram）：用节点和连线表达层级关系，适合组织架构、家谱等分支结构。
- 箱线图（Box Plot）：用箱线图表示数据分布和五数概括，适合数据分析。
- K线图/蜡烛图（Candlestick）：用于股票/金融数据的可视化。

---

## ⏰ Cron 表达式指南

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

- `0 9 * * *` — 每天上午 9 点
- `0 0 * * 1` — 每周一午夜
- `35 16 * * *` — 每天下午 4 点 35 分

---

## 🔌 API 使用示例

### 登录

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

### 创建 API Key

```bash
curl -X POST http://localhost:8080/api/apikeys \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_jwt_token>" \
  -d '{
    "name": "我的服务 API Key",
    "expires_at": "2024-12-31T23:59:59Z"
  }'
```

**响应：**
```json
{
  "api_key": "abc123def456ghi789jkl012mno345pqr678stu901vwx234yz",
  "prefix": "abc123def456",
  "name": "我的服务 API Key",
  "expires_at": "2024-12-31T23:59:59Z",
  "created_at": "2024-01-15T10:30:00Z"
}
```

### 使用 API Key 认证
```bash
curl -X GET http://localhost:8080/api/queries \
  -H "Authorization: ApiKey abc123def456ghi789jkl012mno345pqr678stu901vwx234yz"
```

### 列出 API Keys
```bash
curl -X GET http://localhost:8080/api/apikeys \
  -H "Authorization: Bearer <your_jwt_token>"
```

### 吊销 API Key
```bash
curl -X DELETE http://localhost:8080/api/apikeys/1 \
  -H "Authorization: Bearer <your_jwt_token>"
```

### 创建 Webhook
```bash
curl -X POST http://localhost:8080/api/webhooks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_jwt_token>" \
  -d '{
    "name": "报告通知",
    "url": "https://your-app.com/webhooks/reports",
    "events": ["report.generated", "report.failed"],
    "headers": {
      "X-Custom-Header": "custom-value"
    }
  }'
```

### 测试 Webhook
```bash
curl -X POST http://localhost:8080/api/webhooks/1/test \
  -H "Authorization: Bearer <your_jwt_token>"
```

### 列出 Webhook 发送记录
```bash
curl -X GET http://localhost:8080/api/webhooks/1/deliveries \
  -H "Authorization: Bearer <your_jwt_token>"
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

---

## Webhook 事件

### 支持的事件

- `report.generated` — 报告生成成功
- `report.failed` — 报告生成失败
- `webhook.test` — Webhook 测试事件

### 事件数据格式

```json
{
  "event": "report.generated",
  "data": {
    "report_id": 123,
    "report_name": "每日销售报告",
    "schedule_id": 456,
    "schedule_name": "每日销售计划",
    "status": "success",
    "generated_at": "2024-01-15T10:30:00Z",
    "file_size": 1024,
    "download_url": "/api/reports/123/download"
  }
}
```

### Webhook 安全

- **签名验证**: 每个 webhook 都包含 HMAC-SHA256 签名
- **请求头**: 
  - `X-Gobi-Signature`: HMAC 签名
  - `X-Gobi-Timestamp`: Unix 时间戳
  - `X-Gobi-Event`: 事件类型
- **重试机制**: 自动重试，指数退避（3次尝试）
- **发送记录**: 所有发送尝试都会被记录

### 签名验证

```python
import hmac
import hashlib

def verify_signature(payload, signature, timestamp, secret):
    message = f"{timestamp}.{payload}"
    expected = hmac.new(
        secret.encode('utf-8'),
        message.encode('utf-8'),
        hashlib.sha256
    ).hexdigest()
    return hmac.compare_digest(expected, signature)
```

---

## 认证方式

### JWT 认证
使用 `Authorization: Bearer <jwt_token>` 头部进行用户认证。

### API Key 认证
使用 `Authorization: ApiKey <api_key>` 头部进行服务间认证。

**API Key 特性：**
- **安全生成**：32字节随机密钥，使用 bcrypt 哈希
- **前缀索引**：使用密钥前缀进行快速查找
- **过期支持**：可选的过期日期
- **吊销功能**：可以吊销密钥而不删除
- **用户隔离**：用户只能管理自己的密钥
- **管理员权限**：管理员可以管理所有密钥

**安全注意事项：**
- API Key 仅在创建时显示一次
- 请安全存储密钥，切勿提交到版本控制系统
- 生产环境请使用 HTTPS 保护密钥传输
- 定期轮换密钥以增强安全性

---

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
- `Invalid or expired API key` - API Key 无效或已过期

## 安全特性

- 所有接口都需要 JWT 认证
- **API Key 认证用于服务间通信**
- **Webhook 签名验证确保通知安全**
- 使用 bcrypt 加密密码
- **API Key 使用 bcrypt 哈希**
- 用户数据隔离
- **安全的随机密钥生成**
- **自动 Webhook 重试，指数退避** 