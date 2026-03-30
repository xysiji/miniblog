USE `miniblog_db`;
CREATE TABLE `like_record` (
  `id` bigint NOT NULL COMMENT '主键ID',
  `post_id` bigint NOT NULL COMMENT '博文ID',
  `user_id` bigint NOT NULL COMMENT '点赞用户ID',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态：1-点赞，0-取消',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_post_user` (`post_id`, `user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `comment` (
  `id` bigint NOT NULL COMMENT '评论ID',
  `post_id` bigint NOT NULL COMMENT '博文ID',
  `root_id` bigint NOT NULL DEFAULT '0' COMMENT '根评论ID (0表示是一级主评论)',
  `parent_id` bigint NOT NULL DEFAULT '0' COMMENT '父评论ID (回复哪条评论)',
  `user_id` bigint NOT NULL COMMENT '评论用户ID',
  `reply_to_user_id` bigint NOT NULL DEFAULT '0' COMMENT '被回复的最终用户ID',
  `content` text NOT NULL COMMENT '评论内容',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态：1-正常，0-已删除',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_post_id` (`post_id`),
  KEY `idx_root_id` (`root_id`) -- 核心：支撑根据主评论一次性拉取所有子评论
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;