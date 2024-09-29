ALTER TABLE `order_details` DROP COLUMN `expired_at`,
DROP COLUMN `stock_id`,
     DROP FOREIGN KEY fk_order_detail_stock_id;