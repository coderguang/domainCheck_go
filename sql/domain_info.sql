/*
Navicat MySQL Data Transfer

Source Server         : 47.107.177.155
Source Server Version : 50505
Source Host           : 47.107.177.155:3306
Source Database       : hacker

Target Server Type    : MYSQL
Target Server Version : 50505
File Encoding         : 65001

Date: 2019-08-12 21:20:51
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for domain_info
-- ----------------------------
DROP TABLE IF EXISTS `domain_info`;
CREATE TABLE `domain_info` (
  `name` varchar(255) NOT NULL DEFAULT '',
  `zone` varchar(255) NOT NULL DEFAULT '',
  `name_length` int(11) NOT NULL,
  `status` int(255) NOT NULL,
  `create_dt_str` varchar(255) NOT NULL,
  `update_dt_str` varchar(255) NOT NULL,
  `expiry_dt_str` varchar(255) NOT NULL,
  `create_dt` datetime DEFAULT NULL,
  `update_dt` datetime DEFAULT NULL,
  `expiry_dt` datetime DEFAULT NULL,
  PRIMARY KEY (`name`,`zone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
