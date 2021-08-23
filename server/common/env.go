package common

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Env struct {
	Db *mongo.Client
}
