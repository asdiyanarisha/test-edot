package dto

type (
	PayloadAddProduct struct {
		Name        string  `json:"name"`
		Price       float64 `json:"price"`
		Sku         string  `json:"sku"`
		ShopId      int     `json:"shop_id"`
		WarehouseId int     `json:"warehouse_id"`
		Qty         int     `json:"qty"`
	}

	ParameterQuery struct {
		Offset int    `form:"offset"`
		Limit  int    `form:"limit"`
		Search string `form:"search"`
	}

	ProductResponse struct {
		Id    int     `json:"id"`
		Name  string  `json:"name"`
		Price float64 `json:"price"`
		Sku   string  `json:"sku"`
		Shop  string  `json:"shop"`
		Stock int     `json:"stock"`
	}
)
