# ************************************************************
# Sequel Pro SQL dump
# Version 5446
#
# https://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 8.0.21)
# Database: seata_order
# Generation Time: 2021-04-05 09:38:10 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table so_item
# ------------------------------------------------------------

DROP TABLE IF EXISTS `so_item`;

CREATE TABLE `so_item` (
                           `sysno` bigint NOT NULL AUTO_INCREMENT,
                           `so_sysno` bigint DEFAULT NULL,
                           `product_sysno` bigint DEFAULT NULL,
                           `product_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '商品名称',
                           `cost_price` decimal(16,6) DEFAULT NULL COMMENT '成本价',
                           `original_price` decimal(16,6) DEFAULT NULL COMMENT '原价',
                           `deal_price` decimal(16,6) DEFAULT NULL COMMENT '成交价',
                           `quantity` int DEFAULT NULL COMMENT '数量',
                           PRIMARY KEY (`sysno`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单明细表';



# Dump of table so_master
# ------------------------------------------------------------

DROP TABLE IF EXISTS `so_master`;

CREATE TABLE `so_master` (
                             `sysno` bigint NOT NULL,
                             `so_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                             `buyer_user_sysno` bigint DEFAULT NULL COMMENT '下单用户号',
                             `seller_company_code` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '卖家公司编号',
                             `receive_division_sysno` bigint NOT NULL,
                             `receive_address` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                             `receive_zip` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                             `receive_contact` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                             `receive_contact_phone` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                             `stock_sysno` bigint DEFAULT NULL,
                             `payment_type` tinyint DEFAULT NULL COMMENT '支付方式：1，支付宝，2，微信',
                             `so_amt` decimal(16,6) DEFAULT NULL COMMENT '订单总额',
                             `status` tinyint DEFAULT NULL COMMENT '10,创建成功，待支付；30；支付成功，待发货；50；发货成功，待收货；70，确认收货，已完成；90，下单失败；100已作废',
                             `order_date` timestamp NULL DEFAULT NULL COMMENT '下单时间',
                             `paymemt_date` timestamp NULL DEFAULT NULL COMMENT '支付时间',
                             `delivery_date` timestamp NULL DEFAULT NULL COMMENT '发货时间',
                             `receive_date` timestamp NULL DEFAULT NULL COMMENT '发货时间',
                             `appid` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '订单来源',
                             `memo` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
                             `create_user` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                             `gmt_create` timestamp NULL DEFAULT NULL,
                             `modify_user` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                             `gmt_modified` timestamp NULL DEFAULT NULL,
                             PRIMARY KEY (`sysno`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';



# Dump of table undo_log
# ------------------------------------------------------------

DROP TABLE IF EXISTS `undo_log`;

CREATE TABLE `undo_log` (
                            `id` bigint NOT NULL AUTO_INCREMENT,
                            `branch_id` bigint NOT NULL,
                            `xid` varchar(100) NOT NULL,
                            `context` varchar(128) NOT NULL,
                            `rollback_info` longblob NOT NULL,
                            `log_status` int NOT NULL,
                            `log_created` datetime NOT NULL,
                            `log_modified` datetime NOT NULL,
                            `ext` varchar(100) DEFAULT NULL,
                            PRIMARY KEY (`id`),
                            KEY `idx_unionkey` (`xid`,`branch_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
