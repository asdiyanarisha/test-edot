ALTER TABLE `order_details`
    ADD COLUMN `expired_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP AFTER `total`,
    ADD COLUMN `stock_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 AFTER `product_id`,
    ADD CONSTRAINT fk_order_detail_stock_id FOREIGN KEY (stock_id) REFERENCES stock_levels(id) ON DELETE CASCADE;