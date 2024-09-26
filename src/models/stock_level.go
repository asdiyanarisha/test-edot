package models

import "time"

type StockLevel struct {
	ID            int       `gorm:"primaryKey" json:"id"`
	ProductId     int       `json:"product_id"  gorm:"column:product_id"`
	WarehouseId   int       `json:"warehouse_id" gorm:"column:warehouse_id"`
	Stock         int       `json:"stock"  gorm:"column:stock"`
	ReservedStock int       `json:"reserved_stock"  gorm:"column:reserved_stock"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
