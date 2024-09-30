package dto

type (
	PayloadCreateShop struct {
		Name     string `json:"name" gorm:"column:name"`
		Location string `json:"location" gorm:"column:location"`
	}

	ResponseCreateShop struct {
		Id       int    `json:"id" gorm:"column:id"`
		Name     string `json:"name" gorm:"column:name"`
		Location string `json:"location" gorm:"column:location"`
	}
)
