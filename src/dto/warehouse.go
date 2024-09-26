package dto

type (
	PayloadAddWarehouse struct {
		Name     string `json:"name"`
		Location string `json:"location"`
	}

	ParameterQueryWarehouse struct {
		Status string `form:"status"`
	}
)
