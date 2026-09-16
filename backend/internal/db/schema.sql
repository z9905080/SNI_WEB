CREATE TABLE `carousel` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '流水號',
  `image` text NOT NULL COMMENT '圖片網址',
  `url` text NOT NULL COMMENT '網址',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb3;

CREATE TABLE `marquee` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '流水號',
  `text` mediumtext NOT NULL COMMENT '文字內容',
  `color` varchar(20) NOT NULL COMMENT '色碼',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb3;

CREATE TABLE `page_content` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '流水號',
  `page_group_id` bigint NOT NULL COMMENT '群組ID',
  `page_name` varchar(50) NOT NULL COMMENT '分頁名稱',
  `html_context` text NOT NULL COMMENT '內頁資料(html)',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

CREATE TABLE `page_group` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '流水號',
  `group_name` varchar(50) NOT NULL COMMENT '群組名稱',
  `page_sort` text NOT NULL COMMENT '內頁順序',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

CREATE TABLE `user` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '流水號',
  `account` varchar(50) NOT NULL COMMENT '帳號',
  `pwd` varchar(255) NOT NULL COMMENT '密碼（bcrypt）',
  `identity` int NOT NULL COMMENT '1.管理者',
  `user_name` varchar(50) NOT NULL COMMENT '暱稱',
  PRIMARY KEY (`id`),
  KEY `account` (`account`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb3;

CREATE TABLE `user_token` (
  `user_id` bigint NOT NULL COMMENT '使用者ID',
  `token` varchar(255) NOT NULL COMMENT '登入 token（SHA-256）',
  `expire_time` datetime NOT NULL COMMENT '過期時間',
  PRIMARY KEY (`user_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb3;

CREATE TABLE `web_config` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '流水號',
  `data_key` varchar(50) NOT NULL COMMENT 'Key',
  `data_value` text NOT NULL COMMENT '值',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
