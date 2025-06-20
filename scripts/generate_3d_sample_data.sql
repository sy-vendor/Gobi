-- 3D Charts Sample Data | 3D图表示例数据
-- 用于测试3D图表功能的示例数据 (SQLite兼容版本)

-- 创建销售数据表 (用于3D柱状图)
CREATE TABLE IF NOT EXISTS sales_3d (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category TEXT NOT NULL,
    region TEXT NOT NULL,
    amount REAL NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 插入销售数据
INSERT INTO sales_3d (category, region, amount) VALUES
('Electronics', 'North', 15000.00),
('Electronics', 'South', 12000.00),
('Electronics', 'East', 18000.00),
('Electronics', 'West', 14000.00),
('Clothing', 'North', 8000.00),
('Clothing', 'South', 9500.00),
('Clothing', 'East', 11000.00),
('Clothing', 'West', 7500.00),
('Books', 'North', 5000.00),
('Books', 'South', 6500.00),
('Books', 'East', 7200.00),
('Books', 'West', 4800.00),
('Sports', 'North', 12000.00),
('Sports', 'South', 13500.00),
('Sports', 'East', 16000.00),
('Sports', 'West', 11000.00);

-- 创建产品性能数据表 (用于3D散点图和气泡图)
CREATE TABLE IF NOT EXISTS products_3d (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    performance_score REAL NOT NULL,
    price REAL NOT NULL,
    customer_rating REAL NOT NULL,
    product_category TEXT NOT NULL,
    sales_volume INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 插入产品数据
INSERT INTO products_3d (performance_score, price, customer_rating, product_category, sales_volume) VALUES
(8.5, 299.99, 4.2, 'Electronics', 1500),
(7.8, 199.99, 3.9, 'Electronics', 2200),
(9.1, 599.99, 4.5, 'Electronics', 800),
(6.5, 99.99, 3.7, 'Electronics', 3500),
(7.2, 149.99, 4.0, 'Clothing', 2800),
(8.0, 89.99, 4.3, 'Clothing', 3200),
(6.8, 129.99, 3.8, 'Clothing', 1800),
(9.3, 249.99, 4.6, 'Clothing', 950),
(7.5, 19.99, 4.1, 'Books', 5000),
(8.7, 29.99, 4.4, 'Books', 3800),
(6.2, 15.99, 3.6, 'Books', 6500),
(9.0, 39.99, 4.7, 'Books', 1200),
(8.2, 199.99, 4.2, 'Sports', 2100),
(7.6, 159.99, 3.9, 'Sports', 2800),
(9.4, 399.99, 4.8, 'Sports', 600),
(6.9, 119.99, 3.7, 'Sports', 3200);

-- 创建地形数据表 (用于3D曲面图)
CREATE TABLE IF NOT EXISTS terrain_3d (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    longitude REAL NOT NULL,
    latitude REAL NOT NULL,
    elevation REAL NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 插入地形数据 (模拟网格数据)
INSERT INTO terrain_3d (longitude, latitude, elevation) VALUES
-- 第一行
(120.0000, 30.0000, 100.50),
(120.0100, 30.0000, 105.20),
(120.0200, 30.0000, 110.80),
(120.0300, 30.0000, 108.30),
(120.0400, 30.0000, 102.10),
-- 第二行
(120.0000, 30.0100, 103.20),
(120.0100, 30.0100, 108.90),
(120.0200, 30.0100, 115.60),
(120.0300, 30.0100, 112.40),
(120.0400, 30.0100, 106.80),
-- 第三行
(120.0000, 30.0200, 107.80),
(120.0100, 30.0200, 113.50),
(120.0200, 30.0200, 120.20),
(120.0300, 30.0200, 117.90),
(120.0400, 30.0200, 111.30),
-- 第四行
(120.0000, 30.0300, 104.60),
(120.0100, 30.0300, 109.80),
(120.0200, 30.0300, 116.40),
(120.0300, 30.0300, 113.20),
(120.0400, 30.0300, 107.90),
-- 第五行
(120.0000, 30.0400, 101.20),
(120.0100, 30.0400, 106.40),
(120.0200, 30.0400, 112.10),
(120.0300, 30.0400, 109.80),
(120.0400, 30.0400, 103.50);

-- 创建城市数据表 (用于3D气泡图)
CREATE TABLE IF NOT EXISTS cities_3d (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    city_name TEXT NOT NULL,
    gdp REAL NOT NULL,  -- 单位：亿元
    population INTEGER NOT NULL,  -- 单位：万人
    area REAL NOT NULL,   -- 单位：平方公里
    tourism_income REAL NOT NULL, -- 单位：亿元
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 插入城市数据
INSERT INTO cities_3d (city_name, gdp, population, area, tourism_income) VALUES
('北京', 36102.60, 2154, 16410.54, 6224.60),
('上海', 38700.58, 2428, 6340.50, 5357.00),
('广州', 25019.11, 1530, 7434.40, 4454.00),
('深圳', 27670.24, 1344, 1997.47, 1721.00),
('杭州', 16106.00, 1036, 16596.00, 4005.00),
('南京', 14817.95, 931, 6587.02, 2785.00),
('武汉', 15616.06, 1121, 8569.15, 3570.00),
('成都', 17716.70, 1633, 14335.00, 4663.00),
('西安', 10688.28, 1020, 10108.00, 3146.00),
('重庆', 25002.79, 3205, 82402.00, 5734.00),
('天津', 14104.28, 1562, 11966.45, 3120.00),
('苏州', 20170.45, 1075, 8657.32, 2280.00),
('青岛', 12400.56, 949, 11282.00, 1905.00),
('长沙', 12142.52, 839, 11819.00, 1800.00),
('宁波', 12408.70, 854, 9816.00, 1650.00),
('无锡', 12370.48, 659, 4627.47, 1420.00);

-- 创建索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_sales_3d_category_region ON sales_3d(category, region);
CREATE INDEX IF NOT EXISTS idx_products_3d_category ON products_3d(product_category);
CREATE INDEX IF NOT EXISTS idx_terrain_3d_coords ON terrain_3d(longitude, latitude);
CREATE INDEX IF NOT EXISTS idx_cities_3d_name ON cities_3d(city_name);

-- 显示数据统计
SELECT 'sales_3d' as table_name, COUNT(*) as record_count FROM sales_3d
UNION ALL
SELECT 'products_3d' as table_name, COUNT(*) as record_count FROM products_3d
UNION ALL
SELECT 'terrain_3d' as table_name, COUNT(*) as record_count FROM terrain_3d
UNION ALL
SELECT 'cities_3d' as table_name, COUNT(*) as record_count FROM cities_3d; 