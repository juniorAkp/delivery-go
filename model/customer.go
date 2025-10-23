package model

type Customer struct {
	Model
	Username     string  `bson:"username" json:"username" validate:"required,min=2,max=20"`
	Password     string  `bson:"password" json:"password" validate:"required,min=8"`
	Phone        string  `bson:"phone" json:"phone"`
	Code         int     `bson:"code" json:"code"`
	Longitude    float64 `bson:"longitude" json:"longitude"`
	Latitude     float64 `bson:"latitude" json:"latitude"`
	Email        string  `bson:"email" json:"email" validate:"required,email"`
	Token        string  `bson:"token" json:"token"`
	RefreshToken string  `bson:"refreshToken" json:"refreshToken"`
	Rating       float64 `bson:"rating" json:"rating"`
}

type CustomerLogin struct {
	Email    string `bson:"email" json:"email" validate:"required,email"`
	Password string `bson:"password" json:"password"`
}
