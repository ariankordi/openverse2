-- Adminer 4.8.1 MySQL 10.4.33-MariaDB dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

DROP TABLE IF EXISTS `blocks`;
CREATE TABLE `blocks` (
  `block_id` int(11) NOT NULL AUTO_INCREMENT,
  `block_to` int(11) NOT NULL,
  `block_by` int(11) NOT NULL,
  PRIMARY KEY (`block_id`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;


SET NAMES utf8mb4;

DROP TABLE IF EXISTS `communities`;
CREATE TABLE `communities` (
  `community_id` int(11) NOT NULL AUTO_INCREMENT,
  `community_title` int(11) NOT NULL,
  `community_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci NOT NULL,
  `community_description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `community_icon` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `community_banner` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `community_platform` int(1) NOT NULL COMMENT '0 = 3DS, 1 = Wii U, 2 = 3DS/Wii U, 3 or more = N/A',
  `community_type` int(1) NOT NULL,
  `community_perms` int(1) NOT NULL,
  PRIMARY KEY (`community_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf32 COLLATE=utf32_unicode_520_ci;


DROP TABLE IF EXISTS `empathies`;
CREATE TABLE `empathies` (
  `yeah_id` int(11) NOT NULL AUTO_INCREMENT,
  `yeah_type` int(1) NOT NULL,
  `yeah_to` int(11) NOT NULL,
  `yeah_by` int(11) NOT NULL,
  PRIMARY KEY (`yeah_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


DROP TABLE IF EXISTS `favorites`;
CREATE TABLE `favorites` (
  `favorite_id` int(11) NOT NULL AUTO_INCREMENT,
  `favorite_by` int(11) NOT NULL,
  `favorite_to` int(11) NOT NULL,
  PRIMARY KEY (`favorite_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


DROP TABLE IF EXISTS `follows`;
CREATE TABLE `follows` (
  `follow_id` int(11) NOT NULL AUTO_INCREMENT,
  `follow_to` int(11) NOT NULL,
  `follow_by` int(11) DEFAULT NULL,
  PRIMARY KEY (`follow_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


DROP TABLE IF EXISTS `friendships`;
CREATE TABLE `friendships` (
  `friend_id` int(11) NOT NULL AUTO_INCREMENT,
  `friend_date` datetime NOT NULL,
  `friend_to` int(11) NOT NULL,
  `friend_by` int(11) NOT NULL,
  PRIMARY KEY (`friend_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


DROP TABLE IF EXISTS `friend_requests`;
CREATE TABLE `friend_requests` (
  `request_id` int(11) NOT NULL AUTO_INCREMENT,
  `request_to` int(11) NOT NULL,
  `request_by` int(11) NOT NULL,
  `request_date` datetime NOT NULL,
  PRIMARY KEY (`request_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


DROP TABLE IF EXISTS `login_tokens`;
CREATE TABLE `login_tokens` (
  `token_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `token_for` int(11) NOT NULL,
  `token_created` datetime NOT NULL DEFAULT current_timestamp(),
  `token_status` int(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (`token_id`),
  KEY `token_for` (`token_for`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;


DROP TABLE IF EXISTS `messages`;
CREATE TABLE `messages` (
  `message_id` int(11) NOT NULL AUTO_INCREMENT,
  `message_to` int(11) NOT NULL,
  `message_by` int(11) NOT NULL,
  `message_feeling_id` int(1) NOT NULL,
  `message_content` text NOT NULL,
  `message_is_spoiler` int(1) NOT NULL,
  `message_screenshot` varchar(1024) NOT NULL,
  PRIMARY KEY (`message_id`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;


DROP TABLE IF EXISTS `notifications`;
CREATE TABLE `notifications` (
  `notif_id` int(11) NOT NULL AUTO_INCREMENT,
  `notif_type` int(1) NOT NULL,
  `notif_to` int(11) NOT NULL,
  `notif_by` int(11) NOT NULL,
  `notif_by_others` text DEFAULT NULL,
  `notif_topic` int(11) NOT NULL,
  `notif_date` datetime NOT NULL DEFAULT current_timestamp(),
  `notif_read` int(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (`notif_id`),
  KEY `notif_to` (`notif_to`),
  KEY `notif_read` (`notif_read`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


DROP TABLE IF EXISTS `posts`;
CREATE TABLE `posts` (
  `post_id` int(11) NOT NULL AUTO_INCREMENT,
  `post_feeling_id` int(1) NOT NULL,
  `post_content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `post_screenshot` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `post_drawing` text DEFAULT NULL,
  `post_url` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `post_is_spoiler` int(1) DEFAULT 0,
  `post_has_title_depended_value` int(1) NOT NULL DEFAULT 0,
  `post_date` datetime DEFAULT current_timestamp(),
  `post_community` int(11) NOT NULL,
  `post_by` int(11) NOT NULL,
  `post_status` int(1) NOT NULL DEFAULT 0,
  `post_edited` datetime DEFAULT NULL,
  `post_content_before` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `post_screenshot_before` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  PRIMARY KEY (`post_id`),
  KEY `post_community` (`post_community`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


DROP TABLE IF EXISTS `post_reports`;
CREATE TABLE `post_reports` (
  `report_id` int(11) NOT NULL AUTO_INCREMENT,
  `report_to` int(11) NOT NULL,
  `report_by` int(11) NOT NULL,
  `report_type` int(1) NOT NULL,
  `report_body` tinytext CHARACTER SET utf8 COLLATE utf8_unicode_520_ci NOT NULL,
  `report_date` datetime NOT NULL,
  PRIMARY KEY (`report_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;


DROP TABLE IF EXISTS `post_yeahs`;
CREATE TABLE `post_yeahs` (
  `yeah_id` int(11) NOT NULL AUTO_INCREMENT,
  `yeah_post` int(11) NOT NULL,
  `yeah_by` int(11) NOT NULL,
  `yeah_date` datetime NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`yeah_id`),
  KEY `yeah_post` (`yeah_post`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


DROP TABLE IF EXISTS `replies`;
CREATE TABLE `replies` (
  `reply_id` int(11) NOT NULL AUTO_INCREMENT,
  `reply_to` int(11) NOT NULL,
  `reply_content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `reply_feeling_id` int(1) NOT NULL,
  `reply_screenshot` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `reply_drawing` text DEFAULT NULL,
  `reply_url` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `reply_date` datetime NOT NULL DEFAULT current_timestamp(),
  `reply_edited` datetime DEFAULT NULL,
  `reply_content_before` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `reply_screenshot_before` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `reply_is_spoiler` int(1) NOT NULL DEFAULT 0,
  `reply_by` int(11) NOT NULL,
  `reply_status` int(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (`reply_id`),
  KEY `reply_to` (`reply_to`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


DROP TABLE IF EXISTS `reply_reports`;
CREATE TABLE `reply_reports` (
  `rreport_id` int(11) NOT NULL AUTO_INCREMENT,
  `rreport_to` int(11) NOT NULL,
  `rreport_by` int(11) NOT NULL,
  `rreport_type` int(1) NOT NULL,
  `rreport_body` tinytext NOT NULL,
  `rreport_date` datetime NOT NULL,
  PRIMARY KEY (`rreport_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;


DROP TABLE IF EXISTS `reply_yeahs`;
CREATE TABLE `reply_yeahs` (
  `ryeah_id` int(11) NOT NULL AUTO_INCREMENT,
  `ryeah_reply` int(11) NOT NULL,
  `ryeah_by` int(11) NOT NULL,
  `ryeah_date` datetime NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`ryeah_id`),
  KEY `ryeah_reply` (`ryeah_reply`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


DROP TABLE IF EXISTS `titles`;
CREATE TABLE `titles` (
  `title_id` int(11) NOT NULL AUTO_INCREMENT,
  `title_type` int(1) NOT NULL,
  `title_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci NOT NULL,
  `title_icon` varchar(255) NOT NULL,
  `title_banner` varchar(255) NOT NULL,
  `title_platform` int(1) NOT NULL COMMENT '0 = 3DS, 1 = Wii U, 2 = 3DS/Wii U, 3 or more = N/A',
  PRIMARY KEY (`title_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `user_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci NOT NULL,
  `user_pid` int(11) NOT NULL AUTO_INCREMENT,
  `user_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `user_pass` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `user_email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `user_date` datetime NOT NULL DEFAULT current_timestamp(),
  `user_rank` int(1) NOT NULL DEFAULT 0,
  `user_avatar` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `user_profile_comment` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `user_country` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `user_birthday` date DEFAULT NULL,
  `user_website` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `user_skill` int(1) NOT NULL DEFAULT 0,
  `user_systems` int(1) NOT NULL DEFAULT 0,
  `user_favorite_post` int(11) DEFAULT NULL,
  `user_favorite_post_type` int(1) DEFAULT 0,
  `user_nnid` varchar(16) NOT NULL DEFAULT '',
  `user_ip` varchar(255) NOT NULL,
  `user_code` varchar(4) NOT NULL DEFAULT '',
  `user_email_confirmed` int(1) NOT NULL DEFAULT 0,
  `user_relationship_visibility` int(1) NOT NULL DEFAULT 1,
  PRIMARY KEY (`user_pid`),
  UNIQUE KEY `user_pid` (`user_pid`),
  UNIQUE KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- 2024-05-08 00:27:33
