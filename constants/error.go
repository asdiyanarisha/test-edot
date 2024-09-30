package constants

import "errors"

var (
	ErrorPostAlreadyInserted = errors.New("post already inserted")
	ErrorPostNotFound        = errors.New("post not found")
	RolePayloadInvalid       = errors.New("role payload invalid")
	FormatEmailInvalid       = errors.New("email format invalid")
	FormatPhoneInvalid       = errors.New("phone number format invalid, range 8-15 numeric")
	UserAlreadyInserted      = errors.New("user already inserted")
	UserNotFound             = errors.New("user not found")
	ShopAlreadyInserted      = errors.New("shop already inserted")
	ShopNotFound             = errors.New("shop not found")
	OrderNotFound            = errors.New("order not found")
	InvalidPassword          = errors.New("password invalid")
	BearerExpired            = errors.New("bearer expired")
	ProductAlreadyInserted   = errors.New("product already inserted")
	ProductNotFound          = errors.New("product not found")
	WarehouseAlreadyExisted  = errors.New("warehouse already existed")
	WarehouseNotFound        = errors.New("warehouse not found")
	FromWarehouseNotFound    = errors.New("from warehouse not found")
	ToWarehouseNotFound      = errors.New("to warehouse not found")
	DuplicateProduct         = errors.New("duplicate product")
	StatusNotSamePrevious    = errors.New("status not same previous status")
	StockMustEmpty           = errors.New("for inactive warehouse stock must be empty")
	StockProductEmpty        = errors.New("stock product is empty")
	NotEnoughStockToTransfer = errors.New("not enough stock to transfer")
)
