package models

import "time"

type Warehouse struct {
	ID        int       `json:"id" gorm:"primary_key,column:id"`
	Name      string    `json:"name" gorm:"column:name"`
	Location  string    `json:"location" gorm:"column:location"`
	UserId    int       `json:"user_id,omitempty" gorm:"column:user_id"`
	IsActive  bool      `json:"is_active" gorm:"column:is_active"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}
