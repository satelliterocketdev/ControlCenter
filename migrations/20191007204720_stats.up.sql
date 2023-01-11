CREATE TABLE `worker_stats` (
    `id` INT(11) auto_increment,
    `worker_id` INT(11) NOT NULL,
    `fake_user_id` VARCHAR(255) NOT NULL,
    `chat_id` VARCHAR(255) NOT NULL,
    `sent_messages` INT(11) NOT NULL DEFAULT 0,
    `received_messages` INT(11) NOT NULL DEFAULT 0,
    `date` DATE NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`id`),
   FOREIGN KEY (`worker_id`) REFERENCES users(`id`),
   UNIQUE KEY `worker_chat_uk` (`worker_id`, `fake_user_id`, `chat_id`, `date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
