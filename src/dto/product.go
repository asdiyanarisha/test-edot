package dto

type (
	PayloadAddProduct struct {
		Name   string  `json:"name"`
		Price  float64 `json:"price"`
		Sku    string  `json:"sku"`
		ShopId int     `json:"shop_id"`
	}
)
