package models

import "time"

type (
	Order struct {
		Id        int       `json:"id" gorm:"primaryKey;column:id"`
		OrderNo   string    `json:"order_no" gorm:"column:order_no"`
		UserId    int       `json:"user_id" gorm:"column:user_id"`
		IsPayment bool      `json:"is_payment" gorm:"column:is_payment"`
		IsRelease bool      `json:"is_release" gorm:"column:is_release"`
		Total     float64   `json:"total" gorm:"column:total"`
		ExpiredAt time.Time `json:"expired_at" gorm:"column:expired_at"`
		CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
		UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	}

	OrderDetail struct {
		Id        int       `json:"id" gorm:"primaryKey;column:id"`
		OrderId   int       `json:"order_id" gorm:"column:order_id"`
		ProductId int       `json:"product_id" gorm:"column:product_id"`
		StockId   int       `json:"stock_id" gorm:"column:stock_id"`
		Qty       int       `json:"qty" gorm:"column:qty"`
		Total     float64   `json:"total" gorm:"column:total"`
		ExpiredAt time.Time `json:"expired_at" gorm:"column:expired_at"`
		CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
		UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	}

	OrderWithDetail struct {
		Id        int           `json:"id" gorm:"primaryKey;column:id"`
		UserId    int           `json:"user_id" gorm:"column:user_id"`
		IsPayment bool          `json:"is_payment" gorm:"column:is_payment"`
		Total     float64       `json:"total" gorm:"column:total"`
		ExpiredAt time.Time     `json:"expired_at" gorm:"column:expired_at"`
		Items     []OrderDetail `json:"items" gorm:"-"`
		CreatedAt time.Time     `json:"created_at" gorm:"column:created_at"`
		UpdatedAt time.Time     `json:"updated_at" gorm:"column:updated_at"`
	}
)
