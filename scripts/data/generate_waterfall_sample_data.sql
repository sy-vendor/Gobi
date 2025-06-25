-- Waterfall Chart Sample Data | 瀑布图示例数据
-- 用于测试瀑布图功能的示例数据 (SQLite兼容)

CREATE TABLE IF NOT EXISTS waterfall_demo (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    step TEXT NOT NULL,         -- 步骤名称
    amount REAL NOT NULL,       -- 变化值
    type TEXT NOT NULL,         -- base/increase/decrease
    description TEXT,           -- 步骤说明
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 示例数据：利润拆解
INSERT INTO waterfall_demo (step, amount, type, description) VALUES
('期初余额', 1000, 'base', '年初资金'),
('主营业务收入', 2000, 'increase', '主营业务带来的收入'),
('其他收入', 500, 'increase', '其他来源收入'),
('运营成本', -1200, 'decrease', '日常运营支出'),
('税费', -300, 'decrease', '税收及附加'),
('净利润', 2000, 'base', '年末净利润');

-- 示例查询：
-- SELECT step, amount, type, description FROM waterfall_demo ORDER BY id; 