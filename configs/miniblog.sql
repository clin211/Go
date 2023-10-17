/*
 Navicat Premium Data Transfer

 Source Server         : local
 Source Server Type    : MySQL
 Source Server Version : 50722
 Source Host           : 127.0.0.1:33066
 Source Schema         : miniblog

 Target Server Type    : MySQL
 Target Server Version : 50722
 File Encoding         : 65001

 Date: 17/10/2023 17:36:43
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for post
-- ----------------------------
DROP TABLE IF EXISTS `post`;
CREATE TABLE `post` (
`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`username` varchar(255) NOT NULL COMMENT '用户名',
`postID` varchar(256) NOT NULL COMMENT '帖子ID',
`title` varchar(256) NOT NULL COMMENT '标题',
`content` longtext NOT NULL COMMENT '内容',
`createdAt` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`updatedAt` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (`id`),
UNIQUE KEY `postID` (`postID`),
KEY `idx_username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=141 DEFAULT CHARSET=utf8;

INSERT INTO `post` (`username`, `postID`, `title`, `content`) VALUES
(`Forest`, `0a207d8d-ca19-42ff-8b22-2217925bc53f`, `Vue.js入门指南`, `Vue.js是一款流行的JavaScript框架，用于构建用户界面。它具有简单易学、灵活高效的特点，适用于开发单页面应用和复杂的前端项目。`),
(`Forest`, `0f401e46-c89c-4dc0-b948-bc42dc5aa4f6`, `React.js实战教程`, `React.js是一个用于构建用户界面的JavaScript库，由Facebook开发并维护。它通过组件化的方式，使得前端开发更加模块化和可复用。`),
(`Forest`, `9ee3a3f4-cd96-4ca2-bc23-1716899bb82c`, `Go语言入门教程`, `Go语言是一门简单高效的编程语言，由Google开发。它具有静态类型、垃圾回收、并发编程等特性，适用于构建高性能的后端服务和系统工具。`),
(`Forest`, `0264559f-1e72-4e06-a01a-8e59d592a38a`, `Vue.js组件开发`, `Vue.js的组件是构建用户界面的基本单元，它可以封装HTML、CSS和JavaScript，提供复用和可组合性。通过组件化开发，我们可以更好地管理和维护前端代码。`),
(`Forest`, `f2656a9c-7c47-4098-925a-c09c326bdb6f`, `React.js性能优化指南`, `React.js应用在复杂场景下可能面临性能问题，通过优化可以提升应用的响应速度和用户体验。本文将介绍一些React.js性能优化的技巧和最佳实践。`),
(`Forest`, `bae2aac4-116e-4f3e-98af-f7ab506273d5`, `Go语言并发编程`, `Go语言通过goroutine和channel实现了轻量级的并发编程模型。并发编程在处理高并发、IO密集型任务等场景下非常有用，本文将介绍Go语言并发编程的基本概念和用法。`),
(`Forest`, `f0f5f143-c740-4a0f-a44e-5ad44f4e98da`, `Vue.js全栈开发`, `Vue.js不仅可以用于构建前端应用，还可以与后端技术结合实现全栈开发。本文将介绍Vue.js全栈开发的基本概念和常见技术栈，帮助你构建全面的应用程序。`),
(`Forest`, `57b0f4d8-cacf-4dc8-8b8e-96750098ea8b`, `React.js表单处理`, `React.js提供了强大的表单处理能力，可以方便地处理用户输入、验证和提交。本文将介绍React.js表单处理的方法和常用技巧，帮助你构建交互性强的表单页面。`),
(`Forest`, `9f84c425-dea0-46aa-995f-48a9e5fe1049`, `Go语言Web框架比较`, `Go语言提供了一些优秀的Web框架，如gin、echo等。本文将介绍Go语言Web框架的比较和选择方式，帮助你选择合适的框架。`);
-- ----------------------------
-- Records of post
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(255) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '密码',
  `nickname` varchar(30) NOT NULL COMMENT '昵称',
  `email` varchar(256) NOT NULL COMMENT '电子邮件地址',
  `phone` varchar(16) NOT NULL COMMENT '手机号码',
  `createdAt` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updatedAt` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=MyISAM AUTO_INCREMENT=27 DEFAULT CHARSET=utf8;

-- inset user data
-- INSERT INTO `user` VALUES (`Forest`, `123456`, `forest`, `767425412lin@gmail.com`, `5432109876`);
-- ----------------------------
-- Records of user
-- ----------------------------
BEGIN;
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
