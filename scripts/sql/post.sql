CREATE TABLE `post` (
  `id` bigint NOT NULL COMMENT '全局唯一博文ID',
  `user_id` bigint NOT NULL COMMENT '作者ID (未来的分库分表 Sharding Key)',
  `content` text NOT NULL COMMENT '博文内容',
  `images` json NULL COMMENT '博文图片URL数组', -- ✅ 新增这一行
  `like_count` int NOT NULL DEFAULT '0' COMMENT '点赞数',
  `comment_count` int NOT NULL DEFAULT '0' COMMENT '评论数',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发布时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`) COMMENT '按作者查询博文的辅助索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='微型博客内容表';