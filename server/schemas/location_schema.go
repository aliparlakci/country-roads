package schemas

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LocationSchema struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Display  string             `bson:"display" json:"display"`
	ParentID primitive.ObjectID `bson:"parent,omitempty" json:"parent"`
}
