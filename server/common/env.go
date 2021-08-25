package common

import (
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type Env struct {
	Db  *mongo.Client
	Rdb *redis.Client
}
