package dto

type (
	PayloadCreateOrder struct {
		Items []PayloadCreateOrderItems `json:"items"`
	}

	PayloadCreateOrderItems struct {
		ProductId int `json:"product_id"`
		Qty       int `json:"qty"`
	}
)
