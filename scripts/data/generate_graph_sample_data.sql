-- 关系图/力导向图示例数据
-- 用于测试Graph/Network/Force-directed图表

CREATE TABLE IF NOT EXISTS graph_nodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    group_id INTEGER,
    value REAL,
    category TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS graph_edges (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source INTEGER NOT NULL,
    target INTEGER NOT NULL,
    weight REAL,
    relation TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 节点示例（用户、设备、组织、地点等）
INSERT INTO graph_nodes (id, name, group_id, value, category) VALUES
(1, 'Alice', 1, 10, 'Person'),
(2, 'Bob', 1, 8, 'Person'),
(3, 'Carol', 2, 6, 'Person'),
(4, 'David', 2, 7, 'Person'),
(5, 'Server', 3, 12, 'Device'),
(6, 'Database', 3, 15, 'Device'),
(7, 'Router', 3, 9, 'Device'),
(8, 'Company', 4, 20, 'Organization'),
(9, 'Office', 4, 13, 'Location'),
(10, 'Internet', 5, 30, 'Network'),
(11, 'Eve', 1, 5, 'Person'),
(12, 'Mall', 4, 11, 'Location');

-- 边示例（多种关系/权重）
INSERT INTO graph_edges (source, target, weight, relation) VALUES
(1, 2, 1, 'friend'),
(2, 3, 1, 'colleague'),
(1, 4, 0.5, 'colleague'),
(3, 4, 2, 'friend'),
(1, 5, 2, 'access'),
(2, 5, 1, 'access'),
(4, 6, 1.5, 'query'),
(5, 6, 3, 'connects'),
(7, 5, 2, 'routes'),
(7, 6, 2, 'routes'),
(5, 10, 2.5, 'internet'),
(6, 10, 2.5, 'internet'),
(7, 10, 2.5, 'internet'),
(1, 8, 1, 'member'),
(2, 8, 1, 'member'),
(3, 8, 1, 'member'),
(4, 8, 1, 'member'),
(8, 9, 1, 'located_at'),
(8, 12, 0.8, 'branch'),
(9, 12, 0.5, 'nearby'),
(11, 1, 0.7, 'friend'),
(11, 2, 0.6, 'friend'),
(11, 12, 0.9, 'visit'),
(3, 7, 0.5, 'admin'),
(4, 7, 0.5, 'admin'); 