/*
Source Server         : 10.10.10.10
Source Server Version : 50556
Source Host           : 10.10.10.10:3306
Source Database       : mana

Target Server Type    : MYSQL
Target Server Version : 50556
File Encoding         : 65001

Date:
*/

CREATE DATABASE /*!32312 IF NOT EXISTS*/ mana /*!40100 DEFAULT CHARACTER SET utf8mb4 */;

USE mana;

SET FOREIGN_KEY_CHECKS=0;
-- ----------------------------
-- Table structure for user
-- ----------------------------
-- DROP TABLE IF EXISTS user;
CREATE TABLE user (
  ID int(255) NOT NULL AUTO_INCREMENT,
  USERID varchar(255) NOT NULL COMMENT '用户id',
  USERNAME varchar(255) NOT NULL COMMENT '用户名',
  NICKNAME varchar(255) DEFAULT NULL COMMENT '昵称',
  ROLE int(25) NOT NULL COMMENT '角色',
  PASSWD varchar(255) NOT NULL COMMENT '密码',
  EXPIRES varchar(255) DEFAULT NULL COMMENT '密码过期时间',
  INACTIVE int(25) NOT NULL COMMENT '用户状态',
  CREATETIME datetime DEFAULT NULL COMMENT '创建时间',
  UPDATETIME datetime DEFAULT NULL COMMENT '最近一次密码修改时间',
  PRIMARY KEY (ID)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '用户表';

-- ----------------------------
-- Table structure for user_center
-- ----------------------------
-- DROP TABLE IF EXISTS user_center;
CREATE TABLE user_center (
  ID int(255) NOT NULL AUTO_INCREMENT,
  USERID varchar(255) NOT NULL COMMENT '用户id',
  USERNAME varchar(255) NOT NULL COMMENT '用户名',
  NICKNAME varchar(255) DEFAULT NULL COMMENT '昵称',
  MOBILE varchar(15) DEFAULT '' COMMENT '手机',
  EMAIL varchar(255) DEFAULT '' COMMENT '邮箱',
  DESCRIBES varchar(255) DEFAULT '' COMMENT '描述说明',
  PICTURE varchar(255) DEFAULT '' COMMENT '头像',
  CREATETIME datetime DEFAULT NULL COMMENT '创建时间',
  UPDATETIME datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (ID),
  KEY USERID (USERID),
  KEY USERNAME (USERNAME)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户中心表';

-- ----------------------------
-- Table structure for kube_config
-- ----------------------------
-- DROP TABLE IF EXISTS kube_config;
CREATE TABLE kube_config (
  ID int(255) NOT NULL AUTO_INCREMENT,
  USERID varchar(255) NOT NULL COMMENT '用户id',
  CLUSTER_ALIAS varchar(255) NOT NULL COMMENT '集群别名',
  CLUSTER_USER varchar(255) NOT NULL COMMENT '集群权限用户',
  CURRENT_CONTEXT varchar(255) NOT NULL COMMENT '上下文',
  SERVER varchar(255) NOT NULL COMMENT '集群地址',
  CREATION_TIME datetime DEFAULT NULL COMMENT '创建时间',
  STATUS int(11) NOT NULL COMMENT '配置文件配置状态',
  CERTIFICATE_AUTHORITY_DATA text NOT NULL COMMENT 'CA证书',
  CLIENT_CERTIFICATE_DATA text NOT NULL COMMENT '用户证书',
  CLIENT_KEY_DATA text NOT NULL COMMENT '用户证书私钥',
  PRIMARY KEY (ID)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='集群配置表';

-- ----------------------------
-- Table structure for resource_links
-- ----------------------------
-- DROP TABLE IF EXISTS resource_links;
CREATE TABLE resource_links (
  ID int(11) NOT NULL AUTO_INCREMENT,
  USERID varchar(255) NOT NULL COMMENT '用户id',
  LINKNAME varchar(255) NOT NULL COMMENT '链接名',
  LINKURL varchar(255) NOT NULL COMMENT 'url',
  DESCRIBES varchar(255) DEFAULT NULL COMMENT '说明',
  CATEGORY varchar(255) DEFAULT NULL COMMENT '类别',
  ACTIVATE int(11) NOT NULL DEFAULT '1' COMMENT '状态，1 有效，0 无效',
  CREATETIME varchar(255) DEFAULT NULL COMMENT '创建时间',
  UPDATETIME varchar(255) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (ID)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='导航链接';

-- 用户详情表添加索引，确保不能插入相同用户名的用户
ALTER TABLE user_center ADD UNIQUE(USERNAME)
