
CREATE TABLE `siteparser`.`tbl_article` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `title` VARCHAR(100) NULL,
  `description` VARCHAR(1000) NULL,
  PRIMARY KEY (`id`));

ALTER TABLE `siteparser`.`tbl_article`
CHANGE COLUMN `title` `title` VARCHAR(100) NOT NULL ,
CHANGE COLUMN `description` `description` VARCHAR(1000) NOT NULL ;