-- TreeMap & Sunburst Sample Data | 树状图/旭日图示例数据
-- 用于测试树状图/旭日图功能的层级结构数据 (SQLite兼容版本)

-- 公司部门层级表（用于树状图/旭日图）
CREATE TABLE IF NOT EXISTS department_hierarchy (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    parent_id INTEGER,
    name TEXT NOT NULL,
    value REAL NOT NULL,
    category TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 示例数据：公司部门层级
INSERT INTO department_hierarchy (id, parent_id, name, value, category) VALUES
(1, NULL, '公司', 100, '集团'),
(2, 1, '研发中心', 40, '技术'),
(3, 1, '市场部', 30, '业务'),
(4, 1, '人力资源部', 10, '职能'),
(5, 2, '后端组', 20, '技术'),
(6, 2, '前端组', 10, '技术'),
(7, 2, '测试组', 10, '技术'),
(8, 3, '国内市场', 20, '业务'),
(9, 3, '国际市场', 10, '业务');

-- 产品分类层级表（用于旭日图/树状图）
CREATE TABLE IF NOT EXISTS product_hierarchy (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    parent_id INTEGER,
    name TEXT NOT NULL,
    value REAL NOT NULL,
    category TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 示例数据：产品分类层级
INSERT INTO product_hierarchy (id, parent_id, name, value, category) VALUES
(1, NULL, '全部产品', 200, '总类'),
(2, 1, '电子产品', 80, '电子'),
(3, 1, '服装', 70, '服饰'),
(4, 1, '食品', 50, '食品'),
(5, 2, '手机', 40, '电子'),
(6, 2, '电脑', 40, '电子'),
(7, 3, '男装', 30, '服饰'),
(8, 3, '女装', 40, '服饰'),
(9, 4, '零食', 20, '食品'),
(10, 4, '饮料', 30, '食品'); 