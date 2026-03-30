-- 确保使用正确的数据库
USE `miniblog_db`;

CREATE TABLE `relation` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `follower_id` bigint NOT NULL COMMENT '粉丝ID (谁在关注)',
  `following_id` bigint NOT NULL COMMENT '被关注者ID (关注了谁)',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '关注时间',
  PRIMARY KEY (`id`),
  -- 联合唯一索引：防止 A 关注 B 多次
  UNIQUE KEY `idx_follower_following` (`follower_id`, `following_id`),
  -- 辅助索引：为了快速查询“谁关注了我”
  KEY `idx_following` (`following_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户关注关系表';