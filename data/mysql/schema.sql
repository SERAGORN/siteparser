CREATE TABLE `siteparser`.`tbl_article` (
        id INT UNSIGNED AUTO_INCREMENT NOT NULL PRIMARY KEY,
        title VARCHAR(200) NOT NULL,
        description TEXT NOT NULL,
        source_url VARCHAR(100) NOT NULL,
        site_url VARCHAR(100) NOT NULL,
        FULLTEXT (title,description)
);