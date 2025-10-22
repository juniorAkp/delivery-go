package model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type VehicleType string

const (
	Car   VehicleType = "car"
	Bike  VehicleType = "bike"
	Motor VehicleType = "motor"
)

type Vehicle struct {
	Model

	RegistrationNumber string        `bson:"regNumber" json:"regNumber"`
	Type               VehicleType   `bson:"vehicleType" json:"vehicleType"`
	VehicleModel       string        `bson:"vehicleModel" json:"vehicleModel"`
	CourierId          bson.ObjectID `bson:"courierId" json:"courierId"`

	State  State `bson:"state" json:"state"`
	Active bool  `bson:"active" json:"active"`
}
