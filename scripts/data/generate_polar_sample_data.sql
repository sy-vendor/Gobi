-- Polar Chart Sample Data | 极坐标图示例数据
-- 用于测试极坐标图功能的示例数据 (SQLite兼容)

CREATE TABLE IF NOT EXISTS polar_demo (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    angle TEXT NOT NULL,        -- 角度/类别（如方向、月份等）
    value REAL NOT NULL,        -- 数值
    category TEXT,              -- 分组（可选）
    description TEXT,           -- 说明
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 示例数据1：风向玫瑰图（Wind Rose）
INSERT INTO polar_demo (angle, value, category, description) VALUES
('N', 120, '风速', '北风'),
('NE', 150, '风速', '东北风'),
('E', 180, '风速', '东风'),
('SE', 90, '风速', '东南风'),
('S', 60, '风速', '南风'),
('SW', 80, '风速', '西南风'),
('W', 110, '风速', '西风'),
('NW', 100, '风速', '西北风');

-- 示例数据2：月份销售极坐标图
INSERT INTO polar_demo (angle, value, category, description) VALUES
('1月', 20000, '电子', '1月电子产品销售'),
('2月', 22000, '电子', '2月电子产品销售'),
('3月', 25000, '电子', '3月电子产品销售'),
('4月', 21000, '电子', '4月电子产品销售'),
('5月', 23000, '电子', '5月电子产品销售'),
('6月', 26000, '电子', '6月电子产品销售'),
('7月', 24000, '电子', '7月电子产品销售'),
('8月', 27000, '电子', '8月电子产品销售'),
('9月', 29000, '电子', '9月电子产品销售'),
('10月', 31000, '电子', '10月电子产品销售'),
('11月', 33000, '电子', '11月电子产品销售'),
('12月', 35000, '电子', '12月电子产品销售');

-- 示例查询：
-- SELECT angle, value, category, description FROM polar_demo ORDER BY id; 