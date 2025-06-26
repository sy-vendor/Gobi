-- Rose Chart Sample Data | 玫瑰图示例数据
-- 用于测试玫瑰图功能的示例数据 (SQLite兼容)

CREATE TABLE IF NOT EXISTS rose_demo (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category TEXT NOT NULL,       -- 类别（如方向、月份等）
    value REAL NOT NULL,          -- 数值
    angle REAL,                   -- 角度（可选，通常自动计算）
    color TEXT,                   -- 颜色（可选）
    description TEXT,             -- 说明
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 示例数据1：风向玫瑰图（Wind Rose）
INSERT INTO rose_demo (category, value, angle, color, description) VALUES
('N', 120, 45, '#1890ff', '北风'),
('NE', 150, 45, '#2fc25b', '东北风'),
('E', 180, 45, '#facc14', '东风'),
('SE', 90, 45, '#f5222d', '东南风'),
('S', 60, 45, '#722ed1', '南风'),
('SW', 80, 45, '#13c2c2', '西南风'),
('W', 110, 45, '#eb2f96', '西风'),
('NW', 100, 45, '#fa8c16', '西北风');

-- 示例数据2：月份销售玫瑰图
INSERT INTO rose_demo (category, value, angle, color, description) VALUES
('1月', 20000, 30, '#1890ff', '1月销售'),
('2月', 22000, 30, '#2fc25b', '2月销售'),
('3月', 25000, 30, '#facc14', '3月销售'),
('4月', 21000, 30, '#f5222d', '4月销售'),
('5月', 23000, 30, '#722ed1', '5月销售'),
('6月', 26000, 30, '#13c2c2', '6月销售'),
('7月', 24000, 30, '#eb2f96', '7月销售'),
('8月', 27000, 30, '#fa8c16', '8月销售'),
('9月', 29000, 30, '#a0d911', '9月销售'),
('10月', 31000, 30, '#52c41a', '10月销售'),
('11月', 33000, 30, '#fa541c', '11月销售'),
('12月', 35000, 30, '#eb2f96', '12月销售');

-- 示例数据3：用户活跃度玫瑰图
INSERT INTO rose_demo (category, value, angle, color, description) VALUES
('00:00-06:00', 500, 60, '#1890ff', '凌晨时段'),
('06:00-12:00', 1200, 60, '#2fc25b', '上午时段'),
('12:00-18:00', 1800, 60, '#facc14', '下午时段'),
('18:00-24:00', 1500, 60, '#f5222d', '晚上时段');

-- 示例查询：
-- SELECT category, value, angle, color, description FROM rose_demo ORDER BY id; 