# Gobi 配置管理系统

## 概述

Gobi 配置管理系统提供了完整的配置管理功能，支持多环境配置、环境变量覆盖、配置验证、热重载等特性。

## 特性

- ✅ **多环境支持** - 支持 default、dev、prod、test 等环境
- ✅ **环境变量覆盖** - 支持通过环境变量覆盖配置
- ✅ **配置验证** - 完整的配置项验证和类型检查
- ✅ **热重载** - 支持运行时配置文件变更自动重载
- ✅ **配置导出** - 支持导出为 YAML、JSON、环境变量格式
- ✅ **配置模板** - 提供配置模板生成功能
- ✅ **敏感信息处理** - 自动生成安全的密钥和密码
- ✅ **配置回调** - 支持配置变更时的回调通知

## 快速开始

### 1. 基本使用

```go
package main

import (
    "gobi/config"
    "log"
)

func main() {
    // 加载配置
    if err := config.LoadConfig(); err != nil {
        log.Fatal("Failed to load config:", err)
    }
    
    // 获取配置
    cfg := config.GetConfig()
    
    // 使用配置
    fmt.Printf("Server will start on port: %s\n", cfg.GetServerPort())
}
```

### 2. 环境变量设置

```bash
# 设置环境
export GOBI_ENV=prod

# 覆盖配置项
export GOBI_SERVER_PORT=9090
export GOBI_DATABASE_DSN="host=db.example.com user=app password=secret dbname=gobi"
export GOBI_JWT_SECRET="your-super-secret-jwt-key-here"

# 启动应用
./gobi
```

### 3. 配置变更监听

```go
// 注册配置变更回调
config.OnConfigChange(func(cfg *config.Config) {
    log.Println("Configuration has been updated")
    // 处理配置变更逻辑
})
```

## 配置结构

### Server 配置

```yaml
server:
  port: "8080"                    # 服务器端口
  host: "0.0.0.0"                # 服务器地址
  read_timeout: 30s              # 读取超时时间
  write_timeout: 30s             # 写入超时时间
  max_header_bytes: 1048576      # 最大请求头大小 (1MB)
  graceful_timeout: 30s          # 优雅关闭超时时间
  enable_https: false            # 是否启用HTTPS
  cert_file: ""                  # SSL证书文件路径
  key_file: ""                   # SSL私钥文件路径
```

### JWT 配置

```yaml
jwt:
  secret: "your-jwt-secret"      # JWT密钥 (至少32字符)
  expiration_hours: 24           # Token过期时间 (小时)
  refresh_expiration_hours: 168  # 刷新Token过期时间 (小时)
  issuer: "gobi"                 # Token发行者
  audience: "gobi-users"         # Token受众
  algorithm: "HS256"             # 签名算法
```

### Database 配置

```yaml
database:
  type: "postgres"               # 数据库类型: sqlite, mysql, postgres
  dsn: "connection-string"       # 数据库连接字符串
  max_open_conns: 25            # 最大连接数
  max_idle_conns: 5             # 最大空闲连接数
  conn_max_lifetime: 300s       # 连接最大生命周期
  connection_pool:              # 连接池配置
    max_open_conns: 25
    max_idle_conns: 5
    conn_max_lifetime: 300s
    conn_max_idle_time: 60s
  ssl:                          # SSL配置
    enabled: true
    mode: "require"
    cert_file: ""
    key_file: ""
    ca_file: ""
  retry:                        # 重试配置
    max_retries: 3
    retry_delay: 5s
    backoff_type: "exponential"  # linear, exponential
  migration:                    # 迁移配置
    auto_migrate: true
    path: "./migrations"
    table_name: "schema_migrations"
```

### Security 配置

```yaml
security:
  bcrypt_cost: 12               # bcrypt加密成本 (4-31)
  rate_limit: "100-M"           # 速率限制 (100 requests per minute)
  cors_origins:                 # CORS允许的源
    - "http://localhost:3000"
    - "https://yourdomain.com"
  allowed_headers:              # 允许的请求头
    - "Content-Type"
    - "Authorization"
    - "X-Requested-With"
  allowed_methods:              # 允许的HTTP方法
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
  trusted_proxies:              # 信任的代理地址
    - "127.0.0.1"
    - "::1"
  api_key_length: 32           # API密钥长度
  password_policy:             # 密码策略
    min_length: 8              # 最小长度
    require_uppercase: true    # 要求大写字母
    require_lowercase: true    # 要求小写字母
    require_numbers: true      # 要求数字
    require_symbols: false     # 要求特殊字符
```

### Logging 配置

```yaml
logging:
  level: "info"                 # 日志级别: debug, info, warn, error, fatal, panic
  format: "json"                # 日志格式: json, text
  output: "stdout"              # 输出目标: stdout, stderr, file
  file_path: ""                 # 文件路径 (当output为file时)
  max_size: 100                # 单个日志文件最大大小 (MB)
  max_backups: 3               # 最大备份文件数
  max_age: 28                  # 日志文件保留天数
  compress: true               # 是否压缩备份文件
  console: true                # 是否同时输出到控制台
```

