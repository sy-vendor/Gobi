# 构建阶段
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o gobi ./cmd/server/main.go

# 运行阶段
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/gobi .
COPY config ./config
COPY migrations ./migrations
COPY .env .env
COPY config.yaml config.yaml
EXPOSE 8080
CMD ["/app/gobi"] 