/*
 Navicat Premium Data Transfer

 Source Server         : #1-本地库
 Source Server Type    : MySQL
 Source Server Version : 50726
 Source Host           : localhost:3306
 Source Schema         : go_admin

 Target Server Type    : MySQL
 Target Server Version : 50726
 File Encoding         : 65001

 Date: 14/10/2020 01:15:57
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for admin_article
-- ----------------------------
DROP TABLE IF EXISTS `admin_article`;
CREATE TABLE `admin_article`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `category_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '文章分类ID',
  `describe` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '文章描述',
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '文章内容',
  `status` tinyint(1) UNSIGNED NULL DEFAULT 0 COMMENT '文章状态 1-可用 2-禁用 3-删除',
  `create_time` int(11) UNSIGNED NULL DEFAULT 0 COMMENT '文章创建时间',
  `update_time` int(11) UNSIGNED NULL DEFAULT 0 COMMENT '文章更新时间',
  `tag` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT 'Tags',
  `post_hits` bigint(20) NULL DEFAULT 0 COMMENT '查看数',
  `post_like` bigint(20) NULL DEFAULT 0 COMMENT '点赞数',
  `comment_count` bigint(20) NULL DEFAULT 0 COMMENT '评论数',
  `comment_status` tinyint(3) NOT NULL DEFAULT 1 COMMENT '评论状态;1:允许;0:不允许',
  `more` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '扩展属性,如缩略图;格式为json',
  `source` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '来源',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `category_id`(`category_id`, `status`, `create_time`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_article
-- ----------------------------
INSERT INTO `admin_article` VALUES (1, '文章修改测试', 1, '', '文章修改测试!', 1, 0, 1598887505, '', 0, 0, 0, 0, '', '');
INSERT INTO `admin_article` VALUES (2, '文章修改测试123', 1, '', '文章修改测试!文章修改测试!', 1, 0, 1598887489, '', 0, 0, 0, 0, '', '');

-- ----------------------------
-- Table structure for admin_auth_group
-- ----------------------------
DROP TABLE IF EXISTS `admin_auth_group`;
CREATE TABLE `admin_auth_group`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `title` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户组中文名称',
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '描述信息',
  `status` tinyint(2) NULL DEFAULT 1 COMMENT '用户组状态：为1正常，为0禁用,-1为删除',
  `rules` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '用户组拥有的规则id，多个规则 , 隔开',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_auth_group
-- ----------------------------
INSERT INTO `admin_auth_group` VALUES (1, '测试组', '', 0, '1,2,3,4');

-- ----------------------------
-- Table structure for admin_auth_group_access
-- ----------------------------
DROP TABLE IF EXISTS `admin_auth_group_access`;
CREATE TABLE `admin_auth_group_access`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `group_id` int(11) UNSIGNED NOT NULL COMMENT '用户组id',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `uid`(`uid`, `group_id`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Fixed;

-- ----------------------------
-- Records of admin_auth_group_access
-- ----------------------------
INSERT INTO `admin_auth_group_access` VALUES (5, 1, 1);
INSERT INTO `admin_auth_group_access` VALUES (2, 2, 1);
INSERT INTO `admin_auth_group_access` VALUES (3, 3, 1);
INSERT INTO `admin_auth_group_access` VALUES (4, 4, 1);

-- ----------------------------
-- Table structure for admin_category
-- ----------------------------
DROP TABLE IF EXISTS `admin_category`;
CREATE TABLE `admin_category`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `parent_id` int(10) NULL DEFAULT 0 COMMENT '分类父id',
  `status` tinyint(1) UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态;1:显示;2:隐藏',
  `sort` float(5, 0) NULL DEFAULT 0 COMMENT '排序',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '分类名称',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '分类描述',
  `alias` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '分类别名',
  `list_tpl` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '分类列表模板',
  `one_tpl` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '分类文章页模板',
  `create_time` int(10) UNSIGNED NULL DEFAULT 0 COMMENT '创建时间',
  `icon` char(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '分类图标',
  `thumbnail` int(10) NULL DEFAULT 0 COMMENT '分类封面图',
  `more` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '扩展属性,格式为json',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '栏目表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_category
-- ----------------------------
INSERT INTO `admin_category` VALUES (1, 0, 1, 0, '测试栏目1', '', 'test2', '', '', 1589090707, '', 0, '');

-- ----------------------------
-- Table structure for admin_category_articles
-- ----------------------------
DROP TABLE IF EXISTS `admin_category_articles`;
CREATE TABLE `admin_category_articles`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `articles_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '文章id',
  `category_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '分类id',
  `list_order` float NOT NULL DEFAULT 10000 COMMENT '排序',
  `status` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态,1:发布;0:不发布',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `term_taxonomy_id`(`category_id`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'portal应用 分类文章对应表' ROW_FORMAT = Fixed;

-- ----------------------------
-- Records of admin_category_articles
-- ----------------------------
INSERT INTO `admin_category_articles` VALUES (1, 1, 1, 0, 0);
INSERT INTO `admin_category_articles` VALUES (2, 2, 1, 0, 1);
INSERT INTO `admin_category_articles` VALUES (3, 1, 1, 0, 1);

-- ----------------------------
-- Table structure for admin_rule
-- ----------------------------
DROP TABLE IF EXISTS `admin_rule`;
CREATE TABLE `admin_rule`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '规则名',
  `rule` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '规则，一般是路由。如： /user/list',
  `param` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '规则参数。如： limit=10&page=1',
  `status` tinyint(1) UNSIGNED NULL DEFAULT 0 COMMENT '状态 1-启用 0-禁用',
  `create_time` int(10) UNSIGNED NULL DEFAULT 0,
  `soft` int(10) NULL DEFAULT 0 COMMENT '排序',
  `open` tinyint(1) UNSIGNED NULL DEFAULT 0 COMMENT '是否对所有人员开放 1-是 0-否',
  `add_staff` bigint(20) UNSIGNED NULL DEFAULT 0 COMMENT '添加人员',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_rule
-- ----------------------------
INSERT INTO `admin_rule` VALUES (1, '测试', '/user/list', '', 1, 0, 1, 0, 1);

-- ----------------------------
-- Table structure for admin_system_logs
-- ----------------------------
DROP TABLE IF EXISTS `admin_system_logs`;
CREATE TABLE `admin_system_logs`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `type` tinyint(2) UNSIGNED NOT NULL DEFAULT 1 COMMENT '日志类型；1-登录日志 2-操作日志',
  `create_time` int(11) UNSIGNED NOT NULL DEFAULT 0,
  `IP` char(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '客户端IP',
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '访问路径(去参数',
  `user_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID (如果有',
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '用户名，如果用户ID不存在，则显示用户名',
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '日志详细内容',
  `param` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '附带参数(json',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 18 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_system_logs
-- ----------------------------
INSERT INTO `admin_system_logs` VALUES (1, 1, 1602607874, '127.0.0.1', '/backstage/article/list', 0, '', 'test', '');
INSERT INTO `admin_system_logs` VALUES (2, 1, 1602608234, '127.0.0.1', '/backstage/article/list', 0, '', 'test', '');
INSERT INTO `admin_system_logs` VALUES (3, 1, 1602608296, '127.0.0.1', '/backstage/article/list', 0, '', 'test', '');
INSERT INTO `admin_system_logs` VALUES (4, 1, 1602608329, '127.0.0.1', '/backstage/article/list', 0, '', 'test', '');
INSERT INTO `admin_system_logs` VALUES (5, 1, 1602608577, '127.0.0.1', '/backstage/article/list', 0, '', 'test', '');
INSERT INTO `admin_system_logs` VALUES (6, 1, 1602608644, '127.0.0.1', '/backstage/article/list', 0, '', 'test', '{\"requestBody\":{\"Title\":\"而已\",\"limit\":1},\"requestQuery\":\"\"}');
INSERT INTO `admin_system_logs` VALUES (7, 1, 1602609245, '127.0.0.1', '/backstage/article/list', 0, '', 'RequestInput', '{\"requestBody\":{\"Title\":\"而已\",\"limit\":1},\"requestQuery\":\"\"}');
INSERT INTO `admin_system_logs` VALUES (8, 1, 1602609255, '127.0.0.1', '/backstage/article/list', 0, '', 'RequestInput', '{\"requestBody\":{\"Title\":\"而已\",\"limit\":1},\"requestQuery\":\"\"}');
INSERT INTO `admin_system_logs` VALUES (9, 1, 1602609255, '127.0.0.1', '/backstage/article/list', 0, '', 'RequestInput', '{\"requestBody\":{\"Title\":\"而已\",\"limit\":1},\"requestQuery\":\"\"}');
INSERT INTO `admin_system_logs` VALUES (10, 1, 1602609256, '127.0.0.1', '/backstage/article/list', 0, '', 'RequestInput', '{\"requestBody\":{\"Title\":\"而已\",\"limit\":1},\"requestQuery\":\"\"}');
INSERT INTO `admin_system_logs` VALUES (11, 1, 1602609256, '127.0.0.1', '/backstage/article/list', 0, '', 'RequestInput', '{\"requestBody\":{\"Title\":\"而已\",\"limit\":1},\"requestQuery\":\"\"}');
INSERT INTO `admin_system_logs` VALUES (12, 1, 1602609256, '127.0.0.1', '/backstage/article/list', 0, '', 'RequestInput', '{\"requestBody\":{\"Title\":\"而已\",\"limit\":1},\"requestQuery\":\"\"}');
INSERT INTO `admin_system_logs` VALUES (13, 1, 1602609280, '127.0.0.1', '/backstage/article/list', 0, '', 'RequestInput', '{\"requestBody\":{\"Title\":\"而已\",\"limit\":1},\"requestQuery\":\"\"}');
INSERT INTO `admin_system_logs` VALUES (14, 1, 1602609284, '127.0.0.1', '/backstage/article/list', 0, '', 'RequestInput', '{\"requestBody\":{\"Title\":\"而已\",\"limit\":1},\"requestQuery\":\"\"}');
INSERT INTO `admin_system_logs` VALUES (15, 1, 1602609285, '127.0.0.1', '/backstage/article/list', 0, '', 'RequestInput', '{\"requestBody\":{\"Title\":\"而已\",\"limit\":1},\"requestQuery\":\"\"}');
INSERT INTO `admin_system_logs` VALUES (16, 1, 1602609285, '127.0.0.1', '/backstage/article/list', 0, '', 'RequestInput', '{\"requestBody\":{\"Title\":\"而已\",\"limit\":1},\"requestQuery\":\"\"}');
INSERT INTO `admin_system_logs` VALUES (17, 1, 1602609286, '127.0.0.1', '/backstage/article/list', 0, '', 'RequestInput', '{\"requestBody\":{\"Title\":\"而已\",\"limit\":1},\"requestQuery\":\"\"}');

-- ----------------------------
-- Table structure for admin_user
-- ----------------------------
DROP TABLE IF EXISTS `admin_user`;
CREATE TABLE `admin_user`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_type` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT '用户类型;1:admin;2:会员',
  `sex` tinyint(2) NOT NULL DEFAULT 0 COMMENT '性别;0:保密,1:男,2:女',
  `birthday` int(11) NOT NULL DEFAULT 0 COMMENT '生日',
  `last_login_time` int(11) NOT NULL DEFAULT 0 COMMENT '最后登录时间',
  `last_login_ip` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '最后登录ip',
  `score` int(11) NOT NULL DEFAULT 0 COMMENT '用户积分',
  `gold` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '金币',
  `create_time` int(11) NOT NULL DEFAULT 0 COMMENT '注册时间',
  `update_time` int(11) UNSIGNED NULL DEFAULT 0 COMMENT '信息更新时间',
  `user_status` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT '用户状态;0:禁用,1:正常,2:未验证',
  `user_login` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
  `user_pass` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '登录密码;',
  `user_nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户昵称',
  `user_email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户登录邮箱',
  `user_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户个人网址',
  `avatar` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户头像',
  `signature` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '个性签名',
  `user_activation_key` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '激活码',
  `mobile` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户手机号',
  `position` int(11) NOT NULL DEFAULT 0 COMMENT '层级',
  `lock_time` int(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '登陆错误锁定结束时间',
  `lock_time_start` int(11) NOT NULL DEFAULT 0 COMMENT '登录错误锁定开始时间',
  `error_sum` tinyint(1) UNSIGNED NOT NULL DEFAULT 0 COMMENT '登陆错误次数',
  `first` tinyint(1) UNSIGNED NULL DEFAULT 1 COMMENT '是否首次登录系统',
  `last_edit_pass` int(10) UNSIGNED NULL DEFAULT 0 COMMENT '最后一次修改密码的时间',
  `openid` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '微信openid',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `user_login`(`user_login`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_user
-- ----------------------------
INSERT INTO `admin_user` VALUES (1, 1, 1, 0, 1602605845, '127.0.0.1', 0, 0, 0, 0, 1, 'admin', '###72f96bce79b5ba5645c09c8e98f6d91b', '超级管理员', '', '', 0, '', '', '', 0, 0, 0, 1, 1, 0, '');
INSERT INTO `admin_user` VALUES (2, 1, 0, 0, 0, '', 0, 0, 1587309956, 1587309956, 1, 'zhangdahuan', '###ced8c9360bdb264d8e7fc7090f7bce54', '張大歡', '', '', 0, '', '', '', 0, 0, 0, 0, 0, 0, '');
INSERT INTO `admin_user` VALUES (3, 1, 0, 0, 0, '', 0, 0, 1598121194, 1598121194, 1, 'zhangdahuan1', '###ced8c9360bdb264d8e7fc7090f7bce54', '張大歡1', '', '', 0, '', '', '', 0, 0, 0, 0, 1, 0, '');

SET FOREIGN_KEY_CHECKS = 1;
