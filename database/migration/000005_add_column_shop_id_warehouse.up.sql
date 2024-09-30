ALTER TABLE `warehouses`
    ADD COLUMN `shop_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 AFTER `location`,
    ADD CONSTRAINT fk_warehouse_shop_id FOREIGN KEY (shop_id) REFERENCES shops(id) ON DELETE CASCADE;