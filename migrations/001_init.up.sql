-- 用户表
CREATE TABLE users (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  username VARCHAR(64) NOT NULL UNIQUE,
  email VARCHAR(128) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  role VARCHAR(16) NOT NULL,
  last_login DATETIME,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME
);

-- 数据源表
CREATE TABLE datasources (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT,
  name VARCHAR(255) NOT NULL,
  type VARCHAR(32) NOT NULL,
  host VARCHAR(255),
  port INT,
  database VARCHAR(255),
  username VARCHAR(255),
  password VARCHAR(255),
  description TEXT,
  is_public BOOLEAN,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME
);

-- 查询表
CREATE TABLE queries (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT,
  data_source_id BIGINT,
  name VARCHAR(255),
  sql TEXT,
  description TEXT,
  is_public BOOLEAN,
  exec_count BIGINT DEFAULT 0,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME
);

-- 图表表
CREATE TABLE charts (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  query_id BIGINT,
  user_id BIGINT,
  name VARCHAR(255),
  type VARCHAR(32),
  config TEXT,
  data TEXT,
  description TEXT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME
);

-- Excel 模板表
CREATE TABLE excel_templates (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT,
  name VARCHAR(255),
  template BLOB,
  description TEXT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME
); 