### Cache 配置

```yaml
cache:
  enabled: true                # 是否启用缓存
  ttl: 300s                    # 默认缓存时间
  max_size: 1000              # 最大缓存项数
  strategy:                    # 缓存策略
    simple_query_ttl: 300s     # 简单查询缓存时间
    complex_query_ttl: 600s    # 复杂查询缓存时间
    max_cache_size: 1000       # 最大缓存大小
    hot_cache_enabled: true    # 是否启用热缓存
    hot_cache_ratio: 0.2       # 热缓存比例 (0-1)
    promotion_threshold: 3      # 提升到热缓存的访问次数
    business_hours_start: 9     # 工作时间开始 (小时)
    business_hours_end: 17      # 工作时间结束 (小时)
    adaptive_ttl: true         # 是否启用自适应TTL
    cache_warmup: true         # 是否启用缓存预热
    maintenance_interval: 300s  # 维护间隔
    eviction_policy: "lru"     # 淘汰策略: lru, lfu, fifo
    compression_enabled: false # 是否启用压缩
    metrics_enabled: true      # 是否启用指标收集
```

### Webhook 配置

```yaml
webhook:
  timeout: 30s                 # 请求超时时间
  max_retries: 3               # 最大重试次数
  retry_delay: 5s              # 重试延迟时间
  max_payload: 1048576         # 最大负载大小 (1MB)
  verify_ssl: true             # 是否验证SSL证书
  headers:                     # 默认请求头
    User-Agent: "Gobi-Webhook/1.0"
```

### Monitor 配置

```yaml
monitor:
  enabled: true                # 是否启用监控
  metrics_port: "9090"         # 指标端口
  health_check: true           # 是否启用健康检查
  profiling: false             # 是否启用性能分析
  alerting:                    # 告警配置
    enabled: false             # 是否启用告警
    channels:                  # 告警通道
      - "console"
      - "email"
      - "slack"
    thresholds:                # 告警阈值
      error_rate: 0.05         # 错误率阈值
      response_time: 1000      # 响应时间阈值 (ms)
      memory_usage: 0.8        # 内存使用率阈值
    cooldown: 300s             # 告警冷却时间
```

### API 配置

```yaml
api:
  version: "v1"                # API版本
  prefix: "/api"               # API前缀
  default_limit: 20            # 默认分页大小
  max_limit: 100               # 最大分页大小
  enable_swagger: true         # 是否启用Swagger文档
  enable_metrics: true         # 是否启用API指标
  enable_profiling: false      # 是否启用性能分析
```

## 环境变量

### 环境变量前缀

所有环境变量都以 `GOBI_` 为前缀，配置项中的点号会被转换为下划线。

### 常用环境变量

| 环境变量 | 配置项 | 说明 |
|---------|--------|------|
| `GOBI_ENV` | - | 环境名称 (default, dev, prod, test) |
| `GOBI_CONFIG_PATH` | - | 配置文件路径 |
| `GOBI_SERVER_PORT` | server.port | 服务器端口 |
| `GOBI_SERVER_HOST` | server.host | 服务器地址 |
| `GOBI_JWT_SECRET` | jwt.secret | JWT密钥 |
| `GOBI_JWT_EXPIRATION_HOURS` | jwt.expiration_hours | JWT过期时间 |
| `GOBI_DATABASE_TYPE` | database.type | 数据库类型 |
| `GOBI_DATABASE_DSN` | database.dsn | 数据库连接字符串 |
| `GOBI_SECURITY_BCRYPT_COST` | security.bcrypt_cost | bcrypt成本 |
| `GOBI_SECURITY_RATE_LIMIT` | security.rate_limit | 速率限制 |
| `GOBI_LOGGING_LEVEL` | logging.level | 日志级别 |
| `GOBI_LOGGING_FORMAT` | logging.format | 日志格式 |
| `GOBI_CACHE_ENABLED` | cache.enabled | 是否启用缓存 |
| `GOBI_CACHE_TTL` | cache.ttl | 缓存TTL |

## 配置验证

### 自动验证

配置加载时会自动进行验证，包括：

- 必需字段检查
- 数据类型验证
- 数值范围检查
- 枚举值验证
- 逻辑关系验证

### 手动验证

```go
validator := config.NewConfigValidator()
if err := validator.Validate(cfg); err != nil {
    log.Fatal("Config validation failed:", err)
}
```

## 配置导出

### 导出为YAML

```go
exporter := config.NewConfigExporter()
err := exporter.ExportToYAML(cfg, "config-export.yaml")
```

### 导出为JSON

```go
exporter := config.NewConfigExporter()
err := exporter.ExportToJSON(cfg, "config-export.json")
```

### 导出为环境变量

```go
exporter := config.NewConfigExporter()
err := exporter.ExportToEnv(cfg, ".env")
```

