-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
    name VARCHAR(100) NOT NULL COMMENT '用户名',
    email VARCHAR(100) NOT NULL UNIQUE COMMENT '邮箱',
    age INT COMMENT '年龄',
    avatar VARCHAR(255) COMMENT '头像URL',
    status TINYINT DEFAULT 1 COMMENT '状态：1-正常，0-禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_email (email),
    INDEX idx_created_at (created_at)
) COMMENT='用户表';

-- 分类表
CREATE TABLE IF NOT EXISTS categories (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '分类ID',
    name VARCHAR(50) NOT NULL COMMENT '分类名称',
    description TEXT COMMENT '分类描述',
    parent_id INT DEFAULT 0 COMMENT '父分类ID，0表示顶级分类',
    sort_order INT DEFAULT 0 COMMENT '排序值',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否激活',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_parent_id (parent_id),
    INDEX idx_sort_order (sort_order)
) COMMENT='分类表';

-- 文章表
CREATE TABLE IF NOT EXISTS articles (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '文章ID',
    title VARCHAR(200) NOT NULL COMMENT '文章标题',
    content LONGTEXT COMMENT '文章内容',
    summary VARCHAR(500) COMMENT '文章摘要',
    author_id BIGINT NOT NULL COMMENT '作者ID',
    category_id INT COMMENT '分类ID',
    view_count INT DEFAULT 0 COMMENT '阅读次数',
    like_count INT DEFAULT 0 COMMENT '点赞数',
    status ENUM('draft', 'published', 'archived') DEFAULT 'draft' COMMENT '文章状态',
    published_at TIMESTAMP NULL COMMENT '发布时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL,
    INDEX idx_author_id (author_id),
    INDEX idx_category_id (category_id),
    INDEX idx_status (status),
    INDEX idx_published_at (published_at)
) COMMENT='文章表';

-- 标签表
CREATE TABLE IF NOT EXISTS tags (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '标签ID',
    name VARCHAR(50) NOT NULL UNIQUE COMMENT '标签名称',
    color VARCHAR(7) DEFAULT '#007bff' COMMENT '标签颜色',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'
) COMMENT='标签表';

-- 文章标签关联表（多对多关系）
CREATE TABLE IF NOT EXISTS article_tags (
    article_id BIGINT NOT NULL COMMENT '文章ID',
    tag_id INT NOT NULL COMMENT '标签ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (article_id, tag_id),
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
) COMMENT='文章标签关联表';

-- 插入示例数据
INSERT INTO categories (name, description, parent_id, sort_order) VALUES
('技术', '技术相关文章', 0, 1),
('生活', '生活感悟文章', 0, 2),
('Go语言', 'Go语言相关技术文章', 1, 1),
('前端技术', '前端开发相关文章', 1, 2);

INSERT INTO tags (name, color) VALUES
('Go', '#00ADD8'),
('GORM', '#FF6B6B'),
('数据库', '#4ECDC4'),
('教程', '#45B7D1');

INSERT INTO articles (title, content, summary, author_id, category_id, status, published_at) VALUES
('Gorm Gen 入门教程', '这是一篇关于 Gorm Gen 的详细教程...', '详细介绍了 Gorm Gen 的基础使用方法', 1, 3, 'published', NOW()),
('Go 语言并发编程', '本文介绍 Go 语言的并发编程模式...', 'Go 语言并发编程的最佳实践', 1, 3, 'published', NOW());

INSERT INTO article_tags (article_id, tag_id) VALUES
(1, 1), (1, 2), (1, 4),
(2, 1), (2, 4);