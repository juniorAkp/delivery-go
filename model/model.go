package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Model struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
}
