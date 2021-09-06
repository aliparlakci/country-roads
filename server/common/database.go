package common

import (
	"context"
	"time"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitializeDb(uri, name, username, password string) (database *mongo.Database, close func()) {
	c, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	credentials := options.Credential{
		Username: username,
		Password: password,
	}
	client, err := mongo.Connect(c, options.Client().SetAuth(credentials).ApplyURI(uri))
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

func RedisInitilizer(uri, password string) func(db int) *redis.Client {
	return func(db int) *redis.Client {
		return redis.NewClient(&redis.Options{
			Addr:     uri,
			Password: password,
			DB:       db,
		})
	}
}