## 配置模板

### 生成配置模板

```go
template := config.NewConfigTemplate()
cfg := template.GenerateTemplate("prod")
err := template.SaveTemplate(cfg, "config-template.yaml")
```

## 最佳实践

### 1. 环境分离

- 使用不同的配置文件或环境变量来管理不同环境
- 生产环境不要使用默认密钥
- 敏感信息通过环境变量传递

### 2. 安全配置

```yaml
# 生产环境安全配置示例
security:
  bcrypt_cost: 12              # 高成本加密
  rate_limit: "100-M"          # 限制请求频率
  cors_origins:                # 严格限制CORS
    - "https://yourdomain.com"
  password_policy:             # 强密码策略
    min_length: 12
    require_uppercase: true
    require_lowercase: true
    require_numbers: true
    require_symbols: true
```

### 3. 数据库配置

```yaml
# 生产环境数据库配置示例
database:
  type: "postgres"
  dsn: "${DATABASE_URL}"       # 使用环境变量
  max_open_conns: 100          # 根据负载调整
  max_idle_conns: 10
  conn_max_lifetime: 300s
  ssl:
    enabled: true              # 启用SSL
    mode: "require"
  retry:
    max_retries: 5             # 增加重试次数
    retry_delay: 10s
    backoff_type: "exponential"
```

### 4. 监控配置

```yaml
# 生产环境监控配置示例
monitor:
  enabled: true
  metrics_port: "9090"
  health_check: true
  alerting:
    enabled: true
    channels:
      - "email"
      - "slack"
    thresholds:
      error_rate: 0.01         # 1%错误率告警
      response_time: 500       # 500ms响应时间告警
      memory_usage: 0.7        # 70%内存使用率告警
    cooldown: 600s             # 10分钟冷却时间
```

### 5. 日志配置

```yaml
# 生产环境日志配置示例
logging:
  level: "warn"                # 只记录警告及以上级别
  format: "json"               # 结构化日志
  output: "file"
  file_path: "/var/log/gobi/app.log"
  max_size: 100               # 100MB文件大小
  max_backups: 10             # 保留10个备份
  max_age: 90                 # 保留90天
  compress: true              # 压缩备份文件
  console: false              # 不输出到控制台
```

## 故障排除

### 常见问题

1. **配置加载失败**
   - 检查配置文件路径是否正确
   - 检查配置文件格式是否正确
   - 检查必需的环境变量是否设置

2. **配置验证失败**
   - 查看错误信息，修复相应的配置项
   - 检查数值范围是否正确
   - 检查枚举值是否有效

3. **热重载不工作**
   - 检查文件权限
   - 检查文件系统是否支持文件监听
   - 检查配置文件路径是否正确

4. **环境变量不生效**
   - 检查环境变量名称是否正确
   - 检查环境变量前缀是否为 `GOBI_`
   - 检查环境变量值格式是否正确

### 调试模式

设置环境变量启用调试模式：

```bash
export GOBI_LOGGING_LEVEL=debug
```

## 示例配置

### 开发环境

```yaml
dev:
  server:
    port: "8080"
    host: "0.0.0.0"
  jwt:
    secret: "dev-secret-key-change-in-production"
    expiration_hours: 168
  database:
    type: "mysql"
    dsn: "user:password@tcp(localhost:3306)/gobi_dev"
  security:
    bcrypt_cost: 10
    rate_limit: "1000-M"
  logging:
    level: "debug"
    format: "text"
  cache:
    enabled: true
    ttl: 300s
  monitor:
    enabled: true
    profiling: true
```

### 生产环境

```yaml
prod:
  server:
    port: "8080"
    host: "0.0.0.0"
  jwt:
    secret: "${JWT_SECRET}"  # 从环境变量获取
    expiration_hours: 24
  database:
    type: "postgres"
    dsn: "${DATABASE_URL}"   # 从环境变量获取
    ssl:
      enabled: true
      mode: "require"
  security:
    bcrypt_cost: 12
    rate_limit: "100-M"
    cors_origins:
      - "https://yourdomain.com"
  logging:
    level: "warn"
    format: "json"
    output: "file"
    file_path: "/var/log/gobi/app.log"
  cache:
    enabled: true
    ttl: 600s
  monitor:
    enabled: true
    alerting:
      enabled: true
      channels: ["email", "slack"]
```

### 测试环境

```yaml
test:
  server:
    port: "8081"
    host: "localhost"
  jwt:
    secret: "test-secret"
    expiration_hours: 1
  database:
    type: "sqlite"
    dsn: ":memory:"
  security:
    bcrypt_cost: 4
    rate_limit: "1000-M"
  logging:
    level: "error"
    format: "text"
  cache:
    enabled: false
  monitor:
    enabled: false
```

## 更新日志

### v1.0.0
- 初始版本发布
- 支持多环境配置
- 支持环境变量覆盖
- 支持配置验证
- 支持热重载
- 支持配置导出/导入
- 支持配置模板生成 