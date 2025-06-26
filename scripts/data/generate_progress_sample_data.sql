-- Progress/Circular Progress Chart Sample Data | 进度条/环形进度图示例数据
-- 用于测试进度图功能的示例数据 (SQLite兼容)

CREATE TABLE IF NOT EXISTS progress_demo (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,         -- 任务/项目/目标名称
    value REAL NOT NULL,        -- 进度百分比（0-100）
    category TEXT,              -- 分类
    color TEXT,                 -- 颜色（可选）
    description TEXT,           -- 说明
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 示例数据1：任务进度
INSERT INTO progress_demo (name, value, category, color, description) VALUES
('需求分析', 100, '项目A', '#1890ff', '需求分析已完成'),
('设计', 80, '项目A', '#2fc25b', '设计阶段进行中'),
('开发', 60, '项目A', '#facc14', '开发阶段进行中'),
('测试', 30, '项目A', '#f5222d', '测试阶段尚未开始'),
('上线', 0, '项目A', '#722ed1', '尚未上线');

-- 示例数据2：项目完成率
INSERT INTO progress_demo (name, value, category, color, description) VALUES
('项目A', 70, '年度项目', '#1890ff', '项目A已完成70%'),
('项目B', 45, '年度项目', '#2fc25b', '项目B已完成45%'),
('项目C', 90, '年度项目', '#facc14', '项目C已完成90%');

-- 示例数据3：销售目标完成率
INSERT INTO progress_demo (name, value, category, color, description) VALUES
('Q1销售目标', 85, '销售', '#fa8c16', '第一季度销售目标完成85%'),
('Q2销售目标', 60, '销售', '#13c2c2', '第二季度销售目标完成60%'),
('Q3销售目标', 40, '销售', '#eb2f96', '第三季度销售目标完成40%'),
('Q4销售目标', 20, '销售', '#a0d911', '第四季度销售目标完成20%');

-- 示例查询：
-- SELECT name, value, category, color, description FROM progress_demo ORDER BY id; 