package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Location struct {
	ID      primitive.ObjectID
	Display string
	Parent  primitive.ObjectID
}

func (location Location) String() string {
	return location.Display
}
