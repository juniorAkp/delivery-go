package model

type Order struct {
	Model

	ProductId  string `bson:"productId" json:"productId"`
	CustomerId string `bson:"customerId" json:"customerId"`

	TotalAmount   float64 `bson:"totalAmount" json:"totalAmount"`
	TotalQuantity float64 `bson:"totalQuantity" json:"totalQuantity"`
}
