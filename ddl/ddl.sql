-- bigintじゃなくてintを使用した方が良い
  -- パフォーマンスが良い
CREATE TABLE `account` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL UNIQUE,
  `password_hash` varchar(255) NOT NULL,
  `display_name` varchar(255),
  `avatar` text,
  `header` text,
  `note` text,
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
);

CREATE TABLE `status` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `account_id` bigint(20) NOT NULL,
  `content` text NOT NULL,
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`account_id`) REFERENCES `account` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE `relationship` (
  `follower_id` bigint(20) NOT NULL,
  `followee_id` bigint(20) NOT NULL,
  PRIMARY KEY (`follower_id`, `followee_id`),
  FOREIGN KEY (`follower_id`) REFERENCES `account` (`id`),
  FOREIGN KEY (`followee_id`) REFERENCES `account` (`id`)
);

CREATE TABLE `media` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `account_id` bigint(20) NOT NULL,
  `type` varchar(255) NOT NULL,
  `url` text NOT NULL,
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`account_id`) REFERENCES `account` (`id`)
);

CREATE TABLE `attachment` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `type` bigint(20) NOT NULL,
  `url` text NOT NULL,
  `description` datetime NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`account_id`) REFERENCES `account` (`id`)
);

CREATE TABLE `attachmentBuilding` (
  `attachment_id` bigint(20) NOT NULL,
  `status_id` bigint(20) NOT NULL,
  PRIMARY KEY (`attachment_id`, `status_id`),
  FOREIGN KEY (`attachment_id`) REFERENCES `attachment` (`id`),
  FOREIGN KEY (`status_id`) REFERENCES `status` (`id`)
);