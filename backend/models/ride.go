package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ride struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type string             `bson:"type" json:"type"`
	Date time.Time          `bson:"date" json:"date"`
	From string             `bson:"from" json:"from"`
	To   string             `bson:"to" json:"to"`
}

type RideDTO struct {
	Type string `json:"type"`
	Date int64  `json:"date"`
	From string `json:"from"`
	To   string `json:"to"`
}
