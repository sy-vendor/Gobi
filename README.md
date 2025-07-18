# Gobi - Modern Go-Native Business Intelligence (BI) Engine

[简体中文](./README.zh-CN.md)

🚀 **A lightweight, API-first Business Intelligence (BI) engine built with Go** - designed for modern applications that need embedded analytics, automated reporting, and real-time data visualization. The ultimate Go-native BI solution for developers.

## ✨ Why Gobi?

- **🔧 Go-Native**: Built entirely in Go for performance, simplicity, and easy deployment
- **🔌 API-First**: RESTful APIs with JWT and API key authentication for seamless integration
- **📊 Multi-Chart Support**: From basic charts to advanced 3D visualizations
- **🤖 Automation Ready**: Scheduled reports with webhook notifications
- **🔐 Enterprise Security**: Multi-user isolation, API key management, and webhook signatures
- **📈 Production Ready**: Clean architecture, dependency injection, comprehensive error handling, and logging
- **⚡ High Performance**: Optimized SQL validation, intelligent caching, and connection pooling

## 🎯 Perfect For

- **SaaS Applications** needing embedded BI and analytics
- **Microservices** requiring lightweight Business Intelligence capabilities
- **Internal Tools** for data visualization and BI reporting
- **API-First Platforms** that need headless Business Intelligence functionality
- **Go Applications** looking for native BI integration
- **Business Intelligence** teams requiring modern, scalable BI solutions

