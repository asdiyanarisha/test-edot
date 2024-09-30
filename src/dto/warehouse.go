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

	ResponseWarehouse struct {
		ID       int    `json:"id" gorm:"primary_key,column:id"`
		Name     string `json:"name" gorm:"column:name"`
		Location string `json:"location" gorm:"column:location"`
		UserId   int    `json:"user_id,omitempty" gorm:"column:user_id"`
		IsActive bool   `json:"is_active" gorm:"column:is_active"`
	}
)
