CREATE DATABASE IF NOT EXISTS mud;

CREATE TABLE IF NOT EXISTS `mud`.`registrations` (
  `id` BINARY(16) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `timezone` VARCHAR(45) NOT NULL,
  `password` VARBINARY(256) NOT NULL,
  `shortbio` LONGTEXT NULL,
  `validated` TINYINT NULL,
  PRIMARY KEY (`id`, `name`));


