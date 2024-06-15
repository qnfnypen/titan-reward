CREATE DATABASE IF NOT EXISTS `titan_pledge` DEFAULT CHARACTER SET 'utf8mb4';

USE `titan_pledge`;

CREATE TABLE IF NOT EXISTS `user` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `uuid` varchar(63) NOT NULL DEFAULT '' COMMENT '用户uuid',
    `wallet` varchar(255) NOT NULL DEFAULT '' COMMENT 'titan的钱包地址',
    `created_time` int(10) NOT NULL DEFAULT '0' COMMENT '创建时间',
    `deleted_time` int(10) NOT NULL DEFAULT '0' COMMENT '删除时间',
    `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '提现状态:0-未提现 1-提现中 2-已提现',
    PRIMARY KEY(`id`),
    UNIQUE KEY `uq_uuid` (`uuid`),
    UNIQUE KEY `uq_wallet` (`wallet`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';