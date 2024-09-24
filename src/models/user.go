package models

import "time"

type (
	User struct {
		Id        int       `json:"id"`
		FullName  string    `json:"full_name" gorm:"column:full_name"`
		Password  string    `json:"password" gorm:"column:password"`
		Role      string    `json:"role" gorm:"column:role"`
		Email     string    `json:"email" gorm:"column:email"`
		Phone     string    `json:"phone" gorm:"column:phone"`
		CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
		UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	}
)
