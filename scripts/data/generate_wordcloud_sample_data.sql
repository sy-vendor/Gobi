-- 词云图表示例数据
-- 用于展示文本数据中词语频率的可视化

-- 创建社交媒体热门话题数据表
CREATE TABLE IF NOT EXISTS social_media_topics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    topic TEXT NOT NULL,
    frequency INTEGER NOT NULL,
    category TEXT NOT NULL,
    sentiment TEXT NOT NULL, -- positive, negative, neutral
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 插入社交媒体热门话题数据
INSERT INTO social_media_topics (topic, frequency, category, sentiment) VALUES
-- 科技类话题
('人工智能', 1250, '科技', 'positive'),
('机器学习', 980, '科技', 'positive'),
('深度学习', 750, '科技', 'positive'),
('区块链', 680, '科技', 'neutral'),
('云计算', 920, '科技', 'positive'),
('大数据', 1100, '科技', 'positive'),
('物联网', 580, '科技', 'positive'),
('5G技术', 820, '科技', 'positive'),
('虚拟现实', 450, '科技', 'positive'),
('增强现实', 380, '科技', 'positive'),
('量子计算', 320, '科技', 'positive'),
('边缘计算', 280, '科技', 'neutral'),

-- 商业类话题
('数字化转型', 890, '商业', 'positive'),
('可持续发展', 760, '商业', 'positive'),
('绿色能源', 650, '商业', 'positive'),
('电动汽车', 720, '商业', 'positive'),
('远程办公', 950, '商业', 'positive'),
('混合办公', 680, '商业', 'positive'),
('供应链', 540, '商业', 'neutral'),
('电子商务', 1100, '商业', 'positive'),
('移动支付', 820, '商业', 'positive'),
('数字营销', 780, '商业', 'positive'),
('客户体验', 650, '商业', 'positive'),
('品牌建设', 580, '商业', 'positive'),

-- 健康类话题
('心理健康', 1200, '健康', 'positive'),
('健身运动', 980, '健康', 'positive'),
('营养饮食', 850, '健康', 'positive'),
('睡眠质量', 720, '健康', 'positive'),
('压力管理', 680, '健康', 'positive'),
('冥想放松', 450, '健康', 'positive'),
('瑜伽练习', 380, '健康', 'positive'),
('户外运动', 650, '健康', 'positive'),
('慢性疾病', 420, '健康', 'negative'),
('预防保健', 580, '健康', 'positive'),
('中医养生', 320, '健康', 'positive'),
('现代医学', 750, '健康', 'positive'),

-- 教育类话题
('在线教育', 1100, '教育', 'positive'),
('终身学习', 850, '教育', 'positive'),
('技能提升', 920, '教育', 'positive'),
('职业发展', 780, '教育', 'positive'),
('编程学习', 680, '教育', 'positive'),
('语言学习', 650, '教育', 'positive'),
('创新思维', 580, '教育', 'positive'),
('批判性思维', 420, '教育', 'positive'),
('团队协作', 550, '教育', 'positive'),
('领导力', 480, '教育', 'positive'),
('创业精神', 380, '教育', 'positive'),
('学术研究', 320, '教育', 'positive'),

-- 娱乐类话题
('电影推荐', 950, '娱乐', 'positive'),
('音乐分享', 880, '娱乐', 'positive'),
('游戏体验', 820, '娱乐', 'positive'),
('旅游攻略', 780, '娱乐', 'positive'),
('美食探店', 720, '娱乐', 'positive'),
('时尚穿搭', 680, '娱乐', 'positive'),
('宠物日常', 650, '娱乐', 'positive'),
('摄影技巧', 580, '娱乐', 'positive'),
('读书分享', 520, '娱乐', 'positive'),
('综艺节目', 480, '娱乐', 'neutral'),
('明星动态', 420, '娱乐', 'neutral'),
('网红打卡', 380, '娱乐', 'neutral');

