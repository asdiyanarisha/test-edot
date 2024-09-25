ALTER TABLE `warehouses` DROP COLUMN `is_active`,
DROP COLUMN `user_id`,
     DROP FOREIGN KEY fk_warehouse_user_id;