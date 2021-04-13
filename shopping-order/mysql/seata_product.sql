# ************************************************************
# Sequel Pro SQL dump
# Version 5446
#
# https://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 8.0.21)
# Database: seata_product
# Generation Time: 2021-04-05 09:40:33 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

# Dump of table inventory
# ------------------------------------------------------------

DROP TABLE IF EXISTS `inventory`;

CREATE TABLE `inventory` (
                             `sysno` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
                             `product_sysno` bigint unsigned NOT NULL,
                             `account_qty` int DEFAULT NULL COMMENT '财务库存',
                             `available_qty` int DEFAULT NULL COMMENT '可用库存',
                             `allocated_qty` int DEFAULT NULL COMMENT '分配库存',
                             `adjust_locked_qty` int DEFAULT NULL COMMENT '调整锁定库存',
                             PRIMARY KEY (`sysno`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品库存';

LOCK TABLES `inventory` WRITE;
/*!40000 ALTER TABLE `inventory` DISABLE KEYS */;

INSERT INTO `inventory` (`sysno`, `product_sysno`, `account_qty`, `available_qty`, `allocated_qty`, `adjust_locked_qty`)
VALUES
(1,1,1000000,1000000,0,0);

/*!40000 ALTER TABLE `inventory` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table product
# ------------------------------------------------------------

DROP TABLE IF EXISTS `product`;

CREATE TABLE `product` (
                           `sysno` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
                           `product_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '品名',
                           `product_title` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
                           `product_desc` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '描述',
                           `product_desc_long` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '描述',
                           `default_image_src` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                           `c3_sysno` bigint DEFAULT NULL,
                           `barcode` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                           `length` int DEFAULT NULL,
                           `width` int DEFAULT NULL,
                           `height` int DEFAULT NULL,
                           `weight` float DEFAULT NULL,
                           `merchant_sysno` bigint DEFAULT NULL,
                           `merchant_productid` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                           `status` tinyint NOT NULL DEFAULT '1' COMMENT '1，待上架；2，上架；3，下架；4，售罄下架；5，违规下架',
                           `gmt_create` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                           `create_user` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '创建人',
                           `modify_user` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '修改人',
                           `gmt_modified` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                           PRIMARY KEY (`sysno`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品SKU';

LOCK TABLES `product` WRITE;
/*!40000 ALTER TABLE `product` DISABLE KEYS */;

INSERT INTO `product` (`sysno`, `product_name`, `product_title`, `product_desc`, `product_desc_long`, `default_image_src`, `c3_sysno`, `barcode`, `length`, `width`, `height`, `weight`, `merchant_sysno`, `merchant_productid`, `status`, `gmt_create`, `create_user`, `modify_user`, `gmt_modified`)
VALUES
(1,'北冰洋','北冰洋汽水让人怀念','好喝好喝好好喝','','https://bkimg.cdn.bcebos.com/pic/7e3e6709c93d70cfd13954baf2dcd100bba12b42',0,'',15,5,5,5,1,'',1,'2021-04-05 17:32:36','','','2021-04-05 17:36:30');

/*!40000 ALTER TABLE `product` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table undo_log
# ------------------------------------------------------------

DROP TABLE IF EXISTS `undo_log`;

CREATE TABLE `undo_log` (
                            `id` bigint NOT NULL AUTO_INCREMENT,
                            `branch_id` bigint NOT NULL,
                            `xid` varchar(100) NOT NULL,
                            `context` varchar(128) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
                            `rollback_info` longblob NOT NULL,
                            `log_status` int NOT NULL,
                            `log_created` datetime NOT NULL,
                            `log_modified` datetime NOT NULL,
                            `ext` varchar(100) DEFAULT NULL,
                            PRIMARY KEY (`id`),
                            KEY `idx_unionkey` (`xid`,`branch_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
