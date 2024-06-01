CREATE DATABASE IF NOT EXISTS `titan_reward` DEFAULT CHARACTER SET 'utf8mb4';

USE `titan_reward`;

CREATE TABLE IF NOT EXISTS `user` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `uuid` varchar(63) NOT NULL DEFAULT '' COMMENT '用户uuid',
    `email` varchar(255) NOT NULL DEFAULT '' COMMENT '用户邮箱',
    `wallet_addr` varchar(255) NOT NULL DEFAULT '' COMMENT '用户钱包地址',
    `address` varchar(255) NOT NULL DEFAULT '' COMMENT 'titan的钱包地址',
    `created_at` int(10) NOT NULL DEFAULT '0' COMMENT '创建时间',
    `deleted_at` int(10) NOT NULL DEFAULT '0' COMMENT '删除时间',
    `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '提现状态:0-提现中 1-已到账',
    PRIMARY KEY(`id`),
    UNIQUE KEY `uq_uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';