-- 创建新闻关键词数据表
CREATE TABLE IF NOT EXISTS news_keywords (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    keyword TEXT NOT NULL,
    frequency INTEGER NOT NULL,
    source TEXT NOT NULL,
    date DATE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 插入新闻关键词数据
INSERT INTO news_keywords (keyword, frequency, source, date) VALUES
-- 经济类关键词
('经济增长', 156, '财经新闻', '2024-01-15'),
('通货膨胀', 142, '财经新闻', '2024-01-15'),
('货币政策', 128, '财经新闻', '2024-01-15'),
('股市行情', 135, '财经新闻', '2024-01-15'),
('房地产', 118, '财经新闻', '2024-01-15'),
('消费升级', 95, '财经新闻', '2024-01-15'),
('数字化转型', 112, '财经新闻', '2024-01-15'),
('绿色发展', 88, '财经新闻', '2024-01-15'),
('科技创新', 125, '财经新闻', '2024-01-15'),
('国际贸易', 98, '财经新闻', '2024-01-15'),

-- 科技类关键词
('人工智能', 189, '科技新闻', '2024-01-15'),
('芯片技术', 156, '科技新闻', '2024-01-15'),
('新能源车', 134, '科技新闻', '2024-01-15'),
('5G网络', 145, '科技新闻', '2024-01-15'),
('云计算', 167, '科技新闻', '2024-01-15'),
('大数据', 178, '科技新闻', '2024-01-15'),
('区块链', 123, '科技新闻', '2024-01-15'),
('元宇宙', 98, '科技新闻', '2024-01-15'),
('量子计算', 87, '科技新闻', '2024-01-15'),
('网络安全', 134, '科技新闻', '2024-01-15'),

-- 社会类关键词
('教育改革', 145, '社会新闻', '2024-01-15'),
('医疗健康', 167, '社会新闻', '2024-01-15'),
('环境保护', 134, '社会新闻', '2024-01-15'),
('文化传承', 98, '社会新闻', '2024-01-15'),
('乡村振兴', 112, '社会新闻', '2024-01-15'),
('城市发展', 123, '社会新闻', '2024-01-15'),
('社会保障', 145, '社会新闻', '2024-01-15'),
('就业创业', 156, '社会新闻', '2024-01-15'),
('养老服务', 89, '社会新闻', '2024-01-15'),
('儿童教育', 134, '社会新闻', '2024-01-15');

-- 创建产品评论关键词数据表
CREATE TABLE IF NOT EXISTS product_review_keywords (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    keyword TEXT NOT NULL,
    frequency INTEGER NOT NULL,
    product_category TEXT NOT NULL,
    sentiment TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 插入产品评论关键词数据
INSERT INTO product_review_keywords (keyword, frequency, product_category, sentiment) VALUES
-- 电子产品正面评价
('性能强劲', 234, '电子产品', 'positive'),
('外观精美', 198, '电子产品', 'positive'),
('操作流畅', 187, '电子产品', 'positive'),
('续航持久', 165, '电子产品', 'positive'),
('拍照清晰', 145, '电子产品', 'positive'),
('音质出色', 123, '电子产品', 'positive'),
('性价比高', 178, '电子产品', 'positive'),
('做工精良', 156, '电子产品', 'positive'),
('功能丰富', 134, '电子产品', 'positive'),
('散热良好', 98, '电子产品', 'positive'),

-- 电子产品负面评价
('价格偏高', 89, '电子产品', 'negative'),
('发热严重', 76, '电子产品', 'negative'),
('续航不足', 67, '电子产品', 'negative'),
('系统卡顿', 54, '电子产品', 'negative'),
('拍照模糊', 43, '电子产品', 'negative'),
('音质一般', 38, '电子产品', 'negative'),
('做工粗糙', 45, '电子产品', 'negative'),
('功能单一', 32, '电子产品', 'negative'),
('售后服务差', 28, '电子产品', 'negative'),
('包装简陋', 19, '电子产品', 'negative'),

-- 服装类正面评价
('面料舒适', 167, '服装', 'positive'),
('版型合身', 145, '服装', 'positive'),
('颜色漂亮', 134, '服装', 'positive'),
('做工精细', 123, '服装', 'positive'),
('款式时尚', 156, '服装', 'positive'),
('尺码标准', 98, '服装', 'positive'),
('质量不错', 112, '服装', 'positive'),
('穿着舒适', 134, '服装', 'positive'),
('百搭实用', 89, '服装', 'positive'),
('性价比高', 145, '服装', 'positive'),

-- 服装类负面评价
('尺码偏小', 67, '服装', 'negative'),
('面料粗糙', 54, '服装', 'negative'),
('颜色偏差', 43, '服装', 'negative'),
('做工粗糙', 38, '服装', 'negative'),
('版型不合', 45, '服装', 'negative'),
('容易起球', 32, '服装', 'negative'),
('褪色严重', 28, '服装', 'negative'),
('线头较多', 19, '服装', 'negative'),
('价格虚高', 34, '服装', 'negative'),
('款式过时', 23, '服装', 'negative');

-- 创建索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_social_media_topics_category ON social_media_topics(category);
CREATE INDEX IF NOT EXISTS idx_social_media_topics_sentiment ON social_media_topics(sentiment);
CREATE INDEX IF NOT EXISTS idx_news_keywords_source ON news_keywords(source);
CREATE INDEX IF NOT EXISTS idx_news_keywords_date ON news_keywords(date);
CREATE INDEX IF NOT EXISTS idx_product_review_keywords_category ON product_review_keywords(product_category);
CREATE INDEX IF NOT EXISTS idx_product_review_keywords_sentiment ON product_review_keywords(sentiment);

-- 显示数据统计
SELECT 'social_media_topics' as table_name, COUNT(*) as record_count FROM social_media_topics
UNION ALL
SELECT 'news_keywords' as table_name, COUNT(*) as record_count FROM news_keywords
UNION ALL
SELECT 'product_review_keywords' as table_name, COUNT(*) as record_count FROM product_review_keywords;

-- 示例查询：获取社交媒体热门话题（按频率排序）
SELECT topic, frequency, category, sentiment 
FROM social_media_topics 
ORDER BY frequency DESC 
LIMIT 20;

-- 示例查询：获取新闻关键词（按频率排序）
SELECT keyword, frequency, source, date 
FROM news_keywords 
ORDER BY frequency DESC 
LIMIT 20;

-- 示例查询：获取产品评论关键词（按频率排序）
SELECT keyword, frequency, product_category, sentiment 
FROM product_review_keywords 
ORDER BY frequency DESC 
LIMIT 20; 