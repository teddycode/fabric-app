

SET NAMES utf8;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
--  用户表结构
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id`          int(10) unsigned NOT NULL AUTO_INCREMENT,
  `username`    varchar(50) DEFAULT '' COMMENT '昵称',
  `email`       varchar(50) DEFAULT '' COMMENT '邮箱',
  `role`        int(1)  DEFAULT 0 NOT NULL COMMENT '管理员/普通用户',
  `phone`       varchar(50) DEFAULT '' COMMENT '电话',
  `password`    varchar(50) DEFAULT '' COMMENT '密码',
  `balance`      float(6,2) DEFAULT '0' COMMENT '余额',
  `secret`      varchar(20) NOT NULL DEFAULT '' COMMENT 'jwt动态密钥，注销、修改密码时候改变',
  `created_on`  int(11) unsigned DEFAULT NULL COMMENT '创建时间',
  `modified_on` int(11) unsigned DEFAULT NULL COMMENT '更新时间',
  `deleted_on`  int(11) unsigned DEFAULT '0' COMMENT '删除时间戳',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='用户管理';

-- ----------------------------
--  充值记录表结构
-- ----------------------------
DROP TABLE IF EXISTS `income`;
CREATE TABLE `income` (
    `id`             int(10) unsigned NOT NULL AUTO_INCREMENT,
    `user_id`         int(10) unsigned NOT NULL  COMMENT '用户id',
    `amount`         float NOT NULL COMMENT '数额',
    `method`         varchar(50) DEFAULT '充值方式',
    `created_on`  int(11) unsigned DEFAULT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='余额管理';

-- ----------------------------
--  消费记录表结构
-- ----------------------------
DROP TABLE IF EXISTS `cost`;
CREATE TABLE `cost` (
     `id`             int(10) unsigned NOT NULL AUTO_INCREMENT,
     `user_id`         int(10) unsigned NOT NULL  COMMENT '用户id',
     `amount`         float NOT NULL COMMENT '数额',
     `created_on`     int(11) unsigned DEFAULT NULL COMMENT '创建时间',
     `end_on`         int(11) unsigned DEFAULT NULL COMMENT '结束时间',
     `electric`       int(11) unsigned DEFAULT NULL COMMENT '耗电量',
     PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='余额管理';


-- ----------------------------
--  充值码表结构
-- ----------------------------
DROP TABLE IF EXISTS  `code`;
CREATE  TABLE  `code` (
    `id`             int(10) unsigned NOT NULL AUTO_INCREMENT,
    `amount`         float NOT NULL       COMMENT '数额',
    `code`   varchar(50) NOT NULL COMMENT '充值码',
    `created_by`      varchar(50) NOT NULL COMMENT '创建用户',
    `created_on`     int(11) unsigned DEFAULT NULL COMMENT '创建时间',
    `available`      int(1) unsigned NOT NULL  DEFAULT 0 COMMENT  '是否有效',
    PRIMARY KEY (`id`)
)ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='余额管理';

-- ----------------------------
--  添加用户记录
-- ----------------------------
BEGIN;
INSERT INTO `user` VALUES ('1', 'admin', '105533@qq.com','1','18677331111','123456','0', '', '1', '0', '0');
COMMIT;

-- ----------------------------
--  添加充值记录
-- ----------------------------
BEGIN;
INSERT INTO `income` VALUES ('1', '1', '10','code');
COMMIT;

-- ----------------------------
--  添加消费记录
-- ----------------------------
BEGIN;
INSERT INTO `cost` VALUES ('1', '1', '10','11','10','20');
COMMIT;


-- ----------------------------
--  添加充值码记录
-- ----------------------------
BEGIN;
INSERT INTO `code` VALUES ('1', '100', '1234','admin','123456','1');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