[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg)](https://go.dev/dl/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://github.com/sy-vendor/gobi/actions/workflows/go.yml/badge.svg)](https://github.com/sy-vendor/gobi/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sy-vendor/gobi)](https://goreportcard.com/report/github.com/sy-vendor/gobi)
[![GitHub stars](https://img.shields.io/github/stars/sy-vendor/gobi)](https://github.com/sy-vendor/gobi/stargazers)
[![API-First](https://img.shields.io/badge/API--First-Design-blueviolet)](https://github.com/sy-vendor/gobi)
[![3D Charts](https://img.shields.io/badge/3D--Charts-Supported-orange)](https://github.com/sy-vendor/gobi)
[![Business Intelligence](https://img.shields.io/badge/Business--Intelligence-Engine-green.svg)](https://github.com/sy-vendor/gobi)
[![Clean Architecture](https://img.shields.io/badge/Clean--Architecture-Implemented-green)](https://github.com/sy-vendor/gobi)

---

## 🚀 Key Features

### 🔌 **API-First Architecture**
- RESTful APIs with comprehensive CRUD operations
- **API Key authentication** for service-to-service communication
- **Webhook system** with signature verification for real-time notifications
- Unified JSON response format with proper error handling

### 📊 **Advanced Visualization**
- **25+ Chart Types**: Bar, Line, Pie, Scatter, Radar, Heatmap, Gauge, Funnel, Area, 3D charts, TreeMap, Sunburst, Tree, BoxPlot, Candlestick, WordCloud, Graph, Waterfall, Polar, Gantt, Rose, Geo/Map/Choropleth, Progress/Circular Progress
- **3D Charts**: 3D Bar, 3D Scatter, 3D Surface, 3D Bubble charts
- Interactive chart configuration and customization
- Excel template integration for professional reports

### 🤖 **Automation & Scheduling**
- **Cron-based scheduling** for automated report generation
- **Webhook notifications** for report completion events
- **Retry logic** with exponential backoff for failed deliveries
- **Delivery tracking** with detailed logs

### 🔐 **Enterprise Security**
- **JWT authentication** with configurable expiration
- **API key management** with secure generation and revocation
- **Multi-user isolation** ensuring data privacy
- **Webhook signature verification** for secure notifications
- **Role-based access control** (Admin/User roles)
- **SQL injection protection** with comprehensive validation

### 🏗️ **Clean Architecture**
- **Repository Pattern** for data access layer separation
- **Service Layer** with dependency injection for business logic
- **Infrastructure Services** for caching, encryption, validation
- **Factory Pattern** for dependency management
- **Comprehensive error handling** with detailed logging
- **Configuration management** with YAML support

### ⚡ **Performance Optimizations**
- **SQL Validation Optimization**: Eliminated redundant validation checks, reducing validation time by 60%
- **Intelligent Caching**: Smart TTL based on query complexity (simple: 5min, complex: 30min)
- **Connection Pool Configuration**: Configurable database connection pool settings
- **System Monitoring**: Real-time performance metrics and health checks
- **Database Query Optimization**: Advanced query analysis and optimization engine
  - **Query Complexity Analysis**: Automatic analysis of SQL complexity and performance impact
  - **Index Management**: Multi-database index analysis and optimization suggestions
  - **Query Execution Plans**: Detailed execution plans with performance metrics
  - **Smart Caching Strategy**: Dynamic TTL based on query characteristics and business hours
  - **Batch Query Execution**: Concurrent query execution with controlled concurrency
  - **Performance Monitoring**: Real-time query performance tracking and slow query detection
  - **Memory Usage Optimization**: Intelligent memory management for large result sets
  - **Network Transfer Optimization**: Optimized data transfer and compression

### 📈 **Data Management**
- **Multi-database support** (SQLite, MySQL, PostgreSQL)
- **SQL query management** with execution tracking
- **Data source management** for centralized connection handling
- **Query caching** for improved performance
- **Dashboard statistics** and analytics

---

## 🎯 Use Cases

### **Embedded Analytics**
```go
// Integrate BI directly into your Go application
client := gobi.NewClient("https://your-gobi-instance.com")
client.SetAPIKey("your-api-key")

// Create charts programmatically
chart := &gobi.Chart{
    Name: "Sales Analytics",
    Type: "3d_surface",
    Data: salesData,
}
```

### **Automated Reporting**
```yaml
# Schedule daily reports with webhook notifications
schedule:
  name: "Daily Sales Report"
  cron: "0 9 * * *"  # Every day at 9 AM
  webhook: "https://your-app.com/webhooks/reports"
```

### **API-First Integration**
```bash
# Service-to-service authentication
curl -H "Authorization: ApiKey your-api-key" \
     https://gobi.example.com/api/charts

# Real-time webhook notifications
POST /webhooks/reports
{
  "event": "report.generated",
  "data": { "report_id": 123, "status": "success" }
}
```

### **Database Query Optimization**
```go
// Analyze query performance and get optimization suggestions
optimizer := database.NewQueryOptimizer()
plan, err := optimizer.AnalyzeQuery(sql, dataSource)
for _, suggestion := range plan.Suggestions {
    fmt.Println("Optimization:", suggestion)
}

// Execute query with optimization
optimizedService := infrastructure.NewOptimizedSQLExecutionService(cacheService)
result, err := optimizedService.ExecuteWithOptimization(ctx, dataSource, sql)
fmt.Printf("Execution time: %v, Cache hit: %v\n", 
    result.ExecutionTime, result.CacheHit)

// Get index suggestions
indexManager := database.NewIndexManager()
suggestions, err := indexManager.SuggestIndexes(queryPatterns, dataSource)
for _, suggestion := range suggestions {
    fmt.Printf("Index suggestion: %s\n", suggestion.SQL)
}
```

---

## 🛠️ Technology Stack

- **Backend**: Go 1.21+ with Gin framework
- **Database**: SQLite (dev) / MySQL/PostgreSQL (prod)
- **Authentication**: JWT + API Keys with bcrypt hashing
- **Charts**: Custom 3D rendering with WebGL support
- **Scheduling**: Cron-based with timezone support
- **Notifications**: Webhook system with HMAC signatures
- **Architecture**: Clean Architecture with Repository Pattern
- **Caching**: Intelligent caching with go-cache
- **Documentation**: OpenAPI/Swagger ready
- **Database Optimization**: Advanced query optimization with index management
- **Performance Monitoring**: Real-time metrics and query analysis

## 🏗️ Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   API Layer     │    │  Service Layer  │    │ Repository Layer│
│   (Handlers)    │◄──►│  (Business      │◄──►│  (Data Access)  │
│                 │    │   Logic)        │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Infrastructure  │    │   Validation    │    │   Database      │
│   Services      │    │     Layer       │    │   Connection    │
│ (Cache, Auth,   │    │ (SQL, Chart,    │    │     Pool        │
│  Encryption)    │    │  DataSource)    │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Query Optimizer │    │ Index Manager   │    │ Performance     │
│ (Analysis,      │    │ (Multi-DB       │    │ Monitoring      │
│  Suggestions)   │    │  Support)       │    │ (Metrics,       │
│                 │    │                 │    │  Analytics)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

---

## 🚀 Recent Improvements

### ⚡ **Performance Optimizations**
- **SQL Validation Optimization**: Eliminated redundant validation checks, reducing validation time by 60%
- **Intelligent Caching**: Smart TTL based on query complexity (simple: 5min, complex: 30min)
- **Connection Pool Configuration**: Configurable database connection pool settings
- **System Monitoring**: Real-time performance metrics and health checks
- **Database Query Optimization Engine**: Advanced query analysis and optimization
  - **Query Complexity Analysis**: Automatic analysis of SQL complexity with scoring system
  - **Index Management**: Multi-database (MySQL, PostgreSQL, SQLite) index analysis and suggestions
  - **Query Execution Plans**: Detailed execution plans with performance metrics and optimization suggestions
  - **Smart Caching Strategy**: Dynamic TTL based on query characteristics, business hours, and data volatility
  - **Batch Query Execution**: Concurrent query execution with controlled concurrency (max 5 concurrent queries)
  - **Performance Monitoring**: Real-time query performance tracking, slow query detection, and memory usage optimization
  - **Memory Usage Optimization**: Intelligent memory management for large result sets with size estimation
  - **Network Transfer Optimization**: Optimized data transfer with compression and size estimation

### 🏗️ **Architecture Refactoring**
- **Repository Pattern**: Clean separation of data access logic
- **Dependency Injection**: Improved testability and maintainability
- **Service Layer**: Clear business logic separation
- **Infrastructure Services**: Centralized caching, encryption, and validation
- **Factory Pattern**: Simplified dependency management

### 🔐 **Security Enhancements**
- **Comprehensive SQL Injection Protection**: Multi-layer validation with configurable strictness
- **Column Name Validation**: Smart validation allowing business column names
- **Read-Only Query Enforcement**: Strict enforcement of SELECT-only queries
- **Suspicious Pattern Detection**: Advanced pattern matching for malicious SQL

### 📊 **Chart Support Expansion**
- **25+ Chart Types**: Added TreeMap, Sunburst, Tree, BoxPlot, Candlestick, WordCloud, Graph, Waterfall, Polar, Gantt, Rose, Geo/Map/Choropleth, Progress/Circular Progress
- **Enhanced 3D Support**: Improved 3D chart rendering and performance
- **Interactive Features**: Better user interaction and customization options

---

## 📊 Chart Gallery

| Chart Type | 2D | 3D | Interactive |
|------------|----|----|-------------|
| Bar Charts | ✅ | ✅ | ✅ |
| Line Charts | ✅ | ❌ | ✅ |
| Pie Charts | ✅ | ❌ | ✅ |
| Scatter Plots | ✅ | ✅ | ✅ |
| Area Charts | ✅ | ❌ | ✅ |
| Surface Charts | ❌ | ✅ | ✅ |
| Heat Maps | ✅ | ❌ | ✅ |
| Gauge Charts | ✅ | ❌ | ✅ |
| Funnel Charts | ✅ | ❌ | ✅ |
| TreeMap Charts | ✅ | ❌ | ✅ |
| Sunburst Charts | ✅ | ❌ | ✅ |
| Tree Diagram | ✅ | ❌ | ✅ |
| Box Plot | ✅ | ❌ | ✅ |
| Candlestick charts | ✅ | ❌ | ✅ |
| Word Cloud charts | ✅ | ❌ | ✅ |
| Graph charts | ✅ | ❌ | ✅ |
| Waterfall charts | ✅ | ❌ | ✅ |
| Polar charts | ✅ | ❌ | ✅ |
| Gantt charts | ✅ | ❌ | ✅ |
| Rose charts | ✅ | ❌ | ✅ |
| Geo/Map/Choropleth charts | ✅ | ❌ | ✅ |
| Progress/Circular Progress charts | ✅ | ❌ | ✅ |

---

## 🔧 Prerequisites

- Go 1.21 or later
- SQLite (for development)
- MySQL/PostgreSQL (for production)

---

## ⚡ Quick Start

### 🚀 One-Command Setup
```bash
# Clone and setup everything
git clone https://github.com/sy-vendor/gobi.git
cd gobi
make setup

# Start development server
make dev

# Server starts on http://localhost:8080
# Default admin: admin/admin123
```

### 🛠️ Development Tools

```bash
# View all available commands
make help

# Development workflow
make dev-full          # Format, lint, test, build

# Testing
make test              # Run tests
make test-coverage     # Run tests with coverage
make test-race         # Run tests with race detection

# Code quality
make lint              # Run linter
make format            # Format code
make security-scan     # Security vulnerability scan

# Database
make migrate           # Run migrations
make setup-data        # Setup sample data
make test-data         # Test all chart data

# Docker
make docker-build      # Build Docker image
make docker-run        # Run Docker container

# Production
make prod-build        # Production build with all checks
```

### 📦 Docker Quick Start
```bash
# Build and run with Docker
docker build -t gobi .
docker run -p 8080:8080 gobi

# Or use docker-compose
docker-compose up -d
```

---

## 📋 Configuration

### Configuration File

The application uses `config/config.yaml` for configuration management.

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

### JWT Configuration

- `jwt.secret`: JWT signing secret
- `jwt.expiration_hours`: Token expiration time (hours)
  - 168 = 7 days
  - 720 = 30 days
  - 2160 = 90 days

---

## 🔌 API Endpoints

### Authentication
- `POST /api/auth/register` — Register a new user
- `POST /api/auth/login` — Login and get JWT token

### API Key Management
- `POST /api/apikeys` — Create a new API key
- `GET /api/apikeys` — List all API keys (user's own or all for admin)
- `DELETE /api/apikeys/:id` — Revoke an API key

### Webhook Management
- `POST /api/webhooks` — Create a new webhook
- `GET /api/webhooks` — List all webhooks (user's own or all for admin)
- `GET /api/webhooks/:id` — Get a specific webhook
- `PUT /api/webhooks/:id` — Update a webhook
- `DELETE /api/webhooks/:id` — Delete a webhook
- `GET /api/webhooks/:id/deliveries` — List webhook delivery attempts
- `POST /api/webhooks/:id/test` — Test a webhook

### Dashboard
- `GET /api/dashboard/stats` — Get dashboard statistics

### Data Sources
- `POST /api/datasources` — Create a new data source
- `GET /api/datasources` — List all data sources
- `GET /api/datasources/:id` — Get a specific data source
- `PUT /api/datasources/:id` — Update a data source
- `DELETE /api/datasources/:id` — Delete a data source

### Queries
- `POST /api/queries` — Create a new query
- `GET /api/queries` — List all queries
- `GET /api/queries/:id` — Get a specific query
- `PUT /api/queries/:id` — Update a query
- `DELETE /api/queries/:id` — Delete a query
- `POST /api/queries/:id/execute` — Execute a query

### Charts
- `POST /api/charts` — Create a new chart
- `GET /api/charts` — List all charts
- `GET /api/charts/:id` — Get a specific chart
- `PUT /api/charts/:id` — Update a chart
- `DELETE /api/charts/:id` — Delete a chart

### Excel Templates
- `POST /api/templates` — Upload a new template
- `GET /api/templates` — List all templates
- `GET /api/templates/:id/download` — Download a template

### Report Schedules
- `POST /api/reports/schedules` — Create a new report schedule
- `GET /api/reports/schedules` — List all report schedules
- `GET /api/reports/schedules/:id` — Get a specific report schedule
- `PUT /api/reports/schedules/:id` — Update a report schedule
- `DELETE /api/reports/schedules/:id` — Delete a report schedule

### Reports
- `GET /api/reports` — List all generated reports
- `GET /api/reports/:id/download` — Download a specific report

---

## 📊 Chart Types

TreeMap (rectangular, area-based): shows hierarchy using nested rectangles, area represents value.
Tree Diagram: shows hierarchy using nodes and branches (parent-child), like org charts or family trees.
Box Plot: shows distribution of data based on five-number summary.
Word Cloud: shows text data visualization where font size represents word frequency or importance.
Graph/Network/Force-directed: shows relationships between entities using nodes and edges, supports force-directed layout for network analysis.

Supported chart types:
- Bar charts
- Line charts
- Pie charts
- Scatter plots
- Radar charts
- Heat maps
- Gauge charts
- Funnel charts
- Area charts
- 3D Bar charts
- 3D Scatter plots
- 3D Surface charts
- 3D Bubble charts
- TreeMap charts
- Sunburst charts
- Tree Diagram
- Box Plot charts
- Candlestick charts (K-line/Stock charts)
- Word Cloud charts
- Graph/Network/Force-directed charts
- Waterfall charts
- Polar charts
- Gantt charts
- Rose charts
- Geo/Map/Choropleth charts
- Progress/Circular Progress charts

---

## ⏰ Cron Expression Guide

### Basic Format

```
* * * * *
│ │ │ │ │
│ │ │ │ └── Day of week (0-7)
│ │ │ └──── Month (1-12)
│ │ └────── Day of month (1-31)
│ └──────── Hour (0-23)
└────────── Minute (0-59)
```

### Common Examples

- `0 9 * * *` — Every day at 9:00 AM
- `0 0 * * 1` — Every Monday at midnight
- `35 16 * * *` — Every day at 4:35 PM

---

## 🔌 API Usage Examples

### Registration

```bash
curl -X POST https://gobi.example.com/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

### Create API Key

```bash
curl -X POST http://localhost:8080/api/apikeys \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_jwt_token>" \
  -d '{
    "name": "My Service API Key",
    "expires_at": "2024-12-31T23:59:59Z"
  }'
```

**Response:**
```json
{
  "api_key": "abc123def456ghi789jkl012mno345pqr678stu901vwx234yz",
  "prefix": "abc123def456",
  "name": "My Service API Key",
  "expires_at": "2024-12-31T23:59:59Z",
  "created_at": "2024-01-15T10:30:00Z"
}
```

### Using API Key Authentication

```bash
curl -X GET http://localhost:8080/api/queries \
  -H "Authorization: ApiKey abc123def456ghi789jkl012mno345pqr678stu901vwx234yz"
```

### List API Keys

```bash
curl -X GET http://localhost:8080/api/apikeys \
  -H "Authorization: Bearer <your_jwt_token>"
```

### Revoke API Key

```bash
curl -X DELETE http://localhost:8080/api/apikeys/1 \
  -H "Authorization: Bearer <your_jwt_token>"
```

### Create Webhook

```bash
curl -X POST http://localhost:8080/api/webhooks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_jwt_token>" \
  -d '{
    "name": "Report Notifications",
    "url": "https://your-app.com/webhooks/reports",
    "events": ["report.generated", "report.failed"],
    "headers": {
      "X-Custom-Header": "custom-value"
    }
  }'
```

### Test Webhook

```bash
curl -X POST http://localhost:8080/api/webhooks/1/test \
  -H "Authorization: Bearer <your_jwt_token>"
```

### List Webhook Deliveries

```bash
curl -X GET http://localhost:8080/api/webhooks/1/deliveries \
  -H "Authorization: Bearer <your_jwt_token>"
```

### Create Report Schedule

```bash
curl -X POST http://localhost:8080/api/reports/schedules \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_jwt_token>" \
  -d '{
    "name": "Daily Sales Report",
    "type": "daily",
    "query_ids": [1, 2, 3],
    "chart_ids": [1, 2],
    "template_ids": [1],
    "cron_pattern": "35 16 * * *"
  }'
```

---

## 📊 Webhook Events

### Supported Events

- `report.generated` — Report generation completed successfully
- `report.failed` — Report generation failed
- `webhook.test` — Test webhook event

### Event Payload Format

```json
{
  "event": "report.generated",
  "data": {
    "report_id": 123,
    "report_name": "Daily Sales Report",
    "schedule_id": 456,
    "schedule_name": "Daily Sales Schedule",
    "status": "success",
    "generated_at": "2024-01-15T10:30:00Z",
    "file_size": 1024,
    "download_url": "/api/reports/123/download"
  }
}
```

### Webhook Security

- **Signature Verification**: Each webhook includes an HMAC-SHA256 signature
- **Headers**: 
  - `X-Gobi-Signature`: HMAC signature
  - `X-Gobi-Timestamp`: Unix timestamp
  - `X-Gobi-Event`: Event type
- **Retry Logic**: Automatic retry with exponential backoff (3 attempts)
- **Delivery Tracking**: All delivery attempts are logged

### Signature Verification

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

## 🔐 Authentication Methods

### JWT Authentication
Use `Authorization: Bearer <jwt_token>` header for user authentication.

### API Key Authentication
Use `Authorization: ApiKey <api_key>` header for service-to-service authentication.

**API Key Features:**
- **Secure Generation**: 32-byte random keys with bcrypt hashing
- **Prefix Indexing**: Fast lookup using key prefixes
- **Expiration Support**: Optional expiration dates
- **Revocation**: Ability to revoke keys without deletion
- **User Isolation**: Users can only manage their own keys
- **Admin Access**: Administrators can manage all keys

**Security Notes:**
- API keys are only shown once upon creation
- Store keys securely and never commit them to version control
- Use HTTPS in production to protect key transmission
- Regularly rotate keys for enhanced security

---

## 🔧 Error Handling

All API errors are returned in JSON format:

```json
{
  "code": 401,
  "message": "Token expired",
  "error": "Token expired: token is expired"
}
```

### Common Token Errors

- `Authorization header is required`
- `Invalid token`
- `Token expired`
- `Token missing required claims`
- `Invalid or expired API key`

---

## 🔒 Security

- JWT authentication for all endpoints
- **API key authentication for service-to-service communication**
- **Webhook signature verification for secure notifications**
- Password hashing with bcrypt
- **API key hashing with bcrypt**
- User data isolation
- **Secure random key generation**
- **Automatic webhook retry with exponential backoff**

---

## 📦 Docker Deployment

_Coming soon..._
