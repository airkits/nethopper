CREATE TABLE IF NOT EXISTS `user_0`(
   `uid` BIGINT UNSIGNED NOT NULL,
   `appid` CHAR(64) NOT NULL,
   `openid` CHAR(255) NOT NULL,
   `uuid` CHAR(255) DEFAULT NULL,
   `avatar` CHAR(255) NOT NULL,
   `name` CHAR(32) NOT NULL,
   `password` CHAR(64) DEFAULT NULL,
   `phone` CHAR(20) DEFAULT NULL,
   `gender` INT NOT NULL DEFAULT 0,
   `age` INT NOT NULL DEFAULT 0,
   `gold` BIGINT NOT NULL DEFAULT 0,
   `coin` BIGINT NOT NULL DEFAULT 0,
   `status` INT NOT NULL DEFAULT 0,
   `channel` CHAR(100) NOT NULL DEFAULT "", 
   `loginip` CHAR(100) NOT NULL DEFAULT "", 
   `loginat` DATE,
   `createat` DATE,
   PRIMARY KEY ( `uid` ),
   UNIQUE KEY `idx_appid_oid` ( `appid` , `openid` ),
   KEY `idx_appid_uuid` ( `appid` , `uuid` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `user_1`(
   `uid` BIGINT UNSIGNED NOT NULL,
   `appid` CHAR(64) NOT NULL,
   `openid` CHAR(255) NOT NULL,
   `uuid` CHAR(255) DEFAULT NULL,
   `avatar` CHAR(255) NOT NULL,
   `name` CHAR(32) NOT NULL,
   `password` CHAR(64) DEFAULT NULL,
   `phone` CHAR(20) DEFAULT NULL,
   `gender` INT NOT NULL DEFAULT 0,
   `age` INT NOT NULL DEFAULT 0,
   `gold` BIGINT NOT NULL DEFAULT 0,
   `coin` BIGINT NOT NULL DEFAULT 0,
   `status` INT NOT NULL DEFAULT 0,
   `channel` CHAR(100) NOT NULL DEFAULT "", 
   `loginip` CHAR(100) NOT NULL DEFAULT "", 
   `loginat` DATE,
   `createat` DATE,
   PRIMARY KEY ( `uid` ),
   UNIQUE KEY `idx_appid_oid` ( `appid` , `openid` ),
   KEY `idx_appid_uuid` ( `appid` , `uuid` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `user_2`(
   `uid` BIGINT UNSIGNED NOT NULL,
   `appid` CHAR(64) NOT NULL,
   `openid` CHAR(255) NOT NULL,
   `uuid` CHAR(255) DEFAULT NULL,
   `avatar` CHAR(255) NOT NULL,
   `name` CHAR(32) NOT NULL,
   `password` CHAR(64) DEFAULT NULL,
   `phone` CHAR(20) DEFAULT NULL,
   `gender` INT NOT NULL DEFAULT 0,
   `age` INT NOT NULL DEFAULT 0,
   `gold` BIGINT NOT NULL DEFAULT 0,
   `coin` BIGINT NOT NULL DEFAULT 0,
   `status` INT NOT NULL DEFAULT 0,
   `channel` CHAR(100) NOT NULL DEFAULT "", 
   `loginip` CHAR(100) NOT NULL DEFAULT "", 
   `loginat` DATE,
   `createat` DATE,
   PRIMARY KEY ( `uid` ),
   UNIQUE KEY `idx_appid_oid` ( `appid` , `openid` ),
   KEY `idx_appid_uuid` ( `appid` , `uuid` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `user_3`(
   `uid` BIGINT UNSIGNED NOT NULL,
   `appid` CHAR(64) NOT NULL,
   `openid` CHAR(255) NOT NULL,
   `uuid` CHAR(255) DEFAULT NULL,
   `avatar` CHAR(255) NOT NULL,
   `name` CHAR(32) NOT NULL,
   `password` CHAR(64) DEFAULT NULL,
   `phone` CHAR(20) DEFAULT NULL,
   `gender` INT NOT NULL DEFAULT 0,
   `age` INT NOT NULL DEFAULT 0,
   `gold` BIGINT NOT NULL DEFAULT 0,
   `coin` BIGINT NOT NULL DEFAULT 0,
   `status` INT NOT NULL DEFAULT 0,
   `channel` CHAR(100) NOT NULL DEFAULT "", 
   `loginip` CHAR(100) NOT NULL DEFAULT "", 
   `loginat` DATE,
   `createat` DATE,
   PRIMARY KEY ( `uid` ),
   UNIQUE KEY `idx_appid_oid` ( `appid` , `openid` ),
   KEY `idx_appid_uuid` ( `appid` , `uuid` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `user_4`(
   `uid` BIGINT UNSIGNED NOT NULL,
   `appid` CHAR(64) NOT NULL,
   `openid` CHAR(255) NOT NULL,
   `uuid` CHAR(255) DEFAULT NULL,
   `avatar` CHAR(255) NOT NULL,
   `name` CHAR(32) NOT NULL,
   `password` CHAR(64) DEFAULT NULL,
   `phone` CHAR(20) DEFAULT NULL,
   `gender` INT NOT NULL DEFAULT 0,
   `age` INT NOT NULL DEFAULT 0,
   `gold` BIGINT NOT NULL DEFAULT 0,
   `coin` BIGINT NOT NULL DEFAULT 0,
   `status` INT NOT NULL DEFAULT 0,
   `channel` CHAR(100) NOT NULL DEFAULT "", 
   `loginip` CHAR(100) NOT NULL DEFAULT "", 
   `loginat` DATE,
   `createat` DATE,
   PRIMARY KEY ( `uid` ),
   UNIQUE KEY `idx_appid_oid` ( `appid` , `openid` ),
   KEY `idx_appid_uuid` ( `appid` , `uuid` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `user_5`(
   `uid` BIGINT UNSIGNED NOT NULL,
   `appid` CHAR(64) NOT NULL,
   `openid` CHAR(255) NOT NULL,
   `uuid` CHAR(255) DEFAULT NULL,
   `avatar` CHAR(255) NOT NULL,
   `name` CHAR(32) NOT NULL,
   `password` CHAR(64) DEFAULT NULL,
   `phone` CHAR(20) DEFAULT NULL,
   `gender` INT NOT NULL DEFAULT 0,
   `age` INT NOT NULL DEFAULT 0,
   `gold` BIGINT NOT NULL DEFAULT 0,
   `coin` BIGINT NOT NULL DEFAULT 0,
   `status` INT NOT NULL DEFAULT 0,
   `channel` CHAR(100) NOT NULL DEFAULT "", 
   `loginip` CHAR(100) NOT NULL DEFAULT "", 
   `loginat` DATE,
   `createat` DATE,
   PRIMARY KEY ( `uid` ),
   UNIQUE KEY `idx_appid_oid` ( `appid` , `openid` ),
   KEY `idx_appid_uuid` ( `appid` , `uuid` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `user_6`(
   `uid` BIGINT UNSIGNED NOT NULL,
   `appid` CHAR(64) NOT NULL,
   `openid` CHAR(255) NOT NULL,
   `uuid` CHAR(255) DEFAULT NULL,
   `avatar` CHAR(255) NOT NULL,
   `name` CHAR(32) NOT NULL,
   `password` CHAR(64) DEFAULT NULL,
   `phone` CHAR(20) DEFAULT NULL,
   `gender` INT NOT NULL DEFAULT 0,
   `age` INT NOT NULL DEFAULT 0,
   `gold` BIGINT NOT NULL DEFAULT 0,
   `coin` BIGINT NOT NULL DEFAULT 0,
   `status` INT NOT NULL DEFAULT 0,
   `channel` CHAR(100) NOT NULL DEFAULT "", 
   `loginip` CHAR(100) NOT NULL DEFAULT "", 
   `loginat` DATE,
   `createat` DATE,
   PRIMARY KEY ( `uid` ),
   UNIQUE KEY `idx_appid_oid` ( `appid` , `openid` ),
   KEY `idx_appid_uuid` ( `appid` , `uuid` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `user_7`(
   `uid` BIGINT UNSIGNED NOT NULL,
   `appid` CHAR(64) NOT NULL,
   `openid` CHAR(255) NOT NULL,
   `uuid` CHAR(255) DEFAULT NULL,
   `avatar` CHAR(255) NOT NULL,
   `name` CHAR(32) NOT NULL,
   `password` CHAR(64) DEFAULT NULL,
   `phone` CHAR(20) DEFAULT NULL,
   `gender` INT NOT NULL DEFAULT 0,
   `age` INT NOT NULL DEFAULT 0,
   `gold` BIGINT NOT NULL DEFAULT 0,
   `coin` BIGINT NOT NULL DEFAULT 0,
   `status` INT NOT NULL DEFAULT 0,
   `channel` CHAR(100) NOT NULL DEFAULT "", 
   `loginip` CHAR(100) NOT NULL DEFAULT "", 
   `loginat` DATE,
   `createat` DATE,
   PRIMARY KEY ( `uid` ),
   UNIQUE KEY `idx_appid_oid` ( `appid` , `openid` ),
   KEY `idx_appid_uuid` ( `appid` , `uuid` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE IF NOT EXISTS `oid2uid_0`(
   `openid` CHAR(255) NOT NULL,
   `uid` BIGINT UNSIGNED NOT NULL,
   PRIMARY KEY ( `openid` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `oid2uid_1`(
   `openid` CHAR(255) NOT NULL,
   `uid` BIGINT UNSIGNED NOT NULL,
   PRIMARY KEY ( `openid` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `oid2uid_2`(
   `openid` CHAR(255) NOT NULL,
   `uid` BIGINT UNSIGNED NOT NULL,
   PRIMARY KEY ( `openid` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `oid2uid_3`(
   `openid` CHAR(255) NOT NULL,
   `uid` BIGINT UNSIGNED NOT NULL,
   PRIMARY KEY ( `openid` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `oid2uid_4`(
   `openid` CHAR(255) NOT NULL,
   `uid` BIGINT UNSIGNED NOT NULL,
   PRIMARY KEY ( `openid` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `oid2uid_5`(
   `openid` CHAR(255) NOT NULL,
   `uid` BIGINT UNSIGNED NOT NULL,
   PRIMARY KEY ( `openid` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `oid2uid_6`(
   `openid` CHAR(255) NOT NULL,
   `uid` BIGINT UNSIGNED NOT NULL,
   PRIMARY KEY ( `openid` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `oid2uid_7`(
   `openid` CHAR(255) NOT NULL,
   `uid` BIGINT UNSIGNED NOT NULL,
   PRIMARY KEY ( `openid` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
