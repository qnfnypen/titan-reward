-- explorer
CREATE TABLE IF NOT EXISTS `users` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `username` varchar(255) NOT NULL DEFAULT '' COMMENT '用户名称,邮箱/钱包地址',
    `role` tinyint NOT NULL DEFAULT '0' COMMENT '2的时候为KOL',
    `closed_test_reward` decimal(14,6) NOT NULL DEFAULT '0' COMMENT '封测奖励',
    `huygens_reward` decimal(14,6) NOT NULL DEFAULT '0' COMMENT '惠更斯测试网节点收益',
    `huygens_referral_reward` decimal(14,6) NOT NULL DEFAULT '0' COMMENT '惠更斯测试网邀请好友收益',
    `herschel_reward` decimal(14,6) NOT NULL DEFAULT '0' COMMENT '郝歇尔测试网邀请好友收益',
    `herschel_referral_reward` decimal(14,6) NOT NULL DEFAULT '0' COMMENT '郝歇尔测试网邀请好友收益',
    `reward` double NOT NULL DEFAULT '0' COMMENT '当前网络邀请好友收益',
    `referral_reward` double NOT NULL DEFAULT '0' COMMENT '当前网络邀请好友收益',
    `from_kol_bonus_reward` decimal(14,6) NOT NULL DEFAULT '0' COMMENT '受KOL邀请用户的额外奖励',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_username` (`username`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- quest
CREATE TABLE IF NOT EXISTS `user_mission` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `username` varchar(128) NOT NULL DEFAULT '' COMMENT '用户名称,邮箱/钱包地址',
    `credit` bigint NOT NULL DEFAULT '0' COMMENT '社区任务收益',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户社区收益表';

CREATE TABLE IF NOT EXISTS `invite_log` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `username` varchar(255) NOT NULL DEFAULT '' COMMENT '用户名称,邮箱/钱包地址',
    `credit` bigint NOT NULL DEFAULT '0' COMMENT '社区邀请好友收益',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户社区邀请好友收益表';