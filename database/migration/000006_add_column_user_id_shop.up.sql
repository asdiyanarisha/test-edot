ALTER TABLE `shops`
    ADD COLUMN `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 AFTER `location`,
    ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;