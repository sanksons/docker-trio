-- name: mysql-sync-drop
DROP TABLE `styloko_mysql_syncing`;
-- name: mysql-sync
CREATE TABLE IF NOT EXISTS `styloko_mysql_syncing` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `version` bigint(20) NOT NULL,
  `type` varchar(255) NOT NULL,
  `resource` int(11) NOT NULL,
  `resource_type` varchar(255) NOT NULL,
  `status` enum('init','processing','failed','success','canceled') NOT NULL DEFAULT 'init',
  `trials` int(11) NOT NULL DEFAULT '0',
  `data` blob,
  `created_at` datetime NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `version` (`version`),
  KEY `type` (`type`),
  KEY `identifier` (`resource`),
  KEY `identifier_type` (`resource_type`),
  KEY `status` (`status`),
  KEY `trials` (`trials`),
  KEY `resource_type_status` (`resource_type`,`status`),
  KEY `updated_at` (`updated_at`) 
) ENGINE=InnoDB; 


