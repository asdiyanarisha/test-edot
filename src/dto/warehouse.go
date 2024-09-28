package dto

type (
	PayloadAddWarehouse struct {
		Name     string `json:"name"`
		Location string `json:"location"`
	}

	ParameterQueryWarehouse struct {
		Status string `form:"status"`
	}

	ParameterChangeStatusWarehouse struct {
		WarehouseId int  `json:"warehouse_id"`
		IsActive    bool `json:"is_active"`
	}
)
