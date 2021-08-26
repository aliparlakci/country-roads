package schemas

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RideSchema struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type        string             `bson:"type" json:"type"`
	Date        time.Time          `bson:"date" json:"date" time_format:"unix"`
	Direction   string             `bson:"direction" json:"direction"`
	Destination primitive.ObjectID `bson:"destination" json:"destination"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt,omitempty" time_format:"unix"`
	// Author      primitive.ObjectID `bson:"author" json:"author,omitempty"`
}
