# Gobi - BI Engine MVP

[简体中文](./README.zh-CN.md)

A minimal viable product (MVP) for a Business Intelligence engine built with Go.

[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg)](https://go.dev/dl/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://github.com/sy-vendor/gobi/actions/workflows/go.yml/badge.svg)](https://github.com/sy-vendor/gobi/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sy-vendor/gobi)](https://goreportcard.com/report/github.com/sy-vendor/gobi)
[![GitHub stars](https://img.shields.io/github/stars/sy-vendor/gobi)](https://github.com/sy-vendor/gobi/stargazers)

---

## Features

- SQL query management and execution
- Interactive chart visualization
- Excel template management and export
- User authentication and authorization
- **API key support for service-to-service authentication**
- Data isolation between users
- Dashboard statistics and analytics
- Scheduled report generation
- Enhanced JWT configuration
- Improved error handling

---

## Prerequisites

- Go 1.21 or later
- SQLite (for development)
- MySQL/PostgreSQL (optional)

---

## Quick Start

```bash
git clone https://github.com/sy-vendor/gobi.git
cd gobi
go mod download
go run cmd/server/main.go
```

The server will start on port 8080 by default.

---

## Configuration

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

## API Endpoints

### Authentication
- `POST /api/auth/register` — Register a new user
- `POST /api/auth/login` — Login and get JWT token

### API Key Management
- `POST /api/apikeys` — Create a new API key
- `GET /api/apikeys` — List all API keys (user's own or all for admin)
- `DELETE /api/apikeys/:id` — Revoke an API key

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

## Chart Types

Supported chart types:
- Bar charts
- Line charts
- Pie charts
- Scatter plots
- Radar charts
- Heat maps
- Gauge charts
- Funnel charts
- 3D Bar charts
- 3D Scatter plots
- 3D Surface charts
- 3D Bubble charts

---

## Cron Expression Guide

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

## API Usage Examples

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

## Authentication Methods

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

## Error Handling

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

## Security

- JWT authentication for all endpoints
- **API key authentication for service-to-service communication**
- Password hashing with bcrypt
- **API key hashing with bcrypt**
- User data isolation
- **Secure random key generation**

---

## Docker Deployment

_Coming soon..._
