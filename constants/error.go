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
	InvalidPassword          = errors.New("password invalid")
	BearerExpired            = errors.New("bearer expired")
	ProductAlreadyInserted   = errors.New("product already inserted")
	WarehouseAlreadyExisted  = errors.New("warehouse already existed")
)
