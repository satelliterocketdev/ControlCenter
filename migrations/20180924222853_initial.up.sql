CREATE TABLE `users` (
    `id` INT(11) auto_increment,
    `email` VARCHAR(255) NOT NULL,
    `password` blob NOT NULL,
    `first_name` VARCHAR(255) NOT NULL,
    `last_name` VARCHAR(255) NOT NULL,
    `roles` JSON NOT NULL,
    `is_active` TINYINT(1) NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_email_users` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO users (email, password, first_name, last_name, roles, is_active) VALUES('admin1@i69social.com', AES_ENCRYPT('i69social_pass', 'secret'), 'Admin', 'Test', '["admin"]', 1);
INSERT INTO users (email, password, first_name, last_name, roles, is_active) VALUES('admin2@i69social.com', AES_ENCRYPT('i69social_pass', 'secret'), 'Admin', 'Test', '["admin", "chatter"]', 1);
INSERT INTO users (email, password, first_name, last_name, roles, is_active) VALUES('worker@i69social.com', AES_ENCRYPT('i69social_pass', 'secret'), 'Admin', 'Test', '["chatter"]', 1);
