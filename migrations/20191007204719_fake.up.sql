CREATE TABLE `fake_users` (
    `id` INT(11) auto_increment,
    `external_id` VARCHAR(255) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `owner_id` INT(11) NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`id`),
   FOREIGN KEY (`owner_id`) REFERENCES users(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
