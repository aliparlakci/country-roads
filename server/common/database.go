package common

import (
	"context"
	"time"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitilizeDb(uri, name string) (database *mongo.Database, close func()) {
	c, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(c, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	database = client.Database(name)

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
