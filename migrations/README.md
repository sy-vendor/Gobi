# 数据库迁移使用说明（golang-migrate）

本项目使用 [golang-migrate/migrate](https://github.com/golang-migrate/migrate) 进行数据库结构管理，支持 MySQL、Postgres、SQLite。

---

## 1. 安装 migrate 工具

- macOS:
  ```bash
  brew install golang-migrate
  # 或
  curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.darwin-amd64.tar.gz | tar xvz
  sudo mv migrate /usr/local/bin/
  ```
- 其他平台请参考 [官方文档](https://github.com/golang-migrate/migrate#installation)

---

## 2. 目录结构

```
migrations/
  001_init.up.sql      # 初始化建表
  001_init.down.sql    # 回滚建表
  ...
  README.md            # 本说明
```

---

## 3. 执行迁移

### MySQL
```bash
migrate -path ./migrations -database "mysql://user:pass@tcp(localhost:3306)/gobi" up
```

### SQLite
```bash
migrate -path ./migrations -database "sqlite3://$(pwd)/gobi.db" up
```

### Postgres
```bash
migrate -path ./migrations -database "postgres://user:pass@localhost:5432/gobi?sslmode=disable" up
```

---

## 4. 回滚迁移

```bash
migrate -path ./migrations -database "<your-dsn>" down
```

---

## 5. 新建迁移脚本

```bash
migrate create -ext sql -dir migrations -seq add_new_table
```
会生成：
```
002_add_new_table.up.sql
002_add_new_table.down.sql
```

---

## 6. 常见问题

- **GORM AutoMigrate 与 migrate 冲突？**
  > 生产环境建议关闭 GORM 的 AutoMigrate，仅用 migrate 管理表结构。

- **迁移失败怎么办？**
  > 检查 SQL 语法、数据库连接、权限等。

- **如何查看当前迁移版本？**
  ```bash
  migrate -path ./migrations -database "<your-dsn>" version
  ```

---

更多用法请参考 [golang-migrate 官方文档](https://github.com/golang-migrate/migrate) 