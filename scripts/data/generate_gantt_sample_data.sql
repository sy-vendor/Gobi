-- Gantt Chart Sample Data | 甘特图示例数据
-- 用于测试甘特图功能的示例数据 (SQLite兼容)

CREATE TABLE IF NOT EXISTS gantt_demo (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id TEXT NOT NULL,        -- 任务ID
    task_name TEXT NOT NULL,      -- 任务名称
    start_date TEXT NOT NULL,     -- 开始日期 (YYYY-MM-DD)
    end_date TEXT NOT NULL,       -- 结束日期 (YYYY-MM-DD)
    duration INTEGER NOT NULL,    -- 持续时间（天）
    progress INTEGER NOT NULL,    -- 进度百分比 (0-100)
    status TEXT NOT NULL,         -- 状态 (未开始/进行中/已完成/延期)
    assignee TEXT,                -- 负责人
    dependencies TEXT,            -- 依赖任务ID（逗号分隔）
    project TEXT NOT NULL,        -- 项目名称
    priority TEXT,                -- 优先级 (高/中/低)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 示例数据1：软件开发项目
INSERT INTO gantt_demo (task_id, task_name, start_date, end_date, duration, progress, status, assignee, dependencies, project, priority) VALUES
('TASK-001', '需求分析', '2024-01-01', '2024-01-05', 5, 100, '已完成', '张三', NULL, '电商平台开发', '高'),
('TASK-002', '系统设计', '2024-01-06', '2024-01-15', 10, 80, '进行中', '李四', 'TASK-001', '电商平台开发', '高'),
('TASK-003', '数据库设计', '2024-01-08', '2024-01-12', 5, 100, '已完成', '王五', 'TASK-001', '电商平台开发', '中'),
('TASK-004', '前端开发', '2024-01-16', '2024-02-15', 31, 60, '进行中', '赵六', 'TASK-002,TASK-003', '电商平台开发', '高'),
('TASK-005', '后端开发', '2024-01-16', '2024-02-20', 36, 50, '进行中', '钱七', 'TASK-002,TASK-003', '电商平台开发', '高'),
('TASK-006', '系统测试', '2024-02-21', '2024-03-05', 13, 0, '未开始', '孙八', 'TASK-004,TASK-005', '电商平台开发', '中'),
('TASK-007', '用户验收', '2024-03-06', '2024-03-10', 5, 0, '未开始', '周九', 'TASK-006', '电商平台开发', '低'),
('TASK-008', '系统上线', '2024-03-11', '2024-03-15', 5, 0, '未开始', '吴十', 'TASK-007', '电商平台开发', '高');

-- 示例数据2：建筑项目
INSERT INTO gantt_demo (task_id, task_name, start_date, end_date, duration, progress, status, assignee, dependencies, project, priority) VALUES
('BUILD-001', '土地平整', '2024-01-01', '2024-01-10', 10, 100, '已完成', '工程队A', NULL, '商业大厦建设', '高'),
('BUILD-002', '地基施工', '2024-01-11', '2024-02-10', 31, 90, '进行中', '工程队B', 'BUILD-001', '商业大厦建设', '高'),
('BUILD-003', '主体结构', '2024-02-11', '2024-06-10', 121, 40, '进行中', '工程队C', 'BUILD-002', '商业大厦建设', '高'),
('BUILD-004', '外墙装修', '2024-06-11', '2024-08-10', 61, 0, '未开始', '装修队A', 'BUILD-003', '商业大厦建设', '中'),
('BUILD-005', '内部装修', '2024-07-01', '2024-09-30', 92, 0, '未开始', '装修队B', 'BUILD-003', '商业大厦建设', '中'),
('BUILD-006', '设备安装', '2024-08-01', '2024-10-31', 92, 0, '未开始', '设备队A', 'BUILD-004,BUILD-005', '商业大厦建设', '中'),
('BUILD-007', '竣工验收', '2024-11-01', '2024-11-15', 15, 0, '未开始', '验收组', 'BUILD-006', '商业大厦建设', '高');

-- 示例查询：
-- SELECT task_id, task_name, start_date, end_date, duration, progress, status, assignee, dependencies, project, priority FROM gantt_demo ORDER BY project, start_date; 