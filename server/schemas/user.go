package schemas

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID
	DisplayName string
	Email       string
	SignedUpAt  time.Time
}
