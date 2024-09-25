package dto

type (
	PayloadCreateShop struct {
		Name     string `json:"name" gorm:"column:name"`
		Location string `json:"location" gorm:"column:location"`
	}
)
