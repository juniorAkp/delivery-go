package model

import "go.mongodb.org/mongo-driver/v2/bson"

type State string

const (
	//delivery state
	Pending         State = "pending"
	PendingPickup   State = "pending_pickup"
	AtPickup        State = "at_pickup"
	DeliveryOngoing State = "delivery_ongoing"
	AtDropOff       State = "at_dropOff"
	Cancelled       State = "cancelled"
	Completed       State = "completed"

	//courier state
	AwaitingDispatch State = "awaiting_dispatch"
	Dispatched       State = "dispatched"
	OnTrip           State = "on_trip"
	Offline          State = "offline"

	//vehicle state
	Inactive State = "inactive"
	Active   State = "active"

	//customer state
	Searching State = "searching"
	AtRest    State = "at_rest"
)

type Delivery struct {
	Model
	State                 State         `bson:"state" json:"state"`
	OriginLongitude       float64       `bson:"originLongitude" json:"originLongitude"`
	OriginLatitude        float64       `bson:"originLatitude" json:"originLatitude"`
	DestinationLongitude  float64       `bson:"DestinationLongitude" json:"DestinationLongitude"`
	DestinationLatitude   float64       `bson:"DestinationLatitude" json:"DestinationLatitude"`
	InitialCost           float64       `bson:"initialCost" json:"initialCost"`
	FinalCost             float64       `bson:"finalCost" json:"finalCost"`
	Notes                 string        `bson:"notes" json:"notes"`
	Completed             bool          `bson:"completed" json:"completed"`
	CustomerId            bson.ObjectID `bson:"customerId" json:"customerId"`
	OrderId               bson.ObjectID `bson:"orderId" json:"orderId"`
	CourierId             bson.ObjectID `bson:"courierId" json:"courierId"`
	CourierRating         float64       `bson:"courierRating" json:"courierRating"`
	CourierRatingMessage  string        `bson:"courierRatingMessage" json:"courierRatingMessage"`
	CustomerRating        float64       `bson:"customerRating" json:"customerRating"`
	CustomerRatingMessage string        `bson:"customerRatingMessage" json:"customerRatingMessage"`
}
