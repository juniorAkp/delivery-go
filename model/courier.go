package model

type Courier struct {
	Model
	Username  string   `bson:"username" json:"username" validate:"required,min=2,max=20"`
	Phone     string   `bson:"phone" json:"phone"`
	Vehicle   *Vehicle `bson:"vehicle,omitempty" json:"vehicle,omitempty"`
	Longitude float64  `bson:"longitude" json:"longitude"`
	Latitude  float64  `bson:"latitude" json:"latitude"`
	Email     string   `bson:"email" json:"email" validate:"required,email"`
	PhoneUrl  string   `bson:"photoUrl" json:"photoUrl"`
	Rating    float64  `bson:"rating" json:"rating"`
}
