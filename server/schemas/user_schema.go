package schemas

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSchema struct {
	ID          primitive.ObjectID
	DisplayName string
	Email       string
	SignedUpAt  time.Time
}
