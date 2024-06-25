package models

type PostTag struct {
	ID     int `json:"id"`
	PostId int `json:"id_post" gorm:"column:id_post"`
	TagId  int `json:"id_tag" gorm:"column:id_tag" `
}
