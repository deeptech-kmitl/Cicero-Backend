package order

import "github.com/deeptech-kmitl/Cicero-Backend/modules/users"

type OrderProducts struct {
	Products []*users.Cart
}

type AddOrderReq struct {
	UserId        string         `json:"user_id" form:"user_id" db:"user_id"`
	Total         float64        `json:"total" form:"total" db:"total"`
	Status        string         `json:"status" form:"status" db:"status"`
	Address       *Address       `json:"address" form:"address" db:"address"`
	PaymentDetail *PaymentDetail `json:"payment_detail" form:"payment_detail" db:"payment_detail"`
}

type Address struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Street    string `json:"street"`
	City      string `json:"city"`
	ZipCode   string `json:"zip_code"`
	Country   string `json:"country"`
}

type PaymentDetail struct {
	CardHolder string `json:"card_holder"`
	CardNumber string `json:"card_number"`
	Expired    string `json:"expired"`
	Cvv        string `json:"cvv"`
}
