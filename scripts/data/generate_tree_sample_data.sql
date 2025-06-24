-- Tree Diagram Sample Data | 树形图示例数据
-- 用于测试分支结构树形图功能的层级结构数据 (SQLite兼容版本)

-- 组织架构树形表（用于分支结构树形图）
CREATE TABLE IF NOT EXISTS org_tree (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    parent_id INTEGER,
    name TEXT NOT NULL,
    position TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 示例数据：公司组织架构树
INSERT INTO org_tree (id, parent_id, name, position) VALUES
(1, NULL, '王总', '总经理'),
(2, 1, '李工', '技术总监'),
(3, 1, '张主管', '市场总监'),
(4, 2, '赵工程师', '后端工程师'),
(5, 2, '钱工程师', '前端工程师'),
(6, 3, '孙专员', '国内市场'),
(7, 3, '周专员', '国际市场'); 