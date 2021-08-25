package common

import (
	"context"
	"time"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitilizeDb(uri string) (client *mongo.Client, close func()) {
	c, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(c, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	close = func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}

	return
}

func InitilizeRedis(uri, password string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     uri,
		Password: password,
		DB:       db,
	})
}
