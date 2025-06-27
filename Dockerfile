# 构建阶段
FROM golang:1.23-alpine AS builder
WORKDIR /app

# 安装必要的系统依赖（包括 gcc 和 musl-dev 支持 cgo/sqlite3）
RUN apk add --no-cache git ca-certificates tzdata gcc musl-dev

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o gobi ./cmd/server/main.go

# 运行阶段
FROM alpine:latest
WORKDIR /app

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata sqlite

# 从构建阶段复制二进制文件
COPY --from=builder /app/gobi .

# 复制配置文件和迁移文件
COPY config ./config
COPY migrations ./migrations

# 创建数据目录
RUN mkdir -p /app/data

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

# 运行应用
CMD ["/app/gobi"] 