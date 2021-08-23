package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ride struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type        string             `bson:"type" json:"type"`
	Date        time.Time          `bson:"date" json:"date" time_format:"unix"`
	Direction   string             `bson:"direction" json:"direction"`
	Destination string             `bson:"destination" json:"destination"`
	Author      primitive.ObjectID `bson:"author" json:"author,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt,omitempty" time_format:"unix"`
}

type RideDTO struct {
	Type        string    `bson:"type" json:"type" form:"type" binding:"required,validridetype"`
	Date        time.Time `bson:"date" json:"date" form:"date" time_format:"unix" binding:"required,futuredate"`
	Direction   string    `bson:"direction" json:"direction" form:"direction" binding:"required,validdirection"`
	Destination string    `bson:"destination" json:"destination" form:"destination" binding:"required"`
}
