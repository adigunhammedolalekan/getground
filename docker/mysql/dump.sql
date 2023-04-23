CREATE TABLE `tables` (
  `id` INT NOT NULL auto_increment,
  `capacity` INT,
  `allowed_extras` INT,
  `created_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE `guests` (
  `id` INT NOT NULL auto_increment,
  `table_id` INT NOT NULL,
  `name` VARCHAR(128) NOT NULL,
  `accompanying_guests` INT,
  `created_at` datetime NOT NULL,
  `deleted_at` datetime,
  PRIMARY KEY (`id`)
);
