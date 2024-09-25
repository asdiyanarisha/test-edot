ALTER TABLE `warehouses`
    ADD COLUMN `is_active` TINYINT NOT NULL DEFAULT 0 AFTER `location`,
    ADD COLUMN `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 AFTER `is_active`,
    ADD CONSTRAINT fk_warehouse_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;