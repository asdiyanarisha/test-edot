package models

import "time"

type (
	Product struct {
		Id        int       `json:"id" gorm:"primaryKey,column:id"`
		Name      string    `json:"name" gorm:"column:name"`
		Sku       string    `json:"sku" gorm:"column:sku"`
		Price     float64   `json:"price" gorm:"column:price"`
		ShopId    int       `json:"shop_id" gorm:"column:shop_id"`
		CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
		UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	}

	ProductDetail struct {
		Id        int          `json:"id" gorm:"primaryKey;column:id;index"`
		Name      string       `json:"name" gorm:"column:name"`
		Sku       string       `json:"sku" gorm:"column:sku"`
		Price     float64      `json:"price" gorm:"column:price"`
		ShopId    int          `json:"shop_id" gorm:"column:shop_id"`
		Shop      Shop         `json:"shop" gorm:"foreignKey:shop_id"`
		Stock     []StockLevel `json:"stock" gorm:"foreignKey:product_id;references:Id"`
		CreatedAt time.Time    `json:"created_at" gorm:"column:created_at"`
		UpdatedAt time.Time    `json:"updated_at" gorm:"column:updated_at"`
	}
)

func (ProductDetail) TableName() string {
	return "products"
}
