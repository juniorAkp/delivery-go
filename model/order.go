package model

import "go.mongodb.org/mongo-driver/v2/bson"

type Order struct {
	Model

	ProductId  bson.ObjectID `bson:"productId" json:"productId"`
	CustomerId bson.ObjectID `bson:"customerId" json:"customerId"`
}
