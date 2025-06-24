-- Area Charts Sample Data | 面积图示例数据
-- 用于测试面积图功能的示例数据 (SQLite兼容版本)

-- 创建销售趋势数据表 (用于基础面积图)
CREATE TABLE IF NOT EXISTS sales_trend (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    month TEXT NOT NULL,
    sales_amount REAL NOT NULL,
    product_category TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 插入销售趋势数据
INSERT INTO sales_trend (month, sales_amount, product_category) VALUES
('2024-01', 15000.00, 'Electronics'),
('2024-02', 18000.00, 'Electronics'),
('2024-03', 22000.00, 'Electronics'),
('2024-04', 25000.00, 'Electronics'),
('2024-05', 28000.00, 'Electronics'),
('2024-06', 32000.00, 'Electronics'),
('2024-07', 35000.00, 'Electronics'),
('2024-08', 38000.00, 'Electronics'),
('2024-09', 42000.00, 'Electronics'),
('2024-10', 45000.00, 'Electronics'),
('2024-11', 48000.00, 'Electronics'),
('2024-12', 52000.00, 'Electronics'),
('2024-01', 8000.00, 'Clothing'),
('2024-02', 9500.00, 'Clothing'),
('2024-03', 11000.00, 'Clothing'),
('2024-04', 12500.00, 'Clothing'),
('2024-05', 14000.00, 'Clothing'),
('2024-06', 16000.00, 'Clothing'),
('2024-07', 17500.00, 'Clothing'),
('2024-08', 19000.00, 'Clothing'),
('2024-09', 21000.00, 'Clothing'),
('2024-10', 22500.00, 'Clothing'),
('2024-11', 24000.00, 'Clothing'),
('2024-12', 26000.00, 'Clothing'),
('2024-01', 5000.00, 'Books'),
('2024-02', 6500.00, 'Books'),
('2024-03', 7200.00, 'Books'),
('2024-04', 8500.00, 'Books'),
('2024-05', 9200.00, 'Books'),
('2024-06', 10500.00, 'Books'),
('2024-07', 11500.00, 'Books'),
('2024-08', 12500.00, 'Books'),
('2024-09', 13500.00, 'Books'),
('2024-10', 14500.00, 'Books'),
('2024-11', 15500.00, 'Books'),
('2024-12', 16500.00, 'Books');

-- 创建用户增长数据表 (用于渐变面积图)
CREATE TABLE IF NOT EXISTS user_growth (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL,
    new_users INTEGER NOT NULL,
    total_users INTEGER NOT NULL,
    user_type TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 插入用户增长数据
INSERT INTO user_growth (date, new_users, total_users, user_type) VALUES
('2024-01-01', 150, 150, 'Free'),
('2024-01-02', 180, 330, 'Free'),
('2024-01-03', 220, 550, 'Free'),
('2024-01-04', 250, 800, 'Free'),
('2024-01-05', 280, 1080, 'Free'),
('2024-01-06', 320, 1400, 'Free'),
('2024-01-07', 350, 1750, 'Free'),
('2024-01-08', 380, 2130, 'Free'),
('2024-01-09', 420, 2550, 'Free'),
('2024-01-10', 450, 3000, 'Free'),
('2024-01-11', 480, 3480, 'Free'),
('2024-01-12', 520, 4000, 'Free'),
('2024-01-13', 550, 4550, 'Free'),
('2024-01-14', 580, 5130, 'Free'),
('2024-01-15', 620, 5750, 'Free'),
('2024-01-01', 50, 50, 'Premium'),
('2024-01-02', 65, 115, 'Premium'),
('2024-01-03', 80, 195, 'Premium'),
('2024-01-04', 95, 290, 'Premium'),
('2024-01-05', 110, 400, 'Premium'),
('2024-01-06', 125, 525, 'Premium'),
('2024-01-07', 140, 665, 'Premium'),
('2024-01-08', 155, 820, 'Premium'),
('2024-01-09', 170, 990, 'Premium'),
('2024-01-10', 185, 1175, 'Premium'),
('2024-01-11', 200, 1375, 'Premium'),
('2024-01-12', 215, 1590, 'Premium'),
('2024-01-13', 230, 1820, 'Premium'),
('2024-01-14', 245, 2065, 'Premium'),
('2024-01-15', 260, 2325, 'Premium');

-- 创建市场份额数据表 (用于堆叠面积图)
CREATE TABLE IF NOT EXISTS market_share (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    quarter TEXT NOT NULL,
    market_share REAL NOT NULL,
    company TEXT NOT NULL,
    industry TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 插入市场份额数据
INSERT INTO market_share (quarter, market_share, company, industry) VALUES
('2023-Q1', 25.5, 'Company A', 'Technology'),
('2023-Q2', 26.8, 'Company A', 'Technology'),
('2023-Q3', 28.2, 'Company A', 'Technology'),
('2023-Q4', 29.5, 'Company A', 'Technology'),
('2024-Q1', 30.8, 'Company A', 'Technology'),
('2024-Q2', 32.1, 'Company A', 'Technology'),
('2024-Q3', 33.5, 'Company A', 'Technology'),
('2024-Q4', 34.8, 'Company A', 'Technology'),
('2023-Q1', 18.2, 'Company B', 'Technology'),
('2023-Q2', 19.5, 'Company B', 'Technology'),
('2023-Q3', 20.8, 'Company B', 'Technology'),
('2023-Q4', 22.1, 'Company B', 'Technology'),
('2024-Q1', 23.4, 'Company B', 'Technology'),
('2024-Q2', 24.7, 'Company B', 'Technology'),
('2024-Q3', 26.0, 'Company B', 'Technology'),
('2024-Q4', 27.3, 'Company B', 'Technology'),
('2023-Q1', 15.8, 'Company C', 'Technology'),
('2023-Q2', 16.2, 'Company C', 'Technology'),
('2023-Q3', 16.8, 'Company C', 'Technology'),
('2023-Q4', 17.5, 'Company C', 'Technology'),
('2024-Q1', 18.2, 'Company C', 'Technology'),
('2024-Q2', 18.9, 'Company C', 'Technology'),
('2024-Q3', 19.6, 'Company C', 'Technology'),
('2024-Q4', 20.3, 'Company C', 'Technology'),
('2023-Q1', 40.5, 'Others', 'Technology'),
('2023-Q2', 37.5, 'Others', 'Technology'),
('2023-Q3', 34.2, 'Others', 'Technology'),
('2023-Q4', 30.9, 'Others', 'Technology'),
('2024-Q1', 27.6, 'Others', 'Technology'),
('2024-Q2', 24.3, 'Others', 'Technology'),
('2024-Q3', 20.9, 'Others', 'Technology'),
('2024-Q4', 17.6, 'Others', 'Technology');

-- 创建网站流量数据表 (用于平滑面积图)
CREATE TABLE IF NOT EXISTS website_traffic (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hour INTEGER NOT NULL,
    page_views INTEGER NOT NULL,
    unique_visitors INTEGER NOT NULL,
    bounce_rate REAL NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 插入网站流量数据 (24小时)
INSERT INTO website_traffic (hour, page_views, unique_visitors, bounce_rate) VALUES
(0, 1200, 800, 45.2),
(1, 980, 650, 48.1),
(2, 750, 520, 52.3),
(3, 620, 420, 55.8),
(4, 580, 380, 58.2),
(5, 650, 450, 54.7),
(6, 850, 580, 49.3),
(7, 1200, 820, 42.1),
(8, 1800, 1250, 38.5),
(9, 2200, 1500, 35.2),
(10, 2500, 1700, 32.8),
(11, 2800, 1900, 30.5),
(12, 3000, 2100, 28.9),
(13, 3200, 2250, 27.3),
(14, 3400, 2400, 25.8),
(15, 3600, 2550, 24.2),
(16, 3800, 2700, 22.7),
(17, 4000, 2850, 21.3),
(18, 4200, 3000, 19.8),
(19, 4400, 3150, 18.5),
(20, 4600, 3300, 17.2),
(21, 4800, 3450, 16.1),
(22, 5000, 3600, 15.3),
(23, 4800, 3450, 16.8);

-- 创建收入数据表 (用于累计面积图)
CREATE TABLE IF NOT EXISTS revenue_data (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    month TEXT NOT NULL,
    revenue REAL NOT NULL,
    expenses REAL NOT NULL,
    profit REAL NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 插入收入数据
INSERT INTO revenue_data (month, revenue, expenses, profit) VALUES
('2024-01', 50000.00, 35000.00, 15000.00),
('2024-02', 55000.00, 37000.00, 18000.00),
('2024-03', 60000.00, 39000.00, 21000.00),
('2024-04', 65000.00, 41000.00, 24000.00),
('2024-05', 70000.00, 43000.00, 27000.00),
('2024-06', 75000.00, 45000.00, 30000.00),
('2024-07', 80000.00, 47000.00, 33000.00),
('2024-08', 85000.00, 49000.00, 36000.00),
('2024-09', 90000.00, 51000.00, 39000.00),
('2024-10', 95000.00, 53000.00, 42000.00),
('2024-11', 100000.00, 55000.00, 45000.00),
('2024-12', 105000.00, 57000.00, 48000.00);

-- 创建示例查询语句
-- 这些查询可以直接在Gobi中使用来创建面积图

-- 1. 基础面积图查询 - 销售趋势
-- SELECT month as x, sales_amount as y FROM sales_trend WHERE product_category = 'Electronics' ORDER BY month

-- 2. 堆叠面积图查询 - 按产品类别分组
-- SELECT month as x, sales_amount as y, product_category as category FROM sales_trend ORDER BY month, product_category

-- 3. 用户增长面积图查询
-- SELECT date as x, new_users as y, user_type as category FROM user_growth ORDER BY date, user_type

-- 4. 市场份额堆叠面积图查询
-- SELECT quarter as x, market_share as y, company as category FROM market_share ORDER BY quarter, company

-- 5. 网站流量面积图查询
-- SELECT hour as x, page_views as y FROM website_traffic ORDER BY hour

-- 6. 收入趋势面积图查询
-- SELECT month as x, revenue as y FROM revenue_data ORDER BY month

-- 7. 利润vs支出对比面积图查询
-- SELECT month as x, profit as profit, expenses as expenses FROM revenue_data ORDER BY month

-- 8. 累计收入面积图查询
-- SELECT month as x, SUM(revenue) OVER (ORDER BY month) as cumulative_revenue FROM revenue_data ORDER BY month

-- 注意：SQLite不支持窗口函数，所以累计查询需要使用子查询或其他方式实现 