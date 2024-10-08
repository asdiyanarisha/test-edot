package dto

import "test-edot/src/models"

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

	TransferProductWarehouse struct {
		FromWarehouseId int `json:"from_warehouse_id"`
		ToWarehouseId   int `json:"to_warehouse_id"`
		Qty             int `json:"qty"`
	}

	InitialTransferProduct struct {
		Product       models.Product
		FromWarehouse models.Warehouse
		ToWarehouse   models.Warehouse
	}

	ProductResponse struct {
		Id    int     `json:"id"`
		Name  string  `json:"name"`
		Price float64 `json:"price"`
		Sku   string  `json:"sku"`
		Shop  string  `json:"shop"`
		Stock int     `json:"stock"`
	}

	ProductDetailResponse struct {
		Id            int     `json:"id"`
		Name          string  `json:"name"`
		Price         float64 `json:"price"`
		Sku           string  `json:"sku"`
		Shop          string  `json:"shop"`
		Stock         int     `json:"stock"`
		ReservedStock int     `json:"reserved_stock"`
	}
